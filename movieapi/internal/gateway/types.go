package gateway

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/url"

	"github.com/faizan2786/movie-backend/pkg/discovery"
)

var (
	// ErrNotFound is returned when a resource is "not found".
	ErrNotFound = errors.New("not found")

	// ErrInvalidInput is returned when a service rejects a request as "not valid".
	ErrInvalidInput = errors.New("input not valid")

	// ErrInvalidURL is returned when a service URL is "not valid".
	ErrInvalidURL = errors.New("URL not valid")
)

// GetServiceURL returns a valid url of a specified service from the registry at random.
func GetServiceURL(ctx context.Context, service string, registry discovery.ServiceRegistry, scheme string) (*url.URL, error) {

	// get a service address from the registry at random
	addrs, err := registry.ServiceAddresses(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("cannot get an address for %q service: %w", service, err)
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("cannot get an address for %q service: %w", service, discovery.ErrNotFound)
	}

	addr := addrs[rand.IntN(len(addrs))]
	addr = fmt.Sprintf("%s://%s", scheme, addr)
	u, err := parseURL(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot get a valid URL for %q service: %w", service, err)
	}
	return u, nil
}

// parseURL parses given address as a valid service URL.
func parseURL(addr string) (*url.URL, error) {
	serviceURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, err)
	}

	if err := validateURL(serviceURL); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, err)
	}

	return serviceURL, nil
}

// validateURL validates a service URL.
func validateURL(u *url.URL) error {
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme must be http or https")
	}
	if u.Hostname() == "" {
		return errors.New("hostname can not be empty")
	}
	return nil
}
