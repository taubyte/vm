package service

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/taubyte/go-interfaces/vm"
)

// REF -> https://tinygo.org/docs/concepts/compiler-internals/datatypes/

var _ vm.FunctionInstance = &funcInstance{}

func (f *funcInstance) Timeout(timeout time.Duration) vm.FunctionInstance {
	f.timeout = timeout
	f.ctx, f.ctxC = context.WithTimeout(f.ctx, f.timeout)
	return f
}

func (f *funcInstance) Call(args ...interface{}) vm.Return {
	wasm_args, err := f.golangToWasm(args)
	if err != nil {
		return &wasmReturn{
			err: err,
		}
	}

	rtypes := f.function.Definition().ResultTypes() // TODO: cache this in function
	returns, err := f.function.Call(f.ctx, wasm_args...)
	if err != nil {
		return &wasmReturn{
			err: err,
		}
	}

	return &wasmReturn{
		types:  rtypes,
		values: returns,
	}
}

func (f *funcInstance) golangToWasm(args []interface{}) ([]uint64, error) {
	wasm_args := make([]uint64, len(args))
	for i, arg := range args {
		_arg := reflect.ValueOf(arg)
		switch _arg.Kind() {
		case reflect.Float64, reflect.Float32:
			wasm_args[i] = math.Float64bits(_arg.Float())
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			wasm_args[i] = _arg.Uint()
		case reflect.Int, reflect.Int32, reflect.Int64:
			wasm_args[i] = uint64(_arg.Int())
		default:
			return nil, fmt.Errorf("failed to process arguments %v of type %T", arg, arg)
		}
	}

	return wasm_args, nil
}

func (f *funcInstance) Cancel() error {
	return nil
}