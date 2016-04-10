// Package config represents the configuration for mvm
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/DexterLB/mvm/types"
)

// Config contains the general configuration of mvm
type Config struct {
	FileRoot string   `toml:"file_root"`
	Importer Importer `toml:"importer"`
	Library  Library  `toml:"library"`
}

// Importer contains the configuration for all importers
type Importer struct {
	BufferSize int       `toml:"buffer_size"`
	Osdb       Osdb      `toml:"osdb"`
	Imdb       Imdb      `toml:"imdb"`
	Subtitles  Subtitles `toml:"subtitles"`
}

// Library contains the configuration for the library
type Library struct {
	Database    string `toml:"database"`
	DatabaseDSN string `toml:"database_dsn"`
}

// Osdb contains the configuration related to the opensubtitles.org api
type Osdb struct {
	// Username for opensubtitles.org (leave blank for no user)
	Username string `toml:"username"`
	// Password for opensubtitles.org (leave blank for no password)
	Password string `toml:"password"`
	// MaxRequests is the maximum number of parallel requests to opensubtitles.org
	MaxRequests int `toml:"max_requests"`
	// MaxMoviesPerRequest is the maximum number of movies to ask for in a
	// single request. Currently, opensubtitles limits this to 200
	MaxMoviesPerRequest int `toml:"max_per_request"`
}

// Imdb contains the configuration related to imdb.com
type Imdb struct {
	// MaxRequests is the maximum number of parallel requests to imdb
	MaxRequests int `toml:"max_requests"`
}

// Subtitles contains the configuration for the subtitle downloader
type Subtitles struct {
	Languages            types.Languages `toml:"languages"`
	Filename             types.Template  `toml:"filename"`
	SubtitlesPerLanguage int             `toml:"subtitles_per_language"`
}

// Load loads a configuration file
func Load(filename string) (*Config, error) {
	config := &Config{}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	md, err := toml.DecodeReader(f, &config)
	if err != nil {
		return nil, err
	}
	undecoded := md.Undecoded()

	if len(undecoded) > 0 {
		return nil, fmt.Errorf("unknown config values: %v", undecoded)
	}

	return config, nil
}
