package message

type Message struct {
	From string
	To   string
	Body string
}

func (m *Message) ToString() string {
	return m.From + ": " + m.Body
}
