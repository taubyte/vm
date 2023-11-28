package service

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
	helpers "github.com/taubyte/vm/engines/wasm/wazero/helpers"
	vm "github.com/taubyte/vm/ifaces"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	wasm "github.com/taubyte/vm/ifaces/wasm"
)

var _ wasm.Instance = &instance{}

func (i *instance) Runtime(hostDef *wasm.HostModuleDefinitions) (wasm.Runtime, error) {
	rt := helpers.NewRuntime(i.ctx.Context(), i.config)
	r := &runtime{
		instance:      i,
		wasiStartDone: make(chan bool),
		runtime:       rt,
	}

	hm, err := r.Expose("env")
	if err != nil {
		return nil, fmt.Errorf("exposing `env` failed with: %w", err)
	}

	moduleFunctions := r.defaultModuleFunctions()

	if hostDef != nil {
		moduleFunctions = append(moduleFunctions, hostDef.Functions...)

		if err = hm.Globals(hostDef.Globals...); err != nil {
			return nil, fmt.Errorf("adding global definitions to host module failed with: %w", err)
		}

		if err = hm.Memories(hostDef.Memories...); err != nil {
			return nil, fmt.Errorf("adding memory definitions to host module failed with: %w", err)
		}

	}

	if err = hm.Functions(moduleFunctions...); err != nil {
		return nil, fmt.Errorf("adding function definitions to host module with: %s", err)
	}

	if _, err = hm.Compile(); err != nil {
		return nil, fmt.Errorf("compiling host module failed with: %s", err)

	}

	if _, err = wasi.NewBuilder(r.runtime).Instantiate(r.instance.ctx.Context()); err != nil {
		return nil, fmt.Errorf("instantiating host module failed with: %s", err)
	}

	return r, nil
}

func (i *instance) Stdout() io.Reader {
	return i.output
}

func (i *instance) Stderr() io.Reader {
	return i.outputErr
}

func (i *instance) Filesystem() afero.Fs {
	return i.fs
}

func (i *instance) Close() error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.output.Close()
	i.outputErr.Close()
	return nil
}

func (i *instance) Context() vm.Context {
	return i.ctx
}
