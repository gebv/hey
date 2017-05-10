package hey

import (
	"log"
	"reflect"
	"time"

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
	msgpack.Register(
		reflect.TypeOf(RelatedData{}),
		encodeRelatedData,
		decodeRelatedData,
	)
	msgpack.Register(
		reflect.TypeOf(EventObserver{}),
		encodeEventObserver,
		decodeEventObserver,
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

	if err = e.EncodeInt64(m.CreatedAt.UnixNano()); err != nil {
		return
	}

	if err = e.EncodeInt64(m.UpdatedAt.UnixNano()); err != nil {
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

	var secsCreatedAt int64
	if secsCreatedAt, err = d.DecodeInt64(); err != nil {
		return
	} else {
		m.CreatedAt = time.Unix(0, secsCreatedAt)
	}

	var secsUpdatedAt int64
	if secsUpdatedAt, err = d.DecodeInt64(); err != nil {
		return
	} else {
		m.UpdatedAt = time.Unix(0, secsUpdatedAt)
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

	if err = e.EncodeInt64(m.LastDeliveredTime.UnixNano()); err != nil {
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

	var secsLastDeliveredTime int64
	if secsLastDeliveredTime, err = d.DecodeInt64(); err != nil {
		return
	} else {
		m.LastDeliveredTime = time.Unix(0, secsLastDeliveredTime)
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

func encodeRelatedData(e *msgpack.Encoder, v reflect.Value) (err error) {
	m := v.Interface().(RelatedData)

	if err = e.EncodeSliceLen(4); err != nil {
		return
	}

	if err = e.EncodeString(m.UserID); err != nil {
		return
	}

	if err = e.EncodeString(m.EventID); err != nil {
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

func decodeRelatedData(d *msgpack.Decoder, v reflect.Value) (err error) {
	var l int

	m := v.Addr().Interface().(*RelatedData)
	if l, err = d.DecodeSliceLen(); err != nil {
		return
	}

	if l != 4 {
		return ErrMsgPackConflictFields
	}

	if m.UserID, err = d.DecodeString(); err != nil {
		return
	}
	if m.EventID, err = d.DecodeString(); err != nil {
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

func encodeEventObserver(e *msgpack.Encoder, v reflect.Value) (err error) {
	m := v.Interface().(EventObserver)

	if err = e.EncodeSliceLen(2); err != nil {
		return
	}

	if err = e.Encode(m.Event); err != nil {
		return
	}

	if err = e.Encode(m.RelatedData); err != nil {
		return
	}

	return
}

func decodeEventObserver(d *msgpack.Decoder, v reflect.Value) (err error) {
	var l int

	m := v.Addr().Interface().(*EventObserver)
	if l, err = d.DecodeSliceLen(); err != nil {
		return
	}

	if l != 2 {
		return ErrMsgPackConflictFields
	}

	if err = d.Decode(&m.Event); err != nil {
		return
	}
	if err = d.Decode(&m.RelatedData); err != nil {
		return
	}

	return
}
