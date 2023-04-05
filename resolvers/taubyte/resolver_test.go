package resolver

import (
	"fmt"
	"testing"

	gocontext "context"

	functionSpec "github.com/taubyte/go-specs/function"
	"github.com/taubyte/vm/context"
	"gotest.tools/v3/assert"
)

func TestResolverDFS(t *testing.T) {
	resetVars()

	resolver := newResolver(t)

	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	uri, err := resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name))
	assert.NilError(t, err)
	assert.Equal(t, uri, fmt.Sprintf("dfs:///%s", mockConfig.Cid))
}

func TestResolverHTTP(t *testing.T) {
	resetVars()

	resolver := newResolver(t)

	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	uri, err := resolver.Lookup(ctx, "http/"+testEndPoint)
	assert.NilError(t, err)
	assert.Equal(t, testEndPoint, uri)
}

func TestResolverFS(t *testing.T) {
	resetVars()

	resolver := newResolver(t)
	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	uri, err := resolver.Lookup(ctx, "fs/"+wd)
	assert.NilError(t, err)
	assert.Equal(t, uri, "fs:///"+wd)
}

func TestResolverFailures(t *testing.T) {
	resetVars()

	resolver := newResolver(t)

	ctx, err := context.New(gocontext.Background(), contextOptions...)
	assert.NilError(t, err)

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName("hello_world")); err == nil {
		t.Error("expected error")
	}

	if _, err = resolver.Lookup(ctx, testFunc.Name); err == nil {
		t.Error("expected error")
	}

	if _, err = resolver.Lookup(ctx, "funcs/"+testFunc.Name); err == nil {
		t.Error("expected error")
	}

	ctx, err = context.New(gocontext.Background())
	assert.NilError(t, err)

	if _, err = resolver.Lookup(ctx, functionSpec.ModuleName(testFunc.Name)); err == nil {
		t.Error("expected error")
	}
}
