package server2

import (
	"github.com/tracet51/creditChain/protocol"
	"io"
	"testing"
)

func TestAttachToConnection(t *testing.T) {
	var protocol protocol.Protocol
	protocol = &mockProtocol{}
	transport := &mockTransport{}
	AttachConnection(transport, protocol)
	if protocol.Transport() == nil {
		t.Errorf("Transport on Protocol was nil")
	}
}

type mockTransport struct {
	io.ReadWriteCloser
}

type mockProtocol struct {
	transport io.ReadWriteCloser
}

// ConnectionMade ...
func (protocol *mockProtocol) ConnectionMade(transport io.ReadWriteCloser) (err error) {
	protocol.transport = transport
	return nil
}

// DataReceived ...
func (protocol *mockProtocol) DataReceived(data []byte) (err error) {
	return nil
}

// ConnectionLost ...
func (protocol *mockProtocol) ConnectionLost() (err error) {
	return nil
}

func (protocol *mockProtocol) Transport() io.ReadWriteCloser {
	return protocol.transport
}

func TestNewListenerListensOnPort(t *testing.T) {
	ipAddress := "127.0.0.1:5051"
	listener := NewListener(ipAddress)
	listenerAddress := listener.Addr().String()
	if listenerAddress != ipAddress {
		t.Errorf("Expected Address %v but got %v", ipAddress, listenerAddress)
	}
}