package service

import (
	"bytes"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
)

var MaxOutputCapacity = 10 * 1024

func (s *service) New(ctx vm.Context, config vm.Config) (vm.Instance, error) {
	r := &instance{
		ctx:       ctx,
		service:   s,
		config:    &config,
		fs:        afero.NewMemMapFs(),
		output:    bytes.NewBuffer(make([]byte, 0, MaxOutputCapacity)),
		outputErr: bytes.NewBuffer(make([]byte, 0, MaxOutputCapacity)),
		deps:      make(map[string]vm.SourceModule, 0),
	}

	return r, nil
}

func (s *service) Source() vm.Source {
	return s.source
}

// TODO, improve close method to nicely close down services.
// maybe offer an optional "Node closed what now method."
func (s *service) Close() error {
	s.ctxC()
	return nil
}
