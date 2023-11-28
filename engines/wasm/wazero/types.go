package service

import (
	"context"
	"io"
	"sync"

	"github.com/spf13/afero"
	vm "github.com/taubyte/vm/ifaces"
	wasm "github.com/taubyte/vm/ifaces/wasm"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

/*************** Function Instance ***************/

type funcInstance struct {
	module   *moduleInstance
	function api.Function
}

/*************** Host Module ***************/

type functionDef struct {
	handler wasm.HostFunction
}

type memoryPages struct {
	min   uint64
	max   uint64
	maxed bool
}

type hostModule struct {
	ctx       vm.Context
	name      string
	runtime   *runtime
	functions map[string]functionDef
	memories  map[string]memoryPages
	globals   map[string]interface{}
}

/*************** Instance ***************/

type instance struct {
	ctx       vm.Context
	engine    wasm.Engine
	lock      sync.RWMutex
	fs        afero.Fs
	config    *wasm.Config
	output    io.ReadWriteCloser
	outputErr io.ReadWriteCloser
	deps      map[string]wasm.SourceModule
}

/*************** Module Instance ***************/
type moduleInstance struct {
	parent *runtime
	module api.Module
	ctx    context.Context
}

/*************** Runtime ***************/

type runtime struct {
	instance *instance
	runtime  wazero.Runtime

	wasiStartError error
	wasiStartDone  chan bool
}

/*************** Service ***************/

type engine struct {
	ctx    context.Context
	ctxC   context.CancelFunc
	source wasm.Source
}

/*************** Wasm Return ***************/

type wasmReturn struct {
	err    error
	types  []api.ValueType
	values []uint64
}
