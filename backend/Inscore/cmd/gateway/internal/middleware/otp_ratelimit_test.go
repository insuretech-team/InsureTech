package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func otpOKHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// TestOTPRateLimit_First5Pass verifies that the first maxRequests calls succeed.
func TestOTPRateLimit_First5Pass(t *testing.T) {
	rl := OTPRateLimit(5, 10*time.Minute)
	h := rl(otpOKHandler())

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/otp:send", nil)
		req.RemoteAddr = "192.0.2.1:1234"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code, "request %d should pass", i+1)
	}
}

// TestOTPRateLimit_6thRequestReturns429 verifies that the 6th call within the
// window is rejected with HTTP 429.
func TestOTPRateLimit_6thRequestReturns429(t *testing.T) {
	rl := OTPRateLimit(5, 10*time.Minute)
	h := rl(otpOKHandler())

	ip := "192.0.2.2:5678"

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/otp:send", nil)
		req.RemoteAddr = ip
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	}

	// 6th request must be rate-limited.
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/otp:send", nil)
	req.RemoteAddr = ip
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	require.Equal(t, http.StatusTooManyRequests, w.Code)
}

// TestOTPRateLimit_DifferentIPsAreIndependent verifies that two different IPs
// each get their own bucket.
func TestOTPRateLimit_DifferentIPsAreIndependent(t *testing.T) {
	rl := OTPRateLimit(5, 10*time.Minute)
	h := rl(otpOKHandler())

	// Exhaust the limit for IP A.
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/otp:send", nil)
		req.RemoteAddr = "10.0.0.1:1111"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	}

	// IP B should still be within its own limit.
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/otp:send", nil)
	req.RemoteAddr = "10.0.0.2:2222"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
