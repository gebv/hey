package postgres

import (
	"context"
	"hey/core/interfaces"
	"hey/storage"
	"hey/utils"
	"time"

	uuid "github.com/satori/go.uuid"
)

type ThreadRepository struct {
}

func (r *ThreadRepository) clientIDFromContext(ctx context.Context) uuid.UUID {
	return ClientIDFromContext(ctx)
}

// FindThread returns thread
func (r *ThreadRepository) FindThread(
	ctx context.Context,
	threadID uuid.UUID,
) (interfaces.Thread, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	thread := thread{}

	clientID := r.clientIDFromContext(ctx)
	conn := storage.FromContext(ctx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()
		err = conn.QueryRow(`
			SELECT 
				thread_id,
				client_id, 
				channel_id,
				owners,
				related_event_id,
				parent_thread_id
			FROM threads
			WHERE client_id = $1 AND thread_id = $2
		`, clientID, threadID).Scan(
			&thread.threadID,
			&thread.clientID,
			&thread.channelID,
			&thread.owners,
			&thread.relatedEventID,
			&thread.parentThreadID,
		)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-done:
		return &thread, err
	}
}

// CreateThread create new thread
// waiting in the context of the client ID, linked event and thread IDs
func (r *ThreadRepository) CreateThread(
	ctx context.Context,
	threadID,
	channelID,
	relatedEventID,
	parentThreadID uuid.UUID,
	owners []uuid.UUID,
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
		if err != nil {
			return
		}

		err = r.AddCountEvents(
			ctx,
			threadID,
			0,
		)

		if err != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// AddCountEvents increases the number of events in the thread
func (r *ThreadRepository) AddCountEvents(
	ctx context.Context,
	threadID uuid.UUID,
	count int,
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
		// sql := `UPDATE thread_counters SET
		//     counter_events = counter_events + $1
		//     WHERE client_id = $2 AND  thread_id = $3`
		sql := `INSERT INTO thread_counters (
            client_id,
            thread_id,
            counter_events
        ) VALUES ($1, $2, $3)
        ON CONFLICT (client_id, thread_id)
        DO UPDATE SET
            counter_events = thread_counters.counter_events + EXCLUDED.counter_events`
		_, err = conn.Exec(sql,
			clientID,
			threadID,
			count,
		)
		if err != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// SetUnreadByUser update number of unread events of user
func (r *ThreadRepository) SetUnreadByUser(
	ctx context.Context,
	threadID,
	userID uuid.UUID,
	count int,
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
		sql := `INSERT INTO thread_watchers (
            client_id,
            thread_id,
            user_id,
            unread
        ) VALUES ($1, $2, $3, $4)
        ON CONFLICT (client_id, thread_id, user_id)
        DO UPDATE SET
            unread = thread_watchers.unread + EXCLUDED.unread`
		_, err = conn.Exec(sql,
			clientID,
			threadID,
			userID,
			count,
		)
		if err != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
