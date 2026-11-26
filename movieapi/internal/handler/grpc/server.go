package grpc

import (
	"context"
	"errors"
	"log"

	pb "github.com/faizan2786/movie-backend/gen/movieapipb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/faizan2786/movie-backend/movieapi/internal/controller"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
)

// MovieController is a controller layer for movieapi service.
type MovieController interface {
	// GetMovie returns a movie's metadata and its aggregated rating.
	// A movie without any ratings is returned with a zero value rating.
	//
	// It returns [controller.ErrInvalidID] if the id is not valid and
	// [controller.ErrNotFound] if a movie with given id doesn't exist.
	GetMovie(ctx context.Context, id model.ID) (model.Movie, error)
}

// MovieServer is a grpc server for movieapi service.
type MovieServer struct {
	ctlr MovieController
	pb.UnimplementedMovieServiceServer
}

// NewServer returns a new [MovieServer].
func NewServer(mc MovieController) *MovieServer {
	return &MovieServer{ctlr: mc}
}

// GetMovie returns a movie's details and its aggregated rating.
//
// If the movie hasn't been rated yet, it returns the movie with the rating unset.
func (s *MovieServer) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	log.Printf("GetMovie: level=%s, movie_id=%s, msg=%s", "info", req.GetMovieId(), "request received")

	id := req.GetMovieId()
	if len(id) == 0 {
		log.Printf("GetMovie: level=%s, movie_id=%s, msg=%s", "debug", id, "movie id is null")
		return nil, status.Error(codes.InvalidArgument, "movie_id cannot be empty")
	}

	movie, err := s.ctlr.GetMovie(ctx, model.ID(id))
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrInvalidID):
			log.Printf("GetMovie: level=%s, movie_id=%s, error=%s", "debug", id, err)
			return nil, status.Error(codes.InvalidArgument, "movie_id may contain only letters, digits, and single hyphens between them")
		case errors.Is(err, controller.ErrNotFound):
			log.Printf("GetMovie: level=%s, movie_id=%s, error=%s", "debug", id, err)
			return nil, status.Error(codes.NotFound, "movie not found")
		default:
			log.Printf("GetMovie: level=%s, movie_id=%s, error=%s", "error", id, err)
			return nil, status.Error(codes.Internal, "something went wrong")
		}
	}

	return &pb.GetMovieResponse{Movie: toProto(movie)}, nil
}
