package middleware

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/respond"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// Flush implements http.Flusher interface for SSE streaming
func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

type timeoutResponseWriter struct {
	mu          sync.Mutex
	header      http.Header
	body        bytes.Buffer
	statusCode  int
	wroteHeader bool
	timedOut    bool
}

func newTimeoutResponseWriter() *timeoutResponseWriter {
	return &timeoutResponseWriter{
		header:     make(http.Header),
		statusCode: http.StatusOK,
	}
}

func (tw *timeoutResponseWriter) Header() http.Header {
	return tw.header
}

func (tw *timeoutResponseWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.statusCode = code
	tw.wroteHeader = true
}

func (tw *timeoutResponseWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return len(b), nil
	}
	if !tw.wroteHeader {
		tw.statusCode = http.StatusOK
		tw.wroteHeader = true
	}
	return tw.body.Write(b)
}

func (tw *timeoutResponseWriter) Flush() {}

func (tw *timeoutResponseWriter) Timeout() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.timedOut = true
	tw.header = make(http.Header)
	tw.body.Reset()
}

func (tw *timeoutResponseWriter) WriteTo(w http.ResponseWriter) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return
	}
	for key, values := range tw.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(tw.statusCode)
	_, _ = w.Write(tw.body.Bytes())
}

// Metrics middleware logs request metrics
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Get request ID from context
		requestID := GetRequestID(r.Context())

		// Process request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log metrics with structured logging
		logger.Info("HTTP request",
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", rw.statusCode),
			zap.Duration("duration", duration),
			zap.Int64("bytes_written", rw.written),
			zap.String("user_agent", r.UserAgent()),
		)

		// In production, you'd export these to Prometheus/StatsD:
		// - HTTP request count by method, path, status
		// - Request duration histogram
		// - Response size histogram
		// - Active requests gauge
	})
}

// Timeout middleware enforces request timeout
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create timeout context
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Channel to signal completion
			done := make(chan struct{})
			tw := newTimeoutResponseWriter()

			// Execute handler in goroutine
			go func() {
				next.ServeHTTP(tw, r.WithContext(ctx))
				close(done)
			}()

			// Wait for completion or timeout
			select {
			case <-done:
				tw.WriteTo(w)
				return
			case <-ctx.Done():
				tw.Timeout()
				respond.Error(w, r, http.StatusGatewayTimeout, "DEADLINE_EXCEEDED", "Request timeout")
			}
		})
	}
}
