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
	assert.Equal(uint64(0x450f3f0c98a1f11d), dropFile.OsdbHash)

	/* TODO
	assert.Equal(256, dropFile.ResolutionX)
	assert.Equal(240, dropFile.ResolutionY)
	assert.Equal("yuv410p", dropFile.VideoFormat)
	assert.InDelta(30, dropFile.Framerate)
	assert.InDelta(887, dropFile.VideoBitrate)
	assert.Equal(library.Duration(6.07 * float32(time.Second)), dropFile.Duration)
	*/
}
