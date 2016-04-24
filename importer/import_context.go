package importer

import (
	"fmt"
	"sync"

	"github.com/DexterLB/mvm/config"
	"github.com/DexterLB/mvm/library"
	"github.com/oz/osdb"
)

// Context contains common data for all importers
type Context struct {
	// Close this channel to stop all import workers
	Stop chan struct{}

	Library *library.Library
	Config  *config.Config

	// Channel for unrecoverable pipeline errors, to be read by a human
	Errors chan error

	// Files which have failed to identify correctly during import
	FilesWithErrors []*library.VideoFile

	osdbClient *osdb.Client
	osdbLock   sync.Mutex
}

// NewContext initializes a context with the given library and config
func NewContext(library *library.Library, config *config.Config) *Context {
	context := &Context{
		Stop:    make(chan struct{}),
		Library: library,
		Config:  config,
		Errors:  make(chan error),
	}
	go func() {
		<-context.Stop
		close(context.Errors)
	}()
	return context
}

// Errorf sends an error message to the Errors channel
func (c *Context) Errorf(message string, arguments ...interface{}) {
	c.Errors <- fmt.Errorf(message, arguments)
}
