package discovery

import (
	"context"
	"errors"
)

const (
	RatingServiceName   = "rating"
	MetadataServiceName = "metadata"
	MovieAPIServiceName = "movieapi"
)

// ErrNotFound is returned when no matching service or instance is registered.
var ErrNotFound = errors.New("not found")

// ServiceRegistry defines a service registry interface.
type ServiceRegistry interface {
	// Register registers a unique service instance with an address (host+port) to the registry.
	Register(ctx context.Context, service, instance, hostport string) error

	// Deregister de-registers a service instance from the registry
	// It returns [ErrNotFound] if the instance is not registered.
	Deregister(ctx context.Context, service, instance string) error

	// ServiceAddresses returns a list of address (host+port) for healthy instances of the specified service
	// It returns [ErrNotFound] if service instance has no healthy instances.
	ServiceAddresses(ctx context.Context, service string) ([]string, error)

	// ReportHealthy marks an instance as healthy. Instances that don't report
	// within an implementation-defined interval are excluded from
	// [ServiceRegistry.ServiceAddresses].
	// It returns [ErrNotFound] if the instance is not registered.
	ReportHealthy(ctx context.Context, service, instance string) error
}
