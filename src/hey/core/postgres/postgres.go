package postgres

import (
	"hey/core"
	"hey/storage/postgres"
)

var (
	_ core.EventRepository   = (*postgres.EventRepository)(nil)
	_ core.ThreadRepository  = (*postgres.ThreadRepository)(nil)
	_ core.ChannelRepository = (*postgres.ChannelRepository)(nil)
)
