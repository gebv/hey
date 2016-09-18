package events

import (
	"strings"

	"github.com/satori/go.uuid"
)

// NewUUID returns uuid version 4
func NewUUID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}
