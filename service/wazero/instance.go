package service

import (
	"container/list"
	"fmt"
	"io"
	"time"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
	helpers "github.com/taubyte/vm/helpers/wazero"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

var _ vm.Instance = &instance{}

func (i *instance) runtime(hostDef *vm.HostModuleDefinitions) (*runtime, error) {
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

var MaxRuntimeAge time.Duration = 30 * time.Minute
var MaxRuntimes int = 3

func (i *instance) Runtime(hostDef *vm.HostModuleDefinitions) (vm.Runtime, error) {
	if i.runtimes == nil {
		i.runtimes = list.New()
	}

	runtimeChan := make(chan *runtime, 1)
	errChan := make(chan error, 1)

	go func() {
		for i.runtimes.Len() < MaxRuntimes {
			rt, err := i.runtime(hostDef)
			if err != nil {
				errChan <- err
				return
			}

			i.rtLock.Lock()
			rtEl := i.runtimes.PushBack(rt)
			i.rtLock.Unlock()

			rt.removeFunc = func() *runtime {
				i.rtLock.Lock()
				i.runtimes.Remove(rtEl)
				rt.removeFunc = func() *runtime {
					return nil
				}
				i.rtLock.Unlock()
				return rt
			}

			go func(element *list.Element, runtime *runtime) {
				<-time.After(MaxRuntimeAge)
				i.rtLock.Lock()

				_rt := runtime.removeFunc()
				if _rt != nil {
					_rt.Close()
				}
			}(rtEl, rt)
		}
	}()

	go func() {
		var done bool
		for !done {
			select {
			case <-errChan:
				return
			default:
				if i.runtimes.Len() > 0 {
					i.rtLock.Lock()
					rt := i.runtimes.Front()
					i.rtLock.Unlock()
					_rt := rt.Value.(*runtime)
					_rt.removeFunc()
					runtimeChan <- _rt
				}
			}
		}
	}()

	for {
		select {
		case rt := <-runtimeChan:
			return rt, nil
		case err := <-errChan:
			return nil, err
		}
	}
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
