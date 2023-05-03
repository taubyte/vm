package backend

import (
	"errors"

	goHttp "net/http"

	"github.com/taubyte/go-interfaces/p2p/peer"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/dfs"
	"github.com/taubyte/vm/backend/fs"
	"github.com/taubyte/vm/backend/url"
)

// New returns all available backends
func News(node peer.Node, httpClient goHttp.Client) ([]vm.Backend, error) {
	if node == nil {
		return nil, errors.New("node is nil")
	}

	return []vm.Backend{dfs.New(node), url.New(httpClient)}, nil
}

func NewDev(node peer.Node, httpClient goHttp.Client) ([]vm.Backend, error) {
	if node == nil {
		return nil, errors.New("node is nil")
	}

	return []vm.Backend{dfs.New(node), fs.New(), url.New(httpClient)}, nil
}
