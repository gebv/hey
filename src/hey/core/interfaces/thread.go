package interfaces

import uuid "github.com/satori/go.uuid"

type Thread interface {
	ThreadID() uuid.UUID
	ChannelID() uuid.UUID
	Owners() []uuid.UUID

	ParentThreadID() uuid.UUID
	RelatedEventID() uuid.UUID
}
