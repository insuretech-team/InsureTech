# Rule 02: HTTP Status Code Standards

**Scope:** All REST API endpoints across all 34 services  
**Priority:** 🔴 CRITICAL

---

## The Complete Status Code Decision Tree

```
POST /resource          → Creating a new persistent resource?
  YES → 201 Created
  NO  → 200 OK  (action, query, trigger)

GET /resource           → 200 OK
GET /resource/{id}      → 200 OK (or 404 if not found)

PUT /resource/{id}      → 200 OK (full replace, returns updated resource)
PATCH /resource/{id}    → 200 OK (partial update, returns updated resource)

DELETE /resource/{id}   → 204 No Content (no body)

POST /resource/{id}:action  → 200 OK (custom action, not creating a resource)
```

---

## Full Status Code Reference

| Code | Name | When to Use |
|------|------|------------|
| **200** | OK | Successful GET, successful action POST, successful PUT/PATCH |
| **201** | Created | POST that creates a new persistent resource (returns the created resource) |
| **204** | No Content | DELETE, or actions that return nothing |
| **400** | Bad Request | Malformed JSON, missing required fields, wrong types |
| **401** | Unauthorized | No token, expired token, invalid token |
| **403** | Forbidden | Valid token but insufficient permissions |
| **404** | Not Found | Resource with given ID does not exist |
| **409** | Conflict | Resource already exists (duplicate), version conflict |
| **422** | Unprocessable Entity | Valid JSON but business validation failed |
| **429** | Too Many Requests | Rate limit exceeded |
| **500** | Internal Server Error | Unexpected server error |
| **503** | Service Unavailable | Downstream service unavailable |

---

## 200 vs 201 — Example Endpoint Classification

### Must Return 201 (Creating a new persistent resource)

| Example | Endpoint | Required Code |
|---------|----------|---------------|
| Create payment | `POST /v1/payments` | `201` |
| Submit claim | `POST /v1/claims` | `201` |
| Request quote | `POST /v1/quotes` | `201` |
| Add payment method | `POST /v1/users/{id}/payment-methods` | `201` |
| Start KYC | `POST /v1/users/{id}/kyc` | `201` |
| Start workflow | `POST /v1/workflow-instances` | `201` |
| Upload document | `POST /v1/users/{id}/documents` | `201` |
| Generate API key | `POST /v1/api-keys` | `201` |
| Upload media | `POST /v1/media` | `201` |
| Request endorsement | `POST /v1/endorsements` | `201` |

### Representative 201 endpoints

| Endpoint | Status |
|----------|--------|
| `POST /v1/auth/register` | ✅ 201 |
| `POST /v1/auth/email/register` | ✅ 201 |
| `POST /v1/auth/api-keys` | ✅ 201 |
| `POST /v1/auth/users/{id}/profile` | ✅ 201 |
| `POST /v1/auth/voice-sessions` | ✅ 201 |
| `POST /v1/authz/roles` | ✅ 201 |
| `POST /v1/authz/policies` | ✅ 201 |
| `POST /v1/policies` | ✅ 201 |
| `POST /v1/products` | ✅ 201 |
| `POST /v1/orders` | ✅ 201 |
| `POST /v1/tenants` | ✅ 201 |
| `POST /v1/tasks` | ✅ 201 |
| `POST /v1/tickets` | ✅ 201 |
| `POST /v1/partners` | ✅ 201 |
| `POST /v1/audit-logs` | ✅ 201 |
| `POST /v1/audit-events` | ✅ 201 |
| `POST /v1/compliance-logs` | ✅ 201 |
| `POST /v1/analytics/dashboards` | ✅ 201 |
| `POST /v1/report-schedules` | ✅ 201 |
| `POST /v1/workflow-definitions` | ✅ 201 |
| `POST /v1/faqs` | ✅ 201 |
| `POST /v1/knowledge-base` | ✅ 201 |
| `POST /v1/invoices` (BillingService) | ✅ 201 |
| `POST /v1/commission/payouts` | ✅ 201 |
| `POST /v1/fraud-cases` | ✅ 201 |
| `POST /v1/fraud-rules` | ✅ 201 |

### Must Return 200 (Action endpoints — correct)

| Endpoint | Reason |
|----------|--------|
| `POST /v1/auth/login` | Returns session/token, not a new resource |
| `POST /v1/auth/logout` | Action — revoke session |
| `POST /v1/auth/otp:send` | Action — triggers OTP delivery |
| `POST /v1/auth/otp:verify` | Action — validates OTP |
| `POST /v1/auth/token:refresh` | Action — issues new token |
| `POST /v1/auth/password:change` | Action — mutates existing resource |
| `POST /v1/policies/{id}:cancel` | Action on existing resource |
| `POST /v1/policies/{id}:renew` | Action on existing resource |
| `POST /v1/policies/{id}:issue` | Action on existing resource |
| `POST /v1/payments/{id}:verify` | Action on existing resource |
| `POST /v1/payments/{id}:review` | Action on existing resource |
| `POST /v1/claims/{id}:approve` | Action on existing resource |
| `POST /v1/claims/{id}:reject` | Action on existing resource |
| `POST /v1/products/{id}:activate` | Action on existing resource |
| `POST /v1/ai/chat` | Query/compute action |
| `POST /v1/authz/check` | Query action |
| `POST /v1/analytics/queries:run` | Query action |

---

## 400 vs 422 — Mandatory Distinction

**Both MUST be present on every endpoint.**

| Code | Trigger | Example |
|------|---------|---------|
| **400** | Malformed request — cannot parse | Invalid JSON body, wrong Content-Type |
| **422** | Valid request — business rule failed | Email already exists, premium below minimum, NID invalid |

```json
// 400 Bad Request — structural problem
{
  "success": false,
  "data": null,
  "error": {
    "code": "MALFORMED_REQUEST",
    "message": "Request body is not valid JSON",
    "retryable": false
  }
}

// 422 Unprocessable Entity — validation problem
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "One or more fields are invalid",
    "field_violations": [
      { "field": "email", "message": "Email is already registered" },
      { "field": "phone", "message": "Must be a valid BD mobile number" }
    ],
    "retryable": false
  }
}
```

---

## 201 Created — Required Response Headers

When returning 201, the response MUST include a `Location` header:

```
HTTP/1.1 201 Created
Location: /v1/policies/pol_abc123
Content-Type: application/json

{
  "success": true,
  "data": { "policy_id": "pol_abc123", ... },
  "error": null
}
```

---

## 204 No Content — No Body Rule

```
HTTP/1.1 204 No Content
```
**No body. No envelope. Empty response.**

Use for:
- `DELETE` operations
- Fire-and-forget actions where the result is obvious

Do NOT use 204 when the consumer needs confirmation data back.  
Example: `POST /v1/auth/logout` returns 200 with `{ "success": true, "data": null }` — not 204 — because clients may want to confirm session was terminated.

---

## Generator Fix Required (path_generator.py)

The current logic only returns 201 for `Create*` and `Register*` method names.

```python
# CURRENT (WRONG) - too narrow
if method_name.startswith('Create') or method_name.startswith('Register'):
    return '201'

# REQUIRED - classify by semantic intent
CREATES_RESOURCE = [
    'Create', 'Register', 'Submit', 'Initiate', 'Request',
    'Add', 'Upload', 'Start', 'Generate', 'Issue', 'Open', 'File'
]
ACTION_VERBS = [
    'Verify', 'Validate', 'Approve', 'Reject', 'Cancel', 'Revoke',
    'Activate', 'Deactivate', 'Assign', 'Send', 'Check', 'Calculate',
    'Process', 'Handle', 'Complete', 'Confirm', 'Run', 'Execute',
    'Rotate', 'Refresh', 'Login', 'Logout', 'Change', 'Reset',
    'Renew', 'Settle', 'Reconcile', 'Review', 'Discontinue', 'Mark'
]

def get_success_code(method_name, http_method, path):
    # DELETE always 204
    if http_method == 'delete':
        return '204'
    # Custom actions (:verb suffix in path) always 200
    if ':' in path:
        return '200'
    # POST to collection path without action = create = 201
    if http_method == 'post':
        for verb in CREATES_RESOURCE:
            if method_name.startswith(verb):
                return '201'
    return '200'
```
