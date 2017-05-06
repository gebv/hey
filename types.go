package hey

import "log"

type DataType uint32

const (
	EmptyData DataType = 0
)

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
