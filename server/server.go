package server

import (
	"bufio"
	"log"
	"net"

	"github.com/tracet51/creditChain/message"
	"github.com/tracet51/creditChain/voter"
)

const (
	serverDataFilename = "/tmp/serverMeta.txt"
)

var messenger = message.Messenger{
	InboundMessages:  make(chan message.Message),
	OutboundMessages: make(chan message.Message),
}

// Server handles all incoming and outgoing messages
type Server struct {
	IPAddress   string
	Connections map[string]net.Conn
	Voter       *voter.Voter

	listener net.Listener
}

// CreateServer makes a new server
func CreateServer(voter *voter.Voter) *Server {

	connections := make(map[string]net.Conn)
	server := &Server{
		IPAddress:   "127.0.0.1:5051",
		Connections: connections,
		Voter:       voter,
	}

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

	go server.acceptConnections()

	// Runs the server forever by blocking channels
	for {
		select {
		case inboundMessage := <-messenger.InboundMessages:
			// Determine the Message Payload Type
			for _, connection := range server.Connections {
				connection.Write(message.ToBytes(&inboundMessage))
			}
		case outboundMessage := <-messenger.OutboundMessages:
			log.Fatal(outboundMessage)
		}
	}
}

func registerOpcodes() {
	message.Register((*message.InitMessage)(nil), make(chan message.Message))
}

func (server *Server) acceptConnections() {
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
	server.Connections[newConnection.RemoteAddr().String()] = newConnection
	go server.readMessagesFromConnection(newConnection.RemoteAddr().String())
}

func (server *Server) readMessagesFromConnection(connectionKey string) {
	connection := server.Connections[connectionKey]
	reader := bufio.NewReader(connection)
	for {
		payload, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Client %s is gone\n", connection.RemoteAddr().String())
			connection.Close()
			delete(server.Connections, connection.RemoteAddr().String())
			break
		}
		log.Printf(string(payload))
		message.GetMessage(payload)

		messenger.InboundMessages <- message.MessageInt{}
	}
}
