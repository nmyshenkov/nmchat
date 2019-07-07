package message

import "strings"

// Message - message
type Message struct {
	FromNikname string
	From        int
	ToNikname   string
	To          int
	Body        string
}

// FromStringMessage - get full string answer
func (m *Message) FromStringMessage(addtional ...string) string {
	return m.FromNikname + ": " + m.Body + strings.Join(addtional, " ")
}

// FromByteMessage - get full byte answer
func (m *Message) FromByteMessage(addtional ...string) []byte {
	return []byte(m.FromNikname + ": " + m.Body + strings.Join(addtional, " "))
}
