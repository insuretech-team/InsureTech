package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ApiKeyScopeMiddleware returns middleware that checks the API key has the required scope.
// It calls authn ValidateToken to authenticate the key, then checks the scopes field.
//
// The token is expected as a Bearer token in the Authorization header.
// If the token starts with "eyJ" it is treated as a JWT and the middleware skips scope
// enforcement (passes through to next). Only non-JWT bearer tokens (API keys) are checked.
func ApiKeyScopeMiddleware(authnConn *grpc.ClientConn, requiredScope string) func(http.Handler) http.Handler {
	if authnConn == nil {
		logger.Warn("ApiKeyScopeMiddleware: authn connection is nil, allowing all requests (graceful degradation)")
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	client := authnservicev1.NewAuthServiceClient(authnConn)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r.Header.Get("Authorization"))

			// If token is empty or looks like a JWT (starts with "eyJ"), skip API key scope check.
			if token == "" || strings.HasPrefix(token, "eyJ") {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()

			resp, err := client.ValidateToken(ctx, &authnservicev1.ValidateTokenRequest{
				AccessToken: token,
			})
			if err != nil {
				logger.Warn("ApiKeyScopeMiddleware: ValidateToken error",
					zap.Error(err),
					zap.String("path", r.URL.Path),
				)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if resp == nil || !resp.Valid {
				logger.Warn("ApiKeyScopeMiddleware: invalid API key",
					zap.String("path", r.URL.Path),
				)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// NOTE: Permissions are no longer embedded in the token (RS256 + authz service).
			// API key scope enforcement is now delegated to the AuthZ service via
			// the AuthZMiddleware. Here we only validate the key is active and
			// propagate identity. Fine-grained scope checks happen downstream.
			// For backward compat, check portal matches required scope prefix.
			if requiredScope != "" && resp.Portal != "" {
				if !strings.HasPrefix(requiredScope, resp.Portal+":") && requiredScope != resp.Portal {
					logger.Warn("ApiKeyScopeMiddleware: portal/scope mismatch",
						zap.String("required", requiredScope),
						zap.String("portal", resp.Portal),
						zap.String("path", r.URL.Path),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusForbidden)
					_ = json.NewEncoder(w).Encode(map[string]string{
						"error":    "insufficient scope",
						"required": requiredScope,
					})
					return
				}
			}

			// Propagate user identity from the API key validation.
			ctx = context.WithValue(ctx, "user_id", resp.UserId)
			ctx = context.WithValue(ctx, "session_id", resp.SessionId)
			r.Header.Set("X-User-ID", resp.UserId)
			r.Header.Set("X-Session-ID", resp.SessionId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// containsScope is kept for future use when fine-grained scope lists are added.
// Currently scope enforcement is delegated to the AuthZ service.
func containsScope(scopes []string, scope string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}
