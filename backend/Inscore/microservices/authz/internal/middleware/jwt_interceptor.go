package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// contextKey is an unexported type for context keys in this package.
type contextKey string

const (
	// ClaimsKey is the context key under which *AuthClaims is stored.
	ClaimsKey contextKey = "jwt_claims"
)

// AuthClaims holds the extracted JWT claims injected into the request context.
type AuthClaims struct {
	UserID   string
	PortalID string
	Roles    []string
	Email    string
}

// JWTInterceptor validates Bearer JWT tokens on incoming gRPC requests.
// Methods listed in skipMethods are passed through without authentication.
type JWTInterceptor struct {
	publicKey   *rsa.PublicKey
	skipMethods map[string]bool
}

// NewJWTInterceptor constructs a JWTInterceptor.
// publicKey may be nil — in that case every request is passed through (no-op mode).
// skipMethods is a list of full gRPC method paths that bypass auth (e.g. health check).
func NewJWTInterceptor(publicKey *rsa.PublicKey, skipMethods []string) *JWTInterceptor {
	skips := make(map[string]bool, len(skipMethods))
	for _, m := range skipMethods {
		skips[m] = true
	}
	return &JWTInterceptor{
		publicKey:   publicKey,
		skipMethods: skips,
	}
}

// UnaryInterceptor returns a gRPC unary server interceptor that validates JWTs.
func (i *JWTInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return i.unaryIntercept
}

// StreamInterceptor returns a gRPC stream server interceptor that validates JWTs.
func (i *JWTInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return i.streamIntercept
}

func (i *JWTInterceptor) unaryIntercept(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if i.publicKey == nil || i.skipMethods[info.FullMethod] {
		return handler(ctx, req)
	}
	claims, err := i.extractClaims(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	ctx = context.WithValue(ctx, ClaimsKey, claims)
	return handler(ctx, req)
}

func (i *JWTInterceptor) streamIntercept(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	if i.publicKey == nil || i.skipMethods[info.FullMethod] {
		return handler(srv, ss)
	}
	claims, err := i.extractClaims(ss.Context())
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	ctx := context.WithValue(ss.Context(), ClaimsKey, claims)
	wrapped := &wrappedStream{ServerStream: ss, ctx: ctx}
	return handler(srv, wrapped)
}

// extractClaims parses and validates the Bearer JWT from gRPC metadata.
func (i *JWTInterceptor) extractClaims(ctx context.Context) (*AuthClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	vals := md.Get("authorization")
	if len(vals) == 0 {
		return nil, errors.New("missing authorization header")
	}

	raw := vals[0]
	if !strings.HasPrefix(raw, "Bearer ") {
		return nil, errors.New("authorization header must start with 'Bearer '")
	}
	tokenStr := strings.TrimPrefix(raw, "Bearer ")
	if tokenStr == "" {
		return nil, errors.New("empty bearer token")
	}

	mapClaims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, mapClaims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return i.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token validation failed: " + err.Error())
	}

	claims := &AuthClaims{}

	// sub → UserID
	if sub, err := mapClaims.GetSubject(); err == nil {
		claims.UserID = sub
	}

	// portal → PortalID
	if portal, ok := mapClaims["portal"]; ok {
		claims.PortalID, _ = portal.(string)
	}

	// email → Email
	if email, ok := mapClaims["email"]; ok {
		claims.Email, _ = email.(string)
	}

	// roles → []string (may be []interface{} in JWT)
	if rolesRaw, ok := mapClaims["roles"]; ok {
		switch v := rolesRaw.(type) {
		case []interface{}:
			for _, r := range v {
				if s, ok := r.(string); ok {
					claims.Roles = append(claims.Roles, s)
				}
			}
		case []string:
			claims.Roles = v
		case string:
			if v != "" {
				claims.Roles = []string{v}
			}
		}
	}

	return claims, nil
}

// GetClaims retrieves *AuthClaims from the context. Returns nil if not present.
func GetClaims(ctx context.Context) *AuthClaims {
	v := ctx.Value(ClaimsKey)
	if v == nil {
		return nil
	}
	c, _ := v.(*AuthClaims)
	return c
}

// ParseRSAPublicKeyFromPEM parses an RSA public key from a PEM-encoded string.
// Returns nil, nil when pemStr is empty (no-op / key not configured).
func ParseRSAPublicKeyFromPEM(pemStr string) (*rsa.PublicKey, error) {
	normalized := normalizePEM(pemStr)
	if normalized == "" {
		return nil, nil
	}

	block, _ := pem.Decode([]byte(normalized))
	// Fallback: caller may have provided a file path instead of PEM content.
	if block == nil && !strings.Contains(normalized, "BEGIN PUBLIC KEY") {
		if data, err := os.ReadFile(normalized); err == nil {
			normalized = normalizePEM(string(data))
			block, _ = pem.Decode([]byte(normalized))
		}
	}

	if block == nil {
		return nil, errors.New("failed to decode PEM block from public key string")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse PKIX public key: " + err.Error())
	}
	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not an RSA key")
	}
	return rsaKey, nil
}

func normalizePEM(value string) string {
	s := strings.TrimSpace(value)
	if s == "" {
		return ""
	}

	// Drop wrapping quotes from .env values if present.
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			s = s[1 : len(s)-1]
		}
	}

	// Support escaped PEM payloads from env (e.g. "-----BEGIN...\\n...").
	s = strings.ReplaceAll(s, `\r`, "\r")
	s = strings.ReplaceAll(s, `\n`, "\n")

	return strings.TrimSpace(s)
}

// wrappedStream wraps grpc.ServerStream with a replacement context.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context { return w.ctx }
