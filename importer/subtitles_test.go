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
	for _, subtitle := range allSubtitles {
		switch subtitle.Hash {
		case "322f10e7fee92c86ff46ce17cfbec64b":
			assert.Equal("en", subtitle.Language.String())
			assert.Equal(false, subtitle.HearingImpaired)
			assert.Equal(9997, subtitle.Score/10000)
			assert.Equal(
				fmt.Sprintf(
					"%s/Star.Wars.Episode.4.A.New.Hope.1977.1080p.BrRip.x264.BOKUTOX.YIFY.en.%d.srt",
					tempdir,
					subtitle.Score,
				),
				subtitle.Filename,
			)
		case "5097c0bb857cfcf6b8bc4769da8513f6":
			assert.Equal("en", subtitle.Language.String())
			assert.Equal(false, subtitle.HearingImpaired)
			assert.Equal(9996, subtitle.Score/10000)
			assert.Equal(
				fmt.Sprintf(
					"%s/Star.Wars.Episode.4.A.New.Hope.1977.1080p.BrRip.x264.BOKUTOX.YIFY.en.%d.srt",
					tempdir,
					subtitle.Score,
				),
				subtitle.Filename,
			)
		case "c06e9a01f792b55b79581cead2067fe9":
			assert.Equal("bg", subtitle.Language.String())
			assert.Equal(false, subtitle.HearingImpaired)
			assert.Equal(9999, subtitle.Score/10000)
			assert.Equal(
				fmt.Sprintf(
					"%s/Star.Wars.Episode.4.A.New.Hope.1977.1080p.BrRip.x264.BOKUTOX.YIFY.bg.%d.srt",
					tempdir,
					subtitle.Score,
				),
				subtitle.Filename,
			)
		case "6e393aaaa4d33cff81cef1c9f795110c":
			assert.Equal("bg", subtitle.Language.String())
			assert.Equal(false, subtitle.HearingImpaired)
			assert.Equal(9999, subtitle.Score/10000)
			assert.Equal(
				fmt.Sprintf(
					"%s/Star.Wars.Episode.4.A.New.Hope.1977.1080p.BrRip.x264.BOKUTOX.YIFY.bg.%d.srt",
					tempdir,
					subtitle.Score,
				),
				subtitle.Filename,
			)

		default:
			t.Errorf("unknown subtitle hash: %s", subtitle.Hash)
		}
	}
}
