package dfs

import (
	"bytes"
	"compress/lzw"
	"context"
	"io"
	"testing"

	peer "github.com/taubyte/go-interfaces/p2p/peer/mocks"
	fixtures "github.com/taubyte/vm/fixtures/wasm"
	"gotest.tools/v3/assert"
)

func TestBackEnd(t *testing.T) {
	ctx := context.Background()
	simpleNode := peer.New(ctx)

	backend := New(simpleNode)
	assert.Equal(t, backend.Scheme(), Scheme)

	cid, err := simpleNode.AddFile(bytes.NewReader(fixtures.Recursive))
	assert.NilError(t, err)

	incorrectUris := []string{
		"dfs://" + cid,
		"Dfs://" + cid,
		"DFS://" + cid,
		"dfs:///file/" + cid,
		"dfs:///Fake" + cid,
	}

	for _, uri := range incorrectUris {
		if _, err = backend.Get(uri); err == nil {
			t.Error("expected error")
		}
	}

	dagReader, err := backend.Get("dfs:///" + cid)
	assert.NilError(t, err)

	dfsData, err := io.ReadAll(dagReader)
	assert.NilError(t, err)

	testData, err := io.ReadAll(
		lzw.NewReader(
			bytes.NewBuffer(fixtures.Recursive),
			lzw.LSB,
			8,
		),
	)
	assert.NilError(t, err)

	assert.DeepEqual(t, testData, dfsData)
}
