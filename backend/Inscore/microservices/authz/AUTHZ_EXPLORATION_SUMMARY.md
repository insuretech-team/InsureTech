# AuthZ Microservice Exploration Summary

## 1. Overall Folder Structure

```
authz/
├── cmd/server/
│   └── main.go                          # Entry point
├── internal/
│   ├── acl/                             # Access Control List (placeholder)
│   ├── cache/
│   │   ├── permission_cache.go          # Permission caching layer
│   │   └── permission_cache_test.go
│   ├── config/
│   │   ├── config.go                    # Configuration management
│   │   └── config_test.go
│   ├── domain/
│   │   ├── action_matcher.go            # Action matching logic (*, exact, regex)
│   │   ├── action_matcher_test.go
│   │   └── interfaces.go                # Core domain interfaces
│   ├── enforcer/
│   │   ├── casbin_enforcer.go           # Casbin PERM model wrapper
│   │   └── casbin_enforcer_test.go
│   ├── events/
│   │   ├── consumers.go                 # Event consumers
│   │   ├── consumers_test.go
│   │   ├── publisher.go                 # Event publishing (SIEM integration)
│   │   └── publisher_test.go
│   ├── grpc/
│   │   ├── authz_handler.go             # gRPC handler layer
│   │   └── authz_handler_test.go
│   ├── middleware/
│   │   ├── jwt_interceptor.go           # JWT validation middleware
│   │   ├── jwt_interceptor_test.go
│   │   ├── ratelimit_interceptor.go     # Rate limiting middleware
│   │   └── ratelimit_interceptor_test.go
│   ├── metrics/
│   │   ├── metrics.go                   # Prometheus metrics
│   │   └── metrics_test.go
│   ├── repository/
│   │   ├── audit_repository.go          # Access decision audit logs
│   │   ├── authz_repository.go          # Casbin rules repository
│   │   ├── policy_repository.go         # Policy rules CRUD
│   │   ├── portal_repository.go         # Portal configuration
│   │   ├── role_repository.go           # Roles CRUD
│   │   ├── user_role_repository.go      # User-role assignments
│   │   ├── token_config_repository_live_test.go
│   │   └── [*_live_test.go files]       # Database integration tests
│   ├── seeder/
│   │   ├── portal_seeder.go             # Default roles & policies seeding
│   │   └── portal_seeder_test.go
│   └── service/
│       ├── authz_service.go             # Core business logic
│       ├── apikey_scope_validator.go    # API key scope validation
│       ├── authz_service_errors_live_test.go
│       ├── authz_service_live_test.go
│       └── [more test files]
└── AUTHZ_CODE_REPORT.txt
```

---

## 2. How Permissions Are Defined and Checked

### Permission Model: PERM (Policy, Effect, Role, Matcher)

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/enforcer/casbin_enforcer.go`

#### Built-in PERM Model (lines 33-48):
```go
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

#### Format Conventions:
- **Subject:** `"user:<uuid>"` (e.g., `"user:550e8400-..."`)
- **Domain:** `"portal:tenant_id"` (e.g., `"system:root"`, `"agent:tenant-abc"`, `"b2b:root"`)
- **Object:** `"svc:<service>/<resource>"` (e.g., `"svc:policy/create"`, `"svc:b2b/*"`)
- **Action:** HTTP verb or `"*"` (e.g., `"POST"`, `"GET"`, `"*"`)
- **Effect:** `"allow"` or `"deny"`

#### Permission Check Flow (CheckAccess)

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/service/authz_service.go` (lines 77-216)

```
1. API Key Scope Validation (if present, line 81-99)
   ├─ Parse scopes from request context
   ├─ Validate against requested object/action
   └─ DENY if scope doesn't permit

2. Permission Cache Check (line 144-153)
   ├─ If cached and matches, return cached result
   └─ Otherwise, proceed to Casbin

3. Casbin Enforce Call (line 156)
   ├─ Subject: "user:" + req.UserId
   ├─ Domain: req.Domain (e.g., "system:root")
   ├─ Object: req.Object (e.g., "svc:policy/create")
   └─ Action: req.Action (e.g., "POST")

4. B2B Root Domain Fallback (line 161-170)
   ├─ If DENY in non-root b2b domain
   ├─ Check if user's roles have permissions in "b2b:root"
   └─ ALLOW if matched in root domain

5. Cache Store (line 180-182)
   └─ Store result in permission cache

6. Audit Logging (line 187-208)
   ├─ Log all DENY decisions
   ├─ Log ALLOW if auditAllDecisions=true
   └─ Publish DENY events to SIEM
```

#### Action Matching Logic

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/domain/action_matcher.go`

```go
// Supported formats:
// 1. "*"              → matches any action
// 2. "GET"            → exact case-insensitive match
// 3. Regex pattern    → matched as regexp (e.g., "GET|POST")
// Invalid regex = false (fail-closed)
```

**Implementation (lines 15-26):**
```go
func ActionMatches(requestAction, policyAction string) bool {
    if policyAction == "*" || strings.EqualFold(policyAction, requestAction) {
        return true
    }
    matched, err := regexp.MatchString(policyAction, requestAction)
    return err == nil && matched
}
```

---

## 3. How Superadmin Role Is Handled

### Super Admin Definition

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/seeder/portal_seeder.go` (lines 57-67)

```go
authzentityv1.Portal_PORTAL_SYSTEM: {
    {Name: "super_admin", DisplayName: "Super Administrator", IsSystem: true, Policies: []seedPolicy{
        // Wildcard for all services
        {Object: "svc:*", Action: "*", Effect: "allow"},
        // Explicit b2b policies for guaranteed coverage
        {Object: "svc:b2b/*", Action: "*", Effect: "allow"},
        {Object: "svc:b2b/*", Action: "GET", Effect: "allow"},
        {Object: "svc:b2b/*", Action: "POST", Effect: "allow"},
        {Object: "svc:b2b/*", Action: "PATCH", Effect: "allow"},
        {Object: "svc:b2b/*", Action: "DELETE", Effect: "allow"},
    }},
    // ... other roles
}
```

### Super Admin Assignment

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/seeder/portal_seeder.go` (lines 323-375)

**Function:** `SeedDefaultSystemRoleBindings()`

```go
// Query authn_schema.users for SYSTEM_USER type accounts
WHERE user_type IN ('USER_TYPE_SYSTEM_USER', 'SYSTEM_USER', '4')

// Assign super_admin role (not support)
enforcer.AddRoleForUserInDomain("user:" + userID, "role:super_admin", "system:root")
```

**Key Point:** System users get `super_admin` role for **full b2b access** (line 355 comment).

### Super Admin Capabilities

✅ **Has access to:**
- `svc:*` (all services)
- `svc:b2b/*` (all B2B resources)
- All HTTP verbs: `GET`, `POST`, `PATCH`, `DELETE`, `*`
- **No restrictions on tenants/departments** (permission is domain-scoped, not tenant-scoped)

❌ **NOT found to have special 403 handling** (see section 6 below)

---

## 4. Permission Seeds and Initial Data

### Seed Map Definition

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/seeder/portal_seeder.go` (lines 53-245)

#### Portal-to-Roles Mapping:

```
PORTAL_SYSTEM → [super_admin, admin, support, auditor, readonly]
PORTAL_BUSINESS → [owner, admin, finance, hr, readonly]
PORTAL_B2B → [partner_admin, partner_user, api_client]
PORTAL_AGENT → [senior_agent, agent, agent_trainee]
PORTAL_REGULATOR → [regulator_admin, inspector, auditor]
PORTAL_B2C → [customer]
```

#### Key Roles with Full Permissions:

1. **super_admin (PORTAL_SYSTEM)**
   - `svc:*` → `*` (allow)
   - `svc:b2b/*` → all verbs (allow)

2. **owner (PORTAL_BUSINESS)**
   - `svc:*` → `*` (allow)

3. **partner_admin (PORTAL_B2B)**
   - `svc:b2b/*` → all verbs (allow)
   - Limited to tenant/organization scope

### Seeding Process

**Function:** `SeedPortal()` (lines 377-450)

```
1. Load existing policy keys to avoid duplicates
2. For each role in portal:
   a. Create role if not exists
   b. Create policy rules for each seedPolicy
   c. Add p-type rules to Casbin enforcer
   d. Persist to policy_rules table
3. Create g-type role assignments for system users
```

### Domain Key Format

**Function:** `portalDomainKey()` (lines 649-658)

```go
func portalDomainKey(portal authzentityv1.Portal, tenantID string) string {
    prefix := strings.ToLower(strings.TrimPrefix(portal.String(), "PORTAL_"))
    // e.g., PORTAL_SYSTEM → "system"
    if tenantID == "" {
        tenantID = GlobalTenantID // "root"
    }
    return prefix + ":" + tenantID
    // Result: "system:root", "b2b:tenant-123", etc.
}
```

### Global Tenant ID

**Constant:** `GlobalTenantID = "root"` (line 35)

---

## 5. Routes/Endpoints Checking Organization or Department Access

### Domain-Based Scoping

**No explicit organization/department fields found in code.** Instead, scoping is via:

1. **Domain Field** (e.g., `"system:root"`, `"b2b:tenant-123"`)
   - Tenant ID encoded directly in domain
   - Casbin matcher ensures: `r.dom == p.dom`

2. **B2B Root Domain Fallback**

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/service/authz_service.go` (lines 415-457)

```go
func checkB2BRootDomainFallback(ctx context.Context, subject string, req *CheckAccessRequest) 
    (bool, string, error) {
    
    // If DENY in "b2b:tenant-123", check "b2b:root"
    if !allowed && strings.HasPrefix(req.Domain, "b2b:") && req.Domain != "b2b:root" {
        // 1. Get user's roles in the specific domain
        roles := enforcer.GetRolesForUserInDomain(subject, req.Domain)
        
        // 2. Check if those roles have permissions in "b2b:root"
        rootPolicies := listPolicyRulesByDomain("b2b:root")
        
        // 3. If matched in root, ALLOW
        for each policy in rootPolicies:
            if policy.Subject in roles && keyMatch2(req.Object, policy.Object):
                return true, matchedRule
        
        // 4. If DENY rule in root, DENY
        if policy.Effect == POLICY_EFFECT_DENY:
            return false, "root-fallback-deny"
    }
}
```

**When Used:** Line 161-170 in CheckAccess

**Implication:** Allows "promotion" of tenant-specific policies to root domain, but **not the reverse**.

### Scoped Role Policies

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/service/authz_service.go` (lines 459-506)

```go
func ensureScopedRolePolicies(ctx context.Context, role *Role, domain string) error {
    
    // Only for B2B roles assigned to non-root domains
    if role.Portal != PORTAL_B2B || !strings.HasPrefix(domain, "b2b:") || domain == "b2b:root" {
        return nil
    }
    
    // Copy all "b2b:root" policies for this role to the scoped domain
    rootPolicies := listPolicyRulesByDomain("b2b:root")
    for each policy where policy.Subject == "role:" + role.Name:
        enforcer.AddPolicy(roleSubject, domain, policy.Object, policy.Action, effect)
        policyRepo.Create(ctx, policy)
}
```

**When Called:** Line 320 in AssignRole()

**Purpose:** When assigning a B2B role to a user in a non-root domain, automatically copy root domain policies to the tenant domain.

---

## 6. Any Place Where 403 (Forbidden) Is Returned for Superadmin

### Finding: **NO 403 RETURNS FOR SUPERADMIN**

**Search Result:** `grep -r "403|Forbidden|PermissionDenied"` → **No matches found**

**Code Analysis:**

1. **CheckAccess Response** (lines 210-215 in authz_service.go):
   ```go
   return &authzservicev1.CheckAccessResponse{
       Allowed:     allowed,  // true or false
       Effect:      effect,   // POLICY_EFFECT_ALLOW or POLICY_EFFECT_DENY
       MatchedRule: matchedRule,
       Reason:      reason,
   }, nil
   ```
   - Returns gRPC status **OK** with `Allowed: false` in response
   - NO gRPC error code (would translate to 403 HTTP)

2. **gRPC Handler** (lines 27-45 in authz_handler.go):
   ```go
   func (h *AuthZHandler) CheckAccess(ctx context.Context, req *CheckAccessRequest) 
       (*CheckAccessResponse, error) {
       
       resp, err := h.svc.CheckAccess(ctx, req)
       if err != nil {
           return nil, status.Errorf(codes.Internal, "access check failed: %v", err)
       }
       return resp, nil
   }
   ```
   - Errors return `codes.Internal`, not `codes.PermissionDenied`

3. **Super Admin Permissions:**
   - Seeded with `{Object: "svc:*", Action: "*", Effect: "allow"}`
   - Should **always match** in Casbin for any `(svc:*, *)` request
   - Falls through to default DENY only if **no matching rule** exists

### Potential Issues (if superadmin blocked):

If superadmin gets 403, it would be due to:

1. **Policy Cache Staleness** - If `permission_cache.go` returns stale DENY
2. **Casbin State Corruption** - If g/p rules not loaded
3. **JWT Validation** - Middleware rejecting token (lines 70-85 in jwt_interceptor.go)
4. **API Key Scope Validation** - If using API key with restricted scopes (lines 81-99 in authz_service.go)

---

## 7. Seeders, Policies, Middleware, and Permission Handlers

### A. Seeder

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/seeder/portal_seeder.go`

**Key Functions:**

| Function | Purpose | Lines |
|----------|---------|-------|
| `SeedAllPortals()` | Seeds all 6 portals with GlobalTenantID | 286-317 |
| `SeedPortal()` | Seeds single portal+tenant with roles & policies | 377-450 |
| `SeedDefaultSystemRoleBindings()` | Assigns super_admin to SYSTEM_USER accounts | 323-375 |
| `SeedPortalConfigs()` | Seeds MFA & session configs | 474-537 |
| `SeedTokenConfig()` | Seeds JWT signing key from env | 539-578 |
| `portalDomainKey()` | Converts Portal enum → domain string | 649-658 |

**Entry Point (bootstrapping):**
- Called on startup via cmd/server/main.go
- OR via gRPC `SeedPortalDefaults` RPC (if exposed)

---

### B. Policies

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/repository/policy_repository.go`

**Table:** `authz_schema.policy_rules`

**Columns:**
```sql
policy_id    UUID PRIMARY KEY
subject      VARCHAR (e.g., "role:super_admin")
domain       VARCHAR (e.g., "system:root")
object       VARCHAR (e.g., "svc:policy/*")
action       VARCHAR (e.g., "POST" or "*")
effect       VARCHAR ("allow" or "deny")
condition    VARCHAR (optional, for future conditional ACLs)
description  VARCHAR
is_active    BOOLEAN
created_by   UUID (nullable)
created_at   TIMESTAMP
updated_at   TIMESTAMP
deleted_at   TIMESTAMP (soft delete)
```

**CRUD Operations:**

| Operation | Function | Lines |
|-----------|----------|-------|
| Create | `Create()` | 74-101 |
| Read | `GetByID()` | 103-109 |
| Update | `Update()` | 136-155 |
| Delete (soft) | `SoftDelete()` | 111-115 |
| List | `List()` | 117-134 |

**Sync to Casbin:** When created, `policyRepo.Create()` automatically adds p-type rule to Casbin via gorm-adapter (casbin_rules table).

---

### C. Middleware

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/middleware/jwt_interceptor.go`

#### JWT Validation Middleware (lines 34-103)

```go
type JWTInterceptor struct {
    publicKey   *rsa.PublicKey
    skipMethods map[string]bool  // Methods to skip auth (e.g., health check)
}

// Two interceptors:
// 1. UnaryInterceptor()  - for unary gRPC calls
// 2. StreamInterceptor() - for streaming calls
```

**Flow (lines 70-85):**
```go
1. Check if publicKey is nil (no-op mode)
2. Check if method in skipMethods → bypass
3. Check if call from trusted internal service (line 76)
4. Extract JWT from "authorization" metadata header
5. Validate signature using RSA public key
6. Parse claims: sub (UserID), portal, email, roles
7. Inject AuthClaims into context
8. Return codes.Unauthenticated if token invalid
```

**Trusted Internal Services** (lines 41-44):
```go
var trustedInternalServices = map[string]struct{}{
    "gateway":     {},      // API Gateway
    "b2b-service": {},      // B2B Microservice
}
```

**JWT Claims Extracted (lines 137-160):**
```go
claims := &AuthClaims{
    UserID:   sub,                    // from "sub"
    PortalID: portal,                 // from "portal" claim
    Roles:    roles,                  // from "roles" array claim
    Email:    email,                  // from "email"
}
```

#### Rate Limit Interceptor

**File:** `ratelimit_interceptor.go`

(Details not expanded; likely uses redis or in-memory counter per user/IP)

---

### D. Permission Handlers

**File:** `E:/Projects/InsureTech/backend/inscore/microservices/authz/internal/grpc/authz_handler.go`

**Handler Functions:**

| RPC | Handler | Validation | Lines |
|-----|---------|-----------|-------|
| CheckAccess | `CheckAccess()` | UserId, Domain, Object, Action required | 27-45 |
| BatchCheckAccess | `BatchCheckAccess()` | UserId, Domain, at least 1 check | 47-62 |
| CreateRole | `CreateRole()` | Name, Portal required | 66-78 |
| GetRole | `GetRole()` | RoleId required | 80-89 |
| UpdateRole | `UpdateRole()` | RoleId required | 91-99 |
| DeleteRole | `DeleteRole()` | RoleId required | 101-109 |
| ListRoles | `ListRoles()` | None | 111-119 |
| AssignRole | `AssignRole()` | UserId, RoleId, Domain required | 121-135 |
| RevokeRole | `RevokeRole()` | UserId, RoleId, Domain required | 137-149 |
| ListUserRoles | `ListUserRoles()` | UserId required | 151-160 |
| GetUserPermissions | `GetUserPermissions()` | UserId, Domain required | 162-174 |
| CreatePolicyRule | `CreatePolicyRule()` | Subject, Domain, Object, Action required | 178-187 |
| UpdatePolicyRule | `UpdatePolicyRule()` | PolicyId required | 189-198 |
| DeletePolicyRule | `DeletePolicyRule()` | PolicyId required | 200-209 |
| ListPolicyRules | `ListPolicyRules()` | None | 211-217 |
| GetPortalConfig | `GetPortalConfig()` | Portal required | 221-230 |
| UpdatePortalConfig | `UpdatePortalConfig()` | Portal required | 232-241 |
| ListAccessDecisionAudits | `ListAccessDecisionAudits()` | None | 245-251 |
| InvalidatePolicyCache | `InvalidatePolicyCache()` | None | 255-261 |
| GetJWKS | `GetJWKS()` | None | 265-271 |

**Error Handling Pattern:**
```go
if req.UserId == "" {
    return nil, status.Error(codes.InvalidArgument, "user_id is required")
}
resp, err := h.svc.Method(ctx, req)
if err != nil {
    return nil, status.Errorf(codes.Internal, "operation failed: %v", err)
}
return resp, nil
```

---

## Summary: Key Files and Code Snippets

### Core Files by Responsibility:

| Responsibility | File | Key Functions |
|---|---|---|
| **Permission Checking** | `authz_service.go` | `CheckAccess()` (lines 77-216) |
| **Casbin Integration** | `casbin_enforcer.go` | `Enforce()` (lines 106-123), `AddPolicy()`, `AddRoleForUserInDomain()` |
| **Action Matching** | `action_matcher.go` | `ActionMatches()` (lines 15-26) |
| **Default Permissions** | `portal_seeder.go` | `portalSeedMap` (lines 53-245), `SeedPortal()` (lines 377-450) |
| **JWT Validation** | `jwt_interceptor.go` | `UnaryInterceptor()` (lines 70-85) |
| **Policy Storage** | `policy_repository.go` | `Create()` (lines 74-101), `List()` (lines 117-134) |
| **Role Storage** | `role_repository.go` | `GetByNameAndPortal()` (lines 67-81), `List()` (lines 110-149) |
| **User-Role Assignment** | `user_role_repository.go` | `Assign()` (lines 55-90) |
| **Audit Logging** | `audit_repository.go` | `Create()` audit entries |
| **gRPC Routing** | `authz_handler.go` | All RPC handlers (lines 27-271) |

---

## Observations & Design Patterns

1. **PERM Model with Casbin** - Standard pattern; ensures deny-by-default
2. **Domain Scoping** - Via `"portal:tenant_id"` format; no cross-tenant pollution
3. **Root Domain Fallback** - B2B specific; allows tenant-wide policies in root domain
4. **Seeding on Startup** - All roles, policies pre-created; idempotent
5. **Audit Everything** - All DENYs logged; optional logging of ALLOWs
6. **No 403 Special Handling** - Permission denial is a response field, not an HTTP error
7. **API Key Scope Restriction** - Orthogonal to Casbin; validated before Casbin check
8. **Role Inheritance** - Via g-type rules; user inherits role's permissions
9. **Soft Deletes** - Inactive roles/policies preserved in DB
10. **Event Publishing** - DENYs sent to SIEM; role/policy changes published

---

## No Organization/Department Field Found

Despite searching for "organization" and "department" patterns, the codebase uses:
- **Domain** field for tenant/org scoping (not a separate column)
- No explicit department hierarchy
- Multi-tenancy via domain prefix (e.g., `"b2b:acme-corp"`)

If organization/department access control is needed elsewhere, it's likely:
- Implemented in downstream microservices (policy, claim, user services)
- Using domain field as org ID
- Not in AuthZ core logic
