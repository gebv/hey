package postgres

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (s *Service) CreateChannel(
	ctx context.Context,
	userIDs []uuid.UUID,
) (uuid.UUID, uuid.UUID, error) {
	var channelID = uuid.NewV4()
	var rootThreadID = uuid.NewV4()
	var clientID = ClientIDFromContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
	done := make(chan error, 1)

	defer func() {
		cancel()
		close(done)
	}()

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		tx, err := s.db.Begin()

		if err != nil {
			return
		}

		err = s.channels.CreateChannel(
			tx,
			clientID,
			channelID,
			rootThreadID,
			userIDs,
		)

		if err != nil {
			s.logger.Println("[ERR]", "create channel", err)
			tx.Rollback()
			return
		}

		err = s.threads.CreateThread(
			tx,
			clientID,
			rootThreadID,
			channelID,
			uuid.Nil,
			uuid.Nil,
			userIDs,
		)

		if err != nil {
			s.logger.Println("[ERR]", "create root thread", err)
			tx.Rollback()
			return
		}

		err = tx.Commit()
		return
	}()

	select {
	case <-ctx.Done():
		<-done // TODO: force done
		return uuid.Nil, uuid.Nil, ctx.Err()
	case err := <-done:
		if err != nil {
			s.logger.Println("[ERR]", "create channel, summary", err)
			return uuid.Nil, uuid.Nil, err
		}
		return channelID, rootThreadID, err
	}
}
