package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/handlers"
)

// mockPoliSyncBackend starts a fake PoliSync C# HTTP server that captures
// incoming requests so we can assert identity header forwarding.
func mockPoliSyncBackend(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the headers that PoliSync's AuthInterceptor would read
		resp := map[string]string{
			"x-user-id":    r.Header.Get("X-User-ID"),
			"x-tenant-id":  r.Header.Get("X-Tenant-ID"),
			"x-partner-id": r.Header.Get("X-Partner-ID"),
			"x-user-type":  r.Header.Get("X-User-Type"),
			"x-portal":     r.Header.Get("X-Portal"),
			"x-roles":      r.Header.Get("X-Roles"),
			"x-request-id": r.Header.Get("X-Request-ID"),
			"path":         r.URL.Path,
			"method":       r.Method,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

// TestPoliSyncHandler_ProxyForwardsIdentityHeaders verifies that the gateway's
// PoliSyncHandler correctly forwards all X-* identity headers that the Go auth
// middleware injects after JWT validation.
func TestPoliSyncHandler_ProxyForwardsIdentityHeaders(t *testing.T) {
	backend := mockPoliSyncBackend(t)
	defer backend.Close()

	h := handlers.NewPoliSyncHandlerWithURL(nil, "product-service", backend.URL)

	// Simulate request as if the Go gateway auth middleware has already
	// validated JWT and injected X-* headers
	req := httptest.NewRequest(http.MethodGet, "/v1/products", nil)
	req.Header.Set("X-User-ID", "550e8400-e29b-41d4-a716-446655440000")
	req.Header.Set("X-Tenant-ID", "660e8400-e29b-41d4-a716-446655440001")
	req.Header.Set("X-Partner-ID", "770e8400-e29b-41d4-a716-446655440002")
	req.Header.Set("X-User-Type", "agent")
	req.Header.Set("X-Portal", "b2b")
	req.Header.Set("X-Roles", "policy:read,policy:write")
	req.Header.Set("X-Request-ID", "req-abc-123")

	rr := httptest.NewRecorder()
	h.Proxy().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	tests := []struct {
		key  string
		want string
	}{
		{"x-user-id", "550e8400-e29b-41d4-a716-446655440000"},
		{"x-tenant-id", "660e8400-e29b-41d4-a716-446655440001"},
		{"x-partner-id", "770e8400-e29b-41d4-a716-446655440002"},
		{"x-user-type", "agent"},
		{"x-portal", "b2b"},
		{"x-roles", "policy:read,policy:write"},
		{"x-request-id", "req-abc-123"},
		{"path", "/v1/products"},
		{"method", "GET"},
	}

	for _, tt := range tests {
		t.Run("header_"+tt.key, func(t *testing.T) {
			if got := body[tt.key]; got != tt.want {
				t.Errorf("header %q: got %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

// TestPoliSyncHandler_ProxyForwardsRequestBody verifies POST body is forwarded intact.
func TestPoliSyncHandler_ProxyForwardsRequestBody(t *testing.T) {
	var capturedBody string
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		capturedBody = string(b)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"new-product-id"}`))
	}))
	defer backend.Close()

	h := handlers.NewPoliSyncHandlerWithURL(nil, "product-service", backend.URL)

	payload := `{"name":"Health Shield","category":"HEALTH","base_premium_paisa":5000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/products",
		strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "user-123")
	req.Header.Set("X-Tenant-ID", "tenant-456")

	rr := httptest.NewRecorder()
	h.Proxy().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	if capturedBody != payload {
		t.Errorf("body mismatch: got %q, want %q", capturedBody, payload)
	}
}

// TestPoliSyncHandler_ProxyReturns502WhenBackendDown verifies graceful 502 error
// when the PoliSync service is unreachable.
func TestPoliSyncHandler_ProxyReturns502WhenBackendDown(t *testing.T) {
	// Point to a port with nothing listening
	h := handlers.NewPoliSyncHandlerWithURL(nil, "policy-service", "http://127.0.0.1:19999")

	req := httptest.NewRequest(http.MethodGet, "/v1/policies/some-id", nil)
	req.Header.Set("X-User-ID", "user-123")
	req.Header.Set("X-Tenant-ID", "tenant-456")

	rr := httptest.NewRecorder()
	h.Proxy().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadGateway {
		t.Fatalf("expected 502 Bad Gateway, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("error response is not valid JSON: %v", err)
	}
	if body["code"] != "UNAVAILABLE" {
		t.Errorf("expected code=UNAVAILABLE, got %v", body["code"])
	}
}

// TestPoliSyncHandler_ProxyPreservesPath verifies the URL path is forwarded unchanged.
func TestPoliSyncHandler_ProxyPreservesPath(t *testing.T) {
	var capturedPath string
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer backend.Close()

	paths := []string{
		"/v1/policies/abc-123",
		"/v1/policies/abc-123/nominees",
		"/v1/claims/clm-456/documents",
		"/v1/commission/payouts/p-789",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			h := handlers.NewPoliSyncHandlerWithURL(nil, "policy-service", backend.URL)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			h.Proxy().ServeHTTP(rr, req)
			if capturedPath != path {
				t.Errorf("path mismatch: got %q, want %q", capturedPath, path)
			}
		})
	}
}

// TestPoliSyncHandler_ProxyReturns502WhenServiceNotConfigured verifies that
// using the default constructor with an unknown service name returns 502.
func TestPoliSyncHandler_ProxyReturns502WhenServiceNotConfigured(t *testing.T) {
	h := handlers.NewPoliSyncHandler(nil, "nonexistent-service")

	req := httptest.NewRequest(http.MethodGet, "/v1/anything", nil)
	rr := httptest.NewRecorder()
	h.Proxy().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rr.Code)
	}
	var body map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &body)
	if body["code"] != "UNAVAILABLE" {
		t.Errorf("expected UNAVAILABLE, got %v", body["code"])
	}
}
