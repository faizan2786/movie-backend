package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/faizan2786/movie-backend/pkg/discovery"
)

var maxInActivePeriod = 300 * time.Second

// instance is a service instance
type instance struct {
	address    string
	lastActive time.Time // time in UNIX-epoch when this instance was last active
}

// Registry is an in-memory implementation of [discovery.ServiceRegistry].
// The zero value is an empty registry ready to use.
type Registry struct {
	mu   sync.RWMutex
	data map[string]map[string]*instance // we use pointer to instance struct so that
	//  we can update its field (i.e. lastActive) in-place.
}

// New returns a new instance of [Registry].
func New() *Registry {
	return &Registry{
		data: map[string]map[string]*instance{},
	}
}

func (r *Registry) Register(_ context.Context, serviceName string, instanceId string, hostport string) error {
	_, _, err := discovery.ParseHostPort(hostport)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data == nil {
		r.data = map[string]map[string]*instance{}
	}

	if r.data[serviceName] == nil {
		r.data[serviceName] = map[string]*instance{}
	}

	r.data[serviceName][instanceId] = &instance{
		address:    hostport,
		lastActive: time.Now(),
	}
	return nil
}

func (r *Registry) Deregister(_ context.Context, serviceName string, instanceId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data[serviceName] == nil {
		return fmt.Errorf("%q service is not registered: %w", serviceName, discovery.ErrNotFound)
	}
	if _, ok := r.data[serviceName][instanceId]; !ok {
		return fmt.Errorf("instance %q for service %q is not registered: %w", instanceId, serviceName, discovery.ErrNotFound)
	}

	// delete the instance
	delete(r.data[serviceName], instanceId)

	// delete the service when no instances left
	if len(r.data[serviceName]) == 0 {
		delete(r.data, serviceName)
	}

	return nil
}

// ServiceAddresses returns active instances of a service
func (r *Registry) ServiceAddresses(_ context.Context, serviceName string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.data[serviceName] == nil {
		return nil, fmt.Errorf("%q service is not registered: %w", serviceName, discovery.ErrNotFound)
	}

	// only return instances seen active within last maxInActivePeriod
	var out []string
	for _, instance := range r.data[serviceName] {
		if time.Since(instance.lastActive) <= maxInActivePeriod {
			out = append(out, instance.address)
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no healthy instances of service %q: %w", serviceName, discovery.ErrNotFound)
	}

	return out, nil
}

func (r *Registry) ReportHealthy(_ context.Context, serviceName string, instanceId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data[serviceName] == nil {
		return fmt.Errorf("%q service is not registered: %w", serviceName, discovery.ErrNotFound)
	}
	if _, ok := r.data[serviceName][instanceId]; !ok {
		return fmt.Errorf("instance %q for service %q is not registered: %w", instanceId, serviceName, discovery.ErrNotFound)
	}

	// update its lastActive time
	r.data[serviceName][instanceId].lastActive = time.Now()
	return nil
}

// helper method for testing
func (r *Registry) getLastActiveTime(serviceName string, instanceId string) (time.Time, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.data[serviceName] == nil {
		return time.Time{}, fmt.Errorf("%q service is not registered: %w", serviceName, discovery.ErrNotFound)
	}
	if _, ok := r.data[serviceName][instanceId]; !ok {
		return time.Time{}, fmt.Errorf("instance %q for service %q is not registered: %w", instanceId, serviceName, discovery.ErrNotFound)
	}

	return r.data[serviceName][instanceId].lastActive, nil

}
