# Postman Auth Test Setup

Use this file with the generated Postman artifacts in `api/postman`.

## Files To Import

- [ ] Import `api/postman/InsureTech.postman_collection.json`
- [ ] Import one environment file:
- [ ] `api/postman/InsureTech_local.postman_environment.json`
- [ ] `api/postman/InsureTech_staging.postman_environment.json`
- [ ] `api/postman/InsureTech_production.postman_environment.json`

## What The Collection Already Does

- [ ] Pre-request scripts read values from request variables first, then environment variables
- [ ] If `user_mobile_number` or `user_email` is already set in the environment, developers do not need to type them into every request
- [ ] `POST /v1/auth/otp:send` captures `mobile_otp_id`
- [ ] `POST /v1/auth/otp:verify` captures `mobile_otp_verified`
- [ ] `POST /v1/auth/email/otp:send` captures `email_login_otp_id`
- [ ] Login requests automatically capture `user_id`, `session_id`, and token fields
- [ ] `POST /v1/auth/login` blocks OTP-only login unless `mobile_otp_verified=true`
- [ ] `POST /v1/auth/login` auto-checks expected session type through `login_expect_session_type`

## Minimum Environment Variables

Set these once in the imported environment:

- [ ] `base_url`
- [ ] `user_mobile_number` for mobile OTP / mobile password login
- [ ] `user_email` for email OTP login
- [ ] `device_id` if you want a stable device fingerprint
- [ ] `device_name` if you want readable session metadata

Useful auth-specific variables:

- [ ] `login_password`
- [ ] `login_device_type`
- [ ] `login_expect_session_type`
- [ ] `mobile_otp_code`
- [ ] `email_login_otp_code`

## Recommended Defaults

- [ ] `login_password=""` for OTP-only B2C mobile login
- [ ] `login_device_type=API` for JWT-style app testing
- [ ] `login_expect_session_type=JWT` for mobile/API login
- [ ] `email_otp_type=email_login` for email login OTP flow
- [ ] `mobile_otp_type=login`
- [ ] `mobile_otp_channel=sms`

## Mobile OTP Login Flow

This is the easiest real-world B2C login test.

- [ ] Set `user_mobile_number`
- [ ] Keep `login_password` empty
- [ ] Set `login_device_type` to `API`, `ANDROID`, `IOS`, `MOBILE_ANDROID`, or `MOBILE_IOS`
- [ ] Run `POST /v1/auth/otp:send`
- [ ] Enter the received OTP into `mobile_otp_code`
- [ ] Run `POST /v1/auth/otp:verify`
- [ ] Confirm `mobile_otp_verified=true`
- [ ] Run `POST /v1/auth/login`
- [ ] Confirm `session_type=JWT`
- [ ] Confirm `access_token`, `refresh_token`, and `session_id` were captured
- [ ] Run `GET /v1/auth/session/current`
- [ ] Confirm current session returns a JWT-backed session
- [ ] Run `POST /v1/auth/sessions/refresh`
- [ ] Confirm a new access token is returned

## Email OTP Login Flow

- [ ] Set `user_email`
- [ ] Run `POST /v1/auth/email/otp:send`
- [ ] Enter the received OTP into `email_login_otp_code`
- [ ] Run `POST /v1/auth/email/login`
- [ ] Confirm `session_type=SERVER_SIDE`
- [ ] Confirm `session_token`, `csrf_token`, and `session_id` were captured

## Password Login Flows

JWT-style mobile/API password login:

- [ ] Set `login_password` to a valid password
- [ ] Set `login_device_type=API`
- [ ] Run `POST /v1/auth/login`
- [ ] Confirm `session_type=JWT`

Web password login:

- [ ] Set `login_password` to a valid password
- [ ] Set `login_device_type=WEB`
- [ ] Set `login_expect_session_type=SERVER_SIDE`
- [ ] Run `POST /v1/auth/login`
- [ ] Confirm `session_type=SERVER_SIDE`
- [ ] Confirm `session_token` and `csrf_token` were captured

## Variables Captured Automatically

- [ ] `user_id`
- [ ] `session_id`
- [ ] `access_token`
- [ ] `refresh_token`
- [ ] `session_token`
- [ ] `csrf_token`
- [ ] `mobile_otp_id`
- [ ] `mobile_otp_verified`
- [ ] `email_login_otp_id`

## Common Failure Messages

- [ ] `OTP-only B2C login requires a verified OTP. Run POST /v1/auth/otp:verify first.`
- [ ] Means `login_password=""` was used before `mobile_otp_verified=true`

- [ ] `password is required`
- [ ] Means `device_type=WEB` was used with an empty password

- [ ] `invalid credentials`
- [ ] Means password is wrong or there is no recent verified OTP for OTP-only mobile login

## Quick Smoke Checklist

- [ ] Mobile OTP send works
- [ ] Mobile OTP verify works
- [ ] Mobile OTP login with `password=""` returns JWT
- [ ] WEB login with `password=""` is rejected
- [ ] Email OTP send works
- [ ] Email OTP login returns server-side session
- [ ] JWT session can call `GET /v1/auth/session/current`
- [ ] Refresh token flow works

## Source Of Truth

- Generated collection: `api/postman/InsureTech.postman_collection.json`
- Generated environments: `api/postman/InsureTech_*postman_environment.json`
- Generator script: `api/generator/sync_postman.py`
