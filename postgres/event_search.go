package postgres

import (
	"context"
	"time"

	"github.com/gebv/hey"
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
func (s *Service) FindEvents(
	ctx context.Context,
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
