# AuthZ Microservice - Quick Reference Guide

## File Locations

| Component | Path |
|-----------|------|
| Casbin PERM Model | `internal/enforcer/casbin_enforcer.go:33-48` |
| CheckAccess Logic | `internal/service/authz_service.go:77-216` |
| Super Admin Definition | `internal/seeder/portal_seeder.go:57-67` |
| Action Matcher | `internal/domain/action_matcher.go:15-26` |
| JWT Middleware | `internal/middleware/jwt_interceptor.go:70-85` |
| B2B Fallback | `internal/service/authz_service.go:415-457` |
| Policy Repository | `internal/repository/policy_repository.go` |
| Role Repository | `internal/repository/role_repository.go` |
| User-Role Repository | `internal/repository/user_role_repository.go` |
| gRPC Handlers | `internal/grpc/authz_handler.go:27-271` |

---

## Key Code Snippets

### 1. CheckAccess - Main Permission Check

**File:** `internal/service/authz_service.go:77-216`

```go
func (s *AuthZService) CheckAccess(ctx context.Context, req *CheckAccessRequest) 
    (*CheckAccessResponse, error) {
    
    subject := "user:" + req.UserId
    
    // 1. API Key scope check (if present)
    if req.Context != nil && req.Context.Attributes != nil {
        if !validator.ValidateScope(scopes, req.Object, req.Action) {
            return DENY, nil
        }
    }
    
    // 2. Cache check
    if s.permCache != nil {
        if cached, ok := s.permCache.Get(userId, domain, object, action); ok {
            return cached, nil
        }
    }
    
    // 3. Casbin enforce
    allowed, matchedRule, err := s.enforcer.Enforce(
        ctx, 
        "user:"+req.UserId,    // subject
        req.Domain,             // domain (e.g., "system:root")
        req.Object,             // object (e.g., "svc:policy/create")
        req.Action,             // action (e.g., "POST")
    )
    
    // 4. B2B root fallback (if DENY in tenant domain)
    if !allowed && strings.HasPrefix(req.Domain, "b2b:") && req.Domain != "b2b:root" {
        allowed, matchedRule, _ = s.checkB2BRootDomainFallback(ctx, subject, req)
    }
    
    // 5. Audit & return
    if !allowed || s.auditAllDecisions {
        s.auditRepo.Create(ctx, &AccessDecisionAudit{
            UserId:      req.UserId,
            Domain:      req.Domain,
            Object:      req.Object,
            Action:      req.Action,
            Decision:    DENY/ALLOW,
            MatchedRule: matchedRule,
            DecidedAt:   now,
        })
    }
    
    return &CheckAccessResponse{
        Allowed:     allowed,
        Effect:      POLICY_EFFECT_ALLOW or POLICY_EFFECT_DENY,
        MatchedRule: matchedRule,
        Reason:      "",
    }, nil
}
```

---

### 2. Super Admin Permissions (Seeded)

**File:** `internal/seeder/portal_seeder.go:57-67`

```go
portalSeedMap[Portal_PORTAL_SYSTEM] = []portalRole{
    {
        Name:        "super_admin",
        DisplayName: "Super Administrator",
        IsSystem:    true,
        Policies: []seedPolicy{
            // Grants access to ALL services, ALL actions
            {Object: "svc:*", Action: "*", Effect: "allow"},
            // Explicit B2B coverage
            {Object: "svc:b2b/*", Action: "*", Effect: "allow"},
            {Object: "svc:b2b/*", Action: "GET", Effect: "allow"},
            {Object: "svc:b2b/*", Action: "POST", Effect: "allow"},
            {Object: "svc:b2b/*", Action: "PATCH", Effect: "allow"},
            {Object: "svc:b2b/*", Action: "DELETE", Effect: "allow"},
        },
    },
    // ... other roles (admin, support, auditor, readonly)
}
```

**Assignment to System Users (lines 355-356):**

```go
enforcer.AddRoleForUserInDomain(
    "user:" + systemUserID, 
    "role:super_admin", 
    "system:root",
)
```

---

### 3. Casbin PERM Model

**File:** `internal/enforcer/casbin_enforcer.go:33-48`

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

**Request → Policy Matching:**
1. Check if user has role via g-rule: `g(user:alice, role:admin, system:root)` ✓
2. Check role's domain matches: `r.dom == p.dom` (system:root == system:root) ✓
3. Check object pattern match: `keyMatch2(svc:policy/create, svc:policy/*)` ✓
4. Check action match: `actionMatch(POST, *)` ✓
5. **Result:** ALLOW

---

### 4. Action Matching Logic

**File:** `internal/domain/action_matcher.go:15-26`

```go
// Supported formats:
// 1. "*"              matches ANY action
// 2. "GET"            exact case-insensitive match
// 3. Regex pattern    matched as regexp

func ActionMatches(requestAction, policyAction string) bool {
    requestAction = strings.TrimSpace(requestAction)
    policyAction = strings.TrimSpace(policyAction)
    
    if requestAction == "" || policyAction == "" {
        return false
    }
    
    // Wildcard
    if policyAction == "*" {
        return true
    }
    
    // Exact match (case-insensitive)
    if strings.EqualFold(policyAction, requestAction) {
        return true
    }
    
    // Regex match (invalid regex = false, fail-closed)
    matched, err := regexp.MatchString(policyAction, requestAction)
    return err == nil && matched
}
```

**Examples:**
- `ActionMatches("POST", "*")` → true
- `ActionMatches("POST", "POST")` → true
- `ActionMatches("post", "POST")` → true
- `ActionMatches("GET", "POST")` → false
- `ActionMatches("POST", "GET|POST")` → true (regex)

---

### 5. B2B Root Domain Fallback

**File:** `internal/service/authz_service.go:415-457`

```go
func (s *AuthZService) checkB2BRootDomainFallback(
    ctx context.Context,
    subject string,
    req *CheckAccessRequest,
) (bool, string, error) {
    
    // Only applies to non-root B2B domains
    if !strings.HasPrefix(req.Domain, "b2b:") || req.Domain == "b2b:root" {
        return false, "", nil
    }
    
    // Get user's roles in the tenant domain
    roles, _ := s.enforcer.GetRolesForUserInDomain(subject, req.Domain)
    if len(roles) == 0 {
        return false, "", nil
    }
    
    // Check if those roles have permissions in b2b:root
    rootPolicies, _ := s.listPolicyRulesByDomain(ctx, "b2b:root")
    
    for _, policy := range rootPolicies {
        if policy.Subject not in roles {
            continue
        }
        
        if !keyMatch2(req.Object, policy.Object) || 
           !actionMatch(req.Action, policy.Action) {
            continue
        }
        
        // Found matching policy
        if policy.Effect == POLICY_EFFECT_DENY {
            return false, "root-fallback-deny", nil
        }
        
        return true, "root-fallback: " + policy.Subject, nil
    }
    
    return false, "", nil
}
```

**When Triggered:** Line 161 in CheckAccess after initial DENY in non-root B2B domain

---

### 6. Role Assignment (with Scoped Policy Copy)

**File:** `internal/service/authz_service.go:311-350`

```go
func (s *AuthZService) AssignRole(ctx context.Context, req *AssignRoleRequest) 
    (*AssignRoleResponse, error) {
    
    role, _ := s.roleRepo.GetByID(ctx, req.RoleId)
    domain := req.Domain  // e.g., "b2b:tenant-123"
    
    // 1. Copy B2B root policies to tenant domain
    if err := s.ensureScopedRolePolicies(ctx, role, domain); err != nil {
        return nil, err
    }
    
    // 2. Add g-type rule (user → role in domain)
    subject := "user:" + req.UserId
    roleName := "role:" + role.Name
    s.enforcer.AddRoleForUserInDomain(subject, roleName, domain)
    
    // 3. Persist user_roles row
    ur := &UserRole{
        UserId:     req.UserId,
        RoleId:     req.RoleId,
        Domain:     domain,
        AssignedBy: req.AssignedBy,
        AssignedAt: now,
        ExpiresAt:  req.ExpiresAt,
    }
    s.userRoleRepo.Assign(ctx, ur)
    
    return &AssignRoleResponse{UserRole: ur}, nil
}
```

---

### 7. Policy Creation (Dual Persistence)

**File:** `internal/service/authz_service.go:549-586`

```go
func (s *AuthZService) CreatePolicyRule(ctx context.Context, req *CreatePolicyRuleRequest) 
    (*CreatePolicyRuleResponse, error) {
    
    effect := POLICY_EFFECT_ALLOW
    if req.Effect == POLICY_EFFECT_DENY {
        effect = POLICY_EFFECT_DENY
    }
    
    effectStr := "allow"
    if effect == POLICY_EFFECT_DENY {
        effectStr = "deny"
    }
    
    // 1. Add p-type rule to Casbin enforcer (in-memory + gorm-adapter → DB)
    err := s.enforcer.AddPolicy(
        req.Subject,      // "role:admin"
        req.Domain,       // "system:root"
        req.Object,       // "svc:policy/*"
        req.Action,       // "POST"
        effectStr,        // "allow"
    )
    if err != nil {
        return nil, err
    }
    
    // 2. Persist to policy_rules table (for audit/query)
    pr := &PolicyRule{
        PolicyId:    uuid.New().String(),
        Subject:     req.Subject,
        Domain:      req.Domain,
        Object:      req.Object,
        Action:      req.Action,
        Effect:      effect,
        Condition:   req.Condition,
        Description: req.Description,
        IsActive:    true,
        CreatedBy:   req.CreatedBy,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
    s.policyRepo.Create(ctx, pr)
    
    return &CreatePolicyRuleResponse{Policy: pr}, nil
}
```

---

### 8. JWT Validation Middleware

**File:** `internal/middleware/jwt_interceptor.go:70-85`

```go
func (i *JWTInterceptor) unaryIntercept(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    
    // Skip auth if:
    // - publicKey is nil (no-op mode)
    // - method in skipMethods (e.g., health check)
    // - call from trusted internal service
    if i.publicKey == nil || 
       i.skipMethods[info.FullMethod] || 
       isTrustedInternalCall(ctx) {
        return handler(ctx, req)
    }
    
    // Extract & validate JWT
    claims, err := i.extractClaims(ctx)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }
    
    // Inject claims into context
    ctx = context.WithValue(ctx, ClaimsKey, claims)
    return handler(ctx, req)
}
```

**Claims Extracted:**
```go
type AuthClaims struct {
    UserID   string   // from "sub"
    PortalID string   // from "portal"
    Roles    []string // from "roles"
    Email    string   // from "email"
}
```

---

### 9. Seeding All Portals (Bootstrap)

**File:** `internal/seeder/portal_seeder.go:286-317`

```go
func (s *PortalSeeder) SeedAllPortals(ctx context.Context) error {
    portals := []Portal{
        Portal_PORTAL_SYSTEM,
        Portal_PORTAL_BUSINESS,
        Portal_PORTAL_B2B,
        Portal_PORTAL_AGENT,
        Portal_PORTAL_REGULATOR,
        Portal_PORTAL_B2C,
    }
    
    for _, portal := range portals {
        _, err := s.SeedPortal(ctx, portal, "root", false)
        if err != nil {
            log.Error("seed failed", portal, err)
        } else {
            log.Info("portal seeded", portal)
        }
    }
    
    // Seed configs
    s.SeedPortalConfigs(ctx)      // MFA, session TTLs
    s.SeedTokenConfig(ctx)         // JWT public key
    s.SeedDefaultSystemRoleBindings(ctx)  // Assign super_admin to SYSTEM_USERs
    
    return nil
}
```

---

### 10. Domain Key Format

**File:** `internal/seeder/portal_seeder.go:649-658`

```go
const GlobalTenantID = "root"

func portalDomainKey(portal Portal, tenantID string) string {
    prefix := strings.ToLower(strings.TrimPrefix(portal.String(), "PORTAL_"))
    // PORTAL_SYSTEM → "system"
    // PORTAL_B2B → "b2b"
    // PORTAL_AGENT → "agent"
    
    if tenantID == "" {
        tenantID = GlobalTenantID
    }
    
    return prefix + ":" + tenantID
    // Examples:
    // "system:root"
    // "b2b:tenant-123"
    // "agent:acme-corp"
}
```

---

## Portal Roles Mapping

**File:** `internal/seeder/portal_seeder.go:52-245`

```
PORTAL_SYSTEM → [super_admin*, admin, support, auditor, readonly]
                 (*IsSystem=true)

PORTAL_BUSINESS → [owner, admin, finance, hr, readonly]

PORTAL_B2B → [partner_admin, partner_user, api_client]

PORTAL_AGENT → [senior_agent, agent, agent_trainee]

PORTAL_REGULATOR → [regulator_admin, inspector, auditor]

PORTAL_B2C → [customer]
```

---

## Database Tables

| Table | Schema | Purpose |
|-------|--------|---------|
| roles | authz_schema | Role definitions |
| policy_rules | authz_schema | Human-readable policy rules |
| user_roles | authz_schema | User ↔ Role assignments |
| casbin_rules | authz_schema | Casbin p/g rules (auto-synced) |
| access_decision_audits | authz_schema | Audit log of all access checks |
| portal_configs | authz_schema | MFA, session TTL per portal |
| token_configs | authz_schema | JWT signing key config |

---

## Permission Check Deny-by-Default Flow

```
1. User makes request with:
   - UserId: "alice"
   - Domain: "system:root"
   - Object: "svc:policy/create"
   - Action: "POST"

2. AuthZ CheckAccess is called:
   - Convert to subject: "user:alice"
   - Call Casbin Enforce(user:alice, system:root, svc:policy/create, POST)

3. Casbin Matcher Evaluation:
   - g(user:alice, role:admin, system:root) ?  → User's g-rules
   - role:admin has p-rule (role:admin, system:root, svc:policy/*, POST, allow) ?
   - keyMatch2(svc:policy/create, svc:policy/*) ?  → YES
   - actionMatch(POST, POST) ?  → YES
   - Effect is allow, no deny override  → ALLOW

4. If no matching p-rule found anywhere:
   - No p-rule applies  → DENY (fail-closed)

5. Return Response:
   {Allowed: true/false, Effect: ALLOW/DENY, Reason: "..."}
```

---

## Common Scenarios

### Scenario 1: Super Admin Accessing B2B Service

```
User: system-admin (has super_admin role in system:root)
Request: CheckAccess(
    UserId:  "system-admin",
    Domain:  "system:root",
    Object:  "svc:b2b/organizations",
    Action:  "GET"
)

Casbin Check:
  g(user:system-admin, role:super_admin, system:root) ✓
  p(role:super_admin, system:root, svc:b2b/*, *, allow) ✓
  keyMatch2(svc:b2b/organizations, svc:b2b/*) ✓
  actionMatch(GET, *) ✓

Result: ALLOW ✓
```

### Scenario 2: Tenant User Assigned to Root Domain

```
User: tenant-user (assigned partner_admin in b2b:acme-corp)
Request: CheckAccess(
    UserId:  "tenant-user",
    Domain:  "b2b:root",      // Tries root domain
    Object:  "svc:b2b/departments",
    Action:  "POST"
)

Casbin Check:
  g(user:tenant-user, role:partner_admin, b2b:acme-corp) ✓
  BUT domain mismatch: b2b:acme-corp ≠ b2b:root  ✗

Initial: DENY

Fallback Check (b2b root domain):
  Check if partner_admin has p-rule in b2b:root
  p(role:partner_admin, b2b:root, svc:b2b/*, POST, allow) ?
  → YES (copied during assignment)

Result: ALLOW via fallback ✓
```

### Scenario 3: Unauthorized User

```
User: readonly-user (has readonly role in system:root)
Request: CheckAccess(
    UserId:  "readonly-user",
    Domain:  "system:root",
    Object:  "svc:user/delete",
    Action:  "DELETE"
)

Casbin Check:
  g(user:readonly-user, role:readonly, system:root) ✓
  p(role:readonly, system:root, svc:*, GET, allow) ✓
  BUT actionMatch(DELETE, GET) ✗

No other matching rules.

Result: DENY ✓ (logged to audit, may publish to SIEM)
```

---

## Troubleshooting Tips

| Issue | Likely Cause | Fix |
|-------|--------------|-----|
| Super admin gets DENY | 1) Policy not seeded<br>2) Casbin cache stale<br>3) JWT invalid<br>4) API key scope restricted | 1) Run SeedAllPortals<br>2) Call InvalidatePolicyCache<br>3) Check JWT token<br>4) Check API key scopes |
| Role assignment fails | Missing B2B root policies | Call ensureScopedRolePolicies |
| Permission cache returns stale | Cache TTL too long | Invalidate or reduce TTL |
| Audit not logged | auditAllDecisions=false + ALLOW | Set auditAllDecisions=true |
| SIEM not receiving events | Publisher not configured | Set event broker URL |

