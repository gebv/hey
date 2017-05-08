package hey

import (
	"log"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

type DataType uint32

const (
	EmptyData DataType = 0
)

func (t DataType) EncodeMsgpack(dec *msgpack.Encoder) error {
	if err := dec.EncodeUint32(uint32(t)); err != nil {
		return err
	}
	return nil
}

func (t *DataType) DecodeMsgpack(dec *msgpack.Decoder) error {
	if val, err := dec.DecodeUint32(); err != nil {
		return err
	} else {
		*t = DataType(val)
	}
	return nil
}

var types = make(map[DataType]func() interface{})

func RegDataType(t DataType, f func() interface{}) {
	types[t] = f
}

func FactoryDataObj(t DataType) (interface{}, error) {
	f, exists := types[t]

	if !exists {
		log.Printf("hey: not registred data type DataType(%d)", t)
		return nil, ErrNotRegDataType
	}

	return f(), nil
}
