package i18n

import (
	"fmt"
	"net/url"

	"github.com/taubyte/go-interfaces/vm"
)

func ParseError(uri string, err error) error {
	return fmt.Errorf("parsing uri(`%s`) failed with: %s", uri, err)
}

func SchemeError(uri *url.URL, backend vm.Backend) error {
	return fmt.Errorf("unsupported Scheme `%s` in `%s` expected `%s`", uri.Scheme, uri.String(), backend.Scheme())
}

func RetrieveError(path string, err error, backend vm.Backend) error {
	return fmt.Errorf("retrieving ReadCloser through `%s` backend from `%s` failed with: %s", backend.Scheme(), path, err)
}
