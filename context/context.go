package context

import (
	gocontext "context"

	"github.com/taubyte/go-interfaces/vm"
	spec "github.com/taubyte/go-specs/common"
)

type context struct {
	ctx  gocontext.Context
	ctxC gocontext.CancelFunc

	projectId     string
	applicationId string
	id            string
	branch        string
	commit        string
}

func New(ctx gocontext.Context, options ...Option) (vm.Context, error) {
	c := &context{}
	c.ctx, c.ctxC = gocontext.WithCancel(ctx)
	c.branch = spec.DefaultBranch

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// TODO: Change these options
func (c *context) Context() gocontext.Context {
	return c.ctx
}

func (c *context) Cancel() {
	c.ctxC()
}

func (c *context) Project() string {
	return c.projectId
}

func (c *context) Application() string {
	return c.applicationId
}

func (c *context) Id() string {
	return c.id
}

func (c *context) Branch() string {
	return c.branch
}

func (c *context) Commit() string {
	return c.commit
}
