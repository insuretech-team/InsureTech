package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/middleware"
	"github.com/redis/go-redis/v9"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TokenService handles all token operations with production-grade security.
// JWT signing uses RS256 (RSA 2048-bit). HS256 is not supported.
// Keys are loaded once at startup from PEM files; private key never leaves the service.
type TokenService struct {
	sessionRepo    *repository.SessionRepository
	userRepo       *repository.UserRepository
	config         *config.Config
	eventPublisher *events.Publisher
	metadata       *middleware.MetadataExtractor
	refreshLimiter *refreshRateLimiter
	sessionLimiter *SessionLimiter
	rdb            redis.UniversalClient

	// RS256 key pair — loaded once at startup
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
	keyID         string // kid header value — must match TokenConfig.kid in authz DB
}

// TokenPair represents a JWT access/refresh token pair
type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	SessionID             string
	AccessTokenExpiresIn  time.Duration
	RefreshTokenExpiresIn time.Duration
}

// ServerSideSession represents a server-side session with CSRF protection
type ServerSideSession struct {
	SessionID    string
	SessionToken string // Plain token to be set as HttpOnly cookie
	CSRFToken    string
	ExpiresIn    time.Duration
}

// InsureTechClaims defines the custom JWT claims for InsureTech tokens.
// All claims follow RFC 7519. Custom claims use "ins_" prefix.
type InsureTechClaims struct {
	jwt.RegisteredClaims
	// Standard identity
	UserType string `json:"utp"` // user type: B2C_CUSTOMER | AGENT | SYSTEM_USER | ...
	// AuthZ context — used by gateway to call AuthZ.CheckAccess
	Portal   string `json:"ins_portal"` // portal: system | business | b2b | agent | regulator | b2c
	TenantID string `json:"ins_tenant"` // tenant_id UUID
	DeviceID string `json:"ins_device"` // device fingerprint (for binding validation)
	// Session linkage
	SessionID string `json:"sid"`      // session_id — links to sessions table for revocation
	TokenType string `json:"ins_type"` // "access" | "refresh"
}

// NewTokenService creates a TokenService, loading RS256 keys from PEM files at startup.
func NewTokenService(
	sessionRepo *repository.SessionRepository,
	userRepo *repository.UserRepository,
	cfg *config.Config,
	eventPublisher *events.Publisher,
	metadata *middleware.MetadataExtractor,
) (*TokenService, error) {
	svc := &TokenService{
		sessionRepo:    sessionRepo,
		userRepo:       userRepo,
		config:         cfg,
		eventPublisher: eventPublisher,
		metadata:       metadata,
		refreshLimiter: newRefreshRateLimiter(nil),
		keyID:          cfg.JWT.KeyID,
	}
	if err := svc.loadRSAKeys(); err != nil {
		logger.Errorf("failed to load RS256 keys: %v", err)
		return nil, errors.New("failed to load RS256 keys")
	}
	return svc, nil
}

// NewTokenServiceWithRedis creates a TokenService with a Redis-backed refresh rate limiter
// and JTI blocklist support.
func NewTokenServiceWithRedis(
	sessionRepo *repository.SessionRepository,
	userRepo *repository.UserRepository,
	cfg *config.Config,
	eventPublisher *events.Publisher,
	metadata *middleware.MetadataExtractor,
	rdb redis.UniversalClient,
) (*TokenService, error) {
	svc := &TokenService{
		sessionRepo:    sessionRepo,
		userRepo:       userRepo,
		config:         cfg,
		eventPublisher: eventPublisher,
		metadata:       metadata,
		refreshLimiter: newRefreshRateLimiter(rdb),
		rdb:            rdb,
		keyID:          cfg.JWT.KeyID,
	}
	if err := svc.loadRSAKeys(); err != nil {
		logger.Errorf("failed to load RS256 keys: %v", err)
		return nil, errors.New("failed to load RS256 keys")
	}
	return svc, nil
}

// NewTokenServiceWithSessionLimiter creates a TokenService with both a Redis-backed
// refresh rate limiter and a concurrent session limiter.
// maxSessions ≤ 0 defaults to 5.
func NewTokenServiceWithSessionLimiter(
	sessionRepo *repository.SessionRepository,
	userRepo *repository.UserRepository,
	cfg *config.Config,
	eventPublisher *events.Publisher,
	metadata *middleware.MetadataExtractor,
	rdb redis.UniversalClient,
	maxSessions int,
) (*TokenService, error) {
	svc := &TokenService{
		sessionRepo:    sessionRepo,
		userRepo:       userRepo,
		config:         cfg,
		eventPublisher: eventPublisher,
		metadata:       metadata,
		refreshLimiter: newRefreshRateLimiter(rdb),
		sessionLimiter: NewSessionLimiter(rdb, maxSessions),
		rdb:            rdb,
		keyID:          cfg.JWT.KeyID,
	}
	if err := svc.loadRSAKeys(); err != nil {
		logger.Errorf("failed to load RS256 keys: %v", err)
		return nil, errors.New("failed to load RS256 keys")
	}
	return svc, nil
}

// BlockJTI writes a jti to Redis with TTL = remaining token lifetime.
// Key format: revoked:jti:<token_id>
func (s *TokenService) BlockJTI(ctx context.Context, jti string, ttl time.Duration) error {
	if s.rdb == nil || jti == "" {
		return nil
	}
	if ttl <= 0 {
		return nil
	}
	key := "revoked:jti:" + jti
	return s.rdb.Set(ctx, key, "1", ttl).Err()
}

// isJTIBlocked checks if a jti is in the Redis blocklist.
func (s *TokenService) isJTIBlocked(ctx context.Context, jti string) bool {
	if s.rdb == nil || jti == "" {
		return false
	}
	key := "revoked:jti:" + jti
	val, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return val == "1"
}

// loadRSAKeys reads RSA private and public key PEM files from disk.
// The private key file must be accessible only to the authn service process.
func (s *TokenService) loadRSAKeys() error {
	// Load private key
	privPEM, err := os.ReadFile(s.config.JWT.PrivateKeyPath)
	if err != nil {
		return errors.New("read private key file " + s.config.JWT.PrivateKeyPath + ": " + err.Error())
	}
	privBlock, _ := pem.Decode(privPEM)
	if privBlock == nil {
		return errors.New("failed to decode PEM block from private key file")
	}
	privKey, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
	if err != nil {
		// fallback: try PKCS1
		privKey2, err2 := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
		if err2 != nil {
			return errors.New("parse private key (PKCS8: " + err.Error() + ", PKCS1: " + err2.Error() + ")")
		}
		s.rsaPrivateKey = privKey2
	} else {
		rsaKey, ok := privKey.(*rsa.PrivateKey)
		if !ok {
			return errors.New("private key is not RSA")
		}
		s.rsaPrivateKey = rsaKey
	}

	// Load public key
	pubPEM, err := os.ReadFile(s.config.JWT.PublicKeyPath)
	if err != nil {
		return errors.New("read public key file " + s.config.JWT.PublicKeyPath + ": " + err.Error())
	}
	pubBlock, _ := pem.Decode(pubPEM)
	if pubBlock == nil {
		return errors.New("failed to decode PEM block from public key file")
	}
	pubKeyInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		logger.Errorf("parse public key: %v", err)
		return errors.New("parse public key")
	}
	rsaPubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return errors.New("public key is not RSA")
	}
	s.rsaPublicKey = rsaPubKey
	return nil
}

// PublicKey returns the RSA public key (for JWKS endpoint serving).
func (s *TokenService) PublicKey() *rsa.PublicKey { return s.rsaPublicKey }
func (s *TokenService) KeyID() string             { return s.keyID }

// portalForUserType maps UserType → portal name (used in JWT ins_portal claim).
func portalForUserType(userType string) string {
	switch userType {
	case "USER_TYPE_SYSTEM_USER", authnentityv1.UserType_USER_TYPE_SYSTEM_USER.String():
		return "system"
	case "USER_TYPE_BUSINESS_BENEFICIARY", authnentityv1.UserType_USER_TYPE_BUSINESS_BENEFICIARY.String():
		return "business"
	case "USER_TYPE_PARTNER", authnentityv1.UserType_USER_TYPE_PARTNER.String():
		return "b2b"
	case "USER_TYPE_AGENT", authnentityv1.UserType_USER_TYPE_AGENT.String():
		return "agent"
	case "USER_TYPE_REGULATOR", authnentityv1.UserType_USER_TYPE_REGULATOR.String():
		return "regulator"
	case "USER_TYPE_B2C_CUSTOMER", authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER.String():
		return "b2c"
	default:
		return "b2c"
	}
}

// GenerateServerSideSession creates a server-side session for web portals.
// Session token is bcrypt-hashed for storage; plain token returned for HttpOnly cookie.
func (s *TokenService) GenerateServerSideSession(ctx context.Context, userID, deviceID string, deviceType authnentityv1.DeviceType, ipAddress, userAgent string) (*ServerSideSession, error) {
	sessionID := uuid.New().String()
	expiresIn := s.config.Security.ServerSessionDuration
	if expiresIn == 0 {
		expiresIn = 12 * time.Hour
	}

	sessionToken := uuid.New().String()
	sessionTokenHash, err := bcrypt.GenerateFromPassword([]byte(sessionToken), s.config.Security.BCryptCost)
	if err != nil {
		logger.Errorf("failed to hash session token: %v", err)
		return nil, errors.New("failed to hash session token")
	}

	csrfToken, err := generateSecureRandomString(32)
	if err != nil {
		logger.Errorf("failed to generate CSRF token: %v", err)
		return nil, errors.New("failed to generate CSRF token")
	}

	session := &authnentityv1.Session{
		SessionId:          sessionID,
		UserId:             userID,
		DeviceId:           deviceID,
		DeviceType:         deviceType,
		SessionType:        authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE,
		SessionTokenHash:   string(sessionTokenHash),
		SessionTokenLookup: sessionTokenLookup(sessionToken),
		CsrfToken:          csrfToken,
		IpAddress:          ipAddress,
		UserAgent:          userAgent,
		ExpiresAt:          timestamppb.New(time.Now().Add(expiresIn)),
		LastActivityAt:     timestamppb.New(time.Now()),
		IsActive:           true,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		logger.Errorf("failed to create session: %v", err)
		return nil, errors.New("failed to create session")
	}

	// Sprint 3: Initialize idle timeout sliding window key in Redis
	if s.rdb != nil && s.config.Security.IdleTimeoutDuration > 0 {
		idleKey := "session:idle:" + sessionID
		s.rdb.Set(ctx, idleKey, "1", s.config.Security.IdleTimeoutDuration)
	}

	// Enforce concurrent session limit (evict oldest sessions if over limit).
	if s.sessionLimiter != nil {
		evicted, err := s.sessionLimiter.TrackSession(ctx, userID, sessionID, time.Now().Add(expiresIn))
		if err == nil {
			for _, evictedID := range evicted {
				_ = s.RevokeSession(ctx, evictedID)
			}
		}
	}

	return &ServerSideSession{
		SessionID:    sessionID,
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// GenerateJWT creates RS256-signed JWT access and refresh tokens for mobile/API clients.
// Claims include: sub, utp, ins_portal, ins_tenant, ins_device, sid, ins_type, jti, iss, aud, exp, iat, kid.
func (s *TokenService) GenerateJWT(ctx context.Context, userID, userType, tenantID, deviceID string, deviceType authnentityv1.DeviceType, ipAddress, userAgent string) (*TokenPair, error) {
	accessExpiresIn := s.config.JWT.AccessTokenDuration
	refreshExpiresIn := s.config.JWT.RefreshTokenDuration
	if accessExpiresIn == 0 {
		accessExpiresIn = 15 * time.Minute
	}
	if refreshExpiresIn == 0 {
		refreshExpiresIn = 7 * 24 * time.Hour
	}

	sessionID := uuid.New().String()
	accessJTI := uuid.New().String()
	refreshJTI := uuid.New().String()
	portal := portalForUserType(userType)
	now := time.Now()

	// ── Access Token ─────────────────────────────────────────────────────────
	accessClaims := InsureTechClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    s.config.JWT.Issuer,
			Audience:  jwt.ClaimStrings{s.config.JWT.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExpiresIn)),
			ID:        accessJTI,
		},
		UserType:  userType,
		Portal:    portal,
		TenantID:  tenantID,
		DeviceID:  deviceID,
		SessionID: sessionID,
		TokenType: "access",
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenObj.Header["kid"] = s.keyID
	accessToken, err := accessTokenObj.SignedString(s.rsaPrivateKey)
	if err != nil {
		logger.Errorf("failed to sign access token (RS256): %v", err)
		return nil, errors.New("failed to sign access token (RS256)")
	}

	// ── Refresh Token ─────────────────────────────────────────────────────────
	refreshClaims := InsureTechClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    s.config.JWT.Issuer,
			Audience:  jwt.ClaimStrings{s.config.JWT.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpiresIn)),
			ID:        refreshJTI,
		},
		UserType:  userType,
		Portal:    portal,
		TenantID:  tenantID,
		DeviceID:  deviceID,
		SessionID: sessionID,
		TokenType: "refresh",
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenObj.Header["kid"] = s.keyID
	refreshToken, err := refreshTokenObj.SignedString(s.rsaPrivateKey)
	if err != nil {
		logger.Errorf("failed to sign refresh token (RS256): %v", err)
		return nil, errors.New("failed to sign refresh token (RS256)")
	}

	// ── Persist Session ────────────────────────────────────────────────────────
	session := &authnentityv1.Session{
		SessionId:             sessionID,
		UserId:                userID,
		DeviceId:              deviceID,
		DeviceType:            deviceType,
		SessionType:           authnentityv1.SessionType_SESSION_TYPE_JWT,
		IpAddress:             ipAddress,
		UserAgent:             userAgent,
		IsActive:              true,
		AccessTokenJti:        accessJTI,
		RefreshTokenJti:       refreshJTI,
		ExpiresAt:             timestamppb.New(now.Add(refreshExpiresIn)),
		AccessTokenExpiresAt:  timestamppb.New(now.Add(accessExpiresIn)),
		RefreshTokenExpiresAt: timestamppb.New(now.Add(refreshExpiresIn)),
		LastActivityAt:        timestamppb.New(now),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		logger.Errorf("failed to create session: %v", err)
		return nil, errors.New("failed to create session")
	}

	// Enforce concurrent session limit (evict oldest sessions if over limit).
	if s.sessionLimiter != nil {
		evicted, err := s.sessionLimiter.TrackSession(ctx, userID, sessionID, now.Add(refreshExpiresIn))
		if err == nil {
			for _, evictedID := range evicted {
				_ = s.RevokeSession(ctx, evictedID)
			}
		}
	}

	return &TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  accessExpiresIn,
		RefreshTokenExpiresIn: refreshExpiresIn,
		SessionID:             sessionID,
	}, nil
}

// parseRS256Token parses and validates an RS256-signed JWT using the loaded public key.
func (s *TokenService) parseRS256Token(tokenString string) (*jwt.Token, *InsureTechClaims, error) {
	claims := &InsureTechClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method — RS256 required")
		}
		return s.rsaPublicKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return token, claims, nil
}

// ValidateServerSideSession validates a server-side session token.
// Flow: lookup by sha256 hash → bcrypt verify → expiry check → slide last_activity_at.
func (s *TokenService) ValidateServerSideSession(ctx context.Context, sessionToken string) (*authnservicev1.ValidateTokenResponse, error) {
	if sessionToken == "" {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}

	lookup := sessionTokenLookup(sessionToken)
	session, err := s.sessionRepo.GetByTokenLookup(ctx, lookup)
	if err != nil {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}

	if !session.IsActive || time.Now().After(session.ExpiresAt.AsTime()) {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}

	if session.SessionTokenHash == "" {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(session.SessionTokenHash), []byte(sessionToken)); err != nil {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}

	// Sprint 3: Idle timeout via Redis sliding window
	if s.rdb != nil && s.config.Security.IdleTimeoutDuration > 0 {
		idleKey := "session:idle:" + session.SessionId
		exists, _ := s.rdb.Exists(ctx, idleKey).Result()
		if exists == 0 {
			// Key missing means idle window expired — revoke session
			_ = s.sessionRepo.Revoke(ctx, session.SessionId)
			return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
		}
		// Slide the window
		s.rdb.Expire(ctx, idleKey, s.config.Security.IdleTimeoutDuration)
	} else if s.config.Security.IdleTimeoutDuration > 0 && session.LastActivityAt != nil {
		// Fallback: DB-based idle check
		if time.Since(session.LastActivityAt.AsTime()) > s.config.Security.IdleTimeoutDuration {
			_ = s.sessionRepo.Revoke(ctx, session.SessionId)
			return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
		}
	}

	_ = s.sessionRepo.UpdateLastActivity(ctx, session.SessionId)

	resp := &authnservicev1.ValidateTokenResponse{
		Valid:       true,
		UserId:      session.UserId,
		SessionId:   session.SessionId,
		SessionType: "SERVER_SIDE",
		ExpiresAt:   session.ExpiresAt,
	}

	// Populate portal/tenant/userType from user record
	if s.userRepo != nil {
		if user, err := s.userRepo.GetByID(ctx, session.UserId); err == nil && user != nil {
			resp.UserType = user.UserType.String()
			resp.Portal = portalForUserType(user.UserType.String())
			// TenantID populated from user's tenant field if present
		}
	}

	return resp, nil
}

// ValidateCSRFToken validates a CSRF token for a server-side session.
func (s *TokenService) ValidateCSRFToken(ctx context.Context, sessionID, csrfToken string) (bool, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return false, err
	}
	if !session.IsActive || time.Now().After(session.ExpiresAt.AsTime()) {
		return false, errors.New("session expired or inactive")
	}
	if session.CsrfToken != csrfToken {
		exp := sha256.Sum256([]byte(session.CsrfToken))
		rec := sha256.Sum256([]byte(csrfToken))
		reqMeta := &middleware.RequestMetadata{}
		if s.metadata != nil {
			reqMeta = s.metadata.ExtractAll(ctx)
		}
		if s.eventPublisher != nil {
			_ = s.eventPublisher.PublishCSRFValidationFailed(ctx, session.UserId, sessionID,
				hex.EncodeToString(exp[:]), hex.EncodeToString(rec[:]),
				reqMeta.IPAddress, reqMeta.UserAgent, "", "")
		}
		return false, nil
	}
	return true, nil
}

// ValidateJWT validates an RS256 JWT access token.
// Performs: RS256 signature check + expiry + JTI blocklist check + revocation (session.is_active).
// Returns full ValidateTokenResponse with portal/tenant/device_id/token_id fields.
func (s *TokenService) ValidateJWT(ctx context.Context, tokenString string) (*authnservicev1.ValidateTokenResponse, error) {
	_, claims, err := s.parseRS256Token(tokenString)
	if err != nil {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}

	// Device binding check: if caller sent x-device-id metadata, it must match token claim.
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		mdDeviceIDs := md.Get("x-device-id")
		if len(mdDeviceIDs) > 0 {
			reqDeviceID := mdDeviceIDs[0]
			if reqDeviceID != "" && claims.DeviceID != "" && reqDeviceID != claims.DeviceID {
				return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
			}
		}
	}

	// JTI blocklist check (Redis)
	if s.rdb != nil && s.isJTIBlocked(ctx, claims.ID) {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}

	// Revocation check: verify session is still active in DB
	if claims.SessionID != "" && s.sessionRepo != nil {
		session, err := s.sessionRepo.GetByID(ctx, claims.SessionID)
		if err != nil || session == nil || !session.IsActive {
			return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
		}
	}

	return &authnservicev1.ValidateTokenResponse{
		Valid:       true,
		UserId:      claims.Subject,
		UserType:    claims.UserType,
		Portal:      claims.Portal,
		TenantId:    claims.TenantID,
		DeviceId:    claims.DeviceID,
		TokenId:     claims.ID, // jti
		SessionId:   claims.SessionID,
		SessionType: "JWT",
		ExpiresAt:   timestamppb.New(claims.ExpiresAt.Time),
	}, nil
}

// ValidateJWTStrict validates RS256 JWT and additionally checks JTI against the session DB.
// Use for sensitive operations (password change, withdrawal, admin actions).
func (s *TokenService) ValidateJWTStrict(ctx context.Context, tokenString string) (*authnservicev1.ValidateTokenResponse, error) {
	resp, err := s.ValidateJWT(ctx, tokenString)
	if err != nil || !resp.Valid {
		return resp, err
	}

	// Additionally verify JTI matches stored access_token_jti
	_, claims, _ := s.parseRS256Token(tokenString)
	session, err := s.sessionRepo.GetByID(ctx, claims.SessionID)
	if err != nil {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}
	if session.AccessTokenJti != claims.ID {
		return &authnservicev1.ValidateTokenResponse{Valid: false}, nil
	}

	return resp, nil
}

// RefreshJWT validates an RS256 refresh token and issues a new token pair with rotation.
// Revokes old session and creates a new one (prevents refresh token reuse).
func (s *TokenService) RefreshJWT(ctx context.Context, refreshTokenString string) (*authnservicev1.RefreshTokenResponse, error) {
	_, claims, err := s.parseRS256Token(refreshTokenString)
	if err != nil {
		logger.Errorf("invalid refresh token: %v", err)
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("token is not a refresh token")
	}

	session, err := s.sessionRepo.GetByID(ctx, claims.SessionID)
	if err != nil {
		logger.Errorf("session not found: %v", err)
		return nil, errors.New("session not found")
	}
	if !session.IsActive || time.Now().After(session.ExpiresAt.AsTime()) {
		return nil, errors.New("session expired or revoked")
	}
	if session.RefreshTokenJti != claims.ID {
		return nil, errors.New("refresh token JTI mismatch — possible reuse attack")
	}
	if session.RefreshTokenExpiresAt != nil && time.Now().After(session.RefreshTokenExpiresAt.AsTime()) {
		return nil, errors.New("refresh token expired")
	}

	// Enforce per-user refresh rate limit
	if s.refreshLimiter != nil && !s.refreshLimiter.Allow(session.UserId) {
		return nil, errors.New("refresh token rate limit exceeded")
	}

	// Revoke old session (token rotation)
	if err := s.sessionRepo.Revoke(ctx, session.SessionId); err != nil {
		logger.Errorf("failed to revoke old session: %v", err)
		return nil, errors.New("failed to revoke old session")
	}

	// Resolve latest userType from DB (may have changed since last login)
	userType := claims.UserType
	tenantID := claims.TenantID
	if s.userRepo != nil {
		if u, err := s.userRepo.GetByID(ctx, session.UserId); err == nil && u != nil {
			userType = u.UserType.String()
		}
	}

	newPair, err := s.GenerateJWT(ctx, session.UserId, userType, tenantID, session.DeviceId, session.DeviceType, session.IpAddress, session.UserAgent)
	if err != nil {
		logger.Errorf("failed to generate new tokens: %v", err)
		return nil, errors.New("failed to generate new tokens")
	}

	return &authnservicev1.RefreshTokenResponse{
		AccessToken:           newPair.AccessToken,
		RefreshToken:          newPair.RefreshToken,
		AccessTokenExpiresIn:  int32(newPair.AccessTokenExpiresIn.Seconds()),
		RefreshTokenExpiresIn: int32(newPair.RefreshTokenExpiresIn.Seconds()),
		SessionId:             newPair.SessionID,
	}, nil
}

// RevokeSession marks a session as inactive (soft revoke) and blocks its JTIs in Redis.
func (s *TokenService) RevokeSession(ctx context.Context, sessionID string) error {
	// Read the session first to get JTIs (before revoking)
	if s.rdb != nil && s.sessionRepo != nil {
		if session, err := s.sessionRepo.GetByID(ctx, sessionID); err == nil && session != nil {
			now := time.Now()
			// Block access token JTI
			if session.AccessTokenJti != "" {
				var accessTTL time.Duration
				if session.AccessTokenExpiresAt != nil {
					accessTTL = session.AccessTokenExpiresAt.AsTime().Sub(now)
				}
				if accessTTL <= 0 {
					accessTTL = s.config.JWT.AccessTokenDuration
					if accessTTL == 0 {
						accessTTL = 15 * time.Minute
					}
				}
				_ = s.BlockJTI(ctx, session.AccessTokenJti, accessTTL)
			}
			// Block refresh token JTI
			if session.RefreshTokenJti != "" {
				var refreshTTL time.Duration
				if session.RefreshTokenExpiresAt != nil {
					refreshTTL = session.RefreshTokenExpiresAt.AsTime().Sub(now)
				}
				if refreshTTL <= 0 {
					refreshTTL = s.config.JWT.RefreshTokenDuration
					if refreshTTL == 0 {
						refreshTTL = 7 * 24 * time.Hour
					}
				}
				_ = s.BlockJTI(ctx, session.RefreshTokenJti, refreshTTL)
			}
		}
	}
	return s.sessionRepo.Revoke(ctx, sessionID)
}

// GetJWKS builds a JWKS response from the loaded RSA public key.
func (s *TokenService) GetJWKS(ctx context.Context, req *authnservicev1.GetJWKSRequest) (*authnservicev1.GetJWKSResponse, error) {
	if s.rsaPublicKey == nil {
		return &authnservicev1.GetJWKSResponse{Keys: []*authnservicev1.JWK{}}, nil
	}

	// n = base64url(pubKey.N.Bytes())
	nBytes := s.rsaPublicKey.N.Bytes()
	nEncoded := base64.RawURLEncoding.EncodeToString(nBytes)

	// e = base64url(big-endian bytes of pubKey.E, trimmed leading zeros)
	eBig := new(big.Int).SetInt64(int64(s.rsaPublicKey.E))
	eBytes := eBig.Bytes()
	eEncoded := base64.RawURLEncoding.EncodeToString(eBytes)

	jwk := &authnservicev1.JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: s.keyID,
		N:   nEncoded,
		E:   eEncoded,
	}

	return &authnservicev1.GetJWKSResponse{
		Keys: []*authnservicev1.JWK{jwk},
	}, nil
}

// generateSecureRandomString generates a cryptographically secure random hex string.
func generateSecureRandomString(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		logger.Errorf("failed to generate random bytes: %v", err)
		return "", errors.New("failed to generate random bytes")
	}
	return hex.EncodeToString(b), nil
}
