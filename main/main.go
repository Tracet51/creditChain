package main

import (
	"flag"

	"github.com/tracet51/creditChain/server"
	"github.com/tracet51/creditChain/voter"
)

func main() {

	var port = flag.String("port", "5051", "The port which to open the main TCP Connection")
	var peerAddress = flag.String("peer", "127.0.0.0:5052", "The inital Peer to connect to")

	flag.Parse()

	voter := voter.GetVoter()
	server := server.CreateServer(voter)
	server.RunServer()
}
