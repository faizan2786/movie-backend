// Package server provides helpers for running an HTTP server with graceful shutdown.
package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

// ServeHttp runs the given http server until ctx is cancelled, then calls beforeShutdown
// (e.g. deregister) and drains in-flight requests within drainTimeout,
// force-closing if draining takes too long.
//
// beforeShutdown is also called if the server fails to start or stops serving
// unexpectedly, so callers can rely on it for cleanup.
// It can be nil.
func ServeHttp(ctx context.Context, srv *http.Server, drainTimeout time.Duration, beforeShutdown func()) error {

	// 1) start the server
	errCh := make(chan error, 1)
	go func() {
		log.Printf("starting server at %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	// 2) listen for server error or a shutdown signal (i.e. context cancel)
	select {
	case err := <-errCh:
		// 2-a) server failed to start, perform cleanup
		if beforeShutdown != nil {
			beforeShutdown()
		}
		return fmt.Errorf("serve on %s: %w", srv.Addr, err)
	case <-ctx.Done():
	}

	// 3) shutdown signal is received, perform cleanup
	if beforeShutdown != nil {
		beforeShutdown()
	}

	// 4) gracefully shutdown the server
	log.Printf("shutting down server at %s", srv.Addr)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), drainTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		srv.Close() // draining took too long: force-close remaining connections
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	return nil
}

// ServeGrpc runs the given grpc server on specified listener, until ctx is cancelled, then calls beforeShutdown
// (e.g. deregister) and drains in-flight requests within drainTimeout,
// force-closing if draining takes too long.
//
// beforeShutdown is also called if the server fails to start or stops serving
// unexpectedly, so callers can rely on it for cleanup.
// It can be nil.
func ServeGrpc(ctx context.Context, srv *grpc.Server, l net.Listener, drainTimeout time.Duration, beforeShutdown func()) error {

	// 1) start the server
	errCh := make(chan error, 1)
	go func() {
		log.Printf("starting server at %s", l.Addr())
		errCh <- srv.Serve(l)
	}()

	// 2) listen for server error or a shutdown signal (i.e. context cancel)
	select {
	case err := <-errCh:
		// 2-a) server failed to start, perform cleanup
		if beforeShutdown != nil {
			beforeShutdown()
		}
		return fmt.Errorf("serve on %s: %w", l.Addr(), err)
	case <-ctx.Done():
	}

	// 3) shutdown signal is received, perform cleanup
	if beforeShutdown != nil {
		beforeShutdown()
	}

	// 4) gracefully shutdown the server (or return after requested timeout)
	log.Printf("shutting down server at %s", l.Addr())
	done := make(chan struct{}, 1)
	go func() {
		srv.GracefulStop()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(drainTimeout):
		log.Printf("aborting graceful shutdown after drain timeout: %v", drainTimeout)
		srv.Stop()
		return nil
	}
}
