package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun_ShutsDownOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	srv := &http.Server{Addr: "localhost:0", Handler: http.NewServeMux()}

	called := false
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, srv, time.Second, func() { called = true })
	}()

	cancel()

	select {
	case err := <-done:
		require.NoError(t, err)
		require.True(t, called, "beforeShutdown should be called on cancellation")
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestRun_ReturnsErrorWhenPortInUse(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	defer ln.Close()

	srv := &http.Server{Addr: ln.Addr().String(), Handler: http.NewServeMux()}

	called := false
	err = Run(context.Background(), srv, time.Second, func() { called = true })
	require.Error(t, err)
	require.True(t, called, "beforeShutdown should be called when the server fails to start")
}
