package main

import (
	"log"

	"github.com/DexterLB/mvm/config"
	"github.com/DexterLB/mvm/importer"
	"github.com/DexterLB/mvm/library"
	"github.com/cep21/xdgbasedir"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/voxelbrain/goptions"
)

type Options struct {
	ConfigFile string        `goptions:"-c, --config-file, description='path to the configuration file'"`
	Help       goptions.Help `goptions:"-h, --help, description='show this help message'"`

	goptions.Verbs
	Import struct {
		Path string `goptions:"-p, --path, description='path to import (single video file or a folder)', obligatory"`
	} `goptions:"import"`
}

func parseOptions() *Options {
	options := &Options{}
	goptions.ParseAndFail(options)
	return options
}

func parseConfig(options *Options) *config.Config {
	var err error

	if options.ConfigFile == "" {
		options.ConfigFile, err = xdgbasedir.GetConfigFileLocation("mvm.toml")
		if err != nil {
			log.Fatalf("can't find config file: %s", err)
		}
	}

	config, err := config.Load(options.ConfigFile)
	if err != nil {
		log.Fatalf("can't load config file: %s", err)
	}

	return config
}

func openLibrary(config *config.Config) *library.Library {
	library, err := library.New(
		config.Library.Database, config.Library.DatabaseDSN,
	)

	if err != nil {
		log.Fatalf("unable to initialize library database: %s", err)
	}
	return library
}

func main() {
	options := parseOptions()
	config := parseConfig(options)
	library := openLibrary(config)

	importer := importer.NewContext(library, config)
	go func() {
		for err := range importer.Errors {
			log.Printf("import error: %s", err)
		}
	}()

	importer.Import([]string{options.Import.Path})
}
