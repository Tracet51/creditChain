package tcpServer

import (
	"context"
	"net"
)

type ProtocolFactory func () Protocol
type ConnectionFactory func() <-chan net.Conn

type server struct {
	protocolFactory     ProtocolFactory
	connectionGenerator ConnectionFactory
}

func NewServer(protocolFactory ProtocolFactory, connectionFactory ConnectionFactory) *server {
	return &server{
		protocolFactory:     protocolFactory,
		connectionGenerator: connectionFactory,
	}
}

func (server *server) AcceptConnections(ctx context.Context) (err error) {
	connections := server.connectionGenerator()
	for {
		select {
		case <- ctx.Done():
			return ctx.Err()
		case connection := <-connections:
			protocolContext, _ := context.WithCancel(ctx)
			go InitiateCommunication(protocolContext, connection, server.protocolFactory())
		}
	}
}