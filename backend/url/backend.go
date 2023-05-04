package url

import (
	"bytes"
	"io"
	"net/http"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/backend/i18n"
)

type backend struct {
	http http.Client
}

func New(client http.Client) vm.Backend {
	return &backend{
		http: client,
	}
}

func (b *backend) Get(multiAddr ma.Multiaddr) (io.ReadCloser, error) {
	protocols, err := isMADns(multiAddr)
	if err != nil {
		return nil, err
	}

	uri, err := maUriFormat(multiAddr, protocols)
	if err != nil {
		return nil, err
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
	return "url"
}

func (b *backend) Close() error {
	return nil
}
