package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faizan2786/movie-backend/movieapi/internal/controller"
	"github.com/faizan2786/movie-backend/movieapi/internal/gateway/metadata"
	"github.com/faizan2786/movie-backend/movieapi/internal/gateway/rating"
	"github.com/faizan2786/movie-backend/movieapi/internal/handler"
	"github.com/faizan2786/movie-backend/pkg/discovery"
	registry "github.com/faizan2786/movie-backend/pkg/discovery/memory"
	"github.com/faizan2786/movie-backend/pkg/server"
)

const (
	port = "8083"
	addr = "localhost:" + port

	heartbeatInterval       = 1 * time.Second
	deregisterTimeout       = 3 * time.Second
	gracefulShutdownTimeout = 5 * time.Second
)

var srvName = discovery.MovieAPIServiceName

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop) // stop receiving signal notifications as soon as possible (so a second Ctrl+C forces quit)

	// 1) register the service with registry
	reg := registry.New()
	deregister, err := discovery.RegisterWithHeartbeat(ctx, reg, srvName, addr, heartbeatInterval)
	if err != nil {
		return fmt.Errorf("register service: %w", err)
	}

	// also register metadata and rating services as we are using an in-memory registry
	reg.Register(ctx, discovery.MetadataServiceName, discovery.GenerateInstanceID(), "localhost:8081")
	reg.Register(ctx, discovery.RatingServiceName, discovery.GenerateInstanceID(), "localhost:8082")

	// 2) register API handlers
	mux := registerHandlers(reg)

	// 3) serve until a shutdown signal is received, then deregister & shut down gracefully
	srv := &http.Server{Addr: addr, Handler: mux}

	deregisterFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), deregisterTimeout)
		defer cancel()
		deregister(ctx)
	}
	return server.Run(ctx, srv, gracefulShutdownTimeout, deregisterFunc)
}

func registerHandlers(reg discovery.ServiceRegistry) http.Handler {

	// Intialise the service handler
	metadataGw := metadata.NewGateway(reg)
	ratingGw := rating.NewGateway(reg)
	controller := controller.New(metadataGw, ratingGw)
	handler := handler.New(controller)

	// register API paths with handlers
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies/{id}", handler.GetMovie)
	mux.HandleFunc("PUT /movies/{id}/ratings/{userID}", handler.PutMovieRating)

	return mux
}
