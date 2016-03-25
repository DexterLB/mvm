package importer

import (
	"sync"

	"github.com/DexterLB/mvm/imdb"
	"github.com/DexterLB/mvm/library"
)

type ImdbConfig struct {
	// MaxRequests is the maximum number of parallel requests to imdb
	MaxRequests int `toml:"max_requests"`
}

func (c *Context) ImdbIdentifier(
	shows <-chan *library.Show,
	doneSeries chan<- *library.Series,
	done chan<- *library.Show,
) {
	defer close(done)
	defer close(doneSeries)

	cache := makeSeriesCache()

	maxRequests := c.Config.ImdbConfig.MaxRequests

	wg := sync.WaitGroup{}
	wg.Add(maxRequests)
	for i := 0; i < maxRequests; i++ {
		go func() {
			defer wg.Done()

			for {
				select {
				case show, ok := <-shows:
					if !ok {
						return
					}

					seriesData := c.imdbProcessShow(show)

					var (
						series    *library.Series
						newSeries bool
						err       error
					)

					if seriesData != nil {
						id := seriesData.ID()
						cache.Lock()
						if series, ok = cache.PrevSeries[id]; !ok {
							series, err = c.Library.GetSeriesByImdbID(id)
							if err != nil {
								show.Status.For("imdb_identify").Errorf(
									"Unable to get series from library: %s", err,
								)
								cache.Unlock()
								continue
							}
							cache.PrevSeries[seriesData.ID()] = series
							newSeries = true
						}
						cache.Unlock()

						if newSeries {
							c.imdbProcessSeries(series, seriesData)
						}

						series.Lock()
						series.Episodes = append(series.Episodes, show)
						series.Unlock()
					}

					done <- show
					if newSeries {
						doneSeries <- series
					}
				case <-c.Stop:
					return
				}
			}
		}()
	}
	wg.Wait()
}

func (c *Context) imdbProcessShow(show *library.Show) *imdb.Item {
	connection := imdb.New(show.ImdbID)
	defer connection.Free()

	data, err := connection.AllData()
	if err != nil {
		show.Status.For("imdb_identify").Errorf(
			"Error getting data from imdb: %s", err,
		)
		return nil
	}

	show.Title = data.Title
	show.Year = data.Year
	show.OtherTitles = library.MapStringString(data.OtherTitles)
	show.Duration = library.Duration(data.Duration)
	show.Plot = data.Plot
	show.PlotMedium = data.PlotMedium
	show.PlotLong = data.PlotLong
	show.PosterURL = data.PosterURL
	show.ImdbRating = data.Rating
	show.ImdbVotes = data.Votes
	show.Languages = library.NewLanguages(data.Languages)
	show.ReleaseDate = data.ReleaseDate
	show.Tagline = data.Tagline

	show.Status.For("imdb_identify").Succeed()

	if data.Type == imdb.Episode {
		show.Season = data.SeasonNumber
		show.Episode = data.EpisodeNumber
		return data.Series
	}

	return nil
}

func (c *Context) imdbProcessSeries(series *library.Series, data *imdb.Item) {

}

type seriesCache struct {
	sync.Mutex

	PrevSeries map[int]*library.Series
}

func makeSeriesCache() *seriesCache {
	return &seriesCache{
		PrevSeries: make(map[int]*library.Series),
	}
}
