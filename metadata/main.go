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

	"github.com/faizan2786/movie-backend/metadata/internal/controller"
	"github.com/faizan2786/movie-backend/metadata/internal/handler"
	"github.com/faizan2786/movie-backend/metadata/internal/repository/memory"
	"github.com/faizan2786/movie-backend/pkg/discovery"
	registry "github.com/faizan2786/movie-backend/pkg/discovery/memory"
	"github.com/faizan2786/movie-backend/pkg/server"
)

const (
	port = "8081"
	addr = "localhost:" + port

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop) // stop receiving signal notifications as soon as possible (so a second Ctrl+C forces quit)

	// 1) register the service with registry.
	reg := registry.New()
	deregister, err := discovery.RegisterWithHeartbeat(ctx, reg, srvName, addr, heartbeatInterval)
	if err != nil {
		return fmt.Errorf("register service: %w", err)
	}

	// 2) register API handlers
	mux := registerHandlers()

	// 3) serve until a shutdown signal is received, then deregister & shut down gracefully
	srv := &http.Server{Addr: addr, Handler: mux}

	deregisterFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), deregisterTimeout)
		defer cancel()
		deregister(ctx)
	}
	return server.Run(ctx, srv, gracefulShutdownTimeout, deregisterFunc)
}

func registerHandlers() http.Handler {

	db := memory.NewDB()
	controller := controller.New(db)
	handler := handler.New(controller)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{id}", handler.GetMetadata)
	mux.HandleFunc("PUT /{id}", handler.PutMetadata)

	return mux
}
