package middleware

import (
	"context"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
)

// MetadataExtractor extracts metadata from gRPC context
type MetadataExtractor struct{}

type RequestMetadata = grpcmeta.RequestMetadata

// NewMetadataExtractor creates a new metadata extractor
func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{}
}

// ExtractIPAddress extracts the client IP address from gRPC context
func (m *MetadataExtractor) ExtractIPAddress(ctx context.Context) string {
	return grpcmeta.ExtractIPAddress(ctx)
}

// ExtractUserAgent extracts the user agent from gRPC context
func (m *MetadataExtractor) ExtractUserAgent(ctx context.Context) string {
	return grpcmeta.ExtractUserAgent(ctx)
}

// ExtractDeviceID extracts device ID from metadata (custom header)
func (m *MetadataExtractor) ExtractDeviceID(ctx context.Context) string {
	return grpcmeta.ExtractDeviceID(ctx)
}

// ExtractSessionToken extracts session token from cookie header
func (m *MetadataExtractor) ExtractSessionToken(ctx context.Context) string {
	return grpcmeta.ExtractSessionToken(ctx)
}

// ExtractCSRFToken extracts CSRF token from custom header
func (m *MetadataExtractor) ExtractCSRFToken(ctx context.Context) string {
	return grpcmeta.ExtractCSRFToken(ctx)
}

// ExtractAuthorizationToken extracts Bearer token from Authorization header
func (m *MetadataExtractor) ExtractAuthorizationToken(ctx context.Context) string {
	return grpcmeta.ExtractAuthorizationToken(ctx)
}

// parseCookie parses a cookie string and returns the value for the given name
func parseCookie(cookieStr, name string) string {
	return grpcmeta.ParseCookie(cookieStr, name)
}

// ExtractAll extracts all metadata from context
func (m *MetadataExtractor) ExtractAll(ctx context.Context) *RequestMetadata {
	return grpcmeta.ExtractAll(ctx)
}
