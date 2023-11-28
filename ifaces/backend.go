package vm

import (
	"io"

	ma "github.com/multiformats/go-multiaddr"
)

type Backend interface {
	// Returns the URI scheme the backend supports.
	Scheme() string
	// Get attempts to retrieve the WASM asset.
	Get(multAddr ma.Multiaddr) (io.ReadCloser, error)
	// Close will close the Backend.
	Close() error
}
