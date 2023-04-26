package service

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
	helpers "github.com/taubyte/vm/helpers/wazero"
	"github.com/tetratelabs/wazero"
)

func (i *instance) hostModule(defs *vm.HostModuleDefinitions) (vm.HostModule, error) {
	hm, err := i.Expose("env")
	if err != nil {
		return nil, fmt.Errorf("exposing `env` failed with: %s", err)
	}

	funcDefs := i.defaultModuleFunctions()

	if defs != nil {
		funcDefs = append(funcDefs, defs.Functions...)

		if err = hm.Functions(funcDefs...); err != nil {
			return nil, fmt.Errorf("adding functions to host module failed with: %s", err)
		}

		if err = hm.Globals(defs.Globals...); err != nil {
			return nil, fmt.Errorf("adding globals to host module failed with: %s", err)
		}

		if err = hm.Memories(defs.Memories...); err != nil {
			return nil, fmt.Errorf("adding memories to host module failed with: %s", err)
		}
	}

	return hm, nil
}

func (i *instance) initRuntime() {
	i.runtime = &runtime{
		primitive:     helpers.NewRuntime(i.ctx.Context()),
		wasiStartDone: make(chan bool, 1),
	}

	i.runtime.ctx, i.runtime.ctxC = context.WithCancel(i.ctx.Context())

	go func() {
		<-i.runtime.ctx.Done()
		i.runtime.Close(i.ctx.Context())
	}()
}

func (i *instance) instantiate(name string, module vm.SourceModule) error {
	i.lock.RLock()
	compiled, ok := i.compileMap[name]
	i.lock.RUnlock()
	if !ok {
		var err error
		compiled, err = i.runtime.primitive.CompileModule(i.runtime.ctx, module.Source())
		if err != nil {
			return fmt.Errorf("compiling module `%s` failed with: %s", name, err)
		}

		i.lock.Lock()
		i.compileMap[name] = compiled
		i.lock.Unlock()
	}

	config := wazero.
		NewModuleConfig().
		WithName(name).
		WithStartFunctions(). // don't run _start: we need to start it in a go routine
		WithFS(afero.NewIOFS(i.fs)).
		WithStdout(i.output).
		WithStderr(i.outputErr).
		WithArgs(name).
		WithSysWalltime().
		WithSysNanotime()

	// wazero instance will source the instance in it's source
	// (which is diffrent from our source as it sources instances)
	m, err := i.runtime.primitive.InstantiateModule(i.ctx.Context(), compiled, config)
	if err != nil {
		return fmt.Errorf("instantiating compiled module `%s` failed with: %s", name, err)
	}

	// execute _start and keep it running as long as the module is running
	// this ensures that if the language has a runtime, it'll be running fine
	if module.ImportsFunction("env", "_ready") {
		if _start := m.ExportedFunction("_start"); _start != nil {
			go func() {
				_, i.runtime.wasiStartError = _start.Call(i.runtime.ctx)
			}()

			// wait for any runtime initialization
			<-i.runtime.wasiStartDone
		}
	}

	return nil
}

func (i *instance) defaultModuleFunctions() []*vm.HostModuleFunctionDefinition {
	return []*vm.HostModuleFunctionDefinition{
		{
			Name: "_ready",
			Handler: func(ctx context.Context) {
				i.runtime.wasiStartDone <- true
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