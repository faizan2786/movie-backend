package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/faizan2786/movie-backend/movieapi/internal/gateway"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -source=controller.go -destination=mock_controller/mock.go

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

	// create a new error group
	errgrp, gctx := errgroup.WithContext(ctx)

	// fetch the metadata (in parallel)
	errgrp.Go(func() error {
		metadata, err := c.mgw.GetMetadata(gctx, id)
		switch {
		case errors.Is(err, gateway.ErrNotFound):
			return ErrNotFound
		case err != nil:
			return fmt.Errorf("GetMetadata: %w", err)
		default:
			movie.Metadata = metadata
		}
		return nil // return no error on a successful fetch
	})

	// fetch the rating (in parallel)
	errgrp.Go(func() error {
		rating, err := c.rgw.GetMovieRating(gctx, id)
		switch {
		case errors.Is(err, gateway.ErrNotFound):
			// a movie may exist without any ratings
		case err != nil:
			return fmt.Errorf("GetMovieRating: %w", err)
		default:
			movie.Rating = rating
		}
		return nil // return no error on a successful fetch (or on an empty rating)
	})

	// wait until all fetches finish (or an error is returned)
	if err := errgrp.Wait(); err != nil {
		return movie, err
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
