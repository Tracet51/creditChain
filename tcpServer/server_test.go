package tcpServer

import (
	"context"
	"github.com/golang/mock/gomock"
	"net"
	"sync"
	"testing"
)

func TestNewServer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	tcpServer, _, _ := createServer(controller)

	if tcpServer == nil {
		t.Error("server should not be nil")
	}
}

func createServer(controller *gomock.Controller) (*server, *MockProtocol, *MockConn) {
	protocol := NewMockProtocol(controller)
	protocolFactory := func() Protocol {
		return protocol

	}

	connection := NewMockConn(controller)
	connectionFactory := func() <-chan net.Conn {
		connectionChannel := make(chan net.Conn, 1)
		connectionChannel <- connection
		return connectionChannel
	}

	return NewServer(protocolFactory, connectionFactory), protocol, connection
}

func TestServer_AcceptConnections(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	tcpServer, protocol, connection := createServer(controller)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	defer waitGroup.Wait()

	ctx := context.Background()
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	protocol.
		EXPECT().
		ConnectionMade(gomock.Eq(connection)).
		DoAndReturn(func(args ...interface{}) error {
			cancelFunc()
			waitGroup.Done()
			return nil
		})

	protocol.
		EXPECT().
		Transport().
		AnyTimes()

	tcpServer.AcceptConnections(cancelCtx)
}


