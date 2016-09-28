package postgres

// import (
// 	"context"
// 	"hey/storage"
// 	"hey/utils"
// 	"testing"

// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// )

// func TestThreadRepository_simple_RelatedEventsAndCountEvents(t *testing.T) {
// 	var clientID = uuid.NewV4()
// 	var threadID = uuid.NewV4()

// 	tx, err := db.(storage.BeginTX).Begin()
// 	assert.NoError(t, err)

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "__conn", tx)
// 	ctx = context.WithValue(ctx, clientIDContextKey, clientID)

// 	repo := &ThreadRepository{}
// 	// create related entities
// 	err = repo.CreateThread(
// 		ctx,
// 		threadID,
// 		uuid.NewV4(),
// 		uuid.NewV4(),
// 		uuid.NewV4(),
// 		[]uuid.UUID{uuid.NewV4(), uuid.NewV4()},
// 	)
// 	assert.NoError(t, err)

// 	err = repo.AddCountEvents(ctx, threadID, 1)
// 	assert.NoError(t, err)
// 	err = repo.AddCountEvents(ctx, threadID, 1)
// 	assert.NoError(t, err)
// 	err = repo.AddCountEvents(ctx, threadID, 1)
// 	assert.NoError(t, err)

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
// 	var gotCount int64

// 	err = db.QueryRow(`SELECT
// 		counter_events
// 	FROM thread_counters WHERE
//     client_id = $1 AND thread_id = $2`,
// 		clientID,
// 		threadID,
// 	).Scan(
// 		&gotCount,
// 	)
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, gotCount, 3)

// 	// update count events

// 	tx, err = db.(storage.BeginTX).Begin()
// 	assert.NoError(t, err)

// 	ctx = context.WithValue(ctx, "__conn", tx)
// 	var userID = uuid.NewV4()

// 	err = repo.SetUnreadByUser(
// 		ctx,
// 		threadID,
// 		userID,
// 		20,
// 	)
// 	assert.NoError(t, err)

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

// 	err = db.QueryRow(`SELECT
// 		unread
// 	FROM thread_watchers WHERE
//     client_id = $1 AND thread_id = $2`,
// 		clientID,
// 		threadID,
// 	).Scan(
// 		&gotCount,
// 	)
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, gotCount, 20)
// }

// func TestThreadRepository_simple_create(t *testing.T) {
// 	var clientID = uuid.NewV4()
// 	var relatedEventID = uuid.NewV4()
// 	var parentThreadID = uuid.NewV4()

// 	tx, err := db.(storage.BeginTX).Begin()
// 	assert.NoError(t, err)

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "__conn", tx)
// 	ctx = context.WithValue(ctx, clientIDContextKey, clientID)

// 	repo := &ThreadRepository{}
// 	var creatorID,
// 		threadID,
// 		channelID = uuid.NewV4(), uuid.NewV4(), uuid.NewV4()
// 	var someOwnerID = uuid.NewV4()

// 	err = repo.CreateThread(
// 		ctx,
// 		threadID,  // new thread ID
// 		channelID, // channe ID
// 		relatedEventID,
// 		parentThreadID,
// 		[]uuid.UUID{creatorID, someOwnerID}, // creator ID (ref. to users)
// 	)
// 	assert.NoError(t, err, "Create thread")

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

// 	var gotThreadID,
// 		gotChannelID,
// 		gotRelatedEventID,
// 		gotParentThreadID,
// 		gotClientID uuid.UUID

// 	var gotOwners utils.UUIDS
// 	var owners = utils.UUIDS{creatorID, someOwnerID}

// 	err = db.QueryRow(`SELECT
// 		thread_id,
// 		client_id,
// 		channel_id,
// 		owners,
// 		related_event_id,
//         parent_thread_id
// 	FROM threads WHERE thread_id = $1`, threadID).
// 		Scan(
// 			&gotThreadID,
// 			&gotClientID,
// 			&gotChannelID,
// 			&gotOwners,
// 			&gotRelatedEventID,
// 			&gotParentThreadID,
// 		)
// 	assert.NoError(t, err)
// 	assert.Equal(t, gotThreadID, threadID)
// 	assert.Equal(t, gotChannelID, channelID)
// 	assert.Equal(t, gotRelatedEventID, relatedEventID)
// 	assert.Equal(t, gotParentThreadID, parentThreadID)
// 	assert.Equal(t, gotClientID, clientID)
// 	assert.Equal(t, gotOwners, owners)
// }
