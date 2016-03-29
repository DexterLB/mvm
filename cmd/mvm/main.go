package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DexterLB/mvm/config"
	"github.com/DexterLB/mvm/importer"
	"github.com/DexterLB/mvm/library"
	"github.com/cep21/xdgbasedir"
	"github.com/codegangsta/cli"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func parseConfig(c *cli.Context) *config.Config {
	var err error

	filename := c.GlobalString("config-file")

	if filename == "" {
		filename, err = xdgbasedir.GetConfigFileLocation("mvm.toml")
		if err != nil {
			log.Fatalf("can't find config file: %s", err)
		}
	}

	config, err := config.Load(filename)
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

func runImport(c *cli.Context) {
	if c.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "please supply a filename\n")
	}

	config := parseConfig(c)
	library := openLibrary(config)

	importer := importer.NewContext(library, config)
	go func() {
		for err := range importer.Errors {
			log.Printf("import error: %s", err)
		}
	}()

	importer.Import([]string(c.Args()))
}

func main() {
	app := cli.NewApp()
	app.Name = "mvm"
	app.Usage = "identify, manipulate and search movies and series"

	app.Commands = []cli.Command{
		{
			Name:    "import",
			Aliases: []string{"imp", "i"},
			Usage:   "import a video file into the library",
			Action:  runImport,
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-file, c",
			Value: "",
			Usage: "path to the mvm configuration file",
		},
	}

	app.RunAndExitOnError()
}
