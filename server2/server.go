package server2

import (
	"github.com/tracet51/creditChain/protocol"
	"io"
	"log"
	"net"
)

func ListenForConnections(listener net.Listener) io.ReadWriteCloser {

	connection, err := listener.Accept()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Peer connected with address: " + connection.RemoteAddr().String())
	return connection
}

func NewListener(ipAddress string) net.Listener {
	listener, err := net.Listen("tcp", ipAddress)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Listening for connection on %v", listener.Addr().String())
	return listener
}

func AttachConnection(transport io.ReadWriteCloser, protocol protocol.Protocol) {
	protocol.ConnectionMade(transport)

}