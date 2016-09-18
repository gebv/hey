package events

import (
	"bytes"
	"context"
	"hey/storage"
	"time"
)

// type EventDTO struct {
// 	EventID string

// 	ClientID  int64  // ref to clients
// 	ChannelID int64  // ref to channels
// 	ThreadID  string // ref to threads

// 	ParentThreadID string // ref to threads
// 	ParentEventID  int64  // ref to events
// 	BranchThreadID string // ref to threads

// 	Creator int64  // ref to users, owner of the event
// 	Data    []byte // payload

// 	Flags []string // internal flags

// 	ExtProps map[string]interface{} // external properties
// 	ExtFlags []string               // external flags

// 	IsRemoved bool
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

type EventService struct {
	ac            AccessControl
	threadService ThreadService
	eventStore    EventStore
	threadStore   ThreadStore
}

// CreateEvent create a new event
func (e *EventService) CreateEvent(
	ctx context.Context,
	data *bytes.Buffer,
) (string, error) {
	eventID := NewUUID()
	return eventID,
		e.createEvent(
			ctx,
			data,
			e.ThreadIDFromContext(ctx),
			eventID,
		)
}

func (e *EventService) createEvent(
	ctx context.Context,
	data *bytes.Buffer,
	threadID, newEventID string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	var done = make(chan error, 1)
	var allowed chan error

	defer func() {
		cancel()
		close(done)
		close(allowed)
	}()

	var creatorID = e.CreateIDFromContext(ctx)
	allowed = e.AllowedCreateEvent(
		ctx,
		creatorID,
		threadID,
	)

	go func() {
		if err := <-allowed; err != nil {
			done <- err
			return
		}

		done <- e.insertEvent(
			ctx,
			data,
			threadID,
			newEventID,
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// insertEvent add a new event
func (e *EventService) insertEvent(
	ctx context.Context,
	payload *bytes.Buffer,
	eventID,
	threadID string,
) error {
	channelID := e.ChannelIDFromContext(ctx)
	clientID := e.ClientIDFromContext(ctx)
	creatorID := e.CreateIDFromContext(ctx)

	currentThread, err := e.threadStore.FindThread(ctx, channelID, threadID)
	if err != nil {
		return err
	}

	var sql = `INSERT INTO events (
		event_id, 
		client_id,
		channel_id,
		thread_id,

		creator,

		data,

--		props,

		parent_thread_id,
		parent_event_id,
--		branch_thread_id,
		
--		flags,
--		ext_flags,

--		is_removed,
		created_at,
 		updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		-- , $11, $12, $13, $14, $15
	)`

	_, err = e.conn(ctx).
		Exec(sql,
			eventID,
			clientID,
			channelID,
			threadID,
			creatorID,
			payload.Bytes(),
			// map[string]interface{}{}, // TODO: props
			currentThread.ParentThreadID(),
			currentThread.RelatedEventID(), // as parent event id
			// []string{},
			// []string{},
			// false, // default false
			time.Now(),
			time.Now(),
		)
	return err
}

// addBranchEvent add a branch to the event
func (e *EventService) addBranchEvent(
	ctx context.Context,
	eventID, branchThreadID string,
) error {
	sql := `UPDATE events SET 
		branch_thread_id = $1 AND 
		updated_at = $2 
	WHERE event_id = $3`

	_, err := e.conn(ctx).
		Exec(sql,
			branchThreadID,
			time.Now(),
			eventID,
		)

	return err
}

func (e *EventService) conn(ctx context.Context) storage.DB {
	if conn, ok := ctx.Value("__dbconn").(storage.DB); ok {
		return conn
	}
	return nil
}

func (e *EventService) ClientIDFromContext(ctx context.Context) int64 {
	return ctx.Value("client_id").(int64)
}

func (e *EventService) ChannelIDFromContext(ctx context.Context) int64 {
	return ctx.Value("channel_id").(int64)
}

func (e *EventService) ThreadIDFromContext(ctx context.Context) string {
	return ctx.Value("thread_id").(string)
}

func (e *EventService) CreateIDFromContext(ctx context.Context) int64 {
	return ctx.Value("creator_id").(int64)
}

// channels

type ChannelStore interface {
	FindChannel(ctx context.Context, channelID int64) (Channel, error)
}

type Channel interface {
	RootThreadID() string
}

// threads

type ThreadStore interface {
	FindThread(ctx context.Context, channelID int64, threadID string) (Thread, error)
}

type Thread interface {
	ChannelID() int64
	ParentThreadID() string
	RelatedEventID() string
}

// events

type EventStore interface {
	FindEvent(ctx context.Context, eventID string) (Event, error)
}

type Event interface {
	ThreadID() string
	ChannelID() int64
}
