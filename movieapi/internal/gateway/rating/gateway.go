package rating

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/faizan2786/movie-backend/movieapi/internal/gateway"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
	"github.com/faizan2786/movie-backend/pkg/discovery"
)

const defaultScheme string = "http"

// Gateway is a gateway to communicate with the rating service.
type Gateway struct {
	// the http client
	client *http.Client

	// registry to discover the rating service instances
	registry discovery.ServiceRegistry
}

// NewGateway returns a new [Gateway] to the rating service discovered via given registry.
func NewGateway(r discovery.ServiceRegistry) *Gateway {
	return &Gateway{
		client: &http.Client{
			Timeout: 5 * time.Second,
			// disable redirects (because, it is an internal service call)
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		registry: r,
	}
}

// GetMovieRating returns aggregated [model.Rating] for a movie with given id.
// It returns [gateway.ErrNotFound] if the rating doesn't exists.
func (g *Gateway) GetMovieRating(ctx context.Context, movieID model.ID) (model.RatingValue, error) {
	var rating model.RatingValue

	serviceURL, err := gateway.GetServiceURL(ctx, discovery.RatingServiceName, g.registry, defaultScheme)
	if err != nil {
		return rating, err
	}

	// make a GET request to /movies/<movieID> endpoint
	u := serviceURL.JoinPath("movies", url.PathEscape(string(movieID)))
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
	err = json.UnmarshalRead(res.Body, &rating)
	if err != nil {
		return rating, fmt.Errorf("json response decode error: %w", err)
	}
	return rating, nil
}

// SubmitMovieRating submits movie ratings for users.
// It returns [gateway.ErrInvalidInput] if the rating service rejects the rating as not valid.
func (g *Gateway) SubmitMovieRating(ctx context.Context, movieID model.ID, rating model.Rating) error {
	// marshall rating to send it as a request body
	data, err := json.Marshal(rating)
	if err != nil {
		return fmt.Errorf("json request encode error: %w", err)
	}
	body := bytes.NewReader(data)

	serviceURL, err := gateway.GetServiceURL(ctx, discovery.RatingServiceName, g.registry, defaultScheme)
	if err != nil {
		return err
	}

	// make a PUT request to /movies/<movieID> endpoint with above body
	u := serviceURL.JoinPath("movies", url.PathEscape(string(movieID)))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), body)
	if err != nil {
		return fmt.Errorf("can not create a new HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Printf("sending PUT request to %q", u)
	res, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("PUT request error: %w", err)
	}
	defer res.Body.Close()

	// read the status code
	if res.StatusCode == http.StatusBadRequest {
		return gateway.ErrInvalidInput
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("non 2xx response code: %d; response: %v", res.StatusCode, res)
	}
	return nil
}
