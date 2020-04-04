package protocol

import (
	"testing"
)

func TestConnectionMadeSavesTransport(t *testing.T) {
	protocol := &Protocol{}
	transport := mockTransport{}

	protocol.ConnectionMade(&transport)
	if protocol.Transport == nil {
		t.Error("Protocol should not be nil")
	}
}

func TestDataSendsMessage(t *testing.T) {
	transport := &mockTransport{}
	protocol := &Protocol{Transport: transport}

	testCases := []struct {
		Name string
		Data []byte
	}{
		{Name: "Regular Test String", Data: []byte("This is data")},
		{Name: "Nil Data", Data: nil},
	}

	for _, testCase := range testCases {
		protocol.DataRecieved(testCase.Data)
		if len(transport.Memory) != len(testCase.Data) {
			t.Errorf("%v failed, expected %v, got %v", testCase.Name, len(transport.Memory), len(testCase.Data))
		}
	}
}

func TestConnectionLostClosesTransport(t *testing.T) {
	transport := &mockTransport{Calls: make(map[string]int, 0)}
	protocol := &Protocol{Transport: transport}

	err := protocol.ConnectionLost()

	if transport.Calls["Close"] != 1 || err != nil {
		t.Errorf("Close never called on transport")
	}
}

type mockTransport struct {
	Memory []byte
	Calls  map[string]int
}

func (transport *mockTransport) Read(holder []byte) (bytesRead int, err error) {
	message := "Hello World!\n"
	copy(holder, []byte(message))
	return len(message), nil
}

func (transport *mockTransport) Write(b []byte) (n int, err error) {
	transport.Memory = make([]byte, len(b))
	copy(transport.Memory, b)
	return len(transport.Memory), nil
}

func (transport *mockTransport) Close() error {
	transport.Calls["Close"]++
	return nil

}
