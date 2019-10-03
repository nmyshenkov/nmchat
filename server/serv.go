package server

import (
	"bufio"
	"log"
	"net"
	cl "nmchat/client"
	msg "nmchat/message"
	"strconv"
	"time"
)

// Init chat server.
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
					// TODO:: check that client isActive
					client.IncomingChan <- message
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

	// start to lisen channels
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
		cl := cl.NewClient(s.getNextClientID())
		s.AddNewClient(cl)
		go s.handleClient(conn, cl)
	}
}

func (s *Server) handleClient(conn net.Conn, client *cl.Client) error {
	defer func() {
		log.Printf("Closing connection from %v", conn.RemoteAddr())
		s.DelClient(client)
		conn.Close()
	}()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	w.WriteString("Write !help to show commands\n Your nikname: " + client.Name + "\n")
	w.Flush()

	// listen incoming messages
	go func() {
		for {
			select {
			case message := <-client.IncomingChan:
				w.Write(message.FromByteMessage("\n"))
				w.Flush()
			}
		}
	}()

	scanner := bufio.NewScanner(r)
	for {
		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				log.Printf("%v(%v)", err, conn.RemoteAddr())
				return err
			}
			break
		}
		text := scanner.Text()

		switch text {
		case "":
			continue
		case "!exit":
			log.Printf("Client %s decited to exit", conn.RemoteAddr().String())
			return nil
		case "!help":
			text = HELP
		case "!list":
			text = "Users in the chat:\n"
			for _, client := range s.clients {
				text += "\t" + client.Name + " with ID: " + client.GetTextID() + "\n"
			}
		default:
			// TODO: will think about message send design
			// rewrite to function with error interface
			isMessage, cmd, arg := msg.TrimText(text)
			switch {
			case isMessage:
				id, _ := strconv.Atoi(cmd)
				c := s.getActiveClientByID(id)
				if c != nil {
					message := client.NewMessage(c, arg)
					s.SendMessage(&message)
					// clear message
					text = ""
				} else {
					text = "Client not found!!"
				}
			case cmd == "!name":
				client.Name = arg
				text = "Name changed!"
			default:
				message := client.NewBroadcastMessage(text)
				s.SendMessage(&message)
				// clear message
				text = ""
			}
		}

		w.WriteString(text + "\n")
		w.Flush()
	}
	return nil
}
