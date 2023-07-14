package service

import (
	"io"

	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
)

var MaxOutputCapacity = 10 * 1024

type pipe struct {
	io.ReadCloser
	io.WriteCloser
	io.Closer
}

func newPipe() *pipe {
	p := &pipe{}
	p.ReadCloser, p.WriteCloser = io.Pipe()
	return p
}

func (p *pipe) Close() {
	p.WriteCloser.Close()
	p.ReadCloser.Close()
}

func (s *service) New(ctx vm.Context, config vm.Config) (vm.Instance, error) {
	r := &instance{
		ctx:       ctx,
		service:   s,
		config:    &config,
		fs:        afero.NewMemMapFs(),
		output:    newPipe(),
		outputErr: newPipe(),
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
