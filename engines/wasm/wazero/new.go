package service

import (
	"context"

	wasm "github.com/taubyte/vm/ifaces/wasm"
)

var _ wasm.Engine = &engine{}

func New(ctx context.Context, source wasm.Source) wasm.Engine {
	s := &engine{}
	s.ctx, s.ctxC = context.WithCancel(ctx)
	s.source = source
	return s
}
