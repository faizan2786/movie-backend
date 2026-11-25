package grpcutil

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/faizan2786/movie-backend/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewConn returns a new grpc connection to one of the service instances.
func NewConn(ctx context.Context, service string, r discovery.ServiceRegistry) (*grpc.ClientConn, error) {
	// get metadata service addresses and choose one at random
	addrs, err := r.ServiceAddresses(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service addresses from registry: %w", err)
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no service address found for %q", service)
	}
	addr := addrs[rand.IntN(len(addrs))]

	// create a new grpc client connection
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection: %w", err)
	}

	return conn, nil
}
