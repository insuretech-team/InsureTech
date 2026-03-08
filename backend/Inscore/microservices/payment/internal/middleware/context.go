// Package middleware provides gRPC middleware utilities for the payment microservice.
package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

// RequestContext holds the caller identity extracted from incoming gRPC metadata.
// The gateway (auth_middleware.go) sets these headers after validating the JWT or
// server-side session. Services must NOT re-validate tokens — trust these values.
type RequestContext struct {
	UserID         string // x-user-id — UUID of the authenticated user
	TenantID       string // x-tenant-id — multi-tenancy isolation key
	Portal         string // x-portal normalized: "system","b2b","b2c","agent","business","regulator"
	SessionID      string // x-session-id — UUID of the current session
	SessionType    string // x-session-type: "SERVER_SIDE" (web portal) or "JWT" (mobile/API)
	TokenID        string // x-token-id — JWT JTI (for revocation checks)
	DeviceID       string // x-device-id — device fingerprint (JWT only)
	UserType       string // x-user-type: "B2C_CUSTOMER","SYSTEM_USER","AGENT","PARTNER", etc.
	OrganisationID string // x-business-id — B2B org UUID (injected by B2BContextMiddleware)
	TraceID        string // x-request-id — distributed tracing ID from gateway RequestID middleware
}

// ExtractRequestContext reads caller identity from incoming gRPC metadata.
// Returns an empty RequestContext if metadata is missing (e.g. in tests without context setup).
func ExtractRequestContext(ctx context.Context) RequestContext {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return RequestContext{}
	}
	return RequestContext{
		UserID:         firstMD(md, "x-user-id"),
		TenantID:       firstMD(md, "x-tenant-id"),
		Portal:         normPortal(firstMD(md, "x-portal")),
		SessionID:      firstMD(md, "x-session-id"),
		SessionType:    firstMD(md, "x-session-type"),
		TokenID:        firstMD(md, "x-token-id"),
		DeviceID:       firstMD(md, "x-device-id"),
		UserType:       firstMD(md, "x-user-type"),
		OrganisationID: firstMD(md, "x-business-id"),
		TraceID:        firstMD(md, "x-request-id"),
	}
}

// ActorUserID returns the portal user who triggered the action.
func (r RequestContext) ActorUserID() string {
	return r.UserID
}

// IsSystemPortal returns true when the request comes from the system admin portal.
func (r RequestContext) IsSystemPortal() bool {
	return r.Portal == "system"
}

// IsB2B returns true when the request carries a B2B organisation context.
func (r RequestContext) IsB2B() bool {
	return r.Portal == "b2b" && r.OrganisationID != ""
}

// IsMobileOrAPI returns true when the session type is JWT (mobile app or API client).
func (r RequestContext) IsMobileOrAPI() bool {
	return strings.EqualFold(r.SessionType, "JWT")
}

// IsWebPortal returns true when the session type is server-side (web portal).
func (r RequestContext) IsWebPortal() bool {
	return strings.EqualFold(r.SessionType, "SERVER_SIDE")
}

// normPortal strips the "PORTAL_" prefix and lowercases the result.
func normPortal(raw string) string {
	return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(raw), "PORTAL_"))
}
