package main

import (
	"nmchat/server"
)

func main() {

	server := server.Init("")
	server.LisenTelnet()

	//server.Start()

}
