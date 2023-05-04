package url

import (
	"errors"
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/taubyte/vm/backend/i18n"
	resolv "github.com/taubyte/vm/resolvers/taubyte"
)

func isMADns(multiAddr ma.Multiaddr) (protocols []ma.Protocol, err error) {
	protocols = multiAddr.Protocols()
	switch protocols[0].Code {
	case ma.P_DNS, ma.P_DNS4, ma.P_DNS6:
		return
	default:
		return nil, i18n.MultiAddrCompliant(multiAddr, "url")
	}
}

func maUriFormat(multiAddr ma.Multiaddr, protocols []ma.Protocol) (uri string, err error) {
	var http, host, path string

	for _, protocol := range protocols {
		switch protocol.Code {
		case ma.P_HTTP:
			http = "http"
		case ma.P_HTTPS:
			http = "https"
		case ma.P_DNS, ma.P_DNS4, ma.P_DNS6:
			host, err = multiAddr.ValueForProtocol(protocol.Code)
			if err != nil {
				return "", i18n.ParseProtocol(protocol.Name, err)
			}
		case resolv.P_PATH:
			path, err = multiAddr.ValueForProtocol(protocol.Code)
			if err != nil {
				return "", i18n.ParseProtocol(protocol.Name, err)
			}
		}
	}

	if len(http) == 0 || len(host) == 0 {
		return "", errors.New("multi address does not include host or http(s) specification")
	}

	return fmt.Sprintf("%s://%s%s", http, host, path), nil
}
