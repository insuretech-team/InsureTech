package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// AuthzService handles authz-related API calls
type AuthzService struct {
	Client Client
}

// GetJWKS GetJWKS — serves the RS256 public key set for JWT verification
func (s *AuthzService) GetJWKS(ctx context.Context) error {
	path := "/.well-known/jwks.json"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ListAccessDecisionAudits List access decision audits
func (s *AuthzService) ListAccessDecisionAudits(ctx context.Context) error {
	path := "/v1/authz/audits"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CheckAccess CheckAccess — single authorization check (gateway + per-service interceptor)
func (s *AuthzService) CheckAccess(ctx context.Context, req *models.CheckAccessRequest) error {
	path := "/v1/authz/check"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// BatchCheckAccess BatchCheckAccess — check multiple (sub, dom, obj, act) tuples in one call
func (s *AuthzService) BatchCheckAccess(ctx context.Context, req *models.BatchCheckAccessRequest) error {
	path := "/v1/authz/check:batch"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListPolicyRules List policy rules
func (s *AuthzService) ListPolicyRules(ctx context.Context) error {
	path := "/v1/authz/policies"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreatePolicyRule Create policy rule
func (s *AuthzService) CreatePolicyRule(ctx context.Context, req *models.PolicyRuleCreationRequest) error {
	path := "/v1/authz/policies"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdatePolicyRule Update policy rule
func (s *AuthzService) UpdatePolicyRule(ctx context.Context, policyId string, req *models.PolicyRuleUpdateRequest) error {
	path := "/v1/authz/policies/{policy_id}"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeletePolicyRule Delete policy rule
func (s *AuthzService) DeletePolicyRule(ctx context.Context, policyId string) error {
	path := "/v1/authz/policies/{policy_id}"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// ListPortalConfigs List portal configs
func (s *AuthzService) ListPortalConfigs(ctx context.Context) error {
	path := "/v1/authz/portals/configs"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetPortalConfig Get portal config
func (s *AuthzService) GetPortalConfig(ctx context.Context, portal string) error {
	path := "/v1/authz/portals/{portal}/config"
	path = strings.ReplaceAll(path, "{portal}", portal)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdatePortalConfig Update portal config
func (s *AuthzService) UpdatePortalConfig(ctx context.Context, portal string, req *models.PortalConfigUpdateRequest) error {
	path := "/v1/authz/portals/{portal}/config"
	path = strings.ReplaceAll(path, "{portal}", portal)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// ListRoles List roles
func (s *AuthzService) ListRoles(ctx context.Context) error {
	path := "/v1/authz/roles"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateRole Create role
func (s *AuthzService) CreateRole(ctx context.Context, req *models.RoleCreationRequest) error {
	path := "/v1/authz/roles"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetRole Get role
func (s *AuthzService) GetRole(ctx context.Context, roleId string) error {
	path := "/v1/authz/roles/{role_id}"
	path = strings.ReplaceAll(path, "{role_id}", roleId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateRole Update role
func (s *AuthzService) UpdateRole(ctx context.Context, roleId string, req *models.RoleUpdateRequest) error {
	path := "/v1/authz/roles/{role_id}"
	path = strings.ReplaceAll(path, "{role_id}", roleId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeleteRole Delete role
func (s *AuthzService) DeleteRole(ctx context.Context, roleId string) error {
	path := "/v1/authz/roles/{role_id}"
	path = strings.ReplaceAll(path, "{role_id}", roleId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// GetUserPermissions GetUserPermissions — resolves all effective permissions for a user in a domain
func (s *AuthzService) GetUserPermissions(ctx context.Context, userId string) error {
	path := "/v1/authz/users/{user_id}/permissions"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ListUserRoles List user roles
func (s *AuthzService) ListUserRoles(ctx context.Context, userId string) error {
	path := "/v1/authz/users/{user_id}/roles"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// AssignRole AssignRole — assign a role to a user within domain (portal:tenant_id)
func (s *AuthzService) AssignRole(ctx context.Context, userId string, req *models.RoleAssignmentRequest) error {
	path := "/v1/authz/users/{user_id}/roles"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RemoveRole Remove role
func (s *AuthzService) RemoveRole(ctx context.Context, userId string, roleId string) error {
	path := "/v1/authz/users/{user_id}/roles/{role_id}"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	path = strings.ReplaceAll(path, "{role_id}", roleId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

