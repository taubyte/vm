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
	if err := r.runtime.Close(context.TODO()); err != nil {
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
	r.lock.Lock()
	defer r.lock.Unlock()
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

		// handle imports
		for _, dep := range module.Imports() {
			if dep == "env" {
				continue
			}
			_, err := r.module(dep)
			if err != nil {
				return nil, fmt.Errorf("loading module `%s` dependency `%s` failed with: %s", name, dep, err)
			}

		}

		// then start
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
		module: modInst,
		ctx:    r.ctx,
	}, nil
}

func (r *runtime) instantiate(name string, module vm.SourceModule) error {
	r.instance.lock.RLock()
	compiled, ok := r.instance.compileMap[name]
	r.instance.lock.RUnlock()
	if !ok {
		var err error
		compiled, err = r.runtime.CompileModule(r.ctx, module.Source())
		if err != nil {
			return fmt.Errorf("compiling module `%s` failed with: %s", name, err)
		}

		r.instance.lock.Lock()
		r.instance.compileMap[name] = compiled
		r.instance.lock.Unlock()
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

	// wazero instance will source the instance in it's source
	// (which is diffrent from our source as it sources instances)
	m, err := r.runtime.InstantiateModule(r.ctx, compiled, config)
	if err != nil {
		return fmt.Errorf("instantiating compiled module `%s` failed with: %s", name, err)
	}

	// execute _start and keep it running as long as the module is running
	// this ensures that if the language has a runtime, it'll be running fine
	if _start := m.ExportedFunction("_start"); _start != nil {
		if module.ImportsFunction("env", "_ready") {
			go func() {
				_, r.wasiStartError = _start.Call(r.ctx)
			}()

			// wait for any runtime initialization
			<-r.wasiStartDone
		} else {
			_start.Call(r.ctx)
		}
	}

	return nil
}

func (r *runtime) defaultModuleFunctions() []*vm.HostModuleFunctionDefinition {
	return []*vm.HostModuleFunctionDefinition{
		{
			Name: "_ready",
			Handler: func(ctx context.Context) {
				r.wasiStartDone <- true
				<-r.ctx.Done()
			},
		},
		{
			Name: "_sleep",
			Handler: func(ctx context.Context, dur int64) {
				time.Sleep(time.Duration(dur))
			},
		},
	}
}
