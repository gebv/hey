package postgres

import (
	"context"
	"hey/storage"
	"hey/utils"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestChannelRepository_simple_RelatedEventsAndCountEvents(t *testing.T) {
	var clientID = uuid.NewV4()

	tx, err := db.(storage.BeginTX).Begin()
	assert.NoError(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "__conn", tx)
	ctx = context.WithValue(ctx, clientIDContextKey, clientID)

	repo := &ChannelRepository{}
	var creatorID,
		rootThradID,
		channelID = uuid.NewV4(), uuid.NewV4(), uuid.NewV4()
	var someOwnerID = uuid.NewV4()

	err = repo.CreateChannel(
		ctx,
		channelID,
		rootThradID,
		[]uuid.UUID{creatorID, someOwnerID},
	)
	assert.NoError(t, err, "Create channel")

	err = repo.AddCountEvents(ctx, channelID, 1)
	assert.NoError(t, err)
	err = repo.AddCountEvents(ctx, channelID, 1)
	assert.NoError(t, err)
	err = repo.AddCountEvents(ctx, channelID, 1)
	assert.NoError(t, err)

	if err != nil {
		tx.Rollback()
		t.Fatal("TX", err)
	}

	err = tx.Commit()
	assert.NoError(t, err)

	if err != nil {
		t.FailNow()
	}

	// check in the database
	var gotCount int64

	err = db.QueryRow(`SELECT
		counter_events 
	FROM channel_counters WHERE 
    client_id = $1 AND channel_id = $2`,
		clientID,
		channelID,
	).Scan(
		&gotCount,
	)
	assert.NoError(t, err)
	assert.EqualValues(t, gotCount, 3)

	// update count events

	tx, err = db.(storage.BeginTX).Begin()
	assert.NoError(t, err)

	ctx = context.WithValue(ctx, "__conn", tx)
	var userID = uuid.NewV4()

	err = repo.SetUnreadByUser(
		ctx,
		channelID,
		userID,
		20,
	)
	assert.NoError(t, err)

	if err != nil {
		tx.Rollback()
		t.Fatal("TX", err)
	}

	err = tx.Commit()
	assert.NoError(t, err)

	if err != nil {
		t.FailNow()
	}

	// check in the database

	err = db.QueryRow(`SELECT
		unread 
	FROM channel_watchers WHERE 
    client_id = $1 AND channel_id = $2`,
		clientID,
		channelID,
	).Scan(
		&gotCount,
	)
	assert.NoError(t, err)
	assert.EqualValues(t, gotCount, 20)
}

func TestChannelRepository_simple(t *testing.T) {
	var clientID = uuid.NewV4()

	tx, err := db.(storage.BeginTX).Begin()
	assert.NoError(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "__conn", tx)
	ctx = context.WithValue(ctx, clientIDContextKey, clientID)

	repo := &ChannelRepository{}
	var creatorID,
		rootThradID,
		channelID = uuid.NewV4(), uuid.NewV4(), uuid.NewV4()
	var someOwnerID = uuid.NewV4()

	err = repo.CreateChannel(
		ctx,
		channelID,                           // channe ID
		rootThradID,                         // root thread ID
		[]uuid.UUID{creatorID, someOwnerID}, // creator ID (ref. to users)
	)
	assert.NoError(t, err, "Create channel")

	if err != nil {
		tx.Rollback()
		t.Fatal("TX", err)
	}

	err = tx.Commit()
	assert.NoError(t, err)

	if err != nil {
		t.FailNow()
	}

	// check in the database

	var gotChannelID, gotRootThreadID, gotClientID uuid.UUID
	var gotOwners utils.UUIDS
	var owners = utils.UUIDS{creatorID, someOwnerID}

	err = db.QueryRow(`SELECT
		channel_id,
		client_id,
		owners,
		root_thread_id 
	FROM channels WHERE channel_id = $1`, channelID).
		Scan(
			&gotChannelID,
			&gotClientID,
			&gotOwners,
			&gotRootThreadID,
		)
	assert.NoError(t, err)
	assert.Equal(t, gotChannelID, channelID)
	assert.Equal(t, gotRootThreadID, rootThradID)
	assert.Equal(t, gotClientID, clientID)
	assert.Equal(t, gotOwners, owners)
}
