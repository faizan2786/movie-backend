package discovery

import (
	"context"
	"log"
	"time"
	"uuid"
)

// GenerateInstanceID returns a unique instance ID.
func GenerateInstanceID() string {
	return uuid.New().String()
}

// DeregisterFunc is a function to deregister a service from [ServiceRegistry].
type DeregisterFunc func(context.Context)

// RegisterWithHeartbeat registers the service and sends periodic heartbeats with specified period.
// It returns a [DeregisterFunc].
func RegisterWithHeartbeat(ctx context.Context, r ServiceRegistry, serviceName string, hostPort string, period time.Duration) (DeregisterFunc, error) {
	// register the service with registry with periodic health-checks
	instanceID := GenerateInstanceID()
	err := r.Register(context.Background(), serviceName, instanceID, hostPort)
	if err != nil {
		return nil, err
	}

	// send heartbeat
	go func() {
		tCh := time.Tick(period)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tCh:
				err := r.ReportHealthy(ctx, serviceName, instanceID)
				if err != nil {
					log.Printf("%s failed to register heartbeat: %s", serviceName, err)
				}
			}
		}
	}()

	deregister := func(ctx context.Context) {
		log.Printf("deregister instance for %s", serviceName)
		if err := r.Deregister(ctx, serviceName, instanceID); err != nil {
			log.Printf("deregister: %v", err)
		}
	}

	return deregister, nil
}
