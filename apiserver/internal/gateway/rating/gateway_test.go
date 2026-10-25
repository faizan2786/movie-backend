package rating_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faizan2786/movie-backend/apiserver/internal/gateway"
	"github.com/faizan2786/movie-backend/apiserver/internal/gateway/rating"
	"github.com/faizan2786/movie-backend/apiserver/pkg/model"
	"github.com/stretchr/testify/require"
)

const (
	testMovieID = "123"

	testDataGET = model.RatingValue(7.5)
)

var testDataPUT = model.Rating{
	UserID: model.ID("1"),
	Value:  model.RatingValue(6),
}

// test server for GET requests
func newTestServerGET(t *testing.T, statusCode int, sendInvalidResp bool) *httptest.Server {
	t.Helper()

	const expectedPath = "/movies/" + testMovieID
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// verify that the expected request is received from the gateway
		if r.Method != http.MethodGet || r.URL.Path != expectedPath {
			t.Errorf("got request %s %s; want GET %s", r.Method, r.URL.Path, expectedPath)
		}

		// write the given response code
		w.WriteHeader(statusCode)
		if sendInvalidResp {
			json.NewEncoder(w).Encode(`{`) // send invalid data
		} else {
			json.NewEncoder(w).Encode(testDataGET)
		}
	}
	server := httptest.NewTestServer(t, handler)
	server.Start()
	t.Cleanup(server.Close)
	return server
}

// test server for PUT requests
func newTestServerPUT(t *testing.T, statusCode int) *httptest.Server {
	t.Helper()

	const expectedPath = "/movies/" + testMovieID

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// verify that the expected request is received from the gateway
		if r.Method != http.MethodPut || r.URL.Path != expectedPath {
			t.Errorf("got request %s %s; want PUT %s", r.Method, r.URL.Path, expectedPath)
		}

		// verify the request body is correct
		var data model.Rating
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			t.Errorf("got: PUT request body decode error = %s; want: no errors", err)
		}
		defer r.Body.Close()

		if data != testDataPUT {
			t.Errorf("got PUT request body = %+v; want = %+v", data, testDataPUT)
		}

		// write the given response code
		w.WriteHeader(statusCode)
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
			addr:    "http://localhost:8082",
			wantErr: nil,
		},
		{
			name:    "invalid address",
			addr:    "localhost:8082",
			wantErr: gateway.ErrInvalidURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := rating.NewGateway(tc.addr)
			if tc.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}

func TestGetMovieRatingValid(t *testing.T) {
	server := newTestServerGET(t, http.StatusOK, false)
	gw, err := rating.NewGateway(server.URL)
	require.NoError(t, err)

	data, err := gw.GetMovieRating(t.Context(), testMovieID)
	require.NoError(t, err)
	require.Equal(t, testDataGET, data)
}

func TestGetMovieRatingWithErrors(t *testing.T) {
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
			server := newTestServerGET(t, tc.code, tc.invalidResp)
			gw, err := rating.NewGateway(server.URL)
			require.NoError(t, err)

			_, err = gw.GetMovieRating(t.Context(), testMovieID)
			require.Error(t, err)

			// check expected error
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				//t.Log(err)
			}

			// check unexpected errors
			for _, e := range tc.wantNotErrs {
				require.NotErrorIs(t, err, e)
			}

		})
	}
}

func TestSubmitMovieRatingValid(t *testing.T) {
	server := newTestServerPUT(t, http.StatusOK)
	gw, err := rating.NewGateway(server.URL)
	require.NoError(t, err)

	err = gw.SubmitMovieRating(t.Context(), testMovieID, testDataPUT)
	require.NoError(t, err)
}

func TestPUTMovieRatingWithErrors(t *testing.T) {
	testCases := []struct {
		name string
		code int // response code that needs to be sent by the server
	}{
		{
			name: "not found",
			code: http.StatusNotFound,
		},
		{
			name: "not authorized",
			code: http.StatusUnauthorized,
		},
		{
			name: "server error",
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServerPUT(t, tc.code)
			gw, err := rating.NewGateway(server.URL)
			require.NoError(t, err)

			err = gw.SubmitMovieRating(t.Context(), testMovieID, testDataPUT)
			require.Error(t, err)

			// check error logs
			//t.Log(err)
		})
	}
}
