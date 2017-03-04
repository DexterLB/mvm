package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

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

type userState int

const (
	success userState = iota
	retry
	abort
)

func manualImport(
	c *cli.Context,
	importer *importer.Context,
	shows chan<- *library.ShowWithFile,
	file *library.VideoFile,
) userState {
	input := bufio.NewScanner(os.Stdin)

	fmt.Printf(
		"What do you want to do? Manually enter [imdb id or link], [f] forget the file, [d] delete the file: ",
	)

	if !input.Scan() {
		return abort
	}
	text := input.Text()

	switch text {
	case "f":
		log.Printf("not implemented")
		return retry
	case "d":
		log.Printf("not implemented")
		return retry
	case "a":
		return abort
	default:
		imdbID, err := strconv.Atoi(text)
		if err != nil {
			fmt.Printf("unable to read imdb id: %s\n", err)
			return retry
		}
		show, err := importer.Library.GetShowByImdbID(imdbID)
		if err != nil {
			fmt.Printf("unable to get parse imdb id: %s\n", err)
			return retry
		}
		fmt.Printf("adding show with imdb id %d\n", imdbID)

		file.OsdbError = nil

		show.Files = append(show.Files, file)
		shows <- &library.ShowWithFile{
			Show: show,
			File: file,
		}

		return success
	}
}

func fixFileErrors(c *cli.Context, importer *importer.Context) {
	files := importer.FilesWithErrors
	fmt.Printf(
		"%d of the files have errors. Let's walk through them:\n", len(files),
	)

	shows := make(chan *library.ShowWithFile, importer.Config.Importer.BufferSize)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		importer.ProcessShows(shows)
		wg.Done()
	}()
	defer wg.Wait()
	defer close(shows)

	for i := range files {
		fmt.Printf("[%s]\n", files[i].Path)
		if files[i].ImportError != nil {
			fmt.Printf(" import error: %s\n", *files[i].ImportError)
		}
		if files[i].OsdbError != nil {
			fmt.Printf("%s\n", *files[i].OsdbError)
		}

		for {
			switch manualImport(c, importer, shows, files[i]) {
			case abort:
				fmt.Printf("aborting.\n")
				return
			case retry:
				continue
			case success:
				break
			}
		}
	}
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
	defer close(importer.Stop)

	if len(importer.FilesWithErrors) > 0 {
		if c.GlobalBool("non-interactive") {
			log.Printf("warning: there have been errors while importing some files.")
		} else {
			fixFileErrors(c, importer)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "mvm"
	app.Usage = "identify, manipulate and search movies and series"

	app.Commands = []cli.Command{
		{
			Name:      "import",
			Aliases:   []string{"imp", "i"},
			Usage:     "import video files into the library",
			ArgsUsage: "<filename> [filename2] ...",
			Action:    runImport,
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-file, c",
			Value: "",
			Usage: "path to the mvm configuration file",
		},
		cli.BoolFlag{
			Name:  "non-interactive, n",
			Usage: "commands will not display interactive prompts for humans",
		},
	}

	app.RunAndExitOnError()
}
