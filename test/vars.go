package test

import (
	"os"

	"github.com/taubyte/go-interfaces/services/tns/mocks"
	structureSpec "github.com/taubyte/go-specs/structure"
	"github.com/taubyte/utils/id"
	"github.com/taubyte/vm/context"
)

var (
	TestFunc         structureSpec.Function
	MockConfig       mocks.InjectConfig
	MockGlobalConfig mocks.InjectConfig
	ContextOptions   []context.Option

	TestEndPoint = "https://ping.examples.tau.link/ping"
	Wd           string
)

func ResetVars() (err error) {
	TestFunc = structureSpec.Function{
		Id:      id.Generate(),
		Name:    "basic",
		Type:    "http",
		Memory:  10000,
		Timeout: 100000000,
		Method:  "GET",
		Source:  ".",
		Call:    "tou32",
		Paths:   []string{"/ping"},
		Domains: []string{"somDomain"},
	}

	MockConfig = mocks.InjectConfig{
		Branch:      "master",
		Commit:      "head_commit",
		Project:     id.Generate(),
		Application: id.Generate(),
		Cid:         id.Generate(),
	}

	MockGlobalConfig = MockConfig
	MockGlobalConfig.Application = ""

	ContextOptions = []context.Option{
		context.Application(MockConfig.Application),
		context.Project(MockConfig.Project),
		context.Resource(TestFunc.Id),
		context.Branch(MockConfig.Branch),
		context.Commit(MockConfig.Commit),
	}

	if Wd, err = os.Getwd(); err != nil {
		return err
	}

	return nil
}
