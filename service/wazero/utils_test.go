package service

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/taubyte/go-interfaces/vm"
	functionSpec "github.com/taubyte/go-specs/function"
	fixtures "github.com/taubyte/vm/fixtures/wasm"
	sources "github.com/taubyte/vm/sources/taubyte"
	"github.com/taubyte/vm/test_utils"
)

var (
	theAnswer uint32 = 42

	u32RetVal uint32
	i32RetVal int32
	f32RetVal float32
	f64RetVal float64

	controlRetVal string = "hello world"
)

func newService() (vm.Context, vm.Service, error) {
	test_utils.ResetVars()
	_, loader, _, _, _, err := test_utils.Loader(bytes.NewReader(fixtures.Artifact))
	if err != nil {
		return nil, nil, err
	}

	source := sources.New(loader)
	ctx, err := test_utils.Context()
	if err != nil {
		return nil, nil, err
	}

	return ctx, New(ctx.Context(), source), nil
}

func newBasicInstance() (vm.Instance, error) {
	ctx, service, err := newService()
	if err != nil {
		return nil, err
	}

	return service.New(ctx)
}

func newLoadedInstance() (vm.Instance, error) {
	instance, err := newBasicInstance()
	if err != nil {
		return nil, err
	}

	if err := instance.Load(
		&vm.HostModuleDefinitions{
			Functions: []*vm.HostModuleFunctionDefinition{testFunc},
			Memories:  []*vm.HostModuleMemoryDefinition{mockMemoryDef},
			Globals:   []*vm.HostModuleGlobalDefinition{mockGlobalDef},
		}); err != nil {
		return nil, err
	}

	return instance, err
}

func newModuleInstance() (vm.ModuleInstance, error) {
	instance, err := newLoadedInstance()
	if err != nil {
		return nil, err
	}

	return instance.Module(functionSpec.ModuleName(test_utils.TestFunc.Name))

}

func newFuncs(functionNames []string) (map[string]vm.FunctionInstance, error) {
	mi, err := newModuleInstance()
	if err != nil {
		return nil, err
	}

	functions := make(map[string]vm.FunctionInstance, 0)
	for _, name := range functionNames {
		function, err := mi.Function(name)
		if err != nil {
			return nil, err
		}

		functions[name] = function
	}

	return functions, nil
}

func compareError(retrieved, expected interface{}) error {
	return fmt.Errorf("got `%d` expected `%d`", retrieved, expected)
}

func callFuncs(functionNames []string) error {
	functions, err := newFuncs(functionNames)
	if err != nil {
		return err
	}

	for name, function := range functions {
		ret := function.Call(theAnswer)
		if ret.Error() != nil {
			return err
		}

		switch name {
		case "tou32":
			if err = ret.Reflect(&u32RetVal); err != nil {
				return err
			}

			if u32RetVal != theAnswer {
				return compareError(u32RetVal, theAnswer)
			}

		case "tof32":
			if err = ret.Reflect(&f32RetVal); err != nil {
				return err
			}

			if f32RetVal != float32(theAnswer) {
				return compareError(f32RetVal, theAnswer)
			}

		case "toi32":
			if err = ret.Reflect(&i32RetVal); err != nil {
				return err
			}

			if i32RetVal != int32(theAnswer) {
				return compareError(i32RetVal, theAnswer)
			}

		case "tof64":
			if err = ret.Reflect(&f64RetVal); err != nil {
				return err
			}

			if f64RetVal != float64(theAnswer) {
				return compareError(f64RetVal, theAnswer)
			}
		}

		if err = function.Cancel(); err != nil {
			return err
		}
	}

	return nil
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

var mockMemoryDef = &vm.HostModuleMemoryDefinition{
	Name: "mock",
	Pages: struct {
		Min   uint64
		Max   uint64
		Maxed bool
	}{
		Min:   0,
		Max:   10,
		Maxed: false,
	},
}

var mockGlobalDef = &vm.HostModuleGlobalDefinition{
	Name:  "mock",
	Value: "hello_world",
}

var testFunc = &vm.HostModuleFunctionDefinition{
	Name: "_test",
	Handler: func(ctx context.Context, val uint32) uint32 {
		return val
	},
}
