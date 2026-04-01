# InsureTech Postman Collection Guide

## Overview

The InsureTech Postman collection provides comprehensive API testing and documentation for all 34+ microservices. It features:

- **1000+ auto-generated endpoints** from OpenAPI specification
- **Automatic environment variable propagation** — tokens captured from responses automatically populate in subsequent requests
- **Dual authentication flows** — B2C JWT (mobile) and B2B Server-Side Session (web/system)
- **Intelligent test scripts** — Every endpoint validates InsureTech API Rules (01, 02, 03)
- **Complete OTP workflows** — Mobile OTP, Email OTP, Password flows all testable end-to-end
- **Resource ID auto-capture** — Policy, Claim, Order IDs automatically captured from CRUD responses
- **Newman CLI ready** — Run full test suites headless in CI/CD pipelines

## File Structure

```
api/postman/
├── InsureTech.postman_collection.json          # Main collection (all endpoints)
├── auth_smoke.postman_collection.json           # Auth flow smoke tests
├── b2c_authz_enforcement.postman_collection.json # Authorization tests
├── InsureTech_local.postman_environment.json       # Local environment (localhost:8080)
├── InsureTech_staging.postman_environment.json     # Staging environment
├── InsureTech_production.postman_environment.json  # Production environment
├── InsureTech_mock.postman_environment.json        # Mock server environment
└── InsureTech_newman_test.postman_environment.json # Newman CLI test defaults
```

## Quick Start

### 1. Import into Postman

#### Option A: Desktop Postman
```bash
1. Open Postman
2. Click "Import" (top left)
3. Select api/postman/InsureTech.postman_collection.json
4. Click "Import"
5. Environment → Import → api/postman/InsureTech_local.postman_environment.json
6. Click environment dropdown → Select "InsureTech — Local"
```

#### Option B: Postman Web
```
1. Sign in to https://go.postman.co
2. Collections → Import → Upload Files
3. Select InsureTech.postman_collection.json
4. Environments → Import → InsureTech_local.postman_environment.json
```

### 2. Configure Environment Variables

In the selected environment, set:

```
base_url                = http://localhost:8080      # Your API server
user_mobile_number      = +1-234-567-8900           # Test phone number
user_email              = test@example.com           # Test email
login_device_type       = ANDROID                    # or IOS, API, WEB
login_password          = YourPassword123!           # For password-based login
```

### 3. Run Your First Request

**B2C Mobile OTP Flow (Recommended for testing):**

1. Open Collection → AuthService → POST `/v1/auth/otp:send`
2. Click **Send**
3. Response should show `"success": true` with `"data": {"otp_id": "..."}`
4. The `otp_id` is automatically captured → check environment
5. Set `mobile_otp_code` in environment (check your delivery channel for the code)
6. Run POST `/v1/auth/otp:verify` → Should return `"verified": true`
7. Run POST `/v1/auth/login` → Should return JWT with `access_token` and `refresh_token`
8. Both tokens are automatically captured!
9. Now you can run any protected endpoint — bearer token injected automatically

## Authentication Flows

### Flow 1: B2C Mobile OTP → JWT

**When to use:** Mobile apps (iOS, Android), API clients

**Sequence:**
```
1. POST /v1/auth/otp:send
   Input: mobile_number, device_id, device_type (ANDROID/IOS/API)
   Output: otp_id ← auto-captured
   
2. POST /v1/auth/otp:verify
   Input: mobile_number, otp_code (from SMS), otp_id, device_id, device_type
   Output: verified=true ← auto-captured
   
3. POST /v1/auth/login
   Input: mobile_number, password="", otp_id, device_id, device_type=ANDROID/IOS/API
   Output: access_token, refresh_token ← auto-captured, session_id
   
4. GET /v1/auth/session/current [+Bearer access_token]
   Output: Current session details
   
5. POST /v1/auth/sessions/refresh [+refresh_token]
   Output: New access_token ← auto-captured
   
6. POST /v1/auth/logout [+Bearer access_token]
   Output: session_revoked=true
```

**Environment Variables Used:**
- `user_mobile_number` (input)
- `mobile_otp_code` (manual input from SMS)
- `access_token` (auto-captured)
- `refresh_token` (auto-captured)
- `user_id` (auto-captured)
- `session_id` (auto-captured)

### Flow 2: Web Portal Password → Server-Side Session

**When to use:** Web portal, admin dashboard, internal tools

**Sequence:**
```
1. POST /v1/auth/login
   Input: mobile_number, password, device_type=WEB, device_id
   Output: session_token, csrf_token ← auto-captured, session_id, user_id
   Header: Set-Cookie: session_token=...
   
2. GET /v1/auth/session/current [+X-Session-Token: {{session_token}}, +X-CSRF-Token: {{csrf_token}}]
   Output: Current session details
   
3. POST /v1/auth/logout [+X-Session-Token, +X-CSRF-Token, body: {session_id}]
   Output: session_revoked=true
```

**Key Differences from Mobile:**
- `device_type=WEB` ONLY (never ANDROID/IOS/API for this flow)
- Returns `session_token` and `csrf_token` instead of JWT
- Password MUST NOT be empty
- No `access_token` or `refresh_token` returned
- Uses X-Session-Token and X-CSRF-Token headers

**Environment Variables Used:**
- `user_mobile_number` (input)
- `login_password` (input, must not be empty)
- `session_token` (auto-captured)
- `csrf_token` (auto-captured)
- `session_id` (auto-captured)
- `user_id` (auto-captured)

### Flow 3: Email OTP → Server-Side Session

**When to use:** Email-based login, account recovery

**Sequence:**
```
1. POST /v1/auth/email/otp:send
   Input: email, otp_type="email_login"
   Output: otp_id ← auto-captured
   
2. POST /v1/auth/email/login
   Input: email, otp_code (from email), otp_id, device_type, device_id
   Output: session_token, csrf_token ← auto-captured, session_id, user_id
   
3. POST /v1/auth/logout [+X-Session-Token, +X-CSRF-Token]
   Output: session_revoked=true
```

**Important:** Email OTP login NEVER returns `access_token` or `refresh_token` — always server-side session.

## Auto-Capturing Variables

The collection automatically extracts and stores variables from successful responses:

### Authentication Tokens
- `access_token` — JWT token for B2C mobile (captured from login response)
- `refresh_token` — Token for getting new access_token (captured from login response)
- `session_token` — Cryptographic token for server-side session (captured from login response)
- `csrf_token` — CSRF protection token (captured from login response, required in X-CSRF-Token header)
- `session_id` — UUID of current session (captured, used in logout request body)

### User Identity
- `user_id` — Current user ID (auto-captured from login)
- `user_mobile_number` — User phone number (auto-captured if returned)
- `user_email` — User email address (auto-captured if returned)
- `tenant_id` — User's tenant/organization ID (auto-captured if returned)

### OTP Flows
- `mobile_otp_id` — ID from mobile OTP send (auto-captured)
- `email_otp_id` — ID from email OTP send (auto-captured)
- `email_login_otp_id` — ID from email login OTP (auto-captured)

### Resource IDs (CRUD Operations)
All of these are auto-captured from POST/PATCH responses:
- `policy_id` ← captured from POST /v1/policies
- `claim_id` ← captured from POST /v1/claims
- `order_id` ← captured from POST /v1/orders
- `payment_id` ← captured from POST /v1/payments
- `product_id` ← captured from POST /v1/products
- `quote_id` ← captured from POST /v1/quotes
- `ticket_id` ← captured from POST /v1/support/tickets
- `partner_id` ← captured from POST /v1/partners
- `kyc_id` ← captured from POST /v1/kyc/verifications
- `invoice_id` ← captured from POST /v1/invoices
- `document_id` ← captured from POST /v1/documents

You can also reference `last_policy_id`, `last_claim_id`, etc. for more explicit chaining.

## Test Scripts

Every endpoint includes automated tests that verify:

### Rule 01: Response Envelope
```javascript
✓ Response has success, data, error fields
✓ If success=true, error must be null
✓ If success=false, error must be non-null with code and message
```

### Rule 02: HTTP Status Codes
```javascript
✓ Success responses: 200, 201, or 204
✓ Error responses: 400, 401, 403, 404, 422, 429, 500, 503
```

### Rule 03: Error Details
```javascript
✓ All errors have error.code (string)
✓ All errors have error.message
✓ 422 errors have error.field_violations array
```

### Additional Checks
```javascript
✓ Response time < 5000ms
✓ Content-Type: application/json
✓ Resource IDs auto-captured for subsequent requests
```

## Newman CLI Usage

### Installation
```bash
npm install -g newman
# or use npx to run without global install
```

### Run Collection Locally
```bash
newman run api/postman/InsureTech.postman_collection.json \
  -e api/postman/InsureTech_local.postman_environment.json \
  --reporters cli,json \
  --reporter-json-export results.json
```

### Run with Custom Environment Variables
```bash
newman run api/postman/InsureTech.postman_collection.json \
  -e api/postman/InsureTech_local.postman_environment.json \
  --env-var base_url=http://localhost:9090 \
  --env-var user_mobile_number="+1-555-0123" \
  --reporters cli,json
```

### Run Specific Folder (e.g., AuthService only)
```bash
newman run api/postman/InsureTech.postman_collection.json \
  -e api/postman/InsureTech_local.postman_environment.json \
  --folder "AuthService" \
  --reporters cli
```

### Run Auth Smoke Tests Only
```bash
newman run api/postman/auth_smoke.postman_collection.json \
  -e api/postman/InsureTech_newman_test.postman_environment.json \
  --reporters cli,htmlextra \
  --reporter-htmlextra-export auth_report.html
```

### Run with Delay Between Requests
```bash
newman run api/postman/InsureTech.postman_collection.json \
  -e api/postman/InsureTech_local.postman_environment.json \
  --delay-request 500 \
  --reporters cli,json
```

### Run with Data File (CSV/JSON)
```bash
# Create test_data.json with array of test cases
newman run api/postman/InsureTech.postman_collection.json \
  -e api/postman/InsureTech_local.postman_environment.json \
  -d test_data.json \
  --reporters cli
```

## CI/CD Integration

### GitHub Actions
```yaml
name: API Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Newman
        run: npm install -g newman
      
      - name: Run Postman Tests
        run: |
          newman run api/postman/InsureTech.postman_collection.json \
            -e api/postman/InsureTech_local.postman_environment.json \
            --env-var base_url=http://localhost:8080 \
            --reporters cli,json \
            --reporter-json-export results.json
      
      - name: Upload Results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: postman-results
          path: results.json
```

### GitLab CI
```yaml
test:postman:
  image: node:18
  script:
    - npm install -g newman
    - newman run api/postman/InsureTech.postman_collection.json \
        -e api/postman/InsureTech_local.postman_environment.json \
        --env-var base_url=http://api:8080 \
        --reporters cli,json \
        --reporter-json-export results.json
  artifacts:
    paths:
      - results.json
    when: always
```

## Troubleshooting

### "Bearer token required" but token is set
**Check:**
1. Environment has `access_token` set (should appear in red/secret type)
2. Request includes pre-request script (should see "Auth injection" in console)
3. Run a login endpoint first to capture the token
4. Verify token not expired (refresh if needed)

### "Invalid CSRF token" for Web session
**Check:**
1. `csrf_token` is set in environment
2. Request includes `X-CSRF-Token` header (should see in Headers tab)
3. CSRF token was captured from login (check collection's POST login test script)
4. CSRF token not expired (re-login if expired)

### "OTP code mismatch"
**Check:**
1. OTP code copied correctly (no leading/trailing spaces)
2. Code not expired (typical TTL 5-10 minutes)
3. Code matches the OTP ID from send request
4. Device ID matches between send and verify

### "Mobile number required" but it's set
**Check:**
1. Environment variable is: `user_mobile_number` (not `mobile_number`)
2. Format matches backend expectations (usually: `+1-234-567-8900` or `+11234567890`)
3. Variable is enabled (blue toggle next to variable name)

### Newman fails but Postman works
**Check:**
1. `base_url` is set in environment (defaults to `{{base_url}}`)
2. Server is reachable from CI/CD environment
3. POSTMAN_API_KEY not required (it's only for uploading collections)
4. Check Newman verbose output: `newman run ... -v`

## Manual Variable Entry

If auto-capture doesn't work, you can manually set variables:

1. Click **Environment** (top right in Postman)
2. Select **InsureTech — Local**
3. Click **Edit**
4. Find the variable (e.g., `access_token`)
5. Click in the **Current Value** column
6. Paste the token
7. Click **Save**

## Advanced: Custom Environment Variables

Create additional environments by duplicating `InsureTech_local.postman_environment.json`:

```json
{
  "id": "unique-uuid-here",
  "name": "InsureTech — Custom",
  "values": [
    {"key": "base_url", "value": "https://api.custom.com", "enabled": true},
    {"key": "api_key", "value": "sk_...", "enabled": true},
    // ... other variables ...
  ]
}
```

Then import this file into Postman and select it.

## Performance Optimization

### Reduce Test Execution Time
```bash
# Skip tests, just run requests
newman run api/postman/InsureTech.postman_collection.json \
  -e api/postman/InsureTech_local.postman_environment.json \
  --reporters cli \
  --timeout-request 5000
```

### Parallel Execution
Newman runs requests sequentially by default. For parallel:
```bash
# Run multiple times with different data
for i in {1..5}; do
  newman run api/postman/InsureTech.postman_collection.json \
    -e api/postman/InsureTech_local.postman_environment.json \
    --env-var user_mobile_number="+1-555-000$i" &
done
wait
```

### Run Subset of Tests
```bash
# Only run AuthService endpoints
newman run api/postman/InsureTech.postman_collection.json \
  -e api/postman/InsureTech_local.postman_environment.json \
  --folder "AuthService" \
  --reporters cli
```

## Postman API Integration

If `POSTMAN_API_KEY` is set in `.env`, the collection is automatically uploaded:

```bash
# Set in .env
POSTMAN_API_KEY=your-api-key-here
POSTMAN_WORKSPACE_ID=your-workspace-id    # (optional)
POSTMAN_COLLECTION_ID=existing-id         # (optional, for updates)

# Then run
python api/generator/sync_postman.py --upload
```

Access your collection at: `https://go.postman.co/collections`

## Further Reading

- [Postman Documentation](https://learning.postman.com/)
- [Newman CLI Reference](https://github.com/postmanlabs/newman)
- [InsureTech API Rules](./documentation/API_RULES.md)
- [Authentication Architecture](./documentation/AUTHENTICATION.md)

