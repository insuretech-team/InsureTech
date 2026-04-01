# Analysis: authServiceGetCurrentSession Function

## Overview
The `authServiceGetCurrentSession` function retrieves the current user's active session from the InsureTech API. Here's exactly how it works:

---

## 1. Function Definition (sdk.gen.ts, line 702-706)

```typescript
/**
 * Get current user's active session
 */
export const authServiceGetCurrentSession = <ThrowOnError extends boolean = false>(
    options?: Options<AuthServiceGetCurrentSessionData, ThrowOnError>
) => (options?.client ?? client).get<
    AuthServiceGetCurrentSessionResponses, 
    AuthServiceGetCurrentSessionErrors, 
    ThrowOnError
>({
    security: [{ scheme: 'bearer', type: 'http' }],
    url: '/v1/auth/session/current',
    ...options
});
```

### Key Details:
- **HTTP Method**: GET
- **URL Endpoint**: `/v1/auth/session/current`
- **Authentication**: Bearer token (HTTP authorization header)
- **Generic Type Parameters**:
  - `AuthServiceGetCurrentSessionResponses` - Success response type
  - `AuthServiceGetCurrentSessionErrors` - Error response types
  - `ThrowOnError` - Whether to throw on errors (default: false)

---

## 2. Request Type (types.gen.ts, line 27088-27103)

```typescript
export type AuthServiceGetCurrentSessionData = {
    body?: never;
    path?: never;
    query?: {
        /**
         * Page number (1-based)
         */
        page?: number;
        /**
         * Number of items per page
         */
        page_size?: number;
    };
    url: '/v1/auth/session/current';
};
```

### Request Characteristics:
- **No request body** required
- **No path parameters**
- **Optional query parameters**: `page` and `page_size` for pagination
- Simple GET request with bearer token

---

## 3. Response Types (types.gen.ts)

### 3a. Success Response (lines 27140-27145)

```typescript
export type AuthServiceGetCurrentSessionResponses = {
    /**
     * Successful response
     */
    200: CurrentSessionRetrievalResponse;
};
```

**Status Code 200 (Success)** returns: `CurrentSessionRetrievalResponse`

### 3b. Error Responses (lines 27105-27136)

```typescript
export type AuthServiceGetCurrentSessionErrors = {
    /**
     * Bad request - Malformed request body or invalid parameters
     */
    400: ApiResponse & {
        data?: CurrentSessionRetrievalResponse;
    };
    /**
     * Unauthorized - Valid authentication token required
     */
    401: ApiResponse & {
        data?: CurrentSessionRetrievalResponse;
    };
    /**
     * Forbidden - Insufficient permissions for this operation
     */
    403: ApiResponse & {
        data?: CurrentSessionRetrievalResponse;
    };
    /**
     * Too Many Requests - Rate limit exceeded. Retry after the indicated delay
     */
    429: ApiResponse & {
        data?: CurrentSessionRetrievalResponse;
    };
    /**
     * Internal server error - Unexpected server-side error
     */
    500: ApiResponse & {
        data?: CurrentSessionRetrievalResponse;
    };
};
```

**Error Status Codes**:
- **400** - Bad request
- **401** - Unauthorized (missing/invalid token)
- **403** - Forbidden (insufficient permissions)
- **429** - Rate limit exceeded
- **500** - Internal server error

All errors wrapped in `ApiResponse` envelope with optional data field.

---

## 4. Response Payload Type: CurrentSessionRetrievalResponse

### Definition (types.gen.ts, lines 3485-3491)

```typescript
export type CurrentSessionRetrievalResponse = {
    session?: Session;
    /**
     * user_type lets callers determine the role without an extra user lookup.
     * Same format as ValidateTokenResponse.user_type: e.g. "SYSTEM_USER"
     */
    user_type?: string;
};
```

**Returns**:
- `session`: The `Session` object (optional)
- `user_type`: User's role/type as string (e.g., "SYSTEM_USER"), optional

### Session Object Type (types.gen.ts, lines 17604-17676)

```typescript
export type Session = {
    /**
     * @inject_tag: gorm:"primaryKey;column:session_id;not null"
     */
    session_id?: string;
    /**
     * @inject_tag: gorm:"column:user_id;not null"
     */
    user_id?: string;
    /**
     * Session type determines authentication method (SERVER_SIDE for web, JWT for mobile)
     * @inject_tag: gorm:"column:session_type;not null;serializer:proto_enum"
     */
    session_type?: SessionType;
    /**
     * JWT ID for access token (only for JWT sessions)
     * @inject_tag: gorm:"column:access_token_jti"
     */
    access_token_jti?: string;
    /**
     * JWT ID for refresh token (only for JWT sessions)
     * @inject_tag: gorm:"column:refresh_token_jti"
     */
    refresh_token_jti?: string;
    /**
     * @inject_tag: gorm:"column:access_token_expires_at;serializer:proto_timestamp"
     */
    access_token_expires_at?: string;
    /**
     * @inject_tag: gorm:"column:refresh_token_expires_at;serializer:proto_timestamp"
     */
    refresh_token_expires_at?: string;
    /**
     * Session expiry: 12 hours for SERVER_SIDE, 7 days for JWT
     * @inject_tag: gorm:"column:expires_at;not null;serializer:proto_timestamp"
     */
    expires_at?: string;
    /**
     * @inject_tag: gorm:"column:ip_address"
     */
    ip_address?: string;
    /**
     * @inject_tag: gorm:"column:user_agent"
     */
    user_agent?: string;
    /**
     * @inject_tag: gorm:"column:device_id"
     */
    device_id?: string;
    /**
     * @inject_tag: gorm:"column:device_name"
     */
    device_name?: string;
    /**
     * @inject_tag: gorm:"column:device_type;not null;serializer:proto_enum"
     */
    device_type?: DeviceType;
    /**
     * @inject_tag: gorm:"column:created_at;not null;serializer:proto_timestamp"
     */
    created_at?: string;
    /**
     * @inject_tag: gorm:"column:last_activity_at;not null;serializer:proto_timestamp"
     */
    last_activity_at?: string;
    /**
     * @inject_tag: gorm:"column:is_active;not null"
     */
    is_active?: boolean;
};
```

**Session Object Fields**:
- `session_id`: Unique session identifier
- `user_id`: Associated user ID
- `session_type`: SERVER_SIDE or JWT
- `access_token_jti`, `refresh_token_jti`: JWT token IDs
- `access_token_expires_at`, `refresh_token_expires_at`: Token expiration timestamps
- `expires_at`: Overall session expiry (12 hours for SERVER_SIDE, 7 days for JWT)
- `ip_address`: IP address of session origin
- `user_agent`: Browser/client user agent string
- `device_id`, `device_name`: Device identifiers
- `device_type`: DeviceType enum
- `created_at`: Session creation timestamp
- `last_activity_at`: Last activity timestamp
- `is_active`: Whether session is currently active

---

## 5. ApiResponse Envelope Structure (types.gen.ts, lines 1088-1103)

```typescript
export type ApiResponse = {
    /**
     * true when the operation succeeded, false on any error.
     */
    success: boolean;
    /**
     * Typed response payload on success. null on failure or no-content actions
     * (logout, delete, etc.). Concrete type is specified per-endpoint via allOf composition.
     */
    data?: unknown;
    /**
     * Error details on failure. Always null on success.
     */
    error?: Error;
    /**
     * Response metadata: request tracing, pagination, timestamps.
     */
    meta?: ResponseMetadata;
};
```

**Envelope Fields**:
- `success`: boolean indicating operation success/failure
- `data`: Typed payload (null on failure)
- `error`: Error details (null on success)
- `meta`: Metadata like request tracing, pagination

---

## 6. Response Interceptor & Unwrapping (client-wrapper.ts, lines 42-74)

```typescript
// ── Unwrap ApiResponse envelope ─────────────────────────────────────────
// The gateway wraps every response as { success, data, error, meta }.
// hey-api puts the parsed JSON into result.data, so without this
// interceptor consumers would need result.data.data to reach the payload.
// By replacing the Response body with just the inner "data" field we make
// result.data === T directly — no double-wrap.
c.interceptors.response.use(async (response) => {
    const ct = response.headers.get('content-type') ?? '';
    if (!ct.includes('application/json')) return response;
    // Clone so we can read the body without consuming the original.
    const text = await response.clone().text();
    if (!text) return response;
    try {
        const envelope = JSON.parse(text);
        // Only unwrap if it looks like our standard ApiResponse envelope.
        if (
            typeof envelope === 'object' &&
            envelope !== null &&
            'success' in envelope &&
            'data' in envelope
        ) {
            // Success: unwrap envelope.data so result.data === T
            // Error: unwrap envelope.error so result.error has gateway error details
            const inner = envelope.success ? envelope.data : envelope.error;
            return new Response(JSON.stringify(inner ?? {}), {
                status: response.status,
                statusText: response.statusText,
                headers: response.headers,
            });
        }
    } catch { /* not JSON — pass through */ }
    return response;
});
```

### How the Interceptor Works:

1. **Intercepts Response**: Catches all HTTP responses
2. **Checks Content-Type**: Only processes JSON responses
3. **Parses Envelope**: Checks if response is an ApiResponse envelope
4. **Unwraps Payload**: 
   - On success (200): Extracts `envelope.data` → becomes `result.data`
   - On error: Extracts `envelope.error` → becomes `result.error`
5. **Creates New Response**: Returns modified response with unwrapped data
6. **Fallthrough**: Non-JSON or non-envelope responses pass through unchanged

**Key Benefit**: Consumers get `result.data` directly (e.g., `CurrentSessionRetrievalResponse`) instead of `result.data.data`. Single-level data access without double-wrapping.

---

## 7. Complete Call Flow

```
┌─────────────────────────────────────────────────────────────┐
│ Client Code                                                 │
│ authServiceGetCurrentSession({ client })                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ SDK Function (sdk.gen.ts)                                   │
│ - Calls: client.get<Responses, Errors, ThrowOnError>()      │
│ - URL: /v1/auth/session/current                             │
│ - Auth: Bearer token                                        │
│ - Method: GET                                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ HTTP Request                                                │
│ GET /v1/auth/session/current                                │
│ Authorization: Bearer <token>                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ API Gateway Response (Raw)                                  │
│ {                                                           │
│   "success": true,                                          │
│   "data": {                                                 │
│     "session": { ... Session object ... },                  │
│     "user_type": "SYSTEM_USER"                              │
│   },                                                        │
│   "error": null,                                            │
│   "meta": { ... }                                           │
│ }                                                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Response Interceptor (client-wrapper.ts)                    │
│ - Detects ApiResponse envelope (success + data fields)      │
│ - Extracts envelope.data                                    │
│ - Creates new Response with unwrapped payload               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ hey-api Client Processing                                   │
│ - Parses JSON → result.data                                 │
│ - Types as CurrentSessionRetrievalResponse                  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Returned to Consumer                                        │
│ result.data = CurrentSessionRetrievalResponse {             │
│   session?: Session,                                        │
│   user_type?: string                                        │
│ }                                                           │
└─────────────────────────────────────────────────────────────┘
```

---

## 8. Usage Example

```typescript
import { createInsureTechClient, authServiceGetCurrentSession } from '@lifeplus/insuretech-sdk';

const client = createInsureTechClient({
    apiKey: 'your-api-key',
    baseUrl: 'https://api.insuretech.com'
});

// Call the function
const response = await authServiceGetCurrentSession({ client });

if (response.data) {
    // Access the session directly (thanks to interceptor unwrapping)
    console.log('Session ID:', response.data.session?.session_id);
    console.log('User Type:', response.data.user_type);
    console.log('Is Active:', response.data.session?.is_active);
    console.log('Expires At:', response.data.session?.expires_at);
} else if (response.error) {
    console.error('Error:', response.error);
}
```

---

## Summary

**What it does**:
- Makes a simple GET request to `/v1/auth/session/current` with bearer token authentication
- Returns the current user's active session with session metadata and user type

**Key Processing**:
1. Bearer token authentication required
2. Gateway wraps response in `ApiResponse` envelope
3. Response interceptor automatically unwraps envelope
4. Consumer receives `CurrentSessionRetrievalResponse` directly
5. Session includes all metadata: IDs, types, tokens, expiration, device info, timestamps, activity status

**Response Unwrapping**: Eliminates double-wrapping by automatically extracting inner payload from the `{ success, data, error, meta }` envelope, making data access cleaner (single-level instead of double-level access).
