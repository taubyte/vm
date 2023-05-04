package url

import (
	"io"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	"gotest.tools/v3/assert"
)

func TestBackend(t *testing.T) {
	backend := New()
	assert.Equal(t, backend.Scheme(), "url")

	httpUrl := "/dns4/ping.examples.tau.link/https/path/ping"
	incorrectUris := []string{
		"/dns4/ping.examples.tau",
		"/file//tmp/test",
		"/file/tmp/test",
	}

	for _, uri := range incorrectUris {
		mAddr, err := ma.NewMultiaddr(uri)
		if err != nil {
			t.Error(err)
			return
		}

		if _, err := backend.Get(mAddr); err == nil {
			t.Error("expected error")
		}
	}

	// Missing Coverage: Not sure how to get error for read all on successful http get without adding a mock http client
	mAddr, err := ma.NewMultiaddr(httpUrl)
	assert.NilError(t, err)

	httpReader, err := backend.Get(mAddr)
	assert.NilError(t, err)

	data, err := io.ReadAll(httpReader)
	assert.NilError(t, err)

	assert.DeepEqual(t, data, []byte("PONG EXAMPLE"))

	if err = backend.Close(); err != nil {
		t.Error(err)
	}
}
