package helpers

import (
	"context"

	"github.com/tetratelabs/wazero"
)

func NewRuntime(ctx context.Context) wazero.Runtime {
	return wazero.NewRuntimeWithConfig(
		ctx,
		wazero.NewRuntimeConfig(),
	)
}
