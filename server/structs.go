package server

import (
	"log"
	cl "nmchat/client"
	msg "nmchat/message"
)

type Server struct {
	messages  []msg.Message
	sendMsgCh chan *msg.Message
	clients   map[int]*cl.Client
	addCliCh  chan *cl.Client
	delCliCh  chan *cl.Client
	doneCh    chan bool
	errCh     chan error
}

func (s *Server) Start() {
	for {
		select {
		case client := <-s.addCliCh:
			s.clients[client.ID] = client

		case client := <-s.delCliCh:
			delete(s.clients, client.ID)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
