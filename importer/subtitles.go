package importer

import "github.com/DexterLB/mvm/library"

func (c *Context) SubtitleDownloader(
	files <-chan *library.VideoFile,
	subtitles chan<- *library.Subtitle,
	done chan<- *library.VideoFile,
) {
	defer close(done)
	defer close(subtitles)
}
