package repository

import "errors"

// ErrNotFound a "resource not found" error.
var ErrNotFound error = errors.New("not found")

// ErrAlreadyExists represents a "resource already exists" error.
var ErrAlreadyExists error = errors.New("already exists")

// ErrNullId represents a "resource id is null/empty" error.
var ErrNullId error = errors.New("id can not be empty")
