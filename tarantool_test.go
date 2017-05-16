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

	thread := &Thread{}
	err = chrono.NewThread(thread)
	assert.NoError(t, err)

	user1 := &User{}
	err = chrono.NewUser(user1)
	assert.NoError(t, err)

	err = chrono.Observe(user1.UserID, thread.ThreadID)
	assert.NoError(t, err)

	event1 := &Event{
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event2 := &Event{
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event3 := &Event{
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event4 := &Event{
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event5 := &Event{
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
	assert.Equal(t, event4.EventID, events[0].EventID)
	assert.Equal(t, event3.EventID, events[1].EventID)
	assert.Equal(t, event2.EventID, events[2].EventID)
}

func TestRecentActivityThreadline(t *testing.T) {
	chrono, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, chrono)

	thread := &Thread{
		ThreadlineEnabled: true,
	}
	err = chrono.NewThread(thread)
	assert.NoError(t, err)

	thread2 := &Thread{}
	err = chrono.NewThread(thread2)
	assert.NoError(t, err)
	assert.NotEmpty(t, thread2.ThreadID)

	user1 := &User{}
	err = chrono.NewUser(user1)
	assert.NoError(t, err)

	user2 := &User{}
	err = chrono.NewUser(user2)
	assert.NoError(t, err)

	user3 := &User{}
	err = chrono.NewUser(user3)
	assert.NoError(t, err)

	err = chrono.Observe(user1.UserID, thread.ThreadID)
	assert.NoError(t, err)
	err = chrono.Observe(user3.UserID, thread.ThreadID)
	assert.NoError(t, err)

	observers, err := chrono.Observers(thread.ThreadID, 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(observers))

	err = chrono.Ignore(user1.UserID, thread.ThreadID)
	assert.NoError(t, err)

	observers, err = chrono.Observers(thread.ThreadID, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(observers))
	assert.Equal(t, user3.UserID, observers[0].UserID)

	err = chrono.Observe(user3.UserID, thread2.ThreadID)
	assert.NoError(t, err)

	observes, err := chrono.Observes(user3.UserID, 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(observes))

	// невозможно проверить
	now := time.Now()
	err = chrono.MarkAsDelivered(user3.UserID, thread.ThreadID, now)
	assert.NoError(t, err)

	event1 := &Event{
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event2 := &Event{
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event3 := &Event{
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}
	event4 := &Event{
		ThreadID:  thread.ThreadID,
		CreatedAt: time.Now(),
	}

	for _, ev := range []*Event{event1, event2, event3, event4} {
		err = chrono.NewEvent(ev)
		assert.NoError(t, err)
	}

	events, err := chrono.RecentActivity(user3.UserID, thread.ThreadID, 3, 3)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, event4.EventID, events[0].EventID)
	assert.Equal(t, event3.EventID, events[1].EventID)
	assert.Equal(t, event2.EventID, events[2].EventID)

}
