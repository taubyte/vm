package fixtures

import (
	"bytes"
	"compress/lzw"
	_ "embed"
	"io"
)

var (
	//go:embed recursive.wasm
	NonCompressRecursive []byte // non compressed
)

var Recursive []byte // compressed

func init() {
	buf := bytes.NewBuffer(nil)
	sbuf := bytes.NewBuffer(NonCompressRecursive)
	w := lzw.NewWriter(buf, lzw.LSB, 8)
	io.Copy(w, sbuf)
	w.Close()
	Recursive = buf.Bytes()
}
