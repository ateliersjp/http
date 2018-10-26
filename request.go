package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

// Request represents an HTTP request to be sent.
type Request struct {
	Version    string
	Method     string
	Hostname   string
	Port       string
	RequestURI string
	Headers    *Headers
	Body       io.Reader

	conn net.Conn
}

// Close closes the request's connection.
func (r *Request) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

// URL fills in the Hostname, Port and RequestURI fields in r from a url string.
func (r *Request) URL(addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	r.Hostname = u.Hostname()
	r.Port = u.Port()
	r.RequestURI = u.RequestURI()

	return nil
}

// Fprint writes the formatted request to w.
func (r *Request) Fprint(w io.Writer) error {
	// Host header is written with the value extracted from the url.
	// Non-standard port is added to the host value.
	host := r.Hostname
	if r.Port != "" && r.Port != "80" {
		host += fmt.Sprintf(":%v", r.Port)
	}
	r.Headers.Add("Host", host)

	// Write request request line.
	path := r.RequestURI
	if path == "" {
		path = "/"
	}
	_, err := fmt.Fprintf(w, "%v %v HTTP/1.0\r\n", r.Method, path)
	if err != nil {
		return fmt.Errorf("could not write request line: %v", err)
	}

	// Write header lines.
	err = r.Headers.Fprint(w)
	if err != nil {
		return fmt.Errorf("could not write headers: %v", err)
	}

	// Write empty line to signal end of headers.
	_, err = fmt.Fprintf(w, "\r\n")
	if err != nil {
		return fmt.Errorf("could not write: %v", err)
	}

	// Write request body.
	if r.Body != nil {
		_, err := io.Copy(w, r.Body)
		if err != nil {
			return fmt.Errorf("could not write data: %v", err)
		}
	}

	return nil
}

// ReadRequest parses and reads a request from r.
func ReadRequest(r io.Reader) (*Request, error) {
	req := &Request{
		Headers: &Headers{},
	}

	// Reader is used to read request line by line.
	reader := bufio.NewReader(r)

	// Request body is read from the remaining data in the reader.
	req.Body = reader

	// Read first line of request (request line).
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("could not read request line: %v", err)
	}
	rl := strings.Split(strings.TrimSpace(line), " ")
	if len(rl) < 3 {
		return nil, fmt.Errorf("could not parse request line: %v", line)
	}
	req.Method = rl[0]
	req.RequestURI = rl[1]
	req.Version = rl[2]

	// Read and parse header lines.
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read header line: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		req.Headers.AddRaw(line)
	}

	return req, nil
}