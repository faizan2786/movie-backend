package consul

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/faizan2786/movie-backend/pkg/discovery"
	consulapi "github.com/hashicorp/consul/api"
)

// Registry is a Consul implementation of [discovery.ServiceRegistry].
type Registry struct {
	client *consulapi.Client
}

// New returns a new instance of [Registry].
func New(addr string) (*Registry, error) {
	config := consulapi.DefaultConfig()
	config.Address = addr
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create consul client: %w", err)
	}

	return &Registry{
		client: client,
	}, nil
}

// Register implements [discovery.ServiceRegistry].
func (r *Registry) Register(_ context.Context, service string, instance string, hostport string) error {
	host, port, err := discovery.ParseHostPort(hostport)
	if err != nil {
		return err
	}

	// register service with the local consul agent (i.e. the consul agent running locally)
	err = r.client.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		Name:    service,
		ID:      instance,
		Address: host,
		Port:    port,
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        instance,
			TTL:                            "5s", // check turns critical if no TTL update (heartbeat) arrives within this time
			Status:                         consulapi.HealthPassing,
			DeregisterCriticalServiceAfter: "1m", // deregister the instance once its check has been critical for 1m
		},
	})
	if err != nil {
		return fmt.Errorf("register service %q instance %q: %w", service, instance, err)
	}

	return nil
}

// Deregister implements [discovery.ServiceRegistry].
func (r *Registry) Deregister(ctx context.Context, service string, instance string) error {
	err := r.client.Agent().ServiceDeregister(instance)
	if isNotFoundErr(err) {
		return fmt.Errorf("deregister instance %q: %w", instance, discovery.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("deregister service instance %q: %w", instance, err)
	}
	return nil
}

// ServiceAddresses implements [discovery.ServiceRegistry].
func (r *Registry) ServiceAddresses(ctx context.Context, service string) ([]string, error) {
	var addrs []string

	// get entries for healthy services
	entries, _, err := r.client.Health().Service(service, "", true /* passingOnly */, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch healthy services: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no healthy instances of service %q: %w", service, discovery.ErrNotFound)
	}

	for _, e := range entries {
		instance := e.Service
		addrs = append(addrs, fmt.Sprintf("%s:%d", instance.Address, instance.Port))
	}
	return addrs, nil
}

// ReportHealthy implements [discovery.ServiceRegistry].
func (r *Registry) ReportHealthy(ctx context.Context, service string, instance string) error {
	err := r.client.Agent().UpdateTTL(instance, "", consulapi.HealthPassing)
	if isNotFoundErr(err) {
		return fmt.Errorf("heport instance %q health: %w", instance, discovery.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("instance %q healthcheck for %q service: %w", instance, service, err)
	}
	return nil
}

// check if an error is an HTTP NotFound error
func isNotFoundErr(err error) bool {
	var se consulapi.StatusError
	return errors.As(err, &se) && se.Code == http.StatusNotFound
}
