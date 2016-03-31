package importer

import (
	"sync"

	"github.com/DexterLB/mvm/library"
	"github.com/eapache/channels"
)

func (c *Context) Import(paths []string) {
	bufSize := c.Config.Importer.BufferSize

	filenames := make(chan string, bufSize)
	go c.WalkPaths(paths, filenames)

	files := make(chan *library.VideoFile, bufSize)
	go c.FileInfo(filenames, files)

	shows := make(chan *library.Show, bufSize)
	identifiedFiles := make(chan *library.VideoFile, bufSize)
	go c.OsdbIdentifier(files, shows, identifiedFiles)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		c.saveAll(c.filterFilesWithErrors(identifiedFiles))
		wg.Done()
	}()
	go func() {
		c.ProcessShows(shows)
		wg.Done()
	}()
	wg.Wait()

	close(c.Stop)
}

func (c *Context) ProcessShows(shows <-chan *library.Show) {
	bufSize := c.Config.Importer.BufferSize

	identifiedSeries := make(chan *library.Series, bufSize)
	identifiedShows := make(chan *library.Show, bufSize)
	go c.ImdbIdentifier(shows, identifiedSeries, identifiedShows)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		c.saveAll(identifiedShows)
		wg.Done()
	}()
	go func() {
		c.saveAll(identifiedSeries)
		wg.Done()
	}()
	wg.Wait()
}

func (c *Context) saveAll(genericChannel interface{}) {
	channel := channels.Wrap(genericChannel).Out()
	for item := range channel {
		c.Library.Save(item)
	}
}

func (c *Context) filterFilesWithErrors(files <-chan *library.VideoFile) chan<- *library.VideoFile {
	out := make(chan *library.VideoFile)

	go func() {
		defer close(out)

		for file := range files {
			out <- file
			if file.ImportError != nil || file.OsdbError != nil {
				c.FilesWithErrors = append(c.FilesWithErrors, file)
			}
		}
	}()

	return out
}
