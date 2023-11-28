package vm

import (
	"io"
)

type Loader interface {
	// Load resolves the module, then loads the module using a Backend
	Load(ctx Context, name string) (io.ReadCloser, error)
}
