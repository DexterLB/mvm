package importer

import (
	"sync"

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

					series := c.imdbProcessShow(show)
					var newSeries bool

					if series != nil {
						series.Lock()
						newSeries = !cache.Set(series.ImdbID)
						if newSeries {
							c.imdbProcessSeries(series)
						}
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

func (c *Context) imdbProcessShow(show *library.Show) *library.Series {
	return nil
}

func (c *Context) imdbProcessSeries(series *library.Series) {

}

type seriesCache struct {
	sync.Mutex

	PrevIDs map[int]struct{}
}

func makeSeriesCache() *seriesCache {
	return &seriesCache{
		PrevIDs: make(map[int]struct{}),
	}
}

func (i *seriesCache) Set(id int) bool {
	i.Lock()
	defer i.Unlock()

	if _, ok := i.PrevIDs[id]; ok {
		i.PrevIDs[id] = struct{}{}
		return false
	}
	return true
}
