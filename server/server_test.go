package server

import "testing"

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
