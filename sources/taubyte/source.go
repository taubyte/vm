package source

import (
	"fmt"
	"io"

	"github.com/taubyte/go-interfaces/vm"
	module "github.com/taubyte/vm/modules/source"
)

// TODO: add caching .. or maybe use mmap
type source struct {
	loader vm.Loader
}

var _ vm.Source = &source{}

func New(loader vm.Loader) vm.Source {
	return &source{
		loader: loader,
	}
}

func (s *source) Module(ctx vm.Context, name string) (vm.SourceModule, error) {
	_reader, err := s.loader.Load(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("Loading module `%s` failed with %w", name, err)
	}

	_source, err := io.ReadAll(_reader)
	_reader.Close()
	if err != nil {
		return nil, fmt.Errorf("Reading module `%s` failed with %w", name, err)
	}

	return module.New(_source)
}
