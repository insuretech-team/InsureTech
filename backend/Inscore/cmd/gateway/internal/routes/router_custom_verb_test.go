package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCustomVerbCompatMiddleware_RewritesLegacyCustomVerbPaths(t *testing.T) {
	var gotPath string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	h := customVerbCompatMiddleware(next)
	req := httptest.NewRequest(http.MethodPost, "/v1/partners/partner-123:verify", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if gotPath != "/v1/partners/partner-123/verify" {
		t.Fatalf("expected rewritten path %q, got %q", "/v1/partners/partner-123/verify", gotPath)
	}
}

func TestCustomVerbCompatMiddleware_DoesNotRewriteStaticColonPath(t *testing.T) {
	var gotPath string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	h := customVerbCompatMiddleware(next)
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/otp:verify", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if gotPath != "/v1/auth/otp:verify" {
		t.Fatalf("expected path to remain unchanged, got %q", gotPath)
	}
}
