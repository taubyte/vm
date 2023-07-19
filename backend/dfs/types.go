package dfs

import (
	"io"

	"github.com/taubyte/go-interfaces/vm"
	peer "github.com/taubyte/p2p/peer"
)

var _ vm.Backend = &backend{}

type backend struct {
	node peer.Node
}

type zWasmReadCloser struct {
	dag        io.ReadCloser
	unCompress io.ReadCloser
}

type zipReadCloser struct {
	io.ReadCloser
	parent io.Closer
}
