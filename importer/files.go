package importer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/DexterLB/mvm/library"
	"github.com/DexterLB/mvm/types"
	"github.com/DexterLB/osdb"
)

// FileInfo takes filenames and constructs video file data
func (c *Context) FileInfo(filenames <-chan string, files chan<- *library.VideoFile) {
	defer close(files)

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
				file.ImportError = types.Errorf(
					"unable to get file size: %s", err,
				)
				continue
			}

			hash, err := osdb.Hash(filename)
			if err != nil {
				file.ImportError = types.Errorf(
					"unable to calculate file hash: %s", err,
				)
				continue
			}
			file.OsdbHash = types.BigUint64(hash)

			file.ImportError = nil
			files <- file
		case <-c.Stop:
			return
		}
	}
}

// WalkPaths recursively searches for video files in the given directories
// and sends them on the channel. Non-folder paths are sent as-are.
func (c *Context) WalkPaths(paths []string, filenames chan<- string) {
	defer close(filenames)

	// TODO: actually walk directories
	for i := range paths {
		filenames <- paths[i]
	}
}

func filesize(filename string) (uint64, error) {
	file, err := os.Open(filename)
	defer func() {
		_ = file.Close()
	}()

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
