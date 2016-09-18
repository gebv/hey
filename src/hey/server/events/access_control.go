package events

import "context"

// AllowedCreateEvent returns error if the user is not allowed to create event
func (e *EventService) AllowedCreateEvent(
	ctx context.Context,
	creatorID int64,
	threadID string,
) chan error {

	return e.ac.AllowedUserCreateEvent(
		ctx,
		creatorID,
		threadID,
	)
}

// AccessControl verification of the access to create event
type AccessControl interface {
	AllowedUserCreateEvent(
		ctx context.Context,
		creatorID int64,
		threadID string,
	) chan error
}
