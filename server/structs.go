package server

import (
	"log"
	msg "nmchat/message"
)

type Server struct {
	messages  []msg.Message
	sendMsgCh chan *msg.Message
	doneCh    chan bool
	errCh     chan error
}

func (s *Server) Start() {
	for {
		select {
		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
