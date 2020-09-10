package tcpServer

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestTCPProtocol_ConnectionMade_SavesTransport(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	transport := NewMockConn(controller)
	protocol := &TCPProtocol{connection: transport}

	address := NewMockAddr(controller)
	address.EXPECT().String().Return("Hello").MaxTimes(1)
	transport.EXPECT().RemoteAddr().MaxTimes(1).Return(address)

	protocol.ConnectionMade(transport)
	if protocol.connection == nil {
		t.Error("TCPProtocol should not be nil")
	}
}

func TestTCPProtocol_DataReceived_SendsMessage(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	testCases := []struct {
		Name string
		Data []byte
	}{
		{Name: "Regular Test String", Data: []byte("This is data")},
		{Name: "Nil Data", Data: nil},
	}

	for _, testCase := range testCases {
		transport := NewMockConn(controller)
		protocol := &TCPProtocol{connection: transport}

		transport.EXPECT().
			Write(gomock.Eq(testCase.Data)).
			MaxTimes(1)

		mockAddress := NewMockAddr(controller)
		address := "127.0.0.1"
		gomock.InOrder(
			mockAddress.EXPECT().String().Return(address).MaxTimes(1),
			mockAddress.EXPECT().String().Return(address).MaxTimes(1),
		)
		transport.EXPECT().RemoteAddr().Return(mockAddress).MaxTimes(2)
		transport.EXPECT().Write(gomock.Eq(testCase.Data)).MaxTimes(1)

		protocol.DataReceived(testCase.Data)
	}
}

func TestTCPProtocol_ConnectionLost_ClosesTransport(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	address := NewMockAddr(controller)
	address.EXPECT().String().Return("127.0.0.1").Times(1)

	transport := NewMockConn(controller)
	transport.EXPECT().RemoteAddr().Times(1).Return(address)
	transport.EXPECT().Close().Times(1)

	protocol := &TCPProtocol{connection: transport}

	err := protocol.ConnectionLost()

	if err != nil {
		controller.T.Fatalf("Error")
	}

}

func TestInitiateCommunication_ReturnsWhenContextDone(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	connection := NewMockConn(controller)

	protocol := NewMockProtocol(controller)
	gomock.InOrder(
		protocol.EXPECT().ConnectionMade(gomock.Eq(connection)).Times(1),
		protocol.EXPECT().Transport().Return(connection).Do(func() {
			cancel()
		}),
	)

	InitiateCommunication(ctx, connection, protocol)
}

func TestInitiateCommunication_ReadsEndOfPayload(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx, _ := context.WithCancel(context.Background())
	data := []byte("Hello World\n")
	connection := NewMockConn(controller)
	gomock.InOrder(
		connection.EXPECT().Read(gomock.Any()).Return(len(data), nil).Do(func(buf []byte) {
			copy(buf, data)
		}).Times(1),
		connection.EXPECT().Read(gomock.Any()).Return(0, errors.New("test error")).Times(1).Do(func(args ...interface{}) {
		}),
	)

	protocol := NewMockProtocol(controller)
	gomock.InOrder(
		protocol.EXPECT().ConnectionMade(gomock.Eq(connection)).Times(1),
		protocol.EXPECT().Transport().Return(connection),
		protocol.EXPECT().DataReceived(gomock.Any()).Times(1),
		protocol.EXPECT().ConnectionLost().Times(1),
	)

	InitiateCommunication(ctx, connection, protocol)
}
