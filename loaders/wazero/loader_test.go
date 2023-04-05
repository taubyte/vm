package loader

import (
	"bytes"
	gocontext "context"
	"io"
	"testing"

	peer "github.com/taubyte/go-interfaces/p2p/peer/mocks"
	tns "github.com/taubyte/go-interfaces/services/tns/mocks"
	functionSpec "github.com/taubyte/go-specs/function"
	structureSpec "github.com/taubyte/go-specs/structure"
	"github.com/taubyte/utils/id"
	"github.com/taubyte/vm/backend/dfs"
	"github.com/taubyte/vm/context"
	fixtures "github.com/taubyte/vm/fixtures/wasm"
	resolvers "github.com/taubyte/vm/resolvers/taubyte"
	"gotest.tools/v3/assert"
)

var (
	testFunc = structureSpec.Function{
		Id:      id.Generate(),
		Name:    "basicFunc",
		Type:    "http",
		Memory:  10000,
		Timeout: 100000000,
		Method:  "GET",
		Source:  ".",
		Call:    "basic",
		Paths:   []string{"/ping"},
		Domains: []string{"somDomain"},
	}

	mockConfig = tns.InjectConfig{
		Branch:      "master",
		Commit:      "head_commit",
		Project:     id.Generate(),
		Application: id.Generate(),
	}

	contextOptions = []context.Option{
		context.Application(mockConfig.Application),
		context.Project(mockConfig.Project),
		context.Resource(testFunc.Id),
		context.Branch(mockConfig.Branch),
		context.Commit(mockConfig.Commit),
	}
)

func TestLoader(t *testing.T) {
	goctx := gocontext.Background()

	simpleNode := peer.New(goctx)
	backend := dfs.New(simpleNode)

	cid, err := simpleNode.AddFile(bytes.NewReader(fixtures.Recursive))
	assert.NilError(t, err)

	tnsClient := tns.New()

	mockConfig.Cid = cid
	if err = tnsClient.Inject(&testFunc, mockConfig); err != nil {
		t.Error(err)
		return
	}

	resolver := resolvers.New(tnsClient)

	loader := New(resolver, backend)

	ctx, err := context.New(goctx, contextOptions...)
	assert.NilError(t, err)

	moduleName := functionSpec.ModuleName(testFunc.Name)

	reader, err := loader.Load(ctx, moduleName)
	assert.NilError(t, err)

	source, err := io.ReadAll(reader)
	assert.NilError(t, err)

	assert.DeepEqual(t, fixtures.NonCompressRecursive, source)

	// Test Failures
	err = simpleNode.DeleteFile(cid)
	assert.NilError(t, err)

	if _, err = loader.Load(ctx, moduleName); err == nil {
		t.Error("expected error")
	}

	loader = New(resolver)
	if _, err = loader.Load(ctx, moduleName); err == nil {
		t.Error("expected error")
	}

	if _, err = loader.Load(ctx, testFunc.Name); err == nil {
		t.Error("expected error")
	}
}
