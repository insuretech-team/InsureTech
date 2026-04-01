# Rule 04: Security & Authentication Standards

**Scope:** All REST API endpoints across all 34 services  
**Priority:** 🔴 CRITICAL

---

## Why Explicit Security Declaration Is Required

The global `openapi.yaml` may define security schemes, but every endpoint still needs explicit security behavior. Without that:

- SDK generators produce clients with no auth headers
- iOS/Android devs must guess which endpoints need tokens
- Swagger UI "Authorize" applies globally — no per-endpoint visibility
- Frontend DI containers cannot auto-inject tokens correctly

---

## Security Scheme Definition (in openapi.yaml)

```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: >
        JWT access token obtained from POST /v1/auth/login or 
        POST /v1/auth/email/login. Include as:
        Authorization: Bearer <token>
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
      description: >
        API Key for B2B/partner integrations. Obtained from 
        POST /v1/auth/api-keys or POST /v1/partners/{id}/credentials:rotate
    TenantAuth:
      type: apiKey
      in: header
      name: X-Tenant-ID
      description: >
        Tenant identifier for multi-tenant requests. Required alongside
        BearerAuth for tenant-scoped operations.
```

---

## Security Classification Rules

### 🌍 PUBLIC Endpoints (no auth required)
Declare `security: []` explicitly to override global auth:

```yaml
/v1/auth/register:
  post:
    security: []    # explicitly public
    ...

/v1/auth/email/register:
  post:
    security: []

/v1/auth/login:
  post:
    security: []

/v1/auth/email/login:
  post:
    security: []

/v1/auth/otp:send:
  post:
    security: []

/v1/auth/otp:verify:
  post:
    security: []

/v1/auth/otp:resend:
  post:
    security: []

/v1/auth/email/otp:send:
  post:
    security: []

/v1/auth/email/verify:
  post:
    security: []

/v1/auth/email/password:reset-request:
  post:
    security: []

/v1/auth/email/password:reset:
  post:
    security: []

/v1/auth/password:reset:
  post:
    security: []

/v1/auth/.well-known/jwks.json:
  get:
    security: []

/v1/products:
  get:
    security: []    # product catalogue is public browsing

/v1/products/{product_id}:
  get:
    security: []

/v1/products:search:
  get:
    security: []
```

### 🔐 BEARER AUTH Endpoints (JWT token required)
All other endpoints default to BearerAuth. Set globally:

```yaml
# In openapi.yaml — global default
security:
  - BearerAuth: []
```

### 🔑 API KEY Endpoints (partner/B2B integrations)

```yaml
/v1/partners/{partner_id}/credentials:rotate:
  post:
    security:
      - ApiKeyAuth: []

/v1/payments/webhook/{provider}:
  post:
    security: []    # Webhooks authenticated via HMAC signature in body/header
```

### 🏢 TENANT-SCOPED Endpoints

```yaml
/v1/tenants/{tenant_id}:
  get:
    security:
      - BearerAuth: []
        TenantAuth: []   # Both required (AND logic)
```

---

## Required Security Headers

### All Authenticated Requests
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
X-Request-ID: req_uuid_here          # optional but recommended for tracing
```

### Multi-Tenant Requests
```
Authorization: Bearer <token>
X-Tenant-ID: tenant_abc123
Content-Type: application/json
```

### Partner/B2B Requests
```
X-API-Key: ak_live_abc123xyz
Content-Type: application/json
```

### Webhook Endpoints
```
X-Webhook-Signature: sha256=<hmac_signature>
X-Webhook-Timestamp: 1705312800
Content-Type: application/json
```

---

## Token Lifecycle

```
Register/Login → access_token (15 min) + refresh_token (30 days)
                      ↓
             Use access_token for all API calls
                      ↓
             Token expires → POST /v1/auth/token:refresh
                      ↓
             New access_token issued
                      ↓
             Logout → POST /v1/auth/logout (revokes refresh token)
```

### Token Refresh Response (what client receives)
```json
{
  "success": true,
  "data": {
    "access_token": "eyJ...",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "ref_...",
    "refresh_expires_in": 2592000
  },
  "error": null
}
```

---

## Per-Endpoint Security Annotation in YAML

```yaml
# Protected endpoint example
/v1/policies:
  post:
    summary: Create policy
    operationId: PolicyService_CreatePolicy
    security:
      - BearerAuth: []    # ← REQUIRED on every protected endpoint
    tags:
      - PolicyService
    x-roles-required:     # ← custom extension for documentation
      - AGENT
      - ADMIN
    ...

# Public endpoint example  
/v1/auth/login:
  post:
    summary: Login
    operationId: AuthService_Login
    security: []           # ← explicitly no auth
    tags:
      - AuthService
    ...
```

---

## Role-Based Access Control Documentation

Use `x-roles-required` extension to document which roles can access each endpoint:

```yaml
x-roles-required:
  - SUPER_ADMIN      # Full system access
  - ADMIN            # Tenant admin
  - AGENT            # Insurance agent
  - CUSTOMER         # End user
  - PARTNER          # B2B partner
  - UNDERWRITER      # Underwriting staff
  - CLAIMS_ADJUSTER  # Claims processing staff
  - FINANCE          # Finance/billing staff
  - AUDITOR          # Read-only audit access
```

---

## Generator Fix Required (path_generator.py)

```python
# Public endpoints — no auth
PUBLIC_OPERATIONS = {
    'AuthService_Register', 'AuthService_Login', 'AuthService_SendOTP',
    'AuthService_VerifyOTP', 'AuthService_ResendOTP',
    'AuthService_RegisterEmailUser', 'AuthService_EmailLogin',
    'AuthService_SendEmailOTP', 'AuthService_VerifyEmail',
    'AuthService_RequestPasswordResetByEmail', 'AuthService_ResetPasswordByEmail',
    'AuthService_ResetPassword', 'AuthService_GetJWKS',
}

def get_security(operation_id, path):
    if operation_id in PUBLIC_OPERATIONS:
        return []           # security: []
    if '/webhook/' in path:
        return []           # webhooks use HMAC, not bearer
    return [{'BearerAuth': []}]   # default: bearer auth
```
