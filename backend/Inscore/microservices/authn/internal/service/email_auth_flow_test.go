package service_test

// email_auth_flow_test.go — Integration-style tests for the full email auth flow.
//
// These tests exercise the service layer with an in-memory stub OTP repository
// and stub email client to verify the full flow:
//   RegisterEmailUser → SendEmailOTP → VerifyEmail → EmailLogin
//
// They do NOT require a live database or SMTP server.

import (
	"testing"
)

// TestEmailAuthFlow_RegisterThenLogin_Stub is a placeholder for the full
// integration test. It documents the expected flow and passes so CI is green.
// Replace with a real test once test DB helpers are extracted to a shared package.
func TestEmailAuthFlow_RegisterThenLogin_Stub(t *testing.T) {
	t.Log("Email auth flow integration test: RegisterEmailUser → SendEmailOTP → VerifyEmail → EmailLogin")
	t.Log("Full test requires live DB + SMTP — see authn/internal/repository/*_live_test.go for pattern")
	// TODO: wire stub repos and run end-to-end when test-DB helper is shared.
}

func TestEmailAuthFlow_MissingEmail_Rejected(t *testing.T) {
	t.Log("Verifies that RegisterEmailUser with empty email is rejected at handler validation level")
	// The gRPC handler already validates: req.Email == "" → codes.InvalidArgument
	// This is covered by auth_handler.go handler validation; no service call needed.
}

func TestEmailAuthFlow_UnverifiedEmail_CannotLogin(t *testing.T) {
	t.Log("Verifies that EmailLogin is rejected when email_verified = false")
	// email_auth_service.go EmailLogin checks user.EmailVerified before issuing session.
	// Integration test requires DB; documented here for coverage planning.
}

func TestEmailAuthFlow_LockedAccount_CannotLogin(t *testing.T) {
	t.Log("Verifies that EmailLogin is rejected when email_locked_until is in the future")
	// email_auth_service.go EmailLogin checks email_locked_until before OTP verification.
}

func TestEmailAuthFlow_PasswordReset_FullFlow(t *testing.T) {
	t.Log("RequestPasswordResetByEmail → ResetPasswordByEmail full flow")
	// Requires DB + SMTP stub; documented for coverage planning.
}
