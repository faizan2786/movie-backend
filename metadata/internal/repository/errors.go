package repository

import "errors"

// ErrNotFound is returned when a resource is "not found".
var (
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is returned when a resource "already exists".
	ErrAlreadyExists = errors.New("already exists")

	// ErrNullId is returned when a provided resource ID is empty.
	ErrNullId = errors.New("id can not be empty")
)
