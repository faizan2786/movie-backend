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
func (db *DB) GetMetadata(_ context.Context, id string) (model.Metadata, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	var data model.Metadata
	var ok bool
	if data, ok = db.data[id]; !ok {
		return data, repository.ErrNotFound
	}
	return data, nil
}

// UpsertMetadata inserts metadata for a movie, replacing any existing record with the same ID,
// and reports whether a new record was created.
//
// If a record already exists, check callback is called before replacing the record;
// a non-nil check error aborts the update.
//
// It returns [repository.ErrNullId] if the ID is empty.
func (db *DB) UpsertMetadata(_ context.Context, data model.Metadata, check func(existing model.Metadata) error) (created bool, err error) {
	if len(data.ID) == 0 {
		return false, repository.ErrNullId
	}

	db.m.Lock()
	defer db.m.Unlock()

	existing, ok := db.data[data.ID]
	// if data exists, call check (if it is not nil)
	if ok && check != nil {
		if err := check(existing); err != nil {
			return false, err
		}
	}

	// create/update the existing record
	db.data[data.ID] = data
	return !ok, nil // created is false when data was already existing
}
