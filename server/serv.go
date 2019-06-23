package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	cl "nmchat/client"
	msg "nmchat/message"
	"strings"
	"time"
)

//Init chat server.
func Init(addr string) *Server {
	messages := []*msg.Message{}
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
		addCh,
		delCh,
		doneCh,
		errCh,
		time.Minute,
	}
}

func (s *Server) lisenChannels() {

	for {
		select {
		case client := <-s.addCliCh:
			s.clients[client.ID] = client

		case client := <-s.delCliCh:
			delete(s.clients, client.ID)

		case msg := <-s.sendMsgCh:
			log.Println("Send:", msg)
			s.messages = append(s.messages, msg)
			fmt.Printf("%+v\n", s.messages)
			//s.SendMessage(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			log.Println("Server stoped")
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

func (s *Server) Start() {

	defer func() {
		s.Done()
	}()

	addr := s.Addr
	if addr == "" {
		addr = ":3333"
	}

	go s.lisenChannels()

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
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) error {
	defer func() {
		log.Printf("Closing connection from %v", conn.RemoteAddr())
		conn.Close()
	}()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	scanr := bufio.NewScanner(r)
	for {
		scanned := scanr.Scan()
		if !scanned {
			if err := scanr.Err(); err != nil {
				log.Printf("%v(%v)", err, conn.RemoteAddr())
				return err
			}
			break
		}
		text := scanr.Text()

		switch text {
		case "":
			continue
		case "!exit":
			log.Printf("Client %v decited to exit", conn.RemoteAddr())
			conn.Close()
			return nil
		case "!help":
			text = "Commands:\n\t !exit - exit from chat\n\t !help - print info about command"
		default:
			msg := msg.Message{
				Body: scanr.Text(),
			}
			s.SendMessage(&msg)
			text = strings.ToUpper(text)
		}

		w.WriteString(text + "\n")
		w.Flush()
	}
	return nil
}
