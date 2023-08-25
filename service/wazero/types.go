package service

import (
	"container/list"
	"context"
	"io"
	"sync"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
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
	runtime   *runtime
	functions map[string]functionDef
	memories  map[string]memoryPages
	globals   map[string]interface{}
}

/*************** Instance ***************/

type instance struct {
	ctx       vm.Context
	service   vm.Service
	lock      sync.RWMutex
	fs        afero.Fs
	config    *vm.Config
	output    io.ReadWriteCloser
	outputErr io.ReadWriteCloser
	deps      map[string]vm.SourceModule
	runtimes  *list.List
	rtLock    sync.Mutex
}

/*************** Module Instance ***************/
type moduleInstance struct {
	parent *runtime
	module api.Module
	ctx    context.Context
}

/*************** Runtime ***************/

type runtime struct {
	instance   *instance
	runtime    wazero.Runtime
	removeFunc func() *runtime

	wasiStartError error
	wasiStartDone  chan bool
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
	types  []api.ValueType
	values []uint64
}
