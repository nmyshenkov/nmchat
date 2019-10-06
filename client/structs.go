package client

import (
	msg "nmchat/message"
	"strconv"
)

// Client - client's struct
type Client struct {
	ID           int
	Name         string
	IncomingChan chan *msg.Message
}

// NewClient - create new client
func NewClient(ID int) *Client {
	client := Client{}
	client.initClient(ID)
	return &client
}

// InitClient - init client
func (cl *Client) initClient(ID int) {
	cl.ID = ID
	cl.Name = "User" + strconv.Itoa(ID)
	cl.IncomingChan = make(chan *msg.Message)
}

// GetTextID - return string ID
func (cl *Client) GetTextID() string {
	return strconv.Itoa(cl.ID)
}

func (cl *Client) NewBroadcastMessage(msgBody string) msg.Message {
	return msg.Message{
		From:        cl.ID,
		FromNikname: cl.Name,
		Body:        msgBody,
	}
}

//  NewMessage - new message to
func (cl *Client) NewMessage(to *Client, msgBody string) msg.Message {
	return msg.Message{
		To:          to.ID,
		ToNikname:   to.Name,
		From:        cl.ID,
		FromNikname: cl.Name,
		Body:        msgBody,
	}
}

func (cl *Client) ChangeName(newName string) {
	cl.Name = newName
}
