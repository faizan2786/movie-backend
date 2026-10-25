package metadata

import (
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

// GetMetadata returns [model.Metadata] for a movie with given id.
// It returns [gateway.ErrNotFound] if the movie doesn't exists.
func (g *Gateway) GetMetadata(ctx context.Context, id model.ID) (model.Metadata, error) {
	var data model.Metadata

	// make a GET request with ID to the service
	u := g.url.JoinPath(url.PathEscape(string(id)))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return data, fmt.Errorf("can not create a new HTTP request: %w", err)
	}

	log.Printf("sending GET request to %q", u)
	res, err := g.client.Do(req)
	if err != nil {
		return data, fmt.Errorf("GET request error: %w", err)
	}
	defer res.Body.Close()

	// read the status code
	if res.StatusCode == http.StatusNotFound {
		return data, gateway.ErrNotFound
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return data, fmt.Errorf("non 2xx response code: %d; response: %v", res.StatusCode, res)
	}

	// read & decode the json response body
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("json response decode error: %w", err)
	}
	return data, nil
}
