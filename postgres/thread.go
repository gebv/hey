package postgres

import (
	"context"

	"github.com/gebv/hey"
	uuid "github.com/satori/go.uuid"
)

func (s *Service) FindThreadByName(
	ctx context.Context,
	channelID uuid.UUID,
	name string,
) (hey.Thread, error) {
	// TODO: testings

	clientID := ClientIDFromContext(ctx)
	return s.threads.FindThreadByName(
		clientID,
		channelID,
		name,
	)
}
