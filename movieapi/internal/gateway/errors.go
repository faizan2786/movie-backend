package gateway

import (
	"errors"
)

var (
	// ErrNotFound is returned when a resource is "not found".
	ErrNotFound = errors.New("not found")

	// ErrInvalidInput is returned when a service rejects a request as "not valid".
	ErrInvalidInput = errors.New("input not valid")

	// ErrInvalidURL is returned when a service URL is "not valid".
	ErrInvalidURL = errors.New("URL not valid")
)
