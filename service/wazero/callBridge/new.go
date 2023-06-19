package callbridge

import (
	"github.com/taubyte/go-interfaces/vm"
	"github.com/tetratelabs/wazero/api"
)

func New(module api.Module) vm.Module {
	return &callContext{wazero: module}
}
