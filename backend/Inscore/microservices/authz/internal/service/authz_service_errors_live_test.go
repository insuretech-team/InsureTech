package service

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/stretchr/testify/require"
)

func TestAuthZService_Live_CheckAccessErrorAndBatchErrorReason(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, false)
	userID := uuid.New().String()
	domain := genSvcDomain("errcheck")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	// Invalid regex in action pattern -> enforce error path.
	_, err := svc.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{
		Subject:   "user:" + userID,
		Domain:    domain,
		Object:    "svc:err/regex",
		Action:    "[",
		Effect:    authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		CreatedBy: uuid.New().String(),
	})
	require.NoError(t, err)

	_, err = svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID, Domain: domain, Object: "svc:err/regex", Action: "GET",
	})
	require.Error(t, err)

	batchResp, err := svc.BatchCheckAccess(ctx, &authzservicev1.BatchCheckAccessRequest{
		UserId: userID,
		Domain: domain,
		Checks: []*authzservicev1.AccessCheckTuple{
			{Object: "svc:err/regex", Action: "GET"},
		},
	})
	require.NoError(t, err)
	require.Len(t, batchResp.Results, 1)
	require.False(t, batchResp.Results[0].Allowed)
	require.NotEmpty(t, batchResp.Results[0].Reason)
}

func TestAuthZService_Live_RoleAssignmentErrorPaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, false)
	userID := uuid.New().String()
	domain := genSvcDomain("errrole")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	_, err := svc.GetRole(ctx, &authzservicev1.GetRoleRequest{RoleId: uuid.New().String()})
	require.Error(t, err)

	_, err = svc.UpdateRole(ctx, &authzservicev1.UpdateRoleRequest{RoleId: uuid.New().String(), Name: "x"})
	require.Error(t, err)

	_, err = svc.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId: userID, RoleId: uuid.New().String(), Domain: domain, AssignedBy: uuid.New().String(),
	})
	require.Error(t, err)

	createRoleResp, err := svc.CreateRole(ctx, &authzservicev1.CreateRoleRequest{
		Name: "err_role_" + uuid.New().String()[:8], Portal: authzentityv1.Portal_PORTAL_SYSTEM, CreatedBy: uuid.New().String(),
	})
	require.NoError(t, err)
	roleID := createRoleResp.Role.RoleId
	t.Cleanup(func() {
		_ = dbConn.Exec(`DELETE FROM authz_schema.user_roles WHERE role_id = ?`, roleID).Error
		_ = dbConn.Exec(`DELETE FROM authz_schema.roles WHERE role_id = ?`, roleID).Error
	})

	_, err = svc.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId: userID, RoleId: roleID, Domain: domain, AssignedBy: uuid.New().String(),
	})
	require.NoError(t, err)

	// Duplicate assignment should fail in real Casbin enforcer.
	_, err = svc.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId: userID, RoleId: roleID, Domain: domain, AssignedBy: uuid.New().String(),
	})
	require.Error(t, err)

	// No assignment for this user -> casbin remove path errors.
	_, err = svc.RemoveRole(ctx, &authzservicev1.RemoveRoleRequest{
		UserId: uuid.New().String(), RoleId: roleID, Domain: domain,
	})
	require.Error(t, err)
}

func TestAuthZService_Live_PolicyAndPortalErrorPaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, false)
	userID := uuid.New().String()
	domain := genSvcDomain("errpol")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	// Duplicate policy should fail at AddPolicy.
	req := &authzservicev1.CreatePolicyRuleRequest{
		Subject:   "user:" + userID,
		Domain:    domain,
		Object:    "svc:dup/policy",
		Action:    "GET",
		Effect:    authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		CreatedBy: uuid.New().String(),
	}
	_, err := svc.CreatePolicyRule(ctx, req)
	require.NoError(t, err)
	_, err = svc.CreatePolicyRule(ctx, req)
	require.Error(t, err)

	// Persist error branch: overlong domain should fail DB write.
	longDomain := "system:" + strings.Repeat("x", 260)
	_, err = svc.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{
		Subject:   "user:" + userID,
		Domain:    longDomain,
		Object:    "svc:oversize/domain",
		Action:    "GET",
		Effect:    authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		CreatedBy: uuid.New().String(),
	})
	require.Error(t, err)
	t.Cleanup(func() { _ = dbConn.Exec(`DELETE FROM authz_schema.casbin_rules WHERE v1 = ?`, longDomain).Error })

	_, err = svc.UpdatePolicyRule(ctx, &authzservicev1.UpdatePolicyRuleRequest{
		PolicyId: uuid.New().String(), Action: "POST", IsActive: true,
	})
	require.Error(t, err)

	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = svc.ListPolicyRules(cancelCtx, &authzservicev1.ListPolicyRulesRequest{Domain: domain, ActiveOnly: true})
	require.Error(t, err)

	_, err = svc.DeletePolicyRule(cancelCtx, &authzservicev1.DeletePolicyRuleRequest{PolicyId: uuid.New().String()})
	require.Error(t, err)

	_, err = svc.GetPortalConfig(ctx, &authzservicev1.GetPortalConfigRequest{Portal: authzentityv1.Portal_PORTAL_UNSPECIFIED})
	require.NoError(t, err)

	_, err = svc.UpdatePortalConfig(ctx, &authzservicev1.UpdatePortalConfigRequest{
		Portal: authzentityv1.Portal_PORTAL_UNSPECIFIED,
	})
	require.NoError(t, err)

	_, err = svc.ListAccessDecisionAudits(cancelCtx, &authzservicev1.ListAccessDecisionAuditsRequest{
		UserId: userID, Domain: domain,
	})
	require.Error(t, err)
}
