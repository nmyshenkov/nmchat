package server

import (
	"net"
	"time"
)

// Conn - telnet connection struct with timeout
type Conn struct {
	net.Conn
	IdleTimeout time.Duration
}

// Write - implements conn.Write with timeout
func (c *Conn) Write(p []byte) (int, error) {
	c.updateDeadline()
	return c.Conn.Write(p)
}

// Read - implements conn.Read with timeout
func (c *Conn) Read(b []byte) (int, error) {
	c.updateDeadline()
	return c.Conn.Read(b)
}

// updateDeadline - update connection timeout
func (c *Conn) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(idleDeadline)
}
