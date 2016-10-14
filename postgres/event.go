package postgres

import (
	"context"
	"time"

	"github.com/gebv/hey/utils"
	uuid "github.com/satori/go.uuid"
)

// CreateEvent create a new event to an existing thread
func (s *Service) CreateEvent(ctx context.Context,
	threadID uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (eventID uuid.UUID, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	var newEventID = uuid.NewV4()
	var clientID = ClientIDFromContext(ctx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		currentThread, err := s.threads.FindThread(
			clientID,
			threadID,
		)

		if err != nil {
			return
		}

		tx, err := s.db.Begin()

		if err != nil {
			return
		}

		// 1. Create event
		// 2. Create threadline

		// 1
		err = s.events.CreateEvent(
			tx,
			clientID,
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

		// 2
		err = s.events.Threadline(
			tx,
			clientID,
			currentThread.ChannelID(),
			threadID,
			newEventID,
		)

		if err != nil {
			return
		}

		err = tx.Commit()
	}()

	select {
	case <-ctx.Done():
		<-done
		return uuid.Nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return uuid.Nil, err
		}
		return newEventID, err
	}
}

// CreateNodalEvent create new nodal event
// waiting ChannelID from context
func (s *Service) CreateNodalEvent(
	ctx context.Context,
	threadID uuid.UUID,
	owners []uuid.UUID,
	creatorID uuid.UUID,
) (uuid.UUID, uuid.UUID, error) {
	return s.createNodalEvent(
		ctx,
		threadID.String(),
		threadID,
		owners,
		creatorID,
		[]byte{},
	)
}

// CreateNodalEvent create new nodal event
// waiting ChannelID from context
func (s *Service) CreateNodalEventWithThreadName(
	ctx context.Context,
	threadName string,
	threadID uuid.UUID,
	owners []uuid.UUID,
	creatorID uuid.UUID,
) (uuid.UUID, uuid.UUID, error) {

	if err := utils.ValidName(threadName); err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return s.createNodalEvent(
		ctx,
		threadName,
		threadID,
		owners,
		creatorID,
		[]byte{},
	)
}

// CreateNodalEventWithData create new nodal event
// waiting ChannelID from context
func (s *Service) CreateNodalEventWithData(
	ctx context.Context,
	threadID uuid.UUID,
	owners []uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (uuid.UUID, uuid.UUID, error) {
	return s.createNodalEvent(
		ctx,
		threadID.String(),
		threadID,
		owners,
		creatorID,
		data,
	)
}

// CreateNodalEventWithThreadNameWithData create new nodal event
// waiting ChannelID from context
func (s *Service) CreateNodalEventWithThreadNameWithData(
	ctx context.Context,
	threadName string,
	threadID uuid.UUID,
	owners []uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (uuid.UUID, uuid.UUID, error) {

	if err := utils.ValidName(threadName); err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return s.createNodalEvent(
		ctx,
		threadName,
		threadID,
		owners,
		creatorID,
		data,
	)
}

func (s *Service) CreateNewBranchEventWithThreadName(
	ctx context.Context,
	threadName string,
	threadID uuid.UUID,
	relatedEventID uuid.UUID, //
	owners []uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (uuid.UUID, uuid.UUID, error) {

	if err := utils.ValidName(threadName); err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return s.createNewBranchEvent(
		ctx,
		threadName,
		threadID,
		relatedEventID,
		owners,
		creatorID,
		data,
	)
}

// CreateNewBranchEvent create a new event in branch
// if the event already has the branch - error
func (s *Service) CreateNewBranchEvent(
	ctx context.Context,
	threadID uuid.UUID,
	relatedEventID uuid.UUID, //
	owners []uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (uuid.UUID, uuid.UUID, error) {
	return s.createNewBranchEvent(
		ctx,
		"",
		threadID,
		relatedEventID,
		owners,
		creatorID,
		data,
	)
}

// private

func (s *Service) createNodalEvent(
	ctx context.Context,
	threadName string,
	threadID uuid.UUID,
	owners []uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (uuid.UUID, uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	newEventID := uuid.NewV4()
	newThreadID := uuid.NewV4()
	clientID := ClientIDFromContext(ctx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		currentThread, err := s.threads.FindThread(
			clientID,
			threadID,
		)

		if err != nil {
			return
		}

		tx, err := s.db.Begin()

		if err != nil {
			return
		}

		err = s.events.CreateEvent(
			tx,
			clientID,
			newEventID,
			threadID,
			currentThread.ChannelID(),
			creatorID,
			currentThread.ParentThreadID(), // parent thread ID
			currentThread.RelatedEventID(), // parent event ID
			newThreadID,                    // branch thread id
			data,
		)

		if err != nil {
			return
		}

		s.events.CreateThreadline(
			tx,
			clientID,
			currentThread.ChannelID(),
			threadID,
			newEventID,
		)

		// branch thread

		err = s.threads.CreateThreadWithName(
			tx,
			clientID,
			newThreadID,
			threadName,
			currentThread.ChannelID(), // TODO: get channelID
			newEventID,                // related event ID
			threadID,                  // parent thread ID
			owners,
		)

		if err != nil {
			return
		}

		err = tx.Commit()
	}()

	select {
	case <-ctx.Done():
		<-done
		return uuid.Nil, uuid.Nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return uuid.Nil, uuid.Nil, err
		}
		return newThreadID, newEventID, err
	}
}

func (s *Service) createNewBranchEvent(
	ctx context.Context,
	threadName string,
	threadID uuid.UUID,
	relatedEventID uuid.UUID, //
	owners []uuid.UUID,
	creatorID uuid.UUID,
	data []byte,
) (uuid.UUID, uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	newThreadID := uuid.NewV4() // branch thread id
	newEventID := uuid.NewV4()  // event in new branch
	clientID := ClientIDFromContext(ctx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		// 1. Find threadID
		// 2. Set branch thread id for related event
		// 3. Create branch thread
		// 4. Create event

		currentThread, err := s.threads.FindThread(
			clientID,
			threadID,
		)

		if err != nil {
			return
		}

		tx, err := s.db.Begin()

		if err != nil {
			return
		}

		err = s.events.SetBranchThreadID(
			tx,
			clientID,
			relatedEventID,
			newThreadID,
		)

		if err != nil {
			return
		}

		err = s.threads.CreateThreadWithName(
			tx,
			clientID,
			newThreadID,
			threadName,
			currentThread.ChannelID(), // TODO: get channelID
			relatedEventID,            // related event ID
			currentThread.ThreadID(),  // parent thread ID
			owners,
		)

		if err != nil {
			return
		}

		err = s.events.CreateEvent(
			tx,
			clientID,
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

		err = tx.Commit()
	}()

	select {
	case <-ctx.Done():
		<-done
		return uuid.Nil, uuid.Nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return uuid.Nil, uuid.Nil, err
		}
		return newThreadID, newEventID, err
	}
}
