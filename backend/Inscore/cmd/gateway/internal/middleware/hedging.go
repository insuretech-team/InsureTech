package middleware

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// HedgingConfig configures request hedging behavior
type HedgingConfig struct {
	// Delay before sending hedged request
	HedgeDelay time.Duration

	// Maximum number of hedged requests (1 = original + 1 hedge)
	MaxHedges int

	// Only hedge these paths
	AllowedPaths []string
}

// DefaultHedgingConfig returns sensible defaults
func DefaultHedgingConfig() *HedgingConfig {
	return &HedgingConfig{
		HedgeDelay: 50 * time.Millisecond, // Hedge after 50ms
		MaxHedges:  1,                     // One hedge maximum
		AllowedPaths: []string{
			"/v1/products",       // Product reads
			"/v1/products/",      // Individual products
			"/v1/shipping/rates", // Shipping rates
		},
	}
}

// Hedger manages request hedging
type Hedger struct {
	config *HedgingConfig

	// Statistics
	totalRequests  atomic.Int64
	hedgedRequests atomic.Int64
	hedgeWins      atomic.Int64 // Hedged request returned first
	originalWins   atomic.Int64 // Original request returned first
}

// NewHedger creates a new request hedger
func NewHedger(config *HedgingConfig) *Hedger {
	if config == nil {
		config = DefaultHedgingConfig()
	}

	return &Hedger{
		config: config,
	}
}

// shouldHedge determines if a request should be hedged
func (h *Hedger) shouldHedge(r *http.Request) bool {
	// Only hedge GET requests (idempotent)
	if r.Method != http.MethodGet {
		return false
	}

	// Check if path is allowed for hedging
	for _, path := range h.config.AllowedPaths {
		if len(r.URL.Path) >= len(path) && r.URL.Path[:len(path)] == path {
			return true
		}
	}

	return false
}

// hedgedResponseWriter captures response for hedging
type hedgedResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	headersMu   sync.Mutex
	wroteHeader bool
	body        []byte
	bodyMu      sync.Mutex
}

func newHedgedResponseWriter(w http.ResponseWriter) *hedgedResponseWriter {
	return &hedgedResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *hedgedResponseWriter) WriteHeader(statusCode int) {
	w.headersMu.Lock()
	defer w.headersMu.Unlock()

	if !w.wroteHeader {
		w.statusCode = statusCode
		w.wroteHeader = true
	}
}

func (w *hedgedResponseWriter) Write(b []byte) (int, error) {
	w.bodyMu.Lock()
	defer w.bodyMu.Unlock()

	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	// Capture body
	w.body = append(w.body, b...)
	return len(b), nil
}

// Flush implements http.Flusher
func (w *hedgedResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// replay writes captured response to actual response writer
func (w *hedgedResponseWriter) replay(target http.ResponseWriter) {
	w.headersMu.Lock()
	defer w.headersMu.Unlock()
	w.bodyMu.Lock()
	defer w.bodyMu.Unlock()

	// Copy headers
	for key, values := range w.Header() {
		for _, value := range values {
			target.Header().Add(key, value)
		}
	}

	target.WriteHeader(w.statusCode)
	target.Write(w.body)
}

// responseResult holds a hedged response result
type responseResult struct {
	writer *hedgedResponseWriter
	err    error
	index  int // 0 = original, 1+ = hedges
}

// Middleware implements request hedging
func (h *Hedger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.totalRequests.Add(1)

		// Check if request should be hedged
		if !h.shouldHedge(r) {
			next.ServeHTTP(w, r)
			return
		}

		h.hedgedRequests.Add(1)

		// Create context for all requests
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		// Result channel
		resultCh := make(chan responseResult, h.config.MaxHedges+1)

		// Track which request finished first
		var once sync.Once

		// Execute original request
		go func() {
			hrw := newHedgedResponseWriter(w)
			req := r.Clone(ctx)

			next.ServeHTTP(hrw, req)

			once.Do(func() {
				resultCh <- responseResult{writer: hrw, index: 0}
			})
		}()

		// Start hedging timer
		hedgeTimer := time.NewTimer(h.config.HedgeDelay)
		defer hedgeTimer.Stop()

		// Wait for first response or hedge timeout
		select {
		case result := <-resultCh:
			// Original request won
			h.originalWins.Add(1)
			w.Header().Set("X-Hedged", "false")
			w.Header().Set("X-Hedge-Winner", "original")
			result.writer.replay(w)
			cancel() // Cancel any pending hedges
			return

		case <-hedgeTimer.C:
			// Send hedged request
			for i := 0; i < h.config.MaxHedges; i++ {
				go func(hedgeIndex int) {
					hrw := newHedgedResponseWriter(w)
					req := r.Clone(ctx)

					next.ServeHTTP(hrw, req)

					once.Do(func() {
						resultCh <- responseResult{writer: hrw, index: hedgeIndex + 1}
					})
				}(i)
			}

			// Wait for any response
			result := <-resultCh
			if result.index > 0 {
				h.hedgeWins.Add(1)
				w.Header().Set("X-Hedged", "true")
				w.Header().Set("X-Hedge-Winner", "hedged")
			} else {
				h.originalWins.Add(1)
				w.Header().Set("X-Hedged", "true")
				w.Header().Set("X-Hedge-Winner", "original")
			}

			result.writer.replay(w)
			cancel() // Cancel any pending requests
		}
	})
}

// Stats returns hedging statistics
func (h *Hedger) Stats() map[string]interface{} {
	total := h.totalRequests.Load()
	hedged := h.hedgedRequests.Load()
	hedgeWins := h.hedgeWins.Load()
	originalWins := h.originalWins.Load()

	var hedgeRate, hedgeWinRate float64
	if total > 0 {
		hedgeRate = float64(hedged) / float64(total) * 100
	}
	if hedged > 0 {
		hedgeWinRate = float64(hedgeWins) / float64(hedged) * 100
	}

	return map[string]interface{}{
		"total_requests":  total,
		"hedged_requests": hedged,
		"hedge_rate":      hedgeRate,
		"hedge_wins":      hedgeWins,
		"original_wins":   originalWins,
		"hedge_win_rate":  hedgeWinRate,
		"config": map[string]interface{}{
			"delay":      h.config.HedgeDelay.String(),
			"max_hedges": h.config.MaxHedges,
		},
	}
}

// SetHedgeDelay updates the hedge delay (hot reconfiguration)
func (h *Hedger) SetHedgeDelay(delay time.Duration) {
	h.config.HedgeDelay = delay
}

// GetHedgeWinRate returns percentage of hedged requests that won
func (h *Hedger) GetHedgeWinRate() float64 {
	hedged := h.hedgedRequests.Load()
	if hedged == 0 {
		return 0
	}
	return float64(h.hedgeWins.Load()) / float64(hedged) * 100
}
