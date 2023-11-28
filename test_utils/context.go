package test_utils

import (
	gocontext "context"

	"github.com/taubyte/vm/context"
	vm "github.com/taubyte/vm/ifaces"
)

func Context() (vm.Context, error) {
	return context.New(gocontext.Background(), ContextOptions...)
}
