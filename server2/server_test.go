package server2

import (
	"github.com/tracet51/creditChain/protocol"
	"net"
	"testing"
)

func TestNewListenerListensOnPort(t *testing.T) {
	ipAddress := "127.0.0.1:5051"
	listener := NewListener(ipAddress)
	listenerAddress := listener.Addr().String()
	if listenerAddress != ipAddress {
		t.Errorf("Expected Address %v but got %v", ipAddress, listenerAddress)
	}
}

func TestListenForConnections(t *testing.T) {
	testCases := []struct{
		name       string
		listener net.Listener
	}{
		{name: "Correct Connection", listener: &successfulListener{}},
	}

	for _, testCase := range testCases {
		connection := ListenForConnections(testCase.listener)
		if connection == nil {
			t.Errorf("Connection should not be nil")
		}
	}

}

func TestListenForConnectionsPanics(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ListenForConnections should have panicked")
			}
		}()
		// This function should cause a panic
		ListenForConnections(&errorListener{})
	}()
}



func TestAttachToConnection(t *testing.T) {
	var protocol protocol.Protocol
	protocol = &mockProtocol{}
	transport := &mockTransport{}
	InitiateCommunication(transport, protocol)
	if protocol.Transport() == nil {
		t.Errorf("Transport on Protocol was nil")
	}
}


