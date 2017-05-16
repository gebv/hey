package main

import (
	chronograph "github.com/gebv/hey"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

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

func (m *MessageData) EncodeMsgpack(enc *msgpack.Encoder) error {
	enc.EncodeSliceLen(2)
	enc.EncodeString(m.SenderName)
	enc.EncodeString(m.Body)
	return nil
}

func (m *MessageData) DecodeMsgpack(dec *msgpack.Decoder) error {
	dec.DecodeSliceLen()
	m.SenderName, _ = dec.DecodeString()
	m.Body, _ = dec.DecodeString()
	return nil
}
