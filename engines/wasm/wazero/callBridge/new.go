package callbridge

import (
	wasm "github.com/taubyte/vm/ifaces/wasm"
	"github.com/tetratelabs/wazero/api"
)

func New(module api.Module) wasm.Module {
	return &callContext{wazero: module}
}
