package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ConnectionPool manages multiple gRPC connections per service for load distribution
type ConnectionPool struct {
	serviceName string
	address     string
	poolSize    int

	connections         []*PooledConnection
	mu                  sync.RWMutex
	nextIdx             uint64
	healthCheckInterval time.Duration

	dialOpts []grpc.DialOption
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// PooledConnection wraps a gRPC connection with health tracking
type PooledConnection struct {
	conn      *grpc.ClientConn
	healthy   atomic.Bool
	lastCheck time.Time
	mu        sync.RWMutex
}

// PoolConfig configures connection pool behavior
type PoolConfig struct {
	PoolSize            int           // Number of connections to maintain
	MaxConnAge          time.Duration // Maximum connection lifetime
	KeepAliveTime       time.Duration // Keepalive ping interval
	KeepAliveTimeout    time.Duration // Keepalive timeout
	HealthCheckInterval time.Duration // Health check frequency
}

// DefaultPoolConfig returns production-ready pool configuration
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		PoolSize:            2,                 // Further reduced to prevent ping flooding
		MaxConnAge:          15 * time.Minute,  // Shorter lifetime to prevent stale connections
		KeepAliveTime:       120 * time.Second, // Much longer keepalive (2 minutes)
		KeepAliveTimeout:    30 * time.Second,  // Longer timeout
		HealthCheckInterval: 60 * time.Second,  // Less frequent health checks
	}
}

// NewConnectionPool creates a connection pool for a service
func NewConnectionPool(serviceName, address string, cfg *PoolConfig) *ConnectionPool {
	if cfg == nil {
		cfg = DefaultPoolConfig()
	}

	pool := &ConnectionPool{
		serviceName:         serviceName,
		address:             address,
		poolSize:            cfg.PoolSize,
		healthCheckInterval: cfg.HealthCheckInterval,
		connections:         make([]*PooledConnection, 0, cfg.PoolSize),
		dialOpts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                cfg.KeepAliveTime,
				Timeout:             cfg.KeepAliveTimeout,
				PermitWithoutStream: false, // 🚨 CRITICAL: Disable ping without stream to prevent "too_many_pings"
			}),
			grpc.WithDefaultServiceConfig(`{
				"loadBalancingConfig": [{"round_robin":{}}],
				"methodConfig": [{
					"name": [{"service": ""}],
					"retryPolicy": {
						"maxAttempts": 2,
						"initialBackoff": "0.5s",
						"maxBackoff": "2s",
						"backoffMultiplier": 2,
						"retryableStatusCodes": ["UNAVAILABLE"]
					}
				}]
			}`),
			grpc.WithMaxMsgSize(10 * 1024 * 1024), // 10MB
		},
		stopCh: make(chan struct{}),
	}

	return pool
}

// Initialize creates all connections in the pool
func (p *ConnectionPool) Initialize(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	logger.Info("Initializing connection pool",
		zap.String("service", p.serviceName),
		zap.String("address", p.address),
		zap.Int("pool_size", p.poolSize))

	// Create connections concurrently for faster startup
	var wg sync.WaitGroup
	errCh := make(chan error, p.poolSize)
	connCh := make(chan *PooledConnection, p.poolSize)

	for i := 0; i < p.poolSize; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			conn, err := p.createConnection(ctx, idx)
			if err != nil {
				errCh <- fmt.Errorf("connection %d: %w", idx, err)
				return
			}
			connCh <- conn
		}(i)
	}

	wg.Wait()
	close(errCh)
	close(connCh)

	// Collect results
	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	for conn := range connCh {
		p.connections = append(p.connections, conn)
	}

	if len(p.connections) == 0 {
		logger.Error("Failed to create any connections", zap.Any("errors", errors))
		return fmt.Errorf("failed to create any connections: %v", errors)
	}

	if len(errors) > 0 {
		logger.Warn("Some connections failed to initialize",
			zap.String("service", p.serviceName),
			zap.Int("failed", len(errors)),
			zap.Int("succeeded", len(p.connections)))
	}

	logger.Info("Connection pool initialized",
		zap.String("service", p.serviceName),
		zap.Int("connections", len(p.connections)))

	// Start health monitoring
	p.startHealthMonitoring()

	return nil
}

// createConnection creates a single pooled connection
func (p *ConnectionPool) createConnection(ctx context.Context, idx int) (*PooledConnection, error) {
	// Use non-blocking dial for faster startup
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, p.address, p.dialOpts...)
	if err != nil {
		logger.Error("Dial failed", zap.Error(err), zap.String("address", p.address))
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	pooled := &PooledConnection{
		conn:      conn,
		lastCheck: time.Now(),
	}
	pooled.healthy.Store(true)

	logger.Debug("Connection created",
		zap.String("service", p.serviceName),
		zap.Int("index", idx))

	return pooled, nil
}

// GetConnection returns a healthy connection using round-robin
func (p *ConnectionPool) GetConnection() (*grpc.ClientConn, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.connections) == 0 {
		logger.Error("No connections available", zap.String("service", p.serviceName))
		return nil, fmt.Errorf("no connections available for %s", p.serviceName)
	}

	// Round-robin with health checking
	attempts := len(p.connections)
	for i := 0; i < attempts; i++ {
		idx := atomic.AddUint64(&p.nextIdx, 1) % uint64(len(p.connections))
		pooled := p.connections[idx]

		if pooled.healthy.Load() {
			return pooled.conn, nil
		}
	}

	// No healthy connections, return any connection (might recover)
	idx := atomic.LoadUint64(&p.nextIdx) % uint64(len(p.connections))
	return p.connections[idx].conn, nil
}

// GetAllConnections returns all healthy connections for fan-out requests
func (p *ConnectionPool) GetAllConnections() []*grpc.ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var conns []*grpc.ClientConn
	for _, pooled := range p.connections {
		if pooled.healthy.Load() {
			conns = append(conns, pooled.conn)
		}
	}
	return conns
}

// startHealthMonitoring monitors connection health
func (p *ConnectionPool) startHealthMonitoring() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		ticker := time.NewTicker(p.healthCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.checkHealth()
			case <-p.stopCh:
				return
			}
		}
	}()
}

// checkHealth checks all connections
func (p *ConnectionPool) checkHealth() {
	p.mu.RLock()
	connections := make([]*PooledConnection, len(p.connections))
	copy(connections, p.connections)
	p.mu.RUnlock()

	var healthyCount int
	for i, pooled := range connections {
		state := pooled.conn.GetState()
		healthy := state == connectivity.Ready || state == connectivity.Idle

		pooled.mu.Lock()
		pooled.lastCheck = time.Now()
		pooled.mu.Unlock()
		pooled.healthy.Store(healthy)

		if healthy {
			healthyCount++
		} else {
			logger.Debug("Unhealthy connection detected",
				zap.String("service", p.serviceName),
				zap.Int("index", i),
				zap.String("state", state.String()))
		}
	}

	if healthyCount == 0 && len(connections) > 0 {
		logger.Debug("No healthy connections — service may not be running",
			zap.String("service", p.serviceName),
			zap.Int("total", len(connections)))
	}
}

// HealthyCount returns the number of healthy connections
func (p *ConnectionPool) HealthyCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	count := 0
	for _, pooled := range p.connections {
		if pooled.healthy.Load() {
			count++
		}
	}
	return count
}

// TotalCount returns total connections
func (p *ConnectionPool) TotalCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.connections)
}

// Close closes all connections in the pool
func (p *ConnectionPool) Close() error {
	close(p.stopCh)
	p.wg.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()

	logger.Info("Closing connection pool",
		zap.String("service", p.serviceName),
		zap.Int("connections", len(p.connections)))

	var lastErr error
	for i, pooled := range p.connections {
		if err := pooled.conn.Close(); err != nil {
			logger.Error("Failed to close connection",
				zap.String("service", p.serviceName),
				zap.Int("index", i),
				zap.Error(err))
			lastErr = err
		}
	}

	p.connections = nil
	return lastErr
}

// Stats returns pool statistics
func (p *ConnectionPool) Stats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	healthy := 0
	states := make(map[string]int)

	for _, pooled := range p.connections {
		if pooled.healthy.Load() {
			healthy++
		}
		state := pooled.conn.GetState().String()
		states[state]++
	}

	return map[string]interface{}{
		"service": p.serviceName,
		"address": p.address,
		"total":   len(p.connections),
		"healthy": healthy,
		"states":  states,
	}
}
