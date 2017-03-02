package importer

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"testing"

	"github.com/DexterLB/mvm/config"
	"github.com/DexterLB/mvm/library"
	"github.com/DexterLB/mvm/types"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func testContext(t *testing.T) *Context {
	lib, err := library.New("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("unable to initialize library: %s\n", err)
	}

	context := NewContext(lib, &config.Config{
		FileRoot: "./fixtures",
		Importer: config.Importer{
			Osdb: config.Osdb{
				MaxMoviesPerRequest:    200,
				MaxSubtitlesPerRequest: 19,
				MaxRequests:            3,
			},
			Imdb: config.Imdb{
				MaxRequests: 8,
			},
			Subtitles: config.Subtitles{
				Languages: types.MustParseLanguages("en bg"),
				Filename: types.MustParseTemplate(
					"{{.NoExtPath}}.{{.Language}}.{{.Score}}.{{.Format}}",
				),
				SubtitlesPerLanguage: 2,
			},
		},
	})

	go func() {
		for err := range context.Errors {
			t.Fatalf("context error: %s\n", err)
		}
	}()

	return context
}

func md5File(t *testing.T, filename string) (hash string) {
	hasher := md5.New()
	f, err := os.Open(filename)
	if err != nil {
		t.Errorf("can't open file: %s", err)
		return
	}
	defer f.Close()
	_, err = io.Copy(hasher, f)
	if err != nil {
		t.Errorf("can't read from file: %s", err)
		return
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
