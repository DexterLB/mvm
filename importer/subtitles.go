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

			subtitles, err := c.searchSubtitlesForFile(file)
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

func (c *Context) searchSubtitlesForFile(
	file library.ShowWithFile,
) (
	[]*osdb.Subtitle,
	error,
) {
	return nil, fmt.Errorf("not implemented")
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
