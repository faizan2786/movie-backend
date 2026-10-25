package rating

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/faizan2786/movie-backend/apiserver/internal/gateway"
	"github.com/faizan2786/movie-backend/apiserver/pkg/model"
)

// Gateway is a gateway to communicate with the metadata service.
type Gateway struct {
	// the http client
	client *http.Client

	// URL of the metadata service
	url *url.URL
}

// NewGateway returns a new [Gateway] to the metadata service running at given address.
func NewGateway(addr string) (*Gateway, error) {
	serviceURL, err := gateway.ParseURL(addr)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		client: &http.Client{
			Timeout: 5 * time.Second,
			// disable redirects (because, it is an internal service call)
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		url: serviceURL,
	}, nil
}

// GetMovieRating returns aggregated [model.Rating] for a movie with given id.
// It returns [gateway.ErrNotFound] if the rating doesn't exists.
func (g *Gateway) GetMovieRating(ctx context.Context, movieID model.ID) (model.RatingValue, error) {
	var rating model.RatingValue

	// make a GET request to /movies/<movieID> endpoint
	u := g.url.JoinPath("movies", url.PathEscape(string(movieID)))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return rating, fmt.Errorf("can not create a new HTTP request: %w", err)
	}

	log.Printf("sending GET request to %q", u)
	res, err := g.client.Do(req)
	if err != nil {
		return rating, fmt.Errorf("GET request error: %w", err)
	}
	defer res.Body.Close()

	// read the status code
	if res.StatusCode == http.StatusNotFound {
		return rating, gateway.ErrNotFound
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return rating, fmt.Errorf("non 2xx response code: %d; response: %v", res.StatusCode, res)
	}

	// read & decode the json response body
	err = json.NewDecoder(res.Body).Decode(&rating)
	if err != nil {
		return rating, fmt.Errorf("json response decode error: %w", err)
	}
	return rating, nil
}

// SubmitMovieRating submits movie ratings for users.
func (g *Gateway) SubmitMovieRating(ctx context.Context, movieID model.ID, rating model.Rating) error {

	// marshall rating to send it as a request body
	data, err := json.Marshal(rating)
	if err != nil {
		return fmt.Errorf("json request encode error: %w", err)
	}
	body := bytes.NewReader(data)

	// make a PUT request to /movies/<movieID> endpoint with above body
	u := g.url.JoinPath("movies", url.PathEscape(string(movieID)))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), body)
	if err != nil {
		return fmt.Errorf("can not create a new HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Printf("sending PUT request to %q", u)
	res, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("GET request error: %w", err)
	}
	defer res.Body.Close()

	// read the response status code
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("non 2xx response code: %d; response: %v", res.StatusCode, res)
	}

	return nil
}
