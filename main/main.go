package main


func main() {

	voter := GetVoter()
	server := GetServer(voter)
	server.RunServer()
}