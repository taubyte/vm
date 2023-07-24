package service

import (
	"github.com/spf13/afero"
	"github.com/taubyte/go-interfaces/vm"
)

func (s *service) New(ctx vm.Context, config vm.Config) (vm.Instance, error) {
	r := &instance{
		ctx:     ctx,
		service: s,
		config:  &config,
		fs:      afero.NewMemMapFs(),
		deps:    make(map[string]vm.SourceModule, 0),
	}

	switch config.Output {
	case vm.Buffer:
		r.output = newBuffer()
		r.outputErr = newBuffer()
	default:
		r.output = newPipe()
		r.outputErr = newPipe()
	}

	go func() {
		<-ctx.Context().Done()
		r.output.Close()
		r.outputErr.Close()
	}()

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
