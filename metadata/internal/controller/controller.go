package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/faizan2786/movie-backend/metadata/internal/repository"
	"github.com/faizan2786/movie-backend/metadata/pkg/model"
)

// Repository provides a database access layer for the metadata service.
type Repository interface {
	GetMetadata(ctx context.Context, id string) (model.Metadata, error)
	UpsertMetadata(ctx context.Context, data model.Metadata, check func(existing model.Metadata) error) (created bool, err error)
}

// Controller implements business logic to retrieve movie metadata.
type Controller struct {
	repo Repository
}

// New returns a new [Controller].
func New(r Repository) *Controller {
	return &Controller{repo: r}
}

// GetMetadata returns [model.Metadata] for a movie with given id.
// It returns [ErrNotFound] if the movie doesn't exists.
func (c *Controller) GetMetadata(ctx context.Context, id string) (model.Metadata, error) {
	data, err := c.repo.GetMetadata(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return data, ErrNotFound
	}
	if err != nil {
		return data, fmt.Errorf("unknown repository error: %w", err)
	}
	return data, nil
}

// PutMetadata creates or replaces metadata for a movie with given id,
// and reports whether a new record was created.
//
// It returns [ErrNullId] if the id is empty, [ErrNullField] if the title or release year is empty,
// and [ErrImmutableField] if the title or release year differs from the stored metadata.
func (c *Controller) PutMetadata(ctx context.Context, data model.Metadata) (created bool, err error) {
	if len(data.ID) == 0 {
		return false, ErrNullId
	}
	if data.Title == "" {
		return false, fmt.Errorf("%w: title", ErrNullField)
	}
	if data.ReleaseYear == 0 {
		return false, fmt.Errorf("%w: releaseYear", ErrNullField)
	}

	// business logic decides if an existing record can be replaced; passed to the repo
	// so the check runs atomically with the db write
	checkUpdatable := func(existing model.Metadata) error {
		if existing.Title != data.Title || existing.ReleaseYear != data.ReleaseYear {
			return ErrImmutableField
		}
		return nil
	}

	created, err = c.repo.UpsertMetadata(ctx, data, checkUpdatable)
	if errors.Is(err, ErrImmutableField) {
		return false, err
	}
	if err != nil {
		return false, fmt.Errorf("unknown repository error: %w", err)
	}
	return created, nil
}
