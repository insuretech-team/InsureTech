package grpc

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	case contains(lower,
		"not found", "no rows", "record not found", "does not exist",
	):
		return status.Error(codes.NotFound, msg)

	// ── AlreadyExists ───────────────────────────────────────────────────────
	case contains(lower,
		"already exists", "duplicate", "already registered",
	):
		return status.Error(codes.AlreadyExists, msg)

	// ── Unauthenticated ─────────────────────────────────────────────────────
	case contains(lower,
		"invalid token", "token expired", "unauthorized",
		"missing authorization", "empty bearer",
		"token validation failed", "unexpected signing method",
	):
		return status.Error(codes.Unauthenticated, msg)

	// ── PermissionDenied ────────────────────────────────────────────────────
	case contains(lower,
		"forbidden", "permission denied", "not allowed", "access denied",
		"does not belong",
	):
		return status.Error(codes.PermissionDenied, msg)

	// ── InvalidArgument ─────────────────────────────────────────────────────
	case contains(lower,
		"invalid", "required", "must be", "malformed", "unsupported",
		"nil role", "nil policy", "invalid payload",
	):
		return status.Error(codes.InvalidArgument, msg)

	// ── FailedPrecondition ──────────────────────────────────────────────────
	case contains(lower,
		"not configured", "not enabled", "not available",
		"invalid pem block", "public key is not rsa",
	):
		return status.Error(codes.FailedPrecondition, msg)

	// ── Default: Internal ───────────────────────────────────────────────────
	default:
		return status.Error(codes.Internal, msg)
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
