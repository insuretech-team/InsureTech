package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
)

// deviceBindSecret returns the DEVICE_BIND_SECRET env var as bytes.
// Falls back to a compile-time constant ONLY for local dev — MUST be set in production.
func deviceBindSecret() []byte {
	s := os.Getenv("DEVICE_BIND_SECRET")
	if s == "" {
		// Local dev fallback — NOT safe for production
		s = "insuretech-dev-device-bind-secret-000000"
	}
	return []byte(s)
}

// deriveDeviceCredential returns a deterministic HMAC-SHA256 credential bound to
// a specific mobile number + device ID pair. This is the core of WhatsApp-style
// device binding:
//   - Same mobile + same device_id → always same credential ✅
//   - Different device_id → completely different credential ❌ (device-bound)
//   - Different mobile_number → completely different credential ❌ (mobile-bound)
//
// The server re-derives on every login request to compare — no DB storage needed.
func deriveDeviceCredential(mobileNumber, deviceID string) string {
	mac := hmac.New(sha256.New, deviceBindSecret())
	mac.Write([]byte(mobileNumber + ":" + deviceID))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

// isMobileDeviceType returns true for mobile device types that support
// WhatsApp-style device credential binding.
// Only real user-owned devices (Android, iOS) get device credentials.
// "API" is intentionally excluded: API integrations must use explicit API keys
// (see CreateAPIKey/ListAPIKeys) — NOT the mobile device credential flow.
// This prevents an API caller from silently inheriting a user's device binding.
func isMobileDeviceType(deviceType string) bool {
	dt := strings.ToUpper(strings.TrimSpace(deviceType))
	return dt == "ANDROID" || dt == "IOS"
}

// deviceCredentialMatches compares the provided credential against the expected
// derived credential using constant-time comparison (immune to timing attacks).
func deviceCredentialMatches(mobileNumber, deviceID, provided string) bool {
	if mobileNumber == "" || deviceID == "" || provided == "" {
		return false
	}
	expected := deriveDeviceCredential(mobileNumber, deviceID)
	// constant-time compare
	return hmac.Equal([]byte(expected), []byte(provided))
}
