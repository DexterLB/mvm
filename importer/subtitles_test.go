package importer

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/DexterLB/mvm/library"
	"github.com/stretchr/testify/assert"
)

func TestSubtitleDownloader(t *testing.T) {
	context := testContext(t)

	files := make(chan library.ShowWithFile, 5)
	done := make(chan library.ShowWithFile, 5)
	subtitles := make(chan *library.Subtitle, 10)

	go context.SubtitleDownloader(files, subtitles, done)

	movie, err := context.Library.GetShowByImdbID(76759)
	if err != nil {
		t.Errorf("Library error: %s", err)
	}
	movie.Title = "Star Wars: Episode IV - A New Hope"
	movie.Year = 1977

	tempdir, err := ioutil.TempDir("", "mvm_test")
	if err != nil {
		t.Fatalf("can't create temp dir: %s", err)
	}
	defer os.RemoveAll(tempdir)

	file, err := context.Library.GetFileByPath(
		tempdir + "/Star.Wars.Episode.4.A.New.Hope.1977.1080p.BrRip.x264.BOKUTOX.YIFY.mp4",
	)
	if err != nil {
		t.Errorf("Library error: %s", err)
	}
	file.OsdbHash = 12185829445453571797
	file.Size = 1826970305

	movie.Files = []*library.VideoFile{file}

	files <- library.ShowWithFile{Show: movie, File: file}
	close(files)

	doneFile := <-done
	if _, ok := <-done; ok {
		t.Errorf("Done channel not closed after reading all files")
	}

	if doneFile.File != file {
		t.Errorf("Wrong file is done")
	}

	if doneFile.Show != movie {
		t.Errorf("Wrong movie is done")
	}

	if doneFile.File.SubtitlesError != nil {
		t.Fatalf("Subtitles error: %s", *doneFile.File.SubtitlesError)
	}

	var allSubtitles []*library.Subtitle
	for s := range subtitles {
		allSubtitles = append(allSubtitles, s)
	}

	if len(allSubtitles) != 4 {
		t.Fatalf("Downloaded %d subtitles instead of 4", len(allSubtitles))
	}

	assert := assert.New(t)
	assert.Equal("en", allSubtitles[0].Language.String())
	assert.Equal("322f10e7fee92c86ff46ce17cfbec64b", fmt.Sprintf("%s", allSubtitles[0].Hash))
	assert.Equal(false, allSubtitles[0].HearingImpaired)
	assert.Equal(99987157, allSubtitles[0].Score)
	assert.Equal(tempdir+"/Sintel.2010.720p.en.5.srt", allSubtitles[0].Filename)
	// assert.Equal("foobarbaz", md5File(t, allSubtitles[0].Filename))

	assert.Equal("en", allSubtitles[1].Language.String())
	assert.Equal("barfaz", fmt.Sprintf("%s", allSubtitles[1].Hash))
	assert.Equal(true, allSubtitles[1].HearingImpaired)
	assert.Equal(7, allSubtitles[1].Score)
	assert.Equal(tempdir+"/Sintel.2010.720p.en.7.srt", allSubtitles[1].Filename)
	assert.Equal("barfaz", md5File(t, allSubtitles[1].Filename))

	assert.Equal("de", allSubtitles[2].Language.String())
	assert.Equal("dfgfda", fmt.Sprintf("%s", allSubtitles[2].Hash))
	assert.Equal(false, allSubtitles[2].HearingImpaired)
	assert.Equal(3, allSubtitles[2].Score)
	assert.Equal(tempdir+"/Sintel.2010.720p.de.3.srt", allSubtitles[2].Filename)
	assert.Equal("dfgfda", md5File(t, allSubtitles[2].Filename))

	assert.Equal("de", allSubtitles[3].Language.String())
	assert.Equal("dfgfda", fmt.Sprintf("%s", allSubtitles[3].Hash))
	assert.Equal(false, allSubtitles[3].HearingImpaired)
	assert.Equal(9, allSubtitles[3].Score)
	assert.Equal(tempdir+"/Sintel.2010.720p.de.9.srt", allSubtitles[3].Filename)
	assert.Equal("dfgfda", md5File(t, allSubtitles[3].Filename))
}
