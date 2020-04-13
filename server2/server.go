package server2

import (
	"bufio"
	"github.com/tracet51/creditChain/protocol"
	"log"
	"net"
)

// NewListener ...
func NewListener(ipAddress string) net.Listener {
	listener, err := net.Listen("tcp", ipAddress)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%v: Listening for connections", listener.Addr().String())
	return listener
}

// ListenForConnections ...
func ListenForConnections(listener net.Listener) protocol.Transport {

	connection, err := listener.Accept()
	if err != nil {
		log.Fatal(err.Error())
	}
	address := connection.RemoteAddr().String()
	return protocol.NewTCPTransport(connection, address)
}

// InitiateCommunication ...
func InitiateCommunication(transport protocol.Transport, protocol protocol.Protocol) {
	protocol.ConnectionMade(transport)
	reader := bufio.NewReader(protocol.Transport())
	for {
		payload, err := reader.ReadBytes('\n')
		if err != nil {
			protocol.ConnectionLost()
			break
		} else {
			protocol.DataReceived(payload)
		}
	}
}