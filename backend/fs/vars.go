package fs

import (
	"unicode/utf8"
)

var slash, _ = utf8.DecodeLastRuneInString("\u002F")
