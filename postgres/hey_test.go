package postgres

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	_hey "github.com/gebv/hey"

	"github.com/gebv/hey/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

// drop table accounts, balance_changes, channel_counters, channel_watchers, channels, events, orders, thread_counters, thread_watchers, threadline, threads, transactions, users;

func ContextWithClientID() (context.Context, uuid.UUID) {
	clientID := uuid.NewV4()
	return context.WithValue(
			context.Background(),
			clientIDContextKey,
			clientID),
		clientID
}

func TestHey_createChannel(t *testing.T) {

	// ------------------------------------

	hey := NewService(
		db,
		log.New(os.Stderr, "[test hey]", 1),
	)

	ctx, clientID := ContextWithClientID()
	user1 := uuid.NewV4()
	user2 := uuid.NewV4()
	owners := []uuid.UUID{
		user1,
		user2,
	}
	channelID, threadID, err := hey.CreateChannel(
		ctx,
		owners,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, channelID, uuid.Nil)
	assert.NotEqual(t, threadID, uuid.Nil)

	// ------------------------------------

	// check data

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

func TestHey_newEventInThread(t *testing.T) {
	hey := NewService(
		db,
		log.New(os.Stderr, "[test hey]", 1),
	)

	ctx, clientID := ContextWithClientID()
	user1 := uuid.NewV4()
	user2 := uuid.NewV4()
	owners := []uuid.UUID{
		user1,
		user2,
	}
	channelID, threadID, err := hey.CreateChannel(
		ctx,
		owners,
	)

	// ------------------------

	creatorID := uuid.NewV4()

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

	// Create new event in existing thread
	message := []byte("hello from machine")

	newEventID, err := hey.CreateEvent(
		ctx,
		branchThreadID,
		creatorID,
		message,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, newEventID, uuid.Nil)

	// ------------------------------------

	// check data
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
	assert.Equal(t, gotData, message)
	assert.Equal(t, gotParentThreadID, threadID)
	assert.Equal(t, gotParetnEventID, nodalEventID)
	assert.Equal(t, gotBranchThreadID, uuid.Nil)
}

func TestHey_newBranchEvent(t *testing.T) {
	hey := NewService(
		db,
		log.New(os.Stderr, "[test hey]", 1),
	)

	ctx, clientID := ContextWithClientID()
	user1 := uuid.NewV4()
	user2 := uuid.NewV4()
	owners := []uuid.UUID{
		user1,
		user2,
	}
	channelID, threadID, err := hey.CreateChannel(
		ctx,
		owners,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, channelID, uuid.Nil)
	assert.NotEqual(t, threadID, uuid.Nil)

	// Create nodal event

	creatorID := uuid.NewV4()

	branchThreadID, nodalEventID, err := hey.CreateNodalEvent(
		ctx,
		threadID,
		owners,
		creatorID,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, branchThreadID, uuid.Nil)
	assert.NotEqual(t, nodalEventID, uuid.Nil)

	// Create new event in existing thread
	messages := [][]byte{
		[]byte("hello from machine #1"),
		[]byte("hello from machine #2"),
		[]byte("hello from machine #3"),
	}
	var messageIDs = make([]uuid.UUID, len(messages))

	for index, message := range messages {
		newEventID, err := hey.CreateEvent(
			ctx,
			branchThreadID,
			creatorID,
			message,
		)
		assert.NoError(t, err)
		assert.NotEqual(t, newEventID, uuid.Nil)

		messageIDs[index] = newEventID
	}

	// ------------------------------------

	eventID := messageIDs[1]

	newBranchThreadID, newBranchEventid, err := hey.CreateNewBranchEvent(
		ctx,
		branchThreadID,
		eventID, // message #2
		owners,
		creatorID,
		[]byte("first event in a new thread"),
	)
	assert.NoError(t, err)
	assert.NotEqual(t, newBranchThreadID, uuid.Nil)
	assert.NotEqual(t, newBranchEventid, uuid.Nil)

	// ------------------------------------

	// check data

	var gotBranchThreadID uuid.UUID

	err = db.QueryRow(`SELECT
		branch_thread_id
	FROM events
	WHERE client_id = $1 AND event_id = $2`,
		clientID,
		eventID,
	).Scan(
		&gotBranchThreadID,
	)
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, gotBranchThreadID, newBranchThreadID)
}

func TestHey_simple_nodalEvent(t *testing.T) {
	hey := NewService(
		db,
		log.New(os.Stderr, "[test hey]", 1),
	)

	ctx, clientID := ContextWithClientID()
	user1 := uuid.NewV4()
	user2 := uuid.NewV4()
	owners := []uuid.UUID{
		user1,
		user2,
	}
	channelID, threadID, err := hey.CreateChannel(
		ctx,
		owners,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, channelID, uuid.Nil)
	assert.NotEqual(t, threadID, uuid.Nil)

	// ------------------------------------

	creatorID := uuid.NewV4()

	branchThreadID, nodalEventID, err := hey.CreateNodalEvent(
		ctx,
		threadID,
		owners,
		creatorID,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, branchThreadID, uuid.Nil)
	assert.NotEqual(t, nodalEventID, uuid.Nil)

	// ------------------------------------

	// check data

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

func TestHey_simple_search(t *testing.T) {
	hey := NewService(
		db,
		log.New(os.Stderr, "[test hey]", 1),
	)

	ctx, _ := ContextWithClientID()
	user1 := uuid.NewV4()
	user2 := uuid.NewV4()
	owners := []uuid.UUID{
		user1,
		user2,
	}
	channelID, threadID, err := hey.CreateChannel(
		ctx,
		owners,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, channelID, uuid.Nil)
	assert.NotEqual(t, threadID, uuid.Nil)

	creatorID := uuid.NewV4()

	branchThreadID, nodalEventID, err := hey.CreateNodalEvent(
		ctx,
		threadID,
		owners,
		creatorID,
	)
	assert.NoError(t, err)
	assert.NotEqual(t, branchThreadID, uuid.Nil)
	assert.NotEqual(t, nodalEventID, uuid.Nil)

	// create events
	totalEvents := 10
	eventIDS := make([]uuid.UUID, totalEvents)

	for index := 0; index <= totalEvents-1; index++ {

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

	// ------------------------------------

	// Search
	var watcherID = uuid.NewV4()
	var perPage = 3
	var searchResult _hey.SearchResult
	var cursor = ""
	var countItemsLastPage = totalEvents - (perPage * (totalEvents / perPage))

	for page := 0; page <= (totalEvents / perPage); page++ {
		searchResult, err = hey.FindEvents(
			ctx,
			watcherID,
			branchThreadID,
			cursor,
			perPage,
		)
		cursor = searchResult.Cursor() // save for next query

		assert.NoError(t, err)
		// assert.NotEmpty(t, cursor)

		var _perPage = perPage

		if page+1 > totalEvents/perPage {
			// for last page
			_perPage = countItemsLastPage
		}

		// check items per page
		assert.Equal(t, len(searchResult.Events()), _perPage)

		for i := 0; i < _perPage; i++ {
			item := searchResult.Events()[i]

			t.Log(i, len(searchResult.Events()), page*perPage+i)

			expectedItemID := eventIDS[page*perPage+i]
			assert.Equal(
				t,
				uuid.Equal(
					eventIDS[page*perPage+i],
					item.EventID(),
				),
				true,
				"expected %v, got %v",
				expectedItemID,
				item.EventID(),
			)
		}
	}
}
