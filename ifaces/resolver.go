package vm

import (
	ma "github.com/multiformats/go-multiaddr"
)

type Resolver interface {
	// Lookup resolves a module name and returns the uri
	Lookup(ctx Context, name string) (ma.Multiaddr, error)
}
