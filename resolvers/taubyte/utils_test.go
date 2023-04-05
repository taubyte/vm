package resolver

import (
	"os"
	"testing"

	"github.com/taubyte/go-interfaces/services/tns/mocks"
	"github.com/taubyte/go-interfaces/vm"
	structureSpec "github.com/taubyte/go-specs/structure"
	"github.com/taubyte/utils/id"
	"github.com/taubyte/vm/context"
	"gotest.tools/v3/assert"
)

var (
	testFunc       structureSpec.Function
	mockConfig     mocks.InjectConfig
	contextOptions []context.Option

	testEndPoint = "https://ping.examples.tau.link/ping"
	wd           string
)

func resetVars() {
	mockConfig = mocks.InjectConfig{
		Branch:      "master",
		Commit:      "head_commit",
		Project:     id.Generate(),
		Application: id.Generate(),
		Cid:         id.Generate(),
	}

	testFunc = structureSpec.Function{
		Id:      id.Generate(),
		Name:    "basicFunc",
		Type:    "http",
		Memory:  10000,
		Timeout: 100000000,
		Method:  "GET",
		Source:  ".",
		Call:    "basic",
		Paths:   []string{"/ping"},
		Domains: []string{"somDomain"},
	}

	contextOptions = []context.Option{
		context.Application(mockConfig.Application),
		context.Project(mockConfig.Project),
		context.Resource(testFunc.Id),
		context.Branch(mockConfig.Branch),
		context.Commit(mockConfig.Commit),
	}

	var err error
	wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func init() {
	resetVars()
}

func newResolver(t *testing.T) vm.Resolver {
	tnsClient := mocks.New()
	err := tnsClient.Inject(&testFunc, mockConfig)
	assert.NilError(t, err)

	return New(tnsClient)
}
