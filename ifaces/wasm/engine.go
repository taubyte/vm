package vm

import (
	vm "github.com/taubyte/vm/ifaces"
)

type Engine interface {
	New(context vm.Context, config Config) (Instance, error)
	Source() Source
	Close() error
}
