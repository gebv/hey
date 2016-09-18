package server

import (
	"context"
)

type EventDTO interface {
}

type EventService interface {
	CreateEvent(context.Context, EventDTO) (string, error)
}

type Service interface {
	EventService

	CreateChannel()
	ChangeOwnerChannel()

	Threadline(
		ctx context.Context,
		ChannelID int64,
		ThreadID string,
	) // returns the events from the channel\thread
}
