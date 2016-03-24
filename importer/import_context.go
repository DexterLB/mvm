package importer

import (
	"github.com/DexterLB/mvm/library"
)

// Context contains common data for all importers
type Context struct {
	// Close this channel to stop all import workers
	Stop chan struct{}

	Library *library.Library
	Config  *Config

	// Channel for unrecoverable pipeline errors, to be read by a human
	Errors chan error
}

// Config contains the configuration for all importers
type Config struct {
	FileRoot string
}

// NewContext initializes a context with the given library and config
func NewContext(library *library.Library, config *Config) *Context {
	context := &Context{
		Stop:    make(chan struct{}),
		Library: library,
		Config:  config,
		Errors:  make(chan error),
	}
	go func() {
		select {
		case <-context.Stop:
			close(context.Errors)
			return
		}
	}()
	return context
}
