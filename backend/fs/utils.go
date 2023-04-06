package fs

import (
	"encoding/hex"
	"strings"
)

func Encode(val string) string {
	return hex.EncodeToString([]byte(val))
}

func Decode(val string) (string, error) {
	if decoded, err := hex.DecodeString(val); err != nil {
		return "", err
	} else {
		return string(decoded), nil
	}
}

func TrimSlash(val string) string {
	return strings.TrimLeftFunc(val, func(r rune) bool {
		return r == slash
	})
}
