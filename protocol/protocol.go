package protocol

import (
	"log"
	"strings"
)

type Protocol interface {
	ConnectionMade(transport Transport) (err error)
	DataReceived(data []byte) (err error)
	ConnectionLost() (err error)
	Transport() Transport
}

// TCPProtocol ...
type TCPProtocol struct {
	transport Transport
}

// ConnectionMade ...
func (protocol *TCPProtocol) ConnectionMade(transport Transport) (err error) {
	log.Printf("%v: Peer Connected", transport.Address())
	protocol.transport = transport
	return err
}

// DataReceived ...
func (protocol *TCPProtocol) DataReceived(data []byte) (err error) {

	message := strings.TrimRight(string(data), "\r\n")
	log.Printf("%v: Sent: Local Message: %v ", protocol.transport.Address(), message)

	protocol.transport.Write(data)
	log.Printf("Local: Sent: %v Message: %v ", protocol.transport.Address(), message)

	return err
}

// ConnectionLost ...
func (protocol *TCPProtocol) ConnectionLost() (err error) {
	defer protocol.transport.Close()
	log.Printf("%v: Disconnected", protocol.transport.Address())
	return err
}

// Transport ...
func (protocol *TCPProtocol) Transport() Transport {
	return protocol.transport
}
