package dfs

import (
	"context"
	"io"

	peer "github.com/taubyte/go-interfaces/p2p/peer"
	"github.com/taubyte/go-interfaces/vm"
)

var _ vm.Backend = &backend{}

type backend struct {
	ctx  context.Context
	ctxC context.CancelFunc
	node peer.Node
}

type zWasmReadCloser struct {
	dag        io.ReadCloser
	unCompress io.ReadCloser
}
