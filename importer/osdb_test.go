// +build live

package importer

import (
	"testing"

	"github.com/DexterLB/mvm/library"
	"github.com/stretchr/testify/assert"
)

func TestOsdbIdentifier(t *testing.T) {
	context := testContext(t)

	files := make(chan *library.VideoFile, 5)
	done := make(chan *library.VideoFile, 5)
	shows := make(chan *library.Show, 5)

	go context.OsdbIdentifier(files, shows, done)

	file, err := context.Library.GetFileByPath("foo/bar")
	if err != nil {
		t.Errorf("Library error: %s", err)
	}

	file.OsdbHash = 0x09a2c497663259cb

	files <- file
	close(files)

	doneFile := <-done
	if _, ok := <-done; ok {
		t.Errorf("Done channel not closed after reading all files")
	}

	if doneFile != file {
		t.Errorf("Wrong file is done")
	}

	show := <-shows
	if _, ok := <-shows; ok {
		t.Errorf("Shows channel not closed after reading all shows")
	}

	assert := assert.New(t)

	assert.Equal("foo/bar", show.Files[0].Path)
	assert.Equal(403358, show.ImdbID)
}

func TestOsdbIdentifierMultipleShows(t *testing.T) {
	testOsdbIdentifierParallel(t, 3, 200)
}

func TestOsdbIdentifierMultipleRequests(t *testing.T) {
	testOsdbIdentifierParallel(t, 3, 1)
}

func testOsdbIdentifierParallel(t *testing.T, maxRequests int, maxPerRequest int) {
	context := testContext(t)

	context.Config.OsdbConfig.MaxRequests = maxRequests
	context.Config.OsdbConfig.MaxMoviesPerRequest = maxPerRequest

	files := make(chan *library.VideoFile, 5)
	done := make(chan *library.VideoFile, 5)
	shows := make(chan *library.Show, 5)

	go context.OsdbIdentifier(files, shows, done)

	file1, err := context.Library.GetFileByPath("foo/bar")
	if err != nil {
		t.Errorf("Library error: %s", err)
	}

	file1.OsdbHash = 0x09a2c497663259cb

	file2, err := context.Library.GetFileByPath("foo/baz")
	if err != nil {
		t.Errorf("Library error: %s", err)
	}

	file2.OsdbHash = 0x46e33be00464c12e

	files <- file1
	files <- file2
	close(files)

	_ = <-done
	_ = <-done

	if _, ok := <-done; ok {
		t.Errorf("Done channel not closed after reading all files")
	}

	assert := assert.New(t)
	for show := range shows {
		switch show.Files[0].Path {
		case "foo/bar":
			assert.Equal(403358, show.ImdbID)
		case "foo/baz":
			assert.Equal(2816136, show.ImdbID)
		default:
			t.Errorf("Unknown file")
		}
	}
}
