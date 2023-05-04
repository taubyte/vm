package dfs

import (
	"io"

	peer "github.com/taubyte/go-interfaces/p2p/peer"
	"github.com/taubyte/go-interfaces/vm"
)

var _ vm.Backend = &backend{}

type backend struct {
	node peer.Node
}

type zWasmReadCloser struct {
	dag        io.ReadCloser
	unCompress io.ReadCloser
}
