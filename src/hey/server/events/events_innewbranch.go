package events

import (
	"bytes"
	"context"
	"time"
)

// CreateEventInNewBranch create a new event in a new branch
func (e *EventService) CreateEventInNewBranch(
	ctx context.Context,
	data *bytes.Buffer,
	relatedEventID string,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	var done = make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	var newThreadID = NewUUID()
	var newEventID = NewUUID()

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		relatedEvent, err := e.eventStore.FindEvent(ctx, relatedEventID)

		if err != nil {
			return
		}

		err = e.threadService.CreateThread(ctx,
			relatedEventID,
			relatedEvent.ThreadID(),
			newThreadID,
		)

		// TODO: updated relatedEventID
		e.addBranchEvent(
			ctx,
			relatedEventID,
			newThreadID,
		)

		return
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-done:
		if err != nil {
			return "", err
		}
	}

	return newEventID,
		e.createEvent(
			ctx,
			data,
			newThreadID,
			newEventID,
		)
}

// ThreadService

type ThreadService interface {
	CreateThread(ctx context.Context,
		relatedEventID,
		parentThreadID,
		threadID string) error
}
