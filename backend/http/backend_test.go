package http

import (
	"io"
	goHttp "net/http"
	"testing"

	"gotest.tools/v3/assert"
)

func TestBackend(t *testing.T) {
	backend := New(*goHttp.DefaultClient)
	assert.Equal(t, backend.Scheme(), Scheme)

	httpUrl := "https://ping.examples.tau.link/ping"
	incorrectUris := []string{
		"fs:/" + httpUrl,
		"dfs:/" + httpUrl,
		"www.https/" + httpUrl,
		"hello world" + httpUrl,
		"https://ping.examples.tasu.link/ping",
		// ASCII control character for coverage
		string([]byte{0x7f}) + httpUrl,
	}

	for _, uri := range incorrectUris {
		if _, err := backend.Get(uri); err == nil {
			t.Error("expected error")
		}
	}

	// Missing Coverage: Not sure how to get error for read all on successful http get without adding a mock http client
	httpReader, err := backend.Get(httpUrl)
	assert.NilError(t, err)

	data, err := io.ReadAll(httpReader)
	assert.NilError(t, err)

	assert.DeepEqual(t, data, []byte("PONG EXAMPLE"))

	if err = backend.Close(); err != nil {
		t.Error(err)
	}
}
