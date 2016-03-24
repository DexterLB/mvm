package importer

import (
	"testing"

	"github.com/DexterLB/mvm/library"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/orchestrate-io/dvr"
)

func testContext(t *testing.T) *Context {
	lib, err := library.New("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("unable to initialize library: %s\n", err)
	}

	context := NewContext(lib, &Config{
		FileRoot: "./testdata",
	})

	go func() {
		for err := range context.Errors {
			t.Fatalf("context error: %s\n", err)
		}
	}()

	return context
}
