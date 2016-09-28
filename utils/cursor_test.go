package utils

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCursor_simple(t *testing.T) {
	eventID := uuid.NewV4()
	createdAtEvent := time.Now()
	cursor1 := NewCursorFromSource(
		eventID,
		createdAtEvent,
	)

	cursor2 := NewCursorEvents(cursor1.String())
	assert.Equal(t,
		uuid.Equal(cursor2.LastEventID(), cursor1.LastEventID()),
		true,
	)

	assert.Equal(t,
		cursor1.LastCreatedAt(),
		cursor2.LastCreatedAt())

	// not valid cursor
	cursor2 = NewCursorEvents("invalid data")
	assert.Equal(t,
		uuid.Equal(
			cursor2.LastEventID(),
			uuid.Nil,
		),
		true,
	)

	assert.Equal(
		t,
		cursor2.LastCreatedAt(),
		time.Time{})

	// empty cursor
	cursor2 = NewCursorEvents("")
	assert.Equal(t,
		uuid.Equal(
			cursor2.LastEventID(),
			uuid.Nil,
		),
		true,
	)

	assert.Equal(
		t,
		cursor2.LastCreatedAt(),
		time.Time{})
}
