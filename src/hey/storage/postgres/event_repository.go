package postgres

import (
	"context"
	"hey/core/interfaces"
	"hey/storage"
	"hey/utils"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	PerPageMax                       = 100
	PerPageDefault                   = 25
	TimeoutDefault     time.Duration = 50
	FindTimeoutDefault time.Duration = TimeoutDefault * 5
)

type EventRepository struct {
}

func (r *EventRepository) clientIDFromContext(ctx context.Context) uuid.UUID {
	return ClientIDFromContext(ctx)
}

// Event get one event by ID
func (r *EventRepository) Event(
	ctx context.Context,
	eventID uuid.UUID,
) (interfaces.Event, error) {
	// TODO: need test
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	clientID := r.clientIDFromContext(ctx)
	conn := storage.FromContext(ctx)
	event := &event{}

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		conn.QueryRow(`SELECT
			event_id,
			client_id,
			thread_id,
			channel_id,
			creator,
			data,
			parent_thread_id,
			parent_event_id,
			branch_thread_id,
			created_at,
			updated_at 
		FROM events 
		WHERE client_id = $1 AND event_id = $2`,
			clientID,
			eventID,
		).
			Scan(
				&event.eventID,
				&event.clientID,
				&event.threadID,
				&event.channelID,
				&event.creatorID,
				&event.data,
				&event.parentThreadID,
				&event.parentEventID,
				&event.branchThreadID,
				&event.createdAt,
				&event.updatedAt,
			)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-done:
		return event, err
	}
}

// CreateEvent create new event
// waiting in the context of the client ID, channel ID, linked parent
// event and thread IDs
func (r *EventRepository) CreateEvent(
	ctx context.Context,
	eventID,
	threadID,
	channelID,
	creatorID,
	parentThreadID,
	parentEventID,
	branchThreadID uuid.UUID,
	data []byte,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
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
		_, err = conn.Exec(`INSERT INTO events (
            event_id,
            client_id,
            thread_id,
            channel_id,

            creator,

            data,

            parent_thread_id,
            parent_event_id,
            branch_thread_id,

            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			eventID,
			clientID,
			threadID,
			channelID,
			creatorID,
			data, // data
			parentThreadID,
			parentEventID,
			branchThreadID, // branch thread id
			time.Now(),
			time.Now(),
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (r *EventRepository) CreateThreadline(
	ctx context.Context,
	channelID,
	threadID,
	eventID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
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
		_, err = conn.Exec(`INSERT INTO threadline (
            client_id,
            channel_id,
            thread_id,
            event_id,
            created_at
        ) VALUES ($1, $2, $3, $4, $5)`,
			clientID,
			channelID,
			threadID,
			eventID,
			time.Now(),
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// SetBranchThreadID updates the branch thread ID
func (r *EventRepository) SetBranchThreadID(
	ctx context.Context,
	eventID,
	branchThreadID uuid.UUID,
) error {
	// TODO: need test

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
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

		_, err = conn.Exec(`UPDATE events SET
			branch_thread_id = $1,
			updated_at = $2
			WHERE 
				client_id = $3 AND 
				event_id = $4 AND
				branch_thread_id = $5`,
			branchThreadID,
			time.Now(),
			clientID,
			eventID,
			uuid.Nil,
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (r *EventRepository) Threadline(
	ctx context.Context,
	channelID,
	threadID,
	eventID uuid.UUID,
) error {
	// TODO: need test

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
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

		_, err = conn.Exec(`INSERT INTO threadline (
			client_id,
			channel_id,
			thread_id,
			event_id,
			created_at
		) VALUES ($1, $2, $3, $4, $5)`,
			clientID,
			channelID,
			threadID,
			eventID,
			time.Now(),
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (r *EventRepository) DeleteThreadline(
	ctx context.Context,
	eventID uuid.UUID,
) error {
	// TODO: need test

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*TimeoutDefault)
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

		_, err = conn.Exec(`DELETE FROM threadline WHERE 
			client_id = $1 AND 
			event_id = $2`,
			clientID,
			eventID,
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// FindEvents find events
// if the cursor is empty - first page
// else next pages
// if not valid cursor - first page
func (r *EventRepository) FindEvents(
	ctx context.Context,
	threadID uuid.UUID,
	cursorStr string,
	perPage int,
) ([]interfaces.Event, string, error) {
	return r.findEvents(
		ctx,
		threadID,
		cursorStr,
		perPage,
		r,
	)
}

// FindEventsWithProvider search events with custom data provider
// if the cursor is empty - first page
// else next pages
// if not valid cursor - first page
func (r *EventRepository) FindEventsWithProvider(
	ctx context.Context,
	threadID uuid.UUID,
	cursorStr string,
	perPage int,
	provider interfaces.EventProvider,
) ([]interfaces.Event, string, error) {
	return r.findEvents(
		ctx,
		threadID,
		cursorStr,
		perPage,
		provider,
	)
}

func (r *EventRepository) findEvents(
	ctx context.Context,
	threadID uuid.UUID,
	cursorStr string,
	perPage int,
	eventProvider interfaces.EventProvider,
) ([]interfaces.Event, string, error) {
	// TODO: need test

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*FindTimeoutDefault)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	if eventProvider == nil {
		eventProvider = r
	}

	cursor := utils.NewCursorEvents(cursorStr)
	lastEventID := cursor.LastEventID()
	lastEventCreatedAt := cursor.LastCreatedAt()
	cursor = utils.EmptyCursorEvents // reset cursor

	clientID := r.clientIDFromContext(ctx)
	conn := storage.FromContext(ctx)

	if perPage <= 0 || perPage > PerPageMax {
		perPage = PerPageDefault
	}

	// var events = make([]interfaces.Event, perPage)
	var events = []interfaces.Event{}

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		// args for first page
		var args = []interface{}{
			clientID,
			threadID,
			perPage,
		}

		sqlFirstPage := `SELECT
			t.event_id,
			t.created_at
		FROM threadline as t
			WHERE t.client_id = $1 AND t.thread_id = $2
			ORDER BY t.created_at DESC, t.event_id ASC LIMIT $3`
		sqlNextPages := `SELECT
			t.event_id,
			t.created_at
		FROM threadline as t
			WHERE t.client_id = $1 AND t.thread_id = $2 AND
			(t.created_at, t.event_id) < ($3, $4)
			ORDER BY t.created_at DESC, t.event_id ASC LIMIT $5`

		sql := sqlFirstPage

		if !uuid.Equal(lastEventID, uuid.Nil) {
			sql = sqlNextPages
			args = []interface{}{
				clientID,
				threadID,
				lastEventCreatedAt,
				lastEventID,
				perPage,
			}
		}

		rows, err := conn.Query(sql, args...)

		if err != nil {
			return
		}

		var _event interfaces.Event
		var _eventID uuid.UUID
		var _createdAt time.Time

		for rows.Next() {

			if err = rows.Scan(
				&_eventID,
				&_createdAt,
			); err != nil {
				println("[ERR] error scan row", err.Error())
				continue
			}

			_event, err = eventProvider.Event(
				ctx,
				_eventID,
			)

			if err != nil {
				println(
					"[ERR] error find event by ID",
					_eventID.String(),
					err.Error())
				continue
			}

			events = append(events, _event)
		}

		if _event != nil {
			cursor = utils.NewCursorFromSource(
				_event.EventID(),
				_event.CreatedAt(),
			)
		}

	}()

	select {
	case <-ctx.Done():
		return events, "", ctx.Err()
	case err := <-done:
		return events, cursor.String(), err
	}
}
