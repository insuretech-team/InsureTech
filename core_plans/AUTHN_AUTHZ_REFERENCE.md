# AuthN & AuthZ Microservices — Central Reference

> **Single source of truth** for InsureTech Authentication & Authorization architecture, proto-based RPCs, Casbin model, permission seeding, and check flows.

---

## 1. Architecture Overview

Two-tier security architecture with gRPC inter-service communication:

```
Client → Gateway → AuthN.ValidateToken → AuthZ.CheckAccess → Backend Service
                                                              ↳ AuthZ.CheckAccess (defense-in-depth)
```

| Layer | Service | Responsibility | Port |
|-------|---------|---------------|------|
| **Authentication** | `authn` | Identity verification, credentials, sessions, tokens | 50053 |
| **Authorization** | `authz` | Fine-grained access control via Casbin PERM model | 50054 |

**Proto packages:**
- AuthN: `insuretech.authn.services.v1` → Go: `authnservicev1`
- AuthZ: `insuretech.authz.services.v1` → Go: `authzservicev1`

**DB schemas:** `authn_schema` (users, sessions, otps, api_keys, profiles, documents, kyc) · `authz_schema` (casbin_rules, roles, user_roles, policy_rules, portal_configs, token_configs, audits)

---

## 2. AuthN Service

### 2.1 Core RPC Methods

#### Phone/OTP & Password Authentication
| RPC | Purpose |
|-----|---------|
| `Register` | Create user with mobile number |
| `SendOTP` / `ResendOTP` | Send/resend OTP (rate limited: 3/hour) |
| `VerifyOTP` | Verify OTP code |
| `Login` | Authenticate mobile + password → session/JWT |
| `RefreshToken` | Rotate access token from refresh token |
| `Logout` | Invalidate session/token |
| `ChangePassword` / `ResetPassword` | Self-service/admin password change |
| `ValidateToken` | Verify JWT/session token, returns user_id + portal |

#### Email Authentication (Business/System Users)
`RegisterEmailUser` · `SendEmailOTP` · `VerifyEmail` · `EmailLogin` · `RequestPasswordResetByEmail` · `ResetPasswordByEmail`

#### Session Management
`GetSession` · `GetCurrentSession` · `ListSessions` · `RevokeSession` · `RevokeAllSessions` · `ValidateCSRF`

#### API Key Management
`CreateAPIKey` · `ListAPIKeys` · `RevokeAPIKey` · `RotateAPIKey`

#### Profile / Documents / KYC / Voice / TOTP / JWKS
`CreateUserProfile` · `GetUserProfile` · `UpdateUserProfile` · `UploadUserDocument` · `ListUserDocuments` · `InitiateKYC` · `GetKYCStatus` · `ApproveKYC` · `RejectKYC` · `EnableTOTP` · `VerifyTOTP` · `DisableTOTP` · `GetJWKS` · `BiometricAuthenticate` · Voice session RPCs

### 2.2 Key Entities

**User:**
```
user_id (UUID) | mobile_number | email | password_hash (Argon2id) | status | user_type | kyc_verified
```

**Session:**
```
session_id (UUID) | user_id | session_type (SERVER_SIDE|JWT) | session_token | device_type | ip_address | expires_at
```

**OTP:**
```
otp_id | recipient | code (encrypted) | type | max_attempts | attempts_used | expires_at | verified
```

### 2.3 Enums

**UserType:** `B2C_CUSTOMER` · `AGENT` · `BUSINESS_BENEFICIARY` · `SYSTEM_USER` · `PARTNER` · `REGULATOR`

**UserStatus:** `PENDING_VERIFICATION` · `ACTIVE` · `SUSPENDED` · `LOCKED` · `DELETED`

**SessionType:** `SERVER_SIDE` (cookie → web portals) · `JWT` (token → mobile/API)

**DeviceType:** `WEB` → SERVER_SIDE · `MOBILE_ANDROID`/`MOBILE_IOS`/`API`/`DESKTOP` → JWT

### 2.4 Auth Flow Patterns

```
Mobile login:  SMS OTP → JWT tokens (access + refresh)
Web login:     Email/password → Server-side session (httpOnly cookie + CSRF token)
API auth:      API key → JWT scoped token
MFA:           TOTP enrolled → mfa_required=true → VerifyTOTP with mfa_session_token
```

### 2.5 Security Features

- **Password:** Argon2id hashing
- **JWT:** RS256 signing, JWKS public key endpoint
- **PII:** AES-256 encryption + blind indexes for masked lookups
- **Rate limiting:** OTP (3/hour/recipient), login failure → account lock
- **Token revocation:** JTI blocklist in Redis
- **CSRF:** Token validation for server-side sessions
- **Device binding:** Optional per-request validation

---

## 3. AuthZ Service — Casbin PERM Model

### 3.1 Built-in Model Definition

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
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && actionMatch(r.act, p.act)
```

### 3.2 Format Specifications

| Component | Format | Example |
|-----------|--------|---------|
| **Domain** | `portal:tenant_id` | `system:root`, `b2b:tenant-acme`, `agent:root` |
| **Subject** | `user:<uuid>` | `user:550e8400-e29b-41d4-a716-446655440000` |
| **Role** | `role:<name>` | `role:super_admin`, `role:partner_user` |
| **Object** | `svc:<service>/<resource>` | `svc:policy/create`, `svc:b2b/*` |
| **Action** | HTTP verb or `*` | `GET`, `POST`, `PATCH`, `DELETE`, `*` |
| **Effect** | `allow` or `deny` | Deny-by-default (no rule = deny) |

### 3.3 Policy Storage — `authz_schema.casbin_rules`

**P-Rules (Permissions):**
```sql
INSERT INTO casbin_rules (p_type, v0, v1, v2, v3, v4)
VALUES ('p', 'role:super_admin', 'system:root', 'svc:*', '*', 'allow');
-- v0=subject, v1=domain, v2=object, v3=action, v4=effect
```

**G-Rules (Role Assignments):**
```sql
INSERT INTO casbin_rules (p_type, v0, v1, v2, v3, v4)
VALUES ('g', 'user:550e8400-...', 'role:super_admin', 'system:root', NULL, NULL);
-- v0=user, v1=role, v2=domain
```

### 3.4 Matcher Logic

| Function | Meaning |
|----------|---------|
| `g(r.sub, p.sub, r.dom)` | User has role in domain (via g-rule) |
| `r.dom == p.dom` | Domain must match exactly |
| `keyMatch2(r.obj, p.obj)` | Object matches via glob pattern |
| `actionMatch(r.act, p.act)` | `*` → always · case-insensitive exact · regex fallback · invalid regex → deny |

### 3.5 AuthZ RPC Methods

#### Enforcement
- `CheckAccess(user_id, domain, object, action, context?)` → allowed, effect, matched_rule, reason
- `BatchCheckAccess(user_id, domain, checks[], context?)` → results[] (UI pre-filtering)

#### Role Management
`CreateRole` · `GetRole` · `UpdateRole` · `DeleteRole` · `ListRoles`

#### User-Role Assignment
`AssignRole` · `RemoveRole` · `ListUserRoles` · `GetUserPermissions`

#### Policy Rules
`CreatePolicyRule` · `UpdatePolicyRule` · `DeletePolicyRule` · `ListPolicyRules`

#### Portal / Audit / Cache / JWKS
`GetPortalConfig` · `UpdatePortalConfig` · `ListAccessDecisionAudits` · `InvalidatePolicyCache` · `GetJWKS`

---

## 4. Permission Seeding — All 6 Portals

### Portal Enum

```
PORTAL_SYSTEM=1  PORTAL_BUSINESS=2  PORTAL_B2B=3  PORTAL_AGENT=4  PORTAL_REGULATOR=5  PORTAL_B2C=6
```

Domain creation: `portalDomainKey(PORTAL_SYSTEM, "root")` → `"system:root"` (strips `PORTAL_` prefix, lowercases, appends `:` + tenantID; default tenantID = `GlobalTenantID = "root"`)

### PORTAL_SYSTEM → `system:root`

| Role | IsSystem | Policies |
|------|----------|----------|
| `super_admin` | ✅ | `svc:*, *` · `svc:b2b/*, *` · `svc:b2b/*, GET` · `svc:b2b/*, POST` · `svc:b2b/*, PATCH` · `svc:b2b/*, DELETE` |
| `admin` | ✅ | `svc:user/*, *` · `svc:role/*, *` · `svc:policy/*, *` · `svc:storage/*, *` · `svc:audit/*, GET` |
| `support` | ❌ | `svc:user/*, GET` · `svc:claim/*, GET` · `svc:session/*, DELETE` · `svc:b2b/*, GET/POST/PATCH/DELETE` |
| `auditor` | ❌ | `svc:audit/*, GET` · `svc:*, GET` |
| `readonly` | ❌ | `svc:*, GET` |

> **Why explicit B2B rules?** Comment: "keyMatch2 does NOT glob-expand `svc:*` to `svc:b2b/*` so we seed explicit rules to guarantee coverage for every HTTP verb."

### PORTAL_B2B → `b2b:root` (template domain)

| Role | Policies |
|------|----------|
| `partner_admin` | `svc:b2b/*, *` · `svc:partner/*, *` · `svc:policy/*, *` · `svc:user/*, *` |
| `partner_user` | `svc:b2b/*, * + per-verb` · `svc:policy/*, GET` · `svc:claim/*, POST/GET` · `svc:storage/* upload/download/update/delete` |
| `b2b_org_admin` | `svc:b2b/*, * + per-verb` · `svc:user/*, *` · `svc:policy/*, *` · `svc:claim/*, *` · `svc:document/*, *` · `svc:storage/*, *` · `svc:invoice/*, *` · `svc:payment/*, *` · `svc:employee/*, *` · `svc:enrollment/*, *` |
| `api_client` | `svc:policy/quote, POST` · `svc:policy/bind, POST` · `svc:claim/submit, POST` |

### PORTAL_BUSINESS → `business:root`

| Role | Policies |
|------|----------|
| `owner` | `svc:*, *` |
| `admin` | `svc:user/*, *` · `svc:policy/*, *` · `svc:claim/*, *` · `svc:document/*, *` · `svc:storage/*, *` |
| `finance` | `svc:invoice/*, *` · `svc:payment/*, *` · `svc:report/*, GET` |
| `hr` | `svc:employee/*, *` · `svc:enrollment/*, *` |
| `readonly` | `svc:*, GET` |

### PORTAL_AGENT → `agent:root`

| Role | Policies |
|------|----------|
| `senior_agent` | Full policy/claim/customer/document/storage + `svc:commission/*, GET` |
| `agent` | policy GET + quote POST · claim GET · customer GET · document/storage CRUD |
| `agent_trainee` | policy GET · customer GET (read-only) |

### PORTAL_REGULATOR → `regulator:root`

| Role | Policies |
|------|----------|
| `regulator_admin` | `svc:audit/*, *` · `svc:report/*, *` · `svc:compliance/*, *` |
| `inspector` | audit/report/compliance GET |
| `auditor` | audit/report GET |

### PORTAL_B2C → `b2c:root`

| Role | Policies |
|------|----------|
| `customer` | `svc:policy/my/*, GET` · `svc:claim/my/*, *` · `svc:profile/*, *` · `svc:document/my/*, *` · `svc:storage/* upload/download/update/delete` |

### System User Bootstrap

`SeedDefaultSystemRoleBindings()` runs on every authz startup:
1. Queries `authn_schema.users` for `user_type IN ('USER_TYPE_SYSTEM_USER', 'SYSTEM_USER', '4')`
2. Assigns `role:super_admin` in `system:root` (idempotent)
3. Creates g-rule: `g(user:<uuid>, role:super_admin, system:root)`

---

## 5. Permission Check Flow

### system:root — Single-Level (No Fallback)

```
CheckAccess(user:xyz, system:root, svc:policy/create, POST)
  → Casbin Enforce(user:xyz, system:root, svc:policy/create, POST)
  → g-rule: user:xyz → role:? in system:root
  → p-rule: role:? → system:root → svc:policy/create → POST → allow?
  → Result: ALLOW or DENY (NO fallback)
```

### B2B Tenant — Two-Level (With Fallback)

```
CheckAccess(user:xyz, b2b:tenant-acme, svc:b2b/organizations, GET)

Stage 1 — Direct Check:
  → Casbin Enforce(user:xyz, b2b:tenant-acme, ...)
  → If ALLOWED → return ALLOW
  → If DENIED → continue to Stage 2

Stage 2 — b2b:root Fallback (only if domain starts with "b2b:" and ≠ "b2b:root"):
  → Get user's roles in b2b:tenant-acme → [role:partner_user, ...]
  → Load ALL b2b:root policies from DB
  → For each policy: match subject → match object (keyMatch2) → match action
  → DENY effect found → return DENY
  → ALLOW effect found → return ALLOW
  → No match → return DENY
```

### Full CheckAccess Pipeline

```
1. API Key scope validation (if present) → deny-by-default on fail
2. Permission cache lookup (if enabled) → return cached result
3. Casbin enforce on requested domain
4. B2B root fallback (if b2b:* domain, not b2b:root)
5. Store in permission cache
6. Audit log (all denials; all decisions if auditAllDecisions=true)
7. SIEM event publish (on deny)
```

| Aspect | system:root | b2b:tenant-xyz |
|--------|-------------|---------------|
| Fallback | NONE | b2b:root policies |
| Trigger | N/A | Primary denied + domain is `b2b:*` + domain ≠ `b2b:root` |
| Roles checked | User's system:root roles | User's tenant roles vs b2b:root policies |

---

## 6. Key Code Locations

| Task | File | Function |
|------|------|----------|
| Create domain key | `portal_seeder.go` | `portalDomainKey()` |
| Seed all portals | `portal_seeder.go` | `SeedAllPortals()` |
| Assign system users | `portal_seeder.go` | `SeedDefaultSystemRoleBindings()` |
| Check permissions | `authz_service.go` | `CheckAccess()` |
| B2B fallback | `authz_service.go` | `checkB2BRootDomainFallback()` |
| Action matching | `action_matcher.go` | `ActionMatches()` |
| Casbin setup | `casbin_enforcer.go` | `New()` |
| Casbin enforce | `casbin_enforcer.go` | `Enforce()` |
| AuthN handler | `auth_handler.go` | All RPC handlers |
| AuthZ handler | `authz_handler.go` | All RPC handlers |

---

## 7. Portal Enum & Shared Patterns

### Portal Enum (shared across services)
```
PORTAL_UNSPECIFIED=0  PORTAL_SYSTEM=1  PORTAL_BUSINESS=2  PORTAL_B2B=3
PORTAL_AGENT=4        PORTAL_REGULATOR=5  PORTAL_B2C=6
```

### PolicyEffect Enum
```
POLICY_EFFECT_UNSPECIFIED=0  POLICY_EFFECT_ALLOW=1  POLICY_EFFECT_DENY=2
```

### Common Error Proto
```
Error { code: string, message: string, details: Detail[] }
```

### Audit Logging
- All denials logged (always)
- All decisions if `auditAllDecisions=true`
- Fields: user, domain, object, action, result, IP, session ID, timestamp
- SIEM event published on deny

### Key Constants
```go
GlobalTenantID = "root"
Domain         = "portal:tenant_id"
Subject        = "user:<uuid>"
Role           = "role:<name>"
Object         = "svc:<service>/<resource>"
Action         = "HTTP verb or *"
```

---

## 8. Current Status & Roadmap

### AuthN — ~85% Complete
✅ Phone/OTP + Email auth · JWT RS256 · Sessions · MFA/TOTP · API Keys · KYC · PII encryption · Event publishing
⚠️ Email verification link flow · Session idle timeout · API key rotation/scoping
❌ Biometric · Social login · SSO/SAML (future)

### AuthZ — ~75% Complete
✅ Casbin PERM · CheckAccess/BatchCheckAccess · Role/Policy CRUD · Portal config · Audit logging · B2B fallback
⚠️ ABAC conditions · Redis distributed cache · Audit retention
❌ Time/IP-based access · Row-level security · Bulk permission pre-loading (future)

### Priority Roadmap
1. **P0 (Critical):** Fix enforcer test failure · Complete email verification · Integration tests · Gateway optimization
2. **P1 (High):** Permission caching (Redis) · Session validation optimization · Per-service AuthZ interceptor · Monitoring dashboards
3. **P2 (Future):** ABAC custom matchers · IP rate limiting · Anomaly detection · Compliance reports · FIDO2/WebAuthn
