package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// CSRFMiddleware enforces CSRF token validation for SERVER_SIDE sessions.
// It only activates when X-Session-Type == "SERVER_SIDE" and the request
// method is mutating (not GET, HEAD, OPTIONS).
// The session ID is read from X-Session-ID (set by AuthMiddleware) and the
// CSRF token is read from the X-CSRF-Token request header.
func CSRFMiddleware(authnConn *grpc.ClientConn) func(http.Handler) http.Handler {
	var client authnservicev1.AuthServiceClient
	if authnConn != nil {
		client = authnservicev1.NewAuthServiceClient(authnConn)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip safe methods.
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				next.ServeHTTP(w, r)
				return
			}

			// Only enforce for SERVER_SIDE sessions.
			sessionType := r.Header.Get("X-Session-Type")
			if sessionType != "SERVER_SIDE" {
				next.ServeHTTP(w, r)
				return
			}

			csrfToken := r.Header.Get("X-CSRF-Token")
			if csrfToken == "" {
				logger.Warn("CSRF token missing for SERVER_SIDE session", zap.String("path", r.URL.Path))
				http.Error(w, "Forbidden: missing CSRF token", http.StatusForbidden)
				return
			}

			sessionID := r.Header.Get("X-Session-ID")

			if client == nil {
				logger.Error("CSRF middleware has no authn connection")
				http.Error(w, "Forbidden: CSRF validation unavailable", http.StatusForbidden)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			resp, err := client.ValidateCSRF(ctx, &authnservicev1.ValidateCSRFRequest{
				SessionId: sessionID,
				CsrfToken: csrfToken,
			})
			if err != nil {
				logger.Warn("ValidateCSRF gRPC error", zap.Error(err), zap.String("path", r.URL.Path))
				http.Error(w, "Forbidden: CSRF validation failed", http.StatusForbidden)
				return
			}
			if resp == nil || !resp.Valid {
				logger.Warn("Invalid CSRF token", zap.String("path", r.URL.Path))
				http.Error(w, "Forbidden: invalid CSRF token", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
