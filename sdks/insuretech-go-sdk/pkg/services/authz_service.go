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

// CreatePolicyRule Create policy rule
func (s *AuthzService) CreatePolicyRule(ctx context.Context, req *models.PolicyRuleCreationRequest) (*models.PolicyRuleCreationResponse, error) {
	path := "/v1/authz/policies"
	var result models.PolicyRuleCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListPolicyRules List policy rules
func (s *AuthzService) ListPolicyRules(ctx context.Context) (*models.PolicyRulesListingResponse, error) {
	path := "/v1/authz/policies"
	var result models.PolicyRulesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CheckAccess CheckAccess — single authorization check (gateway + per-service interceptor)
func (s *AuthzService) CheckAccess(ctx context.Context, req *models.CheckAccessRequest) (*models.CheckAccessResponse, error) {
	path := "/v1/authz/check"
	var result models.CheckAccessResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePolicyRule Update policy rule
func (s *AuthzService) UpdatePolicyRule(ctx context.Context, policyId string, req *models.PolicyRuleUpdateRequest) (*models.PolicyRuleUpdateResponse, error) {
	path := "/v1/authz/policies/{policy_id}"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.PolicyRuleUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePolicyRule Delete policy rule
func (s *AuthzService) DeletePolicyRule(ctx context.Context, policyId string) error {
	path := "/v1/authz/policies/{policy_id}"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// GetJWKS GetJWKS — serves the RS256 public key set for JWT verification
func (s *AuthzService) GetJWKS(ctx context.Context) (*models.AuthzJWKSRetrievalResponse, error) {
	path := "/.well-known/jwks.json"
	var result models.AuthzJWKSRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPortalConfig Get portal config
func (s *AuthzService) GetPortalConfig(ctx context.Context, portal string) (*models.PortalConfigRetrievalResponse, error) {
	path := "/v1/authz/portals/{portal}/config"
	path = strings.ReplaceAll(path, "{portal}", portal)
	var result models.PortalConfigRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePortalConfig Update portal config
func (s *AuthzService) UpdatePortalConfig(ctx context.Context, portal string, req *models.PortalConfigUpdateRequest) (*models.PortalConfigUpdateResponse, error) {
	path := "/v1/authz/portals/{portal}/config"
	path = strings.ReplaceAll(path, "{portal}", portal)
	var result models.PortalConfigUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAccessDecisionAudits List access decision audits
func (s *AuthzService) ListAccessDecisionAudits(ctx context.Context) (*models.AccessDecisionAuditsListingResponse, error) {
	path := "/v1/authz/audits"
	var result models.AccessDecisionAuditsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchCheckAccess BatchCheckAccess — check multiple (sub, dom, obj, act) tuples in one call
func (s *AuthzService) BatchCheckAccess(ctx context.Context, req *models.BatchCheckAccessRequest) (*models.BatchCheckAccessResponse, error) {
	path := "/v1/authz/check:batch"
	var result models.BatchCheckAccessResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRole Get role
func (s *AuthzService) GetRole(ctx context.Context, roleId string) (*models.RoleRetrievalResponse, error) {
	path := "/v1/authz/roles/{role_id}"
	path = strings.ReplaceAll(path, "{role_id}", roleId)
	var result models.RoleRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateRole Update role
func (s *AuthzService) UpdateRole(ctx context.Context, roleId string, req *models.RoleUpdateRequest) (*models.RoleUpdateResponse, error) {
	path := "/v1/authz/roles/{role_id}"
	path = strings.ReplaceAll(path, "{role_id}", roleId)
	var result models.RoleUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteRole Delete role
func (s *AuthzService) DeleteRole(ctx context.Context, roleId string) error {
	path := "/v1/authz/roles/{role_id}"
	path = strings.ReplaceAll(path, "{role_id}", roleId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// GetUserPermissions GetUserPermissions — resolves all effective permissions for a user in a domain
func (s *AuthzService) GetUserPermissions(ctx context.Context, userId string) (*models.UserPermissionsRetrievalResponse, error) {
	path := "/v1/authz/users/{user_id}/permissions"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.UserPermissionsRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveRole Remove role
func (s *AuthzService) RemoveRole(ctx context.Context, userId string, roleId string) error {
	path := "/v1/authz/users/{user_id}/roles/{role_id}"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	path = strings.ReplaceAll(path, "{role_id}", roleId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// AssignRole AssignRole — assign a role to a user within domain (portal:tenant_id)
func (s *AuthzService) AssignRole(ctx context.Context, userId string, req *models.RoleAssignmentRequest) (*models.RoleAssignmentResponse, error) {
	path := "/v1/authz/users/{user_id}/roles"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.RoleAssignmentResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUserRoles List user roles
func (s *AuthzService) ListUserRoles(ctx context.Context, userId string) (*models.UserRolesListingResponse, error) {
	path := "/v1/authz/users/{user_id}/roles"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.UserRolesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListRoles List roles
func (s *AuthzService) ListRoles(ctx context.Context) (*models.RolesListingResponse, error) {
	path := "/v1/authz/roles"
	var result models.RolesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateRole Create role
func (s *AuthzService) CreateRole(ctx context.Context, req *models.RoleCreationRequest) (*models.RoleCreationResponse, error) {
	path := "/v1/authz/roles"
	var result models.RoleCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

