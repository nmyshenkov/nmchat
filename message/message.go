package message

// Message - message
type Message struct {
	FromNikname string
	From        int
	ToNikname   string
	To          int
	Body        string
}

// FromStringMessage - get full string answer
func (m *Message) FromStringMessage() string {
	return m.FromNikname + ": " + m.Body
}

// FromByteMessage - get full byte answer
func (m *Message) FromByteMessage() []byte {
	return []byte(m.FromNikname + ": " + m.Body)
}
