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

	thread2 := &Thread{}
	err = chrono.NewThread(thread2)
	assert.NoError(t, err)
	assert.NotEmpty(t, thread2.ThreadID)

	user1 := &User{
		UserID: "user2",
	}
	err = chrono.NewUser(user1)
	assert.NoError(t, err)

	user2 := &User{
		UserID: "user3",
	}
	err = chrono.NewUser(user2)
	assert.NoError(t, err)

	user3 := &User{
		UserID: "user4",
	}
	err = chrono.NewUser(user3)
	assert.NoError(t, err)

	err = chrono.Observe(user1.UserID, thread.ThreadID)
	assert.NoError(t, err)
	err = chrono.Observe(user3.UserID, thread.ThreadID)
	assert.NoError(t, err)

	observers, err := chrono.Observers(thread.ThreadID, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(observers))
	assert.Equal(t, user1.UserID, observers[0].UserID)
	assert.Equal(t, user3.UserID, observers[1].UserID)

	err = chrono.Ignore(user1.UserID, thread.ThreadID)
	assert.NoError(t, err)

	observers, err = chrono.Observers(thread.ThreadID, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(observers))
	assert.Equal(t, user3.UserID, observers[0].UserID)

	err = chrono.Observe(user3.UserID, thread2.ThreadID)
	assert.NoError(t, err)

	observes, err := chrono.Observes(user3.UserID, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(observes))
	assert.Equal(t, thread.ThreadID, observes[0].ThreadID)
	assert.Equal(t, thread2.ThreadID, observes[1].ThreadID)

	// невозможно проверить
	now := time.Now()
	err = chrono.MarkAsDelivered(user3.UserID, thread.ThreadID, now)
	assert.NoError(t, err)

	lastTs := time.Now()

	event1 := &Event{
		EventID:   "ev55",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event2 := &Event{
		EventID:   "ev66",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event3 := &Event{
		EventID:   "ev77",
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event4 := &Event{
		EventID:   "ev88",
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
	assert.Equal(t, "ev88", events[0].EventID)
	assert.Equal(t, "ev77", events[1].EventID)
	assert.Equal(t, "ev66", events[2].EventID)

}
