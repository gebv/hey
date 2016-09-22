package utils

import (
	"errors"
	"hey/core/interfaces"
	"time"

	uuid "github.com/satori/go.uuid"
)

type CursorEvents string

var (
	EmptyCursorEvents = CursorEvents("")
	lengthCursor      = 16 + 15
)

func NewCursorEvents(str string) CursorEvents {

	return CursorEvents(str)
}

func NewCursorFromSource(eventID uuid.UUID, createdAt time.Time) CursorEvents {
	res := make([]byte, lengthCursor)
	copy(res, eventID.Bytes())
	b2, err := createdAt.MarshalBinary()
	if err != nil {
		println("[WARNING] marshal time", err.Error())
	}
	copy(res[16:], b2)

	return CursorEvents(Base64(res))
}

func NewCursorEventsFromEvent(e interfaces.Event) CursorEvents {

	return NewCursorFromSource(e.EventID(), e.CreatedAt())
}

func (c CursorEvents) String() string {
	return string(c)
}

func (c CursorEvents) decode() ([]byte, error) {
	dec, err := DecodeBase64(c.String())
	if len(dec) != lengthCursor {
		return []byte{}, errors.New("not valid len cursor")
	}
	return dec, err
}

func (c CursorEvents) LastEventID() uuid.UUID {
	dec, err := c.decode()

	if err != nil {
		println("[WARNING] decode base 64", err.Error())
		return uuid.Nil
	}

	return uuid.FromBytesOrNil(dec[:16])
}

func (c CursorEvents) LastCreatedAt() time.Time {
	dec, err := c.decode()

	if err != nil {
		println("[WARNING] decode base 64", err.Error())
		return time.Time{}
	}
	createdAt := time.Now()
	err = createdAt.UnmarshalBinary(dec[16:])
	if err != nil {
		println("[WARNING] unmarshal binary datetime", err.Error())
		return time.Time{}
	}
	return createdAt
}
