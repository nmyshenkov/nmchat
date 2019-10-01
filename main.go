package main

import (
	"nmchat/server"
)

func main() {

	s := server.Init(":3333")
	s.Start()

}
