package hey

import (
	"log"
	"reflect"

	"gopkg.in/vmihailenco/msgpack.v2"
)

func init() {
	msgpack.Register(
		reflect.TypeOf(Thread{}),
		encodeThread,
		decodeThread,
	)
	msgpack.Register(
		reflect.TypeOf(Event{}),
		encodeEvent,
		decodeEvent,
	)
	msgpack.Register(
		reflect.TypeOf(Threadline{}),
		encodeThreadline,
		decodeThreadline,
	)
	msgpack.Register(
		reflect.TypeOf(Observer{}),
		encodeObserver,
		decodeObserver,
	)
	msgpack.Register(
		reflect.TypeOf(Sources{}),
		encodeSources,
		decodeSources,
	)
	msgpack.Register(
		reflect.TypeOf(User{}),
		encodeUser,
		decodeUser,
	)

}

func encodeThread(e *msgpack.Encoder, v reflect.Value) (err error) {
	m := v.Interface().(Thread)

	if err = e.EncodeSliceLen(3); err != nil {
		return
	}
	if err = e.EncodeString(m.ThreadID); err != nil {
		return
	}
	if err = e.Encode(m.DataType); err != nil {
		return
	}
	if err = e.Encode(m.Data); err != nil {
		return
	}
	return
}

func decodeThread(d *msgpack.Decoder, v reflect.Value) (err error) {
	var l int

	m := v.Addr().Interface().(*Thread)
	if l, err = d.DecodeSliceLen(); err != nil {
		return
	}

	if l != 3 {
		return ErrMsgPackConflictFields
	}

	if m.ThreadID, err = d.DecodeString(); err != nil {
		return
	}
	if err = d.Decode(&m.DataType); err != nil {
		return
	}

	if m.Data, err = FactoryDataObj(m.DataType); err != nil {
		if err = d.Skip(); err != nil {
			return
		}
		log.Printf("hey: not supported data type DataType(%d)", m.DataType)
		return ErrNotRegDataType
	} else if err = d.Decode(&m.Data); err != nil {
		return
	}

	return
}

func encodeEvent(e *msgpack.Encoder, v reflect.Value) (err error) {
	m := v.Interface().(Event)

	if err = e.EncodeSliceLen(6); err != nil {
		return
	}
	if err = e.EncodeString(m.EventID); err != nil {
		return
	}
	if err = e.EncodeString(m.ThreadID); err != nil {
		return
	}
	if err = e.Encode(m.DataType); err != nil {
		return
	}
	if err = e.Encode(m.Data); err != nil {
		return
	}
	if err = e.EncodeTime(m.CreatedAt); err != nil {
		return
	}
	if err = e.EncodeTime(m.UpdatedAt); err != nil {
		return
	}
	return
}

func decodeEvent(d *msgpack.Decoder, v reflect.Value) (err error) {
	var l int

	m := v.Addr().Interface().(*Event)
	if l, err = d.DecodeSliceLen(); err != nil {
		return
	}

	if l != 6 {
		return ErrMsgPackConflictFields
	}

	if m.EventID, err = d.DecodeString(); err != nil {
		return
	}
	if m.ThreadID, err = d.DecodeString(); err != nil {
		return
	}
	if err = d.Decode(&m.DataType); err != nil {
		return
	}

	if m.Data, err = FactoryDataObj(m.DataType); err != nil {
		if err = d.Skip(); err != nil {
			return
		}
		log.Printf("hey: not supported data type DataType(%d)", m.DataType)
		return ErrNotRegDataType
	} else if err = d.Decode(&m.Data); err != nil {
		return
	}

	if m.CreatedAt, err = d.DecodeTime(); err != nil {
		return
	}
	if m.UpdatedAt, err = d.DecodeTime(); err != nil {
		return
	}

	return
}

func encodeThreadline(e *msgpack.Encoder, v reflect.Value) (err error) {
	m := v.Interface().(Threadline)

	if err = e.EncodeSliceLen(3); err != nil {
		return
	}
	if err = e.EncodeString(m.EventID); err != nil {
		return
	}
	if err = e.EncodeString(m.ThreadID); err != nil {
		return
	}
	if err = e.EncodeTime(m.CreatedAt); err != nil {
		return
	}
	return
}

func decodeThreadline(d *msgpack.Decoder, v reflect.Value) (err error) {
	var l int

	m := v.Addr().Interface().(*Threadline)
	if l, err = d.DecodeSliceLen(); err != nil {
		return
	}

	if l != 3 {
		return ErrMsgPackConflictFields
	}

	if m.EventID, err = d.DecodeString(); err != nil {
		return
	}
	if m.ThreadID, err = d.DecodeString(); err != nil {
		return
	}
	if m.CreatedAt, err = d.DecodeTime(); err != nil {
		return
	}

	return
}

func encodeObserver(e *msgpack.Encoder, v reflect.Value) (err error) {
	m := v.Interface().(Observer)

	if err = e.EncodeSliceLen(3); err != nil {
		return
	}
	if err = e.EncodeString(m.UserID); err != nil {
		return
	}
	if err = e.EncodeString(m.ThreadID); err != nil {
		return
	}
	if err = e.EncodeTime(m.LastDeliveredTime); err != nil {
		return
	}
	return
}

func decodeObserver(d *msgpack.Decoder, v reflect.Value) (err error) {
	var l int

	m := v.Addr().Interface().(*Observer)
	if l, err = d.DecodeSliceLen(); err != nil {
		return
	}

	if l != 3 {
		return ErrMsgPackConflictFields
	}

	if m.UserID, err = d.DecodeString(); err != nil {
		return
	}
	if m.ThreadID, err = d.DecodeString(); err != nil {
		return
	}
	if m.LastDeliveredTime, err = d.DecodeTime(); err != nil {
		return
	}

	return
}

func encodeSources(e *msgpack.Encoder, v reflect.Value) (err error) {
	m := v.Interface().(Sources)

	if err = e.EncodeSliceLen(2); err != nil {
		return
	}
	if err = e.EncodeString(m.TargetThreadID); err != nil {
		return
	}
	if err = e.EncodeString(m.SourceThreadID); err != nil {
		return
	}
	return
}

func decodeSources(d *msgpack.Decoder, v reflect.Value) (err error) {
	var l int

	m := v.Addr().Interface().(*Sources)
	if l, err = d.DecodeSliceLen(); err != nil {
		return
	}

	if l != 2 {
		return ErrMsgPackConflictFields
	}

	if m.TargetThreadID, err = d.DecodeString(); err != nil {
		return
	}
	if m.SourceThreadID, err = d.DecodeString(); err != nil {
		return
	}

	return
}

func encodeUser(e *msgpack.Encoder, v reflect.Value) (err error) {
	m := v.Interface().(User)

	if err = e.EncodeSliceLen(3); err != nil {
		return
	}
	if err = e.EncodeString(m.UserID); err != nil {
		return
	}
	if err = e.Encode(m.DataType); err != nil {
		return
	}
	if err = e.Encode(m.Data); err != nil {
		return
	}
	return
}

func decodeUser(d *msgpack.Decoder, v reflect.Value) (err error) {
	var l int

	m := v.Addr().Interface().(*User)
	if l, err = d.DecodeSliceLen(); err != nil {
		return
	}

	if l != 3 {
		return ErrMsgPackConflictFields
	}

	if m.UserID, err = d.DecodeString(); err != nil {
		return
	}
	if err = d.Decode(&m.DataType); err != nil {
		return
	}

	if m.Data, err = FactoryDataObj(m.DataType); err != nil {
		if err = d.Skip(); err != nil {
			return
		}
		log.Printf("hey: not supported data type DataType(%d)", m.DataType)
		return ErrNotRegDataType
	} else if err = d.Decode(&m.Data); err != nil {
		return
	}

	return
}
