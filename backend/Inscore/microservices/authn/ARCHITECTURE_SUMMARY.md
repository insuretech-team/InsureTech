# AuthN Microservice - Architecture Summary

## 1. Overall Structure

### Directory Layout
```
authn/
├── cmd/server/
│   └── main.go                 # Entry point - bootstraps all services
├── internal/
│   ├── apierr/                 # Error definitions and mapping
│   ├── config/                 # Configuration loading (YAML-based)
│   ├── consumers/              # Kafka event consumers (fan-out handlers)
│   ├── domain/                 # Port/adapter interfaces
│   ├── email/                  # SMTP email client
│   ├── events/                 # Event publisher (Kafka)
│   ├── grpc/                   # gRPC server, handlers, interceptors
│   ├── metrics/                # Prometheus metrics
│   ├── middleware/             # Context metadata extraction
│   ├── pii/                    # PII masking utilities
│   ├── repository/             # Data access layer (10+ repos)
│   ├── routes/                 # (Empty - routes defined in grpc/)
│   ├── seeder/                 # Database initialization
│   ├── service/                # Business logic services
│   └── sms/                    # SMS provider (SSL Wireless)
└── test_all_authn.go           # Root test file
```

### Technology Stack
- **Language**: Go 1.19+
- **gRPC Server**: Port 50053 (configurable via services.yaml)
- **Database**: PostgreSQL (GORM with proto-generated models)
- **Message Queue**: Kafka (event publishing & consumption)
- **Cache/Rate Limiting**: Redis (optional, graceful fallback)
- **JWT**: RS256 (RSA 2048-bit, no HS256 support)
- **SMS Provider**: SSL Wireless (Bangladesh)
- **Email**: SMTP
- **Logging**: Zap-based logger
- **Dependency Injection**: Manual (no wire library used)

---

## 2. Endpoints/RPCs Implemented

### All gRPC Methods (79 total)
All methods are on `insuretech.authn.services.v1.AuthService` service.

**Phone/OTP Authentication (Mobile)**
- `Login(LoginRequest) → LoginResponse`
- `Register(RegisterRequest) → RegisterResponse`
- `SendOTP(SendOTPRequest) → SendOTPResponse`
- `VerifyOTP(VerifyOTPRequest) → VerifyOTPResponse`
- `ResendOTP(ResendOTPRequest) → ResendOTPResponse`

**Token & Session Management**
- `ValidateToken(ValidateTokenRequest) → ValidateTokenResponse`
- `RefreshToken(RefreshTokenRequest) → RefreshTokenResponse`
- `Logout(LogoutRequest) → LogoutResponse`
- `GetSession(GetSessionRequest) → GetSessionResponse`
- `ListSessions(ListSessionsRequest) → ListSessionsResponse`
- `RevokeSession(RevokeSessionRequest) → RevokeSessionResponse`
- `RevokeAllSessions(RevokeAllSessionsRequest) → RevokeAllSessionsResponse`
- `GetCurrentSession(GetCurrentSessionRequest) → GetCurrentSessionResponse`
- `ValidateCSRF(ValidateCSRFRequest) → ValidateCSRFResponse`

**Password Management**
- `ChangePassword(ChangePasswordRequest) → ChangePasswordResponse`
- `ResetPassword(ResetPasswordRequest) → ResetPasswordResponse`

**Email Authentication (Web Portal)**
- `RegisterEmailUser(RegisterEmailUserRequest) → RegisterEmailUserResponse`
- `SendEmailOTP(SendEmailOTPRequest) → SendEmailOTPResponse`
- `VerifyEmail(VerifyEmailRequest) → VerifyEmailResponse`
- `EmailLogin(EmailLoginRequest) → EmailLoginResponse`
- `RequestPasswordResetByEmail(RequestPasswordResetByEmailRequest) → RequestPasswordResetByEmailResponse`
- `ResetPasswordByEmail(ResetPasswordByEmailRequest) → ResetPasswordByEmailResponse`

**Biometric Authentication**
- `BiometricAuthenticate(BiometricAuthenticateRequest) → BiometricAuthenticateResponse`

**DLR (Delivery Receipt) Webhook**
- `UpdateDLRStatus(UpdateDLRStatusRequest) → UpdateDLRStatusResponse`

**API Key Management**
- `CreateAPIKey(CreateAPIKeyRequest) → CreateAPIKeyResponse`
- `ListAPIKeys(ListAPIKeysRequest) → ListAPIKeysResponse`
- `RevokeAPIKey(RevokeAPIKeyRequest) → RevokeAPIKeyResponse`
- `RotateAPIKey(RotateAPIKeyRequest) → RotateAPIKeyResponse`

**User Profile**
- `CreateUserProfile(CreateUserProfileRequest) → CreateUserProfileResponse`
- `GetUserProfile(GetUserProfileRequest) → GetUserProfileResponse`
- `UpdateUserProfile(UpdateUserProfileRequest) → UpdateUserProfileResponse`

**User Documents**
- `UploadUserDocument(UploadUserDocumentRequest) → UploadUserDocumentResponse`
- `ListUserDocuments(ListUserDocumentsRequest) → ListUserDocumentsResponse`
- `GetUserDocument(GetUserDocumentRequest) → GetUserDocumentResponse`
- `UpdateUserDocument(UpdateUserDocumentRequest) → UpdateUserDocumentResponse`
- `DeleteUserDocument(DeleteUserDocumentRequest) → DeleteUserDocumentResponse`

**Document Types**
- `ListDocumentTypes(ListDocumentTypesRequest) → ListDocumentTypesResponse`

**KYC Verification**
- `InitiateKYC(InitiateKYCRequest) → InitiateKYCResponse`
- `GetKYCStatus(GetKYCStatusRequest) → GetKYCStatusResponse`
- `SubmitKYCFrame(SubmitKYCFrameRequest) → SubmitKYCFrameResponse`
- `CompleteKYCSession(CompleteKYCSessionRequest) → CompleteKYCSessionResponse`
- `ApproveKYC(ApproveKYCRequest) → ApproveKYCResponse`
- `RejectKYC(RejectKYCRequest) → RejectKYCResponse`

**Document Verification**
- `VerifyDocument(VerifyDocumentRequest) → VerifyDocumentResponse`

**Voice Sessions & Biometric Auth (Sprint 1.10)**
- `CreateVoiceSession(CreateVoiceSessionRequest) → CreateVoiceSessionResponse`
- `GetVoiceSession(GetVoiceSessionRequest) → GetVoiceSessionResponse`
- `EndVoiceSession(EndVoiceSessionRequest) → EndVoiceSessionResponse`
- `InitiateVoiceSession(InitiateVoiceSessionRequest) → InitiateVoiceSessionResponse`
- `SubmitVoiceSample(SubmitVoiceSampleRequest) → SubmitVoiceSampleResponse`
- `VerifyVoiceSession(VerifyVoiceSessionRequest) → VerifyVoiceSessionResponse`

**Profile & Settings**
- `GetProfilePhotoUploadURL(GetProfilePhotoUploadURLRequest) → GetProfilePhotoUploadURLResponse`
- `UpdateNotificationPreferences(UpdateNotificationPreferencesRequest) → UpdateNotificationPreferencesResponse`

**TOTP / 2FA**
- `EnableTOTP(EnableTOTPRequest) → EnableTOTPResponse`
- `VerifyTOTP(VerifyTOTPRequest) → VerifyTOTPResponse`
- `DisableTOTP(DisableTOTPRequest) → DisableTOTPResponse`

**JWKS (Public Key Distribution)**
- `GetJWKS(GetJWKSRequest) → GetJWKSResponse`

---

## 3. Dependency Injection & Wiring

### Pattern: Manual Constructor-Based DI
The service **does NOT use the `wire` library**. Instead, it uses explicit constructor functions and manual initialization in `main.go`.

### Bootstrap Order (main.go)

1. **Logger Initialization**
   ```go
   appLogger.Initialize(appLogger.Config{Level: "info", Format: "text", Output: "console"})
   ```

2. **Configuration Loading**
   - `services.yaml` → port resolution
   - Environment-specific config via `authnconfig.Load()` (YAML files)
   - Returns `*config.Config` with all subsystems (JWT, SMS, Email, Redis, Kafka, Security)

3. **Database Initialization**
   - `db.InitializeManagerForService(dbConfigPath)` — GORM with PostgreSQL
   - `db.GetDB()` returns singleton `*gorm.DB` instance

4. **Repository Layer** (10+ repositories)
   ```go
   sessionRepo := repository.NewSessionRepository(database)
   userRepo := repository.NewUserRepository(database)
   otpRepo := repository.NewOTPRepository(database)
   apiKeyRepo := repository.NewApiKeyRepository(database)
   userProfileRepo := repository.NewUserProfileRepository(database)
   userDocumentRepo := repository.NewUserDocumentRepository(database)
   documentTypeRepo := repository.NewDocumentTypeRepository(database)
   kycRepo := repository.NewKYCVerificationRepository(database)
   voiceRepo := repository.NewVoiceSessionRepository(database)
   ```

5. **Infrastructure Clients**
   - **Kafka Producer**: `producer.NewEventProducerWithRetry()` — with 5 retries, 3s delay, graceful fallback (nil on failure)
   - **SMS Client**: `sms.NewSSLWirelessClient(cfg)`
   - **Email Client**: `email.NewClient(email.Config{...})`
   - **Redis Client** (optional): `redis.NewClient()` → `rdb.Ping()` — graceful fallback if unavailable

6. **Event Publishing**
   ```go
   eventPublisher := events.NewPublisher(kafkaProducer)
   ```

7. **Middleware**
   ```go
   metadataExtractor := middleware.NewMetadataExtractor()
   ```

8. **Service Layer** (business logic)
   ```go
   // TokenService with Redis-backed session limiting
   tokenService, err := service.NewTokenServiceWithSessionLimiter(
       sessionRepo, userRepo, cfg, eventPublisher, metadataExtractor, redisClient, 0)
   
   // OTP Service
   otpService := service.NewOTPService(otpRepo, smsClient, emailClient, cfg, eventPublisher)
   
   // Auth Service (main facade)
   authService := service.NewAuthService(
       tokenService, otpService, userRepo, sessionRepo, otpRepo,
       apiKeyRepo, userProfileRepo, userDocumentRepo, documentTypeRepo,
       kycRepo, voiceRepo, eventPublisher, cfg, metadataExtractor)
   ```

9. **Downstream KYC Client** (optional, Phase B)
   - Checks `cfg.KYC.Enabled` and `cfg.KYC.Address`
   - Can be HTTP (FLVE) or gRPC (internal KYC service)
   - Set via `authService.SetExternalKYCClient(kycClient)`

10. **Kafka Consumer Group** (async event consumption)
    ```go
    consumerGroup := kafkaconsumer.NewConsumerGroup(kafkaconsumer.Config{
        Brokers: kafkaBrokers,
        GroupID: "authn-service-consumer",
        Topics: [...], // SMS DLR, Account Locked, User Registered, etc.
        Handler: fanOut,
        DLQTopic: "authn.dlq",
    })
    go consumerGroup.Start(consumerCtx)
    ```

11. **Admin User Seeder** (idempotent)
    ```go
    seeder.SeedAdminUser(context.Background(), database)
    ```

12. **Background Cleanup Jobs** (30-min tick)
    - Expired sessions cleanup
    - Expired OTPs cleanup (older than 24h)

13. **gRPC Server**
    ```go
    serverConfig := authnGrpc.DefaultServerConfig()
    serverConfig.Host = cfg.Server.Host
    serverConfig.Port = port // from services.yaml
    serverConfig.DB = database
    
    server, err := authnGrpc.NewServer(serverConfig, authService)
    server.Start() // Listens on :50053
    ```

### gRPC Server Setup (grpc/server.go)

**Handler Factory Pattern**
```go
func (s *Server) registerServices() {
    grpc_health_v1.RegisterHealthServer(s.server, s.health)
    authHandler := NewAuthServiceHandler(s.authService)  // Factory
    authnservicev1.RegisterAuthServiceServer(s.server, authHandler)
    reflection.Register(s.server)
}
```

**Handler Implementation** (`auth_handler.go`)
- `AuthServiceHandler` wraps `AuthServiceIface` (interface, not concrete type)
- Each RPC method validates input → delegates to `authService` → translates errors to gRPC codes
- Input validation includes mobile number normalization via `normalizeMobile()`

**Interceptor Chain**
```go
grpc.ChainUnaryInterceptor(defaultUnaryInterceptors()...)
grpc.ChainStreamInterceptor(defaultStreamInterceptors()...)
```

Unary interceptors (from `interceptors.go`):
1. **Recovery** — converts panics to `codes.Internal`
2. **Request ID** — ensures every request has `x-request-id` in context
3. **Logging** — logs method, duration, gRPC code, request ID
4. **Rate Limiting** — per-IP rate limit enforcement (Redis-backed when available)
5. **Authentication** — validates session token from metadata (for auth-required RPCs)

Stream interceptors: Similar chain for bidirectional/server-streaming RPCs

### Service Layer Architecture

**AuthService** (`service/auth_service.go`)
- Main facade with 79 RPC implementations
- Composes: `TokenService`, `OTPService`, all repositories, event publisher
- Supports setting external KYC client: `SetExternalKYCClient(externalKYCClient)`
- Domain interface: `domain.AuthService`

**TokenService** (`service/token_service.go`)
- Generates JWT tokens (RS256, RSA 2048-bit)
- Manages sessions (CRUD, revocation, listing)
- Password hashing (Argon2id or bcrypt)
- CSRF token generation & validation
- Redis-backed JTI blocklist (optional)
- Session limiter (concurrent session enforcement, default 5)

**OTPService** (`service/otp_service.go`)
- OTP generation (6-digit numeric, configurable length)
- SMS & Email delivery (dual-channel)
- Rate limiting (Redis when available, falls back to DB-based CountRecentOTPs)
- Delivery status tracking (DLR from SSL Wireless)

**KYC Services** (`service/kyc_*.go`)
- `KYCOrchestratorService` — orchestrates face liveness checks
- `KYCExternalClient` interface with two implementations:
  - gRPC client (internal KYC service)
  - HTTP client (FLVE — Face Liveness & Verification Engine)

**Other Services**
- `EmailAuthService` — email-based authentication flow
- `BiometricService` — biometric token handling
- `PortalConfigService` — caches portal config from AuthZ (MFA, session limits, TTLs)
- `TOTPService` — Time-based OTP (2FA)

### Domain Interfaces (domain/interfaces.go)

**Primary Port (Inbound)**
```go
type AuthService interface {
    // 30+ methods covering all auth flows
}
```

**Secondary Ports (Outbound)**
```go
type SessionRepository interface { /* CRUD + revocation */ }
type UserRepository interface { /* User CRUD + status updates */ }
type EventPublisher interface { /* 15+ event publishing methods */ }
```

---

## 4. What is Complete vs Incomplete

### ✅ COMPLETE

**Core Authentication**
- ✅ Phone-based OTP auth (SMS) — fully implemented
- ✅ Email-based auth (SMTP) — fully implemented
- ✅ Password reset flows — both SMS & email variants
- ✅ Session management — create, list, revoke, revoke-all
- ✅ Token management — JWT (RS256), refresh, validation
- ✅ CSRF protection — token generation & validation
- ✅ Password hashing — Argon2id + bcrypt
- ✅ OTP rate limiting — Redis-backed + DB fallback
- ✅ Session limiting — concurrent session enforcement (configurable, default 5)

**Token & Security**
- ✅ RS256 JWT signing (RSA 2048-bit)
- ✅ JWKS endpoint (/.well-known/jwks.json)
- ✅ JTI blocklist (Redis-backed)
- ✅ Device binding (JWT includes device_id)
- ✅ Trusted device tracking
- ✅ Refresh token rotation

**User Management**
- ✅ User registration (phone + email)
- ✅ User profiles (address, NID, KYC data)
- ✅ Document uploads (generic document management)
- ✅ Notification preferences
- ✅ TOTP / 2FA setup & verification

**API Key Management**
- ✅ API key creation, listing, revocation, rotation
- ✅ API key usage tracking

**Event Publishing**
- ✅ Kafka-based event streaming
- ✅ 15+ domain events (UserRegistered, LoginSucceeded, PasswordChanged, etc.)
- ✅ Event consumer group with fan-out handlers
- ✅ DLQ (Dead Letter Queue) for failed messages
- ✅ SMS DLR webhook consumption
- ✅ Account locked consumer
- ✅ Portal config update consumer

**Infrastructure**
- ✅ gRPC server (port 50053, configurable)
- ✅ Health check endpoint
- ✅ Reflection enabled (for grpcurl)
- ✅ Graceful shutdown
- ✅ Request ID tracking
- ✅ Structured logging (Zap)
- ✅ Rate limiting (per-IP)
- ✅ Panic recovery

**Database**
- ✅ PostgreSQL via GORM
- ✅ 10+ repository implementations
- ✅ Proto-generated entities (auto-serialization)
- ✅ Automatic migrations (via GORM)
- ✅ Soft deletes where applicable
- ✅ Background cleanup jobs (sessions, OTPs)

**KYC Integration** (Partially Complete)
- ✅ Local KYC verification storage
- ✅ Frame submission & session management
- ✅ Approval/rejection workflows
- ⚠️ External KYC client (gRPC + HTTP) — wired but not heavily tested
- ✅ Document verification endpoints

**SMS & Email**
- ✅ SSL Wireless SMS integration (Bangladesh, BTRC masking)
- ✅ DLR webhook handling
- ✅ SMTP email delivery
- ✅ Email OTP flows

### ⚠️ PARTIAL / IN-PROGRESS

**Voice Biometric Auth** (Sprint 1.10)
- ✅ Voice session CRUD endpoints exist
- ⚠️ Voice sample submission endpoint exists but incomplete service logic
- ⚠️ Voice verification logic not fully implemented
- Status: Scaffolded but service methods may be stubs

**Biometric Authentication**
- ✅ Endpoint exists (`BiometricAuthenticate`)
- ⚠️ Business logic may be incomplete or placeholder
- Status: Endpoint wired but service implementation unclear

**Portal Config Caching** (Sprint 1.9)
- ✅ Consumer listens to `authz.events` topic
- ⚠️ Cache invalidation & refresh logic may need validation
- Status: Infrastructure ready, logic needs review

**Document Types Management**
- ✅ List endpoint exists
- ⚠️ Full CRUD not implemented (no Create, Update, Delete)
- Status: Read-only for now

### ❌ NOT IMPLEMENTED / TODO

**Multi-Factor Authentication (MFA)**
- ❌ OTP + SMS combination not fully orchestrated
- ❌ MFA enforcement based on portal config (partially in EventConsumer)

**Account Lockout & Recovery**
- ❌ Account lockout after N failed attempts (no explicit lockout service)
- ❌ Lockout recovery flow

**Audit Logging**
- ❌ Detailed audit trail not explicitly modeled
- ⚠️ Events serve as basic audit trail but no dedicated audit log storage

**Advanced KYC Features**
- ❌ Liveness detection (delegated to FLVE, not implemented in AuthN)
- ❌ Face matching against ID documents
- ❌ Fraud detection / AML checks

**OAuth 2.0 / OIDC**
- ❌ No OAuth 2.0 authorization code flow
- ❌ No OIDC support
- Status: Out of scope (may be in separate service)

**Social Login**
- ❌ Google, Facebook, Apple login not implemented

**Account Verification**
- ❌ Email verification workflow (partially in RegisterEmailUser, needs validation)
- ❌ Phone verification after registration

**Session Analytics**
- ❌ Session duration tracking
- ❌ Geographic location tracking

**Rate Limiting Granularity**
- ✅ OTP rate limiting (per-user, per-channel)
- ✅ Refresh token rate limiting
- ⚠️ Login attempt tracking (CountRecentOTPs approach, not dedicated login limiter)

---

## 5. Key Design Patterns & Best Practices

### 1. **Graceful Degradation**
- Kafka producer failure → events dropped, service continues
- Redis unavailable → falls back to DB-based rate limiting
- KYC service unreachable → local repository used
- Email/SMS failures → events published for async retry

### 2. **Error Handling**
- Custom `apierr` package maps domain errors to gRPC codes
- All handler methods call `toGRPCError(err)` to convert
- Panic recovery in interceptor (code: `Internal`)

### 3. **Event-Driven Architecture**
- Kafka publisher on every state change (user registered, password changed, etc.)
- Fan-out consumer groups for multi-tenant event processing
- DLQ for failed message handling

### 4. **Repository Pattern**
- Proto-generated entities with GORM tags
- Each aggregate (User, Session, OTP) has dedicated repository
- `db/sql` null handling for optional fields

### 5. **Middleware & Interceptors**
- Request ID injection (trace correlation)
- Structured logging (Zap)
- Rate limiting (per-IP, Redis-backed)
- Authentication (session token validation)

### 6. **Configuration Management**
- YAML-based config (services.yaml, database.yaml)
- Environment variable overrides
- Centralized config struct with subsystem configs

### 7. **Testing Coverage**
- Unit tests: `*_test.go` files for services, repositories
- Live tests: `*_live_test.go` files (integration tests with real DB/Kafka)
- Mocking via interfaces (`AuthServiceIface`)

---

## 6. Known Issues & Technical Debt

1. **No Wire Dependency Injection** — `main.go` is a 500-line bootstrap (consider factoring into init functions)
2. **Email Auth Incomplete** — Email verification flow needs validation
3. **Voice Biometric** — Scaffolded but incomplete
4. **Account Lockout** — Manual tracking in UserRepository, no dedicated service
5. **KYC External Integration** — HTTP client (FLVE) not tested against real FLVE service
6. **Missing OAuth 2.0** — May be in separate service
7. **Audit Trail** — Events serve as audit trail, but no dedicated audit table/service

