package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newCORSHandler() http.Handler {
	// corsMiddleware reads CORS_ALLOWED_ORIGINS env; unset → defaults to localhost:3000,localhost:5173
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return corsMiddleware(inner)
}

func TestCORS_AllowedOrigin_SetsHeaders(t *testing.T) {
	h := newCORSHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	require.Equal(t, "Origin", w.Header().Get("Vary"))
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCORS_AllowedOrigin_Vite_SetsHeaders(t *testing.T) {
	h := newCORSHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_DisallowedOrigin_NoHeaders(t *testing.T) {
	h := newCORSHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	require.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
	// Request should still succeed — CORS only affects browser behaviour
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCORS_NoOriginHeader_NoHeaders(t *testing.T) {
	h := newCORSHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	// No Origin header at all (e.g. same-origin or non-browser request)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCORS_OptionsRequest_ReturnsNoContent(t *testing.T) {
	h := newCORSHandler()

	req := httptest.NewRequest(http.MethodOptions, "/v1/auth/login", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	require.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	require.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
}
