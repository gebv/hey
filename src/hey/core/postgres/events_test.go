package postgres

import (
	"context"
	"hey/core"
	"hey/core/interfaces"
	"hey/storage/postgres"
	"hey/utils"
	"strconv"
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

func TestHeyService_simple_NewEvent(t *testing.T) {
	// db.Exec(`DELETE FROM events;`)

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

	// Create nodal event
	branchThreadID, nodalEventID, err := hey.CreateNodalEvent(
		ctx,
		threadID,
		owners,
		creatorID,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, branchThreadID, uuid.Nil)
	assert.NotEqual(t, nodalEventID, uuid.Nil)

	// Create new event
	newEventID, err := hey.CreateEvent(
		ctx,
		branchThreadID,
		creatorID,
		[]byte("hello machine"),
	)
	assert.NoError(t, err)
	assert.NotEqual(t, newEventID, uuid.Nil)

	// check in the database

	var gotEventID,
		gotClientID,
		gotThreadID,
		gotChannelID,
		gotCreatorID,
		gotParentThreadID,
		gotParetnEventID,
		gotBranchThreadID uuid.UUID
	var gotData []byte

	err = db.QueryRow(`SELECT
		event_id,
		client_id,
		thread_id,
		channel_id,
		creator,
		data,
		parent_thread_id,
		parent_event_id,
		branch_thread_id
	FROM events 
	WHERE client_id = $1 AND event_id = $2`,
		clientID,
		newEventID,
	).Scan(
		&gotEventID,
		&gotClientID,
		&gotThreadID,
		&gotChannelID,
		&gotCreatorID,
		&gotData,
		&gotParentThreadID,
		&gotParetnEventID,
		&gotBranchThreadID,
	)
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, gotEventID, newEventID)
	assert.Equal(t, gotClientID, clientID)
	assert.Equal(t, gotThreadID, branchThreadID)
	assert.Equal(t, gotChannelID, channelID)
	assert.Equal(t, gotCreatorID, creatorID)
	assert.Equal(t, gotData, []byte("hello machine"))
	assert.Equal(t, gotParentThreadID, threadID)
	assert.Equal(t, gotParetnEventID, nodalEventID)
	assert.Equal(t, gotBranchThreadID, uuid.Nil)
}

func TestHeyService_simple_NewBranchEvent(t *testing.T) {
	// db.Exec(`DELETE FROM events;`)

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

	// Create nodal event
	branchThreadID, nodalEventID, err := hey.CreateNodalEvent(
		ctx,
		threadID,
		owners,
		creatorID,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, branchThreadID, uuid.Nil)
	assert.NotEqual(t, nodalEventID, uuid.Nil)

	// Create new event
	newEventID, err := hey.CreateEvent(
		ctx,
		branchThreadID,
		creatorID,
		[]byte("hello machine 1"),
	)
	assert.NoError(t, err)
	assert.NotEqual(t, newEventID, uuid.Nil)

	newEventID, err = hey.CreateEvent(
		ctx,
		branchThreadID,
		creatorID,
		[]byte("hello machine 2"),
	)
	assert.NoError(t, err)
	assert.NotEqual(t, newEventID, uuid.Nil)

	newEventID, err = hey.CreateEvent(
		ctx,
		branchThreadID,
		creatorID,
		[]byte("hello machine 3"),
	)
	assert.NoError(t, err)
	assert.NotEqual(t, newEventID, uuid.Nil)

	//
	newBranchThreadID, newBranchEventid, err := hey.CreateNewBranchEvent(
		ctx,
		branchThreadID,
		newEventID,
		owners,
		creatorID,
		[]byte("first event in a new thread"),
	)
	assert.NoError(t, err)
	assert.NotEqual(t, newBranchThreadID, uuid.Nil)
	assert.NotEqual(t, newBranchEventid, uuid.Nil)

	// check in the database

	var gotBranchThreadID uuid.UUID

	err = db.QueryRow(`SELECT
		branch_thread_id
	FROM events
	WHERE client_id = $1 AND event_id = $2`,
		clientID,
		newEventID,
	).Scan(
		&gotBranchThreadID,
	)
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, gotBranchThreadID, newBranchThreadID)
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

	// Create nodal event
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

func TestHeyService_simple_SearchEvents(t *testing.T) {
	// db.Exec(`DELETE FROM events;`)

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

	// Create nodal event
	branchThreadID, nodalEventID, err := hey.CreateNodalEvent(
		ctx,
		threadID,
		owners,
		creatorID,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, branchThreadID, uuid.Nil)
	assert.NotEqual(t, nodalEventID, uuid.Nil)

	// ticker := time.Tick(time.Millisecond * 10)
	totalEvents := 10
	eventIDS := make([]uuid.UUID, totalEvents)

	for index := 0; index <= totalEvents-1; index++ {
		// select {
		// case <-ticker:
		// }

		// new event
		newEventID, err := hey.CreateEvent(
			ctx,
			branchThreadID,
			creatorID,
			[]byte("hello machine #"+strconv.Itoa(index)),
		)
		assert.NoError(t, err)
		assert.NotEqual(t, newEventID, uuid.Nil)
		eventIDS[index] = newEventID
	}

	// reverse ids
	for i, j := 0, len(eventIDS)-1; i < j; i, j = i+1, j-1 {
		eventIDS[i], eventIDS[j] = eventIDS[j], eventIDS[i]
	}

	// Search
	var perPage = 3
	var events = []interfaces.Event{}
	var cursor = ""
	var countItemsLastPage = totalEvents - (perPage * (totalEvents / perPage))

	for page := 0; page <= (totalEvents / perPage); page++ {
		events, cursor, err = hey.FindEvents(
			ctx,
			branchThreadID,
			cursor,
			perPage,
		)

		assert.NoError(t, err)
		assert.NotEmpty(t, cursor)

		var _perPage = perPage

		if page+1 > totalEvents/perPage {
			// for last page
			_perPage = countItemsLastPage
		}

		// check items per page
		assert.Equal(t, len(events), _perPage)

		for i := 0; i < _perPage; i++ {
			// t.Log(i, len(events), page*perPage+i)
			item := events[i]
			expectedItemID := eventIDS[page*perPage+i]
			assert.Equal(
				t,
				item.EventID(),
				eventIDS[page*perPage+i],
				expectedItemID,
			)
		}
	}

}
