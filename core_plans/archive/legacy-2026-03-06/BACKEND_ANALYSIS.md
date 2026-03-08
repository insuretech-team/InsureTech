# InsureTech Backend Microservices Analysis Report

## Executive Summary

This document provides a comprehensive analysis of three core microservices in the InsureTech backend:
- **AUTHN** (Authentication) - User login, registration, OTP, token management
- **AUTHZ** (Authorization) - Role-based access control using Casbin
- **B2B** (Business-to-Business) - Organization and employee management

All three services follow a clean architecture pattern with clear separation of concerns: gRPC handlers → services → repositories → database.

---

## 1. DIRECTORY STRUCTURE

### AUTHN Microservice
```
backend/inscore/microservices/authn/
├── cmd/server/
│   └── main.go                          # Entry point (507 lines)
├── internal/
│   ├── apierr/                          # Error definitions
│   ├── config/                          # Configuration loading
│   ├── consumers/                       # Kafka event consumers
│   ├── domain/                          # Interfaces
│   ├── email/                           # Email client (SMTP)
│   ├── events/                          # Kafka event publishing
│   ├── grpc/                            # gRPC handler (auth_handler.go)
│   ├── middleware/                      # Metadata extraction
│   ├── pii/                             # PII masking utilities
│   ├── repository/                      # Database access layer (20+ repos)
│   ├── routes/                          # Route definitions
│   ├── seeder/                          # Database seeding
│   ├── service/                         # Business logic (30+ service files)
│   └── sms/                             # SMS client (SSL Wireless)
```

### AUTHZ Microservice
```
backend/inscore/microservices/authz/
├── cmd/server/
│   └── main.go                          # Entry point (295 lines)
├── internal/
│   ├── cache/                           # Permission caching
│   ├── config/                          # Configuration loading
│   ├── domain/                          # Interfaces
│   ├── enforcer/                        # Casbin enforcer
│   ├── events/                          # Kafka producers/consumers
│   ├── grpc/                            # gRPC handler
│   ├── metrics/                         # Prometheus metrics
│   ├── middleware/                      # JWT interceptor
│   ├── repository/                      # Database access (8+ repos)
│   ├── seeder/                          # Portal/policy seeding
│   └── service/                         # Authorization logic
```

### B2B Microservice
```
backend/inscore/microservices/b2b/
├── cmd/server/
│   └── main.go                          # Entry point (302 lines)
├── internal/
│   ├── config/                          # Configuration loading
│   ├── consumers/                       # Event consumers
│   ├── domain/                          # Interfaces
│   ├── events/                          # Event publishing
│   ├── grpc/                            # gRPC handler
│   ├── middleware/                      # AuthZ interceptor
│   ├── repository/                      # Database access (7+ repos)
│   └── service/                         # Business logic (1141 lines)
```

---

## 2. KEY FILES AND THEIR PURPOSE

### Configuration Files

| File | Purpose |
|------|---------|
| `database.yaml` | Primary/backup DB config, failover settings, sync intervals |
| `services.yaml` | gRPC port definitions for all microservices |
| `kafka.yaml` | Kafka broker configuration (Docker Compose format) |
| `logging.yaml` | Log level, format, rotation settings |
| `flve.yaml` | Face Liveness & Verification Engine config |
| `otp_sms.yaml` | SMS provider (SSL Wireless) + OTP settings |
| `s3.yaml` | S3/Spaces storage configuration |
| `storage_layout.yaml` | File storage folder templates & retention |

### Startup Scripts

| Script | Purpose |
|--------|---------|
| `generate.ps1` | Proto code generation + GORM tag injection + TS sync |
| `run_migration.ps1` | Database migration runner (proto-driven) |
| `run_api_pipeline.ps1` | OpenAPI generation & SDK building |
| `start-all.ps1` | Docker Compose orchestration |

---

## 3. HANDLER/ROUTE DEFINITIONS

### AUTHN Handlers (`auth_handler.go`)

**Phone/OTP Flows:**
- `Login(mobile, password)` → Login with OTP validation
- `Register(mobile, password)` → User registration
- `SendOTP(recipient, type)` → Send SMS/Email OTP
- `VerifyOTP(otp_id, code)` → Verify OTP code
- `ResendOTP(original_otp_id)` → Resend OTP with rate limiting

**Token/Session Management:**
- `ValidateToken(access_token|session_id)` → Check token validity
- `RefreshToken(refresh_token)` → Refresh access token
- `Logout(session_id|access_token)` → Invalidate session
- `GetSession(session_id)` → Retrieve session details

**Advanced Features:**
- Email authentication flows
- KYC document verification
- MFA (TOTP) management
- API key management
- Device binding & trust

### AUTHZ Handlers (`authz_handler.go`)

**Core Enforcement:**
- `CheckAccess(user_id, domain, object, action)` → Casbin-based permission check
- `BatchCheckAccess(...)` → Check multiple permissions in one call

**Role Management:**
- `CreateRole(name, portal)` → Create role with permissions
- `GetRole(role_id)` → Retrieve role details
- `UpdateRole(role_id, ...)` → Update role
- `DeleteRole(role_id)` → Delete role
- `ListRoles(portal)` → List all roles

**Assignment:**
- `AssignRole(user_id, role_id, domain)` → Assign role to user
- `RemoveRole(user_id, role_id, domain)` → Remove role

**Policy Management:**
- `CreatePolicyRule(...)` → Create authorization policy
- `ListPolicies(...)` → List policies

### B2B Handlers (`b2b_handler.go`)

**Organization Management:**
- `CreateOrganisation(name, code, ...)` → Create B2B organization
- `GetOrganisation(org_id)` → Get org details
- `ListOrganisations()` → List all orgs
- `UpdateOrganisation(org_id, ...)` → Update org
- `DeleteOrganisation(org_id)` → Delete org
- `ResolveMyOrganisation()` → Get current user's org

**Member Management:**
- `ListOrgMembers(org_id)` → List org members
- `AddOrgMember(org_id, user_id, role)` → Add member
- `AssignOrgAdmin(org_id, user_id)` → Make user admin
- `RemoveOrgMember(org_id, member_id)` → Remove member

**Department & Employees:**
- `CreateDepartment()`, `UpdateDepartment()`, `DeleteDepartment()`
- `ListEmployees()`, `AddEmployee()`, `RemoveEmployee()`

**Purchase Orders & Catalog:**
- `CreatePurchaseOrder()`, `ApprovePurchaseOrder()`, `GetPurchaseOrder()`
- `GetCatalogPlans()`, `GetPlan(plan_id)`

---

## 4. SERVICE/BUSINESS LOGIC

### AUTHN Service (`auth_service.go`)

**Key Responsibilities:**
- User registration with password hashing (Argon2id)
- OTP generation, sending via SMS/Email, verification with rate limiting
- Token lifecycle: generation, refresh, validation, revocation
- Session management with concurrent session limiting (Redis-backed)
- MFA enforcement (TOTP, device binding)
- PII masking for sensitive data
- External KYC orchestration (optional gRPC downstream)
- Email authentication flows
- Password reset & change workflows

**Key Methods:**
```go
Login(ctx, mobile, password) → tokens + session
Register(ctx, mobile, password) → user + session
SendOTP(ctx, recipient, type) → otp_id
VerifyOTP(ctx, otp_id, code) → validated
RefreshToken(ctx, refresh_token) → new_access_token
ValidateToken(ctx, access_token|session_id) → claims
Logout(ctx, session_id) → success
```

**Session Limiter:**
- Enforces max concurrent sessions per user (default 5)
- Uses Redis JTI blocklist for token revocation
- Stores session metadata (IP, device, user agent)

### AUTHZ Service (`authz_service.go`)

**Key Responsibilities:**
- PERM model enforcement using Casbin
- Deny-by-default: missing rule = DENY
- API key scope validation (restricts scopes before Casbin check)
- Permission caching with Redis (configurable TTL)
- Audit logging for all access decisions (configurable)
- Role and policy management
- Portal configuration (MFA requirements, session limits)

**Key Methods:**
```go
CheckAccess(ctx, user_id, domain, object, action) → allowed/denied + reason
BatchCheckAccess(ctx, user_id, domain, checks[]) → results[]
AssignRole(ctx, user_id, role_id, domain) → success
CreateRole(ctx, name, portal, permissions) → role
CreatePolicyRule(ctx, subject, domain, object, action, effect) → rule
```

**PERM Model Structure:**
```
subject = "user:<user_id>"
domain = "portal:tenant_id" or "b2b:org_id" or "system:root"
object = "svc:b2b/*" or "svc:authn/*"
action = "GET", "POST", "PATCH", "DELETE"
```

### B2B Service (`b2b_service.go` - 1141 lines)

**Key Responsibilities:**
- Organization CRUD and membership management
- Department and employee hierarchy
- Purchase order creation and approval workflow
- Catalog plan retrieval with pricing
- Event publishing for org changes
- Authorization checks via downstream AuthZ service

**Key Methods:**
```go
CreateOrganisation(ctx, name, code, email, phone) → org_id
AddOrgMember(ctx, org_id, user_id, role) → member_id
CreatePurchaseOrder(ctx, org_id, plan_id, qty) → po_id
ApprovePurchaseOrder(ctx, po_id) → approval
GetCatalogPlans() → [Seba, Surokkha, Verosa, ...]
```

**Seeded Catalog Plans:**
- **Seba**: Health Insurance - 500 BDT
- **Surokkha**: Health Insurance - 430 BDT
- **Verosa**: Life Insurance - 850 BDT

---

## 5. DATABASE MODELS & QUERIES

### AUTHN Schema

**Core Tables:**
- `authn_schema.users` - User accounts (user_id PK, status, user_type, locked_until)
- `authn_schema.sessions` - Active sessions (session_id PK, user_id FK)
- `authn_schema.otp` - OTP records (otp_id PK, user_id FK, code, delivery_method)
- `authn_schema.api_keys` - API key credentials (api_key_id PK)
- `authn_schema.user_profiles` - Extended profiles
- `authn_schema.user_documents` - Document references
- `authn_schema.kyc_verifications` - KYC data
- `authn_schema.voice_sessions` - Voice authentication

**Access Pattern:**
- Uses raw SQL with proto entity struct scanning
- Proto fields have GORM tags (injected by `inject_gorm_tags.go`)
- Custom scanner functions map SQL rows to proto messages

### AUTHZ Schema

**Core Tables:**
- `authz_schema.roles` - Role definitions (role_id PK, portal FK)
- `authz_schema.user_roles` - User-role assignments (user_id, role_id PK)
- `authz_schema.casbin_rules` - Casbin policy rules (ptype, v0-v5)
- `authz_schema.portal_configs` - Portal settings (MFA, session TTL)
- `authz_schema.access_decision_audits` - Audit log

**Casbin Model:**
```
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act, eft

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

### B2B Schema

**Core Tables:**
- `b2b_schema.organisations` - B2B orgs (organisation_id PK, tenant_id FK)
- `b2b_schema.org_members` - Org members (member_id PK)
- `b2b_schema.departments` - Dept hierarchy
- `b2b_schema.employees` - Employee records
- `b2b_schema.purchase_orders` - PO tracking
- `b2b_schema.catalog_plans` - Available insurance plans

**Access Pattern:**
- Raw SQL with manual row scanning
- COALESCE for nullable fields
- Enum string-to-int conversion for proto enums

---

## 6. AUTHENTICATION & AUTHORIZATION CHECKS

### AUTHN Entry Points

**No explicit auth check** - AUTHN is the auth provider itself. However:
- Incoming requests extract metadata: IP, user agent, device ID, CSRF token
- Requests are traced with `x-request-id`, `x-correlation-id`
- Rate limiting on OTP sends and login attempts

### AUTHZ Entry Points

**JWT Interceptor** (`jwt_interceptor.go`):
```go
publicKey, err := middleware.ParseRSAPublicKeyFromPEM(cfg.Auth.PublicKeyPEM)
// If empty: no-op mode (dev/test)
// If present: validate RS256 JWT signature

skipMethods := [
  "/insuretech.authz.services.v1.AuthZService/CheckPermission",
  "/insuretech.authz.services.v1.AuthZService/CheckAccess",
  "/insuretech.authz.services.v1.AuthZService/GetJWKS",
  "/insuretech.authz.services.v1.AuthZService/ListRoles",
  "/insuretech.authz.services.v1.AuthZService/AssignRole",
  "/insuretech.authz.services.v1.AuthZService/CreatePolicyRule",
  "/grpc.health.v1.Health/Check",
  "/grpc.health.v1.Health/Watch",
]
```

**Rate Limiter:**
- 100 rps steady-state, burst 200 rps

**API Key Scope Validation:**
- Checked BEFORE Casbin to restrict API key scopes
- Scopes parsed from request context attributes
- Format: object/action restrictions

### B2B Entry Points

**AuthZ Interceptor** (`authz_interceptor.go`):
```go
// Extract from metadata:
x-user-id → required
x-business-id (x-organisation-id) → required for most methods

// Map gRPC method to resource/action:
Get*, List* → (svc:b2b/*, GET)
Create*, Add*, Assign* → (svc:b2b/*, POST)
Update* → (svc:b2b/*, PATCH)
Delete*, Remove* → (svc:b2b/*, DELETE)

// Construct domain:
domain = "b2b:{organisation_id}"

// Call downstream AuthZ CheckAccess(user_id, domain, object, action)
```

**Exception:** `ResolveMyOrganisation()` requires no org context (returns user's org)

**Trust Chain:**
- B2B → AuthZ via gRPC (unencrypted, internal only)
- B2B adds `x-internal-service: b2b-service` header for internal calls
- AuthZ recognizes trusted internal services (gateway, b2b-service)

---

## 7. ISSUES IDENTIFIED

### Critical Issues

#### 1. **Unencrypted Downstream gRPC Communication**
- **Location**: B2B → AuthZ, AuthN → KYC
- **Issue**: gRPC connections use `insecure.NewCredentials()`
- **Impact**: Service-to-service calls not encrypted
- **Recommendation**: Use mTLS certificates for internal services

```go
// CURRENT (INSECURE)
grpc.DialContext(..., grpc.WithTransportCredentials(insecure.NewCredentials()))

// SHOULD BE (mTLS)
creds, _ := credentials.NewClientTLSFromFile(certFile, serverName)
grpc.DialContext(..., grpc.WithTransportCredentials(creds))
```

#### 2. **JWT Public Key Not Configured in Production**
- **Location**: `authz/internal/middleware/jwt_interceptor.go` line 200
- **Issue**: If `AUTHZ_JWT_PUBLIC_KEY_PEM` not set, JWT validation is disabled
- **Impact**: AuthZ accepts all requests without token validation in non-prod
- **Recommendation**: Enforce public key in production via startup checks

```go
if publicKey == nil {
    appLogger.Warn("AUTHZ_JWT_PUBLIC_KEY_PEM not set — JWT validation disabled")
    // SHOULD BE: appLogger.Fatal() in production
}
```

#### 3. **No Request Size Limits on gRPC**
- **Location**: All gRPC servers in main.go
- **Issue**: No `grpc.MaxRecvMsgSize`, `grpc.MaxSendMsgSize` configured
- **Impact**: Potential DOS attacks with large payloads
- **Recommendation**: Add message size limits

```go
grpc.NewServer(
    grpc.MaxRecvMsgSize(100*1024*1024), // 100MB
    grpc.MaxSendMsgSize(100*1024*1024),
)
```

#### 4. **Session Limiter Disables Silently on Redis Failure**
- **Location**: `authn/cmd/server/main.go` lines 162-170
- **Issue**: If Redis unavailable, session limiter disabled (no error)
- **Impact**: Concurrent session limits not enforced without Redis
- **Recommendation**: Make Redis mandatory in production or add monitoring alert

```go
if redisClient == nil {
    appLogger.Warn("Redis not available — session limiter disabled")
    // SHOULD BE: appLogger.Fatal() in production
}
```

#### 5. **No Encryption for Sensitive Data at Rest**
- **Location**: Database stores: TOTP secrets, biometric tokens, API keys
- **Issue**: Fields like `totp_secret_enc`, `biometric_token_enc` marked "_enc" but unclear if TDE or column-level encryption
- **Recommendation**: Verify AES encryption is applied; add field-level encryption for PII

### High-Severity Issues

#### 6. **OTP Rate Limiting Depends on Single Clock**
- **Location**: `authn/internal/service/otp_service.go`
- **Issue**: Max resends per hour tracked in DB; clock skew can bypass limits
- **Recommendation**: Use server-side time in DB, add distributed rate limiter (Redis)

#### 7. **No Input Validation for SQL Injection Prevention**
- **Location**: B2B `org_repository.go` - uses parameterized queries ✓ but legacy code
- **Issue**: Some raw SQL queries might be vulnerable if user input not properly parameterized
- **Recommendation**: Audit all raw SQL; prefer ORM for new code

```go
// GOOD (parameterized):
r.db.WithContext(ctx).Where("ptype = ? AND v0 = ?", rule.Ptype, rule.V0)

// RISKY (if legacy code concatenates user input):
query := "SELECT * FROM " + table + " WHERE id = " + id // DON'T DO THIS
```

#### 8. **Kafka Producer Continues Without Broker**
- **Location**: AUTHN/AUTHZ/B2B main.go
- **Issue**: Kafka producer init marked non-fatal; events silently dropped
- **Impact**: Event-driven workflows (user registration, role assignment) fail silently
- **Recommendation**: Make Kafka mandatory in production or implement event retry queue

```go
kafkaProducer, err := producer.NewEventProducerWithRetry(...)
if err != nil {
    appLogger.Warn("Kafka producer init failed — events will be dropped")
    // SHOULD BE: appLogger.Fatal() in production
}
```

### Medium-Severity Issues

#### 9. **No Request Tracing/Correlation IDs in Errors**
- **Location**: All error responses
- **Issue**: Error logs don't include request tracing context
- **Recommendation**: Add context propagation to all responses

```go
// SHOULD INCLUDE:
return nil, status.Errorf(codes.Internal, "request_id=%s: %v", requestID, err)
```

#### 10. **Casbin Model Hardcoded in Code**
- **Location**: `authz/internal/enforcer/casbin_enforcer.go`
- **Issue**: PERM model might be hardcoded; changes require code rebuild
- **Recommendation**: Load model from file (configurable path)

#### 11. **No Pagination Validation**
- **Location**: `b2b/internal/service/b2b_service.go` ListOrganisations, ListEmployees, etc.
- **Issue**: No max page size validation; could request 1M results
- **Recommendation**: Add `max_page_size` config and enforce limits

```go
const maxPageSize = 1000
if req.PageSize > maxPageSize {
    req.PageSize = maxPageSize
}
```

#### 12. **No Graceful Degradation for Optional Services**
- **Location**: `authn/cmd/server/main.go` KYC client (line 206)
- **Issue**: External KYC service integration but no fallback
- **Recommendation**: Define SLA and implement circuit breaker pattern

### Low-Severity Issues

#### 13. **Inconsistent Error Handling Between Services**
- **Location**: AUTHN vs AUTHZ vs B2B
- **Issue**: Error mapping differs; some use codes.Internal, others codes.PermissionDenied
- **Recommendation**: Standardize error codes with error dictionary

#### 14. **No Health Check Readiness Probes**
- **Location**: Kubernetes deployment manifests (not provided)
- **Issue**: Health check only checks DB; should check Kafka, Redis, downstream services
- **Recommendation**: Implement readiness probe that checks all dependencies

#### 15. **Database Connection Pool Too Small**
- **Location**: `database.yaml` - max_open_conns: 15, max_idle_conns: 5
- **Issue**: 15 connections might be insufficient for 3 services under load
- **Recommendation**: Benchmark and increase to 25-50 per service

```yaml
max_open_conns: 50  # 3 services × ~15-20 each
max_idle_conns: 10
conn_max_lifetime: 30m
```

#### 16. **No Secrets Rotation for JWT Keys**
- **Location**: JWT key files loaded at startup
- **Issue**: Key rotation requires service restart
- **Recommendation**: Implement JWKS endpoint with key versioning; support rolling key updates

#### 17. **Missing CORS Headers for gRPC**
- **Location**: gRPC servers (gRPC Web not configured)
- **Issue**: If gRPC Web needed, CORS handling missing
- **Recommendation**: Add grpc-web interceptor if web clients needed

---

## 8. CONFIGURATION SUMMARY

### Environment Variables (Key)

| Variable | Service | Purpose |
|----------|---------|---------|
| `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE` | All | Primary DB |
| `PGHOST2`, `PGPORT2`, etc. | All | Backup DB (Neon) |
| `PGHOST3`, `PGPORT3`, etc. | All | Signal DB |
| `KAFKA_BROKERS` | All | Kafka bootstrap servers |
| `JWT_PRIVATE_KEY_PATH` | AUTHN | RS256 signing key |
| `AUTHZ_JWT_PUBLIC_KEY_PEM` | AUTHZ | RS256 verification key (PEM format) |
| `REDIS_URL` | AUTHN | Redis JTI blocklist + session limiter |
| `SMTP_*` | AUTHN | Email provider (for OTP, password reset) |
| `SSL_*` | AUTHN | SMS provider (SSL Wireless) |

### Port Assignments (services.yaml)

| Service | gRPC Port | HTTP Port |
|---------|-----------|----------|
| authn | 50060 | 50061 |
| authz | 50070 | 50071 |
| b2b | 50110 | 50111 |

---

## 9. DEPENDENCY GRAPH

```
┌─────────────────────────────────────────────────────┐
│                      Gateway/Client                   │
└──────────────┬──────────────────┬────────────────────┘
               │                  │
               v                  v
        ┌──────────────┐   ┌──────────────┐
        │    AUTHN     │   │    AUTHZ     │
        │  (50060)     │   │   (50070)    │
        └──────┬───────┘   └──────┬───────┘
               │                  │
        (JWT Token)        (CheckAccess)
               │                  │
               └────────┬─────────┘
                        │
                        v
              ┌──────────────────┐
              │       B2B        │
              │     (50110)      │
              └──────┬───────────┘
                     │
          (Calls AuthZ for CheckAccess)
                     │
        ┌────────────┼────────────┐
        v            v            v
    ┌────────┐  ┌────────┐  ┌──────────┐
    │Primary │  │Backup  │  │ Signal   │
    │   DB   │  │  DB    │  │   DB     │
    │  (DO)  │  │ (Neon) │  │ (Neon)   │
    └────────┘  └────────┘  └──────────┘

    ┌────────────────────────────────┐
    │      Kafka (Event Bus)          │
    │ (authn, authz, b2b topics)      │
    └────────────────────────────────┘

    ┌────────────┐  ┌──────────────┐
    │   Redis    │  │  S3/Spaces   │
    │(Sessions,  │  │  (Storage)   │
    │ JTI Cache) │  │              │
    └────────────┘  └──────────────┘
```

---

## 10. SECURITY MATRIX

| Layer | Component | Risk Level | Status |
|-------|-----------|------------|--------|
| Transport | gRPC mTLS | **HIGH** | ❌ NOT IMPLEMENTED (insecure) |
| Authentication | RS256 JWT | MEDIUM | ⚠️ Optional in dev |
| Authorization | Casbin PERM | **HIGH** | ✅ Implemented, deny-by-default |
| Session | Concurrent limit | MEDIUM | ⚠️ Redis-dependent |
| Rate Limit | OTP/Login | MEDIUM | ✅ Implemented (BTRC compliant) |
| Data Encryption | At-rest | MEDIUM | ⚠️ Unclear (marked "_enc") |
| Audit | Access decisions | LOW | ✅ Configurable logging |
| PII Protection | Masking | LOW | ✅ Implemented |

---

## 11. DEPLOYMENT CHECKLIST

### Pre-Production

- [ ] Replace insecure gRPC with mTLS
- [ ] Enforce JWT public key in production
- [ ] Make Redis mandatory for session limiting
- [ ] Make Kafka mandatory for event publishing
- [ ] Add gRPC message size limits
- [ ] Configure Casbin model as external file
- [ ] Implement request tracing/correlation IDs
- [ ] Set database connection pool size based on load test
- [ ] Add health check probes for all dependencies
- [ ] Implement circuit breaker for optional services
- [ ] Verify field-level encryption for sensitive data
- [ ] Add pagination limits (max 1000 items)
- [ ] Document all error codes
- [ ] Configure secrets rotation for JWT keys

---

## 12. RECOMMENDATIONS

### Immediate (Critical)

1. **Implement mTLS** for service-to-service communication
2. **Enforce JWT validation** in production
3. **Make Redis mandatory** if using session limiting
4. **Add message size limits** to gRPC servers

### Short-term (1-2 sprints)

5. Audit all raw SQL for injection vulnerabilities
6. Implement request tracing with correlation IDs
7. Add pagination limits validation
8. Implement circuit breaker for external services
9. Configure Casbin model externally

### Medium-term (1-2 months)

10. Implement secrets rotation for JWT keys
11. Add comprehensive health/readiness probes
12. Benchmark and optimize database connection pools
13. Implement distributed rate limiting (Redis-backed)
14. Add gRPC Web support if needed by frontend

---

## 13. PROTO ENTITIES SUMMARY

The system is **proto-first**: all data models are Protobuf messages with GORM tags injected.

**Key Entity Packages:**
- `insuretech.authn.entity.v1` - User, Session, OTP, ApiKey
- `insuretech.authz.entity.v1` - Role, UserRole, CasbinRule
- `insuretech.b2b.entity.v1` - Organisation, OrgMember, Department, Employee, PurchaseOrder
- `insuretech.common.v1` - Common types (Money, Address, etc.)

All generated code includes GORM tags for database mapping.

---

End of Report
