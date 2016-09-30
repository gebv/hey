package postgres

import (
	"log"
	"time"

	"github.com/gebv/hey"

	pg "gopkg.in/jackc/pgx.v2"
)

var (
	_ hey.Hey = (*Service)(nil)

	// PerPageMax max count items of page search
	PerPageMax = 100

	// PerPageDefault default count items of page search
	PerPageDefault = 25

	// TimeoutDefault default timeout of simple procedure
	TimeoutDefault time.Duration = 50

	// FindTimeoutDefault default timeout of search procedure
	FindTimeoutDefault = TimeoutDefault * 5

	DefaultPort     uint16 = 5432
	DefaultMaxConns int    = 10
)

// NewService new hey service
func NewService(
	db *pg.ConnPool,
	logger *log.Logger,
) *Service {
	return &Service{
		db:       db,
		channels: &ChannelRepository{db},
		threads:  &ThreadRepository{db},
		events:   &EventRepository{db},
		logger:   logger,
	}
}

// Service service of threads and events
type Service struct {
	logger   *log.Logger
	db       *pg.ConnPool
	channels *ChannelRepository
	threads  *ThreadRepository
	events   *EventRepository
}
