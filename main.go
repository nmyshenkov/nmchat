package main

import (
	"github.com/nmyshenkov/nmchat/server"
)

func main() {

	s := server.Init(":3333")
	s.Start()

}
