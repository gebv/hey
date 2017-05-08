package hey

import (
	"errors"
	"time"
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
	EventID   string
	ThreadID  string
	DataType  DataType
	Data      interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Threadline - таблица для хранения индекса событий каждой подписки
// выборка происходит в обратном хронологическом порядке по ThreadID
type Threadline struct {
	EventID   string
	ThreadID  string
	UpdatedAt time.Time
}

// Observer это таблица для хранения подписок на трэды.
// используется при создании нового события в трэде,
// для всех подписчиков этого трэда
// создаётся запись в таблице Threadline
type Observer struct {
	UserID   string
	ThreadID string
	//ParentThreadID          string
	//RelatedThreadlineExists bool
	LastTimeStamp time.Time // unix time stamp
	DataType      DataType
	Data          interface{}
}

// User подписчик, обозреватель,
type User struct {
	UserID   string
	DataType DataType
	Data     interface{}
}
