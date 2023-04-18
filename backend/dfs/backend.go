package dfs

import (
	"compress/lzw"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	peer "github.com/taubyte/go-interfaces/p2p/peer"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/i18n"
)

func New(node peer.Node) vm.Backend {
	b := &backend{
		node: node,
	}

	b.ctx, b.ctxC = context.WithCancel(node.Context())
	return b
}

func (b *backend) Get(uri string) (io.ReadCloser, error) {
	_uri, err := url.Parse(uri)
	if err != nil {
		return nil, i18n.ParseError(uri, err)
	}

	if _uri.Scheme != Scheme {
		return nil, i18n.SchemeError(_uri, b)
	}

	if len(_uri.User.String()) != 0 || len(_uri.Host) != 0 {
		return nil, fmt.Errorf("unsupported uri `%s`", uri)
	}

	path := strings.Split(_uri.Path, "/")
	if len(path) != 2 || len(path[0]) != 0 /* the split will run "/" into "" */ {
		return nil, fmt.Errorf("invalid path in uri `%s`", uri)
	}

	// cid is second element in path
	// caching just to make code more readable
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
