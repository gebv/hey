package interfaces

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Thread interface {
	ThreadID() uuid.UUID
	ChannelID() uuid.UUID
	Owners() []uuid.UUID

	ParentThreadID() uuid.UUID
	RelatedEventID() uuid.UUID
}

type Event interface {
	EventID() uuid.UUID
	ChannelID() uuid.UUID
	ThreadID() uuid.UUID
	CreatorID() uuid.UUID
	Data() []byte

	ParentThreadID() uuid.UUID
	ParentEventID() uuid.UUID
	BranchThreadID() uuid.UUID

	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type EventProvider interface {
	Event(ctx context.Context,
		eventID uuid.UUID,
	) (Event, error)
}
