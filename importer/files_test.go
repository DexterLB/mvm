package importer

import (
	"testing"

	"github.com/DexterLB/mvm/library"
	"github.com/stretchr/testify/assert"
)

func TestFileImporter(t *testing.T) {
	context := testContext(t)

	filenames := make(chan string, 5)
	files := make(chan *library.VideoFile, 5)

	go context.FileImporter(filenames, files)

	filenames <- "testdata/drop.avi"
	close(filenames)

	dropFile := <-files

	close(context.Stop)

	assert := assert.New(t)

	assert.Equal(library.Success, dropFile.Status.For("file").Status)
	assert.Equal("drop.avi", dropFile.Path)
	assert.Equal(uint64(675840), dropFile.Size)
}
