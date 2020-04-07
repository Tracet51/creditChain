package main

import (
	"flag"
	"fmt"

	"github.com/tracet51/creditChain/server"
)

func main() {

	var port = flag.String("port", "5051", "The port which to open the main TCP Connection")
	flag.Parse()

	// voter := voter.GetVoter()
	del := delegate{}
	server := server.CreateServer(del, *port)
	server.RunServer()
}

type delegate struct {
}

func (delegate delegate) Vote() {
	fmt.Println("Voted!")
}
