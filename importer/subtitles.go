package importer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/DexterLB/mvm/library"
	"github.com/DexterLB/mvm/types"
	"github.com/DexterLB/osdb"
)

// SubtitleDownloader downloads subtitles for each file, using information
// from its associated show.
func (c *Context) SubtitleDownloader(
	files <-chan library.ShowWithFile,
	subtitles chan<- *library.Subtitle,
	done chan<- library.ShowWithFile,
) {
	defer close(done)
	defer close(subtitles)

	maxRequests := c.Config.Importer.Osdb.MaxRequests

	undownloaded := make(chan *subtitleInfo, c.Config.Importer.BufferSize)
	undownloadedCounts := newSubtitleCounts()

	go func() {
		defer close(undownloaded)

		wg := sync.WaitGroup{}
		wg.Add(maxRequests)
		for i := 0; i < maxRequests; i++ {
			go func() {
				defer wg.Done()
				c.subtitleSearcherWorker(files, undownloaded, done, undownloadedCounts)
			}()
		}
		wg.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(maxRequests)
	for i := 0; i < maxRequests; i++ {
		go func() {
			defer wg.Done()
			c.subtitleDownloaderWorker(undownloaded, subtitles, done, undownloadedCounts)
		}()
	}
	wg.Wait()
}

func (c *Context) subtitleDownloaderWorker(
	undownloaded <-chan *subtitleInfo,
	subtitles chan<- *library.Subtitle,
	done chan<- library.ShowWithFile,
	undownloadedCounts *subtitleCounts,
) {
	var currentSubtitles []*subtitleInfo
	maxSubtitles := c.Config.Importer.Osdb.MaxSubtitlesPerRequest
	for {
		select {
		case us, ok := <-undownloaded:
			if !ok {
				c.downloadSubtitles(currentSubtitles, subtitles, done, undownloadedCounts)
				return
			}

			currentSubtitles = append(currentSubtitles, us)
			if len(currentSubtitles) >= maxSubtitles {
				c.downloadSubtitles(currentSubtitles, subtitles, done, undownloadedCounts)
				currentSubtitles = currentSubtitles[0:0]
			}
		case <-c.Stop:
			return
		}
	}
}

func (c *Context) subtitleSearcherWorker(
	files <-chan library.ShowWithFile,
	undownloaded chan<- *subtitleInfo,
	filesWithNoSubtitles chan<- library.ShowWithFile,
	undownloadedCounts *subtitleCounts,
) {
	for {
		select {
		case file, ok := <-files:
			if !ok {
				return
			}

			subtitles, err := c.searchForSubtitles(file)
			if err != nil {
				file.File.SubtitlesError = types.Errorf("%s", err)
			}

			if len(subtitles) == 0 && undownloadedCounts.Done(file.File.ID) {
				filesWithNoSubtitles <- file
			}

			for i := range subtitles {
				undownloadedCounts.Push(file.File.ID)
				undownloaded <- &subtitleInfo{
					ShowWithFile: file,
					Subtitle:     subtitles[i],
				}
			}
		case <-c.Stop:
			return
		}
	}
}

func (c *Context) downloadSubtitles(
	undownloaded []*subtitleInfo,
	subtitles chan<- *library.Subtitle,
	done chan<- library.ShowWithFile,
	undownloadedCounts *subtitleCounts,
) {
	if len(undownloaded) == 0 {
		return
	}

	defer func() {
		for i := range undownloaded {
			undownloadedCounts.Pop(undownloaded[i].File.ID)
			if undownloadedCounts.Done(undownloaded[i].File.ID) {
				done <- undownloaded[i].ShowWithFile
			}
		}
	}()

	toDownload := make(osdb.Subtitles, len(undownloaded))
	for i := range undownloaded {
		toDownload[i] = *undownloaded[i].Subtitle
	}

	var data []osdb.SubtitleFile

	client, err := c.OsdbClient()
	if err == nil {
		data, err = client.DownloadSubtitles(toDownload)
	}

	if err != nil {
		for i := range undownloaded {
			undownloaded[i].File.SubtitlesError = types.Errorf(
				"unable to download subtitles: %s", // FIXME: what if there's already another error?
				err,
			)
		}
		return
	}

	for i := range data {
		subtitle, err := c.saveSubtitle(&data[i], undownloaded[i])
		if err != nil {
			undownloaded[i].File.SubtitlesError = types.Errorf(
				"unable to save subtitles: %s",
				err,
			)
		} else {
			subtitles <- subtitle
		}
	}
}

func (c *Context) saveSubtitle(
	data *osdb.SubtitleFile,
	info *subtitleInfo,
) (
	*library.Subtitle,
	error,
) {
	language, err := types.ParseLanguage(info.Subtitle.ISO639)
	if err != nil { // TODO: maybe check if the returned language is correct
		return nil, fmt.Errorf("unknown subtitle language")
	}

	score, err := strconv.Atoi(info.Subtitle.SubDownloadsCnt)
	if err != nil {
		return nil, fmt.Errorf("cannot determine subtitle score: %s", err)
	}
	score = 99999999 - score

	description := &struct {
		NoExtPath string
		Language  string
		Score     string
		Format    string
	}{
		NoExtPath: strings.TrimSuffix(info.File.Path, filepath.Ext(info.File.Path)),
		Language:  language.ISO2(),
		Score:     fmt.Sprintf("%08d", score),
		Format:    info.Subtitle.SubFormat,
	}

	filename, err := c.Config.Importer.Subtitles.Filename.On(description)
	if err != nil {
		return nil, fmt.Errorf("unable to determine subtitle filename: %s", err)
	}

	reader, err := data.Reader()
	if err != nil {
		return nil, fmt.Errorf("unable to read subtitle data: %s", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open subtitle file for writing: %s", err)
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = io.Copy(f, reader)
	if err != nil {
		return nil, fmt.Errorf("unble to write subtitle data: %s", err)
	}

	subtitle, err := c.Library.GetSubtitleByFilename(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to create subtitle in library: %s", err)
	}

	subtitle.Hash = info.Subtitle.SubHash
	subtitle.Language = language
	subtitle.HearingImpaired = (info.Subtitle.SubHearingImpaired == "true")
	subtitle.Score = score

	info.File.Lock()
	info.File.Subtitles = append(info.File.Subtitles, subtitle)
	info.File.Unlock()

	return subtitle, nil
}

// searchForSubtitles searches for subtitles for all languages specified
// in the config, and returns the matched subtitle objects. It might
// return valid subtitles _and_ an error if some, but not all of the languages
// fail to execute.
func (c *Context) searchForSubtitles(
	pair library.ShowWithFile,
) (
	[]*osdb.Subtitle,
	error,
) {
	languages := c.Config.Importer.Subtitles.Languages
	// TODO: native languages

	wg := sync.WaitGroup{}
	wg.Add(len(languages))

	errors := make([]string, 0, len(languages))
	errorLock := sync.Mutex{}

	results := make(chan *osdb.Subtitle, c.Config.Importer.BufferSize)
	go func() {
		wg.Wait()
		close(results)
	}()

	// FIXME: this will launch too many requests for >1 languages.
	// it should also not rely on the assumption that the number of
	// languages is less than the number of allowed simoultaneous requests
	for i := range languages {
		go func(errors *[]string, i int) {
			defer wg.Done()
			err := c.searchForSubtitlesWithLanguage(pair, languages[i], results)
			if err != nil {
				errorLock.Lock()
				*errors = append(*errors, fmt.Sprintf("%s", err))
			}
		}(&errors, i)
	}

	var subtitles []*osdb.Subtitle

	for subtitle := range results {
		subtitles = append(subtitles, subtitle)
	}

	if len(errors) > 0 {
		return subtitles, fmt.Errorf(
			"errors while searching for subtitles: %s",
			strings.Join(errors, ", "),
		)
	}

	return subtitles, nil
}

func (c *Context) searchForSubtitlesWithLanguage(
	pair library.ShowWithFile,
	language types.Language,
	results chan<- *osdb.Subtitle,
) error {
	client, err := c.OsdbClient()
	if err != nil {
		return err
	}

	file := pair.File
	show := pair.Show

	// holy shit! the opensubtitles API is ugly!
	params := []interface{}{
		client.Token,
		[]interface{}{
			struct {
				Hash      string `xmlrpc:"moviehash"`
				Size      uint64 `xmlrpc:"moviebytesize"`
				Languages string `xmlrpc:"sublanguageid"`
			}{
				Hash:      fmt.Sprintf("%016x", file.OsdbHash),
				Size:      file.Size,
				Languages: language.ISO3(),
			},
			struct {
				Filename  string `xmlrpc:"tag"`
				Languages string `xmlrpc:"sublanguageid"`
			}{
				Filename:  file.OriginalBasename,
				Languages: language.ISO3(),
			},
			struct {
				ImdbID    int    `xmlrpc:"imdbid"`
				Languages string `xmlrpc:"sublanguageid"`
			}{
				ImdbID:    show.ImdbID,
				Languages: language.ISO3(),
			},
			struct {
				Title     string `xmlrpc:"query"`
				Season    string `xmlrpc:"season"`
				Episode   string `xmlrpc:"episode"`
				Languages string `xmlrpc:"sublanguageid"`
			}{
				Title:     show.Title,
				Season:    fmt.Sprintf("%d", show.Season),
				Episode:   fmt.Sprintf("%d", show.Episode),
				Languages: language.ISO3(),
			}},
		struct {
			NumberOfResults int `xmlrpc:"limit"`
		}{
			NumberOfResults: c.Config.Importer.Subtitles.SubtitlesPerLanguage,
		},
	}

	// FIXME: this is even more retarded. Modify the osdb library.
	monolithicSubtitles, err := client.SearchSubtitles(&params)
	if err != nil {
		return err
	}

	for i := range monolithicSubtitles {
		results <- &monolithicSubtitles[i]
	}

	return nil
}

type subtitleInfo struct {
	library.ShowWithFile

	Subtitle *osdb.Subtitle
}

type subtitleCounts struct {
	sync.Mutex

	counts map[uint]uint
}

func newSubtitleCounts() *subtitleCounts {
	return &subtitleCounts{
		counts: make(map[uint]uint),
	}
}

func (s *subtitleCounts) Push(id uint) {
	s.Lock()
	defer s.Unlock()

	s.counts[id]++
}

func (s *subtitleCounts) Pop(id uint) {
	s.Lock()
	defer s.Unlock()

	s.counts[id]--
}

func (s *subtitleCounts) Done(id uint) bool {
	s.Lock()
	defer s.Unlock()

	return (s.counts[id] == 0)
}
