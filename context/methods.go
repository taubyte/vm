package context

import gocontext "context"

func (c *vmContext) Context() gocontext.Context {
	return c.ctx
}

func (c *vmContext) Close() error {
	c.ctxC()
	return nil
}

func (c *vmContext) Project() string {
	return c.projectId
}

func (c *vmContext) Application() string {
	return c.applicationId
}

func (c *vmContext) Resource() string {
	return c.resourceId
}

func (c *vmContext) Branch() string {
	return c.branch
}

func (c *vmContext) Commit() string {
	return c.commit
}
