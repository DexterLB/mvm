package config

import (
	"testing"

	_ "github.com/orchestrate-io/dvr"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	config, err := Load("./testdata/test_config.toml")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	assert := assert.New(t)
	assert.Equal("/foo/bar", config.FileRoot)

	assert.Equal(50, config.Importer.BufferSize)

	assert.Equal(3, config.Importer.Osdb.MaxRequests)
	assert.Equal(199, config.Importer.Osdb.MaxMoviesPerRequest)
	assert.Equal("foo", config.Importer.Osdb.Username)
	assert.Equal("bar", config.Importer.Osdb.Password)

	assert.Equal("foosql", config.Library.Database)
	assert.Equal("bar", config.Library.DatabaseDSN)

	assert.Equal(16, config.Importer.Imdb.MaxRequests)
}
