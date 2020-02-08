package main

import (
	"github.com/tracet51/creditChain/server"
	"github.com/tracet51/creditChain/voter"
)



func main() {

	voter := voter.GetVoter()
	server := server.GetServer(voter)
	server.RunServer()
}