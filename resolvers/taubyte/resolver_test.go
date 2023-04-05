package resolver_test

import (
	"testing"

	gocontext "context"

	"github.com/taubyte/go-interfaces/services/tns/mocks"
	"github.com/taubyte/go-interfaces/vm"
	functionSpec "github.com/taubyte/go-specs/function"
	"github.com/taubyte/go-specs/methods"
	"github.com/taubyte/vm/context"
	"github.com/taubyte/vm/test"

	"gotest.tools/v3/assert"
)

func basicLookUp(t *testing.T, global bool, module, expectedUri string) (mocks.MockedTns, vm.Resolver, vm.Context) {
	tns, resolver, err := test.Resolver(global)
	assert.NilError(t, err)

	ctx, err := test.Context()
	assert.NilError(t, err)

	if len(module) > 0 {
		uri, err := resolver.Lookup(ctx, module)
		assert.NilError(t, err)

		if len(expectedUri) > 0 {
			assert.Equal(t, uri, expectedUri)
		}
	}

	return tns, resolver, ctx
}

func TestResolverHTTP(t *testing.T) {
	test.ResetVars()
	basicLookUp(t, false, "http/"+test.TestEndPoint, test.TestEndPoint)
}

func TestResolverFS(t *testing.T) {
	test.ResetVars()
	basicLookUp(t, false, "fs/"+test.Wd, "fs:///"+test.Wd)
}

func TestResolverDFS(t *testing.T) {
	test.ResetVars()
	basicLookUp(t, false, functionSpec.ModuleName(test.TestFunc.Name), "dfs:///"+test.MockConfig.Cid)
}

func TestResolverDFSGlobal(t *testing.T) {
	moduleName := functionSpec.ModuleName(test.TestFunc.Name)

	tns, resolver, ctx := basicLookUp(t, true, moduleName, "dfs:///"+test.MockConfig.Cid)

	// Test Failures
	wasmPath, err := functionSpec.Tns().WasmModulePath(ctx.Project(), "", test.TestFunc.Name)
	assert.NilError(t, err)

	// replace the wasm current path with nil, rather than a string array
	err = tns.Push(wasmPath.Slice(), nil)
	assert.NilError(t, err)

	// Current call failure: object retrieved by Current call is expected to be a []string, in this case its nil, thus failing
	if _, err = resolver.Lookup(ctx, moduleName); err == nil {
		t.Error("expected error")
		return
	}

	tns.Delete(wasmPath)
	// TNS Fetch wasm module path Failure: the tns store has been deleted, resulting in failure to fetch.
	if _, err = resolver.Lookup(ctx, moduleName); err == nil {
		t.Error("expected error")
		return
	}
}

func TestResolverDFSFailures(t *testing.T) {
	test.ResetVars()

	tns, resolver, ctx := basicLookUp(t, false, "", "")

	assetHash, err := methods.GetTNSAssetPath(ctx.Project(), ctx.Resource(), ctx.Branch())
	assert.NilError(t, err)

	// Replace asset index with nil value
	err = tns.Push(assetHash.Slice(), nil)
	assert.NilError(t, err)

	moduleName := functionSpec.ModuleName(test.TestFunc.Name)

	// Typecast error: Expected asset cid to be string, but nil value is retrieved
	if _, err = resolver.Lookup(ctx, moduleName); err == nil {
		t.Error("expected error")
		return
	}

	// Delete the assetHash index
	tns.Delete(assetHash)

	// Fetch Error: assetHash index does not exist, thus Fetch fails
	if _, err = resolver.Lookup(ctx, moduleName); err == nil {
		t.Error("expected error")
		return
	}

	wasmPath, err := functionSpec.Tns().WasmModulePath(ctx.Project(), ctx.Application(), test.TestFunc.Name)
	assert.NilError(t, err)

	// Push empty `current`` path to the `current` list
	tns.Push(wasmPath.Slice(), []string{""})

	// Parser Error: `current` path is parsed using regex, an empty string value results in failure of the parser
	if _, err = resolver.Lookup(ctx, moduleName); err == nil {
		t.Error("expected error")
		return
	}

	// Push invalid `current` path, to be used to create asset hash.
	tns.Push(wasmPath.Slice(), []string{"current"})

	// AssetHash Error: the asset hash helper method requires a project Id, resource Id, and branch
	// All are empty thus resulting in failure
	if _, err = resolver.Lookup(ctx, moduleName); err == nil {
		t.Error("expected error")
		return
	}

	// Push multiple `current` values
	tns.Push(wasmPath.Slice(), []string{"current", "current"})

	// Current Path Length Error: There may not be more than one `current` path, thus failure
	if _, err = resolver.Lookup(ctx, moduleName); err == nil {
		t.Error("expected error")
		return
	}

	// Fetch Error: no function "hello_world" has been registered to TNS
	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName("hello_world")); err == nil {
		t.Error("expected error")
		return
	}

	// Create context with no Project,resource, application, branch, or commit
	ctx, err = context.New(gocontext.Background())
	assert.NilError(t, err)

	// WasmModulePathFromModule Error: WasmModulePathFromModule requires a project, and application
	if _, err = resolver.Lookup(ctx, moduleName); err == nil {
		t.Error("expected error")
		return
	}
}

func TestResolverLookupFailures(t *testing.T) {
	test.ResetVars()

	_, resolver, err := test.Resolver(false)
	assert.NilError(t, err)

	ctx, err := test.Context()
	assert.NilError(t, err)

	// Module name should be in convention <type>/<name>
	if _, err = resolver.Lookup(ctx, test.TestFunc.Name); err == nil {
		t.Error("expected error")
		return
	}

	// Module type `funcs` is not recognized, for functions module type is `functions`
	if _, err = resolver.Lookup(ctx, "funcs/"+test.TestFunc.Name); err == nil {
		t.Error("expected error")
		return
	}
}
