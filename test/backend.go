package test

import (
	"bytes"
	"context"
	"io"

	goHttp "net/http"

	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/dfs"
	"github.com/taubyte/vm/backend/fs"
	"github.com/taubyte/vm/backend/http"
	fixtures "github.com/taubyte/vm/fixtures/wasm"

	peer "github.com/taubyte/go-interfaces/p2p/peer/mocks"
)

type testBackend struct {
	vm.Backend
	simple peer.MockedNode
	Cid    string
}

func DFSBackend() *testBackend {
	simpleNode := peer.New(context.Background())

	return &testBackend{
		Backend: dfs.New(simpleNode),
		simple:  simpleNode,
	}
}

func (t *testBackend) Inject(r io.Reader) (*testBackend, error) {
	var err error
	t.Cid, err = t.simple.AddFile(r)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func HTTPBackend() vm.Backend {
	return http.New(*goHttp.DefaultClient)
}

func AllBackends(injectDefault bool) (cid string, simpleNode peer.MockedNode, backends []vm.Backend, err error) {
	dfsBe := DFSBackend()
	if injectDefault {
		if dfsBe, err = dfsBe.Inject(bytes.NewReader(fixtures.Recursive)); err != nil {
			return
		}
	}

	return dfsBe.Cid, dfsBe.simple, []vm.Backend{HTTPBackend(), dfsBe, fs.New()}, nil
}
