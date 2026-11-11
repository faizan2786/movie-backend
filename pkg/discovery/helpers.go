package discovery

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
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
	log.Printf("register instance for %s service", serviceName)
	err := r.Register(ctx, serviceName, instanceID, hostPort)
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
		log.Printf("deregister instance for %s service", serviceName)
		if err := r.Deregister(ctx, serviceName, instanceID); err != nil {
			log.Printf("deregister: %v", err)
		}
	}

	return deregister, nil
}

func ParseHostPort(hostport string) (host string, port int, err error) {
	host, p, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse address: %w", ErrInvalidHostPort)
	}
	if host == "" || p == "" {
		return "", 0, fmt.Errorf("address %q must have both host and port: %w", hostport, ErrInvalidHostPort)
	}

	port, err = strconv.Atoi(p)
	if err != nil {
		return "", 0, fmt.Errorf("port %q must be a numeric value: %w", p, ErrInvalidHostPort)
	}

	return host, port, nil
}
