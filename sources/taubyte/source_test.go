package source

import (
	"bytes"
	"testing"

	functionSpec "github.com/taubyte/go-specs/function"
	fixtures "github.com/taubyte/vm/fixtures/wasm"
	"github.com/taubyte/vm/test"
	"gotest.tools/v3/assert"
)

func TestSource(t *testing.T) {
	test.ResetVars()

	_, loader, _, _, _, err := test.Loader(bytes.NewReader(fixtures.Recursive))
	assert.NilError(t, err)

	source := New(loader)

	ctx, err := test.Context()
	assert.NilError(t, err)

	sourceModule, err := source.Module(ctx, functionSpec.ModuleName(test.TestFunc.Name))
	assert.NilError(t, err)

	sourceData := sourceModule.Source()
	assert.DeepEqual(t, fixtures.NonCompressRecursive, sourceData)

	// Test Failures

	// Load Failure: invalid module name does not follow convention <type>/<name>
	if _, err = source.Module(ctx, test.TestFunc.Name); err == nil {
		t.Error("expected error")
	}
}
