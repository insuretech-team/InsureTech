package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/consumers"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/pii"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/mobile"
	apikeyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/apikey/entity/v1"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthService is the production-grade authentication service
type AuthService struct {
	tokenService     *TokenService
	otpService       *OTPService
	userRepo         *repository.UserRepository
	sessionRepo      *repository.SessionRepository
	otpRepo          *repository.OTPRepository
	apiKeyRepo       *repository.ApiKeyRepository
	userProfileRepo  *repository.UserProfileRepository
	userDocumentRepo *repository.UserDocumentRepository
	documentTypeRepo *repository.DocumentTypeRepository
	kycRepo          *repository.KYCVerificationRepository
	externalKYC      ExternalKYCClient
	flveAdapter      FLVEAdapter
	voiceRepo        *repository.VoiceSessionRepository
	eventPublisher   *events.Publisher
	config           *config.Config
	metadata         *middleware.MetadataExtractor
}

// NewAuthService creates the authentication service with all repositories.
// Pass nil for optional repos (apiKeyRepo, userProfileRepo, userDocumentRepo,
// documentTypeRepo) if those RPCs are not needed in a given deployment.
func NewAuthService(
	tokenService *TokenService,
	otpService *OTPService,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	otpRepo *repository.OTPRepository,
	apiKeyRepo *repository.ApiKeyRepository,
	userProfileRepo *repository.UserProfileRepository,
	userDocumentRepo *repository.UserDocumentRepository,
	documentTypeRepo *repository.DocumentTypeRepository,
	kycRepo *repository.KYCVerificationRepository,
	voiceRepo *repository.VoiceSessionRepository,
	eventPublisher *events.Publisher,
	cfg *config.Config,
	metadata *middleware.MetadataExtractor,
) *AuthService {
	return &AuthService{
		tokenService:     tokenService,
		otpService:       otpService,
		userRepo:         userRepo,
		sessionRepo:      sessionRepo,
		otpRepo:          otpRepo,
		apiKeyRepo:       apiKeyRepo,
		userProfileRepo:  userProfileRepo,
		userDocumentRepo: userDocumentRepo,
		documentTypeRepo: documentTypeRepo,
		kycRepo:          kycRepo,
		voiceRepo:        voiceRepo,
		eventPublisher:   eventPublisher,
		config:           cfg,
		metadata:         metadata,
	}
}

// NewAuthServiceWithAPIKey is kept as a convenience alias for callers that do
// not yet pass the profile/document repos. Deprecated: prefer NewAuthService.
func NewAuthServiceWithAPIKey(
	tokenService *TokenService,
	otpService *OTPService,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	otpRepo *repository.OTPRepository,
	apiKeyRepo *repository.ApiKeyRepository,
	eventPublisher *events.Publisher,
	cfg *config.Config,
	metadata *middleware.MetadataExtractor,
) *AuthService {
	return NewAuthService(tokenService, otpService, userRepo, sessionRepo, otpRepo,
		apiKeyRepo, nil, nil, nil, nil, nil, eventPublisher, cfg, metadata)
}

// Login handles user authentication with hybrid session management
func (s *AuthService) Login(ctx context.Context, req *authnservicev1.LoginRequest) (*authnservicev1.LoginResponse, error) {
	// Extract metadata from context
	reqMeta := s.metadata.ExtractAll(ctx)

	// Login-specific rate limit (RATE_LIMIT_LOGIN_PER_MINUTE env, default 5).
	// Checked before any DB lookup to protect against credential stuffing.
	// This is separate from the OTP rate limit (RATE_LIMIT_PER_MINUTE).
	if limit := s.config.Security.RateLimitLoginPerMinute; limit > 0 {
		recentLogins, countErr := s.sessionRepo.CountRecentLoginAttempts(ctx, req.MobileNumber, 1*time.Minute)
		if countErr == nil && recentLogins >= int64(limit) {
			appLogger.Warnf("Login rate limited: %s made %d attempts in last minute from IP %s",
				req.MobileNumber, recentLogins, reqMeta.IPAddress)
			return nil, status.Errorf(codes.ResourceExhausted,
				"too many login attempts — please wait before trying again")
		}
	}
	deviceID := strings.TrimSpace(req.DeviceId)
	if deviceID == "" {
		// Stable fallback fingerprint when client does not send a device id.
		fp := sha256.Sum256([]byte(strings.TrimSpace(reqMeta.UserAgent) + "|" + strings.TrimSpace(reqMeta.IPAddress)))
		deviceID = "fp_" + hex.EncodeToString(fp[:16])
	}

	// 1. Verify Credentials
	mobileToLookup, err := normalizeMobileForLookup(req.MobileNumber)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	user, err := s.userRepo.GetByMobileNumber(ctx, mobileToLookup)

	// B2C OTP-only login: if Password is empty and user does NOT exist yet,
	// check for a recently verified OTP for this mobile number first, then
	// auto-create the user as B2C_CUSTOMER (passwordless registration on first login).
	// This is the mobile-first B2C pattern: OTP verify → auto-register → JWT.
	otpLoginVerified := false
	if err != nil && req.Password == "" {
		// User not found — check for a verified OTP before deciding to create
		recentOTP, otpErr := s.otpRepo.GetRecentlyVerifiedForMobile(ctx, mobileToLookup, 5*time.Minute)
		if otpErr != nil || recentOTP == nil {
			appLogger.Warnf("B2C login: user not found and no verified OTP for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
			_ = s.eventPublisher.PublishLoginFailed(ctx, "", req.MobileNumber, "otp_required", reqMeta.IPAddress, req.DeviceType, reqMeta.UserAgent, 0)
			return nil, errors.New("invalid credentials: please verify OTP before login")
		}
		// Verified OTP found — auto-create B2C user (passwordless first-time login)
		appLogger.Infof("B2C auto-registration: creating new B2C_CUSTOMER for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
		newUser, createErr := s.userRepo.Create(ctx, mobileToLookup, "", "", authnentityv1.UserStatus_USER_STATUS_ACTIVE)
		if createErr != nil {
			appLogger.Errorf("B2C auto-registration failed for mobile %s: %v", req.MobileNumber, createErr)
			return nil, errors.New("failed to create user account")
		}
		// Set user_type to B2C_CUSTOMER explicitly
		newUser.UserType = authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER
		user = newUser
		otpLoginVerified = true
		_ = s.eventPublisher.PublishUserRegistered(ctx, user.UserId, mobileToLookup, "", reqMeta.IPAddress, "b2c_otp_auto_register", "b2c", "root")
		appLogger.Infof("B2C auto-registration success: user_id=%s mobile=%s otp_id=%s from IP %s", user.UserId, req.MobileNumber, recentOTP.OtpId, reqMeta.IPAddress)
	} else if err != nil && req.Password != "" && req.DeviceId != "" && isMobileDeviceType(req.DeviceType) {
		// User not found but a device credential was provided — check if credential matches
		// (try both with and without '+' prefix since OTP recipient may store without '+')
		mobileNoPlus := strings.TrimPrefix(mobileToLookup, "+")
		credMatchPlus := deviceCredentialMatches(mobileToLookup, req.DeviceId, req.Password)
		credMatchRaw := deviceCredentialMatches(mobileNoPlus, req.DeviceId, req.Password)
		if credMatchPlus || credMatchRaw {
			// Valid device credential but no user yet — auto-create B2C user
			appLogger.Infof("B2C device-credential auto-registration: creating new B2C_CUSTOMER for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
			newUser, createErr := s.userRepo.Create(ctx, mobileToLookup, "", "", authnentityv1.UserStatus_USER_STATUS_ACTIVE)
			if createErr != nil {
				appLogger.Errorf("B2C device-credential auto-registration failed for mobile %s: %v", req.MobileNumber, createErr)
				return nil, errors.New("failed to create user account")
			}
			newUser.UserType = authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER
			user = newUser
			otpLoginVerified = true
			s.markTrustedDevice(ctx, user.UserId, req.DeviceId)
			_ = s.eventPublisher.PublishUserRegistered(ctx, user.UserId, mobileToLookup, "", reqMeta.IPAddress, "b2c_device_credential_auto_register", "b2c", "root")
			appLogger.Infof("B2C device-credential auto-registration success: user_id=%s mobile=%s from IP %s", user.UserId, req.MobileNumber, reqMeta.IPAddress)
		} else {
			appLogger.Warnf("Login failed: user not found and device credential mismatch for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
			_ = s.eventPublisher.PublishLoginFailed(ctx, "", req.MobileNumber, "user_not_found", reqMeta.IPAddress, req.DeviceType, reqMeta.UserAgent, 0)
			return nil, errors.New("invalid credentials")
		}
	} else if err != nil {
		// User not found and no device credential — standard error
		appLogger.Warnf("Login failed: user not found for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
		_ = s.eventPublisher.PublishLoginFailed(ctx, "", req.MobileNumber, "user_not_found", reqMeta.IPAddress, req.DeviceType, reqMeta.UserAgent, 0)
		return nil, errors.New("invalid credentials")
	}

	// Sprint 5: Account lockout check — deny early if still locked
	if user.LockedUntil != nil && time.Now().Before(user.LockedUntil.AsTime()) {
		remaining := time.Until(user.LockedUntil.AsTime()).Round(time.Second)
		appLogger.Warnf("Login blocked: account locked for user %s until %s from IP %s", user.UserId, user.LockedUntil.AsTime(), reqMeta.IPAddress)
		_ = s.eventPublisher.PublishLoginFailed(ctx, user.UserId, req.MobileNumber, "account_locked", reqMeta.IPAddress, req.DeviceType, reqMeta.UserAgent, user.LoginAttempts)
		return nil, fmt.Errorf("account is locked. Try again in %s", remaining)
	}

	// WhatsApp-style device credential check: if mobile device type and non-empty password,
	// try to match it as a device credential (HMAC of mobile+device_id) BEFORE
	// falling through to bcrypt password check. This enables seamless re-login
	// on the same device without OTP.
	if req.Password != "" && req.DeviceId != "" && isMobileDeviceType(req.DeviceType) {
		mobileNoPlus := strings.TrimPrefix(mobileToLookup, "+")
		if deviceCredentialMatches(mobileToLookup, req.DeviceId, req.Password) ||
			deviceCredentialMatches(mobileNoPlus, req.DeviceId, req.Password) {
			otpLoginVerified = true
			// Refresh trusted device TTL in Redis
			s.markTrustedDevice(ctx, user.UserId, req.DeviceId)
			appLogger.Infof("Device credential login: user=%s mobile=%s device_id=%s device_type=%s from IP %s",
				user.UserId, req.MobileNumber, req.DeviceId, req.DeviceType, reqMeta.IPAddress)
		}
	}

	// B2C OTP-only login: if Password is empty and user already exists,
	// check for a recently verified OTP for this mobile number.
	if req.Password == "" && !otpLoginVerified {
		recentOTP, otpErr := s.otpRepo.GetRecentlyVerifiedForMobile(ctx, user.MobileNumber, 5*time.Minute)
		if otpErr == nil && recentOTP != nil {
			otpLoginVerified = true
			// Consume the OTP immediately so it cannot be reused by a subsequent
			// login attempt (prevents replay/session-hijack via OTP reuse).
			if expErr := s.otpRepo.ExpireOTP(ctx, recentOTP.OtpId); expErr != nil {
				appLogger.Warnf("B2C OTP login: failed to expire OTP %s after use: %v", recentOTP.OtpId, expErr)
			}
			appLogger.Infof("B2C OTP login: verified OTP %s consumed for existing user %s mobile %s from IP %s", recentOTP.OtpId, user.UserId, req.MobileNumber, reqMeta.IPAddress)
		} else {
			appLogger.Warnf("Login failed: no password and no verified OTP for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
			_ = s.eventPublisher.PublishLoginFailed(ctx, user.UserId, req.MobileNumber, "otp_required", reqMeta.IPAddress, req.DeviceType, reqMeta.UserAgent, 0)
			return nil, errors.New("invalid credentials: password or verified OTP required")
		}
	}

	if !otpLoginVerified {
		valid, needsRehash, verifyErr := verifyPassword(req.Password, user.PasswordHash)
		if verifyErr != nil || !valid {
			appLogger.Warnf("Login failed: invalid password for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
			// Sprint 5: Increment failed attempts and lock if threshold exceeded
			maxLoginAttempts := s.config.Security.LoginMaxAttempts
			if maxLoginAttempts <= 0 {
				maxLoginAttempts = 8
			}
			lockoutDuration := s.config.Security.LoginLockoutDuration
			if lockoutDuration <= 0 {
				lockoutDuration = 10 * time.Minute
			}
			attempts, _ := s.userRepo.IncrementLoginAttempts(ctx, user.UserId)
			if int(attempts) >= maxLoginAttempts {
				_ = s.userRepo.LockAccount(ctx, user.UserId, lockoutDuration)
				appLogger.Warnf("Account locked for user %s after %d failed attempts from IP %s", user.UserId, attempts, reqMeta.IPAddress)
			}
			_ = s.eventPublisher.PublishLoginFailed(ctx, user.UserId, req.MobileNumber, "invalid_password", reqMeta.IPAddress, req.DeviceType, reqMeta.UserAgent, attempts)
			return nil, errors.New("invalid credentials")
		}
		if needsRehash {
			if newHash, err := hashPassword(req.Password); err == nil {
				_ = s.userRepo.UpdatePassword(ctx, user.UserId, newHash)
			}
		}
	}
	// Sprint 5: Reset failed login attempts on successful credential verification
	_ = s.userRepo.ResetLoginAttempts(ctx, user.UserId)

	// Sprint 2: Per-portal MFA enforcement via GlobalPortalConfigCache
	// The cache is populated by the authz.events Kafka consumer (NewPortalConfigUpdatedHandler).
	portalKey := portalConfigKeyForUserType(user.UserType.String())
	if portalKey != "" {
		if portalCfg := consumers.GlobalPortalConfigCache.Get(portalKey); portalCfg != nil && portalCfg.MfaRequired {
			if s.isTrustedDevice(ctx, user.UserId, deviceID) {
				goto issueTokens
			}
			if !user.TotpEnabled {
				// MFA required but TOTP not configured — prompt setup
				return &authnservicev1.LoginResponse{
					UserId:      user.UserId,
					MfaRequired: true,
					MfaMethod:   "TOTP",
				}, nil
			}
			// MFA required and TOTP is configured — gating: return mfa_required=true.
			// Sprint 2.2: Store short-lived MFA session token (Redis, TTL=5m) so that
			// VerifyTOTP can issue real session tokens after TOTP verification.
			mfaToken, mfaErr := s.StoreMFASessionToken(ctx, user.UserId, deviceID, req.DeviceType, reqMeta.IPAddress)
			if mfaErr != nil {
				appLogger.Warnf("StoreMFASessionToken failed for user %s: %v (continuing without token)", user.UserId, mfaErr)
			}
			return &authnservicev1.LoginResponse{
				UserId:          user.UserId,
				MfaRequired:     true,
				MfaMethod:       "TOTP",
				MfaSessionToken: mfaToken,
			}, nil
		}
	}

issueTokens:
	// 2. Parse and validate device type
	deviceType := parseDeviceType(req.DeviceType)

	// Reject unknown/unspecified device_type — clients must send an explicit valid value.
	// This prevents garbage values like the Postman placeholder "string" from silently
	// falling through to JWT. Valid values: WEB, ANDROID, IOS, API, DESKTOP.
	if deviceType == authnentityv1.DeviceType_DEVICE_TYPE_UNSPECIFIED {
		appLogger.Warnf("Login rejected: unknown device_type %q for mobile %s from IP %s",
			req.DeviceType, req.MobileNumber, reqMeta.IPAddress)
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid device_type %q: must be one of WEB, ANDROID, IOS, API, DESKTOP", req.DeviceType)
	}

	// Web-portal-only user types must always use device_type=WEB → SERVER_SIDE session.
	// Issuing JWT to these user types would be a security misconfiguration.
	//
	// JWT-eligible (mobile/app users):  B2C_CUSTOMER, AGENT
	// SERVER_SIDE only (web portal):    SYSTEM_USER, BUSINESS_BENEFICIARY, PARTNER,
	//                                   REGULATOR, BUSINESS_ADMIN, B2B_ORG_ADMIN
	isWebOnlyUserType := func(ut authnentityv1.UserType) bool {
		switch ut {
		case authnentityv1.UserType_USER_TYPE_SYSTEM_USER,
			authnentityv1.UserType_USER_TYPE_BUSINESS_BENEFICIARY,
			authnentityv1.UserType_USER_TYPE_PARTNER,
			authnentityv1.UserType_USER_TYPE_REGULATOR,
			authnentityv1.UserType_USER_TYPE_BUSINESS_ADMIN,
			authnentityv1.UserType_USER_TYPE_B2B_ORG_ADMIN,
			authnentityv1.UserType_USER_TYPE_B2B_BENEFICIARY:
			return true
		default:
			return false
		}
	}
	if isWebOnlyUserType(user.UserType) && deviceType != authnentityv1.DeviceType_DEVICE_TYPE_WEB {
		appLogger.Warnf("Login rejected: web-portal user %s (type=%s) attempted non-WEB login with device_type=%s from IP %s",
			user.UserId, user.UserType, req.DeviceType, reqMeta.IPAddress)
		return nil, status.Errorf(codes.InvalidArgument,
			"user type %s must use device_type=WEB for login (web portal users cannot use JWT)",
			user.UserType.String())
	}

	sessionType := mapDeviceTypeToSessionType(deviceType)

	// 3. Generate session/tokens based on device type
	resp := &authnservicev1.LoginResponse{
		UserId:                 user.UserId,
		User:                   user,
		PasswordChangeRequired: user.PasswordChangeRequired,
	}

	var sessionID string

	// Hybrid Auth Logic: Web -> Server-Side Session; Mobile/App -> JWT
	if sessionType == authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE {
		// Server-side session for web
		serverSession, err := s.tokenService.GenerateServerSideSession(
			ctx,
			user.UserId,
			deviceID,
			deviceType,
			reqMeta.IPAddress,
			reqMeta.UserAgent,
		)
		if err != nil {
			appLogger.Errorf("Failed to create server-side session: %v", err)
			logger.Errorf("failed to create session: %v", err)
			return nil, errors.New("failed to create session")
		}

		resp.SessionToken = serverSession.SessionToken // To be set as HttpOnly cookie
		resp.SessionId = serverSession.SessionID
		resp.CsrfToken = serverSession.CSRFToken
		resp.SessionType = "SERVER_SIDE"
		sessionID = serverSession.SessionID

		appLogger.Infof("Server-side session created for user %s: %s from IP %s", user.UserId, sessionID, reqMeta.IPAddress)

	} else {
		// JWT for mobile/API
		// B2C customers always use the shared 'b2c:root' domain.
		// tenantID must be non-empty so the JWT ins_tenant claim resolves to
		// "b2c:root" (not "b2c:") in CheckAccess domain format "portal:tenant".
		b2cTenantID := "root"
		if user.UserType != authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER &&
			user.UserType != authnentityv1.UserType_USER_TYPE_AGENT {
			b2cTenantID = "" // non-B2C portals resolve tenant from their org
		}
		tokens, err := s.tokenService.GenerateJWT(
			ctx,
			user.UserId,
			user.UserType.String(),
			b2cTenantID,
			deviceID,
			deviceType,
			reqMeta.IPAddress,
			reqMeta.UserAgent,
		)
		if err != nil {
			appLogger.Errorf("Failed to generate JWT tokens: %v", err)
			logger.Errorf("failed to generate tokens: %v", err)
			return nil, errors.New("failed to generate tokens")
		}

		resp.AccessToken = tokens.AccessToken
		resp.RefreshToken = tokens.RefreshToken
		resp.AccessTokenExpiresIn = int32(tokens.AccessTokenExpiresIn.Seconds())
		resp.RefreshTokenExpiresIn = int32(tokens.RefreshTokenExpiresIn.Seconds())
		resp.SessionId = tokens.SessionID
		resp.SessionType = "JWT"
		sessionID = tokens.SessionID

		appLogger.Infof("JWT session created for user %s: %s from IP %s", user.UserId, sessionID, reqMeta.IPAddress)
	}

	// 4. Publish Event
	_ = s.eventPublisher.PublishUserLoggedIn(ctx, user.UserId, sessionID, resp.SessionType, reqMeta.IPAddress, req.DeviceType, reqMeta.UserAgent)

	return resp, nil
}

// Register handles user registration
func (s *AuthService) Register(ctx context.Context, req *authnservicev1.RegisterRequest) (*authnservicev1.RegisterResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Check if user exists using the canonical +8801XXXXXXXXX form.
	mobileToLookup, err := normalizeMobileForLookup(req.MobileNumber)
	if err != nil {
		return nil, errors.New("invalid mobile number")
	}
	existing, _ := s.userRepo.GetByMobileNumber(ctx, mobileToLookup)
	if existing != nil {
		appLogger.Warnf("Registration failed: user already exists for mobile %s", req.MobileNumber)
		return nil, errors.New("user already exists")
	}

	// Validate password strength
	if err := validatePasswordStrength(req.Password); err != nil {
		logger.Errorf("weak password: %v", err)
		return nil, errors.New("weak password")
	}

	// Hash password with Argon2id
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		appLogger.Errorf("Failed to hash password: %v", err)
		logger.Errorf("failed to hash password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	// Create user with proper signature
	user, err := s.userRepo.Create(ctx, mobileToLookup, hashedPassword, req.Email, authnentityv1.UserStatus_USER_STATUS_ACTIVE)
	if err != nil {
		appLogger.Errorf("Failed to create user: %v", err)
		logger.Errorf("failed to create user: %v", err)
		return nil, errors.New("failed to create user")
	}

	// Publish event
	_ = s.eventPublisher.PublishUserRegistered(ctx, user.UserId, mobileToLookup, req.Email, reqMeta.IPAddress, "", "b2c", "root")

	appLogger.Infof("User registered successfully: %s from IP %s", user.UserId, reqMeta.IPAddress)

	return &authnservicev1.RegisterResponse{
		UserId:  user.UserId,
		Message: "Registration successful",
	}, nil
}

// ValidateToken validates either JWT or server-side session token
func (s *AuthService) ValidateToken(ctx context.Context, req *authnservicev1.ValidateTokenRequest) (*authnservicev1.ValidateTokenResponse, error) {
	// If access_token looks like JWT (starts with "eyJ"), validate as JWT
	if req.AccessToken != "" && strings.HasPrefix(req.AccessToken, "eyJ") {
		return s.tokenService.ValidateJWT(ctx, req.AccessToken)
	}

	// Otherwise validate as server-side session
	if req.SessionId != "" {
		return s.tokenService.ValidateServerSideSession(ctx, req.SessionId)
	}

	return &authnservicev1.ValidateTokenResponse{Valid: false}, errors.New("no valid token or session provided")
}

// RefreshToken handles token refresh for JWT sessions
func (s *AuthService) RefreshToken(ctx context.Context, req *authnservicev1.RefreshTokenRequest) (*authnservicev1.RefreshTokenResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	resp, err := s.tokenService.RefreshJWT(ctx, req.RefreshToken)
	if err != nil {
		appLogger.Warnf("Token refresh failed from IP %s: %v", reqMeta.IPAddress, err)
		return nil, err
	}

	// Best-effort event: token refreshed.
	// RefreshTokenResponse does not include user_id, so derive it from the new access token.
	if v, err := s.tokenService.ValidateJWT(ctx, resp.AccessToken); err == nil && v != nil && v.Valid {
		_ = s.eventPublisher.PublishTokenRefreshed(ctx, v.UserId, resp.SessionId, "", "", "", reqMeta.IPAddress, "", reqMeta.UserAgent)
	}

	appLogger.Infof("Token refreshed successfully: session %s from IP %s", resp.SessionId, reqMeta.IPAddress)
	return resp, nil
}

// Logout handles session termination
func (s *AuthService) Logout(ctx context.Context, req *authnservicev1.LogoutRequest) (*authnservicev1.LogoutResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)
	sessionID := req.SessionId

	// If access_token is provided (JWT), extract session_id from it
	if req.AccessToken != "" && sessionID == "" {
		resp, err := s.tokenService.ValidateJWT(ctx, req.AccessToken)
		if err == nil && resp.Valid {
			sessionID = resp.SessionId
		}
	}

	if sessionID == "" {
		return &authnservicev1.LogoutResponse{
			Message:        "No valid session to revoke",
			SessionRevoked: false,
		}, nil
	}

	if err := s.tokenService.RevokeSession(ctx, sessionID); err != nil {
		appLogger.Warnf("Logout: session revocation failed for %s from IP %s: %v", sessionID, reqMeta.IPAddress, err)
		return &authnservicev1.LogoutResponse{
			Message:        "Failed to revoke session",
			SessionRevoked: false,
		}, err
	}

	// Publish logout event (best-effort)
	logoutReason := req.LogoutReason
	if logoutReason == "" {
		logoutReason = "user_initiated"
	}
	userID := ""
	sessionType := ""
	// Derive user/session type
	if req.AccessToken != "" {
		if v, err := s.tokenService.ValidateJWT(ctx, req.AccessToken); err == nil && v != nil && v.Valid {
			userID = v.UserId
			sessionType = v.SessionType
		}
	}
	if userID == "" {
		if sess, err := s.sessionRepo.GetByID(ctx, sessionID); err == nil && sess != nil {
			userID = sess.UserId
			if sess.SessionType == authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE {
				sessionType = "SERVER_SIDE"
			} else {
				sessionType = "JWT"
			}
		}
	}
	_ = s.eventPublisher.PublishUserLoggedOut(ctx, userID, sessionID, sessionType, logoutReason, reqMeta.IPAddress, "")

	appLogger.Infof("Session revoked: %s from IP %s", sessionID, reqMeta.IPAddress)
	return &authnservicev1.LogoutResponse{
		Message:        "Successfully logged out",
		SessionRevoked: true,
	}, nil
}

// SendOTP handles OTP generation and sending
func (s *AuthService) SendOTP(ctx context.Context, req *authnservicev1.SendOTPRequest) (*authnservicev1.SendOTPResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	appLogger.Infof("Sending OTP to %s (type: %s) from IP %s", req.Recipient, req.Type, reqMeta.IPAddress)

	resp, err := s.otpService.SendOTP(ctx, req)
	if err != nil {
		appLogger.Errorf("Failed to send OTP to %s: %v", req.Recipient, err)
		return nil, err
	}

	// Publish OTP sent event (typed proto)
	// Mask recipient for safe publishing
	masked := req.Recipient
	if len(masked) > 6 {
		masked = masked[:3] + "***" + masked[len(masked)-3:]
	}
	provider := ""
	senderID := ""
	providerMessageID := ""
	maskingUsed := req.UseMasking
	channel := req.Channel
	if channel == "" {
		channel = "sms"
	}
	if channel == "sms" {
		provider = "sslwireless"
		// SenderID/providerMessageID are persisted in DB; we don't have them here without querying.
	} else {
		provider = "smtp"
	}
	_ = s.eventPublisher.PublishOTPSent(ctx, resp.OtpId, masked, req.Type, channel, provider, senderID, providerMessageID, maskingUsed)

	return resp, nil
}

// VerifyOTP handles OTP verification
func (s *AuthService) VerifyOTP(ctx context.Context, req *authnservicev1.VerifyOTPRequest) (*authnservicev1.VerifyOTPResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	appLogger.Infof("Verifying OTP %s from IP %s", req.OtpId, reqMeta.IPAddress)

	resp, err := s.otpService.VerifyOTP(ctx, req)
	if err != nil {
		appLogger.Errorf("Failed to verify OTP %s: %v", req.OtpId, err)
		return nil, err
	}

	if resp.Verified {
		appLogger.Infof("OTP %s verified successfully from IP %s", req.OtpId, reqMeta.IPAddress)
		// Fetch OTP to capture attempts count for event.
		attempts := int32(0)
		var otpMobile string
		if otp, err := s.otpRepo.GetByID(ctx, req.OtpId); err == nil && otp != nil {
			attempts = otp.Attempts
			otpMobile = otp.Recipient // phone number stored as recipient on OTP record
		}
		_ = s.eventPublisher.PublishOTPVerified(ctx, req.OtpId, resp.UserId, attempts)

		// WhatsApp-style device binding: if mobile device provided device_id,
		// derive and return a one-time device credential. The app stores this
		// in Android Keystore / iOS Keychain and uses it as `password` on
		// subsequent logins — no OTP required on same device.
		if req.DeviceId != "" && isMobileDeviceType(req.DeviceType) && otpMobile != "" {
			cred := deriveDeviceCredential(otpMobile, req.DeviceId)
			resp.DeviceCredential = cred
			// Mark this device as trusted in Redis (30-day TTL)
			if resp.UserId != "" {
				s.markTrustedDevice(ctx, resp.UserId, req.DeviceId)
			}
			appLogger.Infof("Device credential issued for mobile=%s device_id=%s device_type=%s otp_id=%s",
				otpMobile, req.DeviceId, req.DeviceType, req.OtpId)
		}
	} else {
		appLogger.Warnf("OTP %s verification failed from IP %s: %s", req.OtpId, reqMeta.IPAddress, resp.Message)
	}

	return resp, nil
}

// ChangePassword handles password change
func (s *AuthService) ChangePassword(ctx context.Context, req *authnservicev1.ChangePasswordRequest) (*authnservicev1.ChangePasswordResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Get user
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify current password (proto field is old_password → OldPassword)
	valid, _, verifyErr := verifyPassword(req.OldPassword, user.PasswordHash)
	if verifyErr != nil || !valid {
		appLogger.Warnf("Password change failed for user %s: invalid current password from IP %s", req.UserId, reqMeta.IPAddress)
		return nil, errors.New("invalid current password")
	}

	// Validate new password strength
	if err := validatePasswordStrength(req.NewPassword); err != nil {
		logger.Errorf("weak password: %v", err)
		return nil, errors.New("weak password")
	}

	// Hash new password with Argon2id
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf("failed to hash password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	// Update password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.UpdatePassword(ctx, req.UserId, hashedPassword); err != nil {
		logger.Errorf("failed to update password: %v", err)
		return nil, errors.New("failed to update password")
	}
	if err := s.userRepo.SetPasswordChangeRequired(ctx, req.UserId, false); err != nil {
		logger.Errorf("failed to clear password_change_required: %v", err)
		return nil, errors.New("failed to update password")
	}

	// Revoke all sessions (force re-login)
	_ = s.sessionRepo.RevokeAllByUserID(ctx, req.UserId, "")

	_ = s.eventPublisher.PublishPasswordChanged(ctx, req.UserId, reqMeta.IPAddress, req.UserId)
	appLogger.Infof("Password changed successfully for user %s from IP %s", req.UserId, reqMeta.IPAddress)

	return &authnservicev1.ChangePasswordResponse{
		Message: "Password changed successfully. Please login again.",
	}, nil
}

// ResetPassword handles password reset (with OTP)
func (s *AuthService) ResetPassword(ctx context.Context, req *authnservicev1.ResetPasswordRequest) (*authnservicev1.ResetPasswordResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Get user by mobile number first using the canonical +8801XXXXXXXXX form.
	mobileToLookup, err := normalizeMobileForLookup(req.MobileNumber)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	user, err := s.userRepo.GetByMobileNumber(ctx, mobileToLookup)
	if err != nil {
		appLogger.Warnf("Password reset failed: user not found for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
		return nil, errors.New("invalid credentials")
	}

	// Get latest OTP for this mobile number for reset_password purpose
	otp, err := s.otpRepo.GetLastOTP(ctx, req.MobileNumber)
	if err != nil || otp == nil {
		appLogger.Warnf("Password reset failed: no OTP found for mobile %s from IP %s", req.MobileNumber, reqMeta.IPAddress)
		return nil, errors.New("invalid or expired OTP")
	}

	// Verify OTP code
	verifyResp, err := s.otpService.VerifyOTP(ctx, &authnservicev1.VerifyOTPRequest{
		OtpId: otp.OtpId,
		Code:  req.OtpCode,
	})
	if err != nil || !verifyResp.Verified {
		appLogger.Warnf("Password reset failed: invalid OTP from IP %s", reqMeta.IPAddress)
		return nil, errors.New("invalid or expired OTP")
	}

	// Validate new password strength
	if err := validatePasswordStrength(req.NewPassword); err != nil {
		logger.Errorf("weak password: %v", err)
		return nil, errors.New("weak password")
	}

	// Hash new password with Argon2id
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf("failed to hash password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, user.UserId, hashedPassword); err != nil {
		logger.Errorf("failed to update password: %v", err)
		return nil, errors.New("failed to update password")
	}
	if err := s.userRepo.SetPasswordChangeRequired(ctx, user.UserId, false); err != nil {
		logger.Errorf("failed to clear password_change_required: %v", err)
		return nil, errors.New("failed to update password")
	}

	// Revoke all sessions
	_ = s.sessionRepo.RevokeAllByUserID(ctx, user.UserId, "")

	_ = s.eventPublisher.PublishPasswordResetRequested(ctx, user.UserId, req.MobileNumber, reqMeta.IPAddress, "")
	_ = s.eventPublisher.PublishPasswordChanged(ctx, user.UserId, reqMeta.IPAddress, user.UserId)
	appLogger.Infof("Password reset successfully for user %s from IP %s", user.UserId, reqMeta.IPAddress)

	return &authnservicev1.ResetPasswordResponse{
		Message: "Password reset successfully. Please login with your new password.",
	}, nil
}

// GetSession retrieves a session by ID
func (s *AuthService) GetSession(ctx context.Context, req *authnservicev1.GetSessionRequest) (*authnservicev1.GetSessionResponse, error) {
	session, err := s.sessionRepo.GetByID(ctx, req.SessionId)
	if err != nil {
		logger.Errorf("session not found: %v", err)
		return nil, errors.New("session not found")
	}
	return &authnservicev1.GetSessionResponse{Session: session}, nil
}

// ListSessions lists all sessions for a user
func (s *AuthService) ListSessions(ctx context.Context, req *authnservicev1.ListSessionsRequest) (*authnservicev1.ListSessionsResponse, error) {
	sessions, err := s.sessionRepo.ListByUserID(ctx, req.UserId, req.ActiveOnly, nil)
	if err != nil {
		logger.Errorf("failed to list sessions: %v", err)
		return nil, errors.New("failed to list sessions")
	}
	return &authnservicev1.ListSessionsResponse{Sessions: sessions}, nil
}

// RevokeSession revokes a single session by ID
func (s *AuthService) RevokeSession(ctx context.Context, req *authnservicev1.RevokeSessionRequest) (*authnservicev1.RevokeSessionResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)
	if err := s.sessionRepo.Revoke(ctx, req.SessionId); err != nil {
		logger.Errorf("failed to revoke session: %v", err)
		return nil, errors.New("failed to revoke session")
	}
	appLogger.Infof("Session revoked: %s reason=%q from IP %s", req.SessionId, req.Reason, reqMeta.IPAddress)
	_ = s.eventPublisher.PublishSessionRevoked(ctx, "", req.SessionId, "", "user", req.Reason)
	return &authnservicev1.RevokeSessionResponse{
		Message: "Session revoked successfully",
	}, nil
}

// ValidateCSRF validates a CSRF token for a server-side session
func (s *AuthService) ValidateCSRF(ctx context.Context, req *authnservicev1.ValidateCSRFRequest) (*authnservicev1.ValidateCSRFResponse, error) {
	valid, err := s.tokenService.ValidateCSRFToken(ctx, req.SessionId, req.CsrfToken)
	if err != nil {
		return &authnservicev1.ValidateCSRFResponse{Valid: false, Message: err.Error()}, nil
	}
	if !valid {
		return &authnservicev1.ValidateCSRFResponse{Valid: false, Message: "invalid CSRF token"}, nil
	}
	return &authnservicev1.ValidateCSRFResponse{Valid: true, Message: "CSRF token valid"}, nil
}

// GetCurrentSession resolves the current session from the token in context metadata.
// Supports both JWT (Authorization: Bearer <token>) and server-side session (cookie header).
func (s *AuthService) GetCurrentSession(ctx context.Context, req *authnservicev1.GetCurrentSessionRequest) (*authnservicev1.GetCurrentSessionResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Try JWT first — only if authorization looks like a real Bearer token (starts with eyJ).
	// An empty "authorization" metadata value must be skipped to avoid a failed JWT parse.
	if strings.HasPrefix(reqMeta.Authorization, "eyJ") {
		resp, err := s.tokenService.ValidateJWT(ctx, reqMeta.Authorization)
		if err == nil && resp.Valid {
			session, err := s.sessionRepo.GetByID(ctx, resp.SessionId)
			if err == nil {
				userType := ""
				passwordChangeRequired := false
				if u, uerr := s.userRepo.GetByID(ctx, session.UserId); uerr == nil && u != nil {
					userType = u.UserType.String()
					passwordChangeRequired = u.GetPasswordChangeRequired()
				}
				return &authnservicev1.GetCurrentSessionResponse{
					Session:                session,
					UserType:               userType,
					PasswordChangeRequired: passwordChangeRequired,
				}, nil
			}
		}
	}

	// Try server-side session cookie (session_token=<value> in cookie metadata).
	if reqMeta.SessionToken != "" {
		resp, err := s.tokenService.ValidateServerSideSession(ctx, reqMeta.SessionToken)
		if err == nil && resp.Valid {
			session, err := s.sessionRepo.GetByID(ctx, resp.SessionId)
			if err == nil {
				userType := ""
				passwordChangeRequired := false
				if u, uerr := s.userRepo.GetByID(ctx, session.UserId); uerr == nil && u != nil {
					userType = u.UserType.String()
					passwordChangeRequired = u.GetPasswordChangeRequired()
				}
				return &authnservicev1.GetCurrentSessionResponse{
					Session:                session,
					UserType:               userType,
					PasswordChangeRequired: passwordChangeRequired,
				}, nil
			}
		}
	}

	appLogger.Warnf("GetCurrentSession: no valid session found — sessionToken=%q auth=%q",
		reqMeta.SessionToken != "", reqMeta.Authorization != "")
	return nil, status.Error(codes.NotFound, "no active session found")
}

func (s *AuthService) FindPortalUser(ctx context.Context, req *authnservicev1.FindPortalUserRequest) (*authnservicev1.FindPortalUserResponse, error) {
	identifier := strings.TrimSpace(req.GetIdentifier())
	if identifier == "" {
		return nil, errors.New("identifier is required")
	}

	var (
		user *authnentityv1.User
		err  error
	)

	if strings.Contains(identifier, "@") {
		user, err = s.userRepo.GetByEmail(ctx, strings.ToLower(identifier))
	} else {
		normalized, normalizeErr := normalizeMobileForLookup(identifier)
		if normalizeErr != nil {
			return nil, errors.New("identifier must be a valid email or Bangladesh mobile number")
		}
		user, err = s.userRepo.GetByMobileNumber(ctx, normalized)
	}
	if err != nil {
		return nil, errors.New("user not found")
	}

	summary := &authnservicev1.PortalUserSummary{
		UserId:                 user.GetUserId(),
		Email:                  user.GetEmail(),
		MobileNumber:           user.GetMobileNumber(),
		UserType:               user.GetUserType().String(),
		EmailVerified:          user.GetEmailVerified(),
		PasswordChangeRequired: user.GetPasswordChangeRequired(),
	}

	if s.userProfileRepo != nil {
		if profile, profileErr := s.userProfileRepo.GetByUserID(ctx, user.GetUserId()); profileErr == nil && profile != nil {
			summary.FullName = profile.GetFullName()
			summary.KycVerified = profile.GetKycVerified()
		}
	}

	return &authnservicev1.FindPortalUserResponse{User: summary}, nil
}

func (s *AuthService) SetTemporaryPassword(ctx context.Context, req *authnservicev1.SetTemporaryPasswordRequest) (*authnservicev1.SetTemporaryPasswordResponse, error) {
	if strings.TrimSpace(req.GetUserId()) == "" {
		return nil, errors.New("user_id is required")
	}
	if err := validatePasswordStrength(req.GetTemporaryPassword()); err != nil {
		logger.Errorf("weak temporary password: %v", err)
		return nil, errors.New("weak password")
	}

	user, err := s.userRepo.GetByID(ctx, req.GetUserId())
	if err != nil {
		return nil, errors.New("user not found")
	}

	hashedPassword, err := hashPassword(req.GetTemporaryPassword())
	if err != nil {
		logger.Errorf("failed to hash temporary password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	if err := s.userRepo.UpdatePassword(ctx, user.GetUserId(), hashedPassword); err != nil {
		logger.Errorf("failed to update temporary password: %v", err)
		return nil, errors.New("failed to update password")
	}

	if err := s.userRepo.SetPasswordChangeRequired(ctx, user.GetUserId(), true); err != nil {
		logger.Errorf("failed to mark password_change_required: %v", err)
		return nil, errors.New("failed to update password")
	}

	_ = s.sessionRepo.RevokeAllByUserID(ctx, user.GetUserId(), "")

	reqMeta := s.metadata.ExtractAll(ctx)
	changedBy := strings.TrimSpace(req.GetAssignedBy())
	if changedBy == "" {
		changedBy = grpcmeta.ActorID(ctx, "")
	}
	_ = s.eventPublisher.PublishPasswordChanged(ctx, user.GetUserId(), reqMeta.IPAddress, changedBy)

	return &authnservicev1.SetTemporaryPasswordResponse{
		Message: "Temporary password set successfully",
	}, nil
}

// RevokeAllSessions revokes all sessions for a user (logout from all devices)
func (s *AuthService) RevokeAllSessions(ctx context.Context, req *authnservicev1.RevokeAllSessionsRequest) (*authnservicev1.RevokeAllSessionsResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// IMPORTANT: excludeSessionID must be a session_id, not a session token.
	excludeSessionID := ""
	if req.ExcludeCurrentSession {
		// Try to resolve current session from Authorization (JWT) or Cookie (server-side session)
		if reqMeta.Authorization != "" {
			tok := reqMeta.Authorization
			if len(tok) > 7 && tok[:7] == "Bearer " {
				tok = tok[7:]
			}
			if v, err := s.tokenService.ValidateJWT(ctx, tok); err == nil && v != nil && v.Valid {
				excludeSessionID = v.SessionId
			}
		}

		if excludeSessionID == "" && reqMeta.SessionToken != "" {
			if v, err := s.tokenService.ValidateServerSideSession(ctx, reqMeta.SessionToken); err == nil && v != nil && v.Valid {
				excludeSessionID = v.SessionId
			}
		}
	}

	revoked, err := s.sessionRepo.RevokeAllByUserIDWithCount(ctx, req.UserId, excludeSessionID)
	if err != nil {
		logger.Errorf("failed to revoke sessions: %v", err)
		return nil, errors.New("failed to revoke sessions")
	}

	appLogger.Infof("RevokeAllSessions: sessions revoked for user %s (revoked=%d exclude_session_id=%s) reason=%q from IP %s",
		req.UserId, revoked, excludeSessionID, req.Reason, reqMeta.IPAddress)
	_ = s.eventPublisher.PublishSessionRevoked(ctx, req.UserId, "", "", "user", req.Reason)

	return &authnservicev1.RevokeAllSessionsResponse{
		RevokedCount: int32(revoked),
		Message:      "All sessions revoked successfully",
	}, nil
}

// Helper functions

func parseDeviceType(deviceTypeStr string) authnentityv1.DeviceType {
	// Case-insensitive matching — clients may send "web", "WEB", "Web" etc.
	switch strings.ToUpper(deviceTypeStr) {
	case "WEB":
		return authnentityv1.DeviceType_DEVICE_TYPE_WEB
	case "MOBILE_ANDROID", "ANDROID":
		return authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_ANDROID
	case "MOBILE_IOS", "IOS":
		return authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_IOS
	case "API":
		return authnentityv1.DeviceType_DEVICE_TYPE_API
	case "DESKTOP":
		return authnentityv1.DeviceType_DEVICE_TYPE_DESKTOP
	default:
		// Unknown/garbage device_type (e.g. Postman placeholder "string") must NOT
		// silently fall through to JWT. Return UNSPECIFIED so callers can reject it.
		return authnentityv1.DeviceType_DEVICE_TYPE_UNSPECIFIED
	}
}

func mapDeviceTypeToSessionType(deviceType authnentityv1.DeviceType) authnentityv1.SessionType {
	switch deviceType {
	case authnentityv1.DeviceType_DEVICE_TYPE_WEB:
		return authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE
	default:
		return authnentityv1.SessionType_SESSION_TYPE_JWT
	}
}

func normalizeMobileForLookup(raw string) (string, error) {
	return mobile.NormalizeBangladeshMobileE164(strings.TrimSpace(raw))
}

// portalConfigKeyForUserType maps UserType to the portal config cache key.
// These keys match what the authz Kafka consumer publishes.
func portalConfigKeyForUserType(userType string) string {
	switch userType {
	case authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER.String(), "B2C_CUSTOMER", "USER_TYPE_B2C_CUSTOMER":
		return "PORTAL_B2C"
	case authnentityv1.UserType_USER_TYPE_AGENT.String(), "AGENT", "USER_TYPE_AGENT":
		return "PORTAL_AGENT"
	case authnentityv1.UserType_USER_TYPE_BUSINESS_BENEFICIARY.String(), "BUSINESS_BENEFICIARY", "USER_TYPE_BUSINESS_BENEFICIARY":
		return "PORTAL_BUSINESS"
	case authnentityv1.UserType_USER_TYPE_SYSTEM_USER.String(), "SYSTEM_USER", "USER_TYPE_SYSTEM_USER":
		return "PORTAL_SYSTEM"
	case authnentityv1.UserType_USER_TYPE_PARTNER.String(), "PARTNER", "USER_TYPE_PARTNER":
		return "PORTAL_B2B"
	case authnentityv1.UserType_USER_TYPE_BUSINESS_ADMIN.String(), "BUSINESS_ADMIN", "USER_TYPE_BUSINESS_ADMIN":
		return "PORTAL_B2B"
	case authnentityv1.UserType_USER_TYPE_B2B_ORG_ADMIN.String(), "B2B_ORG_ADMIN", "USER_TYPE_B2B_ORG_ADMIN":
		return "PORTAL_B2B"
	case authnentityv1.UserType_USER_TYPE_B2B_BENEFICIARY.String(), "B2B_BENEFICIARY", "USER_TYPE_B2B_BENEFICIARY":
		return "PORTAL_B2B"
	case authnentityv1.UserType_USER_TYPE_REGULATOR.String(), "REGULATOR", "USER_TYPE_REGULATOR":
		return "PORTAL_REGULATOR"
	default:
		return ""
	}
}

// ================================================================
// BIOMETRIC AUTHENTICATION
// ================================================================

// BiometricAuthenticate authenticates a mobile user using their device-bound
// biometric token. The token is looked up via its HMAC blind index, verified
// by decrypting and constant-time comparing the stored encrypted value.
// On success a new JWT session is issued (always mobile/JWT, never server-side).
func (s *AuthService) BiometricAuthenticate(ctx context.Context, req *authnservicev1.BiometricAuthenticateRequest) (*authnservicev1.BiometricAuthenticateResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Load PII encryptor to compute the blind index for lookup.
	enc, err := pii.NewEncryptorFromEnv()
	if err != nil {
		// If encryption keys are not configured, biometric auth cannot work.
		appLogger.Warnf("BiometricAuthenticate: PII keys not configured: %v", err)
		return nil, errors.New("biometric authentication is not available")
	}

	// Compute blind index of the supplied token to look up the user.
	tokenIdx := enc.BlindIndex(req.BiometricToken)

	// Look up user by biometric_token_idx blind index.
	user, err := s.userRepo.GetByBiometricTokenIdx(ctx, tokenIdx)
	if err != nil {
		appLogger.Warnf("BiometricAuthenticate: token lookup failed from IP %s: %v", reqMeta.IPAddress, err)
		return nil, errors.New("invalid biometric token")
	}

	// Decrypt the stored biometric token and compare.
	storedPlain, err := enc.Decrypt(user.BiometricTokenEnc)
	if err != nil || storedPlain != req.BiometricToken {
		appLogger.Warnf("BiometricAuthenticate: token mismatch for user %s from IP %s", user.UserId, reqMeta.IPAddress)
		return nil, errors.New("invalid biometric token")
	}

	// Verify account is active.
	if user.Status != authnentityv1.UserStatus_USER_STATUS_ACTIVE {
		return nil, errors.New("account is not active")
	}

	// Parse device type (always mobile for biometric).
	deviceType := parseDeviceType(req.DeviceType)
	if deviceType == authnentityv1.DeviceType_DEVICE_TYPE_WEB ||
		deviceType == authnentityv1.DeviceType_DEVICE_TYPE_UNSPECIFIED {
		deviceType = authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_ANDROID
	}

	// Issue JWT tokens.
	tokens, err := s.tokenService.GenerateJWT(
		ctx,
		user.UserId,
		user.UserType.String(),
		"", // tenantID: populated by authz service after role assignment
		req.DeviceId,
		deviceType,
		reqMeta.IPAddress,
		reqMeta.UserAgent,
	)
	if err != nil {
		appLogger.Errorf("BiometricAuthenticate: failed to generate JWT for user %s: %v", user.UserId, err)
		logger.Errorf("failed to generate tokens: %v", err)
		return nil, errors.New("failed to generate tokens")
	}

	_ = s.eventPublisher.PublishUserLoggedIn(ctx, user.UserId, tokens.SessionID, "JWT", reqMeta.IPAddress, req.DeviceType, reqMeta.UserAgent)
	appLogger.Infof("BiometricAuthenticate: success for user %s session %s from IP %s", user.UserId, tokens.SessionID, reqMeta.IPAddress)

	return &authnservicev1.BiometricAuthenticateResponse{
		UserId:                user.UserId,
		SessionId:             tokens.SessionID,
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresIn:  int32(tokens.AccessTokenExpiresIn.Seconds()),
		RefreshTokenExpiresIn: int32(tokens.RefreshTokenExpiresIn.Seconds()),
		SessionType:           "JWT",
		User:                  user,
	}, nil
}

// ================================================================
// DLR (Delivery Report) — internal
// ================================================================

// GetJWKS returns the JWKS (JSON Web Key Set) for the authn service's RSA public key.
// NOTE: GetJWKS RPC was removed from auth_service.proto (API path conflict with well-known domain).
// Kept as internal helper; exposed via HTTP directly by the gateway.
func (s *AuthService) getJWKSInternal(ctx context.Context) (*JWKSResult, error) {
	return s.tokenService.GetJWKSInternal(ctx)
}

// UpdateDLRStatus updates the delivery status of an OTP SMS message.
// Called by the gateway DLR webhook handler (POST /v1/internal/sms/dlr).
func (s *AuthService) UpdateDLRStatus(ctx context.Context, req *authnservicev1.UpdateDLRStatusRequest) (*authnservicev1.UpdateDLRStatusResponse, error) {
	err := s.otpRepo.UpdateDLRStatus(ctx, req.ProviderMessageId, req.Status, req.ErrorCode)
	if err != nil {
		appLogger.Warnf("UpdateDLRStatus: failed for message_id=%s status=%s: %v", req.ProviderMessageId, req.Status, err)
		return &authnservicev1.UpdateDLRStatusResponse{
			Updated: false,
			Message: fmt.Sprintf("failed to update DLR status: %v", err),
		}, nil
	}
	appLogger.Infof("UpdateDLRStatus: message_id=%s status=%s updated", req.ProviderMessageId, req.Status)
	return &authnservicev1.UpdateDLRStatusResponse{
		Updated: true,
		Message: "DLR status updated",
	}, nil
}

// ================================================================
// API KEY MANAGEMENT
// ================================================================

// CreateAPIKey creates a new API key, stores the SHA-256 hash, and returns the
// raw key once (caller must store it securely; it is never retrievable again).
func (s *AuthService) CreateAPIKey(ctx context.Context, req *authnservicev1.CreateAPIKeyRequest) (*authnservicev1.CreateAPIKeyResponse, error) {
	if s.apiKeyRepo == nil {
		logger.Errorf("API key management is not enabled on this server")
		return nil, errors.New("API key management is not enabled on this server")
	}
	// Generate a cryptographically random raw key: "isk_" prefix + 48 random bytes hex.
	rawBytes := make([]byte, 48)
	if _, err := rand.Read(rawBytes); err != nil {
		logger.Errorf("failed to generate API key: %v", err)
		return nil, errors.New("failed to generate API key")
	}
	rawKey := "isk_" + hex.EncodeToString(rawBytes)

	// Hash the raw key for storage (SHA-256; fast lookup; never reversible).
	h := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(h[:])

	keyID := uuid.New().String()

	var expiresAt *timestamppb.Timestamp
	if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt
	}

	ownerType := apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_UNSPECIFIED
	switch strings.ToUpper(req.OwnerType) {
	case "INSURER":
		ownerType = apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INSURER
	case "PARTNER":
		ownerType = apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_PARTNER
	case "INTERNAL":
		ownerType = apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL
	}

	key := &apikeyv1.ApiKey{
		Id:                 keyID,
		KeyHash:            keyHash,
		Name:               req.Name,
		OwnerId:            req.OwnerId,
		OwnerType:          ownerType,
		Scopes:             req.Scopes,
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: req.RateLimitPerMinute,
		ExpiresAt:          expiresAt,
	}

	if err := s.apiKeyRepo.Create(ctx, key); err != nil {
		logger.Errorf("failed to store API key: %v", err)
		return nil, errors.New("failed to store API key")
	}

	appLogger.Infof("CreateAPIKey: key_id=%s owner=%s scopes=%v", keyID, req.OwnerId, req.Scopes)

	return &authnservicev1.CreateAPIKeyResponse{
		KeyId:     keyID,
		RawKey:    rawKey,
		Name:      req.Name,
		Scopes:    req.Scopes,
		ExpiresAt: expiresAt,
	}, nil
}

// ListAPIKeys lists API keys for an owner.
func (s *AuthService) ListAPIKeys(ctx context.Context, req *authnservicev1.ListAPIKeysRequest) (*authnservicev1.ListAPIKeysResponse, error) {
	if s.apiKeyRepo == nil {
		logger.Errorf("API key management is not enabled on this server")
		return nil, errors.New("API key management is not enabled on this server")
	}
	ownerType := apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_UNSPECIFIED
	switch strings.ToUpper(req.OwnerType) {
	case "INSURER":
		ownerType = apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INSURER
	case "PARTNER":
		ownerType = apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_PARTNER
	case "INTERNAL":
		ownerType = apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL
	}

	var statusFilter *apikeyv1.ApiKeyStatus
	if req.ActiveOnly {
		s := apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE
		statusFilter = &s
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 50
	}

	keys, err := s.apiKeyRepo.ListByOwner(ctx, ownerType, req.OwnerId, statusFilter, pageSize, 0)
	if err != nil {
		logger.Errorf("failed to list API keys: %v", err)
		return nil, errors.New("failed to list API keys")
	}

	summaries := make([]*authnservicev1.APIKeySummary, 0, len(keys))
	for _, k := range keys {
		summaries = append(summaries, &authnservicev1.APIKeySummary{
			KeyId:              k.Id,
			Name:               k.Name,
			Scopes:             k.Scopes,
			Status:             k.Status.String(),
			ExpiresAt:          k.ExpiresAt,
			LastUsedAt:         k.LastUsedAt,
			RateLimitPerMinute: k.RateLimitPerMinute,
		})
	}

	return &authnservicev1.ListAPIKeysResponse{Keys: summaries}, nil
}

// RevokeAPIKey revokes an API key by ID.
func (s *AuthService) RevokeAPIKey(ctx context.Context, req *authnservicev1.RevokeAPIKeyRequest) (*authnservicev1.RevokeAPIKeyResponse, error) {
	if s.apiKeyRepo == nil {
		logger.Errorf("API key management is not enabled on this server")
		return nil, errors.New("API key management is not enabled on this server")
	}
	if err := s.apiKeyRepo.Revoke(ctx, req.KeyId); err != nil {
		logger.Errorf("failed to revoke API key %s: %v", req.KeyId, err)
		return nil, errors.New("failed to revoke API key %s")
	}
	appLogger.Infof("RevokeAPIKey: key_id=%s reason=%q", req.KeyId, req.Reason)
	return &authnservicev1.RevokeAPIKeyResponse{Message: "API key revoked"}, nil
}
