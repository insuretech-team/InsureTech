package service

// totp_service.go — Standalone TOTPService for TOTP MFA (Sprint 1.8).
//
// Architecture
// ─────────────
// TOTP state lives on the User entity (totp_secret_enc, totp_enabled columns).
// The full business logic is implemented in stub_service.go as AuthService
// methods (EnableTOTP, VerifyTOTP, DisableTOTP) which are already wired into
// the gRPC handler via service_iface.go.
//
// This file adds:
//  1. TOTPService — a lean, testable struct that exposes pure helper functions
//     independent of the database (useful for unit tests and future extraction).
//  2. EnrollTOTP / ConfirmTOTP delegation shims on AuthService so that callers
//     that use the Sprint-1.8 vocabulary ("Enroll" / "Confirm") are also
//     supported without duplicating logic.
//
// Encryption helpers (aesGCMEncrypt / aesGCMDecrypt / totpEncryptionKey) live
// in stub_service.go and are reused here.

import (
	"context"
	"errors"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/pquerna/otp/totp"
)

// ── Standalone TOTPService ────────────────────────────────────────────────────

// TOTPService is a stateless helper for TOTP key generation and validation.
// It does not hold any repository reference; all DB interactions are handled by
// AuthService which delegates to this service for the cryptographic operations.
type TOTPService struct{}

// NewTOTPService returns a new TOTPService instance.
func NewTOTPService() *TOTPService { return &TOTPService{} }

// GenerateKey creates a new TOTP key for the given issuer and account name.
// SecretSize is 32 bytes (256 bits), period 30 s, 6 digits — RFC 6238 defaults.
// Returns the otpauth:// provisioning URI and the raw base32 secret.
func (t *TOTPService) GenerateKey(issuer, accountName string) (provisioningURI, secret string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		SecretSize:  32,
		Period:      30,
		Digits:      6,
	})
	if err != nil {
		logger.Errorf("generate TOTP key: %v", err)
		return "", "", errors.New("generate TOTP key")
	}
	return key.URL(), key.Secret(), nil
}

// Validate checks a 6-digit TOTP code against a plaintext base32 secret.
// Allows ±1 step (±30 s) clock skew tolerance.
func (t *TOTPService) Validate(code, secret string) (bool, error) {
	valid, err := totp.ValidateCustom(code, secret, time.Now().UTC(), totp.ValidateOpts{
		Period: 30,
		Skew:   1,
		Digits: 6,
	})
	if err != nil {
		logger.Errorf("totp validate: %v", err)
		return false, errors.New("totp validate")
	}
	return valid, nil
}

// ── Package-level helpers (used by tests without a TOTPService instance) ─────

// EnrollTOTPForUser generates a TOTP key for the given issuer and account.
// Returns the otpauth:// provisioning URI and raw base32 secret.
func EnrollTOTPForUser(issuer, accountName string) (provisioningURI, secret string, err error) {
	return NewTOTPService().GenerateKey(issuer, accountName)
}

// ValidateTOTPCode validates a TOTP code against a plaintext base32 secret
// with ±1 step (30 s) tolerance.
func ValidateTOTPCode(code, secret string) (bool, error) {
	return NewTOTPService().Validate(code, secret)
}

// ── AuthService delegation shims ─────────────────────────────────────────────
// These allow callers using the Sprint-1.8 "Enroll/Confirm" vocabulary to
// reach the existing EnableTOTP / VerifyTOTP implementations without
// duplicating any logic.

// EnrollTOTP generates a new TOTP secret for a user, encrypts and stores it
// (unconfirmed), and returns the provisioning URI + raw base32 secret.
// Delegates to AuthService.EnableTOTP.
func (s *AuthService) EnrollTOTP(ctx context.Context, req *authnservicev1.EnableTOTPRequest) (*authnservicev1.EnableTOTPResponse, error) {
	return s.EnableTOTP(ctx, req)
}

// ConfirmTOTP verifies the first TOTP code after enrolment and activates TOTP
// for the user (sets totp_enabled = true).
// Delegates to AuthService.VerifyTOTP which already handles first-time
// activation (totp_enabled false → true on successful validation).
func (s *AuthService) ConfirmTOTP(ctx context.Context, req *authnservicev1.VerifyTOTPRequest) (*authnservicev1.VerifyTOTPResponse, error) {
	return s.VerifyTOTP(ctx, req)
}
