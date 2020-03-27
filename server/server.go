package server

import (
	"bufio"
	"io"
	"log"
	"net"

	"github.com/tracet51/creditChain/message"
)

const (
	serverDataFilename = "/tmp/serverMeta.txt"
)

var messenger = message.Messenger{
	InboundMessages:  make(chan message.Message),
	OutboundMessages: make(chan message.Message),
}

type delegate interface {
	Vote()
}

// Server handles all incoming and outgoing messages
type Server struct {
	IPAddress   string
	Connections map[string]net.Conn
	Delegate    delegate

	listener net.Listener
}

// CreateServer makes a new server
func CreateServer(delegate delegate, port string) *Server {

	connections := make(map[string]net.Conn)
	server := &Server{
		IPAddress:   "127.0.0.1:" + port,
		Connections: connections,
		Delegate:    delegate,
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

	go server.respondToMessages()

	for {
		connection, err := server.listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Peer connected with address: " + connection.RemoteAddr().String())
		server.Connections[connection.RemoteAddr().String()] = connection
		go server.readMessagesFromConnection(connection.RemoteAddr().String())
	}
}

func registerOpcodes() {
	message.Register((*message.InitMessage)(nil), make(chan message.Message))
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
			return
		}
		log.Printf(string(payload))
		message.GetMessage(payload)

		messenger.InboundMessages <- message.MessageInt{}
	}
}

func (server *Server) respondToMessages() {
	for {
		select {
		case inboundMessage := <-messenger.InboundMessages:
			log.Println(inboundMessage)
		case outboundMessage := <-messenger.OutboundMessages:
			log.Fatal(outboundMessage.(message.InitMessage))
		}
	}
}

func streamVoting(connection io.ReadWriter) chan message.Message {
	messageBroker := make(chan message.Message)
	reader := bufio.NewReader(connection)

	go func() {
		for {
			payload, err := reader.ReadBytes('\n')
			if err != nil {
				panic("AHHH 4")
			}

			messageBroker <- payload
		}
	}()

	return messageBroker
}
