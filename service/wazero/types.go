package service

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/utils/wasm"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

/*************** Function Instance ***************/

type funcInstance struct {
	ctx      context.Context
	ctxC     context.CancelFunc
	module   *moduleInstance
	function api.Function
	timeout  time.Duration
}

/*************** Host Module ***************/

type functionDef struct {
	handler vm.HostFunction
}

type memoryPages struct {
	min   uint64
	max   uint64
	maxed bool
}

type hostModule struct {
	ctx       vm.Context
	name      string
	parent    *runtime
	functions map[string]functionDef
	memories  map[string]memoryPages
	globals   map[string]interface{}
}

/*************** Instance ***************/

type instance struct {
	ctx        vm.Context
	service    vm.Service
	lock       sync.RWMutex
	fs         afero.Fs
	output     *bytes.Buffer
	outputErr  *bytes.Buffer
	compileMap map[string]wazero.CompiledModule
	deps       map[string]vm.SourceModule
}

/*************** Module Instance ***************/

type moduleInstance struct {
	module api.Module
	ctx    context.Context
}

/*************** Runtime ***************/

type runtime struct {
	instance       *instance
	runtime        wazero.Runtime
	ctx            context.Context
	ctxC           context.CancelFunc
	wasiStartError error
	wasiStartDone  chan bool
	lock           sync.RWMutex
}

/*************** Service ***************/

type service struct {
	ctx    context.Context
	ctxC   context.CancelFunc
	source vm.Source
}

/*************** Wasm Return ***************/

type wasmReturn struct {
	err    error
	types  []wasm.ValueType
	values []uint64
}
