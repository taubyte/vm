package source

import (
	"fmt"
	"io"

	vm "github.com/taubyte/vm/ifaces"
	wasm "github.com/taubyte/vm/ifaces/wasm"
)

type source struct {
	loader vm.Loader
}

var _ wasm.Source = &source{}

func New(loader vm.Loader) wasm.Source {
	return &source{
		loader: loader,
	}
}

func (s *source) Module(ctx vm.Context, name string) (wasm.SourceModule, error) {
	_reader, err := s.loader.Load(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("loading module `%s` failed with %w", name, err)
	}
	defer _reader.Close()

	_source, err := io.ReadAll(_reader)
	if err != nil {
		return nil, fmt.Errorf("reading module `%s` failed with %w", name, err)
	}

	return _source, nil
}
