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
	"github.com/taubyte/vm/backend/fs"
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
	splitAddress := strings.Split(name, "/")
	if len(splitAddress) < 2 {
		return "", fmt.Errorf("invalid name `%s`", name)
	}

	// Local module relative to project
	if len(splitAddress[0]) > 1 {
		if len(splitAddress) != 2 {
			return "", fmt.Errorf("invalid local module name got `%s` expected `<module-type>/<module-name>`", name)
		}

		moduleType, moduleName := splitAddress[0], splitAddress[1]
		switch moduleType {
		case functionSpec.PathVariable.String(), smartOpSpec.PathVariable.String(), librarySpec.PathVariable.String():
			return internalDFSPath(ctx, s.tns, moduleType, moduleName)
		default:
			return "", fmt.Errorf("unknown local module type: `%s`", moduleType)
		}
	}

	addressType, address := splitAddress[1], strings.Join(splitAddress[2:], "/")

	switch addressType {
	case "url":
		return address, nil
	case "dfs":
		return fmt.Sprintf("dfs:///%s", address), nil
	case "fs":
		return fmt.Sprintf("fs:///%s", fs.Encode(address)), nil
	default:
		return "", fmt.Errorf("unknown mutli-address type: `%s`", addressType)
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

	parser, err := extract.Tns().BasicPath(currentPath[0].String())
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
