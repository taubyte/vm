package resolver

import (
	"fmt"
	"strings"

	"github.com/taubyte/go-interfaces/services/tns"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/go-specs/extract"
	functionSpec "github.com/taubyte/go-specs/function"
	librarySpec "github.com/taubyte/go-specs/library"
	"github.com/taubyte/go-specs/methods"
	smartOpSpec "github.com/taubyte/go-specs/smartops"
)

type resolver struct {
	tns tns.Client
}

var _ vm.Resolver = &resolver{}

func New(client tns.Client) vm.Resolver {
	return &resolver{
		tns: client,
	}
}

func (s *resolver) Lookup(ctx vm.Context, name string) (string, error) {
	// TODO:: DO THIS ALL BETTER USE REGEX

	splitModule := strings.Split(name, "/")
	if len(splitModule) < 2 {
		return "", fmt.Errorf("name should follow convention <moduleType>/<moduleName> got: `%s`", name)
	}

	moduleType, moduleName := splitModule[0], strings.Join(splitModule[1:], "/")

	switch moduleType {
	case functionSpec.PathVariable.String(), smartOpSpec.PathVariable.String(), librarySpec.PathVariable.String():
		return internalDFSPath(ctx, s.tns, moduleType, moduleName)
	case "http":
		return moduleName, nil
	case "fs":
		return fmt.Sprintf("fs:///%s", moduleName), nil
	default:
		return "", fmt.Errorf("unknown module type `%s`", moduleType)
	}

}

func internalDFSPath(ctx vm.Context, tns tns.Client, moduleType string, moduleName string) (string, error) {
	module := moduleType + "/" + moduleName
	project := ctx.Project()
	application := ctx.Application()
	wasmModulePath, err := methods.WasmModulePathFromModule(ctx.Project(), ctx.Application(), moduleType, moduleName)
	if err != nil {
		return "", fmt.Errorf("creating path for module `%s` with app: `%s` in project `%s` failed with: %s", module, application, project, err)
	}

	wasmIndex, err := tns.Fetch(wasmModulePath)
	if err != nil {
		return "", fmt.Errorf("looking up module `%s` with app: `%s` in project `%s` failed with: %s", module, application, project, err)
	}

	currentPath, err := wasmIndex.Current(ctx.Branch())
	// Checks global if cannot find using the application
	if err != nil || len(currentPath) == 0 && len(ctx.Application()) != 0 {
		wasmModulePath, err := methods.WasmModulePathFromModule(project, "", moduleType, moduleName)
		if err != nil {
			return "", fmt.Errorf("creating global module path `%s` in project `%s` failed with: %s", module, project, err)
		}

		wasmIndex, err = tns.Fetch(wasmModulePath)
		if err != nil {
			return "", fmt.Errorf("looking up global module `%s` in project `%s` failed with: %s", module, project, err)
		}

		currentPath, err = wasmIndex.Current(ctx.Branch())
		if err != nil {
			return "", fmt.Errorf("looking up current commit for global module `%s`  in project `%s` failed with: %s", module, project, err)
		}
	}

	if len(currentPath) > 1 {
		return "", fmt.Errorf("current module `%s`  in `%s` returned too many paths, theres an issue with the compiler", module, project)
	}

	object, err := tns.Fetch(currentPath[0])
	if err != nil {
		return "", fmt.Errorf("fetching current commit for module `%s`  in `%s` failed with: %s", module, project, err)
	}

	parser, err := extract.Tns().BasicPath(object.Path().String())
	if err != nil {
		return "", err
	}

	assetHash, err := methods.GetTNSAssetPath(project, parser.Resource(), parser.Branch())
	if err != nil {
		return "", fmt.Errorf("creating asset hash for module `%s` in project `%s` failed with: %s", module, project, err)
	}

	_assetCid, err := tns.Fetch(assetHash)
	if err != nil {
		return "", fmt.Errorf("fetching asset for module `%s` in project `%s` failed with: %s", module, project, err)
	}

	assetCid, ok := _assetCid.Interface().(string)
	if !ok {
		return "", fmt.Errorf("asset type `%T` unexpected for module `%s` in project `%s`", _assetCid.Interface(), module, project)
	}

	return fmt.Sprintf("dfs:///%s", assetCid), nil
}