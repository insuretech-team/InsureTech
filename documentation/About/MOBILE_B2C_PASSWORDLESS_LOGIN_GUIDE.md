# Mobile B2C Passwordless Login Guide

Quick-start guide for testing passwordless login flows when OTP is already verified.

**Use this when:** OTP verification is complete and `device_credential` is already saved.

---

## Pre-Requisites: Environment Setup

Before starting, ensure these environment variables are set from your previous OTP verification:

| Variable | Value | How to get it |
|----------|-------|---------------|
| `user_mobile_number` | `+8801XXXXXXXXX` | Your test phone number |
| `device_id` | `my-android-dev-01` | Same device ID used in OTP verify |
| `device_name` | `Pixel 8 Dev` | Same device name used before |
| `login_device_type` | `ANDROID` or `IOS` | Same device type as before |
| `mobile_otp_verified` | `true` | OTP verify response |
| `device_credential` | (long string) | From OTP verify response |
| `user_id` | `usr_...` | From OTP verify response |

### Required Environment File

Use: `api/postman/InsureTech_production.postman_environment.json`

Set these manually before running:

```json
{
  "user_mobile_number": "+8801347201751",
  "device_id": "my-android-dev-01",
  "device_name": "Pixel 8 Dev",
  "login_device_type": "ANDROID",
  "mobile_otp_verified": "true",
  "device_credential": "YOUR_DEVICE_CREDENTIAL_FROM_OTP_VERIFY",
  "user_id": "YOUR_USER_ID_FROM_OTP_VERIFY"
}
```

---

## Test Flow (5 Steps)

### Step 1: Login Passwordless (OTP)

**Request:** `03 Login Passwordless (OTP)`

**Endpoint:** `POST {{base_url}}/v1/auth/login`

**Request Body:**
```json
{
  "mobile_number": "{{user_mobile_number}}",
  "password": "",
  "device_id": "{{device_id}}",
  "device_type": "{{login_device_type}}",
  "device_name": "{{device_name}}"
}
```

**Expected Result:**
- `200 OK`
- `session_type = JWT`
- Response contains:
  - `access_token`
  - `refresh_token`
  - `session_id`
  - `user_id`

**Postman automatically saves:**
- `access_token`
- `refresh_token`
- `session_id`
- `user_id`
- `session_token` (if provided)

---

### Step 2: Login Passwordless (Device-Bound)

**Request:** `04 Login Passwordless (Device-Bound)`

**Endpoint:** `POST {{base_url}}/v1/auth/login`

**Request Body:**
```json
{
  "mobile_number": "{{user_mobile_number}}",
  "password": "{{device_credential}}",
  "device_id": "{{device_id}}",
  "device_type": "{{login_device_type}}",
  "device_name": "{{device_name}}"
}
```

**Key difference:** `password` is set to `device_credential` from OTP verification.

**Expected Result:**
- `200 OK`
- `session_type = JWT`
- Login succeeds **without OTP** because device is bound
- New tokens are saved to environment

**This proves:** Device binding works and can be used for passwordless re-login.

---

### Step 3: Current Session

**Request:** `05 Current Session`

**Endpoint:** `GET {{base_url}}/v1/auth/session/current`

**Headers:**
- `Authorization: Bearer {{access_token}}`

**Expected Result:**
- `200 OK`
- Returns current active session details

---

### Step 4: Refresh Token

**Request:** `06 Refresh Token`

**Endpoint:** `POST {{base_url}}/v1/auth/token:refresh`

**Request Body:**
```json
{
  "refresh_token": "{{refresh_token}}"
}
```

**Expected Result:**
- `200 OK`
- New `access_token` (saved automatically)
- New `refresh_token` (saved automatically)
- New `session_id` (saved automatically)

**Important:** After this step, the new tokens are used for all subsequent requests.

---

### Step 5: Logout and Re-Login Test

#### 5a. Logout

**Request:** `AuthZ-02 Logout`

**Endpoint:** `POST {{base_url}}/v1/auth/logout`

**Expected Result:**
- `200 OK` or `204 No Content`

#### 5b. Re-Login Passwordless (Device-Bound)

**Request:** `AuthZ-03 Re-Login Passwordless (Device-Bound)`

**Endpoint:** `POST {{base_url}}/v1/auth/login`

**Request Body:**
```json
{
  "mobile_number": "{{user_mobile_number}}",
  "password": "{{device_credential}}",
  "device_id": "{{device_id}}",
  "device_type": "{{login_device_type}}",
  "device_name": "{{device_name}}"
}
```

**Expected Result:**
- `200 OK`
- **No OTP required** - device binding survives logout
- New session tokens are saved

**This proves:** Device binding persists across logout/login cycles.

#### 5c. Final Logout

**Request:** `AuthZ-05 Final Logout`

**Endpoint:** `POST {{base_url}}/v1/auth/logout`

**Expected Result:**
- `200 OK` or `204`

---

## Optional: Profile & Session Tests

After Step 4 (Refresh Token), you can run these:

### Create Profile
**Request:** `Profile-01 Create Profile`
- Endpoint: `POST {{base_url}}/v1/auth/users/{{user_id}}/profile`
- May return `409` if profile already exists (acceptable)

### Get Profile
**Request:** `Profile-02 Get Profile`
- Endpoint: `GET {{base_url}}/v1/auth/users/{{user_id}}/profile`

### Update Profile
**Request:** `Profile-03 Update Profile`
- Endpoint: `PATCH {{base_url}}/v1/auth/users/{{user_id}}/profile`

### List User Sessions
**Request:** `Session-02 List User Sessions`
- Endpoint: `GET {{base_url}}/v1/auth/users/{{user_id}}/sessions`

### Get Session By ID
**Request:** `Session-01 Get Session By ID`
- Endpoint: `GET {{base_url}}/v1/auth/sessions/{{session_id}}`

---

## Quick Newman Command

Run just the passwordless login flow:

```bash
cd api/postman

newman run b2c_authn_authz_suite.postman_collection.json \
  -e InsureTech_production.postman_environment.json \
  --folder "1. B2C AuthN" \
  --request "03 Login Passwordless (OTP),04 Login Passwordless (Device-Bound),05 Current Session,06 Refresh Token"
```

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| `401 password or verified OTP required` | `mobile_otp_verified` is not `true` or OTP expired - re-run OTP flow |
| `401 invalid credentials` | `device_credential` is wrong or expired - re-run OTP verify |
| `403 device not bound` | Device ID changed or device_credential invalid - use same device values from OTP verify |
| `404 user not found` | Wrong `user_mobile_number` - check phone number format |
| Stale token errors | Run `06 Refresh Token` to get fresh tokens |

---

## Files To Use

| File | Purpose |
|------|---------|
| `api/postman/b2c_authn_authz_suite.postman_collection.json` | Main collection |
| `api/postman/InsureTech_production.postman_environment.json` | Environment |

**Note:** If you need to regenerate the collection, use `api/generator/postman_sync_core.py`