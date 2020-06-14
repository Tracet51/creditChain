package tcpServer

import (
	"errors"
	"net"
)

type mockListener struct {
	net.Listener
}

func (listener *mockListener) Close() error {
	return nil
}

func (listener *mockListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   nil,
		Port: 0,
		Zone: "",
	}
}

type errorListener struct {
	mockListener
}

func (listener *errorListener) Accept() (net.Conn, error) {
	return &net.TCPConn{}, errors.New("fake error")
}

type successfulListener struct {
	mockListener
}

func (listener *successfulListener) Accept() (net.Conn, error) {
	return &mockConnection{}, nil
}

type mockConnection struct {
	net.Conn
}

func (connection *mockConnection) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP: []byte("127.0.0.1"),
		Port: 5051,
		Zone: "",
	}
}


type mockTransport struct {
	net.Conn
}

type mockProtocol struct {
	transport net.Conn
}

// ConnectionMade ...
func (protocol *mockProtocol) ConnectionMade(transport net.Conn) (err error) {
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

func (protocol *mockProtocol) Transport() net.Conn {
	return protocol.transport
}

