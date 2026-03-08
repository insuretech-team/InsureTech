// Package middleware provides gRPC middleware utilities for the payment microservice.
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
type AuthZClient interface {
	CheckAccess(ctx context.Context, req *authzservicev1.CheckAccessRequest, opts ...grpc.CallOption) (*authzservicev1.CheckAccessResponse, error)
}

// PaymentAuthZInterceptor enforces Casbin authorization for the payment-service.
// It reads x-user-id, x-portal, x-tenant-id, x-business-id from incoming gRPC metadata
// (set by the gateway after validating the JWT or server-side session) and calls
// AuthZ.CheckAccess as a defense-in-depth layer.
//
// Pattern copied from b2b/internal/middleware/authz_interceptor.go.
type PaymentAuthZInterceptor struct {
	authzClient AuthZClient
}

// NewPaymentAuthZInterceptor creates a new PaymentAuthZInterceptor.
func NewPaymentAuthZInterceptor(authzClient AuthZClient) *PaymentAuthZInterceptor {
	return &PaymentAuthZInterceptor{authzClient: authzClient}
}

// UnaryServerInterceptor returns a gRPC unary interceptor that enforces AuthZ.
func (i *PaymentAuthZInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
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

		// Webhook callbacks arrive without user identity — skip AuthZ for them.
		if isWebhookMethod(info.FullMethod) {
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

		// B2B portal requires org context for non-webhook methods.
		if !isSystemPortal && portalNorm == "b2b" && orgID == "" {
			return nil, status.Error(codes.PermissionDenied, "missing organisation context")
		}

		// Map gRPC method to Casbin object/action.
		resource, action := mapPaymentMethodToResourceAction(info.FullMethod)
		if resource == "" {
			// Public or bootstrap method — no Casbin check.
			return handler(ctx, req)
		}

		domain := resolvePaymentAuthzDomain(md, orgID)

		resp, err := i.authzClient.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
			UserId: userID,
			Domain: domain,
			Object: resource,
			Action: action,
		})
		if err != nil {
			// Fail open in development — log error so it surfaces in monitoring.
			appLogger.Errorf("payment AuthZ check failed (failing open): user=%s method=%s err=%v",
				userID, info.FullMethod, err)
			// TODO: switch to fail-closed once AuthZ service is stable in all envs:
			// return nil, status.Errorf(codes.Unavailable, "authz service unavailable: %v", err)
			return handler(ctx, req)
		}

		if !resp.Allowed {
			appLogger.Warnf("payment access denied: user=%s domain=%s resource=%s action=%s",
				userID, domain, resource, action)
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}

		return handler(ctx, req)
	}
}

// mapPaymentMethodToResourceAction maps a payment-service gRPC method to a Casbin (object, action) pair.
// Returns ("", "") to skip the Casbin check entirely for public/webhook methods.
func mapPaymentMethodToResourceAction(method string) (resource, action string) {
	parts := strings.Split(method, "/")
	if len(parts) < 3 {
		return "", ""
	}
	methodName := parts[len(parts)-1]

	switch {
	// READ operations
	case strings.HasPrefix(methodName, "Get"),
		strings.HasPrefix(methodName, "List"):
		return "svc:payment/*", "GET"

	// INITIATE / SUBMIT / GENERATE / RECONCILE / HANDLE → POST semantics
	case strings.HasPrefix(methodName, "Initiate"),
		strings.HasPrefix(methodName, "Add"),
		strings.HasPrefix(methodName, "Submit"),
		strings.HasPrefix(methodName, "Generate"),
		strings.HasPrefix(methodName, "Reconcile"),
		strings.HasPrefix(methodName, "Handle"):
		return "svc:payment/*", "POST"

	// VERIFY / REVIEW → PATCH semantics (state mutation by staff/system)
	case strings.HasPrefix(methodName, "Verify"),
		strings.HasPrefix(methodName, "Review"):
		return "svc:payment/*", "PATCH"

	default:
		appLogger.Warnf("payment authz: unmapped method %s — skipping Casbin check", methodName)
		return "", ""
	}
}

// resolvePaymentAuthzDomain builds the Casbin domain string from metadata.
func resolvePaymentAuthzDomain(md metadata.MD, orgID string) string {
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
	case "business":
		if tenantID != "" {
			return fmt.Sprintf("business:%s", tenantID)
		}
		return "business:root"
	default:
		if tenantID != "" {
			return fmt.Sprintf("%s:%s", portal, tenantID)
		}
		return "b2c:root"
	}
}

// isSkipMethod returns true for gRPC methods that must always pass through without auth.
func isSkipMethod(method string) bool {
	return method == "/grpc.health.v1.Health/Check" ||
		method == "/grpc.health.v1.Health/Watch" ||
		strings.Contains(method, "grpc.reflection")
}

// isWebhookMethod returns true for gateway webhook callbacks that arrive without user identity.
// These are authenticated via HMAC signature in the service layer, not via JWT/session.
func isWebhookMethod(method string) bool {
	return strings.Contains(method, "HandleGatewayWebhook")
}

// firstMD returns the first non-empty value for a metadata key, or "".
func firstMD(md metadata.MD, key string) string {
	for _, v := range md.Get(key) {
		if v = strings.TrimSpace(v); v != "" {
			return v
		}
	}
	return ""
}
