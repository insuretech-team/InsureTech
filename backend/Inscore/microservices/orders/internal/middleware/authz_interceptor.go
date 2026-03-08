package middleware

import (
	"context"
	"fmt"
	"strings"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthZClient is the interface for calling the AuthZ service.
// Matches the generated gRPC client signature (variadic grpc.CallOption).
type AuthZClient interface {
	CheckAccess(ctx context.Context, req *authzservicev1.CheckAccessRequest, opts ...grpc.CallOption) (*authzservicev1.CheckAccessResponse, error)
}

// OrderAuthZInterceptor enforces Casbin authorization for the orders-service.
// It reads x-user-id, x-portal, x-tenant-id, x-business-id from incoming gRPC metadata
// (set by the gateway after validating the JWT or server-side session) and calls
// AuthZ.CheckAccess as a defense-in-depth layer.
//
// Pattern copied from b2b/internal/middleware/authz_interceptor.go.
type OrderAuthZInterceptor struct {
	authzClient AuthZClient
}

// NewOrderAuthZInterceptor creates a new OrderAuthZInterceptor.
func NewOrderAuthZInterceptor(authzClient AuthZClient) *OrderAuthZInterceptor {
	return &OrderAuthZInterceptor{authzClient: authzClient}
}

// UnaryServerInterceptor returns a gRPC unary interceptor that enforces AuthZ.
func (i *OrderAuthZInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Health checks and reflection always pass through.
		if isSkipMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		userID := firstMD(md, "x-user-id")
		if userID == "" {
			return nil, status.Error(codes.Unauthenticated, "x-user-id required")
		}

		portalRaw := firstMD(md, "x-portal")
		portalNorm := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(portalRaw), "PORTAL_"))
		isSystemPortal := portalNorm == "system"

		orgID := firstMD(md, "x-business-id")

		// Non-system portals without org context: allow bootstrap methods, deny others.
		if !isSystemPortal && orgID == "" {
			if isOrderBootstrapMethod(info.FullMethod) {
				return handler(ctx, req)
			}
			// B2C customers and agents do not require x-business-id.
			// Only B2B portal requires it for non-bootstrap methods.
			if portalNorm == "b2b" {
				return nil, status.Error(codes.PermissionDenied, "missing organisation context")
			}
		}

		// Map gRPC method to Casbin object/action.
		resource, action := mapOrderMethodToResourceAction(info.FullMethod)
		if resource == "" {
			// Bootstrap or public method — no Casbin check.
			return handler(ctx, req)
		}

		domain := resolveOrderAuthzDomain(md, orgID)

		resp, err := i.authzClient.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
			UserId: userID,
			Domain: domain,
			Object: resource,
			Action: action,
		})
		if err != nil {
			// Fail open in development; log the error so it surfaces.
			appLogger.Errorf("orders AuthZ check failed (failing open): user=%s method=%s err=%v",
				userID, info.FullMethod, err)
			// TODO: switch to fail-closed once AuthZ service is stable in all envs:
			// return nil, status.Errorf(codes.Unavailable, "authz service unavailable: %v", err)
			return handler(ctx, req)
		}

		if !resp.Allowed {
			appLogger.Warnf("orders access denied: user=%s domain=%s resource=%s action=%s",
				userID, domain, resource, action)
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}

		return handler(ctx, req)
	}
}

// mapOrderMethodToResourceAction maps an orders-service gRPC method to a Casbin (object, action) pair.
// Returns ("", "") to skip the Casbin check entirely for bootstrap/public methods.
func mapOrderMethodToResourceAction(method string) (resource, action string) {
	parts := strings.Split(method, "/")
	if len(parts) < 3 {
		return "", ""
	}
	methodName := parts[len(parts)-1]

	switch {
	// READ operations
	case strings.HasPrefix(methodName, "Get"),
		strings.HasPrefix(methodName, "List"):
		return "svc:order/*", "GET"

	// CREATE operations
	case strings.HasPrefix(methodName, "Create"):
		return "svc:order/*", "POST"

	// Payment initiation — POST semantics
	case methodName == "InitiatePayment":
		return "svc:order/*", "POST"

	// Confirm payment — system/internal callback; PATCH semantics
	case methodName == "ConfirmPayment":
		return "svc:order/*", "PATCH"

	// UPDATE operations
	case strings.HasPrefix(methodName, "Update"):
		return "svc:order/*", "PATCH"

	// CANCEL / DELETE semantics
	case methodName == "CancelOrder",
		strings.HasPrefix(methodName, "Delete"),
		strings.HasPrefix(methodName, "Remove"):
		return "svc:order/*", "DELETE"

	default:
		appLogger.Warnf("orders authz: unmapped method %s — skipping Casbin check", methodName)
		return "", ""
	}
}

// resolveOrderAuthzDomain builds the Casbin domain string from metadata.
func resolveOrderAuthzDomain(md metadata.MD, orgID string) string {
	portal := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(firstMD(md, "x-portal")), "PORTAL_"))
	tenantID := firstMD(md, "x-tenant-id")

	switch portal {
	case "system":
		return "system:root"
	case "b2b":
		if orgID != "" {
			return fmt.Sprintf("b2b:%s", orgID)
		}
		if tenantID != "" {
			return fmt.Sprintf("b2b:%s", tenantID)
		}
		return "b2b:root"
	case "b2c":
		if tenantID != "" {
			return fmt.Sprintf("b2c:%s", tenantID)
		}
		return "b2c:root"
	case "agent":
		if tenantID != "" {
			return fmt.Sprintf("agent:%s", tenantID)
		}
		return "agent:root"
	default:
		if tenantID != "" {
			return fmt.Sprintf("%s:%s", portal, tenantID)
		}
		return "b2c:root"
	}
}

// isSkipMethod returns true for gRPC methods that must always pass through without auth.
func isSkipMethod(method string) bool {
	skipMethods := map[string]bool{
		"/grpc.health.v1.Health/Check": true,
		"/grpc.health.v1.Health/Watch": true,
	}
	return skipMethods[method]
}

// isOrderBootstrapMethod returns true for methods that are allowed without org context.
// Currently none for orders-service — all methods require full identity.
func isOrderBootstrapMethod(_ string) bool {
	return false
}
