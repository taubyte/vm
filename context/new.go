package context

import (
	gocontext "context"

	"github.com/taubyte/go-interfaces/vm"
	spec "github.com/taubyte/go-specs/common"
)

func New(ctx gocontext.Context, options ...Option) (vm.Context, error) {
	c := &vmContext{}
	c.ctx, c.ctxC = gocontext.WithCancel(ctx)
	c.branch = spec.DefaultBranch

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}
