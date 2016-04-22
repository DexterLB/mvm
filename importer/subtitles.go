package importer

import (
	"sync"

	"github.com/DexterLB/mvm/library"
	"github.com/DexterLB/mvm/types"
	"github.com/oz/osdb"
)

type undownloadedSubtitle struct {
	subtitle *osdb.Subtitle
	forFile  *library.VideoFile
}

func (c *Context) SubtitleDownloader(
	files <-chan *library.VideoFile,
	subtitles chan<- *library.Subtitle,
	done chan<- *library.VideoFile,
) {
	defer close(done)
	defer close(subtitles)

	maxRequests := c.Config.Importer.Osdb.MaxRequests

	undownloaded := make(chan *undownloadedSubtitle, c.Config.Importer.BufferSize)

	go func() {
		defer close(undownloaded)

		wg := sync.WaitGroup{}
		wg.Add(maxRequests)
		for i := 0; i < maxRequests; i++ {
			go func() {
				defer wg.Done()
				c.subtitleSearcherWorker(files, undownloaded)
			}()
		}
		wg.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(maxRequests)
	for i := 0; i < maxRequests; i++ {
		go func() {
			defer wg.Done()
			c.subtitleDownloaderWorker(undownloaded, subtitles, done)
		}()
	}
	wg.Wait()
}

func (c *Context) subtitleDownloaderWorker(
	undownloaded <-chan *undownloadedSubtitle,
	subtitles chan<- *library.Subtitle,
	done chan<- *library.VideoFile,
) {
	for {
		select {
		case us, ok := <-undownloaded:
			if !ok {
				return
			}

			done <- us.forFile
		case <-c.Stop:
			return
		}
	}
}

func (c *Context) subtitleSearcherWorker(
	files <-chan *library.VideoFile,
	undownloaded chan<- *undownloadedSubtitle,
) {
	for {
		select {
		case file, ok := <-files:
			if !ok {
				return
			}

			subtitles, err := c.searchSubtitlesForFile(file)
			if err != nil {
				file.SubtitlesError = types.Errorf(
					"unable to search for subtitles: %s", err,
				)
			}

			for i := range subtitles {
				undownloaded <- &undownloadedSubtitle{
					forFile:  file,
					subtitle: subtitles[i],
				}
			}
		case <-c.Stop:
			return
		}
	}
}
