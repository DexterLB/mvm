package library

import (
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/orchestrate-io/dvr"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	plots := [...]string{
		`Lorem ipsum dolor sit amet`,
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit.`,
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
Ut est arcu, tempor quis accumsan quis, imperdiet ut ex.
Pellentesque vel lobortis est. Vivamus lobortis eleifend dapibus.
Nullam laoreet ipsum sed massa ornare tristique in nec lorem.
In eleifend odio ac accumsan ultrices. Aenean lacinia vel risus quis mattis.
Donec suscipit pretium euismod. 
Etiam sed justo venenatis, interdum tortor quis, aliquet ipsum.
Vestibulum a facilisis lectus.
Fusce aliquam lectus vel vehicula consequat. 
Aenean venenatis, velit rhoncus scelerisque dictum, 
lorem neque auctor velit, id pretium dui dolor ut ex. Sed quis augue.`,
	}

	languages := []string{"ru", "en"}

	lib, err := New("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	isin, err := lib.HasImdbID(999999)
	if err != nil {
		t.Fatal(err)
	}

	if isin {
		t.Fatalf("Movie already in the database?")
	}

	movie, err := lib.GetShowByImdbID(999999)
	if err != nil {
		t.Fatal(err)
	}

	movie.Title = "title"
	movie.Year = 2048
	movie.OtherTitles = map[string]string{
		"foo": "bar",
		"bar": "baz",
	}
	movie.Duration = Duration(3 * time.Minute)
	movie.Plot = plots[0]
	movie.PlotMedium = plots[1]
	movie.PlotLong = plots[2]
	movie.PosterURL = "http://example.com/foo.jpg"
	movie.ImdbRating = 3.14
	movie.ImdbVotes = 42
	movie.Languages = languages
	movie.ReleaseDate = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	movie.EpisodeData = nil
	movie.Tagline = "foo!"

	err = lib.Save(movie)
	if err != nil {
		t.Fatal(err)
	}

	isin, err = lib.HasImdbID(999999)
	if err != nil {
		t.Fatal(err)
	}

	if !isin {
		t.Fatalf("Movie not in the database?")
	}

	movie2, err := lib.GetShowByImdbID(999999)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(999999, movie2.ImdbID)
	assert.Equal("title", movie2.Title)
	assert.Equal(2048, movie2.Year)
	assert.Equal("bar", movie2.OtherTitles["foo"])
	assert.Equal("baz", movie2.OtherTitles["bar"])
	assert.Equal(2, len(movie2.OtherTitles))
	assert.Equal(3*time.Minute, time.Duration(movie2.Duration))
	assert.Equal(plots[0], movie2.Plot)
	assert.Equal(plots[1], movie2.PlotMedium)
	assert.Equal(plots[2], movie2.PlotLong)
	assert.Equal("http://example.com/foo.jpg", movie2.PosterURL)
	assert.InDelta(3.14, movie2.ImdbRating, 0.0001)
	assert.Equal(42, movie2.ImdbVotes)
	assert.Equal(languages, movie2.Languages)
	assert.Equal(
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		movie2.ReleaseDate,
	)

	assert.Nil(movie2.EpisodeData)
	assert.Equal("foo!", movie2.Tagline)
}
