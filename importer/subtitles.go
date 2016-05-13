package importer

import (
	"fmt"
	"sync"

	"github.com/DexterLB/mvm/library"
	"github.com/DexterLB/mvm/types"
	"github.com/oz/osdb"
)

func (c *Context) SubtitleDownloader(
	files <-chan library.ShowWithFile,
	subtitles chan<- *library.Subtitle,
	done chan<- library.ShowWithFile,
) {
	defer close(done)
	defer close(subtitles)

	maxRequests := c.Config.Importer.Osdb.MaxRequests

	undownloaded := make(chan *undownloadedSubtitle, c.Config.Importer.BufferSize)
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
	undownloaded <-chan *undownloadedSubtitle,
	subtitles chan<- *library.Subtitle,
	done chan<- library.ShowWithFile,
	undownloadedCounts *subtitleCounts,
) {
	var currentSubtitles []*undownloadedSubtitle
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
	undownloaded chan<- *undownloadedSubtitle,
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
				file.File.SubtitlesError = types.Errorf(
					"unable to search for subtitles: %s", err,
				)
			}

			if len(subtitles) == 0 && undownloadedCounts.Done(file.File.ID) {
				filesWithNoSubtitles <- file
			}

			for i := range subtitles {
				undownloadedCounts.Push(file.File.ID)
				undownloaded <- &undownloadedSubtitle{
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
	undownloaded []*undownloadedSubtitle,
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
		subtitle, err := c.saveSubtitle(&data[i], undownloaded[i].File)
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
	file *library.VideoFile,
) (
	*library.Subtitle,
	error,
) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Context) searchForSubtitles(
	pair library.ShowWithFile,
) (
	[]*osdb.Subtitle,
	error,
) {
	// need to loop over all needed languages and search in parallel
	return nil, fmt.Errorf("not implemented")
}

func (c *Context) searchForSubtitlesWithLanguage(
	pair library.ShowWithFile,
	language types.Language,
) (
	[]*osdb.Subtitle,
	error,
) {
	client, err := c.OsdbClient()
	if err != nil {
		return nil, err
	}

	file := pair.File
	show := pair.Show

	// FIXME: this is retarded. Modify the osdb library to contain a params struct
	params := []interface{}{
		client.Token,
		[]struct {
			Hash      string `xmlrpc:"moviehash"`
			Size      uint64 `xmlrpc:"moviebytesize"`
			Filename  string `xmlrpc:"tag"`
			ImdbID    int    `xmlrpc:"imdbid"`
			Languages string `xmlrpc:"sublanguageid"`
			Title     string `xmlrpc:"query"`
			Season    string `xmlrpc:"season"`
			Episode   string `xmlrpc:"episode"`
		}{{
			Hash:      fmt.Sprintf("%016x", file.OsdbHash),
			Size:      file.Size,
			Filename:  file.OriginalBasename,
			ImdbID:    show.ImdbID,
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
		return nil, err
	}

	subtitles := make([]*osdb.Subtitle, len(monolithicSubtitles))
	for i := range monolithicSubtitles {
		subtitles[i] = &monolithicSubtitles[i]
	}

	return subtitles, nil
}

type undownloadedSubtitle struct {
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
