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
	registry "github.com/faizan2786/movie-backend/pkg/discovery/consul"
	"github.com/faizan2786/movie-backend/pkg/server"

	"github.com/faizan2786/movie-backend/metadata/internal/repository/memory"
	"github.com/faizan2786/movie-backend/pkg/discovery"
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

	addr := fmt.Sprintf("%s:%d", "localhost", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop) // stop receiving signal notifications as soon as possible (so a second Ctrl+C forces quit)

	// 1) register the service with registry.
	reg, err := registry.New(discovery.DefaultConsulAddr)
	if err != nil {
		return fmt.Errorf("create consul registry: %w", err)
	}

	deregister, err := discovery.RegisterWithHeartbeat(ctx, reg, srvName, addr, heartbeatInterval)
	if err != nil {
		return fmt.Errorf("register service: %w", err)
	}

	// 2) register metadata grpc server
	srv := registerGrpcServer()

	// 3) create a listener for the server process
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 4) run the grpc server with graceful shutdown
	deregisterFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), deregisterTimeout)
		defer cancel()
		deregister(ctx)
	}
	return server.ServeGrpc(ctx, srv, listener, gracefulShutdownTimeout, deregisterFunc)
}

// create a new grpc server and register metadata grpc server implementation
func registerGrpcServer() *grpc.Server {

	db := memory.NewDB()
	controller := controller.New(db)

	srv := grpc.NewServer()
	pb.RegisterMetadataServiceServer(srv, grpcserver.NewServer(controller))

	// Register reflection service on the server.
	reflection.Register(srv)

	return srv
}
