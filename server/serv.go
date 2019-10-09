package server

import (
	"log"
	"net"
	cl "nmchat/client"
	msg "nmchat/message"
	"time"
)

// Init chat server.
func Init(addr string) *Server {
	var messages []*msg.Message
	sendMsgCh := make(chan *msg.Message)
	clients := make(map[int]*cl.Client)
	addCh := make(chan *cl.Client)
	delCh := make(chan *cl.Client)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		addr,
		messages,
		sendMsgCh,
		clients,
		1,
		addCh,
		delCh,
		doneCh,
		errCh,
		time.Minute,
	}
}

func (s *Server) listenChannels() {

	for {
		select {
		case client := <-s.addCliCh:
			s.clients[client.ID] = client

		case client := <-s.delCliCh:
			delete(s.clients, client.ID)

		case message := <-s.sendMsgCh:
			log.Println("Send:", message)
			s.messages = append(s.messages, message)
			switch {
			case message.To > 0:
				c := s.getActiveClientByID(message.To)
				if c != nil {
					c.IncomingChan <- message
				}
			default:
				for _, client := range s.clients {
					if s.clientIsActive(client.ID) {
						client.IncomingChan <- message
					}
				}
			}

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			log.Println("Server stoped")
			return
		}
	}
}

// Start - start server
func (s *Server) Start() {

	defer func() {
		s.Done()
	}()

	addr := s.Addr
	if addr == "" {
		addr = ":3333"
	}

	// start to listen channels
	go s.listenChannels()

	log.Printf("Starting server on %v\n", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Error accepting connection %v.", err)
		return
	}
	defer listener.Close()
	for {
		newConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection %v.", err)
			continue
		}
		conn := &Conn{
			Conn:        newConn,
			IdleTimeout: s.IdleTimeout,
		}
		log.Printf("New connection from %v.", conn.RemoteAddr())
		conn.SetDeadline(time.Now().Add(conn.IdleTimeout))
		client := cl.NewClient(s.getNextClientID())
		s.AddNewClient(client)
		go s.handleClient(conn, client)
	}
}
