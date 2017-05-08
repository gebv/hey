package hey

import "time"

type Chronograph interface {
	// двигаться по lastts
	RecentActivityByLastTS(threadID string, lastts time.Time) ([]EventObserver, error)

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

	// и так далее по каждому из
}
