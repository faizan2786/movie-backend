package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/faizan2786/movie-backend/gen/metadatapb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/faizan2786/movie-backend/metadata/internal/controller"
	grpcserver "github.com/faizan2786/movie-backend/metadata/internal/handler/grpc"
	"github.com/faizan2786/movie-backend/pkg/server"

	"github.com/faizan2786/movie-backend/metadata/internal/repository/memory"
	"github.com/faizan2786/movie-backend/pkg/discovery"
	registry "github.com/faizan2786/movie-backend/pkg/discovery/consul"
)

const (
	defaultPort = 8081

	heartbeatInterval       = 1 * time.Second
	deregisterTimeout       = 3 * time.Second
	gracefulShutdownTimeout = 5 * time.Second
)

var srvName = discovery.MetadataServiceName

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	var port int
	flag.IntVar(&port, "port", defaultPort, "port to run the service at")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop) // stop receiving signal notifications as soon as possible (so a second Ctrl+C forces quit)

	// 1) create and register metadata grpc server
	srv := newGrpcServer()

	// 2) create a listener for the server process
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen on :%d: %w", port, err)
	}

	// 3) register the service with service registry.
	reg, err := registry.New(discovery.DefaultConsulAddr)
	if err != nil {
		return fmt.Errorf("create consul registry: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", "localhost", port)
	deregister, err := discovery.RegisterWithHeartbeat(ctx, reg, srvName, addr, heartbeatInterval)
	if err != nil {
		return fmt.Errorf("register service: %w", err)
	}

	// 4) run the grpc server with graceful shutdown
	deregisterFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), deregisterTimeout)
		defer cancel()
		deregister(ctx)
	}
	return server.ServeGrpc(ctx, srv, listener, gracefulShutdownTimeout, deregisterFunc)
}

// newGrpcServer initialise repo, controller, creates a new metadata grpc server and registers it.
func newGrpcServer() *grpc.Server {

	db := memory.NewDB()
	controller := controller.New(db)

	srv := grpc.NewServer()
	pb.RegisterMetadataServiceServer(srv, grpcserver.NewServer(controller))

	// Register reflection service on the server.
	reflection.Register(srv)

	return srv
}
