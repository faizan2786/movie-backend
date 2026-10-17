package memory

import (
	"context"

	"github.com/faizan2786/movie-backend/metadata/internal/repository"
	"github.com/faizan2786/movie-backend/metadata/pkg/model"
)

// DB is an in-memory database for the metadata service.
type DB struct {
	data map[string]model.Metadata
}

// NewDB returns a new [DB].
func NewDB() *DB {
	return &DB{data: map[string]model.Metadata{}}
}

// GetMetadata returns [model.Metadata] for a movie with given id.
// It returns [repository.ErrNotFound] if the movie doesn't exists.
func (db *DB) GetMetadata(ctx context.Context, id string) (model.Metadata, error) {
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
	if len(data.Id) == 0 {
		return repository.ErrNullId
	}
	if _, ok := db.data[data.Id]; ok {
		return repository.ErrAlreadyExists
	}
	db.data[data.Id] = data
	return nil
}
