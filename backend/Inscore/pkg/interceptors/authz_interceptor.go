// Package interceptors provides gRPC interceptors for cross-cutting concerns.
// authz_interceptor.go — per-service AuthZ enforcement interceptor (Sprint G).
//
// Usage in each microservice's gRPC server setup:
//
//	grpc.NewServer(
//	    grpc.ChainUnaryInterceptor(
//	        interceptors.NewAuthZInterceptor(authzClient, interceptors.DefaultSkipMethods),
//	    ),
//	)
package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
)

// DefaultSkipMethods is the set of fully-qualified gRPC method names that bypass AuthZ.
// Add health checks, JWKS, and internal-only methods here.
var DefaultSkipMethods = map[string]bool{
	"/grpc.health.v1.Health/Check":       true,
	"/grpc.health.v1.Health/Watch":       true,
	"/insuretech.authn.services.v1.AuthService/GetJWKS": true,
}

// AuthZInterceptor is a gRPC unary server interceptor that calls AuthZ.CheckAccess
// for every incoming RPC not in the skip list.
// It reads x-user-id, x-portal, x-tenant-id from incoming gRPC metadata (set by the gateway).
type AuthZInterceptor struct {
	client      authzservicev1.AuthZServiceClient
	skipMethods map[string]bool
}

// NewAuthZInterceptor creates a new AuthZInterceptor.
func NewAuthZInterceptor(client authzservicev1.AuthZServiceClient, skipMethods map[string]bool) grpc.UnaryServerInterceptor {
	i := &AuthZInterceptor{
		client:      client,
		skipMethods: skipMethods,
	}
	return i.Intercept
}

// Intercept performs the AuthZ check and calls the handler if allowed.
func (i *AuthZInterceptor) Intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Skip methods in the allow list
	if i.skipMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	// Extract metadata from incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	userID := firstMD(md, "x-user-id")
	portal := firstMD(md, "x-portal")
	tenantID := firstMD(md, "x-tenant-id")
	sessionID := firstMD(md, "x-session-id")
	ipAddr := firstMD(md, "x-forwarded-for")
	userAgent := firstMD(md, "x-user-agent")

	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "x-user-id header required")
	}

	// Build Casbin tuple from gRPC method
	object, action := methodToObjectAction(info.FullMethod)
	domain := portal
	if tenantID != "" {
		domain = portal + ":" + tenantID
	}

	// Call AuthZ.CheckAccess
	resp, err := i.client.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID,
		Domain: domain,
		Object: object,
		Action: action,
		Context: &authzservicev1.AccessContext{
			SessionId: sessionID,
			IpAddress: ipAddr,
			UserAgent: userAgent,
		},
	})
	if err != nil {
		// AuthZ service unavailable — fail open or closed depending on policy
		// Default: fail CLOSED (deny) for security
		return nil, status.Errorf(codes.Unavailable, "authz service unavailable: %v", err)
	}

	if !resp.Allowed {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: %s", resp.Reason)
	}

	return handler(ctx, req)
}

// methodToObjectAction maps a gRPC full method to (svc:service/resource, ACTION).
// /insuretech.policy.services.v1.PolicyService/CreatePolicy → (svc:policy/create, POST)
func methodToObjectAction(fullMethod string) (object, action string) {
	// fullMethod: /insuretech.{service}.services.v1.{Service}Service/{RPC}
	parts := strings.Split(strings.TrimPrefix(fullMethod, "/"), "/")
	if len(parts) != 2 {
		return "svc:unknown", "*"
	}

	// Extract service name from package: "insuretech.policy.services.v1.PolicyService" → "policy"
	pkg := parts[0] // e.g. "insuretech.policy.services.v1.PolicyService"
	rpc := parts[1] // e.g. "CreatePolicy"

	pkgParts := strings.Split(pkg, ".")
	svcName := ""
	if len(pkgParts) >= 2 {
		svcName = pkgParts[1] // e.g. "policy"
	}

	// Map RPC name to resource + action
	rpcLower := strings.ToLower(rpc)
	resource, httpAction := rpcToResourceAction(rpcLower)

	return "svc:" + svcName + "/" + resource, httpAction
}

// rpcToResourceAction maps an RPC name pattern to (resource, HTTP verb).
func rpcToResourceAction(rpcLower string) (resource, action string) {
	switch {
	case strings.HasPrefix(rpcLower, "create") || strings.HasPrefix(rpcLower, "register") || strings.HasPrefix(rpcLower, "submit") || strings.HasPrefix(rpcLower, "initiate"):
		return strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(rpcLower, "create"), "register"), "submit"), "initiate"), "POST"
	case strings.HasPrefix(rpcLower, "get") || strings.HasPrefix(rpcLower, "list") || strings.HasPrefix(rpcLower, "fetch"):
		return strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(rpcLower, "get"), "list"), "fetch"), "GET"
	case strings.HasPrefix(rpcLower, "update") || strings.HasPrefix(rpcLower, "patch"):
		return strings.TrimPrefix(strings.TrimPrefix(rpcLower, "update"), "patch"), "PUT"
	case strings.HasPrefix(rpcLower, "delete") || strings.HasPrefix(rpcLower, "remove") || strings.HasPrefix(rpcLower, "revoke"):
		return strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(rpcLower, "delete"), "remove"), "revoke"), "DELETE"
	case strings.HasPrefix(rpcLower, "approve"):
		return strings.TrimPrefix(rpcLower, "approve"), "approve"
	case strings.HasPrefix(rpcLower, "reject"):
		return strings.TrimPrefix(rpcLower, "reject"), "reject"
	default:
		return rpcLower, "*"
	}
}

// firstMD returns the first value for a metadata key, or empty string.
func firstMD(md metadata.MD, key string) string {
	vals := md.Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}
