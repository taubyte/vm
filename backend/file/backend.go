package file

import (
	"io"
	"os"

	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/i18n"

	ma "github.com/multiformats/go-multiaddr"
	resolv "github.com/taubyte/vm/resolvers/taubyte"
)

type backend struct{}

func New() vm.Backend {
	return &backend{}
}

func (b *backend) Close() error {
	return nil
}

func (b *backend) Scheme() string {
	return resolv.FILE_PROTOCOL_NAME
}

func (b *backend) Get(multiAddr ma.Multiaddr) (io.ReadCloser, error) {
	protocols := multiAddr.Protocols()
	if protocols[0].Code != resolv.P_FILE {
		return nil, i18n.MultiAddrCompliant(multiAddr, resolv.FILE_PROTOCOL_NAME)
	}

	path, err := multiAddr.ValueForProtocol(resolv.P_FILE)
	if err != nil {
		return nil, i18n.ParseProtocol(resolv.FILE_PROTOCOL_NAME, err)
	}

	// remove extra slash
	path = path[1:]

	file, err := os.Open(path)
	if err != nil {
		return nil, i18n.RetrieveError(path, err, b)
	}

	return file, nil
}
