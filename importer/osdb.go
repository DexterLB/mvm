package importer

import (
	"strconv"
	"sync"

	"github.com/DexterLB/mvm/library"
	"github.com/oz/osdb"
)

type OsdbConfig struct {
	// Username for opensubtitles.org (leave blank for no user)
	Username string `toml:"username"`
	// Password for opensubtitles.org (leave blank for no password)
	Password string `toml:"password"`
	// MaxRequests is the maximum number of parallel requests to opensubtitles.org
	MaxRequests int `toml:"max_requests"`
	// MaxMoviesPerRequest is the maximum number of movies to ask for in a
	// single request. Currently, opensubtitles limits this to 200
	MaxMoviesPerRequest int `toml:"max_per_request"`
}

func (c *Context) OsdbIdentifier(
	files <-chan *library.VideoFile, shows chan<- *library.Show,
	done chan<- *library.VideoFile,
) {
	defer close(done)
	defer close(shows)

	config := &c.Config.OsdbConfig

	client, err := osdb.NewClient()
	if err != nil {
		c.Errorf("Can't initialize osdb client: %s", err)
		return
	}
	err = client.LogIn(config.Username, config.Password, "")
	if err != nil {
		c.Errorf("Can't login to osdb: %s", err)
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
	maxFiles := c.Config.OsdbConfig.MaxMoviesPerRequest

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
		hashes[i] = files[i].OsdbHash
	}
	movies, err := client.BestMoviesByHashes(hashes)
	if err != nil {
		for i := range files {
			files[i].Status.For("osdb_identify").Errorf(
				"Opensubtitles.org error: %s", err,
			)
			done <- files[i]
		}
		return
	}

	for i := range movies {
		id, err := strconv.Atoi(movies[i].Id)
		if err != nil {
			files[i].Status.For("osdb_identify").Errorf(
				"Can't parse show's imdb ID: %s", err,
			)
		} else {
			show, err := c.Library.GetShowByImdbID(id)
			if err != nil {
				files[i].Status.For("osdb_identify").Errorf(
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
