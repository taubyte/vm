package service

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
)

var _ vm.Instance = &instance{}

func (i *instance) Load(hostModuleDefs *vm.HostModuleDefinitions) error {

	return nil
}

func (i *instance) Attach(plugin vm.Plugin) (vm.PluginInstance, vm.ModuleInstance, error) {
	if err := i.checkRuntime(); err != nil {
		return nil, nil, err
	}

	hm := &hostModule{
		ctx:       i.ctx,
		name:      plugin.Name(),
		runtime:   i.runtime,
		functions: make(map[string]functionDef),
		memories:  make(map[string]memoryPages),
		globals:   make(map[string]interface{}),
	}

	pi, err := plugin.New(i)
	if err != nil {
		return nil, nil, fmt.Errorf("creating new plugin instance failed with: %s", err)
	}

	mInst, err := pi.Load(hm)
	if err != nil {
		return nil, nil, fmt.Errorf("loading plugin instance failed with: %s", err)
	}

	return pi, mInst, nil
}

func (i *instance) Module(name string) (vm.ModuleInstance, error) {
	i.runtime.lock.Lock()
	defer i.runtime.lock.Unlock()
	return i.module(name)
}

func (i *instance) module(name string) (vm.ModuleInstance, error) {
	if err := i.checkRuntime(); err != nil {
		return nil, err
	}

	modInst := i.runtime.primitive.Module(name)
	if modInst == nil {
		// assume module was not instantiated
		// get it from source

		i.lock.RLock()
		module, ok := i.deps[name]
		i.lock.RUnlock()

		var err error
		if !ok {
			module, err = i.service.Source().Module(i.ctx, name)
			if err != nil {
				return nil, fmt.Errorf("loading module `%s` failed with: %s", name, err)
			}

			i.lock.Lock()
			i.deps[name] = module
			i.lock.Unlock()
		}

		// handle imports
		deps := make(map[string]bool)
		for _, dep := range module.Imports() {
			if dep == "env" {
				continue
			}
			deps[dep] = true
		}

		for dep := range deps {
			_, err := i.module(dep)
			if err != nil {
				return nil, fmt.Errorf("loading module `%s` dependency `%s` failed with: %s", name, dep, err)
			}
		}

		// then start
		err = i.instantiate(name, module)
		if err != nil {
			return nil, fmt.Errorf("creating an instance of module `%s` failed with: %s", name, err)
		}

		modInst = i.runtime.primitive.Module(name)
		if modInst == nil {
			return nil, fmt.Errorf("unknown error with module `%s`", name)
		}

	}

	return &moduleInstance{
		module: modInst,
		ctx:    i.runtime.ctx,
	}, nil
}

func (i *instance) Expose(name string) (vm.HostModule, error) {
	if err := i.checkRuntime(); err != nil {
		return nil, err
	}

	return &hostModule{
		ctx:       i.ctx,
		name:      name,
		runtime:   i.runtime,
		functions: make(map[string]functionDef),
		memories:  make(map[string]memoryPages),
		globals:   make(map[string]interface{}),
	}, nil
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
	if i.checkRuntime() == nil {
		i.runtime.ctxC()
	}

	return nil
}

func (i *instance) Context() vm.Context {
	return i.ctx
}

func (i *instance) checkRuntime() error {
	if i.runtime == nil {
		return errors.New("runtime not loaded")
	}

	return nil
}
