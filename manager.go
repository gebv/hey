package hey

import "time"

// Manager represents storage methods
type Manager interface {
	// threads
	NewThread(*Thread) error
	GetThread(threadID string) (*Thread, error)
	UpdateThread(*Thread) error
	DeleteThread(threadID string) error

	// thread sources
	//AddSource(dstThread, sourceThread string) error
	//GetSources(threadID string, offset, limit uint32) ([]Thread, error)
	//GetRefers(sourceThreadID string, offset, limit uint32) ([]Thread, error)
	//RemoveSource(dstThread, sourceThread string) error

	// user subscriptions
	Observe(userID, threadID string) error
	Ignore(userID, threadID string) error
	Observers(threadID string, offset, limit uint32) ([]User, error)
	Observes(userID string, offset, limit uint32) ([]Thread, error)
	MarkAsDelivered(userID string, threadID string, times ...time.Time) error

	// threadline
	RecentActivityByLastTS(userID, threadID string, limit uint32, lastts time.Time) ([]Event, error)
	//RecentActivity(userID, threadID string, limit uint32) ([]Event, error)
	RecentActivity(userID, threadID string, offset, limit uint32) ([]Event, error)
	//ThreadlineActivity(userID, threadID string, limit, offset uint32) (events []Event, err error)

	// events
	NewEvent(*Event) error
	GetEvent(id string) (*Event, error)
	GetEvents(ids ...string) ([]Event, error)
	UpdateEvent(ev *Event) error
	DeleteEvent(eventID string) error

	// related data
	SetRelatedData(*RelatedData) error
	GetRelatedDatas(userID string, events ...Event) ([]EventObserver, error)

	// users
	NewUser(*User) error
	GetUser(userID string) (*User, error)
}
