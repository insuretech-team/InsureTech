package resilience

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryPolicy defines retry behavior
type RetryPolicy struct {
	MaxAttempts       int
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
	RetryableErrors   map[codes.Code]bool
	Jitter            bool
}

// DefaultRetryPolicy returns a sensible retry policy
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts:       3,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        2 * time.Second,
		BackoffMultiplier: 2.0,
		RetryableErrors: map[codes.Code]bool{
			codes.Unavailable:       true,
			codes.DeadlineExceeded:  true,
			codes.ResourceExhausted: true,
			codes.Aborted:           true,
			codes.Internal:          true, // Sometimes transient
		},
		Jitter: true,
	}
}

// AggressiveRetryPolicy for critical operations
func AggressiveRetryPolicy() *RetryPolicy {
	policy := DefaultRetryPolicy()
	policy.MaxAttempts = 5
	policy.InitialBackoff = 50 * time.Millisecond
	return policy
}

// ConservativeRetryPolicy for non-critical operations
func ConservativeRetryPolicy() *RetryPolicy {
	policy := DefaultRetryPolicy()
	policy.MaxAttempts = 2
	policy.InitialBackoff = 200 * time.Millisecond
	return policy
}

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// Retry executes a function with retry logic
func Retry(ctx context.Context, policy *RetryPolicy, operation RetryableFunc) error {
	if policy == nil {
		policy = DefaultRetryPolicy()
	}

	var lastErr error
	for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
		// Check context before attempting
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Execute operation
		err := operation(ctx)
		if err == nil {
			// Success
			if attempt > 0 {
				logger.Debug("Operation succeeded after retry",
					zap.Int("attempt", attempt+1))
			}
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !policy.isRetryable(err) {
			logger.Debug("Non-retryable error encountered",
				zap.Error(err))
			return err
		}

		// Last attempt - don't backoff
		if attempt == policy.MaxAttempts-1 {
			break
		}

		// Calculate backoff duration
		backoff := policy.calculateBackoff(attempt)

		logger.Debug("Retrying operation",
			zap.Int("attempt", attempt+1),
			zap.Int("max_attempts", policy.MaxAttempts),
			zap.Duration("backoff", backoff),
			zap.Error(err))

		// Wait with context cancellation support
		select {
		case <-time.After(backoff):
			// Continue to next attempt
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	logger.Error("Max retry attempts exceeded", zap.Int("attempts", policy.MaxAttempts), zap.Error(lastErr))
	return fmt.Errorf("max retry attempts (%d) exceeded: %w", policy.MaxAttempts, lastErr)
}

// isRetryable checks if an error should be retried
func (p *RetryPolicy) isRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Extract gRPC status code
	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error, don't retry by default
		return false
	}

	return p.RetryableErrors[st.Code()]
}

// calculateBackoff computes backoff duration with exponential growth and optional jitter
func (p *RetryPolicy) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: initialBackoff * (multiplier ^ attempt)
	backoff := float64(p.InitialBackoff) * math.Pow(p.BackoffMultiplier, float64(attempt))

	// Cap at max backoff
	if backoff > float64(p.MaxBackoff) {
		backoff = float64(p.MaxBackoff)
	}

	duration := time.Duration(backoff)

	// Add jitter to prevent thundering herd
	if p.Jitter {
		jitter := time.Duration(rand.Int63n(int64(duration / 4))) // +/- 25% jitter
		duration = duration - jitter/2 + jitter
	}

	return duration
}

// RetryWithMetrics wraps retry with metrics tracking
type RetryMetrics struct {
	TotalAttempts     int
	SuccessfulRetries int
	FailedRetries     int
	TotalBackoffTime  time.Duration
}

// RetryWithTracking executes retry with metrics
func RetryWithTracking(ctx context.Context, policy *RetryPolicy, operation RetryableFunc) (*RetryMetrics, error) {
	metrics := &RetryMetrics{}
	startTime := time.Now()

	err := Retry(ctx, policy, func(ctx context.Context) error {
		metrics.TotalAttempts++
		return operation(ctx)
	})

	metrics.TotalBackoffTime = time.Since(startTime)

	if err != nil {
		metrics.FailedRetries = metrics.TotalAttempts - 1
	} else if metrics.TotalAttempts > 1 {
		metrics.SuccessfulRetries = metrics.TotalAttempts - 1
	}

	return metrics, err
}

// AdaptiveRetryPolicy adjusts retry behavior based on observed failures
type AdaptiveRetryPolicy struct {
	policy *RetryPolicy
	mu     sync.RWMutex

	// Metrics for adaptation
	recentFailures   int
	recentSuccesses  int
	windowStart      time.Time
	adaptationWindow time.Duration
}

// NewAdaptiveRetryPolicy creates an adaptive retry policy
func NewAdaptiveRetryPolicy() *AdaptiveRetryPolicy {
	return &AdaptiveRetryPolicy{
		policy:           DefaultRetryPolicy(),
		windowStart:      time.Now(),
		adaptationWindow: 1 * time.Minute,
	}
}

// GetPolicy returns current retry policy (adapts based on recent behavior)
func (a *AdaptiveRetryPolicy) GetPolicy() *RetryPolicy {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Reset window if expired
	if time.Since(a.windowStart) > a.adaptationWindow {
		a.recentFailures = 0
		a.recentSuccesses = 0
		a.windowStart = time.Now()
	}

	// Adapt based on failure rate
	totalRequests := a.recentFailures + a.recentSuccesses
	if totalRequests > 100 {
		failureRate := float64(a.recentFailures) / float64(totalRequests)

		if failureRate > 0.3 {
			// High failure rate - be more aggressive
			a.policy.MaxAttempts = 5
			a.policy.InitialBackoff = 50 * time.Millisecond
		} else if failureRate < 0.05 {
			// Low failure rate - be conservative
			a.policy.MaxAttempts = 2
			a.policy.InitialBackoff = 200 * time.Millisecond
		}
	}

	return a.policy
}

// RecordResult records operation result for adaptation
func (a *AdaptiveRetryPolicy) RecordResult(success bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if success {
		a.recentSuccesses++
	} else {
		a.recentFailures++
	}
}
