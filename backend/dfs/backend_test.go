package dfs_test

import (
	"bytes"
	"compress/lzw"
	"io"
	"testing"

	"github.com/taubyte/vm/backend/dfs"
	fixtures "github.com/taubyte/vm/fixtures/wasm"
	"github.com/taubyte/vm/test_utils"
	"gotest.tools/v3/assert"
)

func TestBackEnd(t *testing.T) {
	backend, err := test_utils.DFSBackend().Inject(bytes.NewReader(fixtures.Recursive))
	assert.NilError(t, err)
	assert.Equal(t, backend.Scheme(), dfs.Scheme)

	incorrectUris := []string{
		"dfs://" + backend.Cid,
		"Dfs://" + backend.Cid,
		"DFS://" + backend.Cid,
		"dfs:///file/" + backend.Cid,
		"dfs:///Fake" + backend.Cid,
		"hello world" + backend.Cid,
		// ASCII control character for coverage
		string([]byte{0x7f}) + backend.Cid,
	}

	for _, uri := range incorrectUris {
		if _, err = backend.Get(uri); err == nil {
			t.Error("expected error")
		}
	}

	dagReader, err := backend.Get("dfs:///" + backend.Cid)
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

	err = dagReader.Close()
	assert.NilError(t, err)

	// Missing Close coverage as dagReader Close does not seem to fail

	backend.Close()
}
