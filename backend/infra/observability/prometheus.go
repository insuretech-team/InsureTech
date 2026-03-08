package observability

import (
	"bufio"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMetrics holds all Prometheus metrics
type PrometheusMetrics struct {
	// HTTP metrics
	RequestsTotal     *prometheus.CounterVec
	RequestDuration   *prometheus.HistogramVec
	RequestsInFlight  prometheus.Gauge
	ResponseSize      *prometheus.HistogramVec
	
	// Circuit breaker metrics
	CircuitBreakerState  *prometheus.GaugeVec
	CircuitBreakerTrips  *prometheus.CounterVec
	
	// Connection pool metrics
	PoolConnectionsTotal   *prometheus.GaugeVec
	PoolConnectionsHealthy *prometheus.GaugeVec
	PoolConnectionsActive  *prometheus.GaugeVec
	
	// Retry metrics
	RetriesTotal   *prometheus.CounterVec
	RetrySuccesses *prometheus.CounterVec
	
	// Rate limit metrics
	RateLimitHits *prometheus.CounterVec
	
	// System metrics
	GoroutineCount   prometheus.Gauge
	MemoryAlloc      prometheus.Gauge
	MemoryTotal      prometheus.Gauge
}

// NewPrometheusMetrics initializes all Prometheus metrics
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		// HTTP metrics
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),
		
		RequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gateway_http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),
		
		ResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100B to 100MB
			},
			[]string{"method", "path"},
		),
		
		// Circuit breaker metrics
		CircuitBreakerState: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gateway_circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
			},
			[]string{"service"},
		),
		
		CircuitBreakerTrips: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_circuit_breaker_trips_total",
				Help: "Total number of circuit breaker trips",
			},
			[]string{"service"},
		),
		
		// Connection pool metrics
		PoolConnectionsTotal: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gateway_pool_connections_total",
				Help: "Total connections in pool",
			},
			[]string{"service"},
		),
		
		PoolConnectionsHealthy: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gateway_pool_connections_healthy",
				Help: "Healthy connections in pool",
			},
			[]string{"service"},
		),
		
		PoolConnectionsActive: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gateway_pool_connections_active",
				Help: "Active connections currently in use",
			},
			[]string{"service"},
		),
		
		// Retry metrics
		RetriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_retries_total",
				Help: "Total number of retry attempts",
			},
			[]string{"service", "method"},
		),
		
		RetrySuccesses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_retry_successes_total",
				Help: "Total number of successful retries",
			},
			[]string{"service", "method"},
		),
		
		// Rate limit metrics
		RateLimitHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_rate_limit_hits_total",
				Help: "Total number of rate limit hits",
			},
			[]string{"ip"},
		),
		
		// System metrics
		GoroutineCount: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gateway_goroutines",
				Help: "Current number of goroutines",
			},
		),
		
		MemoryAlloc: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gateway_memory_alloc_bytes",
				Help: "Currently allocated memory in bytes",
			},
		),
		
		MemoryTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gateway_memory_total_bytes",
				Help: "Total allocated memory in bytes",
			},
		),
	}
}

// MetricsMiddleware instruments HTTP handlers with Prometheus metrics
func (pm *PrometheusMetrics) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Track in-flight requests
		pm.RequestsInFlight.Inc()
		defer pm.RequestsInFlight.Dec()
		
		// Wrap response writer to capture status and size
		wrapped := &metricsRecorder{
			ResponseWriter: w,
			statusCode:     200,
			bytesWritten:   0,
		}
		
		// Process request
		next.ServeHTTP(wrapped, r)
		
		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapped.statusCode)
		
		pm.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		pm.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		pm.ResponseSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(wrapped.bytesWritten))
	})
}

type metricsRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (r *metricsRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *metricsRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += n
	return n, err
}

// Flush implements http.Flusher to support SSE streaming
func (r *metricsRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack implements http.Hijacker to support WebSockets
func (r *metricsRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// UpdateCircuitBreakerMetrics updates circuit breaker metrics
func (pm *PrometheusMetrics) UpdateCircuitBreakerMetrics(service string, state int, tripped bool) {
	pm.CircuitBreakerState.WithLabelValues(service).Set(float64(state))
	if tripped {
		pm.CircuitBreakerTrips.WithLabelValues(service).Inc()
	}
}

// UpdatePoolMetrics updates connection pool metrics
func (pm *PrometheusMetrics) UpdatePoolMetrics(service string, total, healthy, active int) {
	pm.PoolConnectionsTotal.WithLabelValues(service).Set(float64(total))
	pm.PoolConnectionsHealthy.WithLabelValues(service).Set(float64(healthy))
	pm.PoolConnectionsActive.WithLabelValues(service).Set(float64(active))
}

// RecordRetry records retry attempt
func (pm *PrometheusMetrics) RecordRetry(service, method string, success bool) {
	pm.RetriesTotal.WithLabelValues(service, method).Inc()
	if success {
		pm.RetrySuccesses.WithLabelValues(service, method).Inc()
	}
}

// RecordRateLimitHit records rate limit hit
func (pm *PrometheusMetrics) RecordRateLimitHit(ip string) {
	pm.RateLimitHits.WithLabelValues(ip).Inc()
}

// UpdateSystemMetrics updates system-level metrics
func (pm *PrometheusMetrics) UpdateSystemMetrics(goroutines int, memAlloc, memTotal uint64) {
	pm.GoroutineCount.Set(float64(goroutines))
	pm.MemoryAlloc.Set(float64(memAlloc))
	pm.MemoryTotal.Set(float64(memTotal))
}

// Handler returns Prometheus HTTP handler
func (pm *PrometheusMetrics) Handler() http.Handler {
	return promhttp.Handler()
}
