package url

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/i18n"
)

type backend struct {
	http http.Client
}

var (
	Scheme = "http"
)

func New(client http.Client) vm.Backend {
	return &backend{
		http: client,
	}
}

func (b *backend) Close() error {
	return nil
}

func (b *backend) Get(uri string) (io.ReadCloser, error) {
	_uri, err := url.Parse(uri)
	if err != nil {
		return nil, i18n.ParseError(uri, err)
	}

	if _uri.Scheme != Scheme && _uri.Scheme != Scheme+"s" {
		return nil, i18n.SchemeError(_uri, b)
	}

	res, err := b.http.Get(uri)
	if err != nil {
		return nil, i18n.RetrieveError(uri, err, b)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewBuffer(data)), nil
}

func (b *backend) Scheme() string {
	return Scheme
}
