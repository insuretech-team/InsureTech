package routes

import (
	"context"
	"net/http"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Cookie name used by AuthN metadata extractor.
// For web portals we store the *server-side session token* here.
const SessionCookieName = "session_token"

// AuthMiddleware enforces authentication using AuthN.ValidateToken.
// Supports:
// - Authorization: Bearer <jwt>
// - Cookie: session_token=<server-side-session-token>
func AuthMiddleware(authnConn *grpc.ClientConn) func(http.Handler) http.Handler {
	if authnConn == nil {
		logger.Error("Auth middleware created without authn connection")
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Authentication service unavailable", http.StatusServiceUnavailable)
			})
		}
	}

	client := authnservicev1.NewAuthServiceClient(authnConn)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			jwt := bearerToken(r.Header.Get("Authorization"))
			sessionToken := ""
			if c, err := r.Cookie(SessionCookieName); err == nil {
				sessionToken = c.Value
			}

			ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
				"authorization":   []string{r.Header.Get("Authorization")},
				"cookie":          []string{r.Header.Get("Cookie")},
				"x-csrf-token":    []string{r.Header.Get("X-CSRF-Token")},
				"x-device-id":     []string{r.Header.Get("X-Device-Id")},
				"x-forwarded-for": []string{r.Header.Get("X-Forwarded-For")},
				"x-real-ip":       []string{r.Header.Get("X-Real-Ip")},
				"user-agent":      []string{r.UserAgent()},
			})

			resp, err := client.ValidateToken(ctx, &authnservicev1.ValidateTokenRequest{
				AccessToken: jwt,
				SessionId:   sessionToken, // for server-side sessions AuthN reads token from cookie metadata
			})
			if err != nil {
				logger.Warn("ValidateToken failed", zap.Error(err), zap.String("path", r.URL.Path))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if resp == nil || !resp.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Device binding check for JWT: if client supplies X-Device-Id, it must match token claim.
			requestDeviceID := strings.TrimSpace(r.Header.Get("X-Device-Id"))
			if requestDeviceID != "" && resp.SessionType == "JWT" && resp.DeviceId != "" && requestDeviceID != resp.DeviceId {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Propagate all AuthN-validated identity fields as request headers
			// These are consumed by AuthZMiddleware and per-service gRPC interceptors.
			r.Header.Set("X-User-ID", resp.UserId)
			r.Header.Set("X-Session-ID", resp.SessionId)
			r.Header.Set("X-Session-Type", resp.SessionType)
			r.Header.Set("X-User-Type", resp.UserType)
			r.Header.Set("X-Portal", resp.Portal)      // system|business|b2b|agent|regulator|b2c
			r.Header.Set("X-Tenant-ID", resp.TenantId) // tenant UUID
			r.Header.Set("X-Token-ID", resp.TokenId)   // jti — for revocation lookup
			r.Header.Set("X-Device-ID", resp.DeviceId) // device fingerprint

			// Store in context for downstream handlers (including AuthZ middleware)
			ctx = context.WithValue(ctx, "user_id", resp.UserId)
			ctx = context.WithValue(ctx, "session_id", resp.SessionId)
			ctx = context.WithValue(ctx, "session_type", resp.SessionType)
			ctx = context.WithValue(ctx, "user_type", resp.UserType)
			ctx = context.WithValue(ctx, "portal", resp.Portal)
			ctx = context.WithValue(ctx, "tenant_id", resp.TenantId)
			ctx = context.WithValue(ctx, "token_id", resp.TokenId)
			ctx = context.WithValue(ctx, "device_id", resp.DeviceId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	parts := strings.SplitN(v, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return strings.TrimSpace(parts[1])
	}
	return v
}
