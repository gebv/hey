package main

import (
	"encoding/json"

	chronograph "github.com/gebv/hey"
)

const (
	MessageDataType = "1"
)

func NewMessage(senderName string, body string) chronograph.Event {
	return chronograph.Event{
		DataType: MessageDataType,
		Data: MessageData{
			SenderName: senderName,
			Body:       body,
		}.Marshal(),
	}
}

type MessageData struct {
	SenderName string
	Body       string
}

func (n MessageData) Marshal() []byte {
	bts, _ := json.Marshal(n)
	return bts
}
