package importer

import (
	"testing"

	"github.com/DexterLB/mvm/library"
	"github.com/stretchr/testify/assert"
)

func TestImdbIdentifier(t *testing.T) {
	context := testContext(t)

	shows := make(chan library.ShowWithFile, 5)
	done := make(chan library.ShowWithFile, 5)
	doneSeries := make(chan *library.Series, 5)

	go context.ImdbIdentifier(shows, doneSeries, done)

	movie, err := context.Library.GetShowByImdbID(403358)
	if err != nil {
		t.Errorf("Library error: %s", err)
	}

	shows <- library.ShowWithFile{
		Show: movie,
		File: nil,
	}

	close(shows)

	doneMovie := <-done
	if _, ok := <-done; ok {
		t.Errorf("Done channel not closed after reading all shows")
	}

	if doneMovie.Show != movie {
		t.Errorf("Wrong movie is done")
	}

	if doneMovie.Show.ImdbError != nil {
		t.Fatalf("Imdb error: %s", *doneMovie.Show.ImdbError)
	}

	if _, ok := <-doneSeries; ok {
		t.Errorf("Series channel not closed after reading all series")
	}

	assert := assert.New(t)

	assert.Equal(403358, movie.ImdbID)
	assert.Equal("Nochnoy dozor", movie.Title)
	assert.InDelta(6.5, movie.ImdbRating, 0.01)
}

func TestImdbIdentifierMultipleShows(t *testing.T) {
	context := testContext(t)

	shows := make(chan library.ShowWithFile, 5)
	done := make(chan library.ShowWithFile, 5)
	doneSeries := make(chan *library.Series, 5)

	go context.ImdbIdentifier(shows, doneSeries, done)

	movie, err := context.Library.GetShowByImdbID(403358)
	if err != nil {
		t.Errorf("Library error: %s", err)
	}

	episode, err := context.Library.GetShowByImdbID(2816136)
	if err != nil {
		t.Errorf("Library error: %s", err)
	}

	shows <- library.ShowWithFile{
		Show: movie,
		File: nil,
	}

	shows <- library.ShowWithFile{
		Show: episode,
		File: nil,
	}

	close(shows)

	series := <-doneSeries

	if _, ok := <-doneSeries; ok {
		t.Errorf("Series channel not closed after reading all series")
	}

	if series == nil {
		t.Fatalf("Series is nil")
	}

	assert := assert.New(t)

	var (
		moviePresent   bool
		episodePresent bool
	)

	for sf := range done {
		show := sf.Show

		switch show.ImdbID {
		case 403358:
			assert.Equal("Nochnoy dozor", show.Title)
			assert.InDelta(6.5, show.ImdbRating, 0.01)
			moviePresent = true
		case 2816136:
			if len(series.Episodes) == 0 || series.Episodes[0] != show {
				t.Errorf("Series has wrong episode")
			}
			assert.Equal("Two Swords", show.Title)
			assert.Equal(4, show.Season)
			assert.Equal(1, show.Episode)
			episodePresent = true
		default:
			t.Errorf("Unknown show")
		}
	}

	assert.Equal(944947, series.ImdbID)
	assert.Equal("Game of Thrones", series.Title)

	if !episodePresent {
		t.Errorf("Episode not present in identifier output")
	}

	if !moviePresent {
		t.Errorf("Movie not present in identifier output")
	}
}
