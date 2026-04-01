// Package middleware provides gRPC middleware utilities for the orders microservice.
package middleware

import (
	"context"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	"google.golang.org/grpc/metadata"
)

// RequestContext holds the caller identity extracted from incoming gRPC metadata.
// The gateway (auth_middleware.go) sets these headers after validating the JWT or
// server-side session. Services must NOT re-validate tokens — trust these values.
type RequestContext = grpcmeta.RequestContext

// ExtractRequestContext reads caller identity from incoming gRPC metadata.
// Returns an empty RequestContext if metadata is missing (e.g. in tests without context setup).
func ExtractRequestContext(ctx context.Context) RequestContext {
	return grpcmeta.ExtractRequestContext(ctx)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

// normPortal strips the "PORTAL_" prefix and lowercases the result.
// "PORTAL_B2B" → "b2b", "PORTAL_SYSTEM" → "system", "b2c" → "b2c"
func normPortal(raw string) string {
	return grpcmeta.NormalizePortal(raw)
}

// firstMD returns the first non-empty value for a metadata key, or "".
func firstMD(md metadata.MD, key string) string {
	return grpcmeta.First(md, key)
}
