package dfs

import (
	"archive/zip"
	"compress/lzw"
	"context"
	"fmt"
	"io"

	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/taubyte/go-interfaces/p2p/peer"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/go-specs/builders/wasm"
	"github.com/taubyte/vm/backend/errors"
	resolv "github.com/taubyte/vm/resolvers/taubyte"
	"go4.org/readerutil"
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

	ctx, ctxC := context.WithTimeout(b.node.Context(), vm.GetTimeout)
	dagReader, err := b.node.GetFile(ctx, cid)
	ctxC()
	if err != nil {
		return nil, errors.RetrieveError(cid, err, b)
	}

	// Backwards compatibility
	size, _ := readerutil.Size(dagReader)
	zipReader, err := zip.NewReader(
		readerutil.NewBufferingReaderAt(dagReader),
		size,
	)
	if err != nil {
		dagReader.Seek(0, io.SeekStart)
		return &zWasmReadCloser{
			dag:        dagReader,
			unCompress: lzw.NewReader(dagReader, lzw.LSB, 8),
		}, nil
	} else {
		// Trying for both main/artifact.wasm
		reader, err := zipReader.Open(wasm.WasmFile)
		if err != nil {
			reader, err = zipReader.Open(wasm.DeprecatedWasmFile)
			if err != nil {
				return nil, fmt.Errorf("reading wasm file as aritfact and main failed with: %s", err)
			}
		}

		return &zipReadCloser{
			parent:     dagReader,
			ReadCloser: reader,
		}, nil
	}
}

func (b *backend) Scheme() string {
	return Scheme
}

func (b *backend) Close() error {
	b.node = nil
	return nil
}
