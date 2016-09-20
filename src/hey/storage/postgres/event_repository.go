package postgres

import (
	"context"
	"hey/storage"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	parentEventIDContextKey  = "ParentEventID"
	branchThreadIDContextKey = "BranchThreadID"
	channelIDContextKey      = "ChannelID"
)

type EventRepository struct {
}

func (r *EventRepository) parentEventIDFromContext(
	ctx context.Context,
) uuid.UUID {
	return getUUIDFromContext(parentEventIDContextKey, ctx)
}

func (r *EventRepository) parentThreadIDFromContext(
	ctx context.Context,
) uuid.UUID {
	return getUUIDFromContext(parentThreadIDContextKey, ctx)
}

func (r *EventRepository) branchThreadIDFromContext(
	ctx context.Context,
) uuid.UUID {
	return getUUIDFromContext(branchThreadIDContextKey, ctx)
}

func (r *EventRepository) channelIDFromContext(
	ctx context.Context,
) uuid.UUID {
	return getUUIDFromContext(channelIDContextKey, ctx)
}

func (r *EventRepository) clientIDFromContext(ctx context.Context) uuid.UUID {
	return ClientIDFromContext(ctx)
}

// CreateEvent
func (r *EventRepository) CreateEvent(
	ctx context.Context,
	threadID,
	eventID,
	creatorID uuid.UUID,
	data []byte,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	channelID := r.channelIDFromContext(ctx)
	clientID := r.clientIDFromContext(ctx)
	parentEventID := r.parentEventIDFromContext(ctx)
	parentThreadID := r.parentThreadIDFromContext(ctx)
	conn := storage.FromContext(ctx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()
		_, err = conn.Exec(`INSERT INTO events (
            event_id,
            client_id,
            thread_id,
            channel_id,

            creator,

            data,

            parent_thread_id,
            parent_event_id,
            branch_thread_id,

            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			eventID,
			clientID,
			threadID,
			channelID,
			creatorID,
			data, // data
			parentThreadID,
			parentEventID,
			uuid.Nil, // branch thread id
			time.Now(),
			time.Now(),
		)

		// TODO: create related entities
		/*
		   1. threadline
		   2. thread_watchers
		*/
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
