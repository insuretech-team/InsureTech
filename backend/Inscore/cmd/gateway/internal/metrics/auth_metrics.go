// Package metrics provides Prometheus metrics for the Gateway.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Combined Auth Middleware Metrics
	
	// CombinedAuthRequests counts combined auth requests
	CombinedAuthRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Name:      "combined_auth_requests_total",
			Help:      "Total combined auth requests by portal and outcome.",
		},
		[]string{"portal", "result"}, // result: "success" | "authn_failed" | "authz_failed" | "error"
	)
	
	// CombinedAuthLatency measures combined auth latency
	CombinedAuthLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "gateway",
			Name:      "combined_auth_latency_seconds",
			Help:      "Combined auth validation latency in seconds.",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"portal", "cached"}, // cached: "true" | "false"
	)
	
	// CombinedAuthCacheHits tracks cache hit/miss ratio
	CombinedAuthCacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Name:      "combined_auth_cache_hits_total",
			Help:      "Total combined auth cache hits and misses.",
		},
		[]string{"hit"}, // hit: "true" | "false"
	)
	
	// CombinedAuthCacheSize tracks cache size
	CombinedAuthCacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gateway",
			Name:      "combined_auth_cache_size",
			Help:      "Current size of the combined auth cache.",
		},
	)
	
	// Permission Preloader Metrics
	
	// PermissionPreloads counts permission preload requests
	PermissionPreloads = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Name:      "permission_preloads_total",
			Help:      "Total permission preload requests by portal and status.",
		},
		[]string{"portal", "status"}, // status: "success" | "error"
	)
	
	// PermissionPreloadLatency measures permission preload latency
	PermissionPreloadLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "gateway",
			Name:      "permission_preload_latency_seconds",
			Help:      "Permission preload latency in seconds.",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"portal"},
	)
	
	// PermissionPreloadCacheHits tracks preloader cache hits
	PermissionPreloadCacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Name:      "permission_preload_cache_hits_total",
			Help:      "Total permission preload cache hits and misses.",
		},
		[]string{"portal", "hit"}, // hit: "true" | "false"
	)
	
	// PermissionsPreloaded tracks the number of permissions loaded per user
	PermissionsPreloaded = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "gateway",
			Name:      "permissions_preloaded_count",
			Help:      "Number of permissions preloaded per user.",
			Buckets:   []float64{1, 5, 10, 25, 50, 100, 250, 500},
		},
		[]string{"portal"},
	)
	
	// Circuit Breaker Metrics
	
	// CircuitBreakerState tracks circuit breaker state
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "gateway",
			Name:      "circuit_breaker_state",
			Help:      "Circuit breaker state (0=closed, 1=open, 2=half-open).",
		},
		[]string{"service"}, // service: "authn" | "authz"
	)
	
	// CircuitBreakerTrips counts circuit breaker trips
	CircuitBreakerTrips = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Name:      "circuit_breaker_trips_total",
			Help:      "Total circuit breaker trips.",
		},
		[]string{"service"},
	)
	
	// CircuitBreakerRequests counts requests through circuit breaker
	CircuitBreakerRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Name:      "circuit_breaker_requests_total",
			Help:      "Total requests through circuit breaker by service and result.",
		},
		[]string{"service", "result"}, // result: "success" | "failure" | "rejected"
	)
	
	// Device Binding Metrics
	
	// DeviceBindingChecks counts device binding validation checks
	DeviceBindingChecks = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Name:      "device_binding_checks_total",
			Help:      "Total device binding validation checks.",
		},
		[]string{"result"}, // result: "matched" | "mismatched" | "not_required"
	)
)

// RecordCombinedAuthRequest records a combined auth request
func RecordCombinedAuthRequest(portal, result string, latencySeconds float64, cached bool) {
	CombinedAuthRequests.WithLabelValues(portal, result).Inc()
	
	cachedStr := "false"
	if cached {
		cachedStr = "true"
	}
	CombinedAuthLatency.WithLabelValues(portal, cachedStr).Observe(latencySeconds)
}

// RecordCombinedAuthCacheHit records a cache hit/miss
func RecordCombinedAuthCacheHit(hit bool) {
	hitStr := "false"
	if hit {
		hitStr = "true"
	}
	CombinedAuthCacheHits.WithLabelValues(hitStr).Inc()
}

// UpdateCombinedAuthCacheSize updates the cache size gauge
func UpdateCombinedAuthCacheSize(size float64) {
	CombinedAuthCacheSize.Set(size)
}

// RecordPermissionPreload records a permission preload request
func RecordPermissionPreload(portal, status string, latencySeconds float64, permissionCount int) {
	PermissionPreloads.WithLabelValues(portal, status).Inc()
	PermissionPreloadLatency.WithLabelValues(portal).Observe(latencySeconds)
	if status == "success" {
		PermissionsPreloaded.WithLabelValues(portal).Observe(float64(permissionCount))
	}
}

// RecordPermissionPreloadCacheHit records a preloader cache hit/miss
func RecordPermissionPreloadCacheHit(portal string, hit bool) {
	hitStr := "false"
	if hit {
		hitStr = "true"
	}
	PermissionPreloadCacheHits.WithLabelValues(portal, hitStr).Inc()
}

// UpdateCircuitBreakerState updates the circuit breaker state gauge
// state: 0=closed, 1=open, 2=half-open
func UpdateCircuitBreakerState(service string, state float64) {
	CircuitBreakerState.WithLabelValues(service).Set(state)
}

// RecordCircuitBreakerTrip records a circuit breaker trip
func RecordCircuitBreakerTrip(service string) {
	CircuitBreakerTrips.WithLabelValues(service).Inc()
}

// RecordCircuitBreakerRequest records a request through the circuit breaker
func RecordCircuitBreakerRequest(service, result string) {
	CircuitBreakerRequests.WithLabelValues(service, result).Inc()
}

// RecordDeviceBindingCheck records a device binding validation
func RecordDeviceBindingCheck(result string) {
	DeviceBindingChecks.WithLabelValues(result).Inc()
}
