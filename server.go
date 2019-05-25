package main

import (
	"bufio"
	"log"
	"net"
)

// Server handles all incoming and outgoing messages
type Server struct {
	IPAddress         string
	AllConnections    map[net.Conn]int32
	NewConnections    chan net.Conn
	DeadConnections   chan net.Conn
	InboundMessages   chan Message
	OutboundMessage   chan Message
	ConnectionCounter int32
	Node              *Node
}

// RunServer spins up a new TCP server
func (server *Server) RunServer() {

	listener, err := net.Listen("tcp", server.IPAddress)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("TCP Server Started at: " + server.IPAddress)

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				log.Fatal(err.Error())
			}

			server.NewConnections <- connection
		}
	}()

	for {
		select {
		// We get a new connection
		case connection := <-server.NewConnections:

			// Add the connection to the list of connections
			server.AllConnections[connection] = server.ConnectionCounter
			server.ConnectionCounter++

			// Read messages from connection
			go func(connection net.Conn, connectionCounter int32) {

				reader := bufio.NewReader(connection)
				for {
					payload, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					message := Message{Payload: payload, NodeID: connection.RemoteAddr().String()}
					server.InboundMessages <- message
				}
				server.DeadConnections <- connection
			}(connection, server.ConnectionCounter)

		case inboundMessage := <-server.InboundMessages:
			// Determine the Message Payload Type
			log.Fatal(inboundMessage)
		case outboundMessage := <-server.OutboundMessage:
			log.Fatal(outboundMessage)
		case deadConnection := <-server.DeadConnections:
			log.Printf("Client %v is gone\n", server.AllConnections[deadConnection])
			deadConnection.Close()
			delete(server.AllConnections, deadConnection)
		}
	}
}
