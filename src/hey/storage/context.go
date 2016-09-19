package storage

import (
	"context"
	"log"
)

// FromContext returns a connection to the database
func FromContext(ctx context.Context) DB {
	if db, ok := ctx.Value("__conn").(DB); ok {
		return db
	}
	log.Fatalln("[FAIL]", "empty db connection")
	return nil
}
