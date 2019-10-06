package server

import (
	cl "nmchat/client"
	msg "nmchat/message"
)

// AddNewClient - add client
func (s *Server) AddNewClient(cl *cl.Client) {
	s.addCliCh <- cl
}

// DelClient - del client
func (s *Server) DelClient(cl *cl.Client) {
	s.delCliCh <- cl
}

// SendMessage - send message to chat
func (s *Server) SendMessage(msg *msg.Message) {
	s.sendMsgCh <- msg
}

// Done - done channel
func (s *Server) Done() {
	s.doneCh <- true
}

// Err - error channel
func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) getNextClientID() int {
	id := s.clNextID
	s.clNextID++
	return id
}

func (s *Server) getActiveClientByID(id int) *cl.Client {
	if client, ok := s.clients[id]; ok {
		return client
	}
	return nil
}

func (s *Server) checkExistNikname(cl *cl.Client, newName string) bool {
	for _, client := range s.clients {
		// if it current client - slip
		if client.ID == cl.ID {
			continue
		}

		if client.Name == newName {
			return true
		}
	}

	return false
}
