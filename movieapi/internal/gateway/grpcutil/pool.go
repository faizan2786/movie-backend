package grpcutil

import (
	"fmt"
	"log"
	"sync"

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
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// try getting an existing connection
	if conn, ok := cp.pool[addr]; ok {
		return conn, nil
	}

	// add a new connection if it doesn't exist
	conn, err := newConnection(addr)
	if err != nil {
		return nil, err
	}

	// add connection to the pool
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
