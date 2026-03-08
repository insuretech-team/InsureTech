# AuthN Microservice - Complete Code Analysis

## 1. JWT TOKEN GENERATION & CLAIMS

### Location: `internal/service/token_service.go` (Lines 332-441)

#### JWT Claims Structure (InsureTechClaims)
```go
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
```

#### JWT Signing Method
- **Algorithm**: RS256 (RSA 2048-bit)
- **HS256 NOT SUPPORTED** - Explicitly removed in codebase
- **Private Key**: Loaded from PEM file at startup (`JWT_PRIVATE_KEY_PATH`)
- **Public Key**: Loaded from PEM file at startup (`JWT_PUBLIC_KEY_PATH`)
- **KeyID (kid header)**: Configurable via `JWT_KEY_ID` env var

#### GenerateJWT Function (Lines 334-441)
Generates access and refresh token pair with the following claims:

**Access Token Claims:**
```
- sub (Subject): userID
- utp (UserType): User type string
- ins_portal: Portal name derived from userType
- ins_tenant: tenantID (populated by authz service after role assignment)
- ins_device: deviceID (fingerprint or client-provided)
- sid (SessionID): Links to sessions table for revocation
- ins_type: "access"
- jti (JWT ID): Unique identifier for this token
- iss (Issuer): "insuretech-authn" (from config)
- aud (Audience): "insuretech-api" (from config)
- exp (Expiration): Now + AccessTokenDuration (default 15 minutes)
- iat (Issued At): Current time
- kid (Key ID): From config (must match TokenConfig.kid in authz DB)
```

**Refresh Token Claims:**
- Same as access token but with `ins_type: "refresh"`
- Longer expiration: RefreshTokenDuration (default 7 days)
- Different JTI for token rotation detection

#### Token Pair Returned to Client
```go
type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	SessionID             string
	AccessTokenExpiresIn  time.Duration
	RefreshTokenExpiresIn time.Duration
}
```

#### portalForUserType Mapping (Lines 245-263)
Maps UserType to Portal name used in JWT claims:
- `USER_TYPE_SYSTEM_USER` → "system"
- `USER_TYPE_BUSINESS_BENEFICIARY` → "business"
- `USER_TYPE_PARTNER` → "b2b"
- `USER_TYPE_AGENT` → "agent"
- `USER_TYPE_REGULATOR` → "regulator"
- `USER_TYPE_B2C_CUSTOMER` → "b2c"
- Default: "b2c"

---

## 2. LOGIN HANDLER & RESPONSE

### Location: `internal/grpc/auth_handler.go` (Lines 41-55)

```go
func (h *AuthServiceHandler) Login(ctx context.Context, req *authnservicev1.LoginRequest) (*authnservicev1.LoginResponse, error) {
	normalized, err := normalizeMobile(req.MobileNumber)
	if err != nil {
		return nil, err
	}
	req.MobileNumber = normalized
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	resp, err := h.authService.Login(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}
```

### Login Service Implementation: `internal/service/auth_service.go` (Lines 100-258)

#### Input Validation
1. Mobile number normalized to E.164 with '+' prefix (e.g., +8801347210751)
2. Password verification required
3. Device fingerprint generated if not provided:
   - Formula: `fp_` + hex(SHA256(userAgent + "|" + ipAddress)[:16])`

#### Login Flow
1. **Fetch user** by mobile number from DB
2. **Account lockout check** - if `locked_until` > now, deny with remaining time
3. **Password verification** - bcrypt compare with optional rehashing for Argon2id upgrade
4. **Reset login attempts** on success
5. **MFA enforcement check** - Per-portal MFA requirement via `GlobalPortalConfigCache`
   - If MFA required and TOTP not configured: return `mfa_required=true, mfa_method="TOTP"`
   - If MFA required and TOTP configured: return `mfa_required=true` + store MFA session token (5m TTL in Redis)

#### LoginResponse Structure (for JWT clients - Mobile/API)
```
- UserId: user.UserId
- User: Full user proto entity
- SessionId: Generated session UUID
- AccessToken: RS256 signed JWT
- RefreshToken: RS256 signed JWT
- AccessTokenExpiresIn: Duration in seconds (default 15*60)
- RefreshTokenExpiresIn: Duration in seconds (default 7*24*3600)
- SessionType: "JWT"
- MfaRequired: Boolean (if true, client must call VerifyTOTP next)
- MfaMethod: "TOTP" (if MFA required)
- MfaSessionToken: Short-lived token for VerifyTOTP (if MFA required)
```

#### LoginResponse Structure (for Web - Server-Side Session)
```
- UserId: user.UserId
- User: Full user proto entity
- SessionId: Generated session UUID
- SessionToken: Plain token (HttpOnly cookie)
- CsrfToken: For CSRF protection
- SessionType: "SERVER_SIDE"
- MfaRequired: Boolean (if true)
- MfaMethod: "TOTP" (if required)
- MfaSessionToken: MFA session token (if required)
```

#### Device Type → Session Type Mapping (Lines 727-734)
```go
func mapDeviceTypeToSessionType(deviceType authnentityv1.DeviceType) authnentityv1.SessionType {
	switch deviceType {
	case authnentityv1.DeviceType_DEVICE_TYPE_WEB:
		return authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE
	default:
		return authnentityv1.SessionType_SESSION_TYPE_JWT
	}
}
```

#### Events Published
- `PublishLoginFailed` (on credentials invalid, account locked)
- `PublishUserLoggedIn` (on successful login)

---

## 3. USER TYPE & ORG/BUSINESS ID STORAGE

### User Entity Fields: `internal/repository/user_repository.go` (Lines 15-155)

**Database Columns in `authn_schema.users`:**
```
- user_id (UUID, PRIMARY KEY)
- mobile_number (E.164 format with +)
- email
- password_hash (Argon2id)
- status (USER_STATUS enum: ACTIVE, SUSPENDED, DELETED, PENDING_VERIFICATION)
- user_type (USER_TYPE enum)
- created_at
- updated_at
- wallet_balance (int64 paisa)
- email_verified (boolean)
- email_verified_at (timestamp)
- email_login_attempts (int32)
- login_attempts (int32)
- last_login_at (timestamp)
- last_login_session_type (string)
- totp_enabled (boolean)
- totp_secret_enc (encrypted)
- locked_until (timestamp - account lockout)
- email_locked_until (timestamp - email auth lockout)
- notification_preference (string)
- preferred_language (string)
- biometric_token_enc (encrypted)
```

**NOTE: There are NO org_id, business_id, or role fields in the User table.**

These fields are NOT stored in AuthN:
- **org_id**: Populated by AuthZ service AFTER login (stored in AuthZ, not AuthN)
- **business_id**: Same as org_id, populated by AuthZ
- **role**: Populated by AuthZ service based on user assignment

The JWT `ins_tenant` claim is initially empty at login:
```go
// Line 231 in auth_service.go:
"", // tenantID: populated by authz service after role assignment
```

---

## 4. USER TYPES & B2B ADMIN / SUPERADMIN CLASSIFICATION

### UserType Enum Values (from seeder.go and email_auth_service.go)

Located in proto: `gen/go/insuretech/authn/entity/v1/user.pb.go`

**Available User Types:**
```
1. USER_TYPE_B2C_CUSTOMER - Mobile app customer
2. USER_TYPE_AGENT - Insurance agent
3. USER_TYPE_BUSINESS_BENEFICIARY - Business portal user
4. USER_TYPE_SYSTEM_USER - System/Admin user (seeded via ADMIN_EMAIL, SYSTEM_USER type)
5. USER_TYPE_PARTNER - B2B partner
6. USER_TYPE_B2B_ORG_ADMIN - B2B organization admin (seeded via B2B_ADMIN)
7. USER_TYPE_REGULATOR - Regulatory user
```

### B2B Admin User Setup: `internal/seeder/seeder.go` (Lines 105-171)

```go
// SeedB2bAdminUser creates a USER_TYPE_B2B_ORG_ADMIN
admin := &authnentityv1.User{
	UserId:             uuid.NewString(),
	MobileNumber:       adminMobile,
	Email:              adminEmail,
	PasswordHash:       string(hash),
	Status:             authnentityv1.UserStatus_USER_STATUS_ACTIVE,
	UserType:           authnentityv1.UserType_USER_TYPE_B2B_ORG_ADMIN,
	EmailVerified:      true,
	EmailVerifiedAt:    timestamppb.Now(),
	EmailLoginAttempts: 0,
	CreatedAt:          timestamppb.Now(),
	UpdatedAt:          timestamppb.Now(),
}
```

**Environment Variables:**
```
B2B_ADMIN=email@example.com
B2B_ADMIN_MOBILE=+8801XXXXXXXXX
B2B_ADMIN_PASSWARD=password
```

### System User (Superadmin) Setup: `internal/seeder/seeder.go` (Lines 31-103)

```go
admin := &authnentityv1.User{
	UserId:             uuid.NewString(),
	MobileNumber:       adminMobile,
	Email:              adminEmail,
	PasswordHash:       string(hash),
	Status:             authnentityv1.UserStatus_USER_STATUS_ACTIVE,
	UserType:           authnentityv1.UserType_USER_TYPE_SYSTEM_USER,
	EmailVerified:      true,
	EmailVerifiedAt:    timestamppb.Now(),
	EmailLoginAttempts: 0,
	CreatedAt:          timestamppb.Now(),
	UpdatedAt:          timestamppb.Now(),
}
```

**Environment Variables:**
```
ADMIN_EMAIL=admin@example.com
ADMIN_MOBILE=+8801XXXXXXXXX
ADMIN_PASSWORD=password (or ADMIN_PASSWARD for legacy)
```

### Email Auth User Type Restrictions (Lines 30-37 of email_auth_service.go)

```go
func isEmailAuthUser(userType authnentityv1.UserType) bool {
	return userType == authnentityv1.UserType_USER_TYPE_BUSINESS_BENEFICIARY ||
		userType == authnentityv1.UserType_USER_TYPE_SYSTEM_USER ||
		userType == authnentityv1.UserType_USER_TYPE_AGENT ||
		userType == authnentityv1.UserType_USER_TYPE_B2B_ORG_ADMIN
}
```

Only these user types can register and login via email:
- BUSINESS_BENEFICIARY
- SYSTEM_USER (Superadmin)
- AGENT
- B2B_ORG_ADMIN (B2B Admin)

B2C_CUSTOMER can ONLY use mobile OTP login.

---

## 5. SESSION MODEL

### Session Entity: `internal/repository/session_repository.go` (Lines 12-79)

**Database Table: `authn_schema.sessions`**

```go
type Session struct {
	SessionId                string                  // UUID, PRIMARY KEY
	UserId                   string                  // Foreign key to users.user_id
	DeviceId                 string                  // Fingerprint or client device ID
	DeviceType               authnentityv1.DeviceType // DEVICE_TYPE_WEB, MOBILE_ANDROID, MOBILE_IOS, API, DESKTOP
	SessionType              authnentityv1.SessionType // SESSION_TYPE_JWT or SESSION_TYPE_SERVER_SIDE
	SessionTokenHash         string                  // bcrypt hash (for server-side sessions only)
	SessionTokenLookup       string                  // SHA256 hex (for deterministic server-side session lookup)
	AccessTokenJti           string                  // JTI of access token (for JWT sessions)
	RefreshTokenJti          string                  // JTI of refresh token (for JWT sessions)
	AccessTokenExpiresAt     *timestamppb.Timestamp  // When access token expires
	RefreshTokenExpiresAt    *timestamppb.Timestamp  // When refresh token expires
	ExpiresAt                *timestamppb.Timestamp  // When entire session expires
	IpAddress                string                  // Client IP
	UserAgent                string                  // Client user agent
	DeviceName               string                  // Device name (optional)
	CreatedAt                *timestamppb.Timestamp  // Session creation time
	LastActivityAt           *timestamppb.Timestamp  // Last activity (sliding window expiration)
	IsActive                 bool                    // Soft delete flag
	CsrfToken                string                  // CSRF token (server-side sessions only)
}
```

### Session Types

#### SERVER_SIDE Session (Web Portal)
**Flow:** GenerateServerSideSession (Lines 265-330 of token_service.go)

```go
// Key components:
SessionId:               uuid.New().String()
SessionToken:           uuid.New().String() // Returned to client for HttpOnly cookie
SessionTokenHash:       bcrypt.GenerateFromPassword(sessionToken, cost)
SessionTokenLookup:     sha256.Sum256(sessionToken) // For deterministic lookup
CsrfToken:              generateSecureRandomString(32)
SessionType:            SESSION_TYPE_SERVER_SIDE
DeviceType:             DEVICE_TYPE_WEB
ExpiresIn:              12 hours (default, configurable)
```

**Validation Flow (Lines 460-521):**
1. Lookup session by sha256 hash of token
2. Bcrypt verify the token
3. Check expiry
4. Check Redis idle timeout (5m default sliding window)
5. Update last_activity_at for sliding expiration

#### JWT Session (Mobile/API)
**Flow:** GenerateJWT (Lines 334-441)

```go
// Key components:
SessionId:              uuid.New().String()
AccessTokenJti:         uuid.New().String() // Unique token ID
RefreshTokenJti:        uuid.New().String()
SessionType:            SESSION_TYPE_JWT
DeviceType:             DEVICE_TYPE_MOBILE_ANDROID, MOBILE_IOS, API, or DESKTOP
AccessTokenExpiresAt:   Now + 15 minutes (default)
RefreshTokenExpiresAt:  Now + 7 days (default)
ExpiresAt:              Now + RefreshTokenDuration
```

**Validation Flow (Lines 549-594):**
1. Parse RS256 JWT using public key
2. Check device binding (if x-device-id header provided, must match token claim)
3. Check JTI blocklist in Redis
4. Check session is still active in DB
5. Return ValidateTokenResponse with all claims

### Revocation & Token Rotation

**RevokeSession (Lines 681-717):**
```go
// Before marking session as inactive:
1. Read session to get JTIs
2. Block access token JTI in Redis: revoked:jti:<token_id> (TTL = remaining lifetime)
3. Block refresh token JTI in Redis
4. Mark session.is_active = false in DB
```

**RefreshJWT (Lines 619-678):**
```go
1. Validate refresh token
2. Check JTI matches session.refresh_token_jti (prevents reuse attacks)
3. Revoke old session (token rotation)
4. Generate new token pair with new JTIs
5. Create new session record
```

### Concurrent Session Limiting

**SessionLimiter (referenced in token_service.go):**
```go
// Enforces max concurrent sessions per user (default 5)
// When limit exceeded, evicts oldest sessions
if s.sessionLimiter != nil {
	evicted, err := s.sessionLimiter.TrackSession(ctx, userID, sessionID, expiresAt)
	if err == nil {
		for _, evictedID := range evicted {
			_ = s.RevokeSession(ctx, evictedID)
		}
	}
}
```

---

## 6. USER PROFILE ENDPOINT (GetUserProfile)

### GetUserProfile Handler: `internal/grpc/auth_handler.go` (Lines 461-470)

```go
func (h *AuthServiceHandler) GetUserProfile(ctx context.Context, req *authnservicev1.GetUserProfileRequest) (*authnservicev1.GetUserProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.GetUserProfile(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}
```

### GetUserProfile Service: `internal/service/user_profile_service.go` (Lines 71-84)

```go
func (s *AuthService) GetUserProfile(ctx context.Context, req *authnservicev1.GetUserProfileRequest) (*authnservicev1.GetUserProfileResponse, error) {
	if s.userProfileRepo == nil {
		return nil, errors.New("user profile repository not configured")
	}

	profile, err := s.userProfileRepo.GetByUserID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("profile not found: %v", err)
		return nil, errors.New("profile not found")
	}

	return &authnservicev1.GetUserProfileResponse{Profile: profile}, nil
}
```

### UserProfile Entity Structure

**Table: `authn_schema.user_profiles`** (PK = user_id)

```
- user_id (UUID, PRIMARY KEY, FK to users.user_id)
- full_name (string)
- date_of_birth (timestamp)
- gender (GENDER enum)
- occupation (string)
- employer (string)
- address_line1 (string, required)
- address_line2 (string)
- city (string, required)
- district (string, required)
- division (string, required)
- country (string)
- postal_code (string)
- permanent_address (string)
- nid_number (string, required)
- marital_status (string)
- emergency_contact_name (string)
- emergency_contact_number (string)
- profile_photo_url (string)
- kyc_verified (boolean)
- created_at (timestamp)
- updated_at (timestamp)
```

### Response Structure
```go
type GetUserProfileResponse struct {
	Profile *UserProfile
}
```

**NOTE: There is NO "GetMe" endpoint that returns both user + profile + roles.**
- GetUserProfile returns only the profile (demographics)
- User data comes from Login response
- Roles/org_id/business_id are retrieved from AuthZ service (not AuthN)

---

## 7. KEY SECURITY FEATURES

### Password Hashing
- **Algorithm**: Argon2id (default cost)
- **Upgrade path**: On login with bcrypt password, automatically rehash with Argon2id
- **Location**: `internal/service/argon2id.go`, `internal/service/password_hash.go`

### Account Lockout
- **Mobile login**: 5 failed attempts → 30 minute lockout (`login_attempts`, `locked_until`)
- **Email login**: 5 failed attempts → 30 minute lockout (`email_login_attempts`, `email_locked_until`)
- **Check on login**: `if user.LockedUntil != nil && time.Now().Before(...)`

### JWT Revocation
- **Method 1**: Session soft delete (is_active = false)
- **Method 2**: JTI blocklist in Redis with TTL = token lifetime
- **Format**: `revoked:jti:<token_id>`

### Server-Side Session Idle Timeout
- **Implementation**: Redis sliding window key: `session:idle:<session_id>`
- **Default**: Disabled (configurable via `IdleTimeoutDuration`)
- **Behavior**: On validation, if key missing → session expired, else extend TTL

### CSRF Protection
- **Tokens**: Generated for server-side sessions
- **Validation**: POST operations must provide matching CSRF token header
- **Constant-time comparison**: sha256 hash comparison for timing-attack resistance

### Device Binding
- **Access token claim**: `ins_device` contains device fingerprint
- **Validation**: On ValidateJWT, if x-device-id header provided, must match token claim
- **Fingerprint generation**: `fp_` + hex(SHA256(userAgent + "|" + ipAddress)[:16])`

### Biometric Authentication
- **Token storage**: Encrypted in `biometric_token_enc` field
- **Blind index**: `biometric_token_idx` for deterministic lookup without decryption
- **Lookup**: GetByBiometricTokenIdx → decrypt and constant-time compare

---

## 8. CONFIGURATION

### JWT Config (`internal/config/config.go`, Lines 74-88)

```go
type JWTConfig struct {
	PrivateKeyPath       string        // JWT_PRIVATE_KEY_PATH
	PublicKeyPath        string        // JWT_PUBLIC_KEY_PATH
	KeyID                string        // JWT_KEY_ID
	AccessTokenDuration  time.Duration // JWT_ACCESS_TOKEN_DURATION (default 15m)
	RefreshTokenDuration time.Duration // JWT_REFRESH_TOKEN_DURATION (default 7d)
	Issuer               string        // JWT_ISSUER (default "insuretech-authn")
	Audience             string        // JWT_AUDIENCE (default "insuretech-api")
}
```

### Security Config (Lines 137-148)

```go
type SecurityConfig struct {
	ServerSessionDuration  time.Duration // Default 12h
	OTPLength              int
	OTPExpiry              time.Duration
	OTPMaxAttempts         int
	OTPCooldown            time.Duration
	BCryptCost             int           // For server-side session token hashing
	RateLimitPerMinute     int
	RateLimitPerDay        int
	IdleTimeoutDuration    time.Duration // Default 0 (disabled)
}
```

---

## 9. EVENT PUBLISHING

Events published to Kafka (topic: `authn.events` by default):

```
- UserRegistered
- UserLoggedIn (includes session_id, session_type)
- UserLoggedOut (includes logout_reason)
- LoginFailed (includes reason, attempt count)
- TokenRefreshed
- PasswordChanged
- PasswordResetRequested
- SessionRevoked
- CSRFValidationFailed
- OTPSent
- OTPVerified
- EmailLoginSucceeded
- EmailLoginFailed
- EmailVerificationSent
- EmailVerified
- PasswordResetByEmailRequested
```

---

## 10. CRITICAL FINDINGS

### What AuthN Handles
✅ User authentication (credentials verification)
✅ JWT token generation with RS256
✅ Session management (both JWT and server-side)
✅ User type classification (6 types)
✅ OTP generation and verification
✅ Password hashing and reset
✅ MFA (TOTP) enforcement per portal
✅ Account lockout and rate limiting
✅ Biometric authentication
✅ API key management

### What AuthN Does NOT Handle
❌ **Role assignment** - Handled by AuthZ service
❌ **org_id/business_id storage** - Handled by AuthZ service
❌ **Tenant assignment** - Handled by AuthZ service
❌ **Permission checks** - Handled by AuthZ service via CheckAccess RPC

### JWT Claims Flow
```
AuthN Issues:                    AuthZ Populates:
├─ sub (user_id)               ├─ ins_tenant (org_id)
├─ utp (user_type)             ├─ role (via RoleAssignment)
├─ ins_portal (derived)         └─ business_id (via RoleAssignment)
├─ ins_device (fingerprint)
├─ sid (session_id)
├─ ins_type (access|refresh)
└─ jti (token_id)
```

### B2B Admin vs Superadmin
- **Superadmin**: `UserType = USER_TYPE_SYSTEM_USER` (seeded via ADMIN_EMAIL)
- **B2B Admin**: `UserType = USER_TYPE_B2B_ORG_ADMIN` (seeded via B2B_ADMIN)
- Both can use email auth (web portal)
- Role/org assignment happens in AuthZ service post-login

