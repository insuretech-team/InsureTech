package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/stretchr/testify/require"
)

func TestAuthZService_Live_PublisherBranches(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveAuthzServiceWithPublisher(t, false)

	userID := uuid.New().String()
	domain := genSvcDomain("pub")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	// CreateRole -> publisher branch
	roleResp, err := svc.CreateRole(ctx, &authzservicev1.CreateRoleRequest{
		Name:      "live_pub_" + uuid.New().String()[:8],
		Portal:    authzentityv1.Portal_PORTAL_SYSTEM,
		CreatedBy: uuid.New().String(),
	})
	require.NoError(t, err)
	roleID := roleResp.Role.RoleId
	t.Cleanup(func() { _ = dbConn.Exec(`DELETE FROM authz_schema.roles WHERE role_id = ?`, roleID).Error })

	// CreatePolicyRule -> publisher branch
	_, err = svc.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{
		Subject:   "role:" + roleResp.Role.Name,
		Domain:    domain,
		Object:    "svc:pub/allow",
		Action:    "GET",
		Effect:    authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		CreatedBy: uuid.New().String(),
	})
	require.NoError(t, err)

	// AssignRole -> publisher branch
	_, err = svc.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId: userID, RoleId: roleID, Domain: domain, AssignedBy: uuid.New().String(),
	})
	require.NoError(t, err)

	// CheckAccess deny with publisher branch for PublishAccessDenied
	denyResp, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID, Domain: domain, Object: "svc:pub/denied", Action: "POST",
		Context: &authzservicev1.AccessContext{
			IpAddress: "127.0.0.1", SessionId: uuid.New().String(),
		},
	})
	require.NoError(t, err)
	require.False(t, denyResp.Allowed)

	// Portal update -> publisher branch
	_, err = svc.UpdatePortalConfig(ctx, &authzservicev1.UpdatePortalConfigRequest{
		Portal:                  authzentityv1.Portal_PORTAL_B2C,
		MfaRequired:             false,
		MfaMethods:              []string{"sms_otp"},
		AccessTokenTtlSeconds:   1000,
		RefreshTokenTtlSeconds:  2000,
		SessionTtlSeconds:       3000,
		IdleTimeoutSeconds:      400,
		AllowConcurrentSessions: true,
		MaxConcurrentSessions:   2,
		UpdatedBy:               uuid.New().String(),
	})
	require.NoError(t, err)

	// Cache invalidate -> publisher branch
	_, err = svc.InvalidatePolicyCache(ctx, &authzservicev1.InvalidatePolicyCacheRequest{Domain: domain})
	require.NoError(t, err)
}

