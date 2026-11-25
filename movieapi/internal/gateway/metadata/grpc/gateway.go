package grpc

import (
	"context"
	"errors"

	pb "github.com/faizan2786/movie-backend/gen/metadatapb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/faizan2786/movie-backend/pkg/discovery"

	"github.com/faizan2786/movie-backend/movieapi/internal/gateway"
	"github.com/faizan2786/movie-backend/movieapi/internal/gateway/grpcutil"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
)

// Gateway is a gateway to communicate with the metadata service.
type Gateway struct {
	registry discovery.ServiceRegistry
	pool     *grpcutil.ConnPool
}

// NewGateway returns a new [Gateway] to the metadata service discovered via given registry.
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

// GetMetadata returns [model.Metadata] for a movie with given id.
// It returns [gateway.ErrNotFound] if the movie doesn't exist.
func (g *Gateway) GetMetadata(ctx context.Context, id model.ID) (model.Metadata, error) {
	conn, err := grpcutil.GetServiceConn(ctx, discovery.MetadataServiceName, g.registry, g.pool)
	if err != nil {
		return model.Metadata{}, err
	}

	// create a client for the Metadata RPC server
	client := pb.NewMetadataServiceClient(conn)

	ctxTimeout, cancel := context.WithTimeout(ctx, gateway.ClientTimeout)
	defer cancel()

	res, err := client.GetMetadata(ctxTimeout, &pb.GetMetadataRequest{MovieId: string(id)})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return model.Metadata{}, gateway.ErrNotFound
		case codes.InvalidArgument:
			return model.Metadata{}, gateway.ErrInvalidInput
		default:
			return model.Metadata{}, err
		}
	}

	metadataPB := res.GetMetadata()
	if metadataPB == nil {
		return model.Metadata{}, errors.New("metadata server returned empty response")
	}

	// convert metadata proto to the domain model
	return fromProto(metadataPB), nil
}
