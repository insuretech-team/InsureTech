# Mobile B2C Testing Guide

Simple step-by-step guide for testing the mobile B2C auth flow with Postman.

Last validated: March 26, 2026

## What This Guide Covers

1. Send OTP
2. Verify OTP
3. Login passwordless with verified OTP
4. Login again passwordless with the bound device credential
5. Check current session
6. Refresh token
7. Get session by ID
8. Create profile
9. Get profile
10. Update profile
11. List user sessions
12. Check AuthZ
13. Logout
14. Login again after logout with the same bound device and no OTP

This is the flow that was tested live against production.

## Files To Use

Collection:

- `api/postman/b2c_authn_authz_suite.postman_collection.json`

Environment:

- `api/postman/InsureTech_production.postman_environment.json`

Source of truth:

- `api/generator/postman_sync_core.py`

If you change the Postman flow, change the generator and rerun the API pipeline. Do not hand-edit the generated collection.

## Requests In The Suite

### 1. B2C AuthN

1. `01 OTP Send`
2. `02 OTP Verify`
3. `03 Login Passwordless (OTP)`
4. `04 Login Passwordless (Device-Bound)`
5. `05 Current Session`
6. `06 Refresh Token`

### 2. B2C Profile & Sessions

1. `Session-01 Get Session By ID`
2. `Profile-01 Create Profile`
3. `Profile-02 Get Profile`
4. `Profile-03 Update Profile`
5. `Session-02 List User Sessions`

### 3. B2C AuthZ

1. `AuthZ-01 Check Access`
2. `AuthZ-02 Logout`
3. `AuthZ-03 Re-Login Passwordless (Device-Bound)`
4. `AuthZ-04 Current Session After Re-Login`
5. `AuthZ-05 Final Logout`

## Endpoint Map

### Auth Endpoints

1. `01 OTP Send`
   Endpoint: `POST {{base_url}}/v1/auth/otp:send`
2. `02 OTP Verify`
   Endpoint: `POST {{base_url}}/v1/auth/otp:verify`
3. `03 Login Passwordless (OTP)`
   Endpoint: `POST {{base_url}}/v1/auth/login`
4. `04 Login Passwordless (Device-Bound)`
   Endpoint: `POST {{base_url}}/v1/auth/login`
5. `05 Current Session`
   Endpoint: `GET {{base_url}}/v1/auth/session/current`
6. `06 Refresh Token`
   Endpoint: `POST {{base_url}}/v1/auth/token:refresh`

### Profile & Session Endpoints

1. `Session-01 Get Session By ID`
   Endpoint: `GET {{base_url}}/v1/auth/sessions/{{session_id}}`
2. `Profile-01 Create Profile`
   Endpoint: `POST {{base_url}}/v1/auth/users/{{user_id}}/profile`
3. `Profile-02 Get Profile`
   Endpoint: `GET {{base_url}}/v1/auth/users/{{user_id}}/profile`
4. `Profile-03 Update Profile`
   Endpoint: `PATCH {{base_url}}/v1/auth/users/{{user_id}}/profile`
5. `Session-02 List User Sessions`
   Endpoint: `GET {{base_url}}/v1/auth/users/{{user_id}}/sessions`

### AuthZ Endpoints

1. `AuthZ-01 Check Access`
   Endpoint: `POST {{base_url}}/v1/authz/check`
2. `AuthZ-02 Logout`
   Endpoint: `POST {{base_url}}/v1/auth/logout`
3. `AuthZ-03 Re-Login Passwordless (Device-Bound)`
   Endpoint: `POST {{base_url}}/v1/auth/login`
4. `AuthZ-04 Current Session After Re-Login`
   Endpoint: `GET {{base_url}}/v1/auth/session/current`
5. `AuthZ-05 Final Logout`
   Endpoint: `POST {{base_url}}/v1/auth/logout`

## Before You Start

Open the production environment and set these values:

| Variable | Example |
|---|---|
| `base_url` | `https://api.labaidinsuretech.com` |
| `user_mobile_number` | `+8801XXXXXXXXX` |
| `device_id` | `my-android-dev-01` |
| `device_name` | `Pixel 8 Dev` |
| `login_device_type` | `ANDROID` or `IOS` |

Keep the device stable while testing:

- use the same `device_id`
- use the same `device_name`
- use the same `login_device_type`

Leave these alone unless you have a special reason:

- `mobile_otp_type = login`
- `mobile_otp_channel = sms`
- `mobile_otp_use_masking = true`

Do not set `login_password` manually for this flow.

## Step By Step

### Step 1: Run `01 OTP Send`

Endpoint:

- `POST {{base_url}}/v1/auth/otp:send`

Expected result:

- `200 OK`
- Postman saves `mobile_otp_id`
- Postman saves `otp_id`

If SMS delivery works, the user should receive a 6-digit OTP.

### Step 2: Enter OTP And Run `02 OTP Verify`

Endpoint:

- `POST {{base_url}}/v1/auth/otp:verify`

Set:

- `mobile_otp_code`

Expected result:

- `200 OK`
- Postman sets `mobile_otp_verified = true`
- response returns `device_credential`
- Postman saves `device_credential`

This step is what binds the device for future passwordless re-login.

### Step 3: Run `03 Login Passwordless (OTP)`

Endpoint:

- `POST {{base_url}}/v1/auth/login`

This request logs in with:

- `mobile_number`
- empty `password`
- same `device_id`
- same mobile `device_type`

Expected result:

- `200 OK`
- `session_type = JWT`
- response returns `access_token`
- response returns `refresh_token`
- response returns `session_id`
- response returns `user_id`

Postman now saves those values automatically.

### Step 4: Run `04 Login Passwordless (Device-Bound)`

Endpoint:

- `POST {{base_url}}/v1/auth/login`

This is the important mobile re-login test.

This request logs in with:

- the same `mobile_number`
- `password = {{device_credential}}`
- the same `device_id`
- the same mobile `device_type`

Expected result:

- `200 OK`
- `session_type = JWT`
- device-bound passwordless login succeeds without OTP

This proves the device is bound and can be used for passwordless re-login.

### Step 5: Run `05 Current Session`

Endpoint:

- `GET {{base_url}}/v1/auth/session/current`

Expected result:

- `200 OK`

### Step 6: Run `06 Refresh Token`

Endpoint:

- `POST {{base_url}}/v1/auth/token:refresh`

Expected result:

- `200 OK`
- new `access_token`
- new `refresh_token`
- new `session_id`

Postman now saves the refreshed values automatically, so later requests use the latest session.

### Step 7: Run `Session-01 Get Session By ID`

Endpoint:

- `GET {{base_url}}/v1/auth/sessions/{{session_id}}`

Expected result:

- `200 OK`

### Step 8: Run `Profile-01 Create Profile`

Endpoint:

- `POST {{base_url}}/v1/auth/users/{{user_id}}/profile`

Expected result:

- `201 Created` on first run
- `409 Conflict` on later runs if the profile already exists

The suite treats both as acceptable.

### Step 9: Run `Profile-02 Get Profile`

Endpoint:

- `GET {{base_url}}/v1/auth/users/{{user_id}}/profile`

Expected result:

- `200 OK`

### Step 10: Run `Profile-03 Update Profile`

Endpoint:

- `PATCH {{base_url}}/v1/auth/users/{{user_id}}/profile`

Expected result:

- `200 OK`

### Step 11: Run `Session-02 List User Sessions`

Endpoint:

- `GET {{base_url}}/v1/auth/users/{{user_id}}/sessions`

Expected result:

- `200 OK`

### Step 12: Run `AuthZ-01 Check Access`

Endpoint:

- `POST {{base_url}}/v1/authz/check`

Expected result:

- `200 OK`

### Step 13: Run `AuthZ-02 Logout`

Endpoint:

- `POST {{base_url}}/v1/auth/logout`

Expected result:

- `200 OK` or `204`

This logs out the current JWT session.

### Step 14: Test Re-Login After Logout

After logout, run these again:

1. `AuthZ-03 Re-Login Passwordless (Device-Bound)`
   Endpoint: `POST {{base_url}}/v1/auth/login`
2. `AuthZ-04 Current Session After Re-Login`
   Endpoint: `GET {{base_url}}/v1/auth/session/current`
3. `AuthZ-05 Final Logout`
   Endpoint: `POST {{base_url}}/v1/auth/logout`

Expected result:

- re-login succeeds without OTP
- current session succeeds
- logout succeeds again

This is the key proof that device binding survives logout and the same device can log in again passwordless.

## Live-Tested Result

The full production run passed for `+8801347201751` with the same device:

- OTP send passed
- OTP verify passed
- OTP-based passwordless login passed
- device-bound passwordless login passed
- current session passed
- refresh token passed
- get session by ID passed
- create profile returned acceptable existing-profile result
- get profile passed
- update profile passed
- list user sessions passed
- authz check passed
- logout passed
- post-logout device-bound passwordless re-login passed
- post-logout session check passed
- final logout passed

## Common Mistakes

| Problem | Meaning | Fix |
|---|---|---|
| `01 OTP Send` returns `400 unsupported channel` | wrong channel value | use `sms`, not enum text like `OTP_CHANNEL_SMS` |
| `03 Login Passwordless (OTP)` fails with `password or verified OTP required` | OTP was not verified or has expired | rerun OTP send and verify |
| `04 Login Passwordless (Device-Bound)` fails | device is not bound or device values changed | use the same `device_id` and rerun OTP verify |
| `Profile-01 Create Profile` returns `409` | profile already exists for that user | this is acceptable for repeat runs |
| `AuthZ` or `Logout` fails right after refresh | stale token/session in environment | make sure refresh response values are being saved |
| post-logout re-login fails | device binding is broken or device values changed | retry with the same `device_id`; if it still fails, treat as backend bug |

## Newman Example

