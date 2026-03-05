package grpc

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapError is the exported alias of toGRPCError, available for use in tests and
// external packages that need to inspect gRPC error code mapping behaviour.
var MapError = toGRPCError

// toGRPCError converts a service-layer error into a properly coded gRPC status error.
// It inspects the error message string to map well-known domain errors to the correct
// gRPC status codes, falling back to codes.Internal for unexpected errors.
func toGRPCError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	lower := strings.ToLower(msg)

	switch {
	// Already a gRPC status error — pass through unchanged.
	case isGRPCStatus(err):
		return err

	// ── NotFound ────────────────────────────────────────────────────────────
	case contains(lower, "not found", "no rows", "record not found", "does not exist"):
		return status.Error(codes.NotFound, msg)

	// ── AlreadyExists ───────────────────────────────────────────────────────
	case contains(lower, "already exists", "duplicate", "already registered", "already verified"):
		return status.Error(codes.AlreadyExists, msg)

	// ── Unauthenticated ─────────────────────────────────────────────────────
	case contains(lower,
		"invalid credentials", "invalid password", "invalid token",
		"token expired", "session expired", "session not found",
		"invalid otp", "otp expired", "otp not found",
		"invalid csrf", "csrf mismatch",
		"invalid api key", "api key not found", "api key revoked",
		"biometric token mismatch", "biometric not enrolled",
		"unauthorized",
	):
		return status.Error(codes.Unauthenticated, msg)

	// ── PermissionDenied ────────────────────────────────────────────────────
	case contains(lower, "forbidden", "permission denied", "not allowed", "access denied"):
		return status.Error(codes.PermissionDenied, msg)

	// ── ResourceExhausted (rate limiting) ───────────────────────────────────
	case contains(lower, "rate limit", "too many", "too many requests", "quota exceeded"):
		return status.Error(codes.ResourceExhausted, msg)

	// ── InvalidArgument ─────────────────────────────────────────────────────
	case contains(lower,
		"invalid", "required", "must be", "malformed",
		"password too short", "weak password", "invalid email",
		"invalid phone", "invalid mobile",
	):
		return status.Error(codes.InvalidArgument, msg)

	// ── FailedPrecondition ──────────────────────────────────────────────────
	case contains(lower,
		"not verified", "email not verified", "phone not verified",
		"account disabled", "account locked", "account suspended",
	):
		return status.Error(codes.FailedPrecondition, msg)

	// ── Unavailable (downstream/SMS/email failures) ──────────────────────────
	case contains(lower, "sms failed", "email failed", "provider error", "send failed"):
		return status.Error(codes.Unavailable, msg)

	// ── Default: Internal ───────────────────────────────────────────────────
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

// contains returns true if s contains any of the given substrings.
func contains(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// isGRPCStatus returns true if err is already a gRPC status error.
func isGRPCStatus(err error) bool {
	_, ok := status.FromError(err)
	return ok && status.Code(err) != codes.OK && status.Code(err) != codes.Unknown
}
