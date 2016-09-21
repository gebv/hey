package postgres

import (
	"context"
	"hey/core"
	"hey/storage/postgres"
	"hey/utils"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestHeyService_simple_workspace(t *testing.T) {
	repoThreads := &postgres.ThreadRepository{}
	repoChannels := &postgres.ChannelRepository{}
	repoEvents := &postgres.EventRepository{}

	var clientID = uuid.NewV4()
	var owners = []uuid.UUID{
		uuid.NewV4(),
	}

	ctx := context.WithValue(context.Background(), "ClientID", clientID)

	hey := core.NewHeyService(
		db,
		repoEvents,
		repoThreads,
		repoChannels,
	)

	channelID, threadID, err := hey.CreateChannel(
		ctx,
		owners,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, channelID, uuid.Nil)
	assert.NotEqual(t, threadID, uuid.Nil)

	// check root thread
	ctx = context.WithValue(ctx, "__conn", db)

	thread, err := repoThreads.FindThread(ctx, threadID)
	assert.NoError(t, err)
	assert.Equal(t, thread.ThreadID(), threadID)
	assert.Equal(t, thread.ChannelID(), channelID)

	// check in the database

	var gotChannelID,
		gotClientID,
		gotRootThreadID uuid.UUID

	var gotOWners utils.UUIDS

	err = db.QueryRow(`SELECT
		channel_id,
		client_id,
		owners,
		root_thread_id
	FROM channels 
	WHERE client_id = $1 AND channel_id = $2`,
		clientID,
		channelID,
	).Scan(
		&gotChannelID,
		&gotClientID,
		&gotOWners,
		&gotRootThreadID,
	)
	assert.NoError(t, err)
	assert.Equal(t, gotChannelID, channelID)
	assert.Equal(t, gotClientID, clientID)
	assert.Equal(t, gotRootThreadID, threadID)
	assert.EqualValues(t, gotOWners, owners)

	var gotThreadID,
		gotRelatedEventID,
		gotParentThreadID uuid.UUID

	err = db.QueryRow(`SELECT
		thread_id,
		client_id,
		channel_id,
		owners,
		related_event_id,
		parent_thread_id
	FROM threads 
	WHERE client_id = $1 AND thread_id = $2`,
		clientID,
		threadID,
	).Scan(
		&gotThreadID,
		&gotClientID,
		&gotChannelID,
		&gotOWners,
		&gotRelatedEventID,
		&gotParentThreadID,
	)
	assert.NoError(t, err)
	assert.Equal(t, gotThreadID, threadID)
	assert.Equal(t, gotClientID, clientID)
	assert.Equal(t, gotChannelID, channelID)
	assert.EqualValues(t, gotOWners, owners)
	assert.Equal(t, gotRelatedEventID, uuid.Nil)
	assert.Equal(t, gotParentThreadID, uuid.Nil)
}

func TestHeyService_simple_NodalEvent(t *testing.T) {
	repoThreads := &postgres.ThreadRepository{}
	repoChannels := &postgres.ChannelRepository{}
	repoEvents := &postgres.EventRepository{}

	var clientID = uuid.NewV4()
	var creatorID = uuid.NewV4()
	var owners = []uuid.UUID{
		creatorID,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "ClientID", clientID)
	// ctx = context.WithValue(context.Background(), "__conn", db)

	hey := core.NewHeyService(
		db,
		repoEvents,
		repoThreads,
		repoChannels,
	)

	channelID, threadID, err := hey.CreateChannel(
		ctx,
		owners,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, channelID, uuid.Nil)
	assert.NotEqual(t, threadID, uuid.Nil)

	// Create event
	branchThreadID, nodalEventID, err := hey.CreateNodalEvent(
		ctx,
		threadID,
		owners,
		creatorID,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, branchThreadID, uuid.Nil)
	assert.NotEqual(t, nodalEventID, uuid.Nil)

	// check in the database

	var gotThreadID,
		gotClientID,
		gotChannelID,
		gotRelatedEventID,
		gotParentThreadID uuid.UUID
	var gotOwnerIDs utils.UUIDS

	err = db.QueryRow(`SELECT
		thread_id,
		client_id,
		channel_id,
		owners,
		related_event_id,
		parent_thread_id
	FROM threads 
	WHERE client_id = $1 AND thread_id = $2`,
		clientID,
		branchThreadID,
	).Scan(
		&gotThreadID,
		&gotClientID,
		&gotChannelID,
		&gotOwnerIDs,
		&gotRelatedEventID,
		&gotParentThreadID,
	)
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, gotThreadID, branchThreadID)
	assert.Equal(t, gotClientID, clientID)
	assert.Equal(t, gotChannelID, channelID)
	assert.EqualValues(t, gotOwnerIDs, owners)
	assert.Equal(t, gotRelatedEventID, nodalEventID)
	assert.Equal(t, gotParentThreadID, threadID)

	var gotEventID,
		gotCreatorID,
		gotParentEventID,
		gotBranchThreadID uuid.UUID

	err = db.QueryRow(`SELECT
		event_id,
		client_id,
		thread_id,
		channel_id,
		creator,
		parent_thread_id,
		parent_event_id,
		branch_thread_id
	FROM events
	WHERE client_id = $1 AND event_id = $2`,
		clientID,
		nodalEventID,
	).Scan(
		&gotEventID,
		&gotClientID,
		&gotThreadID,
		&gotChannelID,
		&gotCreatorID,
		&gotParentThreadID,
		&gotParentEventID,
		&gotBranchThreadID,
	)
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, gotEventID, nodalEventID)
	assert.Equal(t, gotClientID, clientID)
	assert.Equal(t, gotThreadID, threadID)
	assert.Equal(t, gotChannelID, channelID)
	assert.Equal(t, gotCreatorID, creatorID)
	assert.Equal(t, gotParentThreadID, uuid.Nil)
	assert.Equal(t, gotParentEventID, uuid.Nil)
	assert.Equal(t, gotBranchThreadID, branchThreadID)
}
