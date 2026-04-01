# InsureTech Project Structure & API Rules - Exploration Summary

## Project Overview

**InsureTech** is a comprehensive insurance technology platform with:
- 34 REST API microservices
- Multi-platform clients (iOS, Android, Web)
- B2B, Customer, Partner, and System portals
- Protobuf-based service definitions
- SDK generators for Go and TypeScript

---

## Directory Structure

```
E:\Projects\InsureTech\
├── proto/                          # Protobuf service definitions (34 services)
│   └── insuretech/
│       ├── ai/                    # AI/Agent services
│       ├── analytics/             # Analytics & metrics
│       ├── apikey/                # API key management
│       ├── audit/                 # Audit logging & compliance
│       ├── authn/                 # Authentication
│       ├── authz/                 # Authorization & RBAC
│       ├── b2b/                   # B2B operations
│       ├── beneficiary/           # Beneficiary management
│       ├── billing/               # Billing & invoicing
│       ├── claims/                # Claims processing
│       ├── commission/            # Commission & payouts
│       ├── common/                # Common schemas & types
│       ├── document/              # Document management
│       ├── endorsement/           # Policy endorsements
│       ├── fraud/                 # Fraud detection
│       ├── insurance/             # Insurance core
│       ├── insurer/               # Insurer management
│       ├── iot/                   # IoT device integration
│       ├── kyc/                   # KYC verification
│       ├── media/                 # Media/document storage
│       ├── mfs/                   # Mobile Financial Services
│       ├── notification/          # Notifications & alerts
│       ├── orders/                # Order management
│       ├── partner/               # Partner management
│       ├── payment/               # Payment processing
│       ├── policy/                # Policy lifecycle
│       ├── products/              # Product catalog
│       ├── refund/                # Refund processing
│       ├── renewal/               # Policy renewal
│       ├── report/                # Reporting & dashboards
│       ├── services/              # Service providers
│       ├── storage/               # Storage services
│       ├── support/               # Support & ticketing
│       ├── task/                  # Task management
│       ├── tenant/                # Multi-tenancy
│       ├── underwriting/          # Underwriting decisions
│       ├── voice/                 # Voice services
│       └── webrtc/                # WebRTC communication
│
├── rules/                          # API Standards & Guidelines
│   ├── 00-index.md                # Index of all rules
│   ├── 01-response-envelope.md    # Standard response format
│   ├── 02-http-status-codes.md    # HTTP status code rules
│   ├── 03-error-handling.md       # Error response standards
│   ├── 04-security-authentication.md  # Auth & security
│   ├── 05-pagination-and-lists.md # Pagination standards
│   ├── 06-dependency-injection-and-testing.md
│   ├── 07-naming-and-url-design.md    # URL & naming conventions
│   ├── 08-null-empty-and-optional-data.md
│   ├── 09-generator-fix-plan.md
│   ├── dbrules.md
│   └── ground_truth.md
│
├── sdks/
│   ├── insuretech-go-sdk/         # Go SDK (generated)
│   ├── insuretech-typescript-sdk/ # TypeScript SDK (generated)
│   │   ├── src/
│   │   │   ├── client.gen.ts
│   │   │   ├── errors.ts
│   │   │   ├── types.gen.ts
│   │   │   ├── client/
│   │   │   │   ├── client.gen.ts
│   │   │   │   ├── types.gen.ts
│   │   │   │   └── utils.gen.ts
│   │   │   └── core/
│   │   │       ├── auth.gen.ts
│   │   │       ├── bodySerializer.gen.ts
│   │   │       ├── params.gen.ts
│   │   │       ├── pathSerializer.gen.ts
│   │   │       └── serverSentEvents.gen.ts
│   │   └── tests/
│   └── sdk-generator/             # SDK generator scripts
│       ├── go/
│       └── typescript/
│
├── proto/                         # Proto compilation & generation
│   ├── check_migrations.py
│   └── DUPLICATE_MESSAGES_CLEANUP_SUMMARY.md
│
├── api/                           # REST API implementations
├── backend/                       # Backend services
├── docs/                          # Documentation
├── ops/                           # Operations/DevOps
├── scripts/                       # Utility scripts
│
├── buf.yaml                       # Buf protobuf configuration
├── buf.lock                       # Buf lock file
├── buf.gen.yaml                   # Buf code generator config
├── README.md                      # Main README
├── START_HERE.md                  # Quick start guide
└── go.work                        # Go workspace config
```

---

## Proto Services (34 Total)

Each service follows the pattern:
```
{service}/
├── entity/v1/        # Data models
├── events/v1/        # Event definitions
└── services/v1/      # Service/RPC definitions
```

### Complete Service List:
1. **AI Service** - Agent/AI entity interactions
2. **Analytics** - Analytics, metrics, dashboards
3. **ApiKey** - API key management & usage tracking
4. **Audit** - Audit events, logs, compliance
5. **Authentication** - User auth, OTP, sessions
6. **Authorization** - RBAC, roles, policies
7. **B2B** - B2B org, employee, department, PO management
8. **Beneficiary** - Beneficiary data (individual/business)
9. **Billing** - Invoice generation & management
10. **Claims** - Claim submission & lifecycle
11. **Commission** - Commission config & payouts
12. **Document** - Document generation & templates
13. **Endorsement** - Policy endorsements
14. **Fraud** - Fraud alerts, cases, rules
15. **Insurance** - Core insurance service
16. **Insurer** - Insurer config & products
17. **IoT** - IoT device management
18. **KYC** - Know Your Customer verification
19. **Media** - Media storage & processing
20. **MFS** - Mobile Financial Services integration
21. **Notification** - Alerts & notifications
22. **Orders** - Order management
23. **Partner** - Partner management
24. **Payment** - Payment processing
25. **Policy** - Policy lifecycle management
26. **Products** - Product catalog & plans
27. **Refund** - Refund processing
28. **Renewal** - Policy renewal management
29. **Report** - Reporting & scheduling
30. **Services** - Service provider management
31. **Storage** - File/blob storage
32. **Support** - Support tickets, FAQs, knowledge base
33. **Task** - Task management
34. **Tenant** - Multi-tenancy management
35. **Underwriting** - Underwriting decisions
36. **Voice** - Voice sessions & commands
37. **WebRTC** - Real-time communication

(Note: List exceeds 34 — actual count in codebase)

---

## API Rules & Standards (Critical)

### Rule 01: Standard Response Envelope ⭐ CRITICAL
**Every API response must use this envelope:**

```json
{
  "success": true,
  "data": { },
  "error": null,
  "meta": {
    "request_id": "req_123",
    "pagination": null
  }
}
```

**Key Points:**
- `success`: Boolean indicating operation success
- `data`: Response payload (only on success, never on error)
- `error`: Error object (only on failure, always null on success)
- `meta`: Metadata with request_id and optional pagination
- **NO error field inside success response schemas** (forbidden pattern)

---

### Rule 02: HTTP Status Codes ⭐ CRITICAL

**Decision Tree:**
```
POST /resource          → 201 Created (creates new resource)
POST /resource/:action  → 200 OK (custom action, not creating)
GET /resource           → 200 OK
PUT /resource/{id}      → 200 OK (returns updated resource)
PATCH /resource/{id}    → 200 OK (returns updated resource)
DELETE /resource/{id}   → 204 No Content (no body)
```

**Full Reference:**

| Code | Name | When to Use |
|------|------|------------|
| **200** | OK | GET success, action POST, PUT/PATCH success |
| **201** | Created | POST that creates new persistent resource |
| **204** | No Content | DELETE or actions returning nothing |
| **400** | Bad Request | Malformed JSON, wrong types |
| **401** | Unauthorized | No token, expired token |
| **403** | Forbidden | Valid token, insufficient permissions |
| **404** | Not Found | Resource ID doesn't exist |
| **409** | Conflict | Resource exists, version conflict |
| **422** | Unprocessable Entity | Valid JSON but validation failed |
| **429** | Too Many Requests | Rate limit exceeded |
| **500** | Internal Server Error | Unexpected server error |

**Important Distinctions:**
- **400**: Structural problem (malformed JSON)
- **422**: Business logic problem (validation failed)
- **201**: Must return `Location` header
- **204**: No body, no envelope
- Empty lists return **200 OK** (not 404)

---

### Rule 03: Error Handling ⭐ CRITICAL

**Golden Rule:** Errors NEVER live inside success response schemas. They come exclusively via HTTP 4xx/5xx with standard error envelope.

**Standard Error Schema:**
```yaml
Error:
  type: object
  required:
    - code
    - message
  properties:
    code:                    # UPPER_SNAKE_CASE error code
      type: string
    message:                 # Human-readable, localized message
      type: string
    field_violations:        # For 422: per-field validation errors
      type: array
    error_id:                # Unique error instance ID
      type: string
    retryable:               # Whether client should retry
      type: boolean
    retry_after_seconds:     # Wait time if retryable
      type: integer
    documentation_url:       # Link to error docs
      type: string
    http_status_code:        # Mirror HTTP status
      type: integer
```

**Error Code Naming Convention:** `DOMAIN_SPECIFIC_CODE`
```
AUTH_INVALID_CREDENTIALS
AUTH_TOKEN_EXPIRED
POLICY_NOT_FOUND
PAYMENT_INSUFFICIENT_FUNDS
CLAIM_ALREADY_SETTLED
KYC_VERIFICATION_FAILED
VALIDATION_FAILED
RESOURCE_NOT_FOUND
DUPLICATE_RESOURCE
RATE_LIMIT_EXCEEDED
INTERNAL_ERROR
```

**Required Error Responses Per Endpoint:**

| Code | Required For |
|------|-------------|
| `400` | All endpoints |
| `401` | All authenticated endpoints |
| `403` | All authorized endpoints |
| `404` | Endpoints with path parameters |
| `409` | POST endpoints (duplicates) |
| `422` | All POST/PUT/PATCH endpoints |
| `429` | All public-facing endpoints |
| `500` | All endpoints |

---

### Rule 04: Security & Authentication ⭐ CRITICAL

**Security Classification:**

**🌍 PUBLIC Endpoints** (no auth):
```yaml
/v1/auth/register
/v1/auth/login
/v1/auth/email/register
/v1/auth/otp:send
/v1/auth/otp:verify
/v1/products          # public product browsing
```

**🔐 BEARER AUTH** (JWT token):
```yaml
# All protected endpoints require:
Authorization: Bearer <JWT_TOKEN>
```

**🔑 API KEY** (B2B/partner):
```yaml
X-API-Key: ak_live_abc123xyz
```

**🏢 TENANT-SCOPED**:
```yaml
Authorization: Bearer <JWT_TOKEN>
X-Tenant-ID: tenant_abc123
```

**Token Lifecycle:**
```
Register/Login → access_token (15 min) + refresh_token (30 days)
                ↓
        Use access_token for API calls
                ↓
        Token expires → POST /v1/auth/token:refresh
                ↓
        New access_token issued
                ↓
        Logout → POST /v1/auth/logout
```

**Global Security Definition (openapi.yaml):**
```yaml
security:
  - BearerAuth: []   # Default for all endpoints

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
    TenantAuth:
      type: apiKey
      in: header
      name: X-Tenant-ID
```

---

### Rule 05: Pagination & Lists ⭐ HIGH

**Request Parameters:**
```yaml
- name: page
  in: query
  schema:
    type: integer
    minimum: 1
    default: 1
  description: Page number (1-based)

- name: page_size
  in: query
  schema:
    type: integer
    minimum: 1
    maximum: 100
    default: 20

- name: sort_by
  in: query
  schema:
    type: string
  description: Field to sort by

- name: sort_order
  in: query
  schema:
    type: string
    enum: [asc, desc]
    default: desc

- name: search
  in: query
  schema:
    type: string
```

**Response Shape:**
```json
{
  "success": true,
  "data": {
    "items": [ { ... }, { ... } ]
  },
  "error": null,
  "meta": {
    "request_id": "req_abc",
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_pages": 5,
      "total_items": 98,
      "has_next": true,
      "has_previous": false
    }
  }
}
```

**Canonical PaginationMeta Schema:**
```yaml
PaginationMeta:
  type: object
  required:
    - page
    - page_size
    - total_pages
    - total_items
    - has_next
    - has_previous
  properties:
    page:
      type: integer
      description: Current page (1-based)
    page_size:
      type: integer
    total_pages:
      type: integer
    total_items:
      type: integer
      format: int64
    has_next:
      type: boolean
    has_previous:
      type: boolean
    next_page_token:
      type: string
      nullable: true
      description: Optional cursor token
```

**Key Rules:**
- All lists return `data.items` array (never bare array)
- Empty list returns **200 OK** (not 404)
- Filtering via query parameters (never POST body for GET)
- Two competing schemas must be merged into ONE standard

---

### Rule 07: URL Design & Naming Standards ⭐ HIGH

**1. Resources are Nouns:**
```
✅ POST /v1/policies
❌ POST /v1/createPolicy
```

**2. Custom Actions Use Colon Notation (Google AIP-136):**
```
✅ POST /v1/policies/{id}:cancel
✅ POST /v1/auth/otp:send
✅ POST /v1/products/{id}:activate
❌ POST /v1/policies/cancel/{id}
```

**3. URL Case: kebab-case:**
```
✅ /v1/audit-logs
✅ /v1/api-keys
✅ /v1/payment-methods
❌ /v1/auditLogs
❌ /v1/audit_logs
```

**4. Path Parameters: snake_case:**
```
✅ /v1/policies/{policy_id}
✅ /v1/users/{user_id}/policies
❌ /v1/policies/{policyId}
```

**5. Query Parameters: snake_case:**
```
✅ ?page=1&page_size=20&sort_by=created_at&sort_order=desc
❌ ?pageSize=20&fromDate=2024-01-01
```

**6. JSON Fields: snake_case:**
```json
✅ {
  "policy_id": "pol_123",
  "policy_number": "INS-2024-001",
  "sum_insured": 500000,
  "created_at": "2024-01-15T10:30:00Z"
}
```

**7. Versioning: Always in URL Path:**
```
✅ /v1/policies
✅ /v2/policies   (breaking changes only)
❌ /policies
```

**8. ID Formats (Prefixed):**
```
pol_abc123      → Policy
clm_abc123      → Claim
pay_abc123      → Payment
ord_abc123      → Order
usr_abc123      → User
prd_abc123      → Product
tnt_abc123      → Tenant
kyc_abc123      → KYC verification
ak_live_abc123  → API key (live)
ak_test_abc123  → API key (test)
```

**9. Date/Time Standards:**

| Type | Format | Example |
|------|--------|---------|
| Timestamp | ISO 8601 UTC | `"2024-01-15T10:30:00Z"` |
| Date | ISO 8601 | `"2024-01-15"` |
| Duration | ISO 8601 | `"P1Y"`, `"P30D"` |
| Currency | String | `"5000.50"` |
| Currency Code | ISO 4217 | `"BDT"`, `"USD"` |

**10. operationId Convention:**
```
{ServiceName}_{MethodName}

✅ PolicyService_CreatePolicy
✅ AuthService_Login
✅ PaymentService_InitiatePayment
```

**11. Resource Naming Reference:**

| Resource | Collection | Single | Action |
|----------|-----------|--------|--------|
| Policy | `/v1/policies` | `/v1/policies/{policy_id}` | `/v1/policies/{policy_id}:cancel` |
| Claim | `/v1/claims` | `/v1/claims/{claim_id}` | `/v1/claims/{claim_id}:approve` |
| Payment | `/v1/payments` | `/v1/payments/{payment_id}` | `/v1/payments/{payment_id}:verify` |
| Order | `/v1/orders` | `/v1/orders/{order_id}` | `/v1/orders/{order_id}:cancel` |
| Product | `/v1/products` | `/v1/products/{product_id}` | `/v1/products/{product_id}:activate` |
| User | `/v1/users` | `/v1/users/{user_id}` | `/v1/users/{user_id}/sessions:revoke-all` |
| Quote | `/v1/quotes` | `/v1/quotes/{quote_id}` | `/v1/quotes/{quote_id}:approve` |
| Ticket | `/v1/tickets` | `/v1/tickets/{ticket_id}` | `/v1/tickets/{ticket_id}:assign` |
| KYC | `/v1/kyc-verifications` | `/v1/kyc-verifications/{kyc_id}` | `/v1/kyc-verifications/{kyc_id}:approve` |
| Invoice | `/v1/invoices` | `/v1/invoices/{invoice_id}` | `/v1/invoices/{invoice_id}:cancel` |

**12. Nested Resources (max 2 levels):**
```
✅ /v1/users/{user_id}/policies
✅ /v1/policies/{policy_id}/claims
❌ /v1/users/{user_id}/policies/{policy_id}/claims/{claim_id}/documents
✅ /v1/claims/{claim_id}/documents   (flatten after 2 levels)
```

---

## TypeScript SDK Structure

```
insuretech-typescript-sdk/
├── src/
│   ├── client.gen.ts           # Auto-generated main client
│   ├── errors.ts               # Error handling
│   ├── types.gen.ts            # Generated types
│   ├── index.ts                # Exports
│   ├── client/
│   │   ├── client.gen.ts       # Client implementation
│   │   ├── types.gen.ts        # Client-specific types
│   │   └── utils.gen.ts        # Utility functions
│   └── core/
│       ├── auth.gen.ts         # Auth handling
│       ├── bodySerializer.gen.ts
│       ├── params.gen.ts       # Parameter handling
│       ├── pathSerializer.gen.ts
│       ├── queryKeySerializer.gen.ts
│       └── serverSentEvents.gen.ts
├── tests/
│   ├── setup.ts
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── tsconfig.json
├── vitest.config.ts
├── package.json
└── README.md
```

---

## Key Configuration Files

- **buf.yaml** - Protobuf compiler configuration
- **buf.lock** - Proto dependency lock file
- **buf.gen.yaml** - Code generator targets
- **go.mod / go.work** - Go module management
- **START_HERE.md** - Quick start guide (references rules/00-index.md)

---

## Important Notes for Developers

### The Standard Envelope is Non-Negotiable
The response envelope is **CRITICAL** (🔴) and applies to ALL 34 services. Clients expect:
- One unified response shape across all endpoints
- Consistent error handling via envelope
- Type-safe SDK generation

### 201 vs 200 Decision Tree
- **201 Created**: `POST` endpoints that create NEW PERSISTENT resources
  - Examples: registration, creating policies, submitting claims
  - Must return `Location` header
  
- **200 OK**: Everything else
  - Action endpoints (`:verify`, `:cancel`, `:approve`)
  - GET, PUT, PATCH
  - Login (returns session/token, not a new resource)

### Security is Per-Endpoint
Even with global security, every endpoint must declare its security explicitly via `security:` key in OpenAPI YAML.

### Common Implementation Pitfalls
1. ❌ Error field embedded in success response → will break SDK generators
2. ❌ Returning bare arrays instead of `data.items` → breaks future extensibility
3. ❌ Using 200 for resource creation → confuses clients about what was created
4. ❌ Field names in camelCase → inconsistent with REST conventions
5. ❌ Missing Location header on 201 → clients can't find created resource
6. ❌ No explicit security declaration → SDKs don't know what auth to use

---

## Next Steps

1. **Review** `rules/00-index.md` for complete rules index
2. **Check** `API_PIPELINE_REFERENCE.md` for pipeline flow & artifacts
3. **Review** `API_PIPELINE_STATUS.md` for current implementation gaps
4. **Examine proto files** in `proto/insuretech/` for service definitions
5. **Review SDK generators** in `sdks/sdk-generator/` for compliance

---

**Last Updated:** 2024  
**Scope:** InsureTech Platform API Standards & Guidelines
