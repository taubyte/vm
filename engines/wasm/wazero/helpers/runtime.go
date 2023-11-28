package helpers

import (
	"context"
	"sync"

	wasm "github.com/taubyte/vm/ifaces/wasm"
	"github.com/tetratelabs/wazero"
)

var lock sync.Mutex

func NewRuntime(ctx context.Context, config *wasm.Config) wazero.Runtime {
	lock.Lock()
	defer lock.Unlock()
	if config.MemoryLimitPages == 0 {
		config.MemoryLimitPages = wasm.MemoryLimitPages
	}

	return wazero.NewRuntimeWithConfig(
		ctx,
		wazero.NewRuntimeConfig().
			WithCloseOnContextDone(true).
			WithDebugInfoEnabled(true).
			WithMemoryLimitPages(config.MemoryLimitPages).
			WithCompilationCache(cache),
	)
}
