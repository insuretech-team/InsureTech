# Rule 07: URL Design & Naming Standards

**Scope:** All REST API endpoints  
**Priority:** 🟠 HIGH

---

## URL Design Principles

### 1. Resources are Nouns, Never Verbs

```
✅ POST /v1/policies          (create)
✅ GET  /v1/policies          (list)
✅ GET  /v1/policies/{id}     (get one)
✅ PATCH /v1/policies/{id}    (update)
✅ DELETE /v1/policies/{id}   (delete)

❌ POST /v1/createPolicy
❌ GET  /v1/getPolicy
❌ POST /v1/deletePolicy
```

### 2. Custom Actions Use Colon Notation (Google AIP-136)

```
✅ POST /v1/policies/{id}:cancel
✅ POST /v1/policies/{id}:renew
✅ POST /v1/policies/{id}:issue
✅ POST /v1/auth/otp:send
✅ POST /v1/auth/otp:verify
✅ POST /v1/products/{id}:activate
✅ POST /v1/products/{id}:deactivate
✅ POST /v1/payments/{id}:verify
✅ POST /v1/payments/{id}:submit-proof

❌ POST /v1/policies/cancel/{id}
❌ PUT  /v1/policies/{id}/cancel
❌ POST /v1/cancelPolicy/{id}
```

### 3. URL Case: kebab-case for Paths

```
✅ /v1/audit-logs
✅ /v1/api-keys
✅ /v1/payment-methods
✅ /v1/knowledge-base
✅ /v1/report-schedules
✅ /v1/workflow-definitions
✅ /v1/voice-sessions

❌ /v1/auditLogs
❌ /v1/audit_logs
❌ /v1/AuditLogs
```

### 4. Path Parameters: snake_case

```
✅ /v1/policies/{policy_id}
✅ /v1/users/{user_id}/policies
✅ /v1/claims/{claim_id}
✅ /v1/payments/{payment_id}

❌ /v1/policies/{policyId}
❌ /v1/policies/{id}      (too generic — use resource name)
```

### 5. Query Parameters: snake_case

```
✅ ?page=1&page_size=20&sort_by=created_at&sort_order=desc
✅ ?from_date=2024-01-01&to_date=2024-03-31
✅ ?user_id=usr_123&product_id=prod_456

❌ ?pageSize=20
❌ ?fromDate=2024-01-01
```

### 6. Versioning: Always in URL Path

```
✅ /v1/policies
✅ /v2/policies   (breaking changes only)

❌ /policies                    (no version)
❌ Accept: application/vnd.api+json;version=1  (header versioning)
```

---

## Resource Naming Reference

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

---

## Nested Resources — When to Nest

**Nest when resource only makes sense in context of parent:**

```
✅ /v1/users/{user_id}/policies           (a user's policies)
✅ /v1/users/{user_id}/payment-methods    (a user's payment methods)
✅ /v1/policies/{policy_id}/claims        (claims under a policy)
✅ /v1/tickets/{ticket_id}/messages       (messages in a ticket)
✅ /v1/quotes/{quote_id}/health-declaration
```

**Do NOT nest more than 2 levels deep:**

```
❌ /v1/users/{user_id}/policies/{policy_id}/claims/{claim_id}/documents
✅ /v1/claims/{claim_id}/documents    (flatten after 2 levels)
```

---

## JSON Field Naming: snake_case

All request and response JSON fields use `snake_case`:

```json
✅ {
  "policy_id": "pol_123",
  "policy_number": "INS-2024-001",
  "sum_insured": 500000,
  "premium_amount": { "amount": "5000", "currency": "BDT" },
  "effective_date": "2024-01-15",
  "expiry_date": "2025-01-14",
  "created_at": "2024-01-15T10:30:00Z"
}

❌ {
  "policyId": "pol_123",
  "PolicyNumber": "INS-2024-001",
  "SumInsured": 500000
}
```

---

## ID Format Standards

All IDs MUST be prefixed with resource type:

```
pol_abc123      → Policy
clm_abc123      → Claim
pay_abc123      → Payment
ord_abc123      → Order
usr_abc123      → User
prd_abc123      → Product
tnt_abc123      → Tenant
kyc_abc123      → KYC verification
inv_abc123      → Invoice
tkt_abc123      → Support ticket
quot_abc123     → Quote
ak_live_abc123  → API key (live)
ak_test_abc123  → API key (test)
```

---

## Date & Time Standards

| Type | Format | Example |
|------|--------|---------|
| Timestamp | ISO 8601 UTC | `"2024-01-15T10:30:00Z"` |
| Date only | ISO 8601 | `"2024-01-15"` |
| Duration | ISO 8601 | `"P1Y"` (1 year), `"P30D"` (30 days) |
| Currency amount | string (avoid float precision) | `"5000.50"` |
| Currency code | ISO 4217 | `"BDT"`, `"USD"` |

```json
✅ {
  "created_at": "2024-01-15T10:30:00Z",
  "effective_date": "2024-01-15",
  "duration": "P1Y",
  "premium": { "amount": "5000.50", "currency": "BDT" }
}

❌ {
  "created_at": 1705312200,        // unix timestamp — not human readable
  "premium": 5000.50               // float precision issues
}
```

---

## operationId Convention

Format: `{ServiceName}_{MethodName}`

```
✅ PolicyService_CreatePolicy
✅ AuthService_Login
✅ PaymentService_InitiatePayment
✅ ClaimService_SubmitClaim

❌ createPolicy
❌ POST_policies
❌ policy-create
```

---

## Tags Convention

Each service = one tag. One tag per endpoint.

```yaml
tags:
  - name: AuthService
    description: Authentication & session management
  - name: PolicyService
    description: Insurance policy lifecycle
  - name: PaymentService
    description: Payment processing & methods
  - name: ClaimService
    description: Claims submission & management
```
