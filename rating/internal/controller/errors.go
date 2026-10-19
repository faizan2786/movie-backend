package controller

import "errors"

// ErrNotFound a "resource not found" error.
var (
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists represents a "resource already exists" error.
	ErrAlreadyExists = errors.New("already exists")

	// ErrNullId represents a "resource id is null/empty" error.
	ErrNullId = errors.New("id can not be empty")
)
