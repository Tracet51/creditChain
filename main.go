package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/google/uuid"
)

func main() {

	id, _ := uuid.NewUUID()
	node := &Node{
		ID:        id,
		QuorumSet: make(map[string]*Node)}

	server := &Server{
		IPAddress:         "127.0.0.1:5051",
		AllConnections:    make(map[net.Conn]int32),
		NewConnections:    make(chan net.Conn),
		DeadConnections:   make(chan net.Conn),
		InboundMessages:   make(chan Message),
		OutboundMessage:   make(chan Message),
		ConnectionCounter: 0,
		Node:              node}

	server.RunServer()
}

type Message struct {
	Payload string
	NodeID  string
}

func (message *Message) Write() []byte {
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err.Error())
	}
	return data
}
