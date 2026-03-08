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

// AuthZClient interface for authorization checks
type AuthZClient interface {
	CheckAccess(ctx context.Context, req *authzservicev1.CheckAccessRequest) (*authzservicev1.CheckAccessResponse, error)
}

// AuthZInterceptor provides authorization enforcement for B2B service
type AuthZInterceptor struct {
	authzClient AuthZClient
}

func NewAuthZInterceptor(authzClient AuthZClient) *AuthZInterceptor {
	return &AuthZInterceptor{
		authzClient: authzClient,
	}
}

// UnaryServerInterceptor checks authorization for unary RPCs
func (i *AuthZInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		userIDs := md.Get("x-user-id")
		if len(userIDs) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing user_id")
		}
		userID := userIDs[0]

		// Detect portal: PORTAL_SYSTEM users (super_admin) must bypass org context checks.
		// The gateway sets x-portal = "PORTAL_SYSTEM" for system user JWTs.
		portalRaw := firstNonEmpty(md.Get("x-portal"))
		portalNorm := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(portalRaw, "PORTAL_")))
		isSystemPortal := portalNorm == "system"

		orgIDs := md.Get("x-business-id")
		organisationID := ""
		for _, candidate := range orgIDs {
			candidate = strings.TrimSpace(candidate)
			if candidate != "" {
				organisationID = candidate
				break
			}
		}

		// For non-system portals with no org context:
		// - Methods that allow no-org calls (e.g. ResolveMyOrganisation) pass through without Casbin.
		// - All other methods require org context to proceed.
		// System portal users always proceed to Casbin check using system:root domain.
		if !isSystemPortal && organisationID == "" {
			if isNoOrgContextMethod(info.FullMethod) {
				// These methods explicitly allow no-org calls without auth check
				// (e.g. ResolveMyOrganisation which bootstraps the session).
				return handler(ctx, req)
			}
			return nil, status.Error(codes.PermissionDenied, "missing organisation context")
		}

		// Map gRPC method to the Casbin object/action used by seeded policies.
		resource, action := mapMethodToResourceAction(info.FullMethod)
		if resource == "" {
			// No Casbin check required for this method (e.g. ResolveMyOrganisation).
			return handler(ctx, req)
		}

		// Resolve the Casbin domain and run the authorization check.
		domain := resolveAuthzDomain(md, organisationID)
		resp, err := i.authzClient.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
			UserId: userID,
			Domain: domain,
			Object: resource,
			Action: action,
		})
		if err != nil {
			appLogger.Errorf("AuthZ check failed: %v", err)
			return nil, status.Error(codes.Internal, "authorization check failed")
		}

		if !resp.Allowed {
			appLogger.Warnf("Access denied: user=%s, domain=%s, resource=%s, action=%s",
				userID, domain, resource, action)
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}

		return handler(ctx, req)
	}
}

// mapMethodToResourceAction maps gRPC method names to resource and action.
// Returns ("", "") to skip the Casbin check entirely (only for truly public/no-auth methods).
func mapMethodToResourceAction(method string) (resource, action string) {
	// Method format: /insuretech.b2b.services.v1.B2BService/MethodName
	parts := strings.Split(method, "/")
	if len(parts) < 3 {
		return "", ""
	}
	methodName := parts[len(parts)-1]

	switch {
	// ResolveMyOrganisation is a bootstrap call — no authz needed (user just logged in).
	case strings.HasPrefix(methodName, "ResolveMyOrganisation"):
		return "", ""

	// READ operations
	case strings.HasPrefix(methodName, "Get"), strings.HasPrefix(methodName, "List"):
		return "svc:b2b/*", "GET"

	// CREATE/ADD/ASSIGN operations
	case strings.HasPrefix(methodName, "Create"),
		strings.HasPrefix(methodName, "Add"),
		strings.HasPrefix(methodName, "Assign"):
		return "svc:b2b/*", "POST"

	// UPDATE operations
	case strings.HasPrefix(methodName, "Update"):
		return "svc:b2b/*", "PATCH"

	// DELETE/REMOVE operations
	case strings.HasPrefix(methodName, "Delete"),
		strings.HasPrefix(methodName, "Remove"):
		return "svc:b2b/*", "DELETE"

	default:
		appLogger.Warnf("Unknown method for authz mapping: %s", methodName)
		return "", ""
	}
}

func resolveAuthzDomain(md metadata.MD, organisationID string) string {
	portal := firstNonEmpty(md.Get("x-portal"))
	tenantID := firstNonEmpty(md.Get("x-tenant-id"))

	portal = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(portal, "PORTAL_")))
	switch portal {
	case "system":
		return "system:root"
	case "b2b":
		if strings.TrimSpace(organisationID) != "" {
			return fmt.Sprintf("b2b:%s", organisationID)
		}
		if tenantID != "" {
			return fmt.Sprintf("b2b:%s", tenantID)
		}
		return "b2b:root"
	default:
		if strings.TrimSpace(organisationID) != "" {
			return fmt.Sprintf("b2b:%s", organisationID)
		}
		if tenantID != "" {
			return fmt.Sprintf("%s:%s", portal, tenantID)
		}
		return "b2b:root"
	}
}

func firstNonEmpty(values []string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

// isNoOrgContextMethod returns true for methods that are allowed to proceed
// without an x-business-id header AND without a Casbin check.
// Only truly public/bootstrap methods belong here — all others must go through
// the normal Casbin path (system:root for super_admin, b2b:{org_id} for b2b admin).
func isNoOrgContextMethod(method string) bool {
	// ResolveMyOrganisation is the session bootstrap call: the gateway calls it
	// immediately after login to discover the user's org_id. It cannot require
	// an org_id in the request because that's exactly what it returns.
	noAuthMethods := []string{
		"ResolveMyOrganisation",
	}
	for _, m := range noAuthMethods {
		if strings.Contains(method, m) {
			return true
		}
	}
	return false
}
