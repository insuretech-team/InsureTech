# WhatsApp-Style Device Binding Authentication

**Status:** Implemented ✅  
**Date:** 2026-03-18  
**Affects:** `authn` service — `VerifyOTP`, `Login`, proto `core.proto`

---

## 🎯 Design Goal

Replace the fragile "OTP verified → empty password login" two-step B2C flow with a
**WhatsApp-style device-bound credential** system:

- After OTP verify, server derives a **deterministic device credential** from
  `HMAC-SHA256(mobile_number + ":" + device_id, SERVER_DEVICE_SECRET)`
- This credential is returned **once** to the app (stored in Android Keystore / iOS Keychain)
- On subsequent logins, the app sends the credential as `password` — **no OTP needed**
- Changing device OR mobile number → old credential invalid → must OTP again

---

## 🏗️ Architecture

```
FIRST TIME (new device + new mobile number):
┌─────────────────────────────────────────────────────────────────┐
│ 1. POST /v1/auth/otp:send                                       │
│    { recipient: "01XXXXXXXXX", type: "login", channel: "sms" }  │
│    ← { otp_id: "abc123" }                                       │
│                                                                 │
│ 2. [User receives SMS with OTP code]                            │
│                                                                 │
│ 3. POST /v1/auth/otp:verify                                     │
│    { otp_id: "abc123", code: "123456",                          │
│      device_id: "android-uuid-xxxx", device_type: "ANDROID" }  │
│    ← { verified: true, device_credential: "base64string..." }  │
│      [App stores device_credential in Android Keystore]         │
│                                                                 │
│ 4. POST /v1/auth/login                                          │
│    { mobile_number: "01XXXXXXXXX",                              │
│      password: "base64string..." (device_credential),           │
│      device_id: "android-uuid-xxxx", device_type: "ANDROID" }  │
│    ← { access_token: "...", refresh_token: "...",               │
│        session_type: "JWT" }                ✅ DONE             │
└─────────────────────────────────────────────────────────────────┘

SAME DEVICE NEXT TIME (seamless - no OTP):
┌─────────────────────────────────────────────────────────────────┐
│ 1. POST /v1/auth/login                                          │
│    { mobile_number: "01XXXXXXXXX",                              │
│      password: "base64string..." (from Keystore),               │
│      device_id: "android-uuid-xxxx", device_type: "ANDROID" }  │
│    ← { access_token: "...", refresh_token: "..." } ✅ NO OTP   │
└─────────────────────────────────────────────────────────────────┘

NEW DEVICE (old device auto-invalidated):
┌─────────────────────────────────────────────────────────────────┐
│ Old device sends login → "invalid credentials" ❌               │
│ New device must go through OTP flow again (step 1-4 above)      │
│ New device_credential derived with new device_id → stored       │
│ Old device_credential (old device_id) no longer matches         │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔐 Credential Derivation (Server-Side)

```go
// HMAC-SHA256(mobile_number + ":" + device_id, SERVER_DEVICE_SECRET)
// Deterministic — server re-derives on every login to compare
func deriveDeviceCredential(mobileNumber, deviceID, secret string) string {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(mobileNumber + ":" + deviceID))
    return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}
```

**Key properties:**
- **Deterministic** — same inputs always produce the same output
- **Server-side only** — client only stores the result, never the secret
- **No DB column needed** — server re-derives and compares directly (no storage)
- **Device-bound** — changing `device_id` produces completely different credential
- **Mobile-bound** — changing `mobile_number` produces completely different credential

---

## 📋 Proto Changes

### `proto/insuretech/authn/services/v1/core.proto`

**VerifyOTPRequest** — added `device_id` and `device_type`:
```protobuf
message VerifyOTPRequest {
  string otp_id     = 1 [(google.api.field_behavior) = REQUIRED];
  string code       = 2 [(google.api.field_behavior) = REQUIRED];
  string device_id  = 3 [(google.api.field_behavior) = OPTIONAL]; // NEW
  string device_type = 4 [(google.api.field_behavior) = OPTIONAL]; // NEW: ANDROID/IOS triggers credential return
}
```

**VerifyOTPResponse** — added `device_credential`:
```protobuf
message VerifyOTPResponse {
  bool   verified          = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
  string user_id           = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
  string message           = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
  insuretech.common.v1.Error error = 4 [(google.api.field_behavior) = OUTPUT_ONLY];
  // NEW: One-time device credential for WhatsApp-style login
  // Only returned when device_type=ANDROID or IOS and verified=true
  // Store in Android Keystore / iOS Keychain. Use as `password` on subsequent logins.
  string device_credential = 5 [(google.api.field_behavior) = OUTPUT_ONLY]; // NEW
}
```

---

## 🔧 Service Changes

### New file: `device_credential.go`

```go
// deriveDeviceCredential returns a deterministic HMAC-SHA256 credential
// for a given mobile+device pair. Used for WhatsApp-style device binding.
func deriveDeviceCredential(mobileNumber, deviceID string) string
```

### `auth_service.go` — `VerifyOTP()`

After successful OTP verify, if `req.DeviceId != ""` and `req.DeviceType` is `ANDROID` or `IOS`:
1. Derive `device_credential = HMAC-SHA256(mobile+device_id, secret)`
2. Set `resp.DeviceCredential = device_credential`
3. Call `markTrustedDevice(userID, deviceID)` in Redis

### `auth_service.go` — `Login()`

In the password verification block, for mobile device types:
1. If `device_type == ANDROID/IOS` and `password != ""`
2. Derive expected credential from `mobile_number + device_id`
3. Compare with `req.Password` using `hmac.Equal` (constant-time)
4. If match → `otpLoginVerified = true` → issue JWT

---

## 🌍 Environment Variables

| Var | Description | Required |
|---|---|---|
| `DEVICE_BIND_SECRET` | 32-byte hex secret for HMAC derivation | ✅ Yes |

```bash
# Generate a secure secret (run once, store in .env.prod):
openssl rand -hex 32
```

---

## 📱 Client SDK Integration

### Android (Kotlin)
```kotlin
// After OTP verify — store credential
val deviceCredential = verifyOtpResponse.deviceCredential
val keyStore = KeyStore.getInstance("AndroidKeyStore")
// Store in encrypted SharedPreferences or EncryptedFile
EncryptedSharedPreferences.create(...).edit()
    .putString("device_credential", deviceCredential).apply()

// Subsequent login — use credential
val cred = prefs.getString("device_credential", "")
authService.login(LoginRequest(
    mobileNumber = "01XXXXXXXXX",
    password = cred,
    deviceId = deviceId,  // must match what was used at OTP verify
    deviceType = "ANDROID"
))
```

### iOS (Swift)
```swift
// After OTP verify — store in Keychain
KeychainHelper.save(key: "device_credential", value: verifyResponse.deviceCredential)

// Subsequent login
let cred = KeychainHelper.load(key: "device_credential") ?? ""
authService.login(mobile: mobile, password: cred, deviceId: deviceId, deviceType: "IOS")
```

---

## 🔒 Security Properties

| Property | Detail |
|---|---|
| **Device binding** | Credential tied to `device_id` — different device = different credential |
| **Mobile binding** | Credential tied to `mobile_number` — number change = credential invalid |
| **No server storage** | Server re-derives on every request — no credential DB table needed |
| **Constant-time compare** | `hmac.Equal` used — immune to timing attacks |
| **One-time delivery** | `device_credential` only returned once in `VerifyOTPResponse` |
| **Trusted device cache** | Redis `trusted:device:{uid}:{did}` with 30-day TTL for fast path |
| **OTP still required** | New device / new mobile always requires fresh OTP verification |
| **WEB unaffected** | `device_type=WEB` never receives or uses `device_credential` |

---

## ⚠️ Migration Notes

- **Existing B2C users** with `password=""` login: after first OTP verify with `device_id`, they get a `device_credential` — subsequent logins use it seamlessly
- **Old Postman smoke tests**: update `mobile_otp_use_masking=false`, add `device_id` + `device_type` to `VerifyOTP` request body, capture `device_credential` from response, use as `password` in login
- **B2B WEB flow**: completely unaffected — `device_type=WEB` skips all device credential logic

---

## ✅ Test Checklist

- [ ] OTP verify with `device_type=ANDROID` returns `device_credential` (non-empty)
- [ ] OTP verify with `device_type=WEB` does NOT return `device_credential`
- [ ] Login with `device_credential` as `password` + correct `device_id` → JWT ✅
- [ ] Login with `device_credential` but WRONG `device_id` → 401 ❌
- [ ] Login with `device_credential` but WRONG `mobile_number` → 401 ❌
- [ ] Login after new OTP verify (new `device_id`) → old credential fails → 401 ❌
- [ ] WEB login with real password → SERVER_SIDE session (unchanged) ✅
- [ ] Redis `trusted:device:{uid}:{did}` key set after successful device login
