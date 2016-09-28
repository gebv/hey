package postgres

import (
	"time"

	"github.com/gebv/hey"
	"github.com/gebv/hey/utils"

	uuid "github.com/satori/go.uuid"
	pg "gopkg.in/jackc/pgx.v2"
)

type EventRepository struct {
	db *pg.ConnPool
}

// FindEvent get one event by ID
func (r *EventRepository) FindEvent(
	clientID,
	eventID uuid.UUID,
) (hey.Event, error) {

	event := event{}

	err := r.db.QueryRow(`SELECT
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
	return &event, err
}

// CreateEvent create new event
// waiting in the context of the client ID, channel ID, linked parent
// event and thread IDs
func (r *EventRepository) CreateEvent(
	tx *pg.Tx,
	clientID,
	eventID,
	threadID,
	channelID,
	creatorID,
	parentThreadID,
	parentEventID,
	branchThreadID uuid.UUID,
	data []byte,
) error {
	_, err := tx.Exec(`INSERT INTO events (
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
	return err
}

func (r *EventRepository) CreateThreadline(
	tx *pg.Tx,
	clientID,
	channelID,
	threadID,
	eventID uuid.UUID,
) error {
	_, err := tx.Exec(`INSERT INTO threadline (
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

	return err
}

// SetBranchThreadID updates the branch thread ID
func (r *EventRepository) SetBranchThreadID(
	tx *pg.Tx,
	clientID,
	eventID,
	branchThreadID uuid.UUID,
) error {
	_, err := tx.Exec(`UPDATE events SET
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
	return err
}

func (r *EventRepository) Threadline(
	tx *pg.Tx,
	clientID,
	channelID,
	threadID,
	eventID uuid.UUID,
) error {
	_, err := tx.Exec(`INSERT INTO threadline (
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
	return err
}

func (r *EventRepository) DeleteThreadline(
	tx *pg.Tx,
	clientID,
	eventID uuid.UUID,
) error {
	_, err := tx.Exec(`DELETE FROM threadline WHERE 
			client_id = $1 AND 
			event_id = $2`,
		clientID,
		eventID,
	)

	return err
}

// FindEvents find events
// if the cursor is empty - first page
// else next pages
// if not valid cursor - first page
func (r *EventRepository) FindEvents(
	clientID,
	watcherID,
	threadID uuid.UUID,
	cursorStr string,
	perPage int,
) ([]hey.Event, string, error) {
	return r.findEvents(
		clientID,
		watcherID,
		threadID,
		cursorStr,
		perPage,
	)
}

func (r *EventRepository) findEvents(
	clientID,
	watcherID,
	threadID uuid.UUID,
	cursorStr string,
	perPage int,
) ([]hey.Event, string, error) {

	cursor := utils.NewCursorEvents(cursorStr)
	lastEventID := cursor.LastEventID()
	lastEventCreatedAt := cursor.LastCreatedAt()
	cursor = utils.EmptyCursorEvents // reset cursor

	if perPage <= 0 || perPage > PerPageMax {
		perPage = PerPageDefault
	}

	var events = make([]hey.Event, perPage)
	// var events []hey.Event

	// find

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

	rows, err := r.db.Query(sql, args...)

	if err != nil {
		return events, "", err
	}

	var _event hey.Event
	var _eventID uuid.UUID
	var _createdAt time.Time
	var i = 0

	for rows.Next() {

		if err = rows.Scan(
			&_eventID,
			&_createdAt,
		); err != nil {
			println("[ERR] error scan row", err.Error())
			continue
		}

		_event, err = r.FindEvent(
			clientID,
			_eventID,
		)

		if err != nil {
			println(
				"[ERR] error find event by ID",
				_eventID.String(),
				err.Error())
			continue
		}

		events[i] = _event
		i++
	}

	if _event != nil {
		cursor = utils.NewCursorFromSource(
			_event.EventID(),
			_event.CreatedAt(),
		)
	}

	if i < perPage-1 {
		// обрезать хвост у массива событий, если кол-во найденных меньше
		// чем ожидаемое кол-во
		events = events[:i]
	}

	return events, cursor.String(), err
}
