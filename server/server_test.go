package server

import (
	"net"
	"testing"
	"time"
)

type mockVoter struct {
}

func (voter mockVoter) Vote() {

}

func TestCreatesNewServer(t *testing.T) {
	port := "5051"
	ipAddress := "127.0.0.1:" + port
	server := CreateServer(mockVoter{}, port)
	if server.IPAddress != ipAddress {
		t.Errorf("Sent IP Address %v but created server with IP %v", ipAddress, server.IPAddress)
	}
}

func TestStreamsVoting(t *testing.T) {
	mockConnection := MockConnection{}
	broker := streamVoting(mockConnection)

	for i := 0; i < 3; i++ {
		generatedMessage := <-broker
		if generatedMessage == nil {
			t.Error("Unexpected nil message")
		}
		i++
	}
}

func TastStreamingReturnsMessage(t *testing.T) {
	mockConnection := MockConnection{}
	broker := streamVoting(mockConnection)
}

type MockConnection struct {
}

func (mock MockConnection) Read(b []byte) (n int, err error) {
	copy(b, "Hello world! \n")
	return 14, nil
}

func (mock MockConnection) Write(b []byte) (n int, err error) {
	return 4, nil
}

func (mock MockConnection) Close() error {
	return nil

}

func (mock MockConnection) LocalAddr() net.Addr {
	return nil
}

func (mock MockConnection) RemoteAddr() net.Addr {
	return nil
}

func (mock MockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (mock MockConnection) SetReadDeadline(t time.Time) error {
	return nil

}

func (mock MockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
