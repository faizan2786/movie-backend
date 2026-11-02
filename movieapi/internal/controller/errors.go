package controller

import "errors"

var (
	// ErrNotFound is returned when a resource is "not found".
	ErrNotFound = errors.New("not found")

	// ErrInvalidID is returned when an ID is "not valid".
	ErrInvalidID = errors.New("id not valid")

	// ErrInvalidRating is returned when a rating is "not valid".
	ErrInvalidRating = errors.New("rating not valid")
)
