package websocket

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFromBinary(t *testing.T) {
	inpMessageType := MessageTypeKeepalive
	payload := "{\"name\":\"David\"}"
	m, err := FromBinary([]byte(fmt.Sprintf("%s\n%s", inpMessageType, payload)))
	if err != nil {
		t.Fatalf("error formatting binary message: %v", err)
	}
	if m.Type != MessageTypeKeepalive {
		t.Fatalf("expected [ %s ], got [ %s ]", inpMessageType, m.Type)
	}
	formattedPayload := string(m.Payload)
	if formattedPayload != payload {
		t.Fatalf("expected [ %s ], got [ %s ]", payload, formattedPayload)
	}
	_, err = FromBinary([]byte(fmt.Sprintf("%s%s", inpMessageType, payload)))
	if err == nil {
		t.Fatal("no error returned from formatting invalid binary message with no \\n separator")
	}
	_, err = FromBinary([]byte(fmt.Sprintf("%s\n", inpMessageType)))
	if err == nil {
		t.Fatal("no error returned from formatting invalid binary message with empty payload")
	}
	inpMessageType = "bad"
	_, err = FromBinary([]byte(fmt.Sprintf("%s\n%s", inpMessageType, payload)))
	if err == nil {
		t.Fatal("no error returned from formatting invalid binary message with bad type")
	}
}

func TestToBinary(t *testing.T) {
	m1, err := NewMessage(MessageTypeKeepalive, []byte("{\"name\":\"Nikita\""))
	if err != nil {
		t.Fatalf("error creating message: %v", err)
	}
	m2, err := FromBinary(m1.ToBinary())
	if err != nil {
		t.Fatalf("error formatting valid binary message: %v", err)
	}
	if !reflect.DeepEqual(m1, m2) {
		t.Fatalf("message should be equal but they are not. m1 : [ %v ], m2 : [ %v ]", m1, m2)
	}
}

func TestNewMessage(t *testing.T) {
	if _, err := NewMessage(MessageTypeKeepalive, []byte("{\"name\":\"Azriel\"")); err != nil {
		t.Fatalf("error creating message: %v", err)
	}
	if _, err := NewMessage("bad", []byte("{\"name\":\"Azriel\"")); err == nil {
		t.Fatal("no error returned from creating invalid message with bad message type")
	}
}
