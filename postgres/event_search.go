package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/gebv/hey"
	"github.com/gebv/hey/utils"
	uuid "github.com/satori/go.uuid"
)

var (
	_ hey.SearchResult = (*SearchResult)(nil)
)

// SearchResult search result
type SearchResult struct {
	events  []hey.Event
	cursor  string
	hasNext bool
}

func (s SearchResult) Events() []hey.Event {
	return s.events
}

func (s SearchResult) Cursor() string {
	return s.cursor
}

func (s SearchResult) HasNext() bool {
	return s.hasNext
}

// FindEvents find events
// waiting WatcherID (from a user view) from context
func (s *Service) FindEventsByName(
	ctx context.Context,
	watcherID uuid.UUID,
	name string,
	cursorStr string,
	perPage int,
) (hey.SearchResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*FindTimeoutDefault)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	var result hey.SearchResult
	var clientID = ClientIDFromContext(ctx)
	var channel hey.Channel
	var thread hey.Thread

	sn := utils.SpecialName(name)
	if !sn.Valid() {
		return result, errors.New("find events: not valid special name")
	}

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		channel, err = s.channels.FindChannelByName(
			clientID,
			sn.Channel(),
		)

		if err != nil {
			return
		}

		thread, err = s.threads.FindThreadByName(
			clientID,
			channel.ChannelID(),
			sn.Thread(),
		)

		if err != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		<-done
		return result, ctx.Err()
	case err := <-done:
		if err != nil {
			return result, err
		}
		return s.findEvents(
			ctx,
			watcherID,
			thread.ThreadID(),
			cursorStr,
			perPage,
		)
	}

}

// FindEvents find events
// waiting WatcherID (from a user view) from context
func (s *Service) FindEvents(
	ctx context.Context,
	watcherID uuid.UUID,
	threadID uuid.UUID,
	cursorStr string,
	perPage int,
) (hey.SearchResult, error) {
	return s.findEvents(
		ctx,
		watcherID,
		threadID,
		cursorStr,
		perPage,
	)
}

// private

func (s *Service) findEvents(ctx context.Context,
	watcherID uuid.UUID,
	threadID uuid.UUID,
	cursorStr string,
	perPage int,
) (hey.SearchResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*FindTimeoutDefault)
	done := make(chan error, 1)

	defer func() {
		cancel()
		close(done)
	}()

	var result = &SearchResult{}
	var clientID = ClientIDFromContext(ctx)

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		result.events,
			result.cursor,
			err = s.events.
			FindEvents(
				clientID,
				watcherID,
				threadID,
				cursorStr,
				perPage,
			)

	}()

	select {
	case <-ctx.Done():
		<-done // TODO: force done
		return result, ctx.Err()
	case err := <-done:
		if err != nil {
			// s.logger.Println("[ERR]", "create channel, summary", err)
		}
		return result, err
	}
}
