package test

import (
	"io"

	peer "github.com/taubyte/go-interfaces/p2p/peer/mocks"
	tns "github.com/taubyte/go-interfaces/services/tns/mocks"
	"github.com/taubyte/go-interfaces/vm"
	loaders "github.com/taubyte/vm/loaders/wazero"
)

func Loader(injectReader io.Reader) (cid string, loader vm.Loader, resolver vm.Resolver, tns tns.MockedTns, simple peer.MockedNode, err error) {
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
