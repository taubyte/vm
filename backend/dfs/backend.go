package dfs

import (
	"compress/lzw"
	"context"
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
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

	_cid, err := multiAddr.ValueForProtocol(resolv.P_DFS)
	if err != nil {
		return nil, errors.ParseProtocol(resolv.DFS_PROTOCOL_NAME, err)
	}

	__cid, err := cid.Decode(_cid)
	if err != nil {
		return nil, err
	}

	ctx, ctxC := context.WithTimeout(b.node.Context(), vm.GetTimeout)
	defer ctxC()

	ok, err := b.node.DAG().BlockStore().Has(ctx, __cid)
	if !ok || err != nil {
		dagReader, err := b.node.GetFile(ctx, _cid)
		if err != nil {
			return nil, fmt.Errorf("caching CID `%s` failed with:  %w", _cid, err)
		}

		dagReader.Close()
	}

	dagReader, err := b.node.GetFile(ctx, _cid)
	if err != nil {
		return nil, errors.RetrieveError(_cid, err, b)
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
