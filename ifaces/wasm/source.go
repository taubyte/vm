package vm

import (
	vm "github.com/taubyte/vm/ifaces"
)

type Source interface {
	// Module Loads the given module name, and returns the SourceModule
	Module(ctx vm.Context, name string) (SourceModule, error)
}

type SourceModule []byte
