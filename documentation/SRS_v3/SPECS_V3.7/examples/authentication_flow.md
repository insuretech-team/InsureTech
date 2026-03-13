# Authentication Flow Example

## User Registration Flow

```
1. User enters mobile number (+880 1XXX XXXXXX)
2. System validates format
3. System sends OTP via SMS
4. User enters 6-digit OTP code
5. System verifies OTP
6. User creates password (min 8 chars, 1 upper, 1 number, 1 special)
7. System creates user account
8. System issues JWT tokens (access + refresh)
9. User redirected to profile completion
```

## Login Flow

```
1. User enters mobile number + password
2. System validates credentials
3. System checks account status (not locked)
4. System generates JWT access token (15 min) and refresh token (7 days)
5. System creates session record
6. User authenticated - redirect to dashboard
```

## OTP Verification Example

**Request:**
```json
POST /api/v1/auth/otp/verify
{
  "otp_id": "otp_abc123",
  "code": "123456"
}
```

**Response:**
```json
{
  "verified": true,
  "user_id": "user_xyz789",
  "message": "OTP verified successfully"
}
```
