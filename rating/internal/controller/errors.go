package controller

import "errors"

var (
	// ErrNotFound is returned when a resource is "not found".
	ErrNotFound = errors.New("not found")

	// ErrEmptyID is returned when a provided resource ID is empty.
	ErrEmptyID = errors.New("record type, record id and user id cannot be empty")

	// ErrInvalidValue is returned when a provided rating value is out of range.
	ErrInvalidValue = errors.New("rating value is out of range")
)
