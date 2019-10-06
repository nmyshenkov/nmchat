package server

import (
	"bufio"
	"log"
	"net"
	cl "nmchat/client"
	msg "nmchat/message"
	"strconv"
)

func (s *Server) handleClient(conn net.Conn, client *cl.Client) {
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
				return
			}
			break
		}
		text := scanner.Text()

		switch text {
		case "":
			continue
		case "!exit":
			log.Printf("Client %s decited to exit", conn.RemoteAddr().String())
			return
		case "!help":
			text = HELP
		case "!list":
			text = "Users in the chat:\n"
			for _, client := range s.clients {
				text += "\t\t" + client.Name + " with ID: " + client.GetTextID() + "\n"
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
				text = "Nikname already exists!"
				if !s.checkExistNikname(client, arg) {
					client.ChangeName(arg)
					text = "Name changed!"
				}
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
	return
}
