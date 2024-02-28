package test_utils

import (
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/tau/protocols/tns/mocks"

	"github.com/taubyte/go-specs/common"
	functionSpec "github.com/taubyte/go-specs/function"
	resolvers "github.com/taubyte/vm/resolvers/taubyte"
)

func Resolver(global bool) (tnsClient mocks.MockedTns, resolver vm.Resolver, err error) {
	tnsClient = mocks.New()
	config := MockConfig

	if global {
		config = MockGlobalConfig

		var wasmPath *common.TnsPath
		wasmPath, err = functionSpec.Tns().WasmModulePath(config.Project, MockConfig.Application, TestFunc.Name)
		if err != nil {
			return
		}

		if err = tnsClient.Push(wasmPath.Slice(), nil); err != nil {
			return
		}
	}

	if err = tnsClient.Inject(&TestFunc, config); err != nil {
		return
	}

	resolver = resolvers.New(tnsClient)

	return
}
