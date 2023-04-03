package context

type Option func(*context) error

func Project(projectId string) Option {
	return func(ctx *context) error {
		ctx.projectId = projectId
		return nil
	}
}

func Application(applicationId string) Option {
	return func(ctx *context) error {
		ctx.applicationId = applicationId
		return nil
	}
}

func Id(id string) Option {
	return func(ctx *context) error {
		ctx.id = id
		return nil
	}
}

func Branch(branch string) Option {
	return func(ctx *context) error {
		ctx.branch = branch
		return nil
	}
}

func Commit(commit string) Option {
	return func(ctx *context) error {
		ctx.commit = commit
		return nil
	}
}
