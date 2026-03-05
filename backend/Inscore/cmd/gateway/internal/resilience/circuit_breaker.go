package resilience

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int32

const (
	StateClosed   CircuitState = iota // Normal operation
	StateOpen                         // Circuit is open, rejecting requests
	StateHalfOpen                     // Testing if service recovered
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// CircuitBreaker implements the circuit breaker pattern with advanced features
type CircuitBreaker struct {
	name                   string
	maxFailures            int           // Max failures before opening
	maxConsecutiveFailures int           // Consecutive failures to open immediately
	resetTimeout           time.Duration // Time to wait before trying again
	halfOpenRequests       int           // Requests to allow in half-open state

	state               atomic.Int32 // CircuitState
	failures            atomic.Int64
	consecutiveFailures atomic.Int64
	successes           atomic.Int64
	lastFailureTime     atomic.Int64 // Unix timestamp
	halfOpenSuccess     atomic.Int32
	halfOpenFailure     atomic.Int32

	// Metrics
	totalRequests    atomic.Int64
	rejectedRequests atomic.Int64

	onStateChange func(from, to CircuitState)
}

// CircuitBreakerConfig configures circuit breaker behavior
type CircuitBreakerConfig struct {
	Name                   string
	MaxFailures            int
	MaxConsecutiveFailures int
	ResetTimeout           time.Duration
	HalfOpenRequests       int
	OnStateChange          func(from, to CircuitState)
}

// DefaultCircuitBreakerConfig returns production settings
func DefaultCircuitBreakerConfig(name string) *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		Name:                   name,
		MaxFailures:            5,
		MaxConsecutiveFailures: 3,
		ResetTimeout:           30 * time.Second,
		HalfOpenRequests:       3,
	}
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(cfg *CircuitBreakerConfig) *CircuitBreaker {
	if cfg == nil {
		cfg = DefaultCircuitBreakerConfig("default")
	}

	cb := &CircuitBreaker{
		name:                   cfg.Name,
		maxFailures:            cfg.MaxFailures,
		maxConsecutiveFailures: cfg.MaxConsecutiveFailures,
		resetTimeout:           cfg.ResetTimeout,
		halfOpenRequests:       cfg.HalfOpenRequests,
		onStateChange:          cfg.OnStateChange,
	}

	cb.state.Store(int32(StateClosed))

	return cb
}

// Execute runs the function if circuit allows it
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	cb.totalRequests.Add(1)

	// Check if allowed to execute
	if err := cb.beforeRequest(); err != nil {
		cb.rejectedRequests.Add(1)
		return err
	}

	// Execute function with timeout from context
	err := fn(ctx)

	// Record result
	cb.afterRequest(err)

	return err
}

// beforeRequest checks if request should be allowed
func (cb *CircuitBreaker) beforeRequest() error {
	state := CircuitState(cb.state.Load())

	switch state {
	case StateClosed:
		return nil

	case StateOpen:
		// Check if timeout has passed
		lastFailure := time.Unix(0, cb.lastFailureTime.Load())
		if time.Since(lastFailure) > cb.resetTimeout {
			// Transition to half-open
			if cb.state.CompareAndSwap(int32(StateOpen), int32(StateHalfOpen)) {
				cb.halfOpenSuccess.Store(0)
				cb.halfOpenFailure.Store(0)
				cb.notifyStateChange(StateOpen, StateHalfOpen)
				logger.Info("Circuit breaker entering half-open state",
					zap.String("name", cb.name))
				return nil
			}
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		// Allow limited requests in half-open state
		current := cb.halfOpenSuccess.Load() + cb.halfOpenFailure.Load()
		if current >= int32(cb.halfOpenRequests) {
			return ErrTooManyRequests
		}
		return nil

	default:
		return nil
	}
}

// afterRequest records the result of a request
func (cb *CircuitBreaker) afterRequest(err error) {
	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}
}

// onFailure handles a failed request
func (cb *CircuitBreaker) onFailure() {
	cb.failures.Add(1)
	cb.consecutiveFailures.Add(1)
	cb.lastFailureTime.Store(time.Now().UnixNano())

	state := CircuitState(cb.state.Load())

	switch state {
	case StateClosed:
		consecutive := cb.consecutiveFailures.Load()
		total := cb.failures.Load()

		// Open immediately if too many consecutive failures
		if consecutive >= int64(cb.maxConsecutiveFailures) {
			if cb.state.CompareAndSwap(int32(StateClosed), int32(StateOpen)) {
				cb.notifyStateChange(StateClosed, StateOpen)
				logger.Warn("Circuit breaker opened (consecutive failures)",
					zap.String("name", cb.name),
					zap.Int64("consecutive_failures", consecutive))
			}
		} else if total >= int64(cb.maxFailures) {
			if cb.state.CompareAndSwap(int32(StateClosed), int32(StateOpen)) {
				cb.notifyStateChange(StateClosed, StateOpen)
				logger.Warn("Circuit breaker opened",
					zap.String("name", cb.name),
					zap.Int64("failures", total))
			}
		}

	case StateHalfOpen:
		cb.halfOpenFailure.Add(1)
		if cb.state.CompareAndSwap(int32(StateHalfOpen), int32(StateOpen)) {
			cb.notifyStateChange(StateHalfOpen, StateOpen)
			logger.Warn("Circuit breaker reopened (half-open test failed)",
				zap.String("name", cb.name))
		}
	}
}

// onSuccess handles a successful request
func (cb *CircuitBreaker) onSuccess() {
	cb.successes.Add(1)
	cb.consecutiveFailures.Store(0) // Reset consecutive failures

	state := CircuitState(cb.state.Load())

	switch state {
	case StateClosed:
		// Decay failure count on success
		current := cb.failures.Load()
		if current > 0 {
			cb.failures.Store(current - 1)
		}

	case StateHalfOpen:
		success := cb.halfOpenSuccess.Add(1)
		if success >= int32(cb.halfOpenRequests) {
			if cb.state.CompareAndSwap(int32(StateHalfOpen), int32(StateClosed)) {
				cb.failures.Store(0)
				cb.notifyStateChange(StateHalfOpen, StateClosed)
				logger.Info("Circuit breaker closed (recovery successful)",
					zap.String("name", cb.name))
			}
		}
	}
}

// State returns the current state
func (cb *CircuitBreaker) State() CircuitState {
	return CircuitState(cb.state.Load())
}

// IsOpen returns true if circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.State() == StateOpen
}

// Name returns the circuit breaker name
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// Reset manually resets the circuit breaker
func (cb *CircuitBreaker) Reset() {
	oldState := CircuitState(cb.state.Swap(int32(StateClosed)))
	cb.failures.Store(0)
	cb.consecutiveFailures.Store(0)
	cb.halfOpenSuccess.Store(0)
	cb.halfOpenFailure.Store(0)

	if oldState != StateClosed {
		cb.notifyStateChange(oldState, StateClosed)
		logger.Info("Circuit breaker manually reset",
			zap.String("name", cb.name))
	}
}

// Stats returns circuit breaker statistics
func (cb *CircuitBreaker) Stats() map[string]interface{} {
	return map[string]interface{}{
		"name":                 cb.name,
		"state":                cb.State().String(),
		"failures":             cb.failures.Load(),
		"consecutive_failures": cb.consecutiveFailures.Load(),
		"successes":            cb.successes.Load(),
		"total_requests":       cb.totalRequests.Load(),
		"rejected_requests":    cb.rejectedRequests.Load(),
		"error_rate":           cb.calculateErrorRate(),
	}
}

// calculateErrorRate returns current error rate
func (cb *CircuitBreaker) calculateErrorRate() float64 {
	total := cb.totalRequests.Load()
	if total == 0 {
		return 0.0
	}
	failures := cb.failures.Load()
	return float64(failures) / float64(total)
}

// notifyStateChange calls the state change callback
func (cb *CircuitBreaker) notifyStateChange(from, to CircuitState) {
	if cb.onStateChange != nil {
		go cb.onStateChange(from, to)
	}
}

// CircuitBreakerManager manages multiple circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager creates a new manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate gets or creates a circuit breaker for a service
func (m *CircuitBreakerManager) GetOrCreate(serviceName string, cfg *CircuitBreakerConfig) *CircuitBreaker {
	m.mu.RLock()
	cb, exists := m.breakers[serviceName]
	m.mu.RUnlock()

	if exists {
		return cb
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if cb, exists := m.breakers[serviceName]; exists {
		return cb
	}

	// Create new circuit breaker
	if cfg == nil {
		cfg = DefaultCircuitBreakerConfig(serviceName)
	}
	cb = NewCircuitBreaker(cfg)
	m.breakers[serviceName] = cb

	logger.Info("Circuit breaker created",
		zap.String("service", serviceName))

	return cb
}

// Get retrieves a circuit breaker
func (m *CircuitBreakerManager) Get(serviceName string) (*CircuitBreaker, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cb, exists := m.breakers[serviceName]
	return cb, exists
}

// AllStats returns stats for all circuit breakers
func (m *CircuitBreakerManager) AllStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]interface{})
	for name, cb := range m.breakers {
		stats[name] = cb.Stats()
	}
	return stats
}
