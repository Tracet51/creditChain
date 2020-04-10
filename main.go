package main

import (
	"flag"
	"fmt"
	"github.com/tracet51/creditChain/protocol"
	"github.com/tracet51/creditChain/server2"
	"log"
)

func main() {

	log.Println("Starting Credit Chain Server")

	var port = flag.String("port", "5051", "The port which to open the main TCP Connection")
	flag.Parse()

	listener := server2.NewListener("127.0.0.1:" + *port)
	defer listener.Close()
	for {
		connection := server2.ListenForConnections(listener)
		go server2.AttachConnection(connection, &protocol.TCPProtocol{})
	}
}

type delegate struct {
}

func (delegate delegate) Vote() {
	fmt.Println("Voted!")
}
