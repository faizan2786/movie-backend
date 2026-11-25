package gateway

import (
	"errors"
	"time"
)

// ClientTimeout is the max. duration for making network calls to upstream services.
const ClientTimeout time.Duration = 5 * time.Second

var (
	// ErrNotFound is returned when a resource is "not found".
	ErrNotFound = errors.New("not found")

	// ErrInvalidInput is returned when a service rejects a request as "not valid".
	ErrInvalidInput = errors.New("input not valid")

	// ErrInvalidURL is returned when a service URL is "not valid".
	ErrInvalidURL = errors.New("URL not valid")
)
