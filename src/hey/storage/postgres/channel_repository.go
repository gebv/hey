package postgres

import (
	"context"
	"errors"
	"hey/storage"
	"hey/utils"
	"time"

	"github.com/satori/go.uuid"
)

type ChannelRepository struct {
}

func (r *ChannelRepository) clientIDFromContext(ctx context.Context) uuid.UUID {

	return ClientIDFromContext(ctx)
}

// CreateChannel create new channel
// waiting in the context of the client ID (key name 'ClientID')
func (r *ChannelRepository) CreateChannel(
	ctx context.Context,
	channelID,
	rootThreadID uuid.UUID,
	owners []uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
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
		_, err = conn.Exec(`INSERT INTO channels (
            channel_id,
            client_id,
            owners,
            root_thread_id,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6)`,
			channelID,
			clientID,
			(&utils.UUIDS{}).FromArray(owners),
			rootThreadID,
			time.Now(),
			time.Now(),
		)

		// TODO: insert into
		// 1. channel_counters
		// 2. channel_watchers
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// CreateRelatedEntities create the related entities (only for new channel)
func (r *ChannelRepository) CreateRelatedEntities(
	ctx context.Context,
	channelID uuid.UUID,
) error {
	// TODO: create
	/*
		1. channel_counters
		2. channel_watchers
	*/
	return errors.New("not implemented")
}
