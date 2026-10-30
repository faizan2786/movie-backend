package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/faizan2786/movie-backend/apiserver/internal/gateway"
	"github.com/faizan2786/movie-backend/apiserver/pkg/model"
)

// MetadataGateway is an interface for movie metadata.
type MetadataGateway interface {
	GetMetadata(ctx context.Context, id model.ID) (model.Metadata, error)
}

// RatingGateway is an interface for movie ratings.
type RatingGateway interface {
	GetMovieRating(ctx context.Context, movieID model.ID) (model.RatingValue, error)
	SubmitMovieRating(ctx context.Context, movieID model.ID, rating model.Rating) error
}

// Controller implements business logic for the server
type Controller struct {
	mgw MetadataGateway
	rgw RatingGateway
}

// New returns a new [Controller].
func New(mgw MetadataGateway, rgw RatingGateway) *Controller {
	return &Controller{
		mgw: mgw,
		rgw: rgw,
	}
}

// GetMovie returns a movie's data (including its aggregated rating)
// It returns [ErrNotFound] if a movie with given ID doesn't exist.
// A movie without any ratings is returned with a zero rating.
func (c *Controller) GetMovie(ctx context.Context, id model.ID) (model.Movie, error) {
	var movie model.Movie
	if !id.Valid() {
		return movie, ErrInvalidID
	}

	// fetch the metadata
	metadata, err := c.mgw.GetMetadata(ctx, id)
	if errors.Is(err, gateway.ErrNotFound) {
		return movie, ErrNotFound
	} else if err != nil {
		return movie, fmt.Errorf("GetMetadata: %w", err)
	}
	movie.Metadata = metadata

	// fetch the rating
	rating, err := c.rgw.GetMovieRating(ctx, id)
	switch {
	case errors.Is(err, gateway.ErrNotFound):
		// a movie may exist without any ratings
	case err != nil:
		return movie, fmt.Errorf("GetMovieRating: %w", err)
	default:
		movie.Rating = rating
	}

	return movie, nil
}

// SubmitMovieRating submits ratings for a movie for given user.
// It returns [ErrNotFound] if a movie with given ID doesn't exist yet and
// [ErrInvalidRating] if the rating is rejected by the rating service.
func (c *Controller) SubmitMovieRating(ctx context.Context, movieID model.ID, rating model.Rating) error {
	if !(movieID.Valid() && rating.UserID.Valid()) {
		return ErrInvalidID
	}

	// fetch the movie metadata (to check if a movie exists)
	_, err := c.mgw.GetMetadata(ctx, movieID)
	if errors.Is(err, gateway.ErrNotFound) {
		return fmt.Errorf("%w: movie doesn't exist", ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("GetMetadata: %w", err)
	}

	err = c.rgw.SubmitMovieRating(ctx, movieID, rating)
	if errors.Is(err, gateway.ErrInvalidInput) {
		return ErrInvalidRating
	} else if err != nil {
		return fmt.Errorf("SubmitMovieRating: %w", err)
	}
	return nil
}
