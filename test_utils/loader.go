package test_utils

import (
	"io"

	"github.com/taubyte/go-interfaces/vm"
	tns "github.com/taubyte/odo/protocols/tns/mocks"
	"github.com/taubyte/p2p/peer"
	loaders "github.com/taubyte/vm/loaders/wazero"
)

func Loader(injectReader io.Reader) (cid string, loader vm.Loader, resolver vm.Resolver, tns tns.MockedTns, simple peer.Node, err error) {
	var backends []vm.Backend
	cid, simple, backends, err = AllBackends(injectReader)
	if err != nil {
		return
	}

	MockConfig.Cid = cid

	tns, resolver, err = Resolver(false)
	if err != nil {
		return
	}

	loader = loaders.New(resolver, backends...)

	return
}
