package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// MetadataExtractor extracts metadata from gRPC context
type MetadataExtractor struct{}

// NewMetadataExtractor creates a new metadata extractor
func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{}
}

// ExtractIPAddress extracts the client IP address from gRPC context
func (m *MetadataExtractor) ExtractIPAddress(ctx context.Context) string {
	// Try X-Forwarded-For header first (for proxied requests)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			// X-Forwarded-For can be comma-separated, take the first one
			ips := strings.Split(xff[0], ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}

		// Try X-Real-IP header
		if xri := md.Get("x-real-ip"); len(xri) > 0 {
			return xri[0]
		}
	}

	// Fallback to peer address
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}

	return "unknown"
}

// ExtractUserAgent extracts the user agent from gRPC context
func (m *MetadataExtractor) ExtractUserAgent(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			return ua[0]
		}
		// Also check grpc-user-agent
		if gua := md.Get("grpc-user-agent"); len(gua) > 0 {
			return gua[0]
		}
	}
	return "unknown"
}

// ExtractDeviceID extracts device ID from metadata (custom header)
func (m *MetadataExtractor) ExtractDeviceID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if deviceID := md.Get("x-device-id"); len(deviceID) > 0 {
			return deviceID[0]
		}
	}
	return ""
}

// ExtractSessionToken extracts session token from cookie header
func (m *MetadataExtractor) ExtractSessionToken(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if cookies := md.Get("cookie"); len(cookies) > 0 {
			return parseCookie(cookies[0], "session_token")
		}
	}
	return ""
}

// ExtractCSRFToken extracts CSRF token from custom header
func (m *MetadataExtractor) ExtractCSRFToken(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if csrf := md.Get("x-csrf-token"); len(csrf) > 0 {
			return csrf[0]
		}
	}
	return ""
}

// ExtractAuthorizationToken extracts Bearer token from Authorization header
func (m *MetadataExtractor) ExtractAuthorizationToken(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if auth := md.Get("authorization"); len(auth) > 0 {
			// Format: "Bearer <token>"
			parts := strings.SplitN(auth[0], " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				return parts[1]
			}
		}
	}
	return ""
}

// parseCookie parses a cookie string and returns the value for the given name
func parseCookie(cookieStr, name string) string {
	cookies := strings.Split(cookieStr, ";")
	for _, cookie := range cookies {
		parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
		if len(parts) == 2 && parts[0] == name {
			return parts[1]
		}
	}
	return ""
}

// RequestMetadata holds extracted request metadata
type RequestMetadata struct {
	IPAddress     string
	UserAgent     string
	DeviceID      string
	SessionToken  string
	CSRFToken     string
	Authorization string
}

// ExtractAll extracts all metadata from context
func (m *MetadataExtractor) ExtractAll(ctx context.Context) *RequestMetadata {
	return &RequestMetadata{
		IPAddress:     m.ExtractIPAddress(ctx),
		UserAgent:     m.ExtractUserAgent(ctx),
		DeviceID:      m.ExtractDeviceID(ctx),
		SessionToken:  m.ExtractSessionToken(ctx),
		CSRFToken:     m.ExtractCSRFToken(ctx),
		Authorization: m.ExtractAuthorizationToken(ctx),
	}
}
