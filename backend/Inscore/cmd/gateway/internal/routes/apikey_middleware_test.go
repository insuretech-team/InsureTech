package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// ApiKeyScopeMiddleware tests
// ---------------------------------------------------------------------------

// TestApiKeyScopeMiddleware_JWTPassthrough verifies that tokens beginning with
// "eyJ" (JWTs) are passed through without any scope check, regardless of the
// authnConn value.
func TestApiKeyScopeMiddleware_JWTPassthrough(t *testing.T) {
	jwtTokens := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1MSJ9.sig",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.e30.sig",
	}

	for _, tok := range jwtTokens {
		t.Run(tok[:10], func(t *testing.T) {
			called := false
			// Use nil conn — if JWT check works, conn is never used.
			mw := ApiKeyScopeMiddleware(nil, "read:policies")
			handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			}))

			r := httptest.NewRequest(http.MethodGet, "/v1/policies", nil)
			r.Header.Set("Authorization", "Bearer "+tok)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			require.True(t, called, "next handler should be called for JWT tokens")
			require.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestApiKeyScopeMiddleware_NilConnAllowsAll verifies that when authnConn is
// nil the middleware degrades gracefully and allows all requests (returns next).
func TestApiKeyScopeMiddleware_NilConnAllowsAll(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		{"jwt_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.sig"},
		{"api_key", "sk_live_abcdef123456"},
		{"empty", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			mw := ApiKeyScopeMiddleware(nil, "write:claims")
			handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			}))

			r := httptest.NewRequest(http.MethodPost, "/v1/claims", nil)
			if tc.token != "" {
				r.Header.Set("Authorization", "Bearer "+tc.token)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			require.True(t, called, "nil conn should allow all — next should always be called")
			require.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestContainsScope unit-tests the helper directly.
func TestContainsScope(t *testing.T) {
	perms := []string{"read:policies", "write:claims", "admin:users"}

	require.True(t, containsScope(perms, "read:policies"))
	require.True(t, containsScope(perms, "write:claims"))
	require.False(t, containsScope(perms, "delete:all"))
	require.False(t, containsScope(nil, "read:policies"))
	require.False(t, containsScope([]string{}, "read:policies"))
}
