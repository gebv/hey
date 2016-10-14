package hey

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Hey service
type Hey interface {
	// CreateChannel create new channel with root thread
	CreateChannel(
		ctx context.Context,
		userIDs []uuid.UUID,
	) (uuid.UUID, uuid.UUID, error)

	CreateChannelName(
		ctx context.Context,
		name string,
		userIDs []uuid.UUID,
	) (channelID uuid.UUID, rootThreadID uuid.UUID, err error)

	// CreateNodalEvent create new nodal event
	// waiting ChannelID from context
	CreateNodalEvent(
		ctx context.Context,
		threadID uuid.UUID,
		owners []uuid.UUID,
		creatorID uuid.UUID,
	) (newThreadID uuid.UUID, newEventID uuid.UUID, err error)

	CreateNodalEventWithData(
		ctx context.Context,
		threadID uuid.UUID,
		owners []uuid.UUID,
		creatorID uuid.UUID,
		data []byte,
	) (uuid.UUID, uuid.UUID, error)

	CreateNodalEventWithThreadName(
		ctx context.Context,
		threadName string,
		threadID uuid.UUID,
		owners []uuid.UUID,
		creatorID uuid.UUID,
	) (uuid.UUID, uuid.UUID, error)

	CreateNodalEventWithThreadNameWithData(
		ctx context.Context,
		threadName string,
		threadID uuid.UUID,
		owners []uuid.UUID,
		creatorID uuid.UUID,
		data []byte,
	) (uuid.UUID, uuid.UUID, error)

	// CreateNewBranchEvent create a new event in branch
	// if the event already has the branch - error
	CreateNewBranchEvent(
		ctx context.Context,
		threadID uuid.UUID,
		relatedEventID uuid.UUID, //
		owners []uuid.UUID,
		creatorID uuid.UUID,
		data []byte,
	) (newThreadID uuid.UUID, newEventID uuid.UUID, err error)

	CreateNewBranchEventWithThreadName(
		ctx context.Context,
		threadName string,
		threadID uuid.UUID,
		relatedEventID uuid.UUID, //
		owners []uuid.UUID,
		creatorID uuid.UUID,
		data []byte,
	) (uuid.UUID, uuid.UUID, error)

	// CreateEvent create event in existing thread
	CreateEvent(ctx context.Context,
		threadID uuid.UUID,
		creatorID uuid.UUID,
		data []byte,
	) (eventID uuid.UUID, err error)

	// CreateEventInThreadName create event in existing thread name
	// CreateEventInThreadName(ctx context.Context,
	// 	name string,
	// 	creatorID uuid.UUID,
	// 	data []byte,
	// ) (eventID uuid.UUID, err error)

	FindThreadByName(
		ctx context.Context,
		channelID uuid.UUID,
		name string,
	) (Thread, error)

	FindChannelByName(
		ctx context.Context,
		name string,
	) (Channel, error)

	// FindEvents find events
	// waiting WatcherID (from a user view) from context
	FindEvents(
		ctx context.Context,
		watcherID uuid.UUID,
		threadID uuid.UUID,
		cursorStr string,
		perPage int,
	) (SearchResult, error)

	FindEventsByName(
		ctx context.Context,
		watcherID uuid.UUID,
		specialName string,
		cursorStr string,
		perPage int,
	) (SearchResult, error)
}

// SearchResult search result events
type SearchResult interface {
	Events() []Event
	// EventsIterator

	Cursor() string
	HasNext() bool
}

// type EventsIterator interface {
// 	Value() Event
// 	Next() bool
//  Reset() // reset the index
// }

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

type Thread interface {
	ThreadID() uuid.UUID
	ChannelID() uuid.UUID
	Owners() []uuid.UUID
	ParentThreadID() uuid.UUID
	RelatedEventID() uuid.UUID
}

type Channel interface {
	ClientID() uuid.UUID
	ChannelID() uuid.UUID
	Owners() []uuid.UUID
	RootThreadID() uuid.UUID
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

// type EventProvider interface {
// 	FindEvent(
// 		ctx context.Context,
// 		watcherID uuid.UUID,
// 		evetnID uuid.UUID,
// 	) (Event, error)
// }

// type ThreadProvider interface {
// 	FindThread(
// 		ctx context.Context,
// 		threadID uuid.UUID,
// 	) (Thread, error)
// }
