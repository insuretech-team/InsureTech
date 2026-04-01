# B2B Portal Session Authentication Investigation

## Status: RESOLVED (March 16, 2026)

---

## Executive Summary

The B2B portal on the remote server (b2b.labaidinsuretech.com) returned **401 Unauthorized** for every authenticated endpoint (session, profile, employees, departments, purchase-orders, organisations) immediately after successful login. Local development worked perfectly.

**Root Cause:** The `SessionLimiter` in the authn service used **expiry time** as the Redis sorted set score. JWT sessions (from mobile/API testing) had 7-day expiry while SERVER_SIDE web sessions had 12-hour expiry. When the sorted set was full (5 slots), `ZPOPMIN` evicted the member with the **lowest score** — which was always the **brand-new web session** (12h expiry) rather than the older JWT sessions (7-day expiry). The new session was revoked within the same login request.

**Fix:** Changed the sorted set score from expiry time to **creation time**, so `ZPOPMIN` always evicts the oldest-created session regardless of session type or TTL.

---

## 1. Investigation Timeline

### Symptoms Observed
- Login on remote b2b.labaidinsuretech.com succeeded (200 OK, cookies set in browser)
- Immediately after login, all authenticated API calls returned 401
- Gateway logs showed 401 responses for requests from b2b_portal (user_agent: "node")
- Authn service `ValidateToken` gRPC calls returned `code=OK` but gateway still gave 401
- Local development (same Neon DB, local Redis) worked perfectly

### Eliminated Causes
- **Cookie forwarding:** BFF correctly reads `session_token` cookie and forwards via SDK → verified
- **SDK interceptor:** ApiResponse envelope unwrapping works correctly → verified in deployed chunk
- **Docker networking:** b2b_portal container can reach gateway:8080 → healthz/readyz confirmed
- **Nginx proxy:** Correct proxy_pass, headers forwarded properly → config verified
- **Token extraction in login route:** `result.data.session_token` correctly read via `UseProtoNames: true` → verified
- **Redis connectivity:** Redis responds to PING → confirmed
- **Idle timeout:** `IDLE_TIMEOUT_SECONDS` defaults to 0 (disabled) → not the issue

### Key Discovery: ValidateToken Paradox

At the same timestamp, authn logs showed `ValidateToken` returning `grpc.code=OK` while gateway returned 401. The auth middleware code:

```go
// auth_middleware.go — line 62-66
resp, err := client.ValidateToken(ctx, &authnservicev1.ValidateTokenRequest{
    AccessToken: jwt,
    SessionId:   sessionToken,
})
if err != nil { /* logs "ValidateToken failed" */ }
if resp == nil || !resp.Valid {
    // Silent 401 — NO log line for valid=false case
    respond.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "Unauthorized")
    return
}
```

The gRPC call succeeded (`err == nil`, `code=OK`) but the response had `Valid: false`. This meant the session existed but was **already revoked** by the time the validation check ran.

---

## 2. Root Cause: SessionLimiter Evicts New Sessions

### The Bug

**File:** `backend/inscore/microservices/authn/internal/service/session_limiter.go`

The `SessionLimiter` uses a Redis sorted set (`sessions:active:<userID>`) to enforce a per-user limit of 5 concurrent sessions. The **score** was the session's **expiry timestamp**.

```
Redis sorted set (BEFORE fix):
Key: sessions:active:4ba2c039-f208-4937-8a05-008558f40af1

Member (sessionID)                       | Score (expiry)    | Type
-----------------------------------------|-------------------|-------------
e61a923b-fa7b-4f4f-a520-22dca11c1172     | 1774186221 (Mar 22) | JWT (7-day)
73bd977c-c4ea-4e7a-bc27-32c47da28bcc     | 1774190190 (Mar 22) | JWT (7-day)
64253f40-f84e-4380-9293-34139f862f69     | 1774190211 (Mar 22) | JWT (7-day)
859810d4-6890-429f-ad4c-9c6016e8cc60     | 1774190609 (Mar 22) | JWT (7-day)
9a0242de-c698-4c81-8c56-03613ecd28f6     | 1774190637 (Mar 22) | JWT (7-day)
```

When the user logs into the B2B portal (SERVER_SIDE session, 12h expiry):

1. **ZADD** new session with score = `now + 12h` ≈ 1774020000 (Mar 16) → 6 entries
2. **ZREMRANGEBYSCORE** removes nothing — all JWT expiries are in the future
3. **ZCARD** = 6, exceeds limit of 5
4. **ZPOPMIN(1)** evicts the entry with the **lowest score** → the brand-new web session (score ~Mar 16) because 12h < 7 days
5. `RevokeSession()` marks it `is_active = false` in DB
6. Login handler returns the session token to the browser
7. Browser sends the token on next request
8. `ValidateToken` finds the session is revoked → `Valid: false` → **401**

### Database Evidence

```sql
-- All recent SERVER_SIDE sessions: is_active=false, last_activity_at=created_at (never used)
SELECT session_id, is_active, session_type, last_activity_at, created_at
FROM authn_schema.sessions
WHERE user_id = '4ba2c039-f208-4937-8a05-008558f40af1'
ORDER BY created_at DESC LIMIT 10;

-- Result: 9 SERVER_SIDE sessions with is_active=false, last_activity_at=created_at
--         1 old SERVER_SIDE with is_active=true (created before JWT sessions filled Redis)

-- The 5 JWT sessions occupying Redis slots:
-- All is_active=true, session_type=JWT, expires_at=Mar 22 (7-day TTL)
```

### Why Local Worked But Remote Failed

|                    | Local Dev           | Remote (DigitalOcean)     |
|--------------------|---------------------|---------------------------|
| **Database**       | Same Neon DB        | Same Neon DB              |
| **Redis**          | Local Docker Redis  | Remote Docker Redis       |
| **Redis state**    | Clean — only web sessions | 5 JWT sessions from mobile/API testing |
| **Session types**  | Only SERVER_SIDE (12h) | Mixed: JWT (7-day) + SERVER_SIDE (12h) |
| **Limiter effect** | ZPOPMIN evicts truly oldest session | ZPOPMIN evicts newest web session (lowest expiry score) |
| **Result**         | Works ✅            | 401 immediately after login ❌ |

---

## 3. Fix Applied

### Code Change

**File:** `backend/inscore/microservices/authn/internal/service/session_limiter.go`

Changed the sorted set score from **expiry time** to **creation time**:

```go
// BEFORE (buggy):
// Score = expiry time → short-lived sessions get evicted over long-lived ones
if err := sl.rdb.ZAdd(ctx, k, redis.Z{
    Score:  float64(expiry.Unix()),   // ← 12h session has LOWER score than 7-day
    Member: sessionID,
}).Err(); err != nil { ... }

// AFTER (fixed):
// Score = creation time → oldest-created session always evicted first
now := time.Now().UTC()
if err := sl.rdb.ZAdd(ctx, k, redis.Z{
    Score:  float64(now.Unix()),      // ← newest session always has HIGHEST score
    Member: sessionID,
}).Err(); err != nil { ... }
```

Also removed the stale `ZREMRANGEBYSCORE` cleanup step (no longer needed since scores are creation times, not expiry times). Expired sessions are naturally cleaned when they're evicted and the DB reports `is_active=false`.

### Data Cleanup

```bash
# 1. Cleared the stale Redis sorted set
docker exec insuretech-redis redis-cli DEL 'sessions:active:4ba2c039-...'

# 2. Revoked orphaned JWT sessions in DB
UPDATE authn_schema.sessions
SET is_active = false
WHERE user_id = '4ba2c039-...' AND session_type = 'JWT' AND is_active = true;
-- Rows affected: 5
```

### Deployment

Redeployed authn service via `quickerdeploy.sh --services=authn`.

---

## 4. Session Auth Architecture Reference

### Authentication Flow

```
Browser → POST /api/auth/login (Next.js BFF)
  → SDK calls POST /v1/auth/login (gateway)
    → gRPC Login (authn service)
      → Creates SERVER_SIDE session in DB
      → Sets session:idle:<id> in Redis (if idle timeout configured)
      → SessionLimiter.TrackSession() in Redis sorted set
      → Returns session_token in proto JSON body
    ← Gateway sets Set-Cookie: session_token=<token>
    ← Gateway wraps in ApiResponse envelope
  ← SDK interceptor unwraps envelope (result.data = LoginResponse)
  ← Login route reads result.data.session_token
  ← Sets HttpOnly cookie: session_token=<token>
  ← Sets metadata cookies: portal_role, portal_user_id, portal_biz_id
← Browser stores all cookies

Browser → GET /api/auth/session (Next.js BFF)
  → Checks session_token cookie exists
  → SDK calls GET /v1/auth/session/current (gateway, cookie forwarded)
    → AuthMiddleware extracts session_token cookie
    → gRPC ValidateToken(SessionId: token) → authn service
      → SHA256 lookup → bcrypt verify → expiry check → idle check
      → Returns Valid: true + user metadata
    ← Middleware sets X-User-ID, X-Portal, etc. headers
    ← Handler calls GetCurrentSession
  ← SDK interceptor unwraps ApiResponse envelope
  ← Session route converts to PortalSession
  ← Re-mints metadata cookies
← Browser receives session data + refreshed cookies
```

### Key Components

| Component | File | Role |
|-----------|------|------|
| Login route | `b2b_portal/app/api/auth/login/route.ts` | BFF login, sets cookies |
| Session route | `b2b_portal/app/api/auth/session/route.ts` | Validates session, re-mints metadata cookies |
| SDK client factory | `b2b_portal/src/lib/sdk/b2b-sdk-client.ts` | Forwards cookies + portal headers to gateway |
| SDK interceptor | `sdks/insuretech-typescript-sdk/src/client-wrapper.ts` | Unwraps ApiResponse envelope |
| Session headers | `b2b_portal/src/lib/sdk/session-headers.ts` | Resolves x-portal/x-user-id/x-business-id |
| Edge middleware | `b2b_portal/middleware.ts` | Role-based page routing |
| Auth middleware | `backend/.../gateway/internal/routes/auth_middleware.go` | Validates session via gRPC |
| Token service | `backend/.../authn/internal/service/token_service.go` | Session creation + validation |
| Session limiter | `backend/.../authn/internal/service/session_limiter.go` | Per-user concurrent session limit (THE BUG) |

### Cookie Inventory

| Cookie | HttpOnly | Purpose | Set By |
|--------|----------|---------|--------|
| `session_token` | Yes | Backend session auth | Login route |
| `csrf_token` | Yes | CSRF protection | Login route |
| `portal_role` | No | Edge middleware routing | Login + session routes |
| `portal_user_id` | No | x-user-id header injection | Login + session routes |
| `portal_biz_id` | No | x-business-id header injection | Login + session routes |
| `portal_mobile` | No | Contact info display | Login + session routes |
| `portal_email` | No | Contact info display | Login + session routes |

---

## 5. Lessons Learned

1. **Never use expiry time as a sorted set score for eviction.** Different session types have different TTLs. Using expiry as score causes short-lived sessions to be evicted before long-lived ones — the exact opposite of "evict oldest."

2. **Shared Redis state between session types (JWT vs SERVER_SIDE) needs careful handling.** The session limiter counted all session types in the same bucket but didn't account for TTL differences.

3. **Silent 401s are hard to debug.** The auth middleware returned 401 for `resp.Valid == false` without logging. Added this to the investigation notes — consider adding a log line for this case.

4. **Local vs remote divergence from Redis state.** Same DB but different Redis instances means session limiter state diverges. Always test with realistic Redis state (multiple session types, near-limit counts).
