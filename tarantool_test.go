package hey

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecentActivity(t *testing.T) {
	chrono, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, chrono)

	thread := &Thread{
		ThreadID: "ud3",
	}
	err = chrono.NewThread(thread)
	assert.NoError(t, err)

	user1 := &User{
		UserID: "user11",
	}
	err = chrono.NewUser(user1)
	assert.NoError(t, err)

	err = chrono.Observe(user1.UserID, thread.ThreadID)
	assert.NoError(t, err)

	event1 := &Event{
		EventID:   "ev100",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event2 := &Event{
		EventID:   "ev101",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event3 := &Event{
		EventID:   "ev102",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event4 := &Event{
		EventID:   "ev103",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event5 := &Event{
		EventID:   "ev104",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}

	for _, ev := range []*Event{event1, event2, event3, event4, event5} {
		err = chrono.NewEvent(ev)
		assert.NoError(t, err)
	}

	events, err := chrono.RecentActivity(user1.UserID, thread.ThreadID, 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, "ev103", events[0].EventID)
	assert.Equal(t, "ev102", events[1].EventID)
	assert.Equal(t, "ev101", events[2].EventID)
}

func TestRecentActivityByLastTs(t *testing.T) {
	chrono, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, chrono)

	thread := &Thread{
		ThreadID: "ud1",
	}
	err = chrono.NewThread(thread)
	assert.NoError(t, err)

	user1 := &User{
		UserID: "user1",
	}
	err = chrono.NewUser(user1)
	assert.NoError(t, err)

	err = chrono.Observe(user1.UserID, thread.ThreadID)
	assert.NoError(t, err)

	lastTs := time.Now()

	event1 := &Event{
		EventID:   "ev1",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event2 := &Event{
		EventID:   "ev2",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event3 := &Event{
		EventID:   "ev3",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event4 := &Event{
		EventID:   "ev4",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}

	for _, ev := range []*Event{event1, event2, event3, event4} {
		err = chrono.NewEvent(ev)
		assert.NoError(t, err)
	}

	events, err := chrono.RecentActivityByLastTS(user1.UserID, thread.ThreadID, 3, lastTs)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, "ev4", events[0].EventID)
	assert.Equal(t, "ev3", events[1].EventID)
	assert.Equal(t, "ev2", events[2].EventID)

}

func TestRecentActivityThreadline(t *testing.T) {
	chrono, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, chrono)

	thread := &Thread{
		ThreadID:          "ud2",
		ThreadlineEnabled: true,
	}
	err = chrono.NewThread(thread)
	assert.NoError(t, err)

	user1 := &User{
		UserID: "user2",
	}
	err = chrono.NewUser(user1)
	assert.NoError(t, err)

	err = chrono.Observe(user1.UserID, thread.ThreadID)
	assert.NoError(t, err)

	lastTs := time.Now()

	event1 := &Event{
		EventID:   "ev5",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event2 := &Event{
		EventID:   "ev6",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event3 := &Event{
		EventID:   "ev7",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event4 := &Event{
		EventID:   "ev8",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}

	for _, ev := range []*Event{event1, event2, event3, event4} {
		err = chrono.NewEvent(ev)
		assert.NoError(t, err)
	}

	events, err := chrono.RecentActivityByLastTS(user1.UserID, thread.ThreadID, 3, lastTs)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, "ev8", events[0].EventID)
	assert.Equal(t, "ev7", events[1].EventID)
	assert.Equal(t, "ev6", events[2].EventID)

}
