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

	"github.com/faizan2786/movie-backend/pkg/discovery"
	registry "github.com/faizan2786/movie-backend/pkg/discovery/memory"
	"github.com/faizan2786/movie-backend/rating/internal/controller"
	"github.com/faizan2786/movie-backend/rating/internal/handler"
	"github.com/faizan2786/movie-backend/rating/internal/repository/memory"
)

const (
	port = "8082"
	addr = "localhost:" + port

	heartbeatInterval       = 1 * time.Second
	deregisterTimeout       = 3 * time.Second
	gracefulShutdownTimeout = 5 * time.Second
)

type deRegisterFunc = discovery.DeregisterFunc

var srvName = discovery.RatingServiceName

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1) register the service with registry.
	reg := registry.New()
	deregister, err := discovery.RegisterWithHeartbeat(ctx, reg, srvName, addr, heartbeatInterval)
	if err != nil {
		return fmt.Errorf("register service: %w", err)
	}

	// 2) register the API handlers
	mux := registerHandlers()

	// 3) start the server
	cancelledServer, err := startServerUntilCancelled(ctx, mux)
	if err != nil {
		deregister(ctx)
		return fmt.Errorf("cannot start server at port %s: %w", port, err)
	}
	stop() // stop receiving signal notifications as soon as possible (so a second Ctrl+C forces quit)

	// 4) a shutdown signal is received, perform cleanup & shutdown server
	return gracefulShutdown(cancelledServer, deregister)
}

func registerHandlers() http.Handler {

	db := memory.NewDB()
	controller := controller.New(db)
	handler := handler.New(controller)
	mux := http.NewServeMux()

	// register handlers
	mux.HandleFunc("GET /{type}/{id}", handler.GetRating)
	mux.HandleFunc("PUT /{type}/{id}", handler.PutRating)

	return mux
}

func startServerUntilCancelled(ctx context.Context, handler http.Handler) (*http.Server, error) {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		// start the web server
		log.Printf("starting server at port %s\n", port)
		errCh <- server.ListenAndServe()
	}()

	// listen for shutdown signal via context.Done channel
	select {
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return server, nil
	}
}

func gracefulShutdown(server *http.Server, deregister deRegisterFunc) error {

	// first, deregister with a timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), deregisterTimeout)
	deregister(timeoutCtx)
	cancel()

	// then, shutdown gracefully with a timeout
	timeoutCtx, cancel = context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	log.Printf("shutting down %s service", srvName)

	if err := server.Shutdown(timeoutCtx); err != nil {
		server.Close() // apply force shutdown if the graceful shutdown errors
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	return nil
}
