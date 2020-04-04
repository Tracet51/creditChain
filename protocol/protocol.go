package protocol

import (
	"errors"
	"io"
)

// Protocol ...
type Protocol struct {
	Transport io.ReadWriteCloser
}

// ConnectionMade ...
func (protocol *Protocol) ConnectionMade(transport io.ReadWriteCloser) (err error) {

	if protocol.Transport == nil {
		err = errors.New("Transports cannot be nil")
	}
	protocol.Transport = transport

	return err
}

// DataRecieved ...
func (protocol *Protocol) DataRecieved(data []byte) (err error) {

	protocol.Transport.Write(data)
	return err
}

// ConnectionLost ...
func (protocol *Protocol) ConnectionLost() (err error) {
	defer protocol.Transport.Close()
	return err
}
