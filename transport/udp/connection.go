package udp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/g-harel/http/transport/connection"
)

var RouterPort = os.Getenv("ROUTER_PORT")

var _ connection.Connection = &Conn{}

type Conn struct {
	socket *Socket

	sequence    uint32
	peerAddress uint32
	peerPort    uint16

	sending   *bytes.Buffer
	receiving *Window

	// lock sync.Mutex
	// Chan chan Packet
	// Errs chan error

	// sendWindow [][]byte
	// sendChan   chan Packet

	// recvWindow [][]byte
	// recvChan   chan Packet
}

func NewConn(s *Socket, seq uint32, addr uint32, port uint16, window *Window) *Conn {
	c := &Conn{
		socket:      s,
		sequence:    seq,
		peerAddress: addr,
		peerPort:    port,
		sending:     &bytes.Buffer{},
		receiving:   window,
	}

	c.receiving.swallowAck = true
	c.receiving.errChan = make(chan error)
	go func() {
		for {
			log.Printf("error: %v", <-c.receiving.errChan)
		}
	}()

	return c
}

func (c *Conn) Commit() error {
	log.Printf("Comit()\n")
	c.sequence++

	p := &Packet{
		Sequence:    c.sequence,
		PeerAddress: c.peerAddress,
		PeerPort:    c.peerPort,
		Payload:     c.sending.Bytes(),
	}

	err := c.socket.Send(p, Timeout)
	if err != nil {
		return fmt.Errorf("sending packet: %v", err)
	}

	ackPacket, err := c.socket.Receive(Timeout)
	if err != nil {
		return fmt.Errorf("wait for ack packet: %v", err)
	}

	if ackPacket.Type != ACK {
		return fmt.Errorf("check ack packet: not ACK")
	}
	if ackPacket.Sequence != p.Sequence {
		return fmt.Errorf("check ack packet: sequence doesn't match")
	}

	c.sending.Reset()

	return nil
}

func (c *Conn) Read(b []byte) (int, error) {
	p, err := c.receiving.Read(Timeout)
	if err != nil {
		return 0, fmt.Errorf("read packet: %v", err)
	}

	if len(b) < len(p.Payload) {
		return 0, fmt.Errorf("read packet: read buffer too small")
	}

	copy(b, p.Payload)

	return len(p.Payload), io.EOF
}

func (c *Conn) Write(b []byte) (int, error) {
	return c.sending.Write(b)
}

func (c *Conn) Close() error {
	log.Printf("Close()\n")

	// TODO connection close handshake
	return nil
}
