package hey

import "time"

// Manager represents storage methods
type Manager interface {
	// threads
	NewThread(*Thread) error
	GetThread(threadID string) (*Thread, error)
	UpdateThread(*Thread) error
	// 1. Удаляем все записи из events
	DeleteThread(threadID string) error

	// subscriptions
	// подписка конкретного юзера на трэд
	Observe(userID, srcThreadID, desThreadID string) error
	// отписка от трэда
	Ignore(userID, srcThreadID string) error
	ThreadObservers(threadID string) ([]User, error)
	// список трэдов юзера
	Subscriptions(userID string) ([]Thread, error)

	// threadline
	// RecentActivityByLastTS возвращает события позже lastts
	RecentActivityByLastTS(threadID string, limit, lastts time.Time) ([]Event, error)
	// двигаться по limit,offset что предлагает tnt
	RecentActivity(threadID string, limit, offset int) ([]Event, error)

	// events
	// 1. Достаём всех подписчиков трэда.
	// 2. Вставляем в Timeline для всех подписчиков этот eventID
	NewEvent(*Event) error
	GetEvent(id string) error
	GetEvents(ids ...string) ([]Event, error)
	UpdateEvent(ev *Event) error
	DeleteEvent(eventID string) error
}
