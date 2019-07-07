package server

import (
	cl "nmchat/client"
	msg "nmchat/message"
)

//AddNewClient - add client
func (s *Server) AddNewClient(cl *cl.Client) {
	s.addCliCh <- cl
}

//DelClient - del client
func (s *Server) DelClient(cl *cl.Client) {
	s.delCliCh <- cl
}

//SendMessage - send message to chat
func (s *Server) SendMessage(msg *msg.Message) {
	s.sendMsgCh <- msg
}

//Done - done channale
func (s *Server) Done() {
	s.doneCh <- true
}

//Err - error channel
func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) getNextClientID() int {
	id := s.clNextID
	s.clNextID++
	return id
}

func (s *Server) getActiveClientByID(id int) *cl.Client {
	if cl, ok := s.clients[id]; ok {
		return cl
	}
	return nil
}
