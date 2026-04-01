# Postman Variable Propagation Gap Analysis

## Executive Summary
All three collections have **well-managed variable chains** with proper setup/teardown. The primary gaps are:
1. **Email OTP variables** not set in all environments
2. **Token validation logic** assumes JWT presence (edge case)
3. **Minor inconsistencies** between environment files (naming: `session_token` vs implicit cookie handling)

---

## 1. ENVIRONMENT VARIABLES INVENTORY

### InsureTech_newman_test.postman_environment.json
**24 variables** (test/CI focused):
- **Auth tokens:** `access_token`, `refresh_token`, `session_id`, `session_type`
- **User data:** `user_mobile_number`, `user_id`, `device_id`, `device_name`
- **OTP (Mobile only):** `mobile_otp_id`, `mobile_otp_code`, `mobile_otp_verified`, `last_b2c_otp_login`
- **Login:** `login_password`, `login_device_type`, `login_expect_session_type`
- **AuthZ test:** `authz_resource`, `authz_action`
- ❌ **MISSING:** Email OTP variables, `session_token`, `csrf_token`

### InsureTech_local.postman_environment.json
**46 variables** (local dev focused, comprehensive):
- All from newman_test, PLUS:
- **Email OTP:** `email_otp_type`, `email_otp_id`, `email_otp_code`, `email_login_otp_id`, `email_login_otp_code`, `email_verification_otp_id`, `email_verification_otp_code`, `email_verified`
- **Session (WEB):** `session_token`, `csrf_token`
- **Business entity IDs:** `policy_id`, `claim_id`, `payment_id`, `order_id`, `product_id`, `quote_id`, `ticket_id`, `partner_id`, `kyc_id`, `invoice_id`, `document_id`
- **User email:** `user_email`
- **API auth:** `api_key`

### .env.example (Backend defaults)
**Does NOT provide Postman runtime defaults** — it's for backend services:
- Database creds (PGHOST, PGDATABASE, etc.)
- JWT config (JWT_PRIVATE_KEY_PATH, JWT_AUDIENCE, JWT_ACCESS_TOKEN_DURATION=15m)
- OTP config (OTP_VALIDITY_MINUTES=5, OTP_LENGTH=6)
- Service URLs (AUTHN_HOST, KAFKA_BROKERS, REDIS_URL, GOTENBERG_URL)
- **Relevant for manual setup:** `user_mobile_number`, `mobile_otp_code` must come from test execution, not .env

---

## 2. b2c_authn_authz_suite.postman_collection.json

### Variables READ ({{var}} in requests)
| Variable | Used In | Count |
|----------|---------|-------|
| `base_url` | All 15 requests | 15 |
| `user_mobile_number` | OTP Send (AuthN-01) | 1 |
| `mobile_otp_type` | OTP Send (AuthN-01) | 1 |
| `mobile_otp_channel` | OTP Send (AuthN-01) | 1 |
| `mobile_otp_id` | OTP Verify (AuthN-02), Login (AuthN-03) | 2 |
| `mobile_otp_code` | OTP Verify (AuthN-02) | 1 |
| `device_id` | OTP Verify (AuthN-02), Login (AuthN-03), AuthZ pre-request | 3 |
| `access_token` | All AuthZ requests (AuthZ-01..10), AutoLogin pre-request | 12 |
| `user_id` | AuthZ-02 check body | 1 |
| `session_id` | Logout (AuthN-06), AuthZ pre-request reads from env | 1 |
| `login_password` | AuthZ pre-request auto-login | 1 |

### Variables SET (pm.environment.set)
| Variable | Set In | Source |
|----------|--------|--------|
| `mobile_otp_id` | AuthN-01 test script (line 71) | `j.data.otp_id` |
| `mobile_otp_id` | AuthN-01 test (line 51) — **CLEARED on 429/5xx** | Empty string |
| `login_password` | AuthN-02 test script (line 151) | `j.data.device_credential` |
| `access_token` | AuthN-03 test (line 232) | `j.data.access_token` |
| `refresh_token` | AuthN-03 test (line 233) | `j.data.refresh_token` |
| `session_id` | AuthN-03 test (line 234) | `j.data.session_id` |
| `user_id` | AuthN-03 test (line 235) | `j.data.user.user_id` |
| `access_token` | AuthN-05 test (line 340) | `j.data.access_token` (refresh) |
| `refresh_token` | AuthN-05 test (line 341) | `j.data.refresh_token` (refresh) |
| `session_id` | AuthN-05 test (line 342) | `j.data.session_id` (refresh) |
| `access_token` | AuthZ pre-request auto-login (line 577) | `j.data.access_token` |
| `refresh_token` | AuthZ pre-request auto-login (line 578) | `j.data.refresh_token` |
| `session_id` | AuthZ pre-request auto-login (line 579) | `j.data.session_id` |
| `user_id` | AuthZ pre-request auto-login (line 580) | `j.data.user.user_id` |

### Variables MISSING (used but never set in scripts)
**None detected** — all variables used in requests are either:
- Pre-populated in environment (e.g., `user_mobile_number`, `device_id`, `base_url`)
- Set by prior request test scripts

### Chain Gaps ✅ CLEAN
1. **AuthN Chain:** OTP Send → (user sets code) → OTP Verify → Login → Session → Refresh → Logout
   - ✅ `mobile_otp_id` flows: AuthN-01 → AuthN-02, AuthN-03
   - ✅ `login_password` (device_credential) flows: AuthN-02 → AuthN-03, AuthZ auto-login
   - ✅ `access_token` flows: AuthN-03 → AuthN-04..06, AuthZ (with validation)

2. **AuthZ Auto-Login Check:** Lines 542-548
   ```javascript
   const existingToken = pm.environment.get('access_token') || '';
   const isValidJWT = existingToken.split('.').length === 3 && existingToken.length > 100;
   if (isValidJWT) { return; } // Reuse existing token
   ```
   - ✅ Validates JWT structure before reusing
   - ⚠️ **EDGE CASE:** If `access_token` is non-empty but malformed (e.g., "invalid"), it won't pass JWT check and will attempt auto-login

---

## 3. auth_smoke.postman_collection.json

### Variables READ
| Variable | Used In | Count |
|----------|---------|-------|
| `base_url` | All endpoints | 23 |
| `user_mobile_number` | OTP Send (01), OTP Verify (02), Login (03), Web Login (01), Regression (01, 02) | 5 |
| `login_device_type` | Login (03), Regression (02) | 2 |
| `device_id` | OTP Send (01), OTP Verify (02), Login (03), Current Session (04), Refresh (05), Email OTP Send (01), Email OTP Login (02), Regression tests | 8 |
| `mobile_otp_id` | OTP Verify (02), Login (03) | 2 |
| `mobile_otp_code` | OTP Verify (02) | 1 |
| `access_token` | Current Session (04), Refresh (05), AuthZ Check (01) | 3 |
| `refresh_token` | Refresh (05) | 1 |
| `session_id` | Logout (06) | 1 |
| `session_token` | Web Portal Current Session (02), Web Portal Logout (03), AuthZ Check (02) | 3 |
| `csrf_token` | Web Portal Current Session (02), Web Portal Logout (03), AuthZ Check (02) | 3 |
| `user_email` | Email OTP Send (01), Email OTP Login (02) | 2 |
| `email_otp_id` | Email OTP Login (02) | 1 |
| `email_login_otp_code` | Email OTP Login (02), Regression (03) | 2 |
| `authz_resource` | AuthZ Check (01, 02) | 2 |
| `authz_action` | AuthZ Check (01, 02) | 2 |

### Variables SET
| Variable | Set In | Source |
|----------|--------|--------|
| `mobile_otp_id` | 01 OTP Send test (line 71) | `j.data.otp_id` |
| `mobile_otp_id` | 02 OTP Verify test (line 137) | `j.data.otp_id` |
| `access_token` | 03 Login test (line 206) | `j.data.access_token` |
| `refresh_token` | 03 Login test (line 207) | `j.data.refresh_token` |
| `session_id` | 03 Login test (line 208) | `j.data.session_id` |
| `session_type` | 03 Login test (line 209) | `'JWT'` (hardcoded) |
| `access_token` | 05 Refresh Token test (line 329) | `j.data.access_token` |
| `refresh_token` | 05 Refresh Token test (line 330) | `j.data.refresh_token` |
| `session_token` | Web Portal Login test (line 468) | `j.data.session_token` |
| `csrf_token` | Web Portal Login test (line 469) | `j.data.csrf_token` |
| `session_id` | Web Portal Login test (line 470) | `j.data.session_id` |
| `session_type` | Web Portal Login test (line 471) | `'SERVER_SIDE'` (hardcoded) |
| `session_token` | Email OTP Login test (line 734) | `j.data.session_token` |
| `csrf_token` | Email OTP Login test (line 735) | `j.data.csrf_token` |
| `session_id` | Email OTP Login test (line 736) | `j.data.session_id` |
| `access_token` (UNSET) | Logout (JWT) test (line 390) | Cleared |
| `refresh_token` (UNSET) | Logout (JWT) test (line 390) | Cleared |
| `session_token` (UNSET) | Web Logout test (line 593) | Cleared |
| `csrf_token` (UNSET) | Web Logout test (line 593) | Cleared |

### Variables MISSING
**✅ NONE** — all used variables are pre-populated or set by prior requests.

### Chain Gaps ✅ MOSTLY CLEAN
1. **B2C Mobile OTP → JWT:**
   - ✅ 01 OTP Send → sets `mobile_otp_id`
   - ✅ 02 OTP Verify → reads `mobile_otp_id`, `mobile_otp_code` → re-sets `mobile_otp_id`
   - ✅ 03 Login → reads `mobile_otp_id` → sets `access_token`, `refresh_token`, `session_id`
   - ✅ 04 Session → reads `access_token`, `session_id`
   - ✅ 05 Refresh → reads `refresh_token` → updates `access_token`, `refresh_token`
   - ✅ 06 Logout → reads `access_token`, `session_id` → clears both

2. **Web Portal Password → SERVER_SIDE:**
   - ✅ 01 Login → reads `login_password` → sets `session_token`, `csrf_token`, `session_id`
   - ✅ 02 Session → reads `session_token`, `csrf_token`, `session_id`
   - ✅ 03 Logout → reads `session_token`, `csrf_token`, `session_id` → clears tokens

3. **Email OTP → SERVER_SIDE:**
   - ✅ 01 Email OTP Send → sets `email_otp_id`
   - ✅ 02 Email OTP Login → reads `email_otp_id`, `email_login_otp_code` → sets `session_token`, `csrf_token`, `session_id`
   - ✅ Regression-03 validates no JWT returned

4. **AuthZ Permission Checks:**
   - ✅ 01 AuthZ Check (JWT) → reads `access_token`, `authz_resource`, `authz_action`
   - ✅ 02 AuthZ Check (SERVER_SIDE) → reads `session_token`, `csrf_token`, `authz_resource`, `authz_action`

---

## 4. CRITICAL GAPS & RECOMMENDATIONS

### Gap 1: Email OTP Variables Missing in newman_test Environment ⚠️
**Severity:** MEDIUM
- `newman_test.postman_environment.json` does NOT have email OTP variables
- auth_smoke.postman_collection has Email OTP test folder that requires `user_email`, `email_otp_id`, `email_login_otp_code`, `email_verification_otp_code`
- **Impact:** Email OTP section of auth_smoke will fail in Newman CI unless manually populated
- **Fix:** Add to `InsureTech_newman_test.postman_environment.json`:
  ```json
  {
    "key": "user_email",
    "value": "test.user@labaidinsuretech.com",
    "enabled": true,
    "type": "default"
  },
  {
    "key": "email_otp_id",
    "value": "",
    "enabled": true,
    "type": "default"
  },
  {
    "key": "email_login_otp_code",
    "value": "",
    "enabled": true,
    "type": "secret"
  },
  {
    "key": "email_verification_otp_code",
    "value": "",
    "enabled": true,
    "type": "secret"
  },
  {
    "key": "email_verified",
    "value": "false",
    "enabled": true,
    "type": "default"
  }
  ```

### Gap 2: Session Token Handling Inconsistency ⚠️
**Severity:** LOW
- `local` env has explicit `session_token` and `csrf_token` variables
- `newman_test` env does NOT have them (relies on implicit cookie handling?)
- **Impact:** WEB device_type login (auth_smoke section 2) may fail in newman_test because scripts expect to set `session_token`, `csrf_token` but these vars don't exist
- **Fix:** Add `session_token` and `csrf_token` to `InsureTech_newman_test.postman_environment.json`

### Gap 3: JWT Validation Logic in AuthZ Pre-Request ⚠️
**Severity:** LOW (defensive coding)
- **Location:** b2c_authn_authz_suite.json, lines 542-548
- **Issue:** If `access_token` is set to a garbage string (not 3 dot-separated parts), JWT check fails silently and auto-login attempts
- **Current behavior:** Smart retry-on-invalid, acceptable
- **Recommendation:** Add explicit logging or error handling:
  ```javascript
  const existingToken = pm.environment.get('access_token') || '';
  const isValidJWT = existingToken.split('.').length === 3 && existingToken.length > 100;
  
  if (!isValidJWT && existingToken) {
    console.warn('Malformed access_token detected, will attempt auto-login');
  }
  ```

### Gap 4: No Explicit Device Type Validation in Pre-Request ⚠️
**Severity:** LOW
- Both collections allow `login_device_type` to be any string in env
- **Risk:** If accidentally set to "INVALID" or "WEB" in wrong flow, behavior is undefined
- **Current mitigation:** Test scripts validate response `session_type` matches expected value
- **Recommendation:** Add pre-request validation:
  ```javascript
  const devType = pm.environment.get('login_device_type');
  if (!['ANDROID', 'IOS', 'API', 'WEB'].includes(devType)) {
    console.error('Invalid login_device_type: ' + devType);
  }
  ```

### Gap 5: mobile_otp_code Requires Manual Input ⚠️
**Severity:** MEDIUM (by design)
- Both collections expect `mobile_otp_code` to be manually set after OTP send
- **Problem:** Newman CI cannot auto-populate SMS OTPs without SMS API integration
- **Current state:** Collections have pre-request checks that warn if code not set (lines 124, 118-120 in auth_smoke)
- **Acceptable:** This is architectural — test automation should use SMS mock or test API endpoints

---

## 5. ENVIRONMENT READINESS CHECKLIST

### newman_test Environment
| Category | Status | Action |
|----------|--------|--------|
| Auth tokens | ✅ Complete | None |
| Mobile OTP | ✅ Complete | None |
| Email OTP | ❌ MISSING | Add 5 variables |
| Session (WEB) | ❌ MISSING | Add 2 variables |
| User data | ✅ Complete | None |
| Device info | ✅ Complete | None |
| AuthZ params | ✅ Complete | None |

### local Environment
| Category | Status | Action |
|----------|--------|--------|
| All auth flows | ✅ Complete | None |
| Business entities | ✅ Complete (for future use) | None |
| Email verification | ✅ Complete | None |

### Collections
| Collection | Status | Issues |
|-----------|--------|--------|
| b2c_authn_authz_suite | ✅ HEALTHY | No missing variables; proper error handling for rate limits |
| auth_smoke | ✅ HEALTHY | No missing variables; proper pre-request guards and cleanup |

---

## 6. SUMMARY TABLE: Variable Propagation Chains

### B2C Mobile OTP → JWT (Both Collections)
```
AuthN-01: OTP Send
  INPUT:  user_mobile_number, mobile_otp_type, mobile_otp_channel
  OUTPUT: mobile_otp_id ✅

AuthN-02: OTP Verify
  INPUT:  mobile_otp_id, mobile_otp_code (manual)
  OUTPUT: device_credential → login_password ✅

AuthN-03: Login
  INPUT:  user_mobile_number, login_password, device_id
  OUTPUT: access_token, refresh_token, session_id, user_id ✅

AuthN-04..06: Session, Refresh, Logout
  INPUT:  access_token, refresh_token, session_id
  STATUS: ✅ All vars flow correctly
```

### Web Portal Password → SERVER_SIDE (auth_smoke only)
```
Web-01: Login
  INPUT:  user_mobile_number, login_password (pre-set), device_type=WEB
  OUTPUT: session_token, csrf_token, session_id ✅

Web-02..03: Session, Logout
  INPUT:  session_token, csrf_token, session_id
  STATUS: ✅ All vars flow correctly
```

### Email OTP → SERVER_SIDE (auth_smoke only)
```
Email-01: OTP Send
  INPUT:  user_email
  OUTPUT: email_otp_id ✅

Email-02: Email Login
  INPUT:  email_otp_id, email_login_otp_code (manual)
  OUTPUT: session_token, csrf_token, session_id ✅
```

### AuthZ Enforcement (Both Collections)
```
AuthZ Pre-Request (b2c_authn_authz_suite folder)
  INPUT:  access_token (or auto-login: login_password, user_mobile_number, device_id)
  OUTPUT: access_token (refreshed if needed) ✅

AuthZ-01..10: Permission Checks
  INPUT:  access_token, user_id
  STATUS: ✅ All vars available from AuthN chain
```

---

## FINAL RECOMMENDATIONS

### Priority 1 (Do Immediately)
- [ ] Add email OTP variables to `InsureTech_newman_test.postman_environment.json`
- [ ] Add `session_token`, `csrf_token` to `InsureTech_newman_test.postman_environment.json`
- [ ] Test auth_smoke collection with newman_test environment after above fixes

### Priority 2 (Nice to Have)
- [ ] Add defensive pre-request checks for `login_device_type` enum validation
- [ ] Document in collection README that `mobile_otp_code` and `email_*_otp_code` require manual population from SMS/Email
- [ ] Add logging for JWT validation edge cases

### Priority 3 (Documentation)
- [ ] Create a "Variable Flow Diagram" showing OTP → Device Binding → JWT chain
- [ ] Document environment-specific use cases:
  - `newman_test`: CI/CD automation (no email testing)
  - `local`: Full feature testing (all auth flows)
  - `staging`, `production`: Integration testing with real services

---

**Analysis Date:** 2025-01-16  
**Collections Analyzed:** 3  
**Total Variables Tracked:** 46  
**Gap Severity:** 2 MEDIUM, 3 LOW, 0 CRITICAL
