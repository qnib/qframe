package qtypes_messages


type Message struct {
	Base
	Message string
}

func NewMessage(b Base, msg string) Message {
	m := Message{
		Base: b,
		Message: msg,
	}
	m.GenDefaultID()
	return m
}
