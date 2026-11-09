package main

import (
	"context"
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
)

const (
	port            = "8083"
	addr            = "localhost:" + port
	shutdownTimeout = 5 * time.Second
)

type deregisterFunc = discovery.DeregisterFunc

func main() {

	reg := registry.New()

	// register an instance of metadata and rating service
	err := reg.Register(context.Background(), discovery.MetadataServiceName, discovery.GenerateInstanceID(), "localhost:8081")
	if err != nil {
		panic(err)
	}

	log.Printf("register instance for %s", discovery.MovieAPIServiceName)
	err = reg.Register(context.Background(), discovery.RatingServiceName, discovery.GenerateInstanceID(), "localhost:8082")
	if err != nil {
		panic(err)
	}

	// get cancellable context for Ctrl+C / SIGTERM signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	deregister, err := discovery.RegisterWithHeartbeat(ctx, reg, discovery.MovieAPIServiceName, addr, time.Second)
	if err != nil {
		panic(err)
	}

	// Intialise the service handler
	metadataGw := metadata.NewGateway(reg)
	ratingGw := rating.NewGateway(reg)
	controller := controller.New(metadataGw, ratingGw)
	handler := handler.New(controller)

	// register API paths with handlers
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies/{id}", handler.GetMovie)
	mux.HandleFunc("PUT /movies/{id}/ratings/{userID}", handler.PutMovieRating)

	// start the server
	cancelledServer := serveUntilCancelled(ctx, mux, deregister)

	// restore default signal handling so a second Ctrl+C applies a force-quit
	stop()
	shutdown(cancelledServer, deregister)
}

// serveUntilCancelled creates a new http server and start serving requests
// until an explicit cancellation signal is received.
func serveUntilCancelled(ctx context.Context, mux http.Handler, deregister deregisterFunc) *http.Server {
	// start the web server
	log.Printf("starting server at port: %s\n", port)
	srv := &http.Server{Addr: addr, Handler: mux}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		log.Printf("cannot start server at port: %s\n", port)
		deregister(ctx)
		panic(err)
	}

	return srv
}

// shutdown gracefully shuts down the http server
func shutdown(srv *http.Server, deregister deregisterFunc) {
	// deregister first so no new traffic is routed to this service,
	// then stop accepting new requests
	// (drain any in-flight requests by waiting for shutdownTimeout)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	deregister(timeoutCtx)
	log.Println("shutting down movieapi service")
	defer cancel()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Printf("movieapi shutdown: %v", err)
	}
}
