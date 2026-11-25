package grpc

import (
	"context"

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
	// registry to discover the metadata service instances
	registry discovery.ServiceRegistry
}

// NewGateway returns a new [Gateway] to the metadata service discovered via given registry.
func NewGateway(r discovery.ServiceRegistry) *Gateway {
	return &Gateway{registry: r}
}

// GetMetadata returns [model.Metadata] for a movie with given id.
// It returns [gateway.ErrNotFound] if the movie doesn't exists.
func (g *Gateway) GetMetadata(ctx context.Context, id model.ID) (model.Metadata, error) {
	// create a new grpc connection to the metadata service
	conn, err := grpcutil.NewGrpcClient(ctx, g.registry, discovery.MetadataServiceName)
	if err != nil {
		return model.Metadata{}, err
	}
	defer conn.Close()

	// create a client for the Metadata RPC server
	client := pb.NewMetadataServiceClient(conn)

	res, err := client.GetMetadata(ctx, &pb.GetMetadataRequest{MovieId: string(id)})

	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return model.Metadata{}, gateway.ErrNotFound
		default:
			return model.Metadata{}, err
		}
	}

	// convert response proto to the domain model
	var data *model.Metadata
	if data = fromProto(res); data == nil {
		return model.Metadata{}, gateway.ErrNotFound
	}
	return *data, nil
}
