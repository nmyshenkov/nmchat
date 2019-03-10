package main

import (
	"log"
)

type Message struct {
	From string
	To   string
	Body string
}

func (m *Message) ToString() string {
	return m.From + ": " + m.Body
}

type Server struct {
	messages  []Message
	sendMsgCh chan *Message
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

func main() {

}
