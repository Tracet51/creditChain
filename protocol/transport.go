package protocol

import (
	"io"
)

type Transport interface {
	io.ReadWriteCloser
	Address() string
}

type TCPTransport struct {
	io.ReadWriteCloser
	address string
}

func NewTCPTransport(transporter io.ReadWriteCloser, address string) *TCPTransport {
	return &TCPTransport{
		ReadWriteCloser: transporter,
		address:         address,
	}
}

func (transport *TCPTransport) Address() string {
	return transport.address
}
