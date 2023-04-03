package dfs

import (
	"compress/lzw"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	peer "github.com/taubyte/go-interfaces/p2p/peer"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/i18n"
)

// TODO: Move to specs
var (
	Scheme     = "dfs"
	GetTimeout = 3 * time.Second
)

type backend struct {
	ctx  context.Context
	ctxC context.CancelFunc
	node peer.Node
}

var _ vm.Backend = &backend{}

func New(node peer.Node) vm.Backend {
	b := &backend{
		node: node,
	}

	b.ctx, b.ctxC = context.WithCancel(node.Context())
	return b
}

type zWasmReadCloser struct {
	dag        io.ReadCloser
	unCompress io.ReadCloser
}

func (zw *zWasmReadCloser) Close() error {
	err := zw.unCompress.Close()
	if err != nil {
		return fmt.Errorf("closing uncompressed file failed with: %s", err)
	}
	err = zw.dag.Close()
	if err != nil {
		return fmt.Errorf("closing compressed file failed with: %s", err)
	}
	return nil
}

func (zw *zWasmReadCloser) Read(p []byte) (n int, err error) {
	return zw.unCompress.Read(p)
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
