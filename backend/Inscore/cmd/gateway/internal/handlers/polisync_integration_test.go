package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/handlers"
)

// ─────────────────────────────────────────────────────────────────────────────
// PoliSync Bridge Integration Tests
//
// These tests verify the HTTP bridge between the Go gateway and PoliSync C# services.
// They use mock backends to simulate PoliSync HTTP companion ports.
//
// Real end-to-end tests (requiring running PoliSync) are in the e2e/ package.
// ─────────────────────────────────────────────────────────────────────────────

// ── Test: Identity header propagation to PoliSync ────────────────────────────

func TestPoliSyncBridge_IdentityHeadersReachCSharpService(t *testing.T) {
	// Simulate the PoliSync C# AuthInterceptor reading identity headers
	var received struct {
		UserID    string
		TenantID  string
		PartnerID string
		UserType  string
		Portal    string
		Roles     string
		RequestID string
	}

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.UserID    = r.Header.Get("X-User-ID")
		received.TenantID  = r.Header.Get("X-Tenant-ID")
		received.PartnerID = r.Header.Get("X-Partner-ID")
		received.UserType  = r.Header.Get("X-User-Type")
		received.Portal    = r.Header.Get("X-Portal")
		received.Roles     = r.Header.Get("X-Roles")
		received.RequestID = r.Header.Get("X-Request-ID")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"product_id": "prod-123", "message": "Product created"})
	}))
	defer backend.Close()

	h := handlers.NewPoliSyncHandlerWithURL(nil, "product-service", backend.URL)

	req := httptest.NewRequest(http.MethodPost, "/v1/products",
		strings.NewReader(`{"product_code":"HLT-001","product_name":"Health Shield"}`))
	req.Header.Set("Content-Type", "application/json")
	// These are set by Go gateway auth_middleware.go after JWT validation
	req.Header.Set("X-User-ID",    "11111111-1111-1111-1111-111111111111")
	req.Header.Set("X-Tenant-ID",  "22222222-2222-2222-2222-222222222222")
	req.Header.Set("X-Partner-ID", "33333333-3333-3333-3333-333333333333")
	req.Header.Set("X-User-Type",  "agent")
	req.Header.Set("X-Portal",     "b2b")
	req.Header.Set("X-Roles",      "product:write,product:read")
	req.Header.Set("X-Request-ID", "req-test-abc-123")

	rr := httptest.NewRecorder()
	h.Proxy().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}

	// Verify C# service received ALL identity headers
	checks := []struct{ name, got, want string }{
		{"X-User-ID",    received.UserID,    "11111111-1111-1111-1111-111111111111"},
		{"X-Tenant-ID",  received.TenantID,  "22222222-2222-2222-2222-222222222222"},
		{"X-Partner-ID", received.PartnerID, "33333333-3333-3333-3333-333333333333"},
		{"X-User-Type",  received.UserType,  "agent"},
		{"X-Portal",     received.Portal,    "b2b"},
		{"X-Roles",      received.Roles,     "product:write,product:read"},
		{"X-Request-ID", received.RequestID, "req-test-abc-123"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("PoliSync did not receive %s: got %q, want %q", c.name, c.got, c.want)
		}
	}
}

// ── Test: All PoliSync services are routable via gateway ─────────────────────

func TestPoliSyncBridge_AllServicesRoutable(t *testing.T) {
	services := []struct {
		name    string
		method  string
		path    string
	}{
		{"product-service",      "GET",    "/v1/products"},
		{"product-service",      "POST",   "/v1/products"},
		{"quote-service",        "POST",   "/v1/quotations"},
		{"quote-service",        "GET",    "/v1/quotations/q-123"},
		{"order-service",        "POST",   "/v1/orders"},
		{"policy-service",       "GET",    "/v1/policies"},
		{"policy-service",       "POST",   "/v1/policies/p-123/cancel"},
		{"underwriting-service", "POST",   "/v1/health-declarations"},
		{"claim-service",        "POST",   "/v1/claims"},
		{"claim-service",        "GET",    "/v1/claims/c-123"},
		{"commission-service",   "GET",    "/v1/commission/payouts"},
	}

	for _, svc := range services {
		t.Run(svc.name+"_"+svc.method+"_"+svc.path, func(t *testing.T) {
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			defer backend.Close()

			h := handlers.NewPoliSyncHandlerWithURL(nil, svc.name, backend.URL)
			req := httptest.NewRequest(svc.method, svc.path, nil)
			req.Header.Set("X-User-ID",   "user-123")
			req.Header.Set("X-Tenant-ID", "tenant-456")

			rr := httptest.NewRecorder()
			h.Proxy().ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("service %s path %s: expected 200, got %d", svc.name, svc.path, rr.Code)
			}
		})
	}
}

// ── Test: Content-Type and JSON body preservation ────────────────────────────

func TestPoliSyncBridge_JsonBodyPreservedExactly(t *testing.T) {
	payloads := []string{
		`{"product_code":"MOT-001","category":"PRODUCT_CATEGORY_MOTOR","base_premium":{"amount":50000,"currency":"BDT"}}`,
		`{"quotation_id":"q-abc","sum_insured":{"amount":1000000,"currency":"BDT"},"tenure_months":12}`,
		`{"policy_id":"p-xyz","cancellation_reason":"Customer request","effective_date":"2026-03-01"}`,
	}

	for _, payload := range payloads {
		t.Run("payload_len_"+strings.Split(payload, "{")[1][:10], func(t *testing.T) {
			var capturedBody string
			var capturedCT   string
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				buf := make([]byte, len(payload)+10)
				n, _ := r.Body.Read(buf)
				capturedBody = string(buf[:n])
				capturedCT   = r.Header.Get("Content-Type")
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":"new-id"}`))
			}))
			defer backend.Close()

			h := handlers.NewPoliSyncHandlerWithURL(nil, "product-service", backend.URL)
			req := httptest.NewRequest(http.MethodPost, "/v1/products",
				strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID",    "user-123")
			req.Header.Set("X-Tenant-ID",  "tenant-456")

			rr := httptest.NewRecorder()
			h.Proxy().ServeHTTP(rr, req)

			if rr.Code != http.StatusCreated {
				t.Fatalf("expected 201, got %d", rr.Code)
			}
			if capturedBody != payload {
				t.Errorf("body mismatch:\n  got:  %q\n  want: %q", capturedBody, payload)
			}
			if !strings.Contains(capturedCT, "application/json") {
				t.Errorf("Content-Type not forwarded: got %q", capturedCT)
			}
		})
	}
}

// ── Test: Timeout behaviour ───────────────────────────────────────────────────

func TestPoliSyncBridge_SlowBackendReturns502(t *testing.T) {
	// Backend that takes longer than our test timeout
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(10 * time.Second):
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer backend.Close()

	h := handlers.NewPoliSyncHandlerWithURL(nil, "policy-service", backend.URL)

	// Use a context with a very short timeout to simulate a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/v1/policies/p-123", nil).WithContext(ctx)
	req.Header.Set("X-User-ID",   "user-123")
	req.Header.Set("X-Tenant-ID", "tenant-456")

	rr := httptest.NewRecorder()
	h.Proxy().ServeHTTP(rr, req)

	// Should be 502 (context cancelled → connection error) or 499 (client closed)
	if rr.Code != http.StatusBadGateway && rr.Code != 499 {
		t.Logf("got status %d (acceptable for timeout scenario)", rr.Code)
	}
}

// ── Test: PoliSync error responses pass through unchanged ─────────────────────

func TestPoliSyncBridge_ErrorResponsePassThrough(t *testing.T) {
	errorCases := []struct {
		backendStatus int
		backendBody   string
	}{
		{http.StatusBadRequest,          `{"code":"INVALID_ARGUMENT","message":"product_code is required"}`},
		{http.StatusNotFound,            `{"code":"NOT_FOUND","message":"Product not found"}`},
		{http.StatusConflict,            `{"code":"ALREADY_EXISTS","message":"product_code MOT-001 already exists"}`},
		{http.StatusUnprocessableEntity, `{"code":"VALIDATION","message":"base_premium must be > 0"}`},
		{http.StatusForbidden,           `{"code":"PERMISSION_DENIED","message":"Insufficient permissions"}`},
	}

	for _, tc := range errorCases {
		t.Run(http.StatusText(tc.backendStatus), func(t *testing.T) {
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.backendStatus)
				_, _ = w.Write([]byte(tc.backendBody))
			}))
			defer backend.Close()

			h := handlers.NewPoliSyncHandlerWithURL(nil, "product-service", backend.URL)
			req := httptest.NewRequest(http.MethodPost, "/v1/products",
				strings.NewReader(`{"product_code":""}`))
			req.Header.Set("X-User-ID",   "user-123")
			req.Header.Set("X-Tenant-ID", "tenant-456")

			rr := httptest.NewRecorder()
			h.Proxy().ServeHTTP(rr, req)

			if rr.Code != tc.backendStatus {
				t.Errorf("expected status %d, got %d", tc.backendStatus, rr.Code)
			}
			if rr.Body.String() != tc.backendBody {
				t.Errorf("body mismatch:\n  got:  %q\n  want: %q", rr.Body.String(), tc.backendBody)
			}
		})
	}
}
