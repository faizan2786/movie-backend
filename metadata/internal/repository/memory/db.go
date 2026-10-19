package memory

import (
	"context"
	"sync"

	"github.com/faizan2786/movie-backend/metadata/internal/repository"
	"github.com/faizan2786/movie-backend/metadata/pkg/model"
)

// DB is an in-memory database for the metadata service.
type DB struct {
	m    sync.RWMutex
	data map[string]model.Metadata
}

// NewDB returns a new [DB].
func NewDB() *DB {
	return &DB{data: map[string]model.Metadata{}}
}

// GetMetadata returns [model.Metadata] for a movie with given id.
// It returns [repository.ErrNotFound] if the movie doesn't exists.
func (db *DB) GetMetadata(ctx context.Context, id string) (model.Metadata, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	var data model.Metadata
	var ok bool
	if data, ok = db.data[id]; !ok {
		return data, repository.ErrNotFound
	}
	return data, nil
}

// StoreMetadata stores metadata for a movie with given id.
// It returns [repository.ErrAlreadyExists] error when metadata already exists.
func (db *DB) StoreMetadata(ctx context.Context, data model.Metadata) error {
	if len(data.ID) == 0 {
		return repository.ErrNullId
	}

	db.m.Lock()
	defer db.m.Unlock()
	if _, ok := db.data[data.ID]; ok {
		return repository.ErrAlreadyExists
	}
	db.data[data.ID] = data
	return nil
}
