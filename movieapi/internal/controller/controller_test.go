package controller_test

import (
	"testing"

	"github.com/faizan2786/movie-backend/movieapi/internal/controller"
	mock "github.com/faizan2786/movie-backend/movieapi/internal/controller/mock_controller"
	"github.com/faizan2786/movie-backend/movieapi/internal/gateway"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testMovieID     model.ID = "abc-123"
	testUserID               = "1"
	testRatingValue          = model.RatingValue(7.5)
)

var testMetadata = model.Metadata{
	ID:          testMovieID,
	Title:       "The Test Movie",
	ReleaseYear: 2026,
	Description: "A movie for testing",
}

var testMovie = model.Movie{
	Metadata: testMetadata,
	Rating:   testRatingValue,
}

var testRating = model.Rating{
	UserID: testUserID,
	Value:  testRatingValue,
}

func TestGetMovie(t *testing.T) {
	t.Parallel()

	// create a new controller instance with mocked gateways
	ctrl := gomock.NewController(t)
	metadataMock := mock.NewMockMetadataGateway(ctrl)
	ratingMock := mock.NewMockRatingGateway(ctrl)
	c := controller.New(metadataMock, ratingMock)

	// set the mock call expectations with fake values to be returned
	metadataMock.EXPECT().GetMetadata(gomock.Any(), testMovieID).Return(testMetadata, nil)
	ratingMock.EXPECT().GetMovieRating(gomock.Any(), testMovieID).Return(testRatingValue, nil)

	// call the method under test
	data, err := c.GetMovie(t.Context(), testMovieID)
	require.Nil(t, err)
	require.Equal(t, testMovie, data)
}

func TestGetMovie_MetadataErrNotFound(t *testing.T) {
	t.Parallel()

	// create a new controller instance with mocked gateways
	ctrl := gomock.NewController(t)
	metadataMock := mock.NewMockMetadataGateway(ctrl)
	ratingMock := mock.NewMockRatingGateway(ctrl)
	c := controller.New(metadataMock, ratingMock)

	// set the mock call expectations with Error values to be returned
	metadataMock.EXPECT().GetMetadata(gomock.Any(), testMovieID).Return(model.Metadata{}, gateway.ErrNotFound)
	ratingMock.EXPECT().GetMovieRating(gomock.Any(), testMovieID).Return(model.RatingValue(0), nil)

	// call the method under test
	_, err := c.GetMovie(t.Context(), testMovieID)
	require.ErrorIs(t, err, controller.ErrNotFound)
}

func TestGetMovie_RatingErrNotFound(t *testing.T) {
	t.Parallel()

	// create a new controller instance with mocked gateways
	ctrl := gomock.NewController(t)
	metadataMock := mock.NewMockMetadataGateway(ctrl)
	ratingMock := mock.NewMockRatingGateway(ctrl)
	c := controller.New(metadataMock, ratingMock)

	// set the mock call expectations with Error values to be returned
	metadataMock.EXPECT().GetMetadata(gomock.Any(), testMovieID).Return(model.Metadata{}, nil)
	ratingMock.EXPECT().GetMovieRating(gomock.Any(), testMovieID).Return(model.RatingValue(0), gateway.ErrNotFound)

	// call the method under test
	_, err := c.GetMovie(t.Context(), testMovieID)
	require.Nil(t, err) // controller doesn't return error when rating is not found
}

func TestSubmitMovieRating(t *testing.T) {
	t.Parallel()

	// create a new controller instance with mocked gateways
	ctrl := gomock.NewController(t)
	metadataMock := mock.NewMockMetadataGateway(ctrl)
	ratingMock := mock.NewMockRatingGateway(ctrl)
	c := controller.New(metadataMock, ratingMock)

	// set the mock call expectations with fake values to be returned
	metadataMock.EXPECT().GetMetadata(gomock.Any(), testMovieID).Return(model.Metadata{}, nil)
	ratingMock.EXPECT().SubmitMovieRating(gomock.Any(), testMovieID, testRating).Return(nil)

	// call the method under test
	err := c.SubmitMovieRating(t.Context(), testMovieID, testRating)
	require.Nil(t, err)
}

func TestSubmitMovieRating_MetadataErrNotFound(t *testing.T) {
	t.Parallel()

	// create a new controller instance with mocked gateways
	ctrl := gomock.NewController(t)
	metadataMock := mock.NewMockMetadataGateway(ctrl)
	ratingMock := mock.NewMockRatingGateway(ctrl)
	c := controller.New(metadataMock, ratingMock)

	// set the mock call expectations with Error values to be returned
	metadataMock.EXPECT().GetMetadata(gomock.Any(), testMovieID).Return(model.Metadata{}, gateway.ErrNotFound)
	// no calls to rating service should be made if the metadata is not found
	ratingMock.EXPECT().SubmitMovieRating(gomock.Any(), testMovieID, testRating).Return(nil).Times(0)

	// call the method under test
	err := c.SubmitMovieRating(t.Context(), testMovieID, testRating)
	require.ErrorIs(t, err, controller.ErrNotFound)
}

func TestSubmitMovieRating_RatingErrInvalidInput(t *testing.T) {
	t.Parallel()

	// create a new controller instance with mocked gateways
	ctrl := gomock.NewController(t)
	metadataMock := mock.NewMockMetadataGateway(ctrl)
	ratingMock := mock.NewMockRatingGateway(ctrl)
	c := controller.New(metadataMock, ratingMock)

	// set the mock call expectations with Error values to be returned
	metadataMock.EXPECT().GetMetadata(gomock.Any(), testMovieID).Return(model.Metadata{}, nil)
	ratingMock.EXPECT().SubmitMovieRating(gomock.Any(), testMovieID, testRating).Return(gateway.ErrInvalidInput)

	// call the method under test
	err := c.SubmitMovieRating(t.Context(), testMovieID, testRating)
	require.ErrorIs(t, err, controller.ErrInvalidRating)
}
