package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// connect to this socket
	conn, err := net.Dial("tcp", "127.0.0.1:5051")
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text)
		// message, _ := bufio.NewReader(conn).ReadString('\n')
		// fmt.Print("Message from server: " + message)
	}
}
