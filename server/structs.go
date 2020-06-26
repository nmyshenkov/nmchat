package server

import (
	cl "github.com/nmyshenkov/nmchat/client"
	msg "github.com/nmyshenkov/nmchat/message"
	"time"
)

// Server - main struct
type Server struct {
	Addr        string
	messages    []*msg.Message
	sendMsgCh   chan *msg.Message
	clients     map[int]*cl.Client
	clNextID    int
	addCliCh    chan *cl.Client
	delCliCh    chan *cl.Client
	doneCh      chan bool
	errCh       chan error
	IdleTimeout time.Duration
}
