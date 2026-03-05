package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAuthZService_Live_RoleAndListPageSizeBranches(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, false)

	roleName := "live_ps_" + uuid.New().String()[:8]
	createResp, err := svc.CreateRole(ctx, &authzservicev1.CreateRoleRequest{
		Name:        roleName,
		Portal:      authzentityv1.Portal_PORTAL_SYSTEM,
		Description: "page-size role",
		CreatedBy:   uuid.New().String(),
	})
	require.NoError(t, err)
	roleID := createResp.Role.RoleId
	t.Cleanup(func() { _ = dbConn.Exec(`DELETE FROM authz_schema.roles WHERE role_id = ?`, roleID).Error })

	// canceled context -> error path in CreateRole
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = svc.CreateRole(cancelCtx, &authzservicev1.CreateRoleRequest{
		Name:      "cancel_role_" + uuid.New().String()[:8],
		Portal:    authzentityv1.Portal_PORTAL_SYSTEM,
		CreatedBy: uuid.New().String(),
	})
	require.Error(t, err)

	_, err = svc.UpdateRole(ctx, &authzservicev1.UpdateRoleRequest{
		RoleId: roleID,
		Name:   roleName + "_updated",
	})
	require.NoError(t, err)

	// page size normalization branches
	_, err = svc.ListRoles(ctx, &authzservicev1.ListRolesRequest{
		Portal:     authzentityv1.Portal_PORTAL_SYSTEM,
		ActiveOnly: false,
		PageSize:   500,
	})
	require.NoError(t, err)

	_, err = svc.ListRoles(ctx, &authzservicev1.ListRolesRequest{
		Portal:     authzentityv1.Portal_PORTAL_SYSTEM,
		ActiveOnly: false,
		PageSize:   1,
	})
	require.NoError(t, err)
}

func TestAuthZService_Live_AssignWithExpiryAndPolicyDenyBranch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, false)

	userID := uuid.New().String()
	domain := genSvcDomain("expiry")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	roleResp, err := svc.CreateRole(ctx, &authzservicev1.CreateRoleRequest{
		Name:      "live_exp_" + uuid.New().String()[:8],
		Portal:    authzentityv1.Portal_PORTAL_SYSTEM,
		CreatedBy: uuid.New().String(),
	})
	require.NoError(t, err)
	roleID := roleResp.Role.RoleId
	t.Cleanup(func() { _ = dbConn.Exec(`DELETE FROM authz_schema.roles WHERE role_id = ?`, roleID).Error })

	exp := timestamppb.New(time.Now().Add(2 * time.Hour))
	assignResp, err := svc.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		Domain:     domain,
		AssignedBy: uuid.New().String(),
		ExpiresAt:  exp,
	})
	require.NoError(t, err)
	require.NotNil(t, assignResp.UserRole.ExpiresAt)

	denyPol, err := svc.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{
		Subject:   "role:" + roleResp.Role.Name,
		Domain:    domain,
		Object:    "svc:claims/approve",
		Action:    "POST",
		Effect:    authzentityv1.PolicyEffect_POLICY_EFFECT_DENY,
		CreatedBy: uuid.New().String(),
	})
	require.NoError(t, err)
	require.Equal(t, authzentityv1.PolicyEffect_POLICY_EFFECT_DENY, denyPol.Policy.Effect)

	perms, err := svc.GetUserPermissions(ctx, &authzservicev1.GetUserPermissionsRequest{
		UserId: userID,
		Domain: domain,
	})
	require.NoError(t, err)
	require.NotEmpty(t, perms.Permissions)
}

func TestAuthZService_Live_UpdatePolicyDeactivateBranch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, false)

	userID := uuid.New().String()
	domain := genSvcDomain("deact")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	createResp, err := svc.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{
		Subject:   "user:" + userID,
		Domain:    domain,
		Object:    "svc:policy/export",
		Action:    "GET",
		Effect:    authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		CreatedBy: uuid.New().String(),
	})
	require.NoError(t, err)

	checkBefore, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID, Domain: domain, Object: "svc:policy/export", Action: "GET",
	})
	require.NoError(t, err)
	require.True(t, checkBefore.Allowed)

	// deactivate -> UpdatePolicyRule path where no re-add occurs
	_, err = svc.UpdatePolicyRule(ctx, &authzservicev1.UpdatePolicyRuleRequest{
		PolicyId: createResp.Policy.PolicyId,
		IsActive: false,
	})
	require.NoError(t, err)

	checkAfter, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID, Domain: domain, Object: "svc:policy/export", Action: "GET",
	})
	require.NoError(t, err)
	require.False(t, checkAfter.Allowed)
}
