package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"github.com/tracet51/creditChain/message"
)

const (
	serverDataFilename = "/tmp/serverMeta.txt"
)

// Server handles all incoming and outgoing messages
type Server struct {
	IPAddress       string
	AllConnections  map[string]net.Conn
	NewConnections  chan net.Conn
	DeadConnections chan net.Conn
	InboundMessages chan message.Message
	OutboundMessage chan message.Message
	Voter            *Voter

	listener net.Listener
}

// GetServer makes a new server
func GetServer(voter *Voter) *Server {
	
	allConnections := make(map[string]net.Conn)
	server := &Server{
		IPAddress:       "127.0.0.1:5051",
		AllConnections:  allConnections,
		NewConnections:  make(chan net.Conn),
		DeadConnections: make(chan net.Conn),
		InboundMessages: make(chan message.Message),
		OutboundMessage: make(chan message.Message),
		Voter:            voter}

	return server
}

// RunServer spins up a new TCP server
func (server *Server) RunServer() {

	registerOpcodes()

	var err error
	server.listener, err = net.Listen("tcp", server.IPAddress)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer server.listener.Close()
	log.Println("TCP Server Started at: " + server.IPAddress)

	go server.acceptNewConnections()

	// Runs the server forever by blocking channels
	for {
		select {
		// We get a new connection
		case newConnection := <-server.NewConnections:
			server.handleNewConnections(newConnection)
		case inboundMessage := <-server.InboundMessages:
			// Determine the Message Payload Type
			for _, connection := range server.AllConnections {
				connection.Write(message.ToBytes(&inboundMessage))
			}
		case outboundMessage := <-server.OutboundMessage:
			log.Fatal(outboundMessage)
		case deadConnection := <-server.DeadConnections:
			log.Printf("Client %s is gone\n", deadConnection.RemoteAddr().String())
			deadConnection.Close()
			delete(server.AllConnections, deadConnection.RemoteAddr().String())
		}
	}
}

func registerOpcodes() {
	message.Register((*message.InitMessage)(nil), make(chan message.Message))
}

func (server *Server) acceptNewConnections() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Peer connected with address: " + connection.RemoteAddr().String())
		server.handleNewConnections(connection)
	}
}

func (server *Server) handleNewConnections(newConnection net.Conn) {
	// Add the connection to the list of connections
	server.AllConnections[newConnection.RemoteAddr().String()] = newConnection
	go server.readMessagesFromConnection(newConnection.RemoteAddr().String())
}

func (server *Server) readMessagesFromConnection(connectionKey string) {
	connection := server.AllConnections[connectionKey]
	reader := bufio.NewReader(connection)
	for {
		payload, err := reader.ReadBytes('\n')
		if err != nil {
			server.DeadConnections <- connection
			break
		}
		fmt.Println(byte(0))
		log.Printf(string(payload))
		message.GetMessage(payload)
	}
}