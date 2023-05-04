package dfs

import (
	"compress/lzw"
	"context"
	"fmt"
	"io"
	"strings"

	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/taubyte/go-interfaces/p2p/peer"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/i18n"
	resolv "github.com/taubyte/vm/resolvers/taubyte"
)

func New(node peer.Node) vm.Backend {
	b := &backend{
		node: node,
	}

	b.ctx, b.ctxC = context.WithCancel(node.Context())
	return b
}

func (b *backend) Get(multiAddr ma.Multiaddr) (io.ReadCloser, error) {
	protocols := multiAddr.Protocols()
	if protocols[0].Code != resolv.P_DFS {
		return nil, i18n.MultiAddrCompliant(multiAddr, resolv.DFS_PROTOCOL_NAME)
	}

	module, err := multiAddr.ValueForProtocol(resolv.P_DFS)
	if err != nil {
		return nil, i18n.ParseProtocol(resolv.DFS_PROTOCOL_NAME, err)
	}

	path := strings.Split(module, "/")
	if len(path) != 2 || len(path[0]) != 0 {
		return nil, fmt.Errorf("invalid module name expected `/<cid>` got `%s`", module)
	}

	cid := path[1]
	dagReader, err := b.node.GetFile(b.ctx, cid)
	if err != nil {
		return nil, i18n.RetrieveError(cid, err, b)
	}

	return &zWasmReadCloser{
		dag:        dagReader,
		unCompress: lzw.NewReader(dagReader, lzw.LSB, 8),
	}, nil
}

func (b *backend) Scheme() string {
	return Scheme
}

func (b *backend) Close() error {
	b.ctxC()
	return nil
}
