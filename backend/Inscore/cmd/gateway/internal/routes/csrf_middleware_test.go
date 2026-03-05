package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// csrfNext is a simple 200 OK handler that marks itself called.
func csrfNext(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

// TestCSRF_SkipsGET verifies that GET requests bypass CSRF enforcement entirely.
func TestCSRF_SkipsGET(t *testing.T) {
	called := false
	// nil conn: if CSRF were enforced it would 403 (no client).
	h := CSRFMiddleware(nil)(csrfNext(&called))

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/logout", nil)
	req.Header.Set("X-Session-Type", "SERVER_SIDE")
	// deliberately no X-CSRF-Token

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.True(t, called, "next handler should have been called for GET")
	require.Equal(t, http.StatusOK, w.Code)
}

// TestCSRF_SkipsNonServerSide verifies that non-SERVER_SIDE sessions pass through
// without CSRF checking even for mutating methods.
func TestCSRF_SkipsNonServerSide(t *testing.T) {
	called := false
	h := CSRFMiddleware(nil)(csrfNext(&called))

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", nil)
	req.Header.Set("X-Session-Type", "JWT") // not SERVER_SIDE
	// deliberately no X-CSRF-Token

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.True(t, called, "next handler should be called for non-SERVER_SIDE session")
	require.Equal(t, http.StatusOK, w.Code)
}

// TestCSRF_RejectsMissingToken verifies that a POST with SERVER_SIDE session
// and no X-CSRF-Token header returns 403 before attempting any gRPC call.
func TestCSRF_RejectsMissingToken(t *testing.T) {
	called := false
	// nil conn: we expect 403 due to missing token, never reaching gRPC.
	h := CSRFMiddleware(nil)(csrfNext(&called))

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", nil)
	req.Header.Set("X-Session-Type", "SERVER_SIDE")
	req.Header.Set("X-Session-ID", "sess-123")
	// deliberately no X-CSRF-Token

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.False(t, called, "next handler must NOT be called when CSRF token is missing")
	require.Equal(t, http.StatusForbidden, w.Code)
}

// TestCSRF_SkipsHEAD verifies HEAD is treated as a safe method.
func TestCSRF_SkipsHEAD(t *testing.T) {
	called := false
	h := CSRFMiddleware(nil)(csrfNext(&called))

	req := httptest.NewRequest(http.MethodHead, "/v1/auth/logout", nil)
	req.Header.Set("X-Session-Type", "SERVER_SIDE")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

// TestCSRF_SkipsOPTIONS verifies OPTIONS is treated as a safe method.
func TestCSRF_SkipsOPTIONS(t *testing.T) {
	called := false
	h := CSRFMiddleware(nil)(csrfNext(&called))

	req := httptest.NewRequest(http.MethodOptions, "/v1/auth/logout", nil)
	req.Header.Set("X-Session-Type", "SERVER_SIDE")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}
