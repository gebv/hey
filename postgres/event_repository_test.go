package postgres

// import (
// 	"context"
// 	"hey/storage"
// 	"testing"

// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// )

// func TestEventRepository_simple_threadline(t *testing.T) {
// 	var clientID = uuid.NewV4()
// 	var channelID = uuid.NewV4()
// 	var threadID = uuid.NewV4()
// 	var eventID = uuid.NewV4()

// 	tx, err := db.(storage.BeginTX).Begin()
// 	assert.NoError(t, err)

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "__conn", tx)
// 	ctx = context.WithValue(ctx, clientIDContextKey, clientID)

// 	repo := &EventRepository{}
// 	err = repo.CreateThreadline(
// 		ctx,
// 		channelID,
// 		threadID,
// 		eventID,
// 	)
// 	assert.NoError(t, err, "Create threadline")

// 	if err != nil {
// 		tx.Rollback()
// 		t.Fatal("TX", err)
// 	}

// 	err = tx.Commit()
// 	assert.NoError(t, err)

// 	if err != nil {
// 		t.FailNow()
// 	}

// 	// check in the database
// 	var gotClientID,
// 		gotChannelID,
// 		gotThreadID,
// 		gotEventID uuid.UUID

// 	err = db.QueryRow(`SELECT
//         client_id,
//         channel_id,
//         thread_id,
//         event_id
// 	FROM threadline WHERE event_id = $1`, eventID).
// 		Scan(
// 			&gotClientID,
// 			&gotChannelID,
// 			&gotThreadID,
// 			&gotEventID,
// 		)
// 	assert.NoError(t, err)

// 	assert.Equal(t, gotClientID, clientID)
// 	assert.Equal(t, gotChannelID, channelID)
// 	assert.Equal(t, gotThreadID, threadID)
// 	assert.Equal(t, gotEventID, eventID)
// }

// func TestEventRepository_simple(t *testing.T) {
// 	var clientID = uuid.NewV4()
// 	var parentEventID = uuid.NewV4()
// 	var parentThreadID = uuid.NewV4()
// 	var channelID = uuid.NewV4()

// 	tx, err := db.(storage.BeginTX).Begin()
// 	assert.NoError(t, err)

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "__conn", tx)
// 	ctx = context.WithValue(ctx, clientIDContextKey, clientID)

// 	repo := &EventRepository{}
// 	var creatorID,
// 		threadID,
// 		eventID = uuid.NewV4(), uuid.NewV4(), uuid.NewV4()

// 	err = repo.CreateEvent(
// 		ctx,
// 		eventID,
// 		threadID,
// 		channelID,
// 		creatorID,
// 		parentThreadID,
// 		parentEventID,
// 		uuid.Nil,
// 		[]byte("hello"),
// 	)
// 	assert.NoError(t, err, "Create event")

// 	if err != nil {
// 		tx.Rollback()
// 		t.Fatal("TX", err)
// 	}

// 	err = tx.Commit()
// 	assert.NoError(t, err)

// 	if err != nil {
// 		t.FailNow()
// 	}

// 	// check in the database

// 	var gotEventID,
// 		gotClientID,
// 		gotThreadID,
// 		gotChannelID,
// 		gotCreatorID,
// 		gotParentTHreadID,
// 		gotParentEventID uuid.UUID
// 	var gotData []byte

// 	err = db.QueryRow(`SELECT
//         event_id,
//         client_id,
//         thread_id,
//         channel_id,

//         creator,

//         data,

//         parent_thread_id,
//         parent_event_id
// 	FROM events WHERE event_id = $1`, eventID).
// 		Scan(
// 			&gotEventID,
// 			&gotClientID,
// 			&gotThreadID,
// 			&gotChannelID,
// 			&gotCreatorID,
// 			&gotData,
// 			&gotParentTHreadID,
// 			&gotParentEventID,
// 		)
// 	assert.NoError(t, err)
// 	assert.Equal(t, gotEventID, eventID)
// 	assert.Equal(t, gotClientID, clientID)
// 	assert.Equal(t, gotThreadID, threadID)
// 	assert.Equal(t, gotChannelID, channelID)
// 	assert.Equal(t, gotCreatorID, creatorID)
// 	assert.Equal(t, gotData, []byte("hello"))
// 	assert.Equal(t, gotParentTHreadID, parentThreadID)
// 	assert.Equal(t, gotParentEventID, parentEventID)
// }
