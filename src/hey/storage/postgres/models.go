package postgres

import (
	"hey/core/interfaces"
	"hey/utils"

	uuid "github.com/satori/go.uuid"
)

var (
	_ interfaces.Thread = (*thread)(nil)
)

type thread struct {
	threadID       uuid.UUID
	clientID       uuid.UUID
	channelID      uuid.UUID
	owners         utils.UUIDS
	relatedEventID uuid.UUID
	parentThreadID uuid.UUID
}

func (t thread) ThreadID() uuid.UUID {
	return t.threadID
}

func (t thread) ChannelID() uuid.UUID {
	return t.channelID
}

func (t thread) Owners() []uuid.UUID {
	return []uuid.UUID(t.owners)
}

func (t thread) ParentThreadID() uuid.UUID {
	return t.parentThreadID
}

func (t thread) RelatedEventID() uuid.UUID {
	return t.relatedEventID
}
