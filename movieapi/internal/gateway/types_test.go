package gateway_test

import (
	"context"
	"testing"

	"github.com/faizan2786/movie-backend/movieapi/internal/gateway"
	"github.com/faizan2786/movie-backend/pkg/discovery"
	"github.com/faizan2786/movie-backend/pkg/discovery/memory"
	"github.com/stretchr/testify/require"
)

const testService = "test-service"

// emptyRegistry is a [discovery.ServiceRegistry] whose ServiceAddresses
// returns an empty list without an error.
type emptyRegistry struct {
	discovery.ServiceRegistry
}

func (emptyRegistry) ServiceAddresses(context.Context, string) ([]string, error) {
	return nil, nil
}

func TestGetServiceURL(t *testing.T) {
	reg := memory.New()
	err := reg.Register(t.Context(), testService, discovery.GenerateInstanceID(), "localhost:8081")
	require.NoError(t, err)

	u, err := gateway.GetServiceURL(t.Context(), testService, reg, "http")
	require.NoError(t, err)
	require.Equal(t, "http://localhost:8081", u.String())
}

func TestGetServiceURL_WithErrors(t *testing.T) {
	testCases := []struct {
		name     string
		registry func(t *testing.T) discovery.ServiceRegistry
		scheme   string
		wantErr  error
	}{
		{
			name: "service not registered",
			registry: func(t *testing.T) discovery.ServiceRegistry {
				return memory.New()
			},
			scheme:  "http",
			wantErr: discovery.ErrNotFound,
		},
		{
			name: "registry returns no addresses",
			registry: func(t *testing.T) discovery.ServiceRegistry {
				return emptyRegistry{}
			},
			scheme:  "http",
			wantErr: discovery.ErrNotFound,
		},
		{
			name: "address is not a valid URL host",
			registry: func(t *testing.T) discovery.ServiceRegistry {
				reg := memory.New()
				err := reg.Register(t.Context(), testService, discovery.GenerateInstanceID(), "bad host:8081")
				require.NoError(t, err)
				return reg
			},
			scheme:  "http",
			wantErr: gateway.ErrInvalidURL,
		},
		{
			name: "unsupported scheme",
			registry: func(t *testing.T) discovery.ServiceRegistry {
				reg := memory.New()
				err := reg.Register(t.Context(), testService, discovery.GenerateInstanceID(), "localhost:8081")
				require.NoError(t, err)
				return reg
			},
			scheme:  "ftp",
			wantErr: gateway.ErrInvalidURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := gateway.GetServiceURL(t.Context(), testService, tc.registry(t), tc.scheme)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
