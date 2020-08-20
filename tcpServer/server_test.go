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
	tcpServer, _, _ := createServerWithConnection(controller)

	if tcpServer == nil {
		t.Error("server should not be nil")
	}
}

func createServerWithConnection(controller *gomock.Controller) (*server, *MockProtocol, *MockConn) {
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

func createServerWithoutConnection(controller *gomock.Controller) (*server, *MockProtocol) {
	protocol := NewMockProtocol(controller)
	protocolFactory := func() Protocol {
		return protocol

	}

	connectionFactory := func() <-chan net.Conn {
		connectionChannel := make(chan net.Conn, 0)
		return connectionChannel
	}

	return NewServer(protocolFactory, connectionFactory), protocol
}

func TestServer_AcceptConnections(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()


	tcpServer, protocol, connection := createServerWithConnection(controller)

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

func TestServer_AcceptConnections_Closes(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tcpServer, protocol := createServerWithoutConnection(controller)
	cancelContext, cancelFunc := context.WithCancel(context.Background())

	protocol.EXPECT().ConnectionMade(gomock.Any()).MaxTimes(0)

	cancelFunc()

	tcpServer.AcceptConnections(cancelContext)
}


