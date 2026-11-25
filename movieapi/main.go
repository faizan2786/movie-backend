package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faizan2786/movie-backend/movieapi/internal/controller"
	metadatagw "github.com/faizan2786/movie-backend/movieapi/internal/gateway/metadata/grpc"
	ratinggw "github.com/faizan2786/movie-backend/movieapi/internal/gateway/rating"
	"github.com/faizan2786/movie-backend/movieapi/internal/handler"
	"github.com/faizan2786/movie-backend/pkg/discovery"
	registry "github.com/faizan2786/movie-backend/pkg/discovery/consul"
	"github.com/faizan2786/movie-backend/pkg/server"
)

const (
	defaultPort = 8083

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
	var port int
	flag.IntVar(&port, "port", defaultPort, "port to run the service at")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", "localhost", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop) // stop receiving signal notifications as soon as possible (so a second Ctrl+C forces quit)

	// 1) register the service with registry
	reg, err := registry.New(discovery.DefaultConsulAddr)
	if err != nil {
		return fmt.Errorf("create consul registry: %w", err)
	}

	deregister, err := discovery.RegisterWithHeartbeat(ctx, reg, srvName, addr, heartbeatInterval)
	if err != nil {
		return fmt.Errorf("register service: %w", err)
	}

	// 2) create upstream service gateways
	metadataGw := metadatagw.NewGateway(reg)
	defer metadataGw.Close()
	ratingGw := ratinggw.NewGateway(reg)
	// defer ratingGw.Close() // TODO

	// 3) register API handlers
	mux := registerHandlers(metadataGw, ratingGw)

	// 4) serve until a shutdown signal is received, then deregister & shut down gracefully
	srv := &http.Server{Addr: addr, Handler: mux}

	deregisterFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), deregisterTimeout)
		defer cancel()
		deregister(ctx)
	}
	return server.ServeHttp(ctx, srv, gracefulShutdownTimeout, deregisterFunc)
}

func registerHandlers(mgw controller.MetadataGateway, rgw controller.RatingGateway) http.Handler {
	controller := controller.New(mgw, rgw)
	handler := handler.New(controller)

	// register API paths with handlers
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies/{id}", handler.GetMovie)
	mux.HandleFunc("PUT /movies/{id}/ratings/{userID}", handler.PutMovieRating)

	return mux
}
