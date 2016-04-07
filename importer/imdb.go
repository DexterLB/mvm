package importer

import (
	"sync"

	"github.com/DexterLB/mvm/imdb"
	"github.com/DexterLB/mvm/library"
)

// ImdbIdentifier fetches data from imdb for the given shows (they must
// have an ImdbID). For shows which are episodes, it fetches the respective
// series.
func (c *Context) ImdbIdentifier(
	shows <-chan *library.Show,
	doneSeries chan<- *library.Series,
	done chan<- *library.Show,
) {
	defer close(done)
	defer close(doneSeries)

	cache := makeSeriesCache()

	maxRequests := c.Config.Importer.Imdb.MaxRequests

	wg := sync.WaitGroup{}
	wg.Add(maxRequests)
	for i := 0; i < maxRequests; i++ {
		go func() {
			defer wg.Done()
			c.imdbIdentifierWorker(shows, doneSeries, done, cache)
		}()
	}
	wg.Wait()
}

func (c *Context) imdbIdentifierWorker(
	shows <-chan *library.Show,
	doneSeries chan<- *library.Series,
	done chan<- *library.Show,
	cache *seriesCache,
) {
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
						show.ImdbError = library.Errorf(
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
}

func (c *Context) imdbProcessShow(show *library.Show) *imdb.Item {
	connection := imdb.New(show.ImdbID)
	defer connection.Free()

	data, err := connection.AllData()
	if err != nil {
		show.ImdbError = library.Errorf(
			"Error getting data from imdb: %s", err,
		)
		return nil
	}

	imdbSetCommonData(&show.CommonData, data)
	show.ReleaseDate = data.ReleaseDate
	show.Tagline = data.Tagline

	show.ImdbError = nil

	if data.Type == imdb.Episode {
		show.Season = data.SeasonNumber
		show.Episode = data.EpisodeNumber
		return data.Series
	}

	return nil
}

func (c *Context) imdbProcessSeries(series *library.Series, imdbSeries *imdb.Item) {
	defer imdbSeries.Free()

	data, err := imdbSeries.AllData()
	if err != nil {
		series.ImdbError = library.Errorf(
			"Error getting data from imdb: %s", err,
		)
		return
	}

	imdbSetCommonData(&series.CommonData, data)
}

func imdbSetCommonData(commonData *library.CommonData, data *imdb.ItemData) {
	commonData.Title = data.Title
	commonData.Year = data.Year
	commonData.OtherTitles = library.MapStringString(data.OtherTitles)
	commonData.Duration = library.Duration(data.Duration)
	commonData.Plot = data.Plot
	commonData.PlotMedium = data.PlotMedium
	commonData.PlotLong = data.PlotLong
	commonData.PosterURL = data.PosterURL
	commonData.ImdbRating = data.Rating
	commonData.ImdbVotes = data.Votes
	commonData.Languages = library.NewLanguages(data.Languages)
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
