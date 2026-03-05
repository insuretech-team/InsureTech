package resilience

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/pool"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ResilientServiceClient provides resilient gRPC calls with pooling, retries, and circuit breaking
type ResilientServiceClient struct {
	serviceName    string
	pool           *pool.ConnectionPool
	circuitBreaker *CircuitBreaker
	retryPolicy    *RetryPolicy

	// Timeouts
	defaultTimeout time.Duration
}

// ServiceClientConfig configures a resilient service client
type ServiceClientConfig struct {
	ServiceName    string
	Address        string
	PoolConfig     *pool.PoolConfig
	CircuitConfig  *CircuitBreakerConfig
	RetryPolicy    *RetryPolicy
	DefaultTimeout time.Duration
}

// DefaultServiceClientConfig returns production-ready configuration
func DefaultServiceClientConfig(serviceName, address string) *ServiceClientConfig {
	return &ServiceClientConfig{
		ServiceName:    serviceName,
		Address:        address,
		PoolConfig:     pool.DefaultPoolConfig(),
		CircuitConfig:  DefaultCircuitBreakerConfig(serviceName),
		RetryPolicy:    DefaultRetryPolicy(),
		DefaultTimeout: 5 * time.Second,
	}
}

// NewResilientServiceClient creates a resilient service client
func NewResilientServiceClient(cfg *ServiceClientConfig) (*ResilientServiceClient, error) {
	if cfg == nil {
		logger.Error("Config cannot be nil")
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Create connection pool
	connPool := pool.NewConnectionPool(cfg.ServiceName, cfg.Address, cfg.PoolConfig)

	// Create circuit breaker
	cb := NewCircuitBreaker(cfg.CircuitConfig)

	client := &ResilientServiceClient{
		serviceName:    cfg.ServiceName,
		pool:           connPool,
		circuitBreaker: cb,
		retryPolicy:    cfg.RetryPolicy,
		defaultTimeout: cfg.DefaultTimeout,
	}

	return client, nil
}

// Initialize initializes the client (creates connections)
func (c *ResilientServiceClient) Initialize(ctx context.Context) error {
	logger.Info("Initializing resilient service client",
		zap.String("service", c.serviceName))

	return c.pool.Initialize(ctx)
}

// Invoke makes a resilient gRPC call with retry and circuit breaking
func (c *ResilientServiceClient) Invoke(
	ctx context.Context,
	method string,
	args interface{},
	reply interface{},
	opts ...grpc.CallOption,
) error {
	// Add timeout if not present
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	// Execute through circuit breaker
	return c.circuitBreaker.Execute(ctx, func(ctx context.Context) error {
		// Retry logic
		return Retry(ctx, c.retryPolicy, func(ctx context.Context) error {
			// Get connection from pool
			conn, err := c.pool.GetConnection()
			if err != nil {
				logger.Error("Failed to get connection", zap.Error(err))
				return fmt.Errorf("failed to get connection: %w", err)
			}

			// Invoke the actual gRPC call
			return conn.Invoke(ctx, method, args, reply, opts...)
		})
	})
}

// InvokeWithCustomPolicy makes a call with custom retry policy
func (c *ResilientServiceClient) InvokeWithCustomPolicy(
	ctx context.Context,
	method string,
	args interface{},
	reply interface{},
	customRetryPolicy *RetryPolicy,
	opts ...grpc.CallOption,
) error {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	return c.circuitBreaker.Execute(ctx, func(ctx context.Context) error {
		return Retry(ctx, customRetryPolicy, func(ctx context.Context) error {
			conn, err := c.pool.GetConnection()
			if err != nil {
				logger.Error("Failed to get connection", zap.Error(err))
				return fmt.Errorf("failed to get connection: %w", err)
			}
			return conn.Invoke(ctx, method, args, reply, opts...)
		})
	})
}

// InvokeNoRetry makes a call without retry (for idempotent checks)
func (c *ResilientServiceClient) InvokeNoRetry(
	ctx context.Context,
	method string,
	args interface{},
	reply interface{},
	opts ...grpc.CallOption,
) error {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	return c.circuitBreaker.Execute(ctx, func(ctx context.Context) error {
		conn, err := c.pool.GetConnection()
		if err != nil {
			logger.Error("Failed to get connection", zap.Error(err))
			return fmt.Errorf("failed to get connection: %w", err)
		}
		return conn.Invoke(ctx, method, args, reply, opts...)
	})
}

// GetConnection returns a connection (for streaming or custom operations)
func (c *ResilientServiceClient) GetConnection() (*grpc.ClientConn, error) {
	if c.circuitBreaker.IsOpen() {
		return nil, ErrCircuitOpen
	}
	return c.pool.GetConnection()
}

// Health checks if the client is healthy
func (c *ResilientServiceClient) Health() bool {
	// Check circuit breaker state
	if c.circuitBreaker.IsOpen() {
		return false
	}

	// Check pool health
	healthyCount := c.pool.HealthyCount()
	totalCount := c.pool.TotalCount()

	// Healthy if at least 50% connections are healthy
	return healthyCount > 0 && float64(healthyCount)/float64(totalCount) >= 0.5
}

// Stats returns comprehensive client statistics
func (c *ResilientServiceClient) Stats() map[string]interface{} {
	return map[string]interface{}{
		"service":         c.serviceName,
		"pool":            c.pool.Stats(),
		"circuit_breaker": c.circuitBreaker.Stats(),
		"healthy":         c.Health(),
	}
}

// Close closes the client
func (c *ResilientServiceClient) Close() error {
	logger.Info("Closing resilient service client",
		zap.String("service", c.serviceName))

	return c.pool.Close()
}

// ResilientClientManager manages multiple resilient clients
type ResilientClientManager struct {
	clients map[string]*ResilientServiceClient
	mu      sync.RWMutex
}

// NewResilientClientManager creates a manager
func NewResilientClientManager() *ResilientClientManager {
	return &ResilientClientManager{
		clients: make(map[string]*ResilientServiceClient),
	}
}

// RegisterClient registers a new client
func (m *ResilientClientManager) RegisterClient(cfg *ServiceClientConfig) (*ResilientServiceClient, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[cfg.ServiceName]; exists {
		logger.Error("Client already registered", zap.String("service", cfg.ServiceName))
		return nil, fmt.Errorf("client %s already registered", cfg.ServiceName)
	}

	client, err := NewResilientServiceClient(cfg)
	if err != nil {
		logger.Error("Failed to create client", zap.Error(err))
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	m.clients[cfg.ServiceName] = client
	return client, nil
}

// GetClient retrieves a client
func (m *ResilientClientManager) GetClient(serviceName string) (*ResilientServiceClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[serviceName]
	if !exists {
		logger.Error("Client not found", zap.String("service", serviceName))
		return nil, fmt.Errorf("client %s not found", serviceName)
	}
	return client, nil
}

// InitializeAll initializes all registered clients concurrently
func (m *ResilientClientManager) InitializeAll(ctx context.Context) error {
	m.mu.RLock()
	clients := make([]*ResilientServiceClient, 0, len(m.clients))
	for _, client := range m.clients {
		clients = append(clients, client)
	}
	m.mu.RUnlock()

	// Initialize concurrently
	errCh := make(chan error, len(clients))
	var wg sync.WaitGroup

	for _, client := range clients {
		wg.Add(1)
		go func(c *ResilientServiceClient) {
			defer wg.Done()
			if err := c.Initialize(ctx); err != nil {
				errCh <- fmt.Errorf("failed to initialize %s: %w", c.serviceName, err)
			}
		}(client)
	}

	wg.Wait()
	close(errCh)

	// Collect errors
	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		logger.Warn("Some clients failed to initialize",
			zap.Int("failed", len(errors)),
			zap.Int("total", len(clients)))
	}

	return nil
}

// HealthCheck returns health status of all clients
func (m *ResilientClientManager) HealthCheck() map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	health := make(map[string]bool)
	for name, client := range m.clients {
		health[name] = client.Health()
	}
	return health
}

// AllStats returns stats for all clients
func (m *ResilientClientManager) AllStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]interface{})
	for name, client := range m.clients {
		stats[name] = client.Stats()
	}
	return stats
}

// CloseAll closes all clients
func (m *ResilientClientManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			logger.Error("Failed to close client",
				zap.String("service", name),
				zap.Error(err))
			lastErr = err
		}
	}

	m.clients = make(map[string]*ResilientServiceClient)
	return lastErr
}
