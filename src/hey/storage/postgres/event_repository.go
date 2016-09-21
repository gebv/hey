package postgres

import (
	"context"
	"hey/storage"
	"time"

	uuid "github.com/satori/go.uuid"
)

type EventRepository struct {
}

func (r *EventRepository) clientIDFromContext(ctx context.Context) uuid.UUID {
	return ClientIDFromContext(ctx)
}

// CreateEvent create new event
// waiting in the context of the client ID, channel ID, linked parent
// event and thread IDs
func (r *EventRepository) CreateEvent(
	ctx context.Context,
	eventID,
	threadID,
	channelID,
	creatorID,
	parentThreadID,
	parentEventID,
	branchThreadID uuid.UUID,
	data []byte,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	clientID := r.clientIDFromContext(ctx)
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
			branchThreadID, // branch thread id
			time.Now(),
			time.Now(),
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (r *EventRepository) CreateThreadline(
	ctx context.Context,
	channelID,
	threadID,
	eventID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	clientID := r.clientIDFromContext(ctx)
	conn := storage.FromContext(ctx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()
		_, err = conn.Exec(`INSERT INTO threadline (
            client_id,
            channel_id,
            thread_id,
            event_id,
            created_at
        ) VALUES ($1, $2, $3, $4, $5)`,
			clientID,
			channelID,
			threadID,
			eventID,
			time.Now(),
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
