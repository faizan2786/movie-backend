package grpcutil

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/faizan2786/movie-backend/pkg/discovery"
	"google.golang.org/grpc"
)

// GetServiceConn returns a grpc connection to a healthy instance of the specified service at random.
func GetServiceConn(ctx context.Context, service string, r discovery.ServiceRegistry, pool *ConnPool) (*grpc.ClientConn, error) {
	// Get the service addresses
	addrs, err := r.ServiceAddresses(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service addresses from registry: %w", err)
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no addresses found for %q", service)
	}

	// choose one at random and get its grpc connection from the pool
	addr := addrs[rand.IntN(len(addrs))]
	return pool.GetConnection(addr)
}
