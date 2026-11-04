package metadata_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/faizan2786/movie-backend/movieapi/internal/gateway"
	"github.com/faizan2786/movie-backend/movieapi/internal/gateway/metadata"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
	"github.com/faizan2786/movie-backend/pkg/discovery"
	"github.com/faizan2786/movie-backend/pkg/discovery/memory"
	"github.com/stretchr/testify/require"
)

const testID = "123"

var testData = model.Metadata{
	ID:          testID,
	Title:       "Test Movie Title",
	Description: "testing metadata",
	ReleaseYear: 2026,
}

func registryWithTestServer(t *testing.T, statusCode int, sendInvalidResp bool) discovery.ServiceRegistry {
	t.Helper()

	const expectedPath = "/" + testID
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// verify that the expected request is received from the gateway
		if r.Method != http.MethodGet || r.URL.Path != expectedPath {
			t.Errorf("got request %s %s; want GET %s", r.Method, r.URL.Path, expectedPath)
		}

		// write the given response code
		w.WriteHeader(statusCode)
		if sendInvalidResp {
			json.NewEncoder(w).Encode(`{}`) // send invalid data
		} else {
			json.NewEncoder(w).Encode(testData)
		}
	}
	server := httptest.NewTestServer(t, handler)
	server.Start()
	t.Cleanup(server.Close)

	// register the server's address (host+port) with the registry
	reg := memory.New()
	_, serverHostPort, _ := strings.Cut(server.URL, "://")
	err := reg.Register(t.Context(), discovery.MetadataServiceName, discovery.GenerateInstanceID(), serverHostPort)
	if err != nil {
		t.Fatal(err)
	}
	return reg
}

func TestGetMetadataValid(t *testing.T) {
	reg := registryWithTestServer(t, http.StatusOK, false)
	gw := metadata.NewGateway(reg)

	data, err := gw.GetMetadata(t.Context(), testID)
	require.NoError(t, err)
	require.Equal(t, testData, data)
}

func TestGetMetadataWithErrors(t *testing.T) {
	testCases := []struct {
		name        string
		code        int     // response code that needs to be sent by the server
		invalidResp bool    // if server needs to send an invalid response?
		wantErr     error   // which error to expect
		wantNotErrs []error // which errors NOT to expect
	}{
		{
			name:    "not found",
			code:    http.StatusNotFound,
			wantErr: gateway.ErrNotFound,
		},
		{
			name:        "not authorized",
			code:        http.StatusUnauthorized,
			wantNotErrs: []error{gateway.ErrNotFound},
		},
		{
			name:        "server error",
			code:        http.StatusInternalServerError,
			wantNotErrs: []error{gateway.ErrNotFound},
		},
		{
			name:        "invalid response data",
			code:        http.StatusOK,
			invalidResp: true,
			wantNotErrs: []error{gateway.ErrNotFound},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reg := registryWithTestServer(t, tc.code, tc.invalidResp)
			gw := metadata.NewGateway(reg)

			_, err := gw.GetMetadata(t.Context(), testID)
			require.Error(t, err)

			// check expected error
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				t.Log(err)
			}

			// check unexpected errors
			for _, e := range tc.wantNotErrs {
				require.NotErrorIs(t, err, e)
			}

		})
	}
}
