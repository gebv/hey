package models

import (
	"errors"
)

var (
	ErrNotSupported = errors.New("not_supported")
	ErrNotFound = errors.New("not_found")
	ErrNotValid = errors.New("not_valid")
	ErrUnknown = errors.New("unknown")
	ErrNotAllowed = errors.New("not_allowed")
)
