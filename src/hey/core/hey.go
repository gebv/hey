package core

import (
	"context"
	"hey/core/interfaces"
	"hey/storage"

	uuid "github.com/satori/go.uuid"
)

func NewHeyService(
	conn storage.DB,
	eventsRepo EventRepository,
	threadsRepo ThreadRepository,
	channelsRepo ChannelRepository,
) *HeyService {
	return &HeyService{
		conn:         conn,
		eventsRepo:   eventsRepo,
		threadsRepo:  threadsRepo,
		channelsRepo: channelsRepo,
	}
}

type HeyService struct {
	conn         storage.DB
	eventsRepo   EventRepository
	threadsRepo  ThreadRepository
	channelsRepo ChannelRepository
}

// interfaces

type EventRepository interface {
	Event(
		ctx context.Context,
		eventID uuid.UUID,
	) (interfaces.Event, error)

	FindEvents(
		ctx context.Context,
		threadID uuid.UUID,
		cursorStr string,
		perPage int,
	) ([]interfaces.Event, string, error)

	FindEventsWithProvider(
		ctx context.Context,
		threadID uuid.UUID,
		cursorStr string,
		perPage int,
		provider interfaces.EventProvider,
	) ([]interfaces.Event, string, error)

	Threadline(
		ctx context.Context,
		channelID,
		threadID,
		eventID uuid.UUID,
	) error

	DeleteThreadline(
		ctx context.Context,
		eventID uuid.UUID,
	) error

	SetBranchThreadID(
		ctx context.Context,
		eventID,
		branchThreadID uuid.UUID,
	) error

	CreateEvent(
		ctx context.Context,
		eventID,
		threadID,
		channelID,
		creatorID,
		parentThreadID,
		parentEventID,
		branchThreadID uuid.UUID,
		data []byte,
	) error

	CreateThreadline(
		ctx context.Context,
		channelID,
		threadID,
		eventID uuid.UUID,
	) error
}

type ThreadRepository interface {
	FindThread(
		ctx context.Context,
		threadID uuid.UUID) (interfaces.Thread, error)

	CreateThread(
		ctx context.Context,
		threadID,
		channelID,
		relatedEventID,
		parentThreadID uuid.UUID,
		owners []uuid.UUID,
	) error
	AddCountEvents(
		ctx context.Context,
		threadID uuid.UUID,
		count int,
	) error
	SetUnreadByUser(
		ctx context.Context,
		threadID,
		userID uuid.UUID,
		count int,
	) error
}

type ChannelRepository interface {
	CreateChannel(
		ctx context.Context,
		channelID,
		rootThreadID uuid.UUID,
		owners []uuid.UUID,
	) error

	AddCountEvents(
		ctx context.Context,
		channelID uuid.UUID,
		count int,
	) error

	SetUnreadByUser(
		ctx context.Context,
		channelID,
		userID uuid.UUID,
		count int,
	) error
}
