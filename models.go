package hey

import (
	"errors"
	"log"
	"reflect"
	"time"

	uuid "github.com/satori/go.uuid"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var (
	ErrMsgPackConflictFields = errors.New("msgpack serializator: conflict num fields or another error")
	ErrNotRegDataType        = errors.New("not registred data type")
)

// Thread

type Thread struct {
	ThreadID string
	DataType DataType
	Data     interface{}
}

// Events

type Event struct {
	EventID   uuid.UUID
	ThreadID  string
	CreatorID string
	DataType  DataType
	Data      interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Threadline struct {
	EventID   uuid.UUID
	ThreadID  string
	CreatedAt time.Time
	DataType  DataType
	Data      interface{}
}

type Observer struct {
	UserID                  string
	ThreadID                string
	ParentThreadID          string
	RelatedThreadlineExists bool
	LastTimeStamp           time.Time // unix time stamp
	DataType                DataType
	Data                    interface{}
}

type EventObserver struct {
	ObserverID string // related Observer.UserID

	Event      Event
	Threadline Threadline
}

type User struct {
	UserID   string
	DataType DataType
	Data     interface{}
}

// Serializer for thread

func init() {
	msgpack.Register(
		reflect.TypeOf(Thread{}),
		encodeThread,
		decodeThread,
	)
}

func encodeThread(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Thread)

	if err := e.EncodeSliceLen(3); err != nil {
		return err
	}

	//1
	if err := e.EncodeString(m.ThreadID); err != nil {
		return err
	}
	//2
	if err := e.EncodeUint32(uint32(m.DataType)); err != nil {
		return err
	}
	//3
	if err := e.Encode(m.Data); err != nil {
		return err
	}

	return nil
}

func decodeThread(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int

	m := v.Addr().Interface().(*Thread)

	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}

	if l != 3 {
		return ErrMsgPackConflictFields
	}

	//1
	if data, err := d.DecodeString(); err == nil {
		m.ThreadID = data
	} else {
		return err
	}
	//2
	if data, err := d.DecodeUint32(); err == nil {
		m.DataType = DataType(data)
	} else {
		return err
	}

	//5
	m.Data, err = FactoryDataObj(m.DataType)
	if err == nil {
		if err = d.Decode(&m.Data); err != nil {
			return err
		}
	} else {
		d.Skip()
		log.Printf("hey: not supported data type DataType(%d)", m.DataType)
		return ErrNotRegDataType
	}

	return nil
}
