package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/faizan2786/movie-backend/rating/internal/repository"
	"github.com/faizan2786/movie-backend/rating/pkg/model"
)

// Repository provides a database access layer for the metadata service.
type Repository interface {
	GetRating(context.Context, model.RecordType, model.RecordID) (float64, error)
	UpsertRating(context.Context, model.Rating) error
}

// Controller implements business logic to retrieve movie metadata.
type Controller struct {
	repo Repository
}

// New returns a new [Controller].
func New(r Repository) *Controller {
	return &Controller{repo: r}
}

// GetRating returns average ratings for a given record type and record ID.
// It returns [ErrNotFound] when a rating record doesn't exist.
func (c *Controller) GetRating(ctx context.Context, recordType model.RecordType, recordID model.RecordID) (float64, error) {
	rating, err := c.repo.GetRating(ctx, recordType, recordID)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return rating, ErrNotFound
	}
	if err != nil {
		return rating, fmt.Errorf("unknown repository error: %w", err)
	}
	return rating, nil
}

// PutRating stores user rating by its record type and ID. It replaces the existing rating if there is any.
//
// It returns [ErrEmptyID] if the record type, record ID or user ID is empty and
// returns [ErrInvalidValue] if the rating value is out of range.
func (c *Controller) PutRating(ctx context.Context, r model.Rating) error {
	// validate inputs
	if len(r.RecordType) == 0 || len(r.RecordID) == 0 || len(r.UserID) == 0 {
		return ErrEmptyID
	}

	if ok := r.Value.Validate(); !ok {
		return ErrInvalidValue
	}

	err := c.repo.UpsertRating(ctx, r)
	if err != nil {
		return fmt.Errorf("unknown repository error: %w", err)
	}
	return nil
}
