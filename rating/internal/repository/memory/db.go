package memory

import (
	"sync"
)

// DB is an in-memory database for the metadata service.
type DB struct {
	m sync.RWMutex
}

// NewDB returns a new [DB].
func NewDB() *DB {
	return &DB{}
}
