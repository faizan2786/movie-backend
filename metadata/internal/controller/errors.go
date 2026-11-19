package controller

import "errors"

var (
	// ErrNotFound is returned when a resource is "not found".
	ErrNotFound = errors.New("not found")

	// ErrImmutableField is returned when an update tries to change a field that can't be changed.
	ErrImmutableField = errors.New("field can not be changed")

	// ErrNullId is returned when a provided resource ID is empty.
	ErrNullId = errors.New("id can not be empty")

	// ErrNullField is returned when an expected field is empty.
	ErrNullField = errors.New("field can not be empty")
)
