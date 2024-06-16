package http

import (
	"io"

	"golang.org/x/text/transform"
)

func (r *Msg) writeClose(w io.WriteCloser) error {
	defer w.Close()
	return r.Write(w)
}

// Reader returns a Reader to read the formatted message.
func (r *Msg) Reader() io.Reader {
	reader, pipe := io.Pipe()
	go r.writeClose(pipe)
	return reader
}

// Transform returns the transformed message.
func (r *Msg) Transform(t transform.Transformer) (*Msg, error) {
	if t == nil {
		return r, nil
	}
	return ReadMsg(transform.NewReader(r.Reader(), t))
}

// Empty returns true if Msg has no valid headers.
func (r *Msg) Empty() bool {
	return len(r.Headers) == 0
}
