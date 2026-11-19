package repository

import "errors"

// ErrNotFound is returned when a resource is "not found".
var (
	ErrNotFound = errors.New("not found")

	// ErrNullId is returned when a provided resource ID is empty.
	ErrNullId = errors.New("id can not be empty")
)
