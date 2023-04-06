package service

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"

	"github.com/taubyte/go-interfaces/vm"
	wasm "github.com/tetratelabs/wazero/api"
)

var _ vm.Return = &wasmReturn{}

func (r *wasmReturn) Error() error {
	return r.err
}

func (r *wasmReturn) Reflect(args ...interface{}) error {
	if r.err != nil {
		return r.err
	}

	j := 0
	for _, arg := range args {
		valptr := reflect.ValueOf(arg)
		if valptr.Kind() != reflect.Ptr {
			return fmt.Errorf("need to pass a pointer")
		}
		val := valptr.Elem()
		switch val.Kind() {
		case reflect.Float32:
			if r.types[j] != wasm.ValueTypeF64 && r.types[j] != wasm.ValueTypeF32 {
				return fmt.Errorf("can not convert non float32 value to float32")
			}

			var b [8]byte
			binary.LittleEndian.PutUint64(b[:], r.values[j])
			val.SetFloat(float64(math.Float32frombits(binary.LittleEndian.Uint32(b[0:4]))))
			j++
		case reflect.Float64:
			if r.types[j] != wasm.ValueTypeF64 && r.types[j] != wasm.ValueTypeF32 {
				return fmt.Errorf("can not convert `%s` non float64 value to float64", val.Kind().String())
			}

			val.SetFloat(math.Float64frombits(r.values[j]))
			j++
		case reflect.Uint64, reflect.Uint32:
			if r.types[j] != wasm.ValueTypeI64 && r.types[j] != wasm.ValueTypeI32 {
				return fmt.Errorf("can not convert non uint value to uint")
			}
			val.SetUint(r.values[j])
			j++
		case reflect.Int64, reflect.Int32:
			if r.types[j] != wasm.ValueTypeI64 && r.types[j] != wasm.ValueTypeI32 {
				return fmt.Errorf("can not convert non int value to int")
			}
			val.SetInt(int64(r.values[j]))
			j++
		case reflect.Array:
		}
	}

	return nil
}
