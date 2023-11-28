package service

import (
	"github.com/spf13/afero"
	vm "github.com/taubyte/vm/ifaces"
	wasm "github.com/taubyte/vm/ifaces/wasm"
)

func (s *engine) New(ctx vm.Context, config wasm.Config) (wasm.Instance, error) {
	r := &instance{
		ctx:    ctx,
		engine: s,
		config: &config,
		fs:     afero.NewMemMapFs(),
		deps:   make(map[string]wasm.SourceModule, 0),
	}

	switch config.Output {
	case wasm.Buffer:
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

func (s *engine) Source() wasm.Source {
	return s.source
}

// TODO, improve close method to nicely close down services.
// maybe offer an optional "Node closed what now method."
func (s *engine) Close() error {
	s.ctxC()
	return nil
}
