// Package metrics provides Prometheus metrics for the AuthN microservice.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// API Key Rotation Metrics
	
	// APIKeyRotationsTotal counts API key rotations
	APIKeyRotationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authn",
			Name:      "api_key_rotations_total",
			Help:      "Total number of API key rotations.",
		},
		[]string{"owner_type", "status"}, // status: "success" | "failure"
	)
	
	// APIKeyRotationDuration measures rotation operation duration
	APIKeyRotationDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "authn",
			Name:      "api_key_rotation_duration_seconds",
			Help:      "Duration of API key rotation operations in seconds.",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
	)
	
	// APIKeysActive tracks the number of active API keys by status
	APIKeysActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "authn",
			Name:      "api_keys_active",
			Help:      "Number of API keys by status.",
		},
		[]string{"status", "owner_type"}, // status: "active" | "rotating" | "revoked" | "expired"
	)
	
	// Portal Configuration Metrics
	
	// PortalConfigCacheHits tracks portal config cache hits/misses
	PortalConfigCacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authn",
			Name:      "portal_config_cache_hits_total",
			Help:      "Total portal config cache hits and misses.",
		},
		[]string{"portal", "hit"}, // hit: "true" | "false"
	)
	
	// PortalConfigLoadDuration measures config load time from AuthZ
	PortalConfigLoadDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "authn",
			Name:      "portal_config_load_duration_seconds",
			Help:      "Duration to load portal config from AuthZ in seconds.",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"portal"},
	)
	
	// PasswordValidationFailures tracks password validation failures by reason
	PasswordValidationFailures = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authn",
			Name:      "password_validation_failures_total",
			Help:      "Total password validation failures by portal and reason.",
		},
		[]string{"portal", "reason"}, // reason: "too_short" | "no_uppercase" | "no_lowercase" | "no_digit" | "no_symbol"
	)
	
	// Session Metrics
	
	// SessionsCreated counts sessions created by portal and session type
	SessionsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authn",
			Name:      "sessions_created_total",
			Help:      "Total sessions created by portal and type.",
		},
		[]string{"portal", "session_type"}, // session_type: "SERVER_SIDE" | "JWT" | "API_KEY"
	)
	
	// SessionValidationLatency measures session/token validation latency
	SessionValidationLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "authn",
			Name:      "session_validation_latency_seconds",
			Help:      "Session validation latency in seconds.",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25},
		},
		[]string{"session_type", "valid"}, // valid: "true" | "false"
	)
	
	// ActiveSessions tracks current active sessions by portal
	ActiveSessions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "authn",
			Name:      "active_sessions",
			Help:      "Number of active sessions by portal.",
		},
		[]string{"portal", "session_type"},
	)
	
	// MFA Metrics
	
	// MFAChallengesIssued counts MFA challenges issued
	MFAChallengesIssued = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authn",
			Name:      "mfa_challenges_issued_total",
			Help:      "Total MFA challenges issued by portal and method.",
		},
		[]string{"portal", "method"}, // method: "TOTP" | "SMS" | "EMAIL"
	)
	
	// MFAVerificationResults counts MFA verification results
	MFAVerificationResults = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "authn",
			Name:      "mfa_verification_results_total",
			Help:      "Total MFA verification results by portal and outcome.",
		},
		[]string{"portal", "method", "result"}, // result: "success" | "failure"
	)
)

// RecordAPIKeyRotation records an API key rotation event
func RecordAPIKeyRotation(ownerType string, success bool, durationSeconds float64) {
	status := "success"
	if !success {
		status = "failure"
	}
	APIKeyRotationsTotal.WithLabelValues(ownerType, status).Inc()
	APIKeyRotationDuration.Observe(durationSeconds)
}

// UpdateAPIKeyCount updates the gauge for active API keys
func UpdateAPIKeyCount(status, ownerType string, count float64) {
	APIKeysActive.WithLabelValues(status, ownerType).Set(count)
}

// RecordPortalConfigCache records a portal config cache hit/miss
func RecordPortalConfigCache(portal string, hit bool) {
	hitStr := "false"
	if hit {
		hitStr = "true"
	}
	PortalConfigCacheHits.WithLabelValues(portal, hitStr).Inc()
}

// RecordPortalConfigLoad records portal config load duration
func RecordPortalConfigLoad(portal string, durationSeconds float64) {
	PortalConfigLoadDuration.WithLabelValues(portal).Observe(durationSeconds)
}

// RecordPasswordValidationFailure records a password validation failure
func RecordPasswordValidationFailure(portal, reason string) {
	PasswordValidationFailures.WithLabelValues(portal, reason).Inc()
}

// RecordSessionCreation records a session creation event
func RecordSessionCreation(portal, sessionType string) {
	SessionsCreated.WithLabelValues(portal, sessionType).Inc()
}

// RecordSessionValidation records session validation latency
func RecordSessionValidation(sessionType string, valid bool, durationSeconds float64) {
	validStr := "false"
	if valid {
		validStr = "true"
	}
	SessionValidationLatency.WithLabelValues(sessionType, validStr).Observe(durationSeconds)
}

// UpdateActiveSessions updates the active sessions gauge
func UpdateActiveSessions(portal, sessionType string, count float64) {
	ActiveSessions.WithLabelValues(portal, sessionType).Set(count)
}

// RecordMFAChallenge records an MFA challenge issuance
func RecordMFAChallenge(portal, method string) {
	MFAChallengesIssued.WithLabelValues(portal, method).Inc()
}

// RecordMFAVerification records an MFA verification result
func RecordMFAVerification(portal, method string, success bool) {
	result := "success"
	if !success {
		result = "failure"
	}
	MFAVerificationResults.WithLabelValues(portal, method, result).Inc()
}
