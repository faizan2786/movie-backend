package memory

import (
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/faizan2786/movie-backend/pkg/discovery"
	"github.com/stretchr/testify/require"
)

func InitRegistry(t *testing.T) *Registry {
	t.Helper()

	registry := New()
	registry.Register(t.Context(), "service-1", "instance-1", "my-service-1:3081")
	registry.Register(t.Context(), "service-1", "instance-2", "my-service-1:3082")
	registry.Register(t.Context(), "service-2", "instance-1", "my-service-2:3081")
	registry.Register(t.Context(), "service-1", "instance-3", "my-service-1:3083")
	registry.Register(t.Context(), "service-3", "instance-1", "my-service-3:3081")
	registry.Register(t.Context(), "service-2", "instance-2", "my-service-2:3082")

	return registry
}

func TestRegister(t *testing.T) {
	t.Parallel()

	registry := New()

	err := registry.Register(t.Context(), "service-1", "instance-1", "my-service-1:3080")
	require.NoError(t, err)

	err = registry.Register(t.Context(), "service-1", "instance-2", "my-service-1:3081")
	require.NoError(t, err)

	err = registry.Register(t.Context(), "service-2", "instance-1", "my-service-2:3080")
	require.NoError(t, err)

	err = registry.Register(t.Context(), "service-3", "instance-1", "my-service-3:3080")
	require.NoError(t, err)

	// address without a port
	err = registry.Register(t.Context(), "service-4", "instance-1", "my-service-1:")
	require.Error(t, err)

	// address without a host
	err = registry.Register(t.Context(), "service-5", "instance-1", ":5080")
	require.Error(t, err)
}

func TestDeRegister(t *testing.T) {
	t.Parallel()

	registry := InitRegistry(t)

	err := registry.Deregister(t.Context(), "service-1", "instance-1")
	require.NoError(t, err)

	err = registry.Deregister(t.Context(), "service-1", "instance-2")
	require.NoError(t, err)

	err = registry.Deregister(t.Context(), "service-2", "instance-1")
	require.NoError(t, err)

	err = registry.Deregister(t.Context(), "service-3", "instance-1")
	require.NoError(t, err)

	err = registry.Deregister(t.Context(), "service-1", "instance-2") // already de-registered
	require.ErrorIs(t, err, discovery.ErrNotFound)

	err = registry.Deregister(t.Context(), "service-3", "instance-2") // not registered
	require.ErrorIs(t, err, discovery.ErrNotFound)

}
func TestGetInstances(t *testing.T) {
	t.Parallel()

	registry := InitRegistry(t)

	addrs, err := registry.ServiceAddresses(t.Context(), "service-1")
	require.NoError(t, err)
	require.Len(t, addrs, 3)

	addrs, err = registry.ServiceAddresses(t.Context(), "service-2")
	require.NoError(t, err)
	require.Len(t, addrs, 2)

	addrs, err = registry.ServiceAddresses(t.Context(), "service-3")
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	_, err = registry.ServiceAddresses(t.Context(), "service-4") // service is not registered
	require.ErrorIs(t, err, discovery.ErrNotFound)
}

func TestReportHealthy(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		start := time.Now()

		registry := InitRegistry(t)
		activeT, _ := registry.getLastActiveTime("service-1", "instance-1")
		require.WithinDuration(t, start, activeT, 0)

		// perfrom healthcheck after 2 seconds and verify the last active time is updated
		time.Sleep(2 * time.Second)
		err := registry.ReportHealthy(t.Context(), "service-1", "instance-1")
		require.NoError(t, err)

		activeT, err = registry.getLastActiveTime("service-1", "instance-1")
		require.NoError(t, err)
		require.WithinDuration(t, start.Add(2*time.Second), activeT, 0)

		// verify the last active time of an older service
		activeT, err = registry.getLastActiveTime("service-2", "instance-1")
		require.NoError(t, err)
		require.WithinDuration(t, start, activeT, 0)
	})
}

func TestGetInstances_AfterDeregisteringAll(t *testing.T) {
	t.Parallel()

	registry := InitRegistry(t)
	require.NoError(t, registry.Deregister(t.Context(), "service-3", "instance-1"))

	// following deregistering all instances of a service,
	// ServiceAddresses should return an error (instead of an empty list)
	_, err := registry.ServiceAddresses(t.Context(), "service-3")
	require.ErrorIs(t, err, discovery.ErrNotFound)
}

func TestGetInstances_WithInActivity(t *testing.T) {
	t.Parallel()

	// advance time past maxInActivePeriod and verify the returned result
	synctest.Test(t, func(t *testing.T) {
		registry := InitRegistry(t)

		addrs, _ := registry.ServiceAddresses(t.Context(), "service-1")
		require.Len(t, addrs, 3)
		addrs, _ = registry.ServiceAddresses(t.Context(), "service-2")
		require.Len(t, addrs, 2)
		addrs, _ = registry.ServiceAddresses(t.Context(), "service-3")
		require.Len(t, addrs, 1)

		// perform health checks halfway through maxInActivePeriod
		time.Sleep(maxInActivePeriod / 2)
		registry.ReportHealthy(t.Context(), "service-1", "instance-1")
		registry.ReportHealthy(t.Context(), "service-1", "instance-2")
		registry.ReportHealthy(t.Context(), "service-2", "instance-1")

		// further advance just past the maxInActivePeriod mark from the start
		// (the above 3 reported instances will still be within their period)
		time.Sleep(maxInActivePeriod/2 + time.Second)

		// now, the number of active instances for service-1 should be 2,
		// active instances for service-2 should be 1 and
		// service-3 should return NotFound error
		addrs, _ = registry.ServiceAddresses(t.Context(), "service-1")
		require.Len(t, addrs, 2)
		addrs, _ = registry.ServiceAddresses(t.Context(), "service-2")
		require.Len(t, addrs, 1)

		_, err := registry.ServiceAddresses(t.Context(), "service-3")
		require.ErrorIs(t, err, discovery.ErrNotFound)
	})
}

func TestRegistry_Race(t *testing.T) {

	var registry Registry

	wg := sync.WaitGroup{}
	for range 1000 {
		wg.Go(func() {
			// register
			registry.Register(t.Context(), "service-1", "instance-1", "my-service-1:3081")
			registry.Register(t.Context(), "service-1", "instance-2", "my-service-1:3082")
			registry.Register(t.Context(), "service-2", "instance-1", "my-service-2:3081")
			registry.Register(t.Context(), "service-2", "instance-2", "my-service-2:3082")

			// deregister
			registry.Deregister(t.Context(), "service-1", "instance-1")
			registry.Deregister(t.Context(), "service-2", "instance-2")

			// ServiceAddresses
			registry.ServiceAddresses(t.Context(), "service-1")
			registry.ServiceAddresses(t.Context(), "service-2")

			// Health Check
			registry.ReportHealthy(t.Context(), "service-1", "instance-1")
			registry.ReportHealthy(t.Context(), "service-1", "instance-2")
			registry.ReportHealthy(t.Context(), "service-2", "instance-1")

			// ServiceAddresses
			registry.ServiceAddresses(t.Context(), "service-1")
			registry.ServiceAddresses(t.Context(), "service-2")
		})
	}

	wg.Wait()

	svc1, _ := registry.ServiceAddresses(t.Context(), "service-1")
	require.Contains(t, svc1, "my-service-1:3082")
	svc2, _ := registry.ServiceAddresses(t.Context(), "service-2")
	require.Contains(t, svc2, "my-service-2:3081")
}
