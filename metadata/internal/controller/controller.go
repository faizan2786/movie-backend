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
	StoreMetadata(ctx context.Context, data model.Metadata) error
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

// PutMetadata stores metadata for a movie with given id.
// It returns [ErrAlreadyExists] error when metadata already exists.
func (c *Controller) PutMetadata(ctx context.Context, data model.Metadata) error {
	if len(data.ID) == 0 {
		return ErrNullId
	}

	err := c.repo.StoreMetadata(ctx, data)
	if err != nil && errors.Is(err, repository.ErrAlreadyExists) {
		return ErrAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("unknown repository error: %w", err)
	}
	return nil
}
