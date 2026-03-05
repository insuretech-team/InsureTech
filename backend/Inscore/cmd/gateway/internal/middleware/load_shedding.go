package middleware

import (
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

// LoadShedder protects the gateway from overload by limiting concurrent requests
type LoadShedder struct {
	maxConcurrent int64
	current       atomic.Int64
	rejected      atomic.Int64
	total         atomic.Int64
}

// NewLoadShedder creates a new load shedder
func NewLoadShedder(maxConcurrent int) *LoadShedder {
	return &LoadShedder{
		maxConcurrent: int64(maxConcurrent),
	}
}

// Middleware implements load shedding
func (ls *LoadShedder) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment total requests
		ls.total.Add(1)

		// Try to acquire slot
		current := ls.current.Add(1)
		defer ls.current.Add(-1)

		// Check if we're over capacity
		if current > ls.maxConcurrent {
			ls.rejected.Add(1)

			// Set rate limit headers
			w.Header().Set("Retry-After", "1")
			w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(ls.maxConcurrent, 10))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-Load-Shed", "true")

			http.Error(w, "Service temporarily overloaded, please retry", http.StatusServiceUnavailable)
			return
		}

		// Set current load headers
		w.Header().Set("X-Concurrent-Requests", strconv.FormatInt(current-1, 10))
		w.Header().Set("X-Max-Concurrent", strconv.FormatInt(ls.maxConcurrent, 10))

		next.ServeHTTP(w, r)
	})
}

// Stats returns load shedding statistics
func (ls *LoadShedder) Stats() map[string]interface{} {
	current := ls.current.Load()
	rejected := ls.rejected.Load()
	total := ls.total.Load()

	var rejectionRate float64
	if total > 0 {
		rejectionRate = float64(rejected) / float64(total) * 100
	}

	return map[string]interface{}{
		"max_concurrent": ls.maxConcurrent,
		"current":        current,
		"total_requests": total,
		"rejected":       rejected,
		"rejection_rate": rejectionRate,
		"utilization":    float64(current) / float64(ls.maxConcurrent) * 100,
	}
}

// GetCurrent returns current concurrent requests
func (ls *LoadShedder) GetCurrent() int64 {
	return ls.current.Load()
}

// GetRejected returns total rejected requests
func (ls *LoadShedder) GetRejected() int64 {
	return ls.rejected.Load()
}

// SetMaxConcurrent updates the max concurrent limit (hot reconfiguration)
func (ls *LoadShedder) SetMaxConcurrent(max int) {
	atomic.StoreInt64(&ls.maxConcurrent, int64(max))
}

// AdaptiveLoadShedder dynamically adjusts limits based on response times
type AdaptiveLoadShedder struct {
	*LoadShedder
	targetLatency time.Duration
	avgLatency    atomic.Int64 // nanoseconds
	sampleCount   atomic.Int64
}

// NewAdaptiveLoadShedder creates an adaptive load shedder
func NewAdaptiveLoadShedder(initialMax int, targetLatency time.Duration) *AdaptiveLoadShedder {
	return &AdaptiveLoadShedder{
		LoadShedder:   NewLoadShedder(initialMax),
		targetLatency: targetLatency,
	}
}

// Middleware implements adaptive load shedding
func (als *AdaptiveLoadShedder) Middleware(next http.Handler) http.Handler {
	// Start adaptation loop
	go als.adaptLoop()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Use base load shedder
		als.LoadShedder.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			// Record latency
			latency := time.Since(start)
			als.recordLatency(latency)
		})).ServeHTTP(w, r)
	})
}

// recordLatency records request latency for adaptation
func (als *AdaptiveLoadShedder) recordLatency(latency time.Duration) {
	// Simple exponential moving average
	current := time.Duration(als.avgLatency.Load())
	alpha := 0.1 // Smoothing factor

	newAvg := time.Duration(float64(current)*(1-alpha) + float64(latency)*alpha)
	als.avgLatency.Store(int64(newAvg))
	als.sampleCount.Add(1)
}

// adaptLoop periodically adjusts max concurrent based on latency
func (als *AdaptiveLoadShedder) adaptLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Need at least 10 samples to adapt
		if als.sampleCount.Load() < 10 {
			continue
		}

		avgLatency := time.Duration(als.avgLatency.Load())
		currentMax := atomic.LoadInt64(&als.maxConcurrent)

		// If latency is too high, reduce capacity
		if avgLatency > als.targetLatency*2 {
			newMax := int64(float64(currentMax) * 0.9) // Reduce by 10%
			if newMax < 10 {
				newMax = 10 // Minimum 10
			}
			atomic.StoreInt64(&als.maxConcurrent, newMax)
		}

		// If latency is good and we're near capacity, increase
		if avgLatency < als.targetLatency && float64(als.current.Load()) > float64(currentMax)*0.8 {
			newMax := int64(float64(currentMax) * 1.1) // Increase by 10%
			if newMax > 1000 {
				newMax = 1000 // Max 1000
			}
			atomic.StoreInt64(&als.maxConcurrent, newMax)
		}

		// Reset sample count
		als.sampleCount.Store(0)
	}
}

// GetAvgLatency returns average latency
func (als *AdaptiveLoadShedder) GetAvgLatency() time.Duration {
	return time.Duration(als.avgLatency.Load())
}
