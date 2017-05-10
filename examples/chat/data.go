package main

import chronograph "github.com/zhuharev/chronograph"

const (
	MessageDataType chronograph.DataType = 1
)

func init() {
	chronograph.RegDataType(MessageDataType, func() interface{} {
		return &MessageData{}
	})
}

func NewMessage(senderName string, body string) chronograph.Event {
	return chronograph.Event{
		DataType: MessageDataType,
		Data: &MessageData{
			SenderName: senderName,
			Body:       body,
		},
	}
}

type MessageData struct {
	SenderName string
	Body       string
}
