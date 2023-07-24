package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type buffer struct {
	buffer *bytes.Buffer
	io.ReadWriteCloser
}

func (b *buffer) Close() error {
	b.buffer = nil
	return nil
}

func (b *buffer) Read(p []byte) (n int, err error) {
	if b.buffer != nil {
		return b.buffer.Read(p)
	}

	return 0, errors.New("buffer is closed")
}

func (b *buffer) Write(p []byte) (n int, err error) {
	if b.buffer != nil {
		b.buffer.Write(p)
	}

	return 0, errors.New("buffer is closed")
}

func newBuffer() io.ReadWriteCloser {
	return &buffer{
		buffer: bytes.NewBuffer(make([]byte, 0, MaxOutputCapacity)),
	}
}

type pipe struct {
	io.ReadCloser
	io.WriteCloser
}

func newPipe() io.ReadWriteCloser {
	p := &pipe{}
	p.ReadCloser, p.WriteCloser = io.Pipe()
	return p
}

func (p *pipe) Close() error {
	err := p.WriteCloser.Close()
	if err0 := p.ReadCloser.Close(); err0 != nil {
		if err != nil {
			err = fmt.Errorf("%s; %w", err, err0)
		} else {
			err = err0
		}
	}

	return err
}

var MaxOutputCapacity = 10 * 1024
