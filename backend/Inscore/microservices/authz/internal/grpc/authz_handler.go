package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/service"
)

// AuthZHandler implements authzservicev1.AuthZServiceServer by delegating
// to the AuthZService business logic layer.
type AuthZHandler struct {
	authzservicev1.UnimplementedAuthZServiceServer
	svc *service.AuthZService
}

// NewAuthZHandler creates a new AuthZHandler.
func NewAuthZHandler(svc *service.AuthZService) *AuthZHandler {
	return &AuthZHandler{svc: svc}
}

// ── Core Enforcement ─────────────────────────────────────────────────────────

func (h *AuthZHandler) CheckAccess(ctx context.Context, req *authzservicev1.CheckAccessRequest) (*authzservicev1.CheckAccessResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Domain == "" {
		return nil, status.Error(codes.InvalidArgument, "domain is required (format: portal:tenant_id)")
	}
	if req.Object == "" {
		return nil, status.Error(codes.InvalidArgument, "object is required (format: svc:service/resource)")
	}
	if req.Action == "" {
		return nil, status.Error(codes.InvalidArgument, "action is required")
	}
	resp, err := h.svc.CheckAccess(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "access check failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) BatchCheckAccess(ctx context.Context, req *authzservicev1.BatchCheckAccessRequest) (*authzservicev1.BatchCheckAccessResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Domain == "" {
		return nil, status.Error(codes.InvalidArgument, "domain is required")
	}
	if len(req.Checks) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one check is required")
	}
	resp, err := h.svc.BatchCheckAccess(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "batch check failed: %v", err)
	}
	return resp, nil
}

// ── Role Management ──────────────────────────────────────────────────────────

func (h *AuthZHandler) CreateRole(ctx context.Context, req *authzservicev1.CreateRoleRequest) (*authzservicev1.CreateRoleResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Portal == 0 {
		return nil, status.Error(codes.InvalidArgument, "portal is required")
	}
	resp, err := h.svc.CreateRole(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create role failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) GetRole(ctx context.Context, req *authzservicev1.GetRoleRequest) (*authzservicev1.GetRoleResponse, error) {
	if req.RoleId == "" {
		return nil, status.Error(codes.InvalidArgument, "role_id is required")
	}
	resp, err := h.svc.GetRole(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get role failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) UpdateRole(ctx context.Context, req *authzservicev1.UpdateRoleRequest) (*authzservicev1.UpdateRoleResponse, error) {
	if req.RoleId == "" {
		return nil, status.Error(codes.InvalidArgument, "role_id is required")
	}
	resp, err := h.svc.UpdateRole(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update role failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) DeleteRole(ctx context.Context, req *authzservicev1.DeleteRoleRequest) (*authzservicev1.DeleteRoleResponse, error) {
	if req.RoleId == "" {
		return nil, status.Error(codes.InvalidArgument, "role_id is required")
	}
	resp, err := h.svc.DeleteRole(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete role failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) ListRoles(ctx context.Context, req *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error) {
	resp, err := h.svc.ListRoles(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list roles failed: %v", err)
	}
	return resp, nil
}

// ── User-Role Assignment ─────────────────────────────────────────────────────

func (h *AuthZHandler) AssignRole(ctx context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error) {
	if req.UserId == "" || req.RoleId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and role_id are required")
	}
	if req.Domain == "" {
		return nil, status.Error(codes.InvalidArgument, "domain is required (format: portal:tenant_id)")
	}
	resp, err := h.svc.AssignRole(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "assign role failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) RemoveRole(ctx context.Context, req *authzservicev1.RemoveRoleRequest) (*authzservicev1.RemoveRoleResponse, error) {
	if req.UserId == "" || req.RoleId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and role_id are required")
	}
	if req.Domain == "" {
		return nil, status.Error(codes.InvalidArgument, "domain is required")
	}
	resp, err := h.svc.RemoveRole(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "remove role failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) ListUserRoles(ctx context.Context, req *authzservicev1.ListUserRolesRequest) (*authzservicev1.ListUserRolesResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.svc.ListUserRoles(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list user roles failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) GetUserPermissions(ctx context.Context, req *authzservicev1.GetUserPermissionsRequest) (*authzservicev1.GetUserPermissionsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Domain == "" {
		return nil, status.Error(codes.InvalidArgument, "domain is required")
	}
	resp, err := h.svc.GetUserPermissions(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user permissions failed: %v", err)
	}
	return resp, nil
}

// ── Policy Rule Management ────────────────────────────────────────────────────

func (h *AuthZHandler) CreatePolicyRule(ctx context.Context, req *authzservicev1.CreatePolicyRuleRequest) (*authzservicev1.CreatePolicyRuleResponse, error) {
	if req.Subject == "" || req.Domain == "" || req.Object == "" || req.Action == "" {
		return nil, status.Error(codes.InvalidArgument, "subject, domain, object, and action are required")
	}
	resp, err := h.svc.CreatePolicyRule(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create policy rule failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) UpdatePolicyRule(ctx context.Context, req *authzservicev1.UpdatePolicyRuleRequest) (*authzservicev1.UpdatePolicyRuleResponse, error) {
	if req.PolicyId == "" {
		return nil, status.Error(codes.InvalidArgument, "policy_id is required")
	}
	resp, err := h.svc.UpdatePolicyRule(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update policy rule failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) DeletePolicyRule(ctx context.Context, req *authzservicev1.DeletePolicyRuleRequest) (*authzservicev1.DeletePolicyRuleResponse, error) {
	if req.PolicyId == "" {
		return nil, status.Error(codes.InvalidArgument, "policy_id is required")
	}
	resp, err := h.svc.DeletePolicyRule(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete policy rule failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) ListPolicyRules(ctx context.Context, req *authzservicev1.ListPolicyRulesRequest) (*authzservicev1.ListPolicyRulesResponse, error) {
	resp, err := h.svc.ListPolicyRules(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list policy rules failed: %v", err)
	}
	return resp, nil
}

// ── Portal Configuration ──────────────────────────────────────────────────────

func (h *AuthZHandler) GetPortalConfig(ctx context.Context, req *authzservicev1.GetPortalConfigRequest) (*authzservicev1.GetPortalConfigResponse, error) {
	if req.Portal == 0 {
		return nil, status.Error(codes.InvalidArgument, "portal is required")
	}
	resp, err := h.svc.GetPortalConfig(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get portal config failed: %v", err)
	}
	return resp, nil
}

func (h *AuthZHandler) UpdatePortalConfig(ctx context.Context, req *authzservicev1.UpdatePortalConfigRequest) (*authzservicev1.UpdatePortalConfigResponse, error) {
	if req.Portal == authzentityv1.Portal_PORTAL_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "portal is required")
	}
	resp, err := h.svc.UpdatePortalConfig(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update portal config failed: %v", err)
	}
	return resp, nil
}

// ── Audit ─────────────────────────────────────────────────────────────────────

func (h *AuthZHandler) ListAccessDecisionAudits(ctx context.Context, req *authzservicev1.ListAccessDecisionAuditsRequest) (*authzservicev1.ListAccessDecisionAuditsResponse, error) {
	resp, err := h.svc.ListAccessDecisionAudits(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list audits failed: %v", err)
	}
	return resp, nil
}

// ── Cache Invalidation ────────────────────────────────────────────────────────

func (h *AuthZHandler) InvalidatePolicyCache(ctx context.Context, req *authzservicev1.InvalidatePolicyCacheRequest) (*authzservicev1.InvalidatePolicyCacheResponse, error) {
	resp, err := h.svc.InvalidatePolicyCache(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalidate cache failed: %v", err)
	}
	return resp, nil
}

// ── JWKS ──────────────────────────────────────────────────────────────────────

func (h *AuthZHandler) GetJWKS(ctx context.Context, req *authzservicev1.GetJWKSRequest) (*authzservicev1.GetJWKSResponse, error) {
	resp, err := h.svc.GetJWKS(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get JWKS failed: %v", err)
	}
	return resp, nil
}
