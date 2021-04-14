package websocket

type MessageFormatErr struct {
	Message	string
}

func (e *MessageFormatErr) Error() string {
	return e.Message
}
