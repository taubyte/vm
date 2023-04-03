package service

import (
	"context"
	"fmt"

	"github.com/taubyte/go-interfaces/vm"
)

var _ vm.ModuleInstance = &moduleInstance{}

func (m *moduleInstance) Function(name string) (vm.FunctionInstance, error) {
	funcInst := m.module.ExportedFunction(name)
	if funcInst == nil {
		return nil, fmt.Errorf("Function (%s).`%s` does not exist", m.module.Name(), name)
	}

	f := &funcInstance{
		module:   m,
		function: funcInst,
	}

	f.ctx, f.ctxC = context.WithCancel(m.ctx)

	return f, nil
}