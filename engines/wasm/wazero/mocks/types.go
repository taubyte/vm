package mocks

import (
	wasm "github.com/taubyte/vm/ifaces/wasm"
)

type MockedPlugin interface {
	wasm.Plugin
}

type mockPlugin struct {
	InstanceFail bool
}

type MockedPluginInstance interface {
	wasm.PluginInstance
}

type mockPluginInstance struct{}

type MockedModuleInstance interface {
	wasm.ModuleInstance
}

type mockModuleInstance struct {
	wasm.ModuleInstance
}

type MockedModule interface {
	wasm.Module
}

type MockedFunctionInstance interface {
	wasm.FunctionInstance
}

type mockFunctionInstance struct {
	wasm.FunctionInstance
}

type MockedReturn interface {
	wasm.Return
}

type mockReturn struct{ wasm.Return }
