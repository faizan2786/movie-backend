package metadata

import (
	"context"
	"encoding/json"
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

// Gateway is a gateway to communicate with the metadata service.
type Gateway struct {
	// the http client
	client *http.Client

	// registry to discover the metadata service instances
	registry discovery.ServiceRegistry
}

// NewGateway returns a new [Gateway] to the metadata service discovered via given registry.
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

// GetMetadata returns [model.Metadata] for a movie with given id.
// It returns [gateway.ErrNotFound] if the movie doesn't exists.
func (g *Gateway) GetMetadata(ctx context.Context, id model.ID) (model.Metadata, error) {
	var data model.Metadata

	serviceURL, err := gateway.GetServiceURL(ctx, discovery.MetadataServiceName, g.registry, defaultScheme)
	if err != nil {
		return data, err
	}

	// make a GET request with ID to the service
	u := serviceURL.JoinPath(url.PathEscape(string(id)))
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
