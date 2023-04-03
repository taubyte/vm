package service

import (
	"fmt"
	"reflect"

	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/common"
	wasm "github.com/tetratelabs/wazero/api"
)

var moduleType = reflect.TypeOf((*common.Module)(nil)).Elem()
var wazeroModuleType = reflect.TypeOf((*wasm.Module)(nil)).Elem()

func (hm *hostModule) Function(def *vm.HostModuleFunctionDefinition) error {
	if _, exists := hm.functions[def.Name]; exists {
		return fmt.Errorf("function `%s` @ `%s` already defined", def.Name, hm.name)
	}

	tp := reflect.TypeOf(def.Handler)

	count := tp.NumIn()
	_in := make([]reflect.Type, count)

	for i := 0; i < count; i++ {
		in := tp.In(i)
		if in.Kind() == reflect.Interface && in.Implements(moduleType) {
			_in[i] = wazeroModuleType
		} else {
			_in[i] = in
		}
	}

	count = tp.NumOut()
	_out := make([]reflect.Type, count)
	for i := 0; i < count; i++ {
		_out[i] = tp.Out(i)
	}

	_func := reflect.MakeFunc(
		reflect.FuncOf(_in, _out, false),
		func(args []reflect.Value) []reflect.Value {

			for i := 0; i < 2; i++ {
				if len(args) > i && args[i].Kind() == reflect.Interface && args[i].Type().Implements(wazeroModuleType) {
					args[i] = reflect.ValueOf(&callContext{wazero: args[i].Interface().(wasm.Module)})
				}
			}

			return reflect.ValueOf(def.Handler).Call(args)
		})

	hm.functions[def.Name] = functionDef{
		handler: _func.Interface(),
	}
	return nil
}

func (hm *hostModule) Functions(defs []*vm.HostModuleFunctionDefinition) error {
	for _, fdef := range defs {
		err := hm.Function(fdef)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hm *hostModule) Memory(def *vm.HostModuleMemoryDefinition) error {
	if _, exists := hm.memories[def.Name]; exists {
		return fmt.Errorf("memory `%s` @ `%s` already defined", def.Name, hm.name)
	}
	hm.memories[def.Name] = memoryPages{
		min:   def.Pages.Min,
		max:   def.Pages.Max,
		maxed: def.Pages.Maxed,
	}
	return nil
}

func (hm *hostModule) Global(def *vm.HostModuleGlobalDefinition) error {
	if _, exists := hm.globals[def.Name]; exists {
		return fmt.Errorf("global `%s` @ `%s` already defined", def.Name, hm.name)
	}
	hm.globals[def.Name] = def.Value
	return nil
}

func (hm *hostModule) Globals(defs []*vm.HostModuleGlobalDefinition) error {
	for _, gdef := range defs {
		err := hm.Global(gdef)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hm *hostModule) Compile() (vm.ModuleInstance, error) {
	wb := hm.parent.runtime.NewHostModuleBuilder(hm.name)
	// Export functions after translation if needed
	for name, def := range hm.functions {
		wb.NewFunctionBuilder().WithFunc(def.handler).Export(name)
	}

	cm, err := wb.Instantiate(hm.ctx.Context())
	if err != nil {
		return nil, err
	}

	return &moduleInstance{
		module: cm,
		ctx:    hm.ctx.Context(),
	}, nil
}
