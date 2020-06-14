package tcpServer

import (
	"github.com/golang/mock/gomock"
	"net"
	"testing"
)

func TestConnectionMadeSavesTransport(t *testing.T) {
	t.Parallel()
	protocol := &TCPProtocol{}
	transport := MockTransport{}

	protocol.ConnectionMade(&transport)
	if protocol.connection == nil {
		t.Error("TCPProtocol should not be nil")
	}

	gomock.NewController(t)

}

func TestDataSendsMessage(t *testing.T) {
	t.Parallel()
	transport := &MockTransport{}
	protocol := &TCPProtocol{connection: transport}

	testCases := []struct {
		Name string
		Data []byte
	}{
		{Name: "Regular Test String", Data: []byte("This is data")},
		{Name: "Nil Data", Data: nil},
	}

	for _, testCase := range testCases {
		protocol.DataReceived(testCase.Data)
		if len(transport.Memory) != len(testCase.Data) {
			t.Errorf("%v failed, expected %v, got %v", testCase.Name, len(transport.Memory), len(testCase.Data))
		}
	}
}

func TestConnectionLostClosesTransport(t *testing.T) {
	t.Parallel()
	transport := &MockTransport{Calls: make(map[string]int, 0)}
	protocol := &TCPProtocol{connection: transport}

	err := protocol.ConnectionLost()

	if transport.Calls["Close"] != 1 || err != nil {
		t.Errorf("Close never called on connection")
	}
}

type MockTransport struct {
	Memory []byte
	Calls  map[string]int
	net.Conn
}

func (transport *MockTransport) Read(holder []byte) (bytesRead int, err error) {
	message := "Hello World!\n"
	copy(holder, []byte(message))
	return len(message), nil
}

func (transport *MockTransport) Write(b []byte) (n int, err error) {
	transport.Memory = make([]byte, len(b))
	copy(transport.Memory, b)
	return len(transport.Memory), nil
}

func (transport *MockTransport) Close() error {
	transport.Calls["Close"]++
	return nil

}

func (transport *MockTransport) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   []byte("127.0.0.1"),
		Port: 0,
		Zone: "",
	}
}
