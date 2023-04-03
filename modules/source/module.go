package module

import (
	"fmt"

	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/utils/wasm"
	"github.com/taubyte/vm/utils/wasm/binary"
)

type module struct {
	source          []byte
	imports         []string
	importsByModule map[string][]string
	holdRuntime     bool
}

var _ vm.SourceModule = &module{}

func New(source []byte) (vm.SourceModule, error) {
	decoded, err := binary.DecodeModule(source, wasm.Features20220419, wasm.MemorySizer)
	if err != nil {
		return nil, fmt.Errorf("Decoding sections failed with %w", err)
	}

	imports := make([]string, len(decoded.ImportSection))
	importsByModule := make(map[string][]string)
	for i, imp := range decoded.ImportSection {
		imports[i] = imp.Module
		if importsByModule[imp.Module] == nil {
			importsByModule[imp.Module] = make([]string, 0)
		}
		importsByModule[imp.Module] = append(importsByModule[imp.Module], imp.Name)
	}

	return &module{
		source:          source,
		imports:         imports,
		importsByModule: importsByModule,
	}, nil
}

func (m *module) Source() []byte {
	return m.source
}

func (m *module) Imports() []string {
	return m.imports
}

func (m *module) ImportsByModule(name string) []string {
	return m.importsByModule[name]
}

func (m *module) ImportsFunction(module, name string) bool {
	if m.importsByModule[module] == nil {
		return false
	}

	for _, fname := range m.importsByModule[module] {
		if fname == name {
			return true
		}
	}

	return false
}
