// Package middleware provides gRPC middleware utilities for the payment microservice.
package middleware

import (
	"context"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
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

// normPortal strips the "PORTAL_" prefix and lowercases the result.
func normPortal(raw string) string {
	return grpcmeta.NormalizePortal(raw)
}
