package grpc

import (
	"context"

	pb "github.com/faizan2786/movie-backend/gen/ratingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/faizan2786/movie-backend/pkg/discovery"

	"github.com/faizan2786/movie-backend/movieapi/internal/gateway"
	"github.com/faizan2786/movie-backend/movieapi/internal/gateway/grpcutil"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
)

// Gateway is a gateway to communicate with the rating service.
type Gateway struct {
	registry discovery.ServiceRegistry
	pool     *grpcutil.ConnPool
}

// NewGateway returns a new [Gateway] to the rating service discovered via given registry.
func NewGateway(r discovery.ServiceRegistry) *Gateway {
	return &Gateway{
		registry: r,
		pool:     &grpcutil.ConnPool{},
	}
}

// Close closes the grpc connections to upstream service.
func (g *Gateway) Close() {
	g.pool.Close()
}

// GetMovieRating returns the aggregated rating for a movie with given id.
// It returns [gateway.ErrNotFound] if the rating doesn't exist yet.
func (g *Gateway) GetMovieRating(ctx context.Context, movieID model.ID) (float64, error) {
	// get a grpc connection to the rating service
	conn, err := grpcutil.GetServiceConn(ctx, discovery.RatingServiceName, g.registry, g.pool)
	if err != nil {
		return 0, err
	}

	// create a grpc client and a request for the service
	client := pb.NewRatingServiceClient(conn)
	req := &pb.GetAggregatedRatingRequest{
		RecordType: pb.RecordType_RECORD_TYPE_MOVIE,
		RecordId:   string(movieID),
	}

	// call the RPC method
	ctxTimeout, cancel := context.WithTimeout(ctx, gateway.ClientTimeout)
	defer cancel()

	res, err := client.GetAggregatedRating(ctxTimeout, req)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return 0, gateway.ErrNotFound
		case codes.InvalidArgument:
			return 0, gateway.ErrInvalidInput
		default:
			return 0, err
		}
	}

	// return response
	if res.Rating == nil {
		return 0, gateway.ErrNotFound
	}
	return *res.Rating, nil
}

// SubmitMovieRating submits movie ratings for users.
// It returns [gateway.ErrInvalidInput] if the rating service rejects the rating as not valid.
func (g *Gateway) SubmitMovieRating(ctx context.Context, movieID model.ID, rating model.Rating) error {
	// get a grpc connection to the rating service
	conn, err := grpcutil.GetServiceConn(ctx, discovery.RatingServiceName, g.registry, g.pool)
	if err != nil {
		return err
	}

	// create a grpc client and a request for the service
	client := pb.NewRatingServiceClient(conn)
	req := &pb.PutRatingRequest{
		RecordType:  pb.RecordType_RECORD_TYPE_MOVIE,
		RecordId:    string(movieID),
		UserId:      string(rating.UserID),
		RatingValue: uint32(rating.Value),
	}

	// call the RPC method
	ctxTimeout, cancel := context.WithTimeout(ctx, gateway.ClientTimeout)
	defer cancel()
	_, err = client.PutRating(ctxTimeout, req)

	// return response
	if status.Code(err) == codes.InvalidArgument {
		return gateway.ErrInvalidInput
	}
	return err
}
