package postgres

import (
	"context"
	"log"

	uuid "github.com/satori/go.uuid"
)

var (
	clientIDContextKey = "ClientID"
)

func getUUIDFromContext(key string, ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(key).(uuid.UUID); ok {
		return id
	}

	log.Panicln("[FAIL]", "the context does not contain information about the context key", key)

	return uuid.Nil
}

// ClientIDFromContext returns client ID from context
func ClientIDFromContext(ctx context.Context) uuid.UUID {
	return getUUIDFromContext(clientIDContextKey, ctx)
}
