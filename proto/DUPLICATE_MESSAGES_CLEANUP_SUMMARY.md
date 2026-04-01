# Authn Service Proto Files - Duplicate Messages Cleanup

## Summary
Successfully removed duplicate message definitions and RPC methods from the authn service proto files that belonged to other services.

## Changes Made

### File: `insuretech/authn/services/v1/core.proto`

#### Message Definitions Removed:
1. **JWK** (lines 1012-1019) - Belongs to authz service
2. **GetJWKSRequest** (line 1020) - Belongs to authz service
3. **GetJWKSResponse** (lines 1021-1023) - Belongs to authz service
4. **GetVoiceSessionRequest** (lines 882-884) - Belongs to voice service
5. **GetVoiceSessionResponse** (lines 885-894) - Belongs to voice service
6. **EndVoiceSessionRequest** (lines 895-898) - Belongs to voice service
7. **EndVoiceSessionResponse** (lines 900-903) - Belongs to voice service
8. **RejectKYCRequest** (lines 844-848) - Belongs to kyc service
9. **RejectKYCResponse** (lines 849-852) - Belongs to kyc service

**Note:** Authn's own `ApproveKYC` (lines 836-843) remains as it is authn-specific.

### File: `insuretech/authn/services/v1/auth_service.proto`

#### RPC Methods Removed:
1. **GetJWKS** (lines 457-461) - Maps to authz service's GetJWKS RPC
2. **GetVoiceSession** (lines 381-386) - Maps to voice service's GetVoiceSession RPC
3. **EndVoiceSession** (lines 387-393) - Maps to voice service's EndVoiceSession RPC
4. **RejectKYC** (lines 358-363) - Maps to kyc service's RejectKYC RPC

**Note:** Authn's own `CreateVoiceSession` RPC remains as it is used by authn service.

## Remaining Messages in Authn Service

The authn service now retains only its own message definitions:
- All authentication/OTP messages (Register, Login, VerifyOTP, RefreshToken, etc.)
- Email authentication messages
- Biometric authentication messages
- Session management messages
- User profile and document management messages
- API key management messages
- Document type messages
- KYC-specific messages (InitiateKYC, GetKYCStatus, SubmitKYCFrame, CompleteKYCSession, ApproveKYC)
- Document verification messages
- Voice session creation messages (CreateVoiceSession)
- Profile photo upload URL messages
- Notification preference messages
- TOTP/2FA messages
- Voice biometric auth messages (InitiateVoiceSession, SubmitVoiceSample, VerifyVoiceSession)

## Import Statements

**No new imports were added** because:
- The authn service now has NO RPCs that reference types from authz, voice, or kyc services
- It only uses its own message types defined in core.proto
- All imported types (entity definitions, common error) remain unchanged

## Files Modified
1. ✅ `E:\Projects\InsureTech\proto\insuretech\authn\services\v1\core.proto` - 9 message definitions removed
2. ✅ `E:\Projects\InsureTech\proto\insuretech\authn\services\v1\auth_service.proto` - 4 RPC methods removed

## Files Unchanged
- `E:\Projects\InsureTech\proto\insuretech\authz\services\v1\authz_service.proto` - Contains authoritative GetJWKS
- `E:\Projects\InsureTech\proto\insuretech\authz\services\v1\core.proto` - Contains JWK/GetJWKSRequest/GetJWKSResponse with error field
- `E:\Projects\InsureTech\proto\insuretech\voice\services\v1\voice_service.proto` - Contains GetVoiceSession and EndVoiceSession
- `E:\Projects\InsureTech\proto\insuretech\kyc\services\v1\kyc_service.proto` - Contains RejectKYC

## Status
✅ **COMPLETE** - All duplicates removed, files cleaned up successfully.
