package module

import (
	"testing"

	fixtures "github.com/taubyte/vm/fixtures/wasm"
	"gotest.tools/v3/assert"
)

var (
	sourceImportLen = 14
)

func TestSource(t *testing.T) {
	module, err := New(fixtures.NonCompressRecursive)
	assert.NilError(t, err)

	imports := module.Imports()
	assert.DeepEqual(t, module.Source(), fixtures.NonCompressRecursive)
	assert.Equal(t, len(imports), sourceImportLen)

	onlyImport := "wasi_snapshot_preview1"
	for i := 0; i < sourceImportLen; i++ {
		if imports[i] != onlyImport {
			t.Errorf("expected only `%s` as import`", onlyImport)
		}
	}

	assert.Equal(t, sourceImportLen, len(module.ImportsByModule(onlyImport)))
	assert.Equal(t, true, module.ImportsFunction(onlyImport, "fd_write"))

	assert.Equal(t, false, module.ImportsFunction(onlyImport, "hello_world"))
	assert.Equal(t, false, module.ImportsFunction("hello_world", onlyImport))
}
