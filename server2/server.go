package server2

import (
	"github.com/tracet51/creditChain/protocol"
	"log"
	"net"
)

type server struct {
	protocolFactory func () protocol.Protocol
	listener net.Listener
}

func NewServer(ipAddress string, protocolFactory func () protocol.Protocol) *server {
	return &server{
		protocolFactory: protocolFactory,
		listener:        newListener(ipAddress),
	}
}

func newListener(ipAddress string) net.Listener {
	listener, err := net.Listen("tcp", ipAddress)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%v: Listening for connections", listener.Addr().String())
	return listener
}

func (server *server) AcceptConnections() (err error) {
	for {
		transport := server.listenForConnections()
		go protocol.InitiateCommunication(transport, server.protocolFactory())
	}

}

// listenForConnections ...
func (server *server) listenForConnections() protocol.Transport {

	connection, err := server.listener.Accept()
	if err != nil {
		log.Fatal(err.Error())
	}
	address := connection.RemoteAddr().String()
	return protocol.NewTCPTransport(connection, address)
}

func (server *server) CloseConnections() {
	defer server.listener.Close()
}