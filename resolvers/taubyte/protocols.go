package resolver

import (
	ma "github.com/multiformats/go-multiaddr"
)

const (
	P_DFS             = 4242
	DFS_PROTOCOL_NAME = "dfs"

	P_FILE             = 9999
	FILE_PROTOCOL_NAME = "file"

	P_PATH             = 0x2F
	PATH_PROTOCOL_NAME = "path"
)

var ()

var internalProtocols = []ma.Protocol{
	{
		Name:       DFS_PROTOCOL_NAME,
		Code:       P_DFS,
		VCode:      ma.CodeToVarint(P_DFS),
		Size:       ma.LengthPrefixedVarSize,
		Path:       true,
		Transcoder: ma.TranscoderUnix,
	},
	{
		Name:       FILE_PROTOCOL_NAME,
		Code:       P_FILE,
		VCode:      ma.CodeToVarint(P_FILE),
		Size:       ma.LengthPrefixedVarSize,
		Path:       true,
		Transcoder: ma.TranscoderUnix,
	},
	{
		Name:       PATH_PROTOCOL_NAME,
		Code:       P_PATH,
		VCode:      ma.CodeToVarint(P_PATH),
		Size:       ma.LengthPrefixedVarSize,
		Path:       true,
		Transcoder: ma.TranscoderUnix,
	},
}

func init() {
	for _, protocol := range internalProtocols {
		if err := ma.AddProtocol(protocol); err != nil {
			panic(err)
		}
	}
}
