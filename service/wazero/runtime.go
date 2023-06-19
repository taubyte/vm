package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/tetratelabs/wazero"
)

func (r *runtime) Close() error {
	if err := r.runtime.Close(r.instance.ctx.Context()); err != nil {
		return err
	}

	close(r.wasiStartDone)

	return nil
}

func (r *runtime) Expose(name string) (vm.HostModule, error) {
	return &hostModule{
		ctx:       r.instance.ctx,
		name:      name,
		runtime:   r,
		functions: make(map[string]functionDef),
		memories:  make(map[string]memoryPages),
		globals:   make(map[string]interface{}),
	}, nil
}

func (r *runtime) Stdout() io.Reader {
	return r.instance.Stdout()
}

func (r *runtime) Stderr() io.Reader {
	return r.instance.Stderr()
}

func (r *runtime) Attach(plugin vm.Plugin) (vm.PluginInstance, vm.ModuleInstance, error) {
	if plugin == nil {
		return nil, nil, fmt.Errorf("plugin cannot be nil")
	}

	hm := &hostModule{
		ctx:       r.instance.ctx,
		name:      plugin.Name(),
		runtime:   r,
		functions: make(map[string]functionDef),
		memories:  make(map[string]memoryPages),
		globals:   make(map[string]interface{}),
	}

	pi, err := plugin.New(r.instance)
	if err != nil {
		return nil, nil, fmt.Errorf("creating new plugin instance failed with: %s", err)
	}

	minst, err := pi.Load(hm)
	if err != nil {
		return nil, nil, fmt.Errorf("loading plugin instance failed with: %s", err)
	}

	return pi, minst, nil
}

func (r *runtime) Module(name string) (vm.ModuleInstance, error) {
	return r.module(name)
}

func (r *runtime) module(name string) (vm.ModuleInstance, error) {
	modInst := r.runtime.Module(name)
	if modInst == nil {
		r.instance.lock.RLock()
		module, ok := r.instance.deps[name]
		r.instance.lock.RUnlock()
		var err error
		if !ok {
			module, err = r.instance.service.Source().Module(r.instance.ctx, name)
			if err != nil {
				return nil, fmt.Errorf("loading module `%s` failed with: %s", name, err)
			}

			r.instance.lock.Lock()
			r.instance.deps[name] = module
			r.instance.lock.Unlock()
		}

		// TODO: Compiled module L122 should have deps. Use that instead of this.
		// for _, dep := range module.Imports() {
		// 	if dep == "env" {
		// 		continue
		// 	}
		// 	_, err := r.module(dep)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("loading module `%s` dependency `%s` failed with: %s", name, dep, err)
		// 	}
		// }

		err = r.instantiate(name, module)
		if err != nil {
			return nil, fmt.Errorf("creating an instance of module `%s` failed with: %s", name, err)
		}

		modInst = r.runtime.Module(name)
		if modInst == nil {
			return nil, fmt.Errorf("unknown error with module `%s`", name)
		}
	}

	return &moduleInstance{
		parent: r,
		module: modInst,
		ctx:    r.instance.ctx.Context(),
	}, nil
}

func (r *runtime) instantiate(name string, module vm.SourceModule) error {
	compiled, err := r.runtime.CompileModule(r.instance.ctx.Context(), module.Source())
	if err != nil {
		return fmt.Errorf("getting compiled module failed with: %s", err)
	}

	config := wazero.
		NewModuleConfig().
		WithName(name).
		WithStartFunctions(). // don't run _start: we need to start it in a go routine
		WithFS(afero.NewIOFS(r.instance.fs)).
		WithStdout(r.instance.output).
		WithStderr(r.instance.outputErr).
		WithArgs(name).
		WithSysWalltime().
		WithSysNanotime()

	m, err := r.runtime.InstantiateModule(r.instance.ctx.Context(), compiled, config)
	if err != nil {
		return fmt.Errorf("instantiating compiled module `%s` failed with: %s", name, err)
	}

	if _start := m.ExportedFunction("_start"); _start != nil {
		if module.ImportsFunction("env", "_ready") {
			go func() {
				_, r.wasiStartError = _start.Call(r.instance.ctx.Context())
				if r.wasiStartError != nil {
					r.wasiStartDone <- false
				}
			}()

			<-r.wasiStartDone
		} else {
			_start.Call(r.instance.ctx.Context())
		}
	}

	return nil
}

func (r *runtime) defaultModuleFunctions() []*vm.HostModuleFunctionDefinition {
	return []*vm.HostModuleFunctionDefinition{
		{
			Name: "_ready",
			Handler: func(ctx context.Context, module vm.Module) {
				r.wasiStartDone <- true
			},
		},
		{
			Name: "_sleep",
			Handler: func(ctx context.Context, dur int64) {
				select {
				case <-time.After(time.Duration(dur)):
				case <-ctx.Done():
				}
			},
		},
		{
			Name: "_log",
			Handler: func(ctx context.Context, module vm.Module, data uint32, dataLen uint32) {
				msgBuf, _ := module.Memory().Read(data, dataLen)
				fmt.Println(string(msgBuf))
			},
		},
	}
}
