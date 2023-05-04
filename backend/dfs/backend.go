package dfs

import (
	"compress/lzw"
	"context"
	"io"

	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/taubyte/go-interfaces/p2p/peer"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/errors"
	resolv "github.com/taubyte/vm/resolvers/taubyte"
)

func New(node peer.Node) vm.Backend {
	return &backend{
		node: node,
	}
}

func (b *backend) Get(multiAddr ma.Multiaddr) (io.ReadCloser, error) {
	protocols := multiAddr.Protocols()
	if protocols[0].Code != resolv.P_DFS {
		return nil, errors.MultiAddrCompliant(multiAddr, resolv.DFS_PROTOCOL_NAME)
	}

	cid, err := multiAddr.ValueForProtocol(resolv.P_DFS)
	if err != nil {
		return nil, errors.ParseProtocol(resolv.DFS_PROTOCOL_NAME, err)
	}

	ctx, ctxC := context.WithTimeout(context.Background(), vm.GetTimeout)
	dagReader, err := b.node.GetFile(ctx, cid)
	ctxC()
	if err != nil {
		return nil, errors.RetrieveError(cid, err, b)
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
	b.node = nil
	return nil
}
