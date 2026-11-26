package grpc

import (
	pb "github.com/faizan2786/movie-backend/gen/movieapipb"

	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
)

// toProto converts a [model.Movie] to a [pb.Movie].
// A zero rating (i.e. a movie not rated yet) leaves the rating unset.
func toProto(movie model.Movie) *pb.Movie {
	msg := &pb.Movie{
		MovieId:     string(movie.ID),
		Title:       movie.Title,
		ReleaseYear: movie.ReleaseYear,
		Description: movie.Description,
	}
	if movie.Rating != 0 {
		msg.Rating = new(movie.Rating)
	}
	return msg
}
