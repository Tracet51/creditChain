package main

import (
	"context"
	"flag"
	"github.com/tracet51/creditChain/tcpServer"
	"log"
	"net"
)

func main() {

	log.Println("Starting Credit Chain Server")

	var port = flag.String("port", "5051", "The port which to open the main TCP Connection")
	flag.Parse()

	connectionFactory := connectionFactoryBuilder("127.0.0.1:" + *port)
	server := tcpServer.NewServer(protocolFactory, connectionFactory)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	server.AcceptConnections(ctx)
}

func protocolFactory() tcpServer.Protocol {
	return &tcpServer.TCPProtocol{}
}

func connectionFactoryBuilder(ipAddress string) tcpServer.ConnectionFactory {
	listener := newListener(ipAddress)
	return func() <-chan net.Conn {
		connectionGenerator := make(chan net.Conn, 0)

		go func() {
			for {
				connection, err := listener.Accept()
				if err != nil {
					log.Fatal(err.Error())
				}
				connectionGenerator <-connection
			}
		}()
		return connectionGenerator
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
