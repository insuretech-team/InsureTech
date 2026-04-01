# Rule 03: Error Handling Standards

**Scope:** All REST API endpoints across all 34 services  
**Priority:** 🔴 CRITICAL

---

## The Golden Rule

> **Errors NEVER live inside success response schemas.**  
> Errors come exclusively via HTTP 4xx/5xx status codes with the standard error envelope.

---

## Standard Error Schema

```yaml
Error:
  type: object
  required:
    - code
    - message
  properties:
    code:
      type: string
      description: >
        Machine-readable UPPER_SNAKE_CASE error code.
        Used by clients for programmatic error handling.
        Examples: POLICY_NOT_FOUND, EMAIL_ALREADY_REGISTERED,
                  INVALID_PREMIUM_AMOUNT, INSUFFICIENT_PERMISSIONS
    message:
      type: string
      description: >
        Human-readable message in the user's locale.
        Should be clear, actionable, and safe to display directly in UI.
    field_violations:
      type: array
      description: >
        Field-level validation errors. Populated for 422 responses.
        Used by forms to show inline error messages.
      items:
        $ref: '#/components/schemas/FieldViolation'
    error_id:
      type: string
      description: >
        Unique ID for this error instance. Used for support lookups.
        Format: err_{uuid}
    retryable:
      type: boolean
      description: >
        true = client should retry (e.g., network timeout, 503).
        false = retrying will not help (validation error, auth error).
    retry_after_seconds:
      type: integer
      description: Seconds to wait before retrying. Only set when retryable=true.
    documentation_url:
      type: string
      description: Link to docs explaining this error code.
    http_status_code:
      type: integer
      description: Mirrors the HTTP status code for clients that lose HTTP context.

FieldViolation:
  type: object
  required:
    - field
    - message
  properties:
    field:
      type: string
      description: >
        The field path that caused the violation.
        Use dot notation for nested: "address.city"
        Use bracket notation for arrays: "items[0].amount"
    message:
      type: string
      description: Human-readable validation message for this field.
    code:
      type: string
      description: Machine-readable violation code. E.g. REQUIRED, TOO_SHORT, INVALID_FORMAT
    rejected_value:
      type: string
      description: The value that was rejected (stringified).
```

---

## Required Error Responses Per Endpoint

**Every endpoint MUST declare ALL of the following:**

| Code | Required For |
|------|-------------|
| `400` | All endpoints |
| `401` | All authenticated endpoints |
| `403` | All authorized endpoints |
| `404` | Endpoints with path parameters referencing a resource |
| `409` | POST endpoints that may create duplicates |
| `422` | All POST / PUT / PATCH endpoints |
| `429` | All public-facing endpoints |
| `500` | All endpoints |

## Error Code Naming Convention

All error codes MUST be `UPPER_SNAKE_CASE` and domain-prefixed:

```
AUTH_INVALID_CREDENTIALS
AUTH_TOKEN_EXPIRED
AUTH_ACCOUNT_LOCKED
AUTH_OTP_EXPIRED
AUTH_OTP_ALREADY_USED

POLICY_NOT_FOUND
POLICY_ALREADY_CANCELLED
POLICY_RENEWAL_NOT_ALLOWED

PAYMENT_INSUFFICIENT_FUNDS
PAYMENT_GATEWAY_TIMEOUT
PAYMENT_METHOD_NOT_FOUND

CLAIM_ALREADY_SETTLED
CLAIM_DOCUMENT_MISSING
CLAIM_OUTSIDE_COVERAGE

KYC_VERIFICATION_FAILED
KYC_DOCUMENT_EXPIRED
KYC_ALREADY_VERIFIED

VALIDATION_FAILED          # generic for 422
RESOURCE_NOT_FOUND         # generic for 404
DUPLICATE_RESOURCE         # generic for 409
RATE_LIMIT_EXCEEDED        # for 429
INTERNAL_ERROR             # for 500
```

---

## Error Response Examples By Status Code

### 400 Bad Request — Malformed Request
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "MALFORMED_REQUEST",
    "message": "The request body could not be parsed as JSON.",
    "error_id": "err_abc123",
    "retryable": false,
    "http_status_code": 400
  },
  "meta": { "request_id": "req_xyz" }
}
```

### 401 Unauthorized — No/Invalid Token
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "AUTH_TOKEN_EXPIRED",
    "message": "Your session has expired. Please log in again.",
    "error_id": "err_abc124",
    "retryable": false,
    "http_status_code": 401
  },
  "meta": { "request_id": "req_xyz" }
}
```

### 403 Forbidden — Insufficient Permissions
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "INSUFFICIENT_PERMISSIONS",
    "message": "You do not have permission to approve claims.",
    "error_id": "err_abc125",
    "retryable": false,
    "http_status_code": 403
  },
  "meta": { "request_id": "req_xyz" }
}
```

### 404 Not Found
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "POLICY_NOT_FOUND",
    "message": "Policy with ID 'pol_xyz' does not exist.",
    "error_id": "err_abc126",
    "retryable": false,
    "http_status_code": 404
  },
  "meta": { "request_id": "req_xyz" }
}
```

### 409 Conflict — Duplicate Resource
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "AUTH_EMAIL_ALREADY_REGISTERED",
    "message": "An account with this email address already exists.",
    "error_id": "err_abc127",
    "retryable": false,
    "http_status_code": 409
  },
  "meta": { "request_id": "req_xyz" }
}
```

### 422 Unprocessable Entity — Validation Failed
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "One or more fields contain invalid values.",
    "field_violations": [
      {
        "field": "phone_number",
        "message": "Must be a valid Bangladeshi mobile number (e.g. 01XXXXXXXXX)",
        "code": "INVALID_FORMAT",
        "rejected_value": "123456"
      },
      {
        "field": "sum_insured",
        "message": "Sum insured must be between BDT 100,000 and BDT 10,000,000",
        "code": "OUT_OF_RANGE",
        "rejected_value": "50000"
      }
    ],
    "error_id": "err_abc128",
    "retryable": false,
    "http_status_code": 422
  },
  "meta": { "request_id": "req_xyz" }
}
```

### 429 Too Many Requests
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please wait before retrying.",
    "error_id": "err_abc129",
    "retryable": true,
    "retry_after_seconds": 60,
    "http_status_code": 429
  },
  "meta": { "request_id": "req_xyz" }
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An unexpected error occurred. Our team has been notified.",
    "error_id": "err_abc130",
    "retryable": true,
    "retry_after_seconds": 5,
    "documentation_url": "https://docs.insuretech.com/errors/INTERNAL_ERROR",
    "http_status_code": 500
  },
  "meta": { "request_id": "req_xyz" }
}
```

---

## Forbidden Schema Pattern

Remove this pattern from every `*Response` schema:

```yaml
# REMOVE THIS PATTERN from every *Response schema
SomeOperationResponse:
  properties:
    some_field: ...
    error:                      # ← DELETE THIS
      $ref: '#/components/schemas/Error'
    message:                    # ← KEEP or move to data
      type: string
```

The `error` field in a success schema is dead code — it will always be null on 200/201. It confuses SDK generators, breaks type safety, and forces mobile clients to write unnecessary null checks.

---

## Client-Side Error Handling Pattern

```typescript
// TypeScript — ONE interceptor handles ALL errors
axios.interceptors.response.use(
  (response) => response.data.data,  // unwrap data on success
  (error) => {
    const apiError = error.response?.data?.error;
    switch (error.response?.status) {
      case 401: redirectToLogin(); break;
      case 403: showPermissionDenied(); break;
      case 422: showFieldErrors(apiError.field_violations); break;
      case 429: scheduleRetry(apiError.retry_after_seconds); break;
      case 500: showGenericError(apiError.error_id); break;
    }
    return Promise.reject(apiError);
  }
);
```
