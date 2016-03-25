package importer

import (
	"sync"

	"github.com/DexterLB/mvm/library"
	"github.com/eapache/channels"
)

func (c *Context) Import(paths []string) {
	bufSize := c.Config.BufferSize

	filenames := make(chan string, bufSize)
	go c.WalkPaths(paths, filenames)

	files := make(chan *library.VideoFile, bufSize)
	go c.FileInfo(filenames, files)

	shows := make(chan *library.Show, bufSize)
	identifiedFiles := make(chan *library.VideoFile, bufSize)
	go c.OsdbIdentifier(files, shows, identifiedFiles)

	identifiedSeries := make(chan *library.Series, bufSize)
	identifiedShows := make(chan *library.Show, bufSize)
	go c.ImdbIdentifier(shows, identifiedSeries, identifiedShows)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		c.saveAll(identifiedFiles)
		wg.Done()
	}()
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
