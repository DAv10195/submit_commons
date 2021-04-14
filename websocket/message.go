package websocket

import (
	"fmt"
	"github.com/DAv10195/submit_commons/containers"
	"strings"
)

var msgTypes = containers.NewStringSet()

type Message struct {
	Type 	string
	Payload []byte
}

func NewMessage(msgType string, payload []byte) (*Message, error) {
	if !msgTypes.Contains(msgType) {
		return nil, &MessageFormatErr{fmt.Sprintf("invalid message type [ %s ]", msgType)}
	}
	return &Message{msgType, payload}, nil
}

func FromBinary(payload []byte) (*Message, error) {
	payloadStr := string(payload)
	splitIndex := strings.IndexRune(payloadStr, '\n')
	if splitIndex < 0 {
		return nil, &MessageFormatErr{fmt.Sprintf("invalid message. Missing \\n separator between type and payload. Message: [ %s ]", payloadStr)}
	}
	msgType := payloadStr[ : splitIndex]
	if len(payloadStr) - 1 == splitIndex {
		return nil, &MessageFormatErr{fmt.Sprintf("invalid message. Empty payload. Message: [ %s ]", payloadStr)}
	}
	return NewMessage(msgType, []byte(payloadStr[splitIndex + 1 : ]))
}

func (m *Message) ToBinary() []byte {
	var sb strings.Builder
	_, _ = sb.WriteString(m.Type)
	_, _ = sb.WriteRune('\n')
	_, _ = sb.Write(m.Payload)
	return []byte(sb.String())
}

func init() {
	msgTypes.Add(MessageTypeKeepalive)
}
