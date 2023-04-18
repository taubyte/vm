package test_utils

import (
	gocontext "context"

	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/vm/context"
)

func Context() (vm.Context, error) {
	return context.New(gocontext.Background(), ContextOptions...)
}
