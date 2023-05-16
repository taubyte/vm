package dfs_test

import (
	"bytes"
	"compress/lzw"
	"io"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/taubyte/vm/backend/dfs"
	fixtures "github.com/taubyte/vm/fixtures/wasm"
	"github.com/taubyte/vm/test_utils"
	"gotest.tools/v3/assert"
)

func TestBackEnd(t *testing.T) {
	backend, err := test_utils.DFSBackend().Inject(bytes.NewReader(fixtures.Recursive))
	assert.NilError(t, err)
	assert.Equal(t, backend.Scheme(), dfs.Scheme)

	mAddr, err := ma.NewMultiaddr("/dfs/" + backend.Cid)
	assert.NilError(t, err)
	dagReader, err := backend.Get(mAddr)
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
