package importer

import (
	"testing"

	"github.com/DexterLB/mvm/library"
	"github.com/stretchr/testify/assert"
)

func TestFileInfo(t *testing.T) {
	context := testContext(t)

	filenames := make(chan string, 5)
	files := make(chan *library.VideoFile, 5)

	go context.FileInfo(filenames, files)

	filenames <- "fixtures/drop.avi"
	close(filenames)

	dropFile := <-files

	if _, ok := <-files; ok {
		t.Errorf("files channel not closed after reading all files")
	}

	close(context.Stop)

	assert := assert.New(t)

	assert.Nil(dropFile.ImportError)
	assert.Equal("drop.avi", dropFile.Path)
	assert.Equal(uint64(675840), dropFile.Size)
	assert.Equal(uint64(0x450f3f0c98a1f11d), uint64(dropFile.OsdbHash))

	/* TODO
	assert.Equal(256, dropFile.ResolutionX)
	assert.Equal(240, dropFile.ResolutionY)
	assert.Equal("yuv410p", dropFile.VideoFormat)
	assert.InDelta(30, dropFile.Framerate)
	assert.InDelta(887, dropFile.VideoBitrate)
	assert.Equal(types.Duration(6.07 * float32(time.Second)), dropFile.Duration)
	*/
}
