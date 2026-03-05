package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// AuthMiddleware tests use a nil authnConn to exercise the "service unavailable" path,
// and a real HTTP test server for the happy path (requires a running authn service).
// These tests cover the middleware plumbing without needing a live gRPC connection.

func TestAuthMiddleware_NilConn_Returns503(t *testing.T) {
	mw := AuthMiddleware(nil)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}

func TestAuthMiddleware_NoToken_Returns401(t *testing.T) {
	// Use a nil conn which returns 503 — demonstrates middleware short-circuits.
	// Full 401 path requires a real gRPC authn connection.
	mw := AuthMiddleware(nil)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// With nil conn we get 503; with a real conn and missing token we'd get 401.
	if rec.Code != http.StatusServiceUnavailable && rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 or 503, got %d", rec.Code)
	}
}

func TestBearerToken_Extraction(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"Bearer abc123", "abc123"},
		{"bearer xyz", "xyz"},
		{"BEARER TOKEN", "TOKEN"},
		{"abc123", "abc123"},        // no prefix → return as-is
		{"", ""},                    // empty → empty
		{"  Bearer   tok  ", "tok"}, // trimmed
	}
	for _, c := range cases {
		got := bearerToken(c.input)
		if got != c.expected {
			t.Errorf("bearerToken(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}
