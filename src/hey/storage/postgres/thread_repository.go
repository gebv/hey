package postgres

import (
	"context"
	"hey/storage"
	"hey/utils"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	relatedEventIDContextKey = "RelatedEventID"
	parentThreadIDContextKey = "ParentThreadID"
)

type ThreadRepository struct {
}

func (r *ThreadRepository) relatedEventIDFromContext(
	ctx context.Context,
) uuid.UUID {
	return getUUIDFromContext(relatedEventIDContextKey, ctx)
}

func (r *ThreadRepository) parentThreadIDFromContext(
	ctx context.Context,
) uuid.UUID {
	return getUUIDFromContext(parentThreadIDContextKey, ctx)
}

func (r *ThreadRepository) clientIDFromContext(ctx context.Context) uuid.UUID {
	return ClientIDFromContext(ctx)
}

// CreateThread create new thread
// waiting in the context of the client ID, linked event and thread IDs
func (r *ThreadRepository) CreateThread(
	ctx context.Context,
	channelID,
	threadID uuid.UUID,
	owners []uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	clientID := r.clientIDFromContext(ctx)
	relatedEventID := r.relatedEventIDFromContext(ctx)
	parentThreadID := r.parentThreadIDFromContext(ctx)
	conn := storage.FromContext(ctx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()
		_, err = conn.Exec(`INSERT INTO threads (
            thread_id,
            client_id,
            channel_id,

            owners,

            related_event_id,
            parent_thread_id,

            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			threadID,
			clientID,
			channelID,
			(&utils.UUIDS{}).FromArray(owners),
			relatedEventID,
			parentThreadID,
			time.Now(),
			time.Now(),
		)

		// TODO: create related entities
		/*
		   1. thread_counters
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
