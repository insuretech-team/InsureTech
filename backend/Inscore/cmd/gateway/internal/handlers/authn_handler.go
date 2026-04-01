package handlers

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/respond"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// AuthnHandler exposes HTTP endpoints that translate to AuthN gRPC.
// Hybrid auth:
//   - Web portals: SERVER_SIDE session -> sets HttpOnly cookie "session_token" with *session token* (not session_id)
//   - Mobile/B2C: JWT -> tokens in response body
//
// This implementation is based on your archived gateway handlers, but moved into cmd/gateway/internal/handlers.
type AuthnHandler struct {
	client authnservicev1.AuthServiceClient
}

func NewAuthnHandler(conn *grpc.ClientConn) *AuthnHandler {
	return &AuthnHandler{client: authnservicev1.NewAuthServiceClient(conn)}
}

// protoUnmarshal unmarshals JSON into a proto message.
// DiscardUnknown=true: unknown fields are silently ignored instead of causing 500.
// This is the correct REST API behavior — be liberal in what you accept.
// Invalid JSON or type mismatches still return errors (mapped to 400 by callUnary).
var protoUnmarshal = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}.Unmarshal

const (
	sessionCookieName   = "session_token"
	sessionCookiePath   = "/"
	sessionCookieMaxAge = 12 * 60 * 60
)

func (h *AuthnHandler) Register(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RegisterRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.Register(ctx, &req)
	})
}

func (h *AuthnHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.SendOTPRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.SendOTP(ctx, &req)
	})
}

func (h *AuthnHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.VerifyOTPRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.VerifyOTP(ctx, &req)
	})
}

func (h *AuthnHandler) Login(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.LoginRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}

		resp, err := h.client.Login(ctx, &req)
		if err != nil {
			return resp, err
		}

		// Web portals: set HttpOnly cookie with session token.
		// We also keep session_token in the JSON response body so that server-side
		// Next.js BFF route handlers (e.g. b2b_portal login route) can read it
		// without relying on Set-Cookie header forwarding, which is unreliable
		// when the response goes through an SDK interceptor that rewrites the body
		// using `new Response(body, { headers })` — the Fetch API spec forbids
		// Set-Cookie in the Headers constructor, silently dropping it.
		// The token is safe here: the b2b portal's login route is server-side only
		// and immediately re-sets it as its own HttpOnly cookie without exposing
		// it to browser JS.
		if resp != nil && resp.SessionType == "SERVER_SIDE" && resp.SessionToken != "" {
			setSessionCookie(w, resp.SessionToken, sessionCookieMaxAge, r.TLS != nil)
			// NOTE: intentionally NOT clearing resp.SessionToken so the BFF can
			// read it from result.data.session_token on the server side.
			if resp.CsrfToken != "" {
				w.Header().Set("X-CSRF-Token", resp.CsrfToken)
			}
		}

		return resp, nil
	})
}

func (h *AuthnHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RefreshTokenRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.RefreshToken(ctx, &req)
	})
}

func (h *AuthnHandler) Logout(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.LogoutRequest
		// Lenient body
		_ = protoUnmarshal(body, &req)

		resp, err := h.client.Logout(ctx, &req)
		if err == nil {
			clearSessionCookie(w, r.TLS != nil)
		}
		return resp, err
	})
}

func (h *AuthnHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ValidateTokenRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}

		// If caller didn't provide session_id, allow cookie-based validation.
		// Cookie is session_token; middleware forwards cookie anyway.
		return h.client.ValidateToken(ctx, &req)
	})
}

func (h *AuthnHandler) ValidateCSRF(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ValidateCSRFRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ValidateCSRF(ctx, &req)
	})
}

// GetCSRFToken issues a fresh CSRF token for an authenticated session.
// GET /v1/auth/csrf-token — requires auth middleware (session_token cookie or Bearer).
func (h *AuthnHandler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetCurrentSession(ctx, &authnservicev1.GetCurrentSessionRequest{})
	})
}

func (h *AuthnHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ChangePasswordRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ChangePassword(ctx, &req)
	})
}

func (h *AuthnHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ResetPasswordRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ResetPassword(ctx, &req)
	})
}

func (h *AuthnHandler) GetCurrentSession(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetCurrentSession(ctx, &authnservicev1.GetCurrentSessionRequest{})
	})
}

func (h *AuthnHandler) FindPortalUser(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.FindPortalUserRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.FindPortalUser(ctx, &req)
	})
}

func (h *AuthnHandler) SetTemporaryPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.SetTemporaryPasswordRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		req.UserId = userID
		return h.client.SetTemporaryPassword(ctx, &req)
	})
}

func (h *AuthnHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetSession(ctx, &authnservicev1.GetSessionRequest{SessionId: sessionID})
	})
}

func (h *AuthnHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	q := r.URL.Query()

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	activeOnly := q.Get("active_only") == "true"

	req := &authnservicev1.ListSessionsRequest{
		UserId:      userID,
		PageSize:    int32(pageSize),
		PageToken:   q.Get("page_token"),
		SessionType: q.Get("session_type"),
		ActiveOnly:  activeOnly,
		DeviceType:  q.Get("device_type"),
	}

	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListSessions(ctx, req)
	})
}

func (h *AuthnHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		// proto route uses DELETE without body; allow reason via query string.
		req := &authnservicev1.RevokeSessionRequest{
			SessionId: sessionID,
			Reason:    r.URL.Query().Get("reason"),
		}
		return h.client.RevokeSession(ctx, req)
	})
}

func (h *AuthnHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RevokeAllSessionsRequest
		_ = protoUnmarshal(body, &req) // allow empty body
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.RevokeAllSessions(ctx, &req)
	})
}

// BiometricAuthenticate handles mobile biometric login.
func (h *AuthnHandler) BiometricAuthenticate(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.BiometricAuthenticateRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.BiometricAuthenticate(ctx, &req)
	})
}

// CreateAPIKey creates a new API key.
func (h *AuthnHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.CreateAPIKeyRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateAPIKey(ctx, &req)
	})
}

// ListAPIKeys lists API keys for an owner.
func (h *AuthnHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	ownerID := r.URL.Query().Get("owner_id")
	req := &authnservicev1.ListAPIKeysRequest{
		OwnerId:    ownerID,
		OwnerType:  r.URL.Query().Get("owner_type"),
		ActiveOnly: r.URL.Query().Get("active_only") == "true",
	}
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListAPIKeys(ctx, req)
	})
}

// RevokeAPIKey revokes an API key.
func (h *AuthnHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	keyID := r.PathValue("key_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RevokeAPIKeyRequest
		_ = protoUnmarshal(body, &req)
		if req.KeyId == "" {
			req.KeyId = keyID
		}
		return h.client.RevokeAPIKey(ctx, &req)
	})
}

// Email flows
func (h *AuthnHandler) RegisterEmailUser(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RegisterEmailUserRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.RegisterEmailUser(ctx, &req)
	})
}

func (h *AuthnHandler) SendEmailOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.SendEmailOTPRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.SendEmailOTP(ctx, &req)
	})
}

func (h *AuthnHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.VerifyEmailRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.VerifyEmail(ctx, &req)
	})
}

func (h *AuthnHandler) EmailLogin(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.EmailLoginRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}

		resp, err := h.client.EmailLogin(ctx, &req)
		if err != nil {
			return resp, err
		}

		if resp != nil && resp.SessionToken != "" {
			setSessionCookie(w, resp.SessionToken, sessionCookieMaxAge, r.TLS != nil)
			resp.SessionToken = ""
			if resp.CsrfToken != "" {
				w.Header().Set("X-CSRF-Token", resp.CsrfToken)
			}
		}

		return resp, nil
	})
}

func (h *AuthnHandler) EmailPasswordLogin(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.EmailPasswordLoginRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}

		resp, err := h.client.EmailPasswordLogin(ctx, &req)
		if err != nil {
			return resp, err
		}

		if resp != nil && resp.SessionToken != "" {
			setSessionCookie(w, resp.SessionToken, sessionCookieMaxAge, r.TLS != nil)
			if resp.CsrfToken != "" {
				w.Header().Set("X-CSRF-Token", resp.CsrfToken)
			}
		}

		return resp, nil
	})
}

func (h *AuthnHandler) RequestPasswordResetByEmail(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RequestPasswordResetByEmailRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.RequestPasswordResetByEmail(ctx, &req)
	})
}

func (h *AuthnHandler) ResetPasswordByEmail(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ResetPasswordByEmailRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ResetPasswordByEmail(ctx, &req)
	})
}

func (h *AuthnHandler) ProvisionEmployeeUser(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ProvisionEmployeeUserRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ProvisionEmployeeUser(ctx, &req)
	})
}

// ── JWKS ─────────────────────────────────────────────────────────────────────

// JWKS serves the RS256 public key set for JWT verification.
// GET /.well-known/jwks.json — no auth required.
// The public key PEM path is read from JWT_PUBLIC_KEY_PATH env var.
func (h *AuthnHandler) JWKS(w http.ResponseWriter, r *http.Request) {
	pubKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if pubKeyPath == "" {
		pubKeyPath = "/secrets/jwt_rsa_public.pem"
	}
	kid := os.Getenv("JWT_KEY_ID")
	if kid == "" {
		kid = "insuretech-2025-01"
	}

	pemBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		respond.Error(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "JWKS unavailable")
		return
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		respond.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "invalid public key")
		return
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		respond.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "invalid public key format")
		return
	}
	rsaKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		respond.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "not an RSA key")
		return
	}

	nBytes := rsaKey.N.Bytes()
	eVal := rsaKey.E
	eBuf := []byte{byte(eVal >> 16), byte(eVal >> 8), byte(eVal)}
	if eVal < 1<<16 {
		eBuf = eBuf[1:]
	}
	if eVal < 1<<8 {
		eBuf = eBuf[1:]
	}

	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"use": "sig",
				"alg": "RS256",
				"kid": kid,
				"n":   base64.RawURLEncoding.EncodeToString(nBytes),
				"e":   base64.RawURLEncoding.EncodeToString(eBuf),
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_ = json.NewEncoder(w).Encode(jwks)
}

// ── User Profile ──────────────────────────────────────────────────────────────

func (h *AuthnHandler) CreateUserProfile(w http.ResponseWriter, r *http.Request) {
	// BUG FIX: Profile date_of_birth comes as plain date string "1990-05-15" from mobile/web clients.
	// Proto Timestamp expects {"seconds":..., "nanos":...} in JSON. We pre-process the body to
	// convert ISO date strings to proto Timestamp format before unmarshalling.
	// Also: Mobile/web clients may send a nested "address" object — flatten it to proto flat fields.
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		body = normalizeUserProfileBody(body)
		body = convertDateFieldsToProtoTimestamp(body, "date_of_birth")
		var req authnservicev1.CreateUserProfileRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.CreateUserProfile(ctx, &req)
	})
}

// normalizeUserProfileBody flattens a nested "address" object into flat proto fields.
// Mobile/web clients may send:
//   {"address": {"line1": "123 St", "line2": "Apt 4", "city": "Dhaka", "district": "Dhaka",
//                "division": "Dhaka", "country": "BD", "zip_code": "1200"}}
// But the proto has flat fields: address_line1, address_line2, city, district, division, country, zip_code.
func normalizeUserProfileBody(body []byte) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return body
	}
	changed := false

	// Flatten nested "address" object if present
	if addr, ok := m["address"]; ok {
		if addrMap, ok := addr.(map[string]interface{}); ok {
			fieldMap := map[string]string{
				"line1":    "address_line1",
				"line2":    "address_line2",
				"city":     "city",
				"district": "district",
				"division": "division",
				"country":  "country",
				"zip_code": "zip_code",
			}
			for srcKey, dstKey := range fieldMap {
				if v, exists := addrMap[srcKey]; exists {
					if _, alreadySet := m[dstKey]; !alreadySet {
						m[dstKey] = v
					}
				}
			}
			delete(m, "address")
			changed = true
		}
	}

	// Normalize first_name + last_name from full_name if needed
	// Keep full_name in the map so the proto full_name field is also set.
	if fullName, ok := m["full_name"]; ok {
		if _, hasFirst := m["first_name"]; !hasFirst {
			if nameStr, ok := fullName.(string); ok {
				parts := strings.SplitN(nameStr, " ", 2)
				m["first_name"] = parts[0]
				if len(parts) > 1 {
					m["last_name"] = parts[1]
				}
				// Do NOT delete full_name — proto has both full_name and first_name/last_name fields
				changed = true
			}
		}
	}

	if !changed {
		return body
	}
	out, err := json.Marshal(m)
	if err != nil {
		return body
	}
	return out
}

// convertDateFieldsToProtoTimestamp rewrites plain ISO date fields (e.g. "1990-05-15")
// to proto Timestamp JSON format {"seconds": N, "nanos": 0} for the given field names.
// Handles both snake_case (date_of_birth) and camelCase (dateOfBirth) field variants.
// This makes the API client-friendly — clients don't need to know about proto Timestamp format.
func convertDateFieldsToProtoTimestamp(body []byte, fields ...string) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return body
	}
	changed := false
	for _, field := range fields {
		// Build list of candidate key names: snake_case and camelCase
		candidates := []string{field}
		// snake_case → camelCase conversion (date_of_birth → dateOfBirth)
		camel := snakeToCamel(field)
		if camel != field {
			candidates = append(candidates, camel)
		}
		for _, key := range candidates {
			val, ok := m[key]
			if !ok {
				continue
			}
			str, ok := val.(string)
			if !ok {
				continue
			}
			// Try parsing as ISO date (YYYY-MM-DD) or RFC3339
			var t time.Time
			var err error
			if len(str) == 10 { // "1990-05-15"
				t, err = time.Parse("2006-01-02", str)
			} else {
				t, err = time.Parse(time.RFC3339, str)
			}
			if err == nil {
				// protojson accepts both snake_case and camelCase field names.
				// To avoid "duplicate field" error when both date_of_birth and dateOfBirth
				// are present in the map, we delete all candidate keys and write ONLY
				// the snake_case key (protojson's DiscardUnknown=false accepts it).
				rfc3339 := t.UTC().Format(time.RFC3339)
				// Remove all candidate keys first (prevent duplicates)
				for _, c := range candidates {
					delete(m, c)
				}
				// Write only the original snake_case field name
				m[field] = rfc3339
				changed = true
				break // found and converted, no need to check other candidates
			}
		}
	}
	if !changed {
		return body
	}
	out, err := json.Marshal(m)
	if err != nil {
		return body
	}
	return out
}

// snakeToCamel converts snake_case to camelCase (e.g. "date_of_birth" → "dateOfBirth").
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}
	result := parts[0]
	for _, p := range parts[1:] {
		if len(p) > 0 {
			result += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return result
}

func (h *AuthnHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.GetUserProfileRequest
		_ = protoUnmarshal(body, &req) // allow empty body for GET
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.GetUserProfile(ctx, &req)
	})
}

func (h *AuthnHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// BUG-008 FIX: user_id must not be required in body when it's already in the URL path.
	// BUG FIX: date_of_birth converted from ISO string "YYYY-MM-DD" to proto Timestamp format.
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		body = normalizeUserProfileBody(body)
		body = convertDateFieldsToProtoTimestamp(body, "date_of_birth")
		var req authnservicev1.UpdateUserProfileRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.UpdateUserProfile(ctx, &req)
	})
}

func (h *AuthnHandler) GetProfilePhotoUploadURL(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.GetProfilePhotoUploadURLRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.GetProfilePhotoUploadURL(ctx, &req)
	})
}

// GetNotificationPreferences handles GET /v1/auth/users/{user_id}/notification-preferences
// BUG-010 FIX: Implements the missing GET endpoint.
func (h *AuthnHandler) GetNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetNotificationPreferences(ctx, &authnservicev1.GetNotificationPreferencesRequest{
			UserId: userID,
		})
	})
}

func (h *AuthnHandler) UpdateNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	// BUG-009 FIX: user_id must not be required in body when it's already in the URL path.
	// Populate UserId from path param if body doesn't include it.
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.UpdateNotificationPreferencesRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.UpdateNotificationPreferences(ctx, &req)
	})
}

// ── TOTP / 2FA ────────────────────────────────────────────────────────────────

func (h *AuthnHandler) EnableTOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.EnableTOTPRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.EnableTOTP(ctx, &req)
	})
}

func (h *AuthnHandler) VerifyTOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.VerifyTOTPRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.VerifyTOTP(ctx, &req)
	})
}

func (h *AuthnHandler) DisableTOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.DisableTOTPRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.DisableTOTP(ctx, &req)
	})
}

// ── KYC ───────────────────────────────────────────────────────────────────────

func (h *AuthnHandler) InitiateKYC(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.InitiateKYCRequest
		_ = protoUnmarshal(body, &req) // allow empty body
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.InitiateKYC(ctx, &req)
	})
}

func (h *AuthnHandler) GetKYCStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &authnservicev1.GetKYCStatusRequest{
			UserId: userID,
		}
		return h.client.GetKYCStatus(ctx, req)
	})
}

func (h *AuthnHandler) SubmitKYCFrame(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")

	// SubmitKYCFrame accepts multipart/form-data (webcam frame as binary file).
	// Fields: session_id (string), image_data (binary JPEG file).
	// Falls back to JSON body for non-multipart clients.
	ct := r.Header.Get("Content-Type")
	if strings.Contains(ct, "multipart/form-data") {
		ctx := r.Context()
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
			respond.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "failed to parse multipart form")
			return
		}
		sessionID := r.FormValue("session_id")
		var imageData []byte
		// Accept "image_data" or "file" field name
		for _, fieldName := range []string{"image_data", "file"} {
			f, _, err := r.FormFile(fieldName)
			if err == nil {
				imageData, _ = io.ReadAll(f)
				f.Close()
				break
			}
		}
		// Build gRPC context with the same metadata as callUnary
		cookieHeader := r.Header.Get("Cookie")
		if cookieHeader == "" {
			if st := strings.TrimSpace(r.Header.Get("X-Session-Token")); st != "" {
				cookieHeader = "session_token=" + st
			}
		}
		grpcCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		grpcCtx = grpcmeta.WithOutgoingMetadata(grpcCtx,
			"cookie", cookieHeader,
			"x-user-id", r.Header.Get("X-User-ID"),
			"x-tenant-id", r.Header.Get("X-Tenant-ID"),
			"x-portal", r.Header.Get("X-Portal"),
			"x-session-id", r.Header.Get("X-Session-ID"),
			"x-user-type", r.Header.Get("X-User-Type"),
			"x-business-id", r.Header.Get("X-Business-ID"),
			"authorization", r.Header.Get("Authorization"),
		)

		req := &authnservicev1.SubmitKYCFrameRequest{
			UserId:    userID,
			SessionId: sessionID,
			ImageData: imageData,
		}
		resp, err := h.client.SubmitKYCFrame(grpcCtx, req)
		if err != nil {
			respond.GRPCError(w, r, err)
			return
		}
		b, _ := protojson.MarshalOptions{UseProtoNames: true}.Marshal(resp)
		respond.RawProtoJSON(w, r, b, http.StatusOK)
		return
	}

	// JSON fallback
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.SubmitKYCFrameRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.SubmitKYCFrame(ctx, &req)
	})
}

func (h *AuthnHandler) CompleteKYCSession(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.CompleteKYCSessionRequest
		_ = protoUnmarshal(body, &req) // allow empty body
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.CompleteKYCSession(ctx, &req)
	})
}

func (h *AuthnHandler) ApproveKYC(w http.ResponseWriter, r *http.Request) {
	kycID := r.PathValue("kyc_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ApproveKYCRequest
		_ = protoUnmarshal(body, &req) // allow empty body
		if req.KycId == "" {
			req.KycId = kycID
		}
		return h.client.ApproveKYC(ctx, &req)
	})
}

func (h *AuthnHandler) RejectKYC(w http.ResponseWriter, r *http.Request) {
	// RPC removed from proto (API path conflict). Route removed from gateway router.
	respond.Error(w, r, http.StatusNotFound, "NOT_FOUND", "RejectKYC endpoint not available")
}

// ── Documents ─────────────────────────────────────────────────────────────────

func (h *AuthnHandler) UploadUserDocument(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.UploadUserDocumentRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.UserId == "" {
			req.UserId = r.PathValue("user_id")
		}
		return h.client.UploadUserDocument(ctx, &req)
	})
}

func (h *AuthnHandler) ListUserDocuments(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := authnservicev1.ListUserDocumentsRequest{
			UserId:         r.PathValue("user_id"),
			DocumentTypeId: r.URL.Query().Get("document_type_id"),
			PageToken:      r.URL.Query().Get("page_token"),
		}
		if ps := r.URL.Query().Get("page_size"); ps != "" {
			if n, err := strconv.Atoi(ps); err == nil {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListUserDocuments(ctx, &req)
	})
}

func (h *AuthnHandler) ListDocumentTypes(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListDocumentTypes(ctx, &authnservicev1.ListDocumentTypesRequest{})
	})
}

func (h *AuthnHandler) GetUserDocument(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := authnservicev1.GetUserDocumentRequest{
			UserDocumentId: r.PathValue("user_document_id"),
		}
		return h.client.GetUserDocument(ctx, &req)
	})
}

func (h *AuthnHandler) UpdateUserDocument(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.UpdateUserDocumentRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.UserDocumentId == "" {
			req.UserDocumentId = r.PathValue("user_document_id")
		}
		return h.client.UpdateUserDocument(ctx, &req)
	})
}

func (h *AuthnHandler) DeleteUserDocument(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := authnservicev1.DeleteUserDocumentRequest{
			UserDocumentId: r.PathValue("user_document_id"),
		}
		return h.client.DeleteUserDocument(ctx, &req)
	})
}

func (h *AuthnHandler) VerifyDocument(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.VerifyDocumentRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		if req.UserDocumentId == "" {
			req.UserDocumentId = r.PathValue("user_document_id")
		}
		return h.client.VerifyDocument(ctx, &req)
	})
}

func (h *AuthnHandler) CreateVoiceSession(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.CreateVoiceSessionRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateVoiceSession(ctx, &req)
	})
}

func (h *AuthnHandler) GetVoiceSession(w http.ResponseWriter, r *http.Request) {
	// RPC removed from proto (API path conflict with voice domain). Route removed from gateway router.
	respond.Error(w, r, http.StatusNotFound, "NOT_FOUND", "GetVoiceSession endpoint not available")
}

func (h *AuthnHandler) EndVoiceSession(w http.ResponseWriter, r *http.Request) {
	// RPC removed from proto (API path conflict with voice domain). Route removed from gateway router.
	respond.Error(w, r, http.StatusNotFound, "NOT_FOUND", "EndVoiceSession endpoint not available")
}

// --- shared helpers ---

// writeJSONError writes a unified ApiResponse error envelope.
// Deprecated: call respond.Error() directly; kept for any remaining
// internal callsites that haven't been migrated yet.
func writeJSONError(w http.ResponseWriter, r *http.Request, httpStatus int, msg string) {
	code := httpStatusToErrorCode(httpStatus)
	respond.Error(w, r, httpStatus, code, msg)
}

// httpStatusToErrorCode maps common HTTP status codes to error code strings.
func httpStatusToErrorCode(httpStatus int) string {
	switch httpStatus {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHENTICATED"
	case http.StatusForbidden:
		return "PERMISSION_DENIED"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnprocessableEntity:
		return "VALIDATION_ERROR"
	case http.StatusTooManyRequests:
		return "RATE_LIMITED"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	case http.StatusGatewayTimeout:
		return "DEADLINE_EXCEEDED"
	default:
		return "INTERNAL_ERROR"
	}
}

// callUnary reads the request body, calls the given gRPC function, and writes
// a unified ApiResponse envelope — success data is wrapped in data:{...} and
// errors are wrapped in error:{...} with data: null.
func callUnary(w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, body []byte) (proto.Message, error)) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "failed to read request body")
		return
	}

	// Email OTP sends (SMTP) can take up to 30s; use a longer timeout for all unary calls.
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Build cookie header for gRPC metadata.
	// If client sent X-Session-Token (Postman/API), synthesise a cookie so authn
	// services (GetCurrentSession, ValidateCSRF etc.) can extract the session token.
	// Browsers send the actual Cookie header; API clients use X-Session-Token.
	cookieHeader := r.Header.Get("Cookie")
	if cookieHeader == "" {
		if st := strings.TrimSpace(r.Header.Get("X-Session-Token")); st != "" {
			cookieHeader = "session_token=" + st
		}
	}

	// forward metadata (same keys as authn metadata extractor)
	ctx = grpcmeta.WithOutgoingMetadata(ctx,
		"x-forwarded-for", r.Header.Get("X-Forwarded-For"),
		"x-real-ip", r.Header.Get("X-Real-Ip"),
		"user-agent", r.UserAgent(),
		"x-device-id", r.Header.Get("X-Device-Id"),
		"x-csrf-token", r.Header.Get("X-CSRF-Token"),
		"x-user-id", r.Header.Get("X-User-ID"),
		"x-tenant-id", r.Header.Get("X-Tenant-ID"),
		"x-portal", r.Header.Get("X-Portal"),
		"x-session-id", r.Header.Get("X-Session-ID"),
		"x-user-type", r.Header.Get("X-User-Type"),
		"x-business-id", r.Header.Get("X-Business-ID"),
		"x-org-role", r.Header.Get("X-Org-Role"),
		"authorization", r.Header.Get("Authorization"),
		"cookie", cookieHeader,
	)

	msg, err := fn(ctx, body)
	if err != nil {
		st, _ := status.FromError(err)
		logger.Warn("gRPC handler error",
			zap.String("path", r.URL.Path),
			zap.String("grpc_code", st.Code().String()),
			zap.Int("http_status", respond.GRPCCodeToHTTP(st.Code())),
			zap.String("message", st.Message()),
		)
		respond.GRPCError(w, r, err)
		return
	}

	b, mErr := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if mErr != nil {
		respond.Error(w, r, http.StatusInternalServerError, "MARSHAL_ERROR", "failed to marshal response")
		return
	}
	respond.RawProtoJSON(w, r, b, http.StatusOK)
}

func setSessionCookie(w http.ResponseWriter, token string, maxAge int, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     sessionCookiePath,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

func clearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     sessionCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}
