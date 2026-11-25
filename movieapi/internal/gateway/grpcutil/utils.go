package grpcutil

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"

	"github.com/faizan2786/movie-backend/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnPool is a grpc connection pool containing a connection per service instance address.
//
// The zero value is ready to use. A ConnPool is safe for concurrent use.
type ConnPool struct {
	mu   sync.Mutex
	pool map[string]*grpc.ClientConn
}

// GetConnection returns the pooled grpc connection for the given service address,
// creating and caching a new one if none exists yet.
func (cp *ConnPool) GetConnection(addr string) (*grpc.ClientConn, error) {
	// try getting an existing connection
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if conn, ok := cp.pool[addr]; ok {
		return conn, nil
	}

	// add a new connection if it doesn't exist
	conn, err := newConnection(addr)
	if err != nil {
		return nil, err
	}

	if cp.pool == nil {
		cp.pool = map[string]*grpc.ClientConn{}
	}
	cp.pool[addr] = conn
	return conn, nil
}

// newConnection creates a new grpc connection to the given service address.
func newConnection(addr string) (*grpc.ClientConn, error) {
	// create a new grpc client connection
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection for a service instance at %q: %w", addr, err)
	}

	return conn, nil
}

// Close closes all grpc connections in the pool.
func (cp *ConnPool) Close() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for _, conn := range cp.pool {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close grpc connection: %s", err)
		}
	}
}

// GetServiceConn returns a grpc connection to an instance of the specified service at random.
func GetServiceConn(ctx context.Context, service string, r discovery.ServiceRegistry, pool *ConnPool) (*grpc.ClientConn, error) {
	// Get the service addresses
	addrs, err := r.ServiceAddresses(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service addresses from registry: %w", err)
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no addresses found for %q", discovery.MetadataServiceName)
	}

	// choose one at random and get its grpc connection from the pool
	addr := addrs[rand.IntN(len(addrs))]
	return pool.GetConnection(addr)
}
