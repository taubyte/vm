package vm

// TODO
type JSEngine interface {
	// New(context Context, config Config) (Instance, error)
	// Source() Source
	Close() error
}

type PYEngine interface {
	// New(context Context, config Config) (Instance, error)
	// Source() Source
	Close() error
}

type ContainerEngine interface {
	// New(context Context, config Config) (Instance, error)
	// Source() Source
	Close() error
}
