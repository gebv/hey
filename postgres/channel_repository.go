package postgres

import (
	"time"

	"github.com/gebv/hey/utils"
	"github.com/satori/go.uuid"

	pg "gopkg.in/jackc/pgx.v2"
)

type ChannelRepository struct {
	db *pg.ConnPool
}

// CreateChannel create new channel
// waiting in the context of the client ID (key name 'ClientID')
func (r *ChannelRepository) CreateChannel(
	tx *pg.Tx,
	clientID,
	channelID,
	rootThreadID uuid.UUID,
	owners []uuid.UUID,
) error {
	if tx == nil {
		return ErrWantTx
	}
	_, err := tx.Exec(`INSERT INTO channels (
            channel_id,
            client_id,
            owners,
            root_thread_id,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6)`,
		channelID,
		clientID,
		utils.UUIDSFrom(owners),
		rootThreadID,
		time.Now(),
		time.Now(),
	)
	return err
}

// func (r *ChannelRepository) AddCountEvents(
// 	ctx context.Context,
// 	channelID uuid.UUID,
// 	count int,
// ) error {
// 	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
// 	done := make(chan error, 1)
// 	defer func() {
// 		cancel()
// 		close(done)
// 	}()

// 	clientID := r.clientIDFromContext(ctx)
// 	// conn := storage.FromContext(ctx)

// 	go func() {
// 		var err error
// 		defer func() {
// 			done <- err
// 		}()

// 		sql := `INSERT INTO channel_counters (
//             client_id,
//             channel_id,
//             counter_events
//         ) VALUES ($1, $2, $3)
//         ON CONFLICT (client_id, channel_id)
//         DO UPDATE SET
//             counter_events = channel_counters.counter_events + EXCLUDED.counter_events`
// 		_, err = conn.Exec(sql,
// 			clientID,
// 			channelID,
// 			count,
// 		)
// 		if err != nil {
// 			return
// 		}
// 	}()

// 	select {
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	case err := <-done:
// 		return err
// 	}
// }

// // SetUnreadByUser update number of unread events of user
// func (r *ChannelRepository) SetUnreadByUser(
// 	ctx context.Context,
// 	channelID,
// 	userID uuid.UUID,
// 	count int,
// ) error {
// 	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
// 	done := make(chan error, 1)
// 	defer func() {
// 		cancel()
// 		close(done)
// 	}()

// 	clientID := r.clientIDFromContext(ctx)
// 	// conn := storage.FromContext(ctx)

// 	go func() {
// 		var err error
// 		defer func() {
// 			done <- err
// 		}()
// 		sql := `INSERT INTO channel_watchers (
//             client_id,
//             channel_id,
//             user_id,
//             unread
//         ) VALUES ($1, $2, $3, $4)
//         ON CONFLICT (client_id, channel_id, user_id)
//         DO UPDATE SET
//             unread = channel_watchers.unread + EXCLUDED.unread`
// 		_, err = conn.Exec(sql,
// 			clientID,
// 			channelID,
// 			userID,
// 			count,
// 		)
// 		if err != nil {
// 			return
// 		}
// 	}()

// 	select {
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	case err := <-done:
// 		return err
// 	}
// }
