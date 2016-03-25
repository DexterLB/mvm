// +build live

package importer

import (
	"testing"

	"github.com/DexterLB/mvm/library"
	"github.com/stretchr/testify/assert"
)

func TestImdbIdentifier(t *testing.T) {
	context := testContext(t)

	shows := make(chan *library.Show, 5)
	done := make(chan *library.Show, 5)
	doneSeries := make(chan *library.Series, 5)

	go context.ImdbIdentifier(shows, doneSeries, done)

	movie, err := context.Library.GetShowByImdbID(403358)
	if err != nil {
		t.Errorf("Library error: %s", err)
	}

	shows <- movie
	close(shows)

	doneMovie := <-done
	if _, ok := <-done; ok {
		t.Errorf("Done channel not closed after reading all shows")
	}

	if doneMovie != movie {
		t.Errorf("Wrong movie is done")
	}

	if _, ok := <-doneSeries; ok {
		t.Errorf("Series channel not closed after reading all series")
	}

	assert := assert.New(t)

	assert.Equal(403358, movie.ImdbID)
	assert.Equal("Nochnoy Dozor", movie.Title)
	assert.Equal(6.5, movie.ImdbRating)
}

func TestImdbIdentifierMultipleShows(t *testing.T) {
}
