package main

import (
	"fmt"
	"log"

	"github.com/DexterLB/mvm/config"
	"github.com/cep21/xdgbasedir"
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

func main() {
	var err error
	options := &Options{}
	goptions.ParseAndFail(options)

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

	fmt.Printf("config: %v\n", config)
}
