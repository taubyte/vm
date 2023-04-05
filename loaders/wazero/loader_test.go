package loader_test

import (
	"fmt"
	"io"
	"testing"

	functionSpec "github.com/taubyte/go-specs/function"
	fixtures "github.com/taubyte/vm/fixtures/wasm"
	loaders "github.com/taubyte/vm/loaders/wazero"
	"github.com/taubyte/vm/test"
	"gotest.tools/v3/assert"
)

func TestLoader(t *testing.T) {
	test.ResetVars()

	cid, loader, resolver, _, simple, err := test.Loader()
	assert.NilError(t, err)

	ctx, err := test.Context()
	assert.NilError(t, err)

	moduleName := functionSpec.ModuleName(test.TestFunc.Name)

	reader, err := loader.Load(ctx, moduleName)
	assert.NilError(t, err)

	source, err := io.ReadAll(reader)
	assert.NilError(t, err)

	assert.DeepEqual(t, fixtures.NonCompressRecursive, source)

	// Test Failures

	// Delete ipfs stored file
	err = simple.DeleteFile(cid)
	assert.NilError(t, err)

	// No Reader Error: All backends have been checked, but all returned nil readers.
	if _, err = loader.Load(ctx, moduleName); err == nil {
		t.Error("expected error")
	}

	fmt.Println(err)

	// New Loader with no backends
	loader = loaders.New(resolver)
	// Backend Error: Creating a loader with no backends results in failure
	if _, err = loader.Load(ctx, moduleName); err == nil {
		t.Error("expected error")
	}

	// Lookup Error: Attempting to load module that does not follow convention of <type>/<name>
	if _, err = loader.Load(ctx, test.TestFunc.Name); err == nil {
		t.Error("expected error")
	}

}
