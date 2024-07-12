package http

import (
	"net"
)

// Send sends an http request.
func (req *Msg) Send(conn net.Conn) (*Msg, error) {
	// Write request to connection.
	err := req.Write(conn)
	if err != nil {
		return nil, err
	}

	// Read and parse response from connection.
	res, err := ReadMsg(conn)
	if err != nil {
		return nil, err
	}

	return res, nil
}
