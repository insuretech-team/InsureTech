package grpcmeta

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// RequestContext holds caller identity propagated through incoming gRPC metadata.
type RequestContext struct {
	UserID         string
	TenantID       string
	Portal         string
	SessionID      string
	SessionType    string
	TokenID        string
	DeviceID       string
	UserType       string
	OrganisationID string
	TraceID        string
}

// RequestMetadata holds transport-level request metadata useful for auth/audit flows.
type RequestMetadata struct {
	IPAddress     string
	UserAgent     string
	DeviceID      string
	SessionToken  string
	CSRFToken     string
	Authorization string
}

var (
	DefaultActorKeys       = []string{"x-user-id", "x-actor-id", "x-sub", "user-id"}
	DefaultTenantKeys      = []string{"x-tenant-id", "tenant-id", "x-tenant"}
	DefaultCorrelationKeys = []string{"x-correlation-id", "x-request-id", "x-trace-id", "request-id"}
	DefaultForwardKeys     = []string{"x-user-id", "x-tenant-id", "x-portal", "x-user-type", "x-session-id", "x-session-type", "x-token-id", "x-device-id", "x-business-id", "x-request-id", "x-correlation-id", "x-forwarded-for", "x-real-ip", "user-agent", "authorization", "cookie"}
)

// ExtractRequestContext reads common caller identity fields from incoming gRPC metadata.
func ExtractRequestContext(ctx context.Context) RequestContext {
	return RequestContext{
		UserID:         FirstFromContext(ctx, "x-user-id"),
		TenantID:       FirstFromContext(ctx, "x-tenant-id"),
		Portal:         NormalizePortal(FirstFromContext(ctx, "x-portal")),
		SessionID:      FirstFromContext(ctx, "x-session-id"),
		SessionType:    FirstFromContext(ctx, "x-session-type"),
		TokenID:        FirstFromContext(ctx, "x-token-id"),
		DeviceID:       FirstFromContext(ctx, "x-device-id"),
		UserType:       FirstFromContext(ctx, "x-user-type"),
		OrganisationID: FirstFromContext(ctx, "x-business-id"),
		TraceID:        FirstFromContext(ctx, "x-request-id"),
	}
}

func (r RequestContext) ActorUserID() string {
	return r.UserID
}

func (r RequestContext) IsSystemPortal() bool {
	return r.Portal == "system"
}

func (r RequestContext) IsB2B() bool {
	return r.Portal == "b2b" && r.OrganisationID != ""
}

func (r RequestContext) IsMobileOrAPI() bool {
	return strings.EqualFold(r.SessionType, "JWT")
}

func (r RequestContext) IsWebPortal() bool {
	return strings.EqualFold(r.SessionType, "SERVER_SIDE")
}

// First returns the first non-empty metadata value for a key.
func First(md metadata.MD, key string) string {
	for _, value := range md.Get(key) {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

// FirstFromContext returns the first non-empty incoming gRPC metadata value for a key.
func FirstFromContext(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return First(md, key)
}

// FirstOf returns the first non-empty value across multiple metadata keys.
func FirstOf(md metadata.MD, keys ...string) string {
	for _, key := range keys {
		if value := First(md, key); value != "" {
			return value
		}
	}
	return ""
}

// FirstOfFromContext returns the first non-empty incoming metadata value across multiple keys.
func FirstOfFromContext(ctx context.Context, keys ...string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return FirstOf(md, keys...)
}

func WithOutgoingMetadata(ctx context.Context, pairs ...string) context.Context {
	if len(pairs) == 0 {
		return ctx
	}
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		cloned := md.Copy()
		for i := 0; i+1 < len(pairs); i += 2 {
			key := strings.TrimSpace(pairs[i])
			value := strings.TrimSpace(pairs[i+1])
			if key == "" || value == "" {
				continue
			}
			cloned.Set(key, value)
		}
		return metadata.NewOutgoingContext(ctx, cloned)
	}
	filtered := make([]string, 0, len(pairs))
	for i := 0; i+1 < len(pairs); i += 2 {
		key := strings.TrimSpace(pairs[i])
		value := strings.TrimSpace(pairs[i+1])
		if key == "" || value == "" {
			continue
		}
		filtered = append(filtered, key, value)
	}
	if len(filtered) == 0 {
		return ctx
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(filtered...))
}

func ForwardIncomingToOutgoing(ctx context.Context, keys ...string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	if len(keys) == 0 {
		keys = DefaultForwardKeys
	}
	pairs := make([]string, 0, len(keys)*2)
	for _, key := range keys {
		if value := First(md, key); value != "" {
			pairs = append(pairs, key, value)
		}
	}
	return WithOutgoingMetadata(ctx, pairs...)
}

// NormalizePortal strips the PORTAL_ prefix and lowercases the remaining portal name.
func NormalizePortal(raw string) string {
	return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(raw), "PORTAL_"))
}

// ActorID returns the canonical acting user id from incoming metadata using shared fallback aliases.
func ActorID(ctx context.Context, fallback string) string {
	if actorID := FirstOfFromContext(ctx, DefaultActorKeys...); actorID != "" {
		return actorID
	}
	return strings.TrimSpace(fallback)
}

// TenantID returns the canonical tenant id from incoming metadata using shared fallback aliases.
func TenantID(ctx context.Context, fallback string) string {
	if tenantID := FirstOfFromContext(ctx, DefaultTenantKeys...); tenantID != "" {
		return tenantID
	}
	return strings.TrimSpace(fallback)
}

// CorrelationID returns the first request correlation identifier from incoming metadata.
func CorrelationID(ctx context.Context) string {
	return FirstOfFromContext(ctx, DefaultCorrelationKeys...)
}

// ExtractIPAddress returns the best-effort client IP from metadata or peer info.
func ExtractIPAddress(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := First(md, "x-forwarded-for"); xff != "" {
			parts := strings.Split(xff, ",")
			if len(parts) > 0 {
				return normalizeIPAddress(strings.TrimSpace(parts[0]))
			}
		}
		if xri := First(md, "x-real-ip"); xri != "" {
			return normalizeIPAddress(xri)
		}
	}

	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		return normalizeIPAddress(p.Addr.String())
	}
	return "unknown"
}

func normalizeIPAddress(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "unknown"
	}

	if ip := net.ParseIP(value); ip != nil {
		return ip.String()
	}

	if host, _, err := net.SplitHostPort(value); err == nil {
		host = strings.TrimSpace(host)
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
		if host != "" {
			return host
		}
	}

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		host := strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
		if host != "" {
			return host
		}
	}

	return value
}

// ExtractUserAgent returns the caller user agent from gRPC metadata.
func ExtractUserAgent(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := First(md, "user-agent"); ua != "" {
			return ua
		}
		if grpcUA := First(md, "grpc-user-agent"); grpcUA != "" {
			return grpcUA
		}
	}
	return "unknown"
}

func ExtractDeviceID(ctx context.Context) string {
	return FirstFromContext(ctx, "x-device-id")
}

func ExtractSessionToken(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return ParseCookie(First(md, "cookie"), "session_token")
	}
	return ""
}

func ExtractCSRFToken(ctx context.Context) string {
	return FirstFromContext(ctx, "x-csrf-token")
}

func ExtractAuthorizationToken(ctx context.Context) string {
	auth := FirstFromContext(ctx, "authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func ExtractAll(ctx context.Context) *RequestMetadata {
	return &RequestMetadata{
		IPAddress:     ExtractIPAddress(ctx),
		UserAgent:     ExtractUserAgent(ctx),
		DeviceID:      ExtractDeviceID(ctx),
		SessionToken:  ExtractSessionToken(ctx),
		CSRFToken:     ExtractCSRFToken(ctx),
		Authorization: ExtractAuthorizationToken(ctx),
	}
}

// ParseCookie returns a named cookie value from a raw Cookie header.
func ParseCookie(cookieStr, name string) string {
	if cookieStr == "" || name == "" {
		return ""
	}
	for _, cookie := range strings.Split(cookieStr, ";") {
		parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
		if len(parts) == 2 && parts[0] == name {
			return parts[1]
		}
	}
	return ""
}
