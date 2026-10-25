package metadata_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faizan2786/movie-backend/apiserver/internal/gateway"
	"github.com/faizan2786/movie-backend/apiserver/internal/gateway/metadata"
	"github.com/faizan2786/movie-backend/apiserver/pkg/model"
	"github.com/stretchr/testify/require"
)

const testID = "123"

var testData = model.Metadata{
	ID:          testID,
	Title:       "Test Movie Title",
	Description: "testing metadata",
	ReleaseYear: 2026,
}

func getTestServer(t *testing.T, statusCode int, sendInvalidResp bool) *httptest.Server {
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
	return server
}

func TestNewGateway(t *testing.T) {
	testCases := []struct {
		name    string
		addr    string
		wantErr error
	}{
		{
			name:    "valid address",
			addr:    "http://localhost:8081",
			wantErr: nil,
		},
		{
			name:    "invalid address",
			addr:    "localhost:8081",
			wantErr: gateway.ErrInvalidURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := metadata.NewGateway(tc.addr)
			if tc.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}

func TestGetMetadataValid(t *testing.T) {
	server := getTestServer(t, http.StatusOK, false)
	gw, err := metadata.NewGateway(server.URL)
	require.NoError(t, err)

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
			server := getTestServer(t, tc.code, tc.invalidResp)
			gw, err := metadata.NewGateway(server.URL)
			require.NoError(t, err)

			_, err = gw.GetMetadata(t.Context(), testID)
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
