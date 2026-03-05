package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// AdaptiveRateLimiter adjusts rate limits based on backend health and load
type AdaptiveRateLimiter struct {
	// Base configuration
	baseRate  float64
	baseBurst int

	// Current limits (atomic for concurrent access)
	currentRate  atomic.Value // float64
	currentBurst atomic.Int32

	// Per-IP limiters
	limiters sync.Map // map[string]*rate.Limiter

	// Adaptation parameters
	backendHealthy   atomic.Bool
	errorRate        atomic.Value // float64
	adaptationWindow time.Duration

	// Metrics
	totalRequests   atomic.Int64
	limitedRequests atomic.Int64

	// Cleanup
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// AdaptiveRateLimiterConfig configures adaptive rate limiting
type AdaptiveRateLimiterConfig struct {
	BaseRate         float64       // Base requests per second
	BaseBurst        int           // Base burst size
	MinRate          float64       // Minimum rate (when degraded)
	MaxRate          float64       // Maximum rate (when healthy)
	AdaptationWindow time.Duration // How often to adapt
	CleanupInterval  time.Duration // Cleanup stale entries
}

// DefaultAdaptiveConfig returns production-ready config
func DefaultAdaptiveConfig() *AdaptiveRateLimiterConfig {
	return &AdaptiveRateLimiterConfig{
		BaseRate:         100.0, // 100 req/sec baseline
		BaseBurst:        200,
		MinRate:          10.0,  // Throttle to 10 req/sec when degraded
		MaxRate:          500.0, // Scale up to 500 req/sec when healthy
		AdaptationWindow: 10 * time.Second,
		CleanupInterval:  5 * time.Minute,
	}
}

// NewAdaptiveRateLimiter creates an adaptive rate limiter
func NewAdaptiveRateLimiter(cfg *AdaptiveRateLimiterConfig) *AdaptiveRateLimiter {
	if cfg == nil {
		cfg = DefaultAdaptiveConfig()
	}

	rl := &AdaptiveRateLimiter{
		baseRate:         cfg.BaseRate,
		baseBurst:        cfg.BaseBurst,
		adaptationWindow: cfg.AdaptationWindow,
		stopCh:           make(chan struct{}),
	}

	// Initialize current limits
	rl.currentRate.Store(cfg.BaseRate)
	rl.currentBurst.Store(int32(cfg.BaseBurst))
	rl.backendHealthy.Store(true)
	rl.errorRate.Store(0.0)

	// Start adaptation loop
	rl.startAdaptation(cfg)

	// Start cleanup
	rl.startCleanup(cfg.CleanupInterval)

	return rl
}

// Middleware returns HTTP middleware
func (rl *AdaptiveRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.totalRequests.Add(1)

		// Get client IP
		ip := getIP(r)

		// Get or create limiter for this IP
		limiter := rl.getLimiter(ip)

		// Check rate limit
		if !limiter.Allow() {
			rl.limitedRequests.Add(1)

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", limiter.Limit()))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("Retry-After", "1")

			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", limiter.Limit()))

		next.ServeHTTP(w, r)
	})
}

// getLimiter gets or creates a rate limiter for an IP
func (rl *AdaptiveRateLimiter) getLimiter(ip string) *rate.Limiter {
	if limiter, ok := rl.limiters.Load(ip); ok {
		return limiter.(*rate.Limiter)
	}

	// Create new limiter with current rate
	currentRate := rl.currentRate.Load().(float64)
	currentBurst := int(rl.currentBurst.Load())

	limiter := rate.NewLimiter(rate.Limit(currentRate), currentBurst)

	actual, loaded := rl.limiters.LoadOrStore(ip, limiter)
	if loaded {
		return actual.(*rate.Limiter)
	}

	return limiter
}

// UpdateBackendHealth updates backend health status
func (rl *AdaptiveRateLimiter) UpdateBackendHealth(healthy bool, errorRate float64) {
	rl.backendHealthy.Store(healthy)
	rl.errorRate.Store(errorRate)
}

// startAdaptation starts the adaptation loop
func (rl *AdaptiveRateLimiter) startAdaptation(cfg *AdaptiveRateLimiterConfig) {
	rl.wg.Add(1)
	go func() {
		defer rl.wg.Done()
		ticker := time.NewTicker(rl.adaptationWindow)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rl.adapt(cfg)
			case <-rl.stopCh:
				return
			}
		}
	}()
}

// adapt adjusts rate limits based on current conditions
func (rl *AdaptiveRateLimiter) adapt(cfg *AdaptiveRateLimiterConfig) {
	healthy := rl.backendHealthy.Load()
	errorRate := rl.errorRate.Load().(float64)

	var newRate float64
	var newBurst int

	if !healthy || errorRate > 0.2 {
		// Backend unhealthy or high error rate - throttle
		newRate = cfg.MinRate
		newBurst = int(cfg.MinRate * 2)

		logger.Warn("Rate limiter throttling due to backend issues",
			zap.Bool("healthy", healthy),
			zap.Float64("error_rate", errorRate),
			zap.Float64("new_rate", newRate))

	} else if errorRate < 0.01 {
		// Very healthy - increase limits
		newRate = cfg.MaxRate
		newBurst = int(cfg.MaxRate * 2)

		logger.Debug("Rate limiter increasing limits (healthy backend)",
			zap.Float64("error_rate", errorRate),
			zap.Float64("new_rate", newRate))

	} else {
		// Moderate - use base rate with slight adjustment
		adjustment := 1.0 - (errorRate * 2) // Reduce based on error rate
		newRate = cfg.BaseRate * adjustment
		newBurst = int(newRate * 2)
	}

	// Update current limits
	oldRate := rl.currentRate.Load().(float64)
	rl.currentRate.Store(newRate)
	rl.currentBurst.Store(int32(newBurst))

	if oldRate != newRate {
		logger.Info("Rate limits adapted",
			zap.Float64("old_rate", oldRate),
			zap.Float64("new_rate", newRate),
			zap.Int("new_burst", newBurst))
	}

	// Update existing limiters
	rl.updateExistingLimiters(newRate, newBurst)
}

// updateExistingLimiters updates all existing per-IP limiters
func (rl *AdaptiveRateLimiter) updateExistingLimiters(newRate float64, newBurst int) {
	rl.limiters.Range(func(key, value interface{}) bool {
		limiter := value.(*rate.Limiter)
		limiter.SetLimit(rate.Limit(newRate))
		limiter.SetBurst(newBurst)
		return true
	})
}

// startCleanup periodically removes stale limiters
func (rl *AdaptiveRateLimiter) startCleanup(interval time.Duration) {
	rl.wg.Add(1)
	go func() {
		defer rl.wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rl.cleanup()
			case <-rl.stopCh:
				return
			}
		}
	}()
}

// cleanup removes limiters that haven't been used recently
func (rl *AdaptiveRateLimiter) cleanup() {
	// Simple cleanup: remove all limiters (they'll be recreated on next request)
	// In production, you'd want to track last access time
	beforeCount := 0
	rl.limiters.Range(func(key, value interface{}) bool {
		beforeCount++
		return true
	})

	// Clear all
	rl.limiters = sync.Map{}

	if beforeCount > 0 {
		logger.Debug("Rate limiter cleanup complete",
			zap.Int("cleared", beforeCount))
	}
}

// Stats returns rate limiter statistics
func (rl *AdaptiveRateLimiter) Stats() map[string]interface{} {
	activeLimiters := 0
	rl.limiters.Range(func(key, value interface{}) bool {
		activeLimiters++
		return true
	})

	total := rl.totalRequests.Load()
	limited := rl.limitedRequests.Load()

	var limitRate float64
	if total > 0 {
		limitRate = float64(limited) / float64(total)
	}

	return map[string]interface{}{
		"current_rate":     rl.currentRate.Load(),
		"current_burst":    rl.currentBurst.Load(),
		"active_limiters":  activeLimiters,
		"total_requests":   total,
		"limited_requests": limited,
		"limit_rate":       limitRate,
		"backend_healthy":  rl.backendHealthy.Load(),
		"error_rate":       rl.errorRate.Load(),
	}
}

// Close stops the rate limiter
func (rl *AdaptiveRateLimiter) Close() {
	close(rl.stopCh)
	rl.wg.Wait()
}
