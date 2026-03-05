// Package metrics provides Prometheus metrics for the AuthZ microservice.
// Exposes:
//   - authz_decisions_total{portal, decision} — counter
//   - authz_decision_latency_ms{portal} — histogram
//   - authz_cache_hit_ratio{level} — gauge
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// DecisionsTotal counts AuthZ decisions by portal and outcome (allow/deny).
	DecisionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authz",
			Name:      "decisions_total",
			Help:      "Total AuthZ policy decisions partitioned by portal and decision outcome.",
		},
		[]string{"portal", "decision"}, // decision: "allow" | "deny"
	)

	// DecisionLatencyMs measures CheckAccess latency in milliseconds.
	DecisionLatencyMs = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "authz",
			Name:      "decision_latency_ms",
			Help:      "Latency of AuthZ CheckAccess decisions in milliseconds.",
			Buckets:   []float64{0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500},
		},
		[]string{"portal"},
	)

	// CacheHitRatio tracks L1/L2 cache hit ratios.
	CacheHitRatio = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "authz",
			Name:      "cache_hit_ratio",
			Help:      "Cache hit ratio for AuthZ policy cache (L1=in-process, L2=Redis).",
		},
		[]string{"level"}, // level: "L1" | "L2"
	)
	
	// API Key Scope Validation Metrics
	
	// APIScopeValidations counts API key scope validations
	APIScopeValidations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authz",
			Name:      "api_scope_validations_total",
			Help:      "Total API key scope validations by outcome.",
		},
		[]string{"portal", "result"}, // result: "allowed" | "denied"
	)
	
	// APIScopeValidationLatency measures scope validation latency
	APIScopeValidationLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "authz",
			Name:      "api_scope_validation_latency_ms",
			Help:      "API key scope validation latency in milliseconds.",
			Buckets:   []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
	)
	
	// APIScopeDenialReasons tracks reasons for scope denials
	APIScopeDenialReasons = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authz",
			Name:      "api_scope_denial_reasons_total",
			Help:      "Total API key scope denials by reason.",
		},
		[]string{"reason"}, // reason: "no_scopes" | "scope_mismatch"
	)
	
	// Portal Configuration Metrics
	
	// PortalConfigRequests counts portal config requests
	PortalConfigRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authz",
			Name:      "portal_config_requests_total",
			Help:      "Total portal configuration requests.",
		},
		[]string{"portal", "status"}, // status: "success" | "not_found" | "error"
	)
)

// RecordDecision records a CheckAccess decision in Prometheus.
// portal is extracted from the domain string (e.g. "system:root" → "system").
func RecordDecision(domain string, allowed bool, latencyMs float64) {
	portal := extractPortal(domain)
	decision := "allow"
	if !allowed {
		decision = "deny"
	}
	DecisionsTotal.WithLabelValues(portal, decision).Inc()
	DecisionLatencyMs.WithLabelValues(portal).Observe(latencyMs)
}

// UpdateCacheHitRatio updates the cache hit ratio gauge for a given cache level.
func UpdateCacheHitRatio(level string, ratio float64) {
	CacheHitRatio.WithLabelValues(level).Set(ratio)
}

// RecordCacheHit records a cache hit or miss for metrics tracking.
func RecordCacheHit(hit bool) {
	// Simple counter-based approach - can be enhanced with sliding window
	// For now, just tracking in the decision metrics
	// Future: implement proper hit/miss counters and calculate ratio
}

// extractPortal returns the portal portion of a domain string.
// "system:root" → "system", "agent:tenant-abc" → "agent"
func extractPortal(domain string) string {
	for i, c := range domain {
		if c == ':' {
			return domain[:i]
		}
	}
	return domain
}

// RecordAPIScopeValidation records an API key scope validation
func RecordAPIScopeValidation(domain string, allowed bool, latencyMs float64) {
	portal := extractPortal(domain)
	result := "allowed"
	if !allowed {
		result = "denied"
	}
	APIScopeValidations.WithLabelValues(portal, result).Inc()
	APIScopeValidationLatency.Observe(latencyMs)
}

// RecordAPIScopeDenial records a scope denial with reason
func RecordAPIScopeDenial(reason string) {
	if reason == "API key has no scopes defined" {
		APIScopeDenialReasons.WithLabelValues("no_scopes").Inc()
	} else {
		APIScopeDenialReasons.WithLabelValues("scope_mismatch").Inc()
	}
}

// RecordPortalConfigRequest records a portal config request
func RecordPortalConfigRequest(portal, status string) {
	PortalConfigRequests.WithLabelValues(portal, status).Inc()
}
