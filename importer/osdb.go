package importer

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/DexterLB/mvm/library"
	"github.com/DexterLB/mvm/types"
	"github.com/DexterLB/osdb"
)

// OsdbClient returns a logged in Osdb client
func (c *Context) OsdbClient() (*osdb.Client, error) {
	c.osdbLock.Lock()
	defer c.osdbLock.Unlock()

	if c.osdbClient != nil {
		return c.osdbClient, nil
	}

	config := &c.Config.Importer.Osdb

	client, err := osdb.NewClient()
	if err != nil {
		return nil, fmt.Errorf("Can't initialize osdb client: %s", err)
	}
	err = client.LogIn(config.Username, config.Password, "")
	if err != nil {
		return nil, fmt.Errorf("Can't login to osdb: %s", err)
	}

	c.osdbClient = client
	return client, nil
}

// OsdbIdentifier identifies the video files (matching them to shows) using
// the opensubtitles.org database
func (c *Context) OsdbIdentifier(
	files <-chan *library.VideoFile, shows chan<- *library.Show,
	done chan<- *library.VideoFile,
) {
	defer close(done)
	defer close(shows)

	config := &c.Config.Importer.Osdb

	client, err := c.OsdbClient()
	if err != nil {
		c.Errorf("%s", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(config.MaxRequests)
	for i := 0; i < config.MaxRequests; i++ {
		go func() {
			c.osdbIdentifierWorker(files, shows, done, client)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (c *Context) osdbIdentifierWorker(
	files <-chan *library.VideoFile, shows chan<- *library.Show,
	done chan<- *library.VideoFile,
	client *osdb.Client,
) {
	var currentFiles []*library.VideoFile
	maxFiles := c.Config.Importer.Osdb.MaxMoviesPerRequest

	for {
		select {
		case file, ok := <-files:
			if !ok {
				c.osdbProcessFiles(currentFiles, shows, done, client)
				return
			}
			currentFiles = append(currentFiles, file)
			if len(currentFiles) >= maxFiles {
				c.osdbProcessFiles(currentFiles, shows, done, client)
				currentFiles = currentFiles[0:0]
			}
		case <-c.Stop:
			return
		}
	}
}

func (c *Context) osdbProcessFiles(
	files []*library.VideoFile, shows chan<- *library.Show,
	done chan<- *library.VideoFile,
	client *osdb.Client,
) {
	if len(files) == 0 {
		return
	}

	hashes := make([]uint64, len(files))
	for i := range files {
		hashes[i] = uint64(files[i].OsdbHash)
	}
	movies, err := client.BestMoviesByHashes(hashes)
	if err != nil {
		for i := range files {
			files[i].OsdbError = types.Errorf(
				"Opensubtitles.org error: %s", err,
			)
			done <- files[i]
		}
		return
	}

	for i := range movies {
		var (
			err error
			id  int
		)
		if movies[i] == nil {
			err = fmt.Errorf("show not found in opensubtitles.org database")
		} else {
			id, err = strconv.Atoi(movies[i].ID)
			if err != nil {
				err = fmt.Errorf("can't parse imdb id: %s", err)
			}
		}

		if err != nil {
			files[i].OsdbError = types.Errorf(
				"can't identify show: %s", err,
			)
		} else {
			show, err := c.Library.GetShowByImdbID(id)
			if err != nil {
				files[i].OsdbError = types.Errorf(
					"Can't find show's imdb ID: %s", err,
				)
			} else {
				// TODO: episode data
				show.Files = append(show.Files, files[i])
				show.Title = movies[i].Title
				show.Year, _ = strconv.Atoi(movies[i].Year) // FIXME: check error
				shows <- show
			}
		}
		done <- files[i]
	}
}
