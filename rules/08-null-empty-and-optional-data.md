# Rule 08: Null, Empty & Optional Data Standards

**Scope:** All API request and response schemas  
**Priority:** 🟠 HIGH

---

## Core Requirement

Success responses must use structured data with explicit required and nullable behavior.

---

## Core Rules

### Rule 8.1 — Never Omit Fields, Use Explicit Null

**Always include declared fields in the response. If a value is absent, send `null` explicitly — never omit the field.**

```json
// ✅ CORRECT — field present, value null
{
  "success": true,
  "data": {
    "policy_id": "pol_123",
    "cancellation_reason": null,
    "cancelled_at": null
  }
}

// ❌ WRONG — field omitted entirely
{
  "success": true,
  "data": {
    "policy_id": "pol_123"
    // cancellation_reason and cancelled_at missing — client crashes
  }
}
```

**Why:** Omitting fields causes `undefined` in TypeScript, nil crashes in Swift, NullPointerExceptions in Kotlin.

---

### Rule 8.2 — Action Endpoints May Return Null Data

For action endpoints that mutate state and have no meaningful return value, `data` MUST be `null` explicitly — NOT an empty object `{}` or a schema with only a `message` field.

```json
// ✅ CORRECT — action with no return data
POST /v1/auth/logout → 200 OK
{
  "success": true,
  "data": null,
  "error": null,
  "meta": { "request_id": "req_abc" }
}

// ✅ CORRECT — action with confirmation data
POST /v1/policies/{id}:cancel → 200 OK
{
  "success": true,
  "data": {
    "policy_id": "pol_123",
    "status": "CANCELLED",
    "cancelled_at": "2024-01-15T10:30:00Z",
    "refund_amount": { "amount": "1500.00", "currency": "BDT" }
  },
  "error": null,
  "meta": { "request_id": "req_abc" }
}

// ❌ WRONG — thin response with only message
{
  "success": true,
  "data": {
    "message": "Policy cancelled successfully"
  }
}
```

**Reasoning:** A `message` field in `data` is useless for programmatic use. Clients need structured data. Messages belong in `meta.message` if needed at all.

---

### Rule 8.3 — Required vs Optional Field Declaration

Every schema field MUST be marked as either `required` or explicitly `nullable: true` in the OpenAPI spec.

```yaml
PolicyData:
  type: object
  required:
    - policy_id
    - policy_number
    - status
    - product_id
    - effective_date
    - expiry_date
    - premium_amount
    - sum_insured
  properties:
    policy_id:
      type: string
    policy_number:
      type: string
    status:
      type: string
      enum: [ACTIVE, CANCELLED, EXPIRED, PENDING, SUSPENDED]
    cancellation_reason:
      type: string
      nullable: true          # ← explicitly nullable
      description: Set only when status=CANCELLED
    cancelled_at:
      type: string
      format: date-time
      nullable: true          # ← explicitly nullable
    renewal_policy_id:
      type: string
      nullable: true          # ← set only after renewal
```

---

### Rule 8.4 — Classify All Action Endpoints by Return Type

| Category | data value | Example endpoints |
|----------|-----------|-------------------|
| **Creates resource** | Full resource object | POST /v1/policies, POST /v1/claims |
| **Returns computed result** | Result object | POST /v1/products/{id}:calculate-premium |
| **Mutates + confirms** | Updated resource | POST /v1/policies/{id}:cancel |
| **Triggers side effect** | null | POST /v1/auth/logout, POST /v1/notifications/send |
| **Validates/checks** | Boolean/decision | POST /v1/authz/check |

---

### Rule 8.5 — Action Endpoint Patterns That Require Structured Data

These endpoints show the kinds of thin or ambiguous payloads that must be replaced with proper structured data:

| Endpoint | Thin shape | Required data fields |
|----------|---------------|---------------------|
| `POST /v1/auth/logout` | `LogoutResponse { message, error }` | `data: null` |
| `POST /v1/auth/password:change` | `ChangePasswordResponse { message, error }` | `data: null` |
| `POST /v1/auth/password:reset` | `ResetPasswordResponse { message, error }` | `data: null` |
| `POST /v1/auth/otp:send` | `OTPSendingResponse { otp_id, message, error }` | `data: { otp_id, expires_in_seconds }` |
| `POST /v1/auth/otp:verify` | `OTPVerificationResponse { message, error }` | `data: null` |
| `POST /v1/policies/{id}:cancel` | `PolicyCancellationResponse { message, error }` | `data: { policy_id, status, cancelled_at, refund_amount }` |
| `POST /v1/policies/{id}:issue` | `PolicyIssuanceResponse { message, error }` | `data: { policy_id, status, issued_at, document_url }` |
| `POST /v1/claims/{id}:approve` | `ClaimApprovalResponse { message, error }` | `data: { claim_id, status, approved_at, settlement_amount }` |
| `POST /v1/claims/{id}:reject` | `ClaimRejectionResponse { message, error }` | `data: { claim_id, status, rejected_at, reason }` |
| `POST /v1/auth/users/{id}/sessions:revoke-all` | `RevokeAllSessionsResponse { message, error }` | `data: { revoked_count }` |
| `POST /v1/products/{id}:activate` | `ProductActivationResponse { message, error }` | `data: { product_id, status, activated_at }` |
| `POST /v1/payments/{id}:verify` | `PaymentVerificationResponse { message, error }` | `data: { payment_id, status, verified_at }` |
| `PATCH /v1/authz/portals/{portal}/config` | `PortalConfigUpdateResponse { message, error }` | `data: { portal, config, updated_at }` |

---

### Rule 8.6 — No `message` Field in Data

The `message` string field in response data is an anti-pattern for programmatic APIs:

```json
// ❌ WRONG — message in data
{
  "data": {
    "message": "OTP sent successfully to your phone"
  }
}

// ✅ CORRECT — structured data, no message
{
  "data": {
    "otp_id": "otp_abc123",
    "expires_in_seconds": 300,
    "masked_phone": "+880 01***6789"
  }
}
```

**Exception:** `meta.message` is acceptable for informational hints that don't affect logic.

---

### Rule 8.7 — Boolean vs Null for Status Fields

Use enums for status, not booleans:

```json
// ❌ WRONG
{ "is_active": true, "is_cancelled": false }

// ✅ CORRECT
{ "status": "ACTIVE" }
// Status enum: ACTIVE, INACTIVE, CANCELLED, PENDING, SUSPENDED, EXPIRED
```

---

### Rule 8.8 — Empty String vs Null

Never return empty string `""` where `null` is intended:

```json
// ❌ WRONG
{ "cancellation_reason": "" }

// ✅ CORRECT  
{ "cancellation_reason": null }
```

---

## Schema Template for Action Endpoints

```yaml
# Template for action endpoints that return confirmation
PolicyCancellationData:
  type: object
  required:
    - policy_id
    - status
    - cancelled_at
  properties:
    policy_id:
      type: string
    status:
      type: string
      enum: [CANCELLED]
    cancelled_at:
      type: string
      format: date-time
    effective_cancellation_date:
      type: string
      format: date
      nullable: true
    refund_amount:
      $ref: '#/components/schemas/Money'
      nullable: true
      description: Refund amount if applicable

# Template for side-effect only actions (null data)
# Just use data: null — no schema needed
```
