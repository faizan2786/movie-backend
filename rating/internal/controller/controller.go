package controller

import (
	"context"

	"github.com/faizan2786/movie-backend/rating/pkg/model"
)

// Repository provides a database access layer for the metadata service.
type Repository interface {
	GetRating(context.Context, model.RecordType, model.RecordID) (model.RatingSum, model.RatingCount, error)
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
