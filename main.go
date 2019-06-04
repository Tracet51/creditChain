package main

func main() {

	node := CreateNewNode()
	server := CreateNewServer(node)
	server.RunServer()
}
