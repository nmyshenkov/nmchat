package client

import (
	msg "nmchat/message"
	"strconv"
)

//Client - client's struct
type Client struct {
	ID           int
	Name         string
	IncomingChan chan *msg.Message
}

//GetTextID - return string ID
func (cl *Client) GetTextID() string {
	return strconv.Itoa(cl.ID)
}

//InitClient - init client
func (cl *Client) initClient(ID int) {
	cl.ID = ID
	cl.Name = "User" + strconv.Itoa(ID)
	cl.IncomingChan = make(chan *msg.Message)
}

// NewClient - create new client
func NewClient(ID int) *Client {
	client := Client{}
	client.initClient(ID)
	return &client
}
