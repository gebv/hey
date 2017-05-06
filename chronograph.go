package hey

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Chronograph interface {
	// двигаться по lastts
	RecentActivityByLastTS(threadID string, lastts time.Time) ([]EventObserver, error)

	// двигаться по limit,offset что предлагает tnt
	RecentActivity(threadID string, limit, offset int) ([]EventObserver, error)

	NewEvent(
		threadID,
		creatorID string,
		eventID uuid.UUID,
		dataType DataType,
		data interface{},
	) error

	// и так далее по каждому из
}
