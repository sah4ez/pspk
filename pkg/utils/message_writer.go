package utils

type Message struct {
	data string
}

func (m *Message) Write(data []byte) (n int, err error) {
	m.data = string(data)
	return len(m.data), nil
}

func (m *Message) Read() string {
	return m.data
}

func NewMessageWriter() *Message {
	return &Message{}
}
