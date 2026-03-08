package handlers

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

const (
	sessionCookieName   = "session_token"
	sessionCookiePath   = "/"
	sessionCookieMaxAge = 12 * 60 * 60
)

func (h *AuthnHandler) Register(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RegisterRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.Register(ctx, &req)
	})
}

func (h *AuthnHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.SendOTPRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.SendOTP(ctx, &req)
	})
}

func (h *AuthnHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.VerifyOTPRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.VerifyOTP(ctx, &req)
	})
}

func (h *AuthnHandler) Login(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.LoginRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}

		resp, err := h.client.Login(ctx, &req)
		if err != nil {
			return resp, err
		}

		// Web portals: set cookie with *session token* and never expose it in JSON body.
		if resp != nil && resp.SessionType == "SERVER_SIDE" && resp.SessionToken != "" {
			setSessionCookie(w, resp.SessionToken, sessionCookieMaxAge, r.TLS != nil)
			resp.SessionToken = ""
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
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.RefreshToken(ctx, &req)
	})
}

func (h *AuthnHandler) Logout(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.LogoutRequest
		// Lenient body
		_ = protojson.Unmarshal(body, &req)

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
		if err := protojson.Unmarshal(body, &req); err != nil {
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
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ValidateCSRF(ctx, &req)
	})
}

func (h *AuthnHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ChangePasswordRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ChangePassword(ctx, &req)
	})
}

func (h *AuthnHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ResetPasswordRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
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
		_ = protojson.Unmarshal(body, &req) // allow empty body
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
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.BiometricAuthenticate(ctx, &req)
	})
}

// CreateAPIKey creates a new API key.
func (h *AuthnHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.CreateAPIKeyRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
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
		_ = protojson.Unmarshal(body, &req)
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
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.RegisterEmailUser(ctx, &req)
	})
}

func (h *AuthnHandler) SendEmailOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.SendEmailOTPRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.SendEmailOTP(ctx, &req)
	})
}

func (h *AuthnHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.VerifyEmailRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.VerifyEmail(ctx, &req)
	})
}

func (h *AuthnHandler) EmailLogin(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.EmailLoginRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
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

func (h *AuthnHandler) RequestPasswordResetByEmail(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RequestPasswordResetByEmailRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.RequestPasswordResetByEmail(ctx, &req)
	})
}

func (h *AuthnHandler) ResetPasswordByEmail(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ResetPasswordByEmailRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ResetPasswordByEmail(ctx, &req)
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
		http.Error(w, "JWKS unavailable", http.StatusServiceUnavailable)
		return
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		http.Error(w, "invalid public key", http.StatusInternalServerError)
		return
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		http.Error(w, "invalid public key format", http.StatusInternalServerError)
		return
	}
	rsaKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		http.Error(w, "not an RSA key", http.StatusInternalServerError)
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
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.CreateUserProfileRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateUserProfile(ctx, &req)
	})
}

func (h *AuthnHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.GetUserProfileRequest
		_ = protojson.Unmarshal(body, &req) // allow empty body for GET
		if req.UserId == "" {
			req.UserId = userID
		}
		return h.client.GetUserProfile(ctx, &req)
	})
}

func (h *AuthnHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.UpdateUserProfileRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.UpdateUserProfile(ctx, &req)
	})
}

func (h *AuthnHandler) GetProfilePhotoUploadURL(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.GetProfilePhotoUploadURLRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.GetProfilePhotoUploadURL(ctx, &req)
	})
}

func (h *AuthnHandler) UpdateNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.UpdateNotificationPreferencesRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.UpdateNotificationPreferences(ctx, &req)
	})
}

// ── TOTP / 2FA ────────────────────────────────────────────────────────────────

func (h *AuthnHandler) EnableTOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.EnableTOTPRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.EnableTOTP(ctx, &req)
	})
}

func (h *AuthnHandler) VerifyTOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.VerifyTOTPRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.VerifyTOTP(ctx, &req)
	})
}

func (h *AuthnHandler) DisableTOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.DisableTOTPRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.DisableTOTP(ctx, &req)
	})
}

// ── KYC ───────────────────────────────────────────────────────────────────────

func (h *AuthnHandler) InitiateKYC(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.InitiateKYCRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.InitiateKYC(ctx, &req)
	})
}

func (h *AuthnHandler) GetKYCStatus(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.GetKYCStatusRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.GetKYCStatus(ctx, &req)
	})
}

func (h *AuthnHandler) SubmitKYCFrame(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.SubmitKYCFrameRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.SubmitKYCFrame(ctx, &req)
	})
}

func (h *AuthnHandler) CompleteKYCSession(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.CompleteKYCSessionRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CompleteKYCSession(ctx, &req)
	})
}

func (h *AuthnHandler) ApproveKYC(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ApproveKYCRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ApproveKYC(ctx, &req)
	})
}

func (h *AuthnHandler) RejectKYC(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.RejectKYCRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.RejectKYC(ctx, &req)
	})
}

// ── Documents ─────────────────────────────────────────────────────────────────

func (h *AuthnHandler) UploadUserDocument(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.UploadUserDocumentRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
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
		if err := protojson.Unmarshal(body, &req); err != nil {
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
		if err := protojson.Unmarshal(body, &req); err != nil {
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
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateVoiceSession(ctx, &req)
	})
}

func (h *AuthnHandler) GetVoiceSession(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.GetVoiceSessionRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.GetVoiceSession(ctx, &req)
	})
}

func (h *AuthnHandler) EndVoiceSession(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.EndVoiceSessionRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.EndVoiceSession(ctx, &req)
	})
}

// --- shared helpers ---

// writeJSONError writes a clean JSON error response: {"ok":false,"message":"..."}.
func writeJSONError(w http.ResponseWriter, httpStatus int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	b, _ := json.Marshal(map[string]any{"ok": false, "message": msg})
	_, _ = w.Write(b)
}

func callUnary(w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, body []byte) (proto.Message, error)) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// forward metadata (same keys as authn metadata extractor)
	md := metadata.Pairs(
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
		"cookie", r.Header.Get("Cookie"),
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	msg, err := fn(ctx, body)
	if err != nil {
		st, _ := status.FromError(err)
		httpStatus := grpcStatusToHTTP(st.Code())
		logger.Warn("gRPC handler error",
			zap.String("path", r.URL.Path),
			zap.String("grpc_code", st.Code().String()),
			zap.Int("http_status", httpStatus),
			zap.String("message", st.Message()),
		)
		writeJSONError(w, httpStatus, st.Message())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, mErr := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if mErr != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to marshal response")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
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
