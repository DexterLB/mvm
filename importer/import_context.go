package importer

import (
	"github.com/DexterLB/mvm/library"
)

type Context struct {
	// Close this channel to stop all import workers
	Cancel chan struct{}

	Library *library.Library
	Config  *Config
}

type Config struct {
	fileRoot string
}
