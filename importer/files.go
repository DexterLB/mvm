package importer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/DexterLB/mvm/library"
	"github.com/oz/osdb"
)

// FileImporter takes filenames and constructs video file data
func (c *Context) FileInfo(filenames <-chan string, files chan<- *library.VideoFile) {
	for {
		select {
		case filename, ok := <-filenames:
			if !ok {
				return
			}
			relativePath, err := relative(c.Config.FileRoot, filename)
			if err != nil {
				c.Errorf("Invalid filename: %s", err)
				return
			}

			file, err := c.Library.GetFileByPath(relativePath)
			if err != nil {
				c.Errorf("Library error while looking up file: %s", err)
				return
			}

			file.Size, err = filesize(filename)
			if err != nil {
				file.Status.For("file").Errorf("unable to get file size: %s", err)
				continue
			}

			file.OsdbHash, err = osdb.Hash(filename)
			if err != nil {
				file.Status.For("file").Errorf("unable to calculate file hash: %s", err)
				continue
			}

			file.Status.For("file").Succeed()
			files <- file
		case <-c.Stop:
			return
		}
	}
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
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// todo: maybe return something in the form "../../../foo" for files
	// outside the root dir

	return strings.TrimPrefix(
		filepath.Clean(absolutePath),
		filepath.Clean(absoluteRoot)+"/",
	), nil
}
