package protocol

import (
	"errors"
	"io"
)

type Protocol interface {
	ConnectionMade(transport io.ReadWriteCloser) (err error)
	DataReceived(data []byte) (err error)
	ConnectionLost() (err error)
	Transport() io.ReadWriteCloser
}

// TCPProtocol ...
type TCPProtocol struct {
	transport io.ReadWriteCloser
}

// ConnectionMade ...
func (protocol *TCPProtocol) ConnectionMade(transport io.ReadWriteCloser) (err error) {

	if protocol.Transport == nil {
		err = errors.New("transports cannot be nil")
	}
	protocol.transport = transport

	return err
}

// DataReceived ...
func (protocol *TCPProtocol) DataReceived(data []byte) (err error) {

	protocol.transport.Write(data)
	return err
}

// ConnectionLost ...
func (protocol *TCPProtocol) ConnectionLost() (err error) {
	defer protocol.transport.Close()
	return err
}

// Transport ...
func (protocol *TCPProtocol) Transport() io.ReadWriteCloser {
	return protocol.transport
}
