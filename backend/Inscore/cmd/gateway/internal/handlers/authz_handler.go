package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/respond"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// AuthZHandler proxies HTTP requests to the AuthZ gRPC service.
// AuthZ is a standalone isolated gRPC service — accessible by all portals
// (B2B, B2C, system, agent, regulator) without PoliSync HTTP proxy.
type AuthZHandler struct {
	client authzv1.AuthZServiceClient
}

// NewAuthZHandler creates a handler backed by the given authz gRPC connection.
func NewAuthZHandler(conn *grpc.ClientConn) *AuthZHandler {
	return &AuthZHandler{client: authzv1.NewAuthZServiceClient(conn)}
}

// call is the generic helper that forwards identity metadata to gRPC context.
func (h *AuthZHandler) call(w http.ResponseWriter, r *http.Request,
	fn func(ctx context.Context, body []byte) (proto.Message, error)) {

	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
		r.Body.Close()
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Build cookie for server-side session forwarding
	cookieHeader := r.Header.Get("Cookie")
	if cookieHeader == "" {
		if st := r.Header.Get("X-Session-Token"); st != "" {
			cookieHeader = "session_token=" + st
		}
	}

	ctx = grpcmeta.WithOutgoingMetadata(ctx,
		"x-user-id", r.Header.Get("X-User-ID"),
		"x-session-id", r.Header.Get("X-Session-ID"),
		"x-session-type", r.Header.Get("X-Session-Type"),
		"x-user-type", r.Header.Get("X-User-Type"),
		"x-portal", r.Header.Get("X-Portal"),
		"x-tenant-id", r.Header.Get("X-Tenant-ID"),
		"x-email", r.Header.Get("X-Email"),
		"x-roles", r.Header.Get("X-Roles"),
		"authorization", r.Header.Get("Authorization"),
		"cookie", cookieHeader,
		"x-forwarded-for", r.Header.Get("X-Forwarded-For"),
		"user-agent", r.UserAgent(),
	)

	msg, err := fn(ctx, body)
	if err != nil {
		respond.GRPCError(w, r, err)
		return
	}

	b, _ := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	respond.RawProtoJSON(w, r, b, http.StatusOK)
}

// ── Roles ─────────────────────────────────────────────────────────────────────

func (h *AuthZHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.ListRolesRequest{}
		if len(body) > 0 {
			_ = protoUnmarshal(body, req)
		}
		return h.client.ListRoles(ctx, req)
	})
}

func (h *AuthZHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.CreateRoleRequest{}
		if err := protoUnmarshal(body, req); err != nil {
			return nil, err
		}
		return h.client.CreateRole(ctx, req)
	})
}

func (h *AuthZHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		return h.client.GetRole(ctx, &authzv1.GetRoleRequest{
			RoleId: r.PathValue("role_id"),
		})
	})
}

func (h *AuthZHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.UpdateRoleRequest{}
		if err := protoUnmarshal(body, req); err != nil {
			return nil, err
		}
		if req.RoleId == "" {
			req.RoleId = r.PathValue("role_id")
		}
		return h.client.UpdateRole(ctx, req)
	})
}

func (h *AuthZHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		return h.client.DeleteRole(ctx, &authzv1.DeleteRoleRequest{
			RoleId: r.PathValue("role_id"),
		})
	})
}

// ── User Role Assignment ──────────────────────────────────────────────────────

func (h *AuthZHandler) ListUserRoles(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		return h.client.ListUserRoles(ctx, &authzv1.ListUserRolesRequest{
			UserId: r.PathValue("user_id"),
		})
	})
}

func (h *AuthZHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.AssignRoleRequest{}
		if err := protoUnmarshal(body, req); err != nil {
			return nil, err
		}
		if req.UserId == "" {
			req.UserId = r.PathValue("user_id")
		}
		return h.client.AssignRole(ctx, req)
	})
}

func (h *AuthZHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		return h.client.RemoveRole(ctx, &authzv1.RemoveRoleRequest{
			UserId: r.PathValue("user_id"),
			RoleId: r.PathValue("role_id"),
		})
	})
}

// ── Permissions ───────────────────────────────────────────────────────────────

func (h *AuthZHandler) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		return h.client.GetUserPermissions(ctx, &authzv1.GetUserPermissionsRequest{
			UserId: r.PathValue("user_id"),
			Domain: r.URL.Query().Get("domain"),
		})
	})
}

// ── Policy Rules ──────────────────────────────────────────────────────────────

func (h *AuthZHandler) ListPolicies(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		return h.client.ListPolicyRules(ctx, &authzv1.ListPolicyRulesRequest{})
	})
}

func (h *AuthZHandler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.CreatePolicyRuleRequest{}
		if err := protoUnmarshal(body, req); err != nil {
			return nil, err
		}
		return h.client.CreatePolicyRule(ctx, req)
	})
}

func (h *AuthZHandler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.UpdatePolicyRuleRequest{}
		if err := protoUnmarshal(body, req); err != nil {
			return nil, err
		}
		if req.PolicyId == "" {
			req.PolicyId = r.PathValue("policy_id")
		}
		return h.client.UpdatePolicyRule(ctx, req)
	})
}

func (h *AuthZHandler) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		return h.client.DeletePolicyRule(ctx, &authzv1.DeletePolicyRuleRequest{
			PolicyId: r.PathValue("policy_id"),
		})
	})
}

// ── Access Checks ─────────────────────────────────────────────────────────────

func (h *AuthZHandler) CheckAccess(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.CheckAccessRequest{}
		if err := protoUnmarshal(body, req); err != nil {
			return nil, err
		}
		// Inject user_id and domain from JWT-validated headers (set by auth middleware).
		// This prevents clients from spoofing their own user_id.
		if uid := r.Header.Get("X-User-ID"); uid != "" && req.UserId == "" {
			req.UserId = uid
		}
		// Build domain using same logic as authz_middleware.go buildDomain():
		// portal:tenantID, with "root" fallback for empty tenant (e.g. B2C users).
		if req.Domain == "" {
			portal := strings.ToLower(strings.TrimPrefix(r.Header.Get("X-Portal"), "PORTAL_"))
			if portal == "" || portal == "unspecified" {
				portal = "b2c"
			}
			tenantID := r.Header.Get("X-Tenant-ID")
			if portal == "system" {
				req.Domain = "system:root"
			} else {
				if tenantID == "" {
					tenantID = "root"
				}
				req.Domain = portal + ":" + tenantID
			}
		}
		return h.client.CheckAccess(ctx, req)
	})
}

func (h *AuthZHandler) BatchCheckAccess(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.BatchCheckAccessRequest{}
		if err := protoUnmarshal(body, req); err != nil {
			return nil, err
		}
		return h.client.BatchCheckAccess(ctx, req)
	})
}

// ── Audits ────────────────────────────────────────────────────────────────────

func (h *AuthZHandler) ListAudits(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.ListAccessDecisionAuditsRequest{}
		if len(body) > 0 {
			_ = protoUnmarshal(body, req)
		}
		return h.client.ListAccessDecisionAudits(ctx, req)
	})
}

// ── Portal Config ─────────────────────────────────────────────────────────────

func (h *AuthZHandler) ListPortalConfigs(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		return h.client.ListPortalConfigs(ctx, &authzv1.ListPortalConfigsRequest{})
	})
}

// portalFromPath converts a URL path segment like "system" or "PORTAL_SYSTEM"
// to the authz Portal enum value. Case-insensitive.
func portalFromPath(s string) authzentityv1.Portal {
	// Try with PORTAL_ prefix first, then without
	key := strings.ToUpper(s)
	if !strings.HasPrefix(key, "PORTAL_") {
		key = "PORTAL_" + key
	}
	if v, ok := authzentityv1.Portal_value[key]; ok {
		return authzentityv1.Portal(v)
	}
	// Direct numeric or unspecified
	return authzentityv1.Portal_PORTAL_UNSPECIFIED
}

func (h *AuthZHandler) GetPortalConfig(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		portal := portalFromPath(r.PathValue("portal"))
		if portal == authzentityv1.Portal_PORTAL_UNSPECIFIED {
			return nil, fmt.Errorf("unknown portal: %q", r.PathValue("portal"))
		}
		return h.client.GetPortalConfig(ctx, &authzv1.GetPortalConfigRequest{
			Portal: portal,
		})
	})
}

func (h *AuthZHandler) UpdatePortalConfig(w http.ResponseWriter, r *http.Request) {
	h.call(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		req := &authzv1.UpdatePortalConfigRequest{}
		if err := protoUnmarshal(body, req); err != nil {
			return nil, err
		}
		if req.Portal == authzentityv1.Portal_PORTAL_UNSPECIFIED {
			req.Portal = portalFromPath(r.PathValue("portal"))
		}
		return h.client.UpdatePortalConfig(ctx, req)
	})
}
