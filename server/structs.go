package server

import (
	"log"
	cl "nmchat/client"
	msg "nmchat/message"
)

type Server struct {
	messages  []*msg.Message
	sendMsgCh chan *msg.Message
	clients   map[int]*cl.Client
	addCliCh  chan *cl.Client
	delCliCh  chan *cl.Client
	doneCh    chan bool
	errCh     chan error
}

//Init chat server.
func Init(pattern string) *Server {
	messages := []*msg.Message{}
	sendMsgCh := make(chan *msg.Message)
	clients := make(map[int]*cl.Client)
	addCh := make(chan *cl.Client)
	delCh := make(chan *cl.Client)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		messages,
		sendMsgCh,
		clients,
		addCh,
		delCh,
		doneCh,
		errCh,
	}
}

func (s *Server) Start() {
	for {
		select {
		case client := <-s.addCliCh:
			s.clients[client.ID] = client

		case client := <-s.delCliCh:
			delete(s.clients, client.ID)

		case msg := <-s.sendMsgCh:
			log.Println("Send:", msg)
			s.messages = append(s.messages, msg)
			s.SendMessage(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s *Server) AddNewClient(cl *cl.Client) {
	s.addCliCh <- cl
}

func (s *Server) DelClient(cl *cl.Client) {
	s.delCliCh <- cl
}

func (s *Server) SendMessage(msg *msg.Message) {
	s.sendMsgCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}
