package main

import (
	"nmchat/server"
)

func main() {

	server := server.Init(":3333")
	server.Start()

}
