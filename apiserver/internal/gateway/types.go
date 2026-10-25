package gateway

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	// ErrNotFound is returned when a resource is "not found".
	ErrNotFound = errors.New("not found")

	// ErrInvalidURL is returned when a service URL is "not valid".
	ErrInvalidURL = errors.New("URL not valid")
)

// ParseURL parses given address as a valid service URL.
func ParseURL(addr string) (*url.URL, error) {
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
