package core

import (
	"context"
	"errors"
	"hey/storage"
	"time"

	uuid "github.com/satori/go.uuid"
)

// CreateChannel создать канал и поток
func (s *HeyService) CreateChannel(
	ctx context.Context,
	userIDs []uuid.UUID,
) (uuid.UUID, uuid.UUID, error) {
	var channelID = uuid.NewV4()
	var rootThreadID = uuid.NewV4()

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	begintx, ok := s.conn.(storage.BeginTX)
	if !ok {
		return uuid.Nil, uuid.Nil, errors.New("only internal transaction")
	}
	tx, err := begintx.Begin()
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	ctx = context.WithValue(ctx, "__conn", tx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		err = s.channelsRepo.CreateChannel(
			ctx,
			channelID,
			rootThreadID,
			userIDs,
		)
		if err != nil {
			return
		}
		err = s.threadsRepo.CreateThread(
			ctx,
			rootThreadID,
			channelID,
			uuid.Nil,
			uuid.Nil,
			userIDs,
		)
		if err != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		tx.Rollback()
		return uuid.Nil, uuid.Nil, ctx.Err()
	case err := <-done:
		if err != nil {
			tx.Rollback()
			return uuid.Nil, uuid.Nil, err
		}
		return channelID, rootThreadID, tx.Commit()
	}
}

// CreateNodalEvent create new nodal event
// waiting ChannelID from context
func (s *HeyService) CreateNodalEvent(
	ctx context.Context,
	threadID uuid.UUID,
	owners []uuid.UUID,
	creatorID uuid.UUID,
) (uuid.UUID, uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	newEventID := uuid.NewV4()
	newThreadID := uuid.NewV4()

	begintx, ok := s.conn.(storage.BeginTX)
	if !ok {
		return uuid.Nil, uuid.Nil, errors.New("only internal transaction")
	}
	tx, err := begintx.Begin()
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	ctx = context.WithValue(ctx, "__conn", tx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		currentThread, err := s.threadsRepo.FindThread(
			ctx,
			threadID,
		)
		if err != nil {
			return
		}

		err = s.eventsRepo.CreateEvent(
			ctx,
			newEventID,
			threadID,
			currentThread.ChannelID(),
			creatorID,
			currentThread.ParentThreadID(), // parent thread ID
			currentThread.RelatedEventID(), // parent event ID
			newThreadID,                    // branch thread id
			[]byte{},
		)

		if err != nil {
			return
		}

		// branch thread

		err = s.threadsRepo.CreateThread(
			ctx,
			newThreadID,
			currentThread.ChannelID(), // TODO: get channelID
			newEventID,                // related event ID
			threadID,                  // parent thread ID
			owners,
		)

		if err != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		tx.Rollback()
		return uuid.Nil, uuid.Nil, ctx.Err()
	case err := <-done:
		if err != nil {
			tx.Rollback()
			return uuid.Nil, uuid.Nil, err
		}
		return newThreadID, newEventID, tx.Commit()
	}
}

// CreateNewBranchEvent create a new event in branch
// if the event already has the branch - error
func (s *HeyService) CreateNewBranchEvent(
	ctx context.Context,
	threadID uuid.UUID,
	relatedEventID uuid.UUID, //
	owners []uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (uuid.UUID, uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	// TODO: check accessory events to the current thread
	// can only from RealtedEvetnID?

	newThreadID := uuid.NewV4() // branch thread id
	newEventID := uuid.NewV4()  // event in new branch

	begintx, ok := s.conn.(storage.BeginTX)
	if !ok {
		return uuid.Nil, uuid.Nil, errors.New("only internal transaction")
	}
	tx, err := begintx.Begin()
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	ctx = context.WithValue(ctx, "__conn", tx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		err = s.eventsRepo.SetBranchThreadID(
			ctx,
			relatedEventID,
			newThreadID,
		)

		// if error then branch is already

		if err != nil {
			return
		}

		currentThread, err := s.threadsRepo.FindThread(
			ctx,
			threadID,
		)
		if err != nil {
			return
		}

		// branch thread

		err = s.threadsRepo.CreateThread(
			ctx,
			newThreadID,
			currentThread.ChannelID(), // TODO: get channelID
			relatedEventID,            // related event ID
			currentThread.ThreadID(),  // parent thread ID
			owners,
		)

		if err != nil {
			return
		}

		err = s.eventsRepo.CreateEvent(
			ctx,
			newEventID,
			newThreadID,
			currentThread.ChannelID(),
			creatorID,
			currentThread.ThreadID(), // parent thread ID
			relatedEventID,           // parent event ID
			uuid.Nil,                 // branch thread id
			data,
		)

		if err != nil {
			return
		}

	}()

	select {
	case <-ctx.Done():
		tx.Rollback()
		return uuid.Nil, uuid.Nil, ctx.Err()
	case err := <-done:
		if err != nil {
			tx.Rollback()
			return uuid.Nil, uuid.Nil, err
		}
		return newThreadID, newEventID, tx.Commit()
	}
}

// CreateEvent create event in existing branch
func (s *HeyService) CreateEvent(
	ctx context.Context,
	threadID uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	newEventID := uuid.NewV4()

	begintx, ok := s.conn.(storage.BeginTX)
	if !ok {
		return uuid.Nil, errors.New("only internal transaction")
	}
	tx, err := begintx.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	ctx = context.WithValue(ctx, "__conn", tx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		currentThread, err := s.threadsRepo.FindThread(
			ctx,
			threadID,
		)
		if err != nil {
			return
		}

		err = s.eventsRepo.CreateEvent(
			ctx,
			newEventID,
			threadID,
			currentThread.ChannelID(),
			creatorID,
			currentThread.ParentThreadID(), // parent thread ID
			currentThread.RelatedEventID(), // parent event ID
			uuid.Nil,                       // branch thread id
			data,
		)

		if err != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		tx.Rollback()
		return uuid.Nil, ctx.Err()
	case err := <-done:
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		return newEventID, tx.Commit()
	}
}
