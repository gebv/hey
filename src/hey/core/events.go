package core

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (s *HeyService) CreateUserWorkspace(
	ctx context.Context,
	userIDStr string,
) error {
	/*
	   - channel
	   - root thread
	   - timeline thread
	*/

	var channelID = uuid.NewV4()
	var rootThreadID = uuid.NewV4()
	var timeLineHubEventID = uuid.NewV4()
	var timeLineThreadID = uuid.NewV4()
	var userID = uuid.FromStringOrNil(userIDStr)

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	go func() {
		var err error
		defer func() {
			done <- err

			s.channelsRepo.CreateChannel(
				ctx,
				channelID,
				rootThreadID,
				[]uuid.UUID{userID},
			)
			s.threadsRepo.CreateThread(
				ctx,
				channelID,
				rootThreadID,
				[]uuid.UUID{userID},
			)
		}()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
