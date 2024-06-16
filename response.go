package http

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Msg struct {
	Headers    []string
	Body       io.Reader
}

// Write writes the formatted message to w.
func (r *Msg) Write(w io.Writer) error {
	// Write header lines.
	for _, line := range r.Headers {
		_, err := fmt.Fprintf(w, "%v\r\n", line)
		if err != nil {
			return err
		}
	}

	// Write empty line to signal end of headers.
	_, err := fmt.Fprintf(w, "\r\n")
	if err != nil {
		return err
	}

	// Write message body.
	if r.Body != nil {
		_, err := io.Copy(w, r.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadMsg parses and reads a message from r.
func ReadMsg(r io.Reader) (*Msg, error) {
	msg := &Msg{}

	// Reader is used to read message line by line.
	reader := bufio.NewReader(r)

	// Body is read from the remaining data in the reader.
	msg.Body = reader

	// Read and parse header lines.
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		msg.Headers = append(msg.Headers, line)
	}

	return msg, nil
}
