package library

import (
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/orchestrate-io/dvr"
	"github.com/stretchr/testify/assert"
)

func TestShow(t *testing.T) {
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

	isin, err := lib.HasShowWithImdbID(999999)
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

	isin, err = lib.HasShowWithImdbID(999999)
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
	assert.Equal(
		map[string]string{"foo": "bar", "bar": "baz"},
		map[string]string(movie2.OtherTitles),
	)
	assert.Equal(3*time.Minute, time.Duration(movie2.Duration))
	assert.Equal(plots[0], movie2.Plot)
	assert.Equal(plots[1], movie2.PlotMedium)
	assert.Equal(plots[2], movie2.PlotLong)
	assert.Equal("http://example.com/foo.jpg", movie2.PosterURL)
	assert.InDelta(3.14, movie2.ImdbRating, 0.0001)
	assert.Equal(42, movie2.ImdbVotes)
	assert.Equal(languages, []string(movie2.Languages))
	assert.Equal(
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		movie2.ReleaseDate,
	)

	assert.Nil(movie2.EpisodeData)
	assert.Equal("foo!", movie2.Tagline)
}

func TestVideoFile(t *testing.T) {
	lib, err := New("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	isin, err := lib.HasFileWithPath("/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	if isin {
		t.Fatalf("File already in the database?")
	}

	file, err := lib.GetFileByPath("/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	file.Size = 98765432
	file.ResolutionX = 1920
	file.ResolutionY = 1080
	file.OsdbHash = 123456789
	file.Format = "h264"
	file.Duration = Duration(time.Minute * 20)

	file.LastPlayed = time.Date(2012, time.February, 10, 23, 15, 32, 5, time.UTC)
	file.LastPosition = Duration(time.Minute*12 + time.Second*38)

	err = lib.Save(file)
	if err != nil {
		t.Fatal(err)
	}

	isin, err = lib.HasFileWithPath("/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	if !isin {
		t.Fatalf("File not in the database?")
	}

	file2, err := lib.GetFileByPath("/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal("/foo/bar", file2.Path)
	assert.Equal(uint64(98765432), file2.Size)
	assert.Equal(uint(1920), file2.ResolutionX)
	assert.Equal(uint(1080), file2.ResolutionY)
	assert.Equal(uint64(123456789), file2.OsdbHash)
	assert.Equal("h264", file2.Format)
	assert.Equal(time.Minute*20, time.Duration(file2.Duration))
	assert.Equal(
		file2.LastPlayed,
		time.Date(2012, time.February, 10, 23, 15, 32, 5, time.UTC),
	)
	assert.Equal(
		time.Duration(file2.LastPosition), time.Minute*12+time.Second*38,
	)
}
