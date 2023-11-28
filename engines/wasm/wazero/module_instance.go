package service

import (
	"fmt"

	wasm "github.com/taubyte/vm/ifaces/wasm"
)

var _ wasm.ModuleInstance = &moduleInstance{}

func (m *moduleInstance) Function(name string) (wasm.FunctionInstance, error) {
	funcInst := m.module.ExportedFunction(name)
	if funcInst == nil {
		return nil, fmt.Errorf("Function (%s).`%s` does not exist", m.module.Name(), name)
	}

	f := &funcInstance{
		module:   m,
		function: funcInst,
	}

	return f, nil
}

func (m *moduleInstance) Memory() wasm.Memory {
	return m.module.Memory()
}

func (m *moduleInstance) Functions() []wasm.FunctionDefinition {
	defMap := m.module.ExportedFunctionDefinitions()
	defs := make([]wasm.FunctionDefinition, len(defMap))
	for _, def := range defMap {
		defs = append(defs, def)
	}

	return defs
}
