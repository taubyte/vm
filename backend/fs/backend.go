package fs

import (
	"io"
	"net/url"
	"os"

	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/i18n"
)

type backend struct{}

var (
	Scheme = "fs"
)

func New() vm.Backend {
	return &backend{}
}

func (b *backend) Close() error {
	return nil
}

func (b *backend) Get(uri string) (io.ReadCloser, error) {
	_uri, err := url.Parse(uri)
	if err != nil {
		return nil, i18n.ParseError(uri, err)
	}

	if _uri.Scheme != Scheme {
		return nil, i18n.SchemeError(_uri, b)
	}

	file, err := os.Open(_uri.RequestURI())
	if err != nil {
		return nil, i18n.RetrieveError(_uri.RequestURI(), err, b)
	}

	return file, nil
}

func (b *backend) Scheme() string {
	return Scheme
}
