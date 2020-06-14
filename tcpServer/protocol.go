//go:generate mockgen -destination mocks.go . Protocol

package tcpServer

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Protocol interface {
	ConnectionMade(conn net.Conn) (err error)
	DataReceived(data []byte) (err error)
	ConnectionLost() (err error)
	Transport() net.Conn
}

// TCPProtocol ...
type TCPProtocol struct {
	connection net.Conn
}

// ConnectionMade ...
func (protocol *TCPProtocol) ConnectionMade(conn net.Conn) (err error) {
	log.Printf("%v: Peer Connected", conn.RemoteAddr().String())
	protocol.connection = conn
	return err
}

// DataReceived ...
func (protocol *TCPProtocol) DataReceived(data []byte) (err error) {

	message := strings.TrimRight(string(data), "\r\n")
	log.Printf("%v: Sent: Local Message: %v ", protocol.address(), message)

	protocol.connection.Write(data)
	log.Printf("Local: Sent: %v Message: %v ", protocol.address(), message)

	return err
}

// ConnectionLost ...
func (protocol *TCPProtocol) ConnectionLost() (err error) {
	defer protocol.connection.Close()
	log.Printf("%v: Disconnected", protocol.address())
	return err
}

func (protocol TCPProtocol) address() string {
	return protocol.Transport().RemoteAddr().String()
}

// Transport ...
func (protocol *TCPProtocol) Transport() net.Conn {
	return protocol.connection
}

// InitiateCommunication ...
func InitiateCommunication(conn net.Conn, protocol Protocol) {
	protocol.ConnectionMade(conn)
	reader := bufio.NewReader(protocol.Transport())
	for {
		payload, err := reader.ReadBytes('\n')
		if err != nil {
			protocol.ConnectionLost()
			break
		} else {
			protocol.DataReceived(payload)
		}
	}
}
