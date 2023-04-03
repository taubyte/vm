package service

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
	helpers "github.com/taubyte/vm/helpers/wazero"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

var _ vm.Instance = &instance{}

func (i *instance) Runtime(functionDefs ...*vm.HostModuleFunctionDefinition) (vm.Runtime, error) {
	rt := helpers.NewRuntime(i.ctx.Context())
	r := &runtime{
		instance:      i,
		wasiStartDone: make(chan bool, 1),
		runtime:       rt,
	}

	r.ctx, r.ctxC = context.WithCancel(i.ctx.Context())

	go func() {
		<-r.ctx.Done()
		r.Close()
	}()

	hm, err := r.Expose("env")
	if err != nil {
		return nil, fmt.Errorf("exposing `env` failed with: %s", err)
	}

	moduleFunctions := r.defaultModuleFunctions()
	moduleFunctions = append(moduleFunctions, functionDefs...)

	err = hm.Functions(moduleFunctions)
	if err != nil {
		return nil, fmt.Errorf("adding default host module functions failed with: %s", err)
	}

	_, err = hm.Compile()
	if err != nil {
		return nil, fmt.Errorf("compiling host module failed with: %s", err)
	}

	_, err = wasi.NewBuilder(r.runtime).Instantiate(r.ctx)
	if err != nil {
		return nil, fmt.Errorf("instantiating host module failed with: %s", err)
	}

	return r, nil
}

func (r *instance) Stdout() io.Reader {
	return r.output
}

func (r *instance) Stderr() io.Reader {
	return r.outputErr
}

func (r *instance) Filesystem() afero.Fs {
	return r.fs
}

func (r *instance) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return nil
}

func (r *instance) Context() vm.Context {
	return r.ctx
}
