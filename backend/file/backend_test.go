package file

import (
	"io"
	"os"
	"path"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	fixtures "github.com/taubyte/vm/fixtures/wasm"
	vm "github.com/taubyte/vm/ifaces"
	resolv "github.com/taubyte/vm/resolvers/taubyte"
	"gotest.tools/v3/assert"
)

func TestFS(t *testing.T) {
	backend := New()
	assert.Equal(t, backend.Scheme(), resolv.FILE_PROTOCOL_NAME)

	wd, err := os.Getwd()
	assert.NilError(t, err)

	relativePath := "../../fixtures/wasm/recursive.wasm"
	fsPath := path.Join(wd, relativePath)

	testMA(t, backend, "/file/"+fsPath, false)
	testMA(t, backend, "/file/"+relativePath, false)
}

func testMA(t *testing.T, be vm.Backend, raw string, fail bool) {
	mAddr, err := ma.NewMultiaddr(raw)
	assert.NilError(t, err)

	fsReader, err := be.Get(mAddr)
	if fail {
		if err == nil {
			t.Error("expected error")
			return
		}
	}
	assert.NilError(t, err)

	fsData, err := io.ReadAll(fsReader)
	assert.NilError(t, err)

	assert.DeepEqual(t, fsData, fixtures.NonCompressRecursive)
}
