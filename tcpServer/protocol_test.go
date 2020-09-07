package tcpServer

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestConnectionMadeSavesTransport(t *testing.T) {
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

func TestDataSendsMessage(t *testing.T) {
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

		protocol.DataReceived(testCase.Data)
	}
}

func TestConnectionLostClosesTransport(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()
	transport := NewMockConn(controller)
	protocol := &TCPProtocol{connection: transport}

	err := protocol.ConnectionLost()

	if err != nil {
		controller.T.Fatalf("Error")
	}

}
