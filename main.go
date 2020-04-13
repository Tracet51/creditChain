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

	server := server2.NewServer("127.0.0.1:" + *port, func() protocol.Protocol {
		return &protocol.TCPProtocol{}
	})
	defer server.CloseConnections()
	server.AcceptConnections()
}

type delegate struct {
}

func (delegate delegate) Vote() {
	fmt.Println("Voted!")
}
