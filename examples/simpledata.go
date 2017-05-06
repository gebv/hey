package examples

import (
	"reflect"
	"time"

	"github.com/gebv/hey"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

const (
	ThreadSimpleDataType hey.DataType = 1
)

func init() {
	hey.RegDataType(ThreadSimpleDataType, func() interface{} {
		return &ThreadSimpleData{}
	})
}

func NewThreadSimpleData(
	a string,
	b int64,
	c time.Time,
) hey.Thread {
	return hey.Thread{
		DataType: ThreadSimpleDataType,
		Data: &ThreadSimpleData{
			A: a,
			B: b,
			C: c,
		},
	}
}

type ThreadSimpleData struct {
	A string
	B int64
	C time.Time
}

func init() {
	msgpack.Register(
		reflect.TypeOf(ThreadSimpleData{}),
		encodeThreadSimpleData,
		decodeThreadSimpleData,
	)
}

func encodeThreadSimpleData(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(ThreadSimpleData)

	if err := e.EncodeSliceLen(3); err != nil {
		return err
	}

	//1
	if err := e.EncodeString(m.A); err != nil {
		return err
	}
	//2
	if err := e.EncodeInt64(m.B); err != nil {
		return err
	}
	//3
	if err := e.EncodeTime(m.C); err != nil {
		return err
	}

	return nil
}

func decodeThreadSimpleData(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int

	m := v.Addr().Interface().(*ThreadSimpleData)

	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}

	if l != 3 {
		return hey.ErrMsgPackConflictFields
	}

	//1
	if data, err := d.DecodeString(); err == nil {
		m.A = data
	} else {
		return err
	}
	//2
	if data, err := d.DecodeInt64(); err == nil {
		m.B = data
	} else {
		return err
	}
	//3
	if data, err := d.DecodeTime(); err == nil {
		m.C = data
	} else {
		return err
	}

	return nil
}
