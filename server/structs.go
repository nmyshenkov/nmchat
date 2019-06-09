package server

import (
	cl "nmchat/client"
	msg "nmchat/message"
	"time"
)

type Server struct {
	Addr        string
	messages    []*msg.Message
	sendMsgCh   chan *msg.Message
	clients     map[int]*cl.Client
	addCliCh    chan *cl.Client
	delCliCh    chan *cl.Client
	doneCh      chan bool
	errCh       chan error
	IdleTimeout time.Duration
}
