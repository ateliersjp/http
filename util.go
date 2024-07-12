package http

import (
	"io"

	"golang.org/x/text/transform"
)

// Reader returns a Reader to read the formatted message.
func (r *Msg) Reader() io.Reader {
	reader, pipe := io.Pipe()
	go func() {
		defer pipe.Close()
		r.Write(pipe)
	}()
	return reader
}

// Transform returns the transformed message.
func (r *Msg) Transform(t transform.Transformer) (*Msg, error) {
	// no valid Transformer given.
	if t == nil {
		return r, nil
	}

	reader, pipe := io.Pipe()
	writer := transform.NewWriter(pipe, t)
	go func() {
		defer pipe.Close()
		r.Write(writer)
		writer.Close()
	}()
	return ReadMsg(reader)
}
