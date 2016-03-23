package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DexterLB/mvm/library"
)

// FileImporter takes filenames and constructs video file data
func (c *Context) FileImporter(filenames <-chan string, files chan<- *library.VideoFile) error {
	for filename := range filenames {
		relativePath, err := relative(c.Config.FileRoot, filename)
		if err != nil {
			return fmt.Errorf("Invalid filename: %s", err)
		}

		file, err := c.Library.GetFileByPath(relativePath)
		if err != nil {
			return fmt.Errorf("Library error while looking up file: %s")
		}

		file.Size, err = filesize(filename)
		if err != nil {
			file.Status["file"].Errorf("unable to get file size: %s", err)
		}

		files <- file
	}
	return nil
}

func filesize(filename string) (uint64, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return uint64(fi.Size()), nil
}

func relative(root string, path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// todo: maybe return something in the form "../../../foo" for files
	// outside the root dir
	return strings.TrimPrefix(filepath.Clean(root), filepath.Clean(absolute)), nil
}
