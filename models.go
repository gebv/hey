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
	CreatedAt time.Time
}

// Observer это таблица для хранения подписок на трэды.
// Для выбора пописок нужен итератор по UserID.
// Для выбора подписчиков нужен итератор по ThreadID
type Observer struct {
	UserID string
	// ID трэда, на который подписан юзер
	ThreadID          string
	LastDeliveredTime time.Time
}

// Sources таблица источников каджого трэда.
// используется при создании нового события в трэде,
// для всех трэдов-подписчиков (TargetThreadID) этого трэда SourceThreadID
// создаётся запись в таблице Threadline.
type Sources struct {
	TargetThreadID string
	SourceThreadID string
}

// User подписчик, обозреватель,
type User struct {
	UserID   string
	DataType DataType
	Data     interface{}
}
