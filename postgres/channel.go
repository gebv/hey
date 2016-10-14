package postgres

import (
	"context"
	"time"

	"github.com/gebv/hey"
	"github.com/gebv/hey/utils"
	uuid "github.com/satori/go.uuid"
)

func (s *Service) FindChannelByName(
	ctx context.Context,
	name string,
) (hey.Channel, error) {
	// TODO: testings
	var clientID = ClientIDFromContext(ctx)

	return s.channels.FindChannelByName(clientID, name)
}

func (s *Service) CreateChannelName(
	ctx context.Context,
	name string,
	userIDs []uuid.UUID,
) (uuid.UUID, uuid.UUID, error) {
	var channelID = uuid.NewV4()
	var rootThreadID = uuid.NewV4()
	var clientID = ClientIDFromContext(ctx)

	if err := utils.ValidName(name); err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return s.createChannel(
		ctx,
		clientID,
		channelID,
		rootThreadID,
		name,
		userIDs,
	)
}

func (s *Service) CreateChannel(
	ctx context.Context,
	userIDs []uuid.UUID,
) (uuid.UUID, uuid.UUID, error) {
	var channelID = uuid.NewV4()
	var rootThreadID = uuid.NewV4()
	var clientID = ClientIDFromContext(ctx)

	return s.createChannel(
		ctx,
		clientID,
		channelID,
		rootThreadID,
		"",
		userIDs,
	)
}

// private

func (s *Service) createChannel(
	ctx context.Context,
	clientID,
	channelID,
	rootThreadID uuid.UUID,
	name string,
	userIDs []uuid.UUID,
) (uuid.UUID, uuid.UUID, error) {

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

		err = s.channels.CreateChannelWithName(
			tx,
			clientID,
			channelID,
			name,
			rootThreadID,
			userIDs,
		)

		if err != nil {
			s.logger.Println("[ERR]", "create channel", err)
			tx.Rollback()
			return
		}

		err = s.threads.CreateThreadWithName(
			tx,
			clientID,
			rootThreadID,
			name,
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
