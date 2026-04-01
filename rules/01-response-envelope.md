# Rule 01: Standard Response Envelope

**Scope:** All REST API endpoints across all 34 services  
**Applies to:** iOS, Android, Frontend Web, SDK generators, API consumers  
**Priority:** 🔴 CRITICAL

---

## The Rule

**Every API response — success or error — MUST use a standard envelope.**

```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "meta": { ... }
}
```

---

## Envelope Schema

```yaml
ApiResponse:
  type: object
  required:
    - success
  properties:
    success:
      type: boolean
      description: true = operation succeeded, false = operation failed
    data:
      description: >
        Present on success (success=true). 
        MAY be null for operations that return no data (e.g. logout, delete).
        Never present on error responses.
      nullable: true
    error:
      $ref: '#/components/schemas/Error'
      description: >
        Present ONLY on failure (success=false).
        Always null when success=true.
      nullable: true
    meta:
      $ref: '#/components/schemas/ResponseMeta'
      description: >
        Optional metadata. Used for pagination, request tracing, rate limits.
      nullable: true
```

```yaml
ResponseMeta:
  type: object
  properties:
    request_id:
      type: string
      description: Unique request trace ID for debugging
    pagination:
      $ref: '#/components/schemas/PaginationMeta'
      nullable: true
    timestamp:
      type: string
      format: date-time
    api_version:
      type: string

PaginationMeta:
  type: object
  properties:
    page:
      type: integer
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
```

---

## ✅ Correct Examples

### Success with data
```json
HTTP 201 Created
{
  "success": true,
  "data": {
    "policy_id": "pol_abc123",
    "policy_number": "INS-2024-001",
    "status": "ACTIVE"
  },
  "error": null,
  "meta": {
    "request_id": "req_xyz789",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### Success with no data (action endpoints like logout, cancel, revoke)
```json
HTTP 200 OK
{
  "success": true,
  "data": null,
  "error": null,
  "meta": {
    "request_id": "req_xyz789",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### Success with list + pagination
```json
HTTP 200 OK
{
  "success": true,
  "data": {
    "items": [ ... ],
  },
  "error": null,
  "meta": {
    "request_id": "req_xyz789",
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

### Error response
```json
HTTP 422 Unprocessable Entity
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "The request contains invalid fields",
    "field_violations": [
      { "field": "email", "message": "Email already registered" },
      { "field": "phone", "message": "Invalid phone number format" }
    ],
    "error_id": "err_123abc",
    "retryable": false
  },
  "meta": {
    "request_id": "req_xyz789",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

---

## Prohibited Patterns

### Problem 1: Error field embedded inside success response schemas
```yaml
# WRONG
RegistrationResponse:
  properties:
    user_id: string
    message: string
    error:                    # ← ERROR: error inside success response
      $ref: '#/components/schemas/Error'
```

```yaml
# CORRECT
RegistrationResponse:
  properties:
    user_id: string
    otp_id: string
    otp_expires_in_seconds: integer
    # No error field here. Errors come via HTTP 4xx with error envelope.
```

### Problem 2: No envelope — raw schema returned
```json
// WRONG
HTTP 200 OK
{
  "user_id": "usr_123",
  "message": "Registered",
  "otp_sent": true,
  "error": { "code": null }   // ← consumer must check this even on 200
}
```

```json
// CORRECT
HTTP 201 Created
{
  "success": true,
  "data": {
    "user_id": "usr_123",
    "otp_sent": true,
    "otp_id": "otp_abc",
    "otp_expires_in_seconds": 300
  },
  "error": null,
  "meta": { "request_id": "req_xyz" }
}
```

---

## Why This Matters for Clients

| Without Envelope | With Envelope |
|-----------------|---------------|
| Client checks HTTP status AND `.error` field | Client checks `success` boolean only |
| SDK generated models have `error?` on every success type | Clean models with no error contamination |
| Can't mock consistently | Single mock shape for all endpoints |
| DI containers need endpoint-specific null checks | One interceptor handles all responses |
| iOS/Android need custom deserializers per endpoint | One `ApiResponse<T>` generic decoder |

---

## Client Implementation Pattern (All Platforms)

### Swift (iOS)
```swift
struct ApiResponse<T: Decodable>: Decodable {
    let success: Bool
    let data: T?
    let error: ApiError?
    let meta: ResponseMeta?
}

// Usage - same pattern for EVERY endpoint
let response: ApiResponse<PolicyData> = try await apiClient.createPolicy(request)
if response.success, let policy = response.data {
    // use policy
} else if let error = response.error {
    // handle error
}
```

### Kotlin (Android)
```kotlin
data class ApiResponse<T>(
    val success: Boolean,
    val data: T?,
    val error: ApiError?,
    val meta: ResponseMeta?
)

// Retrofit + generic deserializer handles all endpoints identically
```

### TypeScript (Frontend)
```typescript
interface ApiResponse<T> {
  success: boolean;
  data: T | null;
  error: ApiError | null;
  meta: ResponseMeta | null;
}

// Axios interceptor handles once, not per endpoint
axios.interceptors.response.use(response => {
  if (!response.data.success) throw new ApiError(response.data.error);
  return response.data.data;
});
```
