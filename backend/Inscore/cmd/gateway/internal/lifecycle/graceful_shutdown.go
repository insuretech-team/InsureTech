package lifecycle

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
)

// GracefulShutdown manages graceful shutdown with connection draining
type GracefulShutdown struct {
	server         *http.Server
	shutdownSignal atomic.Bool
	inFlightReqs   atomic.Int64

	// Configuration
	drainTimeout    time.Duration
	shutdownTimeout time.Duration
}

// NewGracefulShutdown creates a new graceful shutdown manager
func NewGracefulShutdown(server *http.Server) *GracefulShutdown {
	return &GracefulShutdown{
		server:          server,
		drainTimeout:    30 * time.Second, // Wait 30s for existing requests
		shutdownTimeout: 35 * time.Second, // Force shutdown after 35s
	}
}

// TrackRequest increments in-flight request counter
func (gs *GracefulShutdown) TrackRequest() func() {
	gs.inFlightReqs.Add(1)
	return func() {
		gs.inFlightReqs.Add(-1)
	}
}

// Middleware wraps handlers to track in-flight requests
func (gs *GracefulShutdown) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if shutting down
		if gs.IsShuttingDown() {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Server is shutting down"))
			return
		}

		// Track request
		done := gs.TrackRequest()
		defer done()

		// Process request
		next.ServeHTTP(w, r)
	})
}

// Shutdown performs graceful shutdown with connection draining
func (gs *GracefulShutdown) Shutdown(ctx context.Context) error {
	// Mark as shutting down
	gs.shutdownSignal.Store(true)

	logger.Info("Starting graceful shutdown",
		zap.Int64("in_flight_requests", gs.inFlightReqs.Load()))

	// Step 1: Disable keep-alives (stop accepting new requests on existing connections)
	gs.server.SetKeepAlivesEnabled(false)
	logger.Info("Keep-alives disabled")

	// Step 2: Wait for in-flight requests to complete
	drainCtx, drainCancel := context.WithTimeout(ctx, gs.drainTimeout)
	defer drainCancel()

	logger.Info("Draining connections...",
		zap.Duration("timeout", gs.drainTimeout))

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	drained := false
	for !drained {
		select {
		case <-drainCtx.Done():
			remaining := gs.inFlightReqs.Load()
			if remaining > 0 {
				logger.Warn("Drain timeout reached with requests still in flight",
					zap.Int64("remaining_requests", remaining))
			}
			drained = true

		case <-ticker.C:
			remaining := gs.inFlightReqs.Load()
			if remaining == 0 {
				logger.Info("All requests completed")
				drained = true
			} else {
				logger.Debug("Waiting for requests to complete",
					zap.Int64("remaining", remaining))
			}
		}
	}

	// Step 3: Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gs.shutdownTimeout)
	defer shutdownCancel()

	logger.Info("Shutting down HTTP server")
	if err := gs.server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
		return err
	}

	logger.Info("Graceful shutdown complete")
	return nil
}

// IsShuttingDown returns true if shutdown is in progress
func (gs *GracefulShutdown) IsShuttingDown() bool {
	return gs.shutdownSignal.Load()
}

// InFlightRequests returns current number of in-flight requests
func (gs *GracefulShutdown) InFlightRequests() int64 {
	return gs.inFlightReqs.Load()
}

// ReadinessProbe returns HTTP 200 when ready, 503 when draining
func (gs *GracefulShutdown) ReadinessProbe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if gs.IsShuttingDown() {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status": "draining"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ready"}`))
	}
}

// LivenessProbe returns HTTP 200 when process is alive
func (gs *GracefulShutdown) LivenessProbe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "alive"}`))
	}
}
