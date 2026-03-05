package clients

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dialer provides cached gRPC client connections to backend services.
// Addresses are read from environment variables: <SERVICE>_GRPC_ADDR (e.g., ORDERS_GRPC_ADDR).
// No endpoints are hardcoded.
//
// Example:
//   ORDERS_GRPC_ADDR=localhost:50051
//   CATALOG_GRPC_ADDR=localhost:50052
//
// Use GetConn(ctx, "orders") to retrieve a *grpc.ClientConn.

type Dialer struct {
	mu    sync.RWMutex
	pool  map[string]*grpc.ClientConn
	opts  []grpc.DialOption
	dTout time.Duration
}

// NewDialer creates a new Dialer with sane defaults.
func NewDialer() *Dialer {
	return &Dialer{
		pool: make(map[string]*grpc.ClientConn),
		opts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"pick_first"}`),
		},
		dTout: 5 * time.Second,
	}
}

// GetConn returns a cached *grpc.ClientConn for a service. The service name is transformed
// to its environment variable key: <SERVICE>_GRPC_ADDR. For example, "orders" -> ORDERS_GRPC_ADDR.
func (d *Dialer) GetConn(ctx context.Context, service string) (*grpc.ClientConn, error) {
	key := strings.TrimSpace(service)
	if key == "" {
		return nil, fmt.Errorf("service name is required")
	}

	addrEnv := envKeyForService(key)
	addr := strings.TrimSpace(os.Getenv(addrEnv))
	if addr == "" {
		return nil, fmt.Errorf("missing %s for service %q", addrEnv, key)
	}

	// Fast path: read lock pool
	d.mu.RLock()
	if cc, ok := d.pool[key]; ok && cc != nil {
		d.mu.RUnlock()
		return cc, nil
	}
	d.mu.RUnlock()

	// Establish a new connection with timeout and validation
	dialCtx, cancel := context.WithTimeout(ctx, d.dTout)
	defer cancel()

	// Pre-validate the address format (host:port)
	if _, _, err := net.SplitHostPort(addr); err != nil {
		// If no port present, let grpc handle it; we still try to dial
	}

	cc, err := grpc.DialContext(
		dialCtx,
		addr,
		append(d.opts, grpc.WithBlock())...,
	)
	if err != nil {
		return nil, fmt.Errorf("dial %s (%s): %w", key, addr, err)
	}

	// Cache the connection
	d.mu.Lock()
	d.pool[key] = cc
	d.mu.Unlock()

	return cc, nil
}

// CloseAll closes all cached connections.
func (d *Dialer) CloseAll() error {
	var firstErr error
	d.mu.Lock()
	for k, cc := range d.pool {
		if cc != nil {
			if err := cc.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		delete(d.pool, k)
	}
	d.mu.Unlock()
	return firstErr
}

func envKeyForService(service string) string {
	upper := strings.ToUpper(service)
	upper = strings.ReplaceAll(upper, "-", "_")
	upper = strings.ReplaceAll(upper, ".", "_")
	upper = strings.ReplaceAll(upper, "/", "_")
	return upper + "_GRPC_ADDR"
}
