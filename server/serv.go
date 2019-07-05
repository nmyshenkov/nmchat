package server

import (
	"bufio"
	"log"
	"net"
	cl "nmchat/client"
	msg "nmchat/message"
	"strconv"
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
		1,
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
			if msg.To > 0 {
				if cl, ok := s.clients[msg.To]; ok {
					cl.IncomingChan <- msg
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

//Start - start server
func (s *Server) Start() {

	defer func() {
		s.Done()
	}()

	addr := s.Addr
	if addr == "" {
		addr = ":3333"
	}

	//start to lisen channels
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
		cl := cl.NewClient(s.getNextClientID())
		s.AddNewClient(cl)
		go s.handleClient(conn, cl)
	}
}

func (s *Server) handleClient(conn net.Conn, client *cl.Client) error {
	defer func() {
		log.Printf("Closing connection from %v", conn.RemoteAddr())
		conn.Close()
	}()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	w.WriteString("Write !help to show commands\n Your nikname: " + client.Name + "\n")
	w.Flush()

	//lisen incomming messages
	go func() {
		for {
			select {
			case message := <-client.IncomingChan:
				w.Write(message.FromByteMessage())
				w.Flush()
			}
		}
	}()

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
			text = "Commands:\n\t !exit - exit from chat\n\t !help - print info about command\n\t !list - print users in the chat"
		case "!list":
			text = "Users in the chat:\n"
			for _, client := range s.clients {
				text += "\t" + client.Name + " with ID: " + client.GetTextID() + "\n"
			}
		default:

			//TODO: will think about message send design
			if text[0] == 47 {
				texts := strings.SplitN(text, " ", 2)

				msgTo := texts[0]
				msgBody := texts[1]

				msgTo = strings.Replace(msgTo, "/", "", 1)
				id, _ := strconv.Atoi(msgTo)
				msg := msg.Message{
					To:   id,
					Body: msgBody,
				}
				s.SendMessage(&msg)
			} else {
				text = strings.ToUpper(text)
			}
		}

		w.WriteString(text + "\n")
		w.Flush()
	}
	return nil
}