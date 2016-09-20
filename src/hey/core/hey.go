package core

import (
	"context"
	"hey/storage"

	uuid "github.com/satori/go.uuid"
)

func NewHeyService() *HeyService {
	return &HeyService{}
}

type HeyService struct {
	conn         storage.DB
	eventsRepo   EventRepository
	threadsRepo  ThreadRepository
	channelsRepo ChannelRepository
}

// interfaces

type EventRepository interface {
	CreateChannel(
		ctx context.Context,
		eventID,
		threadID,
		creatorID uuid.UUID,
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
	CreateThread(
		ctx context.Context,
		channelID,
		threadID uuid.UUID,
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
	)
	SetUnreadByUser(
		ctx context.Context,
		channelID,
		userID uuid.UUID,
		count int,
	) error
}
