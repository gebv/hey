package hey

import "time"

// Chronograph represents storage methods
type Chronograph interface {
	// threads
	NewThread(threadID string) error

	NewThreadWithData(threadID string, dataType DataType, data interface{}) error

	// 1. Удаляем все записи из events
	DeleteThread(threadID) error

	// подписка конкретного юзера на трэд
	Observe(userID, threadID string) error

	ThreadObservers(threadID string) ([]User, error)

	// RecentActivityByLastTS возвращает события позже lastts
	RecentActivityByLastTS(threadID string, limit, lastts time.Time) ([]Event, error)
	// двигаться по limit,offset что предлагает tnt
	RecentActivity(threadID string, limit, offset int) ([]Event, error)

	// events
	// 1. Достаём всех подписчиков трэда.
	// 2. Вставляем в Timeline для всех подписчиков этот eventID
	NewEvent(threadID string, eventID string) error

	NewEventWithData(
		threadID string,
		eventID string,
		dataType DataType,
		data interface{},
	) error

	UpdateEvent(ev *Event) error
	DeleteEvent(eventID string) error
}
