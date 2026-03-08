package middleware

import (
	"context"
	"net/http"
	"strings"

	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// B2BContextMiddleware resolves organisation context for B2B requests
type B2BContextMiddleware struct {
	b2bClient b2bservicev1.B2BServiceClient
}

func NewB2BContextMiddleware(conn *grpc.ClientConn) *B2BContextMiddleware {
	return &B2BContextMiddleware{
		b2bClient: b2bservicev1.NewB2BServiceClient(conn),
	}
}

// InjectOrganisationContext resolves the organisation for the authenticated user
// and injects it as x-business-id metadata.
//
// System portal users (superadmin) are exempt: they operate in the system:root
// Casbin domain and do not belong to any organisation, so we skip the lookup
// entirely and let the request through. The downstream b2b gRPC interceptor
// handles system:root via the x-portal=PORTAL_SYSTEM header set by AuthMiddleware.
func (m *B2BContextMiddleware) InjectOrganisationContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// System portal users have no org context — skip resolution entirely.
		// AuthMiddleware already set X-Portal from the validated JWT.
		portal := strings.TrimSpace(r.Header.Get("X-Portal"))
		portalNorm := strings.ToLower(strings.TrimPrefix(portal, "PORTAL_"))
		if portalNorm == "system" {
			next.ServeHTTP(w, r)
			return
		}

		// If the client already supplied a business_id (header or query param), use it.
		if businessID := resolveBusinessID(r); businessID != "" {
			next.ServeHTTP(w, withOrganisationContext(r, businessID, r.Header.Get("X-Org-Role")))
			return
		}

		userIDStr := strings.TrimSpace(r.Header.Get("X-User-ID"))
		if userIDStr == "" {
			if userID, ok := r.Context().Value("user_id").(string); ok {
				userIDStr = strings.TrimSpace(userID)
			}
		}
		if userIDStr == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := m.b2bClient.ResolveMyOrganisation(r.Context(), &b2bservicev1.ResolveMyOrganisationRequest{
			UserId: userIDStr,
		})
		if err != nil {
			http.Error(w, "No organisation found for user", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, withOrganisationContext(r, resp.OrganisationId, resp.Role.String()))
	})
}

func resolveBusinessID(r *http.Request) string {
	if businessID := strings.TrimSpace(r.Header.Get("X-Business-ID")); businessID != "" {
		return businessID
	}
	if businessID := strings.TrimSpace(r.URL.Query().Get("business_id")); businessID != "" {
		r.Header.Set("X-Business-ID", businessID)
		return businessID
	}
	return ""
}

func withOrganisationContext(r *http.Request, organisationID, orgRole string) *http.Request {
	organisationID = strings.TrimSpace(organisationID)
	if organisationID == "" {
		return r
	}

	r.Header.Set("X-Business-ID", organisationID)
	if strings.TrimSpace(orgRole) != "" {
		r.Header.Set("X-Org-Role", orgRole)
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "organisation_id", organisationID)
	if strings.TrimSpace(orgRole) != "" {
		ctx = context.WithValue(ctx, "org_member_role", orgRole)
	}

	md := metadata.New(map[string]string{
		"x-business-id": organisationID,
		"x-user-id":     r.Header.Get("X-User-ID"),
		"x-org-role":    orgRole,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	return r.WithContext(ctx)
}

// InjectUserContext returns a middleware that extracts user_id from session and injects it into context
func InjectUserContext(authClient authnv1.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract session token from cookie
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate session with AuthN service and extract user_id
			resp, err := authClient.ValidateToken(r.Context(), &authnv1.ValidateTokenRequest{
				SessionId: cookie.Value,
			})
			if err != nil || resp == nil || !resp.Valid || resp.UserId == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
