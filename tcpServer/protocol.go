//go:generate mockgen -destination=mock_protocol.go -package=tcpServer . Protocol
//go:generate mockgen -destination=mock_net.go -package=tcpServer net Conn,Addr

package tcpServer

import (
	"bufio"
	"context"
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
	log.Printf("%v -> Local: %v ", protocol.address(), message)

	protocol.connection.Write(data)
	log.Printf("Local -> %v: %v ", protocol.address(), message)

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
func InitiateCommunication(ctx context.Context, conn net.Conn, protocol Protocol) {
	protocol.ConnectionMade(conn)
	ctx, cancelFunc := context.WithCancel(ctx)
	reader := bufio.NewReader(protocol.Transport())
	for {
		select {
		case <-ctx.Done():
			return
		default:
			payload, err := reader.ReadBytes('\n')
			if err != nil {
				protocol.ConnectionLost()
				cancelFunc()
			} else {
				protocol.DataReceived(payload)
			}
		}
	}
}
