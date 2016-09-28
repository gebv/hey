package postgres

import (
	"github.com/gebv/hey/utils"

	"time"

	uuid "github.com/satori/go.uuid"
)

// Thread

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

// Event

type event struct {
	eventID        uuid.UUID
	clientID       uuid.UUID
	threadID       uuid.UUID
	channelID      uuid.UUID
	creatorID      uuid.UUID
	data           []byte
	parentThreadID uuid.UUID
	parentEventID  uuid.UUID
	branchThreadID uuid.UUID
	createdAt      time.Time
	updatedAt      time.Time
}

func (e event) EventID() uuid.UUID {
	return e.eventID
}

func (e event) ClientID() uuid.UUID {
	return e.clientID
}

func (e event) ThreadID() uuid.UUID {
	return e.threadID
}

func (e event) ChannelID() uuid.UUID {
	return e.channelID
}

func (e event) CreatorID() uuid.UUID {
	return e.creatorID
}

func (e event) Data() []byte {
	return e.data
}

func (e event) ParentThreadID() uuid.UUID {
	return e.parentThreadID
}

func (e event) ParentEventID() uuid.UUID {
	return e.parentEventID
}

func (e event) BranchThreadID() uuid.UUID {
	return e.branchThreadID
}

func (e event) CreatedAt() time.Time {
	return e.createdAt
}

func (e event) UpdatedAt() time.Time {
	return e.updatedAt
}
