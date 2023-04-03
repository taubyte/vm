package fs

import (
	"io"
	"os"
	"path"
	"testing"

	fixtures "github.com/taubyte/vm/fixtures/wasm"
	"gotest.tools/v3/assert"
)

func TestFS(t *testing.T) {
	backend := New()
	assert.Equal(t, backend.Scheme(), Scheme)

	wd, err := os.Getwd()
	assert.NilError(t, err)

	fsPath := path.Join(wd, "..", "..", "fixtures", "wasm", "recursive.wasm")

	incorrectUris := []string{
		"fs:/" + fsPath,
		"file:/" + fsPath,
		"dfs/" + fsPath,
		"hello world" + fsPath,
		// ASCII control character for coverage
		string([]byte{0x7f}) + fsPath,
	}

	for _, uri := range incorrectUris {
		if _, err := backend.Get(uri); err == nil {
			t.Errorf("Should have failed getting `%s`", uri)
			return
		}
	}

	fsReader, err := backend.Get("fs://" + fsPath)
	assert.NilError(t, err)

	fsData, err := io.ReadAll(fsReader)
	assert.NilError(t, err)

	assert.DeepEqual(t, fsData, fixtures.NonCompressRecursive)

	if err = backend.Close(); err != nil {
		t.Error(err)
	}
}
