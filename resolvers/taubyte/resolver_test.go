package resolver

import (
	"fmt"
	"testing"

	gocontext "context"

	functionSpec "github.com/taubyte/go-specs/function"
	"github.com/taubyte/go-specs/methods"
	"github.com/taubyte/vm/context"
	"gotest.tools/v3/assert"
)

func TestResolverDFS(t *testing.T) {
	resetVars()

	_, resolver := newResolver(t)

	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	uri, err := resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name))
	assert.NilError(t, err)
	assert.Equal(t, uri, fmt.Sprintf("dfs:///%s", mockConfig.Cid))
}

func TestResolverHTTP(t *testing.T) {
	resetVars()

	_, resolver := newResolver(t)

	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	uri, err := resolver.Lookup(ctx, "http/"+testEndPoint)
	assert.NilError(t, err)
	assert.Equal(t, testEndPoint, uri)
}

func TestResolverFS(t *testing.T) {
	resetVars()

	_, resolver := newResolver(t)
	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	uri, err := resolver.Lookup(ctx, "fs/"+wd)
	assert.NilError(t, err)
	assert.Equal(t, uri, "fs:///"+wd)
}

func TestResolverDFSFailures(t *testing.T) {
	resetVars()

	tns, resolver := newResolver(t)

	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	assetHash, err := methods.GetTNSAssetPath(ctx.Project(), testFunc.Id, ctx.Branch())
	assert.NilError(t, err)

	tns.Push(assetHash.Slice(), nil)
	assert.NilError(t, err)

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
		return
	}

	tns.Delete(assetHash)

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
		return
	}

	wasmPath, err := functionSpec.Tns().WasmModulePath(mockConfig.Project, mockConfig.Application, testFunc.Name)
	assert.NilError(t, err)

	tns.Push(wasmPath.Slice(), []string{""})

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
		return
	}

	tns.Push(wasmPath.Slice(), []string{"current"})

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
		return
	}

	tns.Push(wasmPath.Slice(), []string{"current", "current"})

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
		return
	}

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName("hello_world")); err == nil {
		t.Error("expected error")
		return
	}

	ctx, err = context.New(gocontext.Background())
	assert.NilError(t, err)

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
		return
	}
}

func TestResolverDFSGlobal(t *testing.T) {
	resetVars()

	tns, resolver := newGlobalResolver(t)
	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	wasmPath, err := functionSpec.Tns().WasmModulePath(mockConfig.Project, mockConfig.Application, testFunc.Name)
	assert.NilError(t, err)

	err = tns.Push(wasmPath.Slice(), nil)
	assert.NilError(t, err)

	_, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name))
	assert.NilError(t, err)

	wasmPath, err = functionSpec.Tns().WasmModulePath(mockConfig.Project, "", testFunc.Name)
	assert.NilError(t, err)

	err = tns.Push(wasmPath.Slice(), nil)
	assert.NilError(t, err)

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
		return
	}

	tns.Delete(wasmPath)
	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
		return
	}
}

func TestResolverLookupFailures(t *testing.T) {
	resetVars()

	_, resolver := newResolver(t)
	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	if _, err = resolver.Lookup(ctx, testFunc.Name); err == nil {
		t.Error("expected error")
		return
	}

	if _, err = resolver.Lookup(ctx, "funcs/"+testFunc.Name); err == nil {
		t.Error("expected error")
		return
	}
}
