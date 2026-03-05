package grpc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"testing"

	authzevents "github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/service"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAuthZHandler_ValidationErrors(t *testing.T) {
	h := NewAuthZHandler(nil)
	ctx := context.Background()

	_, err := h.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.BatchCheckAccess(ctx, &authzservicev1.BatchCheckAccessRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.CreateRole(ctx, &authzservicev1.CreateRoleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.GetRole(ctx, &authzservicev1.GetRoleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.UpdateRole(ctx, &authzservicev1.UpdateRoleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.DeleteRole(ctx, &authzservicev1.DeleteRoleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.AssignRole(ctx, &authzservicev1.AssignRoleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.RemoveRole(ctx, &authzservicev1.RemoveRoleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.ListUserRoles(ctx, &authzservicev1.ListUserRolesRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.GetUserPermissions(ctx, &authzservicev1.GetUserPermissionsRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.UpdatePolicyRule(ctx, &authzservicev1.UpdatePolicyRuleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.DeletePolicyRule(ctx, &authzservicev1.DeletePolicyRuleRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.GetPortalConfig(ctx, &authzservicev1.GetPortalConfigRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.UpdatePortalConfig(ctx, &authzservicev1.UpdatePortalConfigRequest{Portal: authzentityv1.Portal_PORTAL_UNSPECIFIED})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

type fakeEnforcer struct{}

func (f *fakeEnforcer) Enforce(context.Context, string, string, string, string) (bool, string, error) {
	return true, "p", nil
}
func (f *fakeEnforcer) AddPolicy(string, string, string, string, string) error { return nil }
func (f *fakeEnforcer) RemovePolicy(string, string, string, string) error       { return nil }
func (f *fakeEnforcer) AddRoleForUserInDomain(string, string, string) error     { return nil }
func (f *fakeEnforcer) DeleteRoleForUserInDomain(string, string, string) error   { return nil }
func (f *fakeEnforcer) GetRolesForUserInDomain(string, string) ([]string, error) {
	return []string{"role:admin"}, nil
}
func (f *fakeEnforcer) GetPermissionsForUserInDomain(string, string) ([][]string, error) {
	return [][]string{{"role:admin", "system:root", "svc:user/*", "GET", "allow"}}, nil
}
func (f *fakeEnforcer) InvalidateCache() error { return nil }

type fakeRoleRepo struct {
	role *authzentityv1.Role
}

func (f *fakeRoleRepo) Create(context.Context, *authzentityv1.Role) (*authzentityv1.Role, error) {
	return &authzentityv1.Role{RoleId: "r1", Name: "admin", Portal: authzentityv1.Portal_PORTAL_SYSTEM}, nil
}
func (f *fakeRoleRepo) GetByID(context.Context, string) (*authzentityv1.Role, error) {
	if f.role != nil {
		return f.role, nil
	}
	return &authzentityv1.Role{RoleId: "r1", Name: "admin", Portal: authzentityv1.Portal_PORTAL_SYSTEM}, nil
}
func (f *fakeRoleRepo) GetByNameAndPortal(context.Context, string, authzentityv1.Portal) (*authzentityv1.Role, error) {
	return nil, errors.New("not found")
}
func (f *fakeRoleRepo) Update(context.Context, *authzentityv1.Role) (*authzentityv1.Role, error) {
	return &authzentityv1.Role{RoleId: "r1", Name: "admin", Portal: authzentityv1.Portal_PORTAL_SYSTEM}, nil
}
func (f *fakeRoleRepo) SoftDelete(context.Context, string) error { return nil }
func (f *fakeRoleRepo) List(context.Context, authzentityv1.Portal, bool, int, int) ([]*authzentityv1.Role, error) {
	return []*authzentityv1.Role{{RoleId: "r1", Name: "admin", Portal: authzentityv1.Portal_PORTAL_SYSTEM}}, nil
}

type fakeUserRoleRepo struct{}

func (f *fakeUserRoleRepo) Assign(context.Context, *authzentityv1.UserRole) (*authzentityv1.UserRole, error) {
	return &authzentityv1.UserRole{UserRoleId: "ur1", UserId: "u1", RoleId: "r1", Domain: "system:root"}, nil
}
func (f *fakeUserRoleRepo) Remove(context.Context, string, string, string) error { return nil }
func (f *fakeUserRoleRepo) ListByUser(context.Context, string, string) ([]*authzentityv1.UserRole, error) {
	return []*authzentityv1.UserRole{{UserRoleId: "ur1", UserId: "u1", RoleId: "r1", Domain: "system:root"}}, nil
}

type fakePolicyRepo struct{}

func (f *fakePolicyRepo) Create(context.Context, *authzentityv1.PolicyRule) (*authzentityv1.PolicyRule, error) {
	return &authzentityv1.PolicyRule{PolicyId: "p1", Subject: "role:admin", Domain: "system:root", Object: "svc:user/*", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW}, nil
}
func (f *fakePolicyRepo) Update(context.Context, *authzentityv1.PolicyRule) (*authzentityv1.PolicyRule, error) {
	return &authzentityv1.PolicyRule{PolicyId: "p1", Subject: "role:admin", Domain: "system:root", Object: "svc:user/*", Action: "POST", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW}, nil
}
func (f *fakePolicyRepo) SoftDelete(context.Context, string) error { return nil }
func (f *fakePolicyRepo) List(context.Context, string, bool, int, int) ([]*authzentityv1.PolicyRule, error) {
	return []*authzentityv1.PolicyRule{{PolicyId: "p1", Subject: "role:admin", Domain: "system:root", Object: "svc:user/*", Action: "GET", Effect: authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW}}, nil
}

type fakePortalRepo struct{}

func (f *fakePortalRepo) GetByPortal(context.Context, authzentityv1.Portal) (*authzentityv1.PortalConfig, error) {
	return &authzentityv1.PortalConfig{Portal: authzentityv1.Portal_PORTAL_SYSTEM, MfaRequired: true}, nil
}
func (f *fakePortalRepo) Upsert(context.Context, *authzentityv1.PortalConfig) (*authzentityv1.PortalConfig, error) {
	return &authzentityv1.PortalConfig{Portal: authzentityv1.Portal_PORTAL_SYSTEM, MfaRequired: true}, nil
}

type fakeAuditRepo struct{}

func (f *fakeAuditRepo) Create(context.Context, *authzentityv1.AccessDecisionAudit) error { return nil }
func (f *fakeAuditRepo) List(context.Context, *authzservicev1.ListAccessDecisionAuditsRequest) ([]*authzentityv1.AccessDecisionAudit, int64, error) {
	return []*authzentityv1.AccessDecisionAudit{{AuditId: "a1", UserId: "u1"}}, 1, nil
}

func testHandler(t *testing.T) *AuthZHandler {
	t.Helper()
	svc := service.New(
		&fakeEnforcer{},
		&fakeRoleRepo{},
		&fakeUserRoleRepo{},
		&fakePolicyRepo{},
		&fakePortalRepo{},
		&fakeAuditRepo{},
		true,
		authzevents.NewPublisher(nil),
	)
	return NewAuthZHandler(svc)
}

func TestAuthZHandler_SuccessPaths(t *testing.T) {
	h := testHandler(t)
	ctx := context.Background()

	_, err := h.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{UserId: "u1", Domain: "system:root", Object: "svc:user/get", Action: "GET"})
	require.NoError(t, err)
	_, err = h.BatchCheckAccess(ctx, &authzservicev1.BatchCheckAccessRequest{UserId: "u1", Domain: "system:root", Checks: []*authzservicev1.AccessCheckTuple{{Object: "svc:user/get", Action: "GET"}}})
	require.NoError(t, err)
	_, err = h.CreateRole(ctx, &authzservicev1.CreateRoleRequest{Name: "admin", Portal: authzentityv1.Portal_PORTAL_SYSTEM, CreatedBy: "u1"})
	require.NoError(t, err)
	_, err = h.GetRole(ctx, &authzservicev1.GetRoleRequest{RoleId: "r1"})
	require.NoError(t, err)
	_, err = h.UpdateRole(ctx, &authzservicev1.UpdateRoleRequest{RoleId: "r1", Name: "admin2"})
	require.NoError(t, err)
	_, err = h.DeleteRole(ctx, &authzservicev1.DeleteRoleRequest{RoleId: "r1"})
	require.NoError(t, err)
	_, err = h.ListRoles(ctx, &authzservicev1.ListRolesRequest{})
	require.NoError(t, err)
	_, err = h.AssignRole(ctx, &authzservicev1.AssignRoleRequest{UserId: "u1", RoleId: "r1", Domain: "system:root", AssignedBy: "u1", ExpiresAt: timestamppb.Now()})
	require.NoError(t, err)
	_, err = h.RemoveRole(ctx, &authzservicev1.RemoveRoleRequest{UserId: "u1", RoleId: "r1", Domain: "system:root"})
	require.NoError(t, err)
	_, err = h.ListUserRoles(ctx, &authzservicev1.ListUserRolesRequest{UserId: "u1", Domain: "system:root"})
	require.NoError(t, err)
	_, err = h.GetUserPermissions(ctx, &authzservicev1.GetUserPermissionsRequest{UserId: "u1", Domain: "system:root"})
	require.NoError(t, err)
	_, err = h.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{Subject: "role:admin", Domain: "system:root", Object: "svc:user/*", Action: "GET"})
	require.NoError(t, err)
	_, err = h.UpdatePolicyRule(ctx, &authzservicev1.UpdatePolicyRuleRequest{PolicyId: "p1", Action: "POST", IsActive: true})
	require.NoError(t, err)
	_, err = h.DeletePolicyRule(ctx, &authzservicev1.DeletePolicyRuleRequest{PolicyId: "p1"})
	require.NoError(t, err)
	_, err = h.ListPolicyRules(ctx, &authzservicev1.ListPolicyRulesRequest{})
	require.NoError(t, err)
	_, err = h.GetPortalConfig(ctx, &authzservicev1.GetPortalConfigRequest{Portal: authzentityv1.Portal_PORTAL_SYSTEM})
	require.NoError(t, err)
	_, err = h.UpdatePortalConfig(ctx, &authzservicev1.UpdatePortalConfigRequest{Portal: authzentityv1.Portal_PORTAL_SYSTEM, MfaRequired: true, UpdatedBy: "u1"})
	require.NoError(t, err)
	_, err = h.ListAccessDecisionAudits(ctx, &authzservicev1.ListAccessDecisionAuditsRequest{})
	require.NoError(t, err)
	_, err = h.InvalidatePolicyCache(ctx, &authzservicev1.InvalidatePolicyCacheRequest{})
	require.NoError(t, err)
}

func TestAuthZHandler_GetJWKS_Success(t *testing.T) {
	h := testHandler(t)
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	require.NoError(t, err)
	p := t.TempDir() + "/pub.pem"
	require.NoError(t, os.WriteFile(p, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), 0600))
	t.Setenv("AUTHZ_JWKS_PUBLIC_KEY_PATH", p)
	t.Setenv("JWT_KEY_ID", "kid-1")

	resp, err := h.GetJWKS(context.Background(), &authzservicev1.GetJWKSRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Keys, 1)
	require.Equal(t, "kid-1", resp.Keys[0].Kid)
}
