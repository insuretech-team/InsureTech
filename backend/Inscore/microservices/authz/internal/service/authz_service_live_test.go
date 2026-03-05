package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/enforcer"
	authzevents "github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	authzSvcDBOnce sync.Once
	authzSvcDB     *gorm.DB
	authzSvcDBErr  error
)

func testAuthzServiceDB(t *testing.T) *gorm.DB {
	t.Helper()

	authzSvcDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}
		authzSvcDBErr = db.InitializeManagerForService(configPath)
		if authzSvcDBErr != nil {
			return
		}
		authzSvcDB = db.GetDB()
	})

	if authzSvcDBErr != nil {
		t.Skipf("skipping live DB test: %v", authzSvcDBErr)
	}
	if authzSvcDB == nil {
		t.Skip("skipping live DB test: db is nil")
	}
	return authzSvcDB
}

func newLiveAuthzService(t *testing.T, auditAll bool) (*AuthZService, *gorm.DB) {
	t.Helper()
	dbConn := testAuthzServiceDB(t)
	enf, err := enforcer.New(dbConn, "")
	require.NoError(t, err)
	svc := New(
		enf,
		repository.NewRoleRepo(dbConn),
		repository.NewUserRoleRepo(dbConn),
		repository.NewPolicyRepo(dbConn),
		repository.NewPortalRepo(dbConn),
		repository.NewAuditRepo(dbConn),
		auditAll,
		nil,
	)
	return svc, dbConn
}

func newLiveAuthzServiceWithPublisher(t *testing.T, auditAll bool) (*AuthZService, *gorm.DB) {
	t.Helper()
	dbConn := testAuthzServiceDB(t)
	enf, err := enforcer.New(dbConn, "")
	require.NoError(t, err)
	pub := authzevents.NewPublisher(nil) // real publisher with nil producer: non-blocking, no mock
	svc := New(
		enf,
		repository.NewRoleRepo(dbConn),
		repository.NewUserRoleRepo(dbConn),
		repository.NewPolicyRepo(dbConn),
		repository.NewPortalRepo(dbConn),
		repository.NewAuditRepo(dbConn),
		auditAll,
		pub,
	)
	return svc, dbConn
}

func genSvcDomain(prefix string) string {
	return fmt.Sprintf("system:%s_%d", prefix, time.Now().UnixNano())
}

func genSvcMobile() string {
	n := time.Now().UnixNano() % 1_000_000_000
	return fmt.Sprintf("+8801%09d", n)
}

func insertAuthnUserForService(t *testing.T, dbConn *gorm.DB, userID string) {
	t.Helper()
	err := dbConn.Exec(
		`INSERT INTO authn_schema.users
		   (user_id, mobile_number, password_hash, status, user_type, created_at, updated_at)
		 VALUES (?, ?, 'test-hash', 'USER_STATUS_ACTIVE', 'USER_TYPE_B2C_CUSTOMER', NOW(), NOW())`,
		userID, genSvcMobile(),
	).Error
	require.NoError(t, err)
}

func cleanupAuthnUserForService(t *testing.T, dbConn *gorm.DB, userID string) {
	t.Helper()
	_ = dbConn.Exec(`DELETE FROM authz_schema.user_roles WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authz_schema.access_decision_audits WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.sessions WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.otps WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.users WHERE user_id = ?`, userID).Error
}

func cleanupDomainArtifacts(t *testing.T, dbConn *gorm.DB, domain string) {
	t.Helper()
	_ = dbConn.Exec(`DELETE FROM authz_schema.access_decision_audits WHERE domain = ?`, domain).Error
	_ = dbConn.Exec(`DELETE FROM authz_schema.user_roles WHERE domain = ?`, domain).Error
	_ = dbConn.Exec(`DELETE FROM authz_schema.policy_rules WHERE domain = ?`, domain).Error
	_ = dbConn.Exec(`DELETE FROM authz_schema.casbin_rules WHERE v1 = ?`, domain).Error
}

func TestAuthZService_Live_CheckAccessAndBatch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, true)

	userID := uuid.New().String()
	domain := genSvcDomain("check")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	createdBy := uuid.New().String()
	polResp, err := svc.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{
		Subject:     "user:" + userID,
		Domain:      domain,
		Object:      "svc:claim/view",
		Action:      "GET",
		Effect:      authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		Description: "direct user allow",
		CreatedBy:   createdBy,
	})
	require.NoError(t, err)
	require.NotNil(t, polResp.Policy)

	checkAllow, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID,
		Domain: domain,
		Object: "svc:claim/view",
		Action: "GET",
		Context: &authzservicev1.AccessContext{
			IpAddress: "127.0.0.1",
			UserAgent: "svc-live-test",
			SessionId: uuid.New().String(),
		},
	})
	require.NoError(t, err)
	require.True(t, checkAllow.Allowed)
	require.Equal(t, authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW, checkAllow.Effect)

	checkDeny, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID,
		Domain: domain,
		Object: "svc:claim/delete",
		Action: "DELETE",
		Context: &authzservicev1.AccessContext{
			IpAddress: "127.0.0.1",
			UserAgent: "svc-live-test",
			SessionId: uuid.New().String(),
		},
	})
	require.NoError(t, err)
	require.False(t, checkDeny.Allowed)
	require.Equal(t, authzentityv1.PolicyEffect_POLICY_EFFECT_DENY, checkDeny.Effect)
	require.NotEmpty(t, checkDeny.Reason)

	batchResp, err := svc.BatchCheckAccess(ctx, &authzservicev1.BatchCheckAccessRequest{
		UserId: userID,
		Domain: domain,
		Checks: []*authzservicev1.AccessCheckTuple{
			{Object: "svc:claim/view", Action: "GET"},
			{Object: "svc:claim/delete", Action: "DELETE"},
		},
	})
	require.NoError(t, err)
	require.Len(t, batchResp.Results, 2)
	require.True(t, batchResp.Results[0].Allowed)
	require.False(t, batchResp.Results[1].Allowed)

	audits, err := svc.ListAccessDecisionAudits(ctx, &authzservicev1.ListAccessDecisionAuditsRequest{
		UserId: userID,
		Domain: domain,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, audits.TotalCount, int32(2))
}

func TestAuthZService_Live_RoleAssignPermissionsLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, false)

	userID := uuid.New().String()
	domain := genSvcDomain("role")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	createdBy := uuid.New().String()
	createRoleResp, err := svc.CreateRole(ctx, &authzservicev1.CreateRoleRequest{
		Name:        "live_role_" + uuid.New().String()[:8],
		Portal:      authzentityv1.Portal_PORTAL_SYSTEM,
		Description: "role for integration",
		CreatedBy:   createdBy,
	})
	require.NoError(t, err)
	role := createRoleResp.Role
	require.NotNil(t, role)
	t.Cleanup(func() {
		_ = dbConn.Exec(`DELETE FROM authz_schema.user_roles WHERE role_id = ?`, role.RoleId).Error
		_ = dbConn.Exec(`DELETE FROM authz_schema.roles WHERE role_id = ?`, role.RoleId).Error
	})

	getRoleResp, err := svc.GetRole(ctx, &authzservicev1.GetRoleRequest{RoleId: role.RoleId})
	require.NoError(t, err)
	require.Equal(t, role.RoleId, getRoleResp.Role.RoleId)

	updateRoleResp, err := svc.UpdateRole(ctx, &authzservicev1.UpdateRoleRequest{
		RoleId:      role.RoleId,
		Description: "updated role description",
	})
	require.NoError(t, err)
	require.Equal(t, "updated role description", updateRoleResp.Role.Description)

	_, err = svc.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{
		Subject:     "role:" + role.Name,
		Domain:      domain,
		Object:      "svc:policy/list",
		Action:      "GET",
		Effect:      authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		Description: "role allow list",
		CreatedBy:   createdBy,
	})
	require.NoError(t, err)

	assignResp, err := svc.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     role.RoleId,
		Domain:     domain,
		AssignedBy: createdBy,
		ExpiresAt:  nil,
	})
	require.NoError(t, err)
	require.Equal(t, userID, assignResp.UserRole.UserId)

	userRolesResp, err := svc.ListUserRoles(ctx, &authzservicev1.ListUserRolesRequest{
		UserId: userID,
		Domain: domain,
	})
	require.NoError(t, err)
	require.NotEmpty(t, userRolesResp.UserRoles)

	permsResp, err := svc.GetUserPermissions(ctx, &authzservicev1.GetUserPermissionsRequest{
		UserId: userID,
		Domain: domain,
	})
	require.NoError(t, err)
	require.NotEmpty(t, permsResp.Permissions)

	accessResp, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID,
		Domain: domain,
		Object: "svc:policy/list",
		Action: "GET",
	})
	require.NoError(t, err)
	require.True(t, accessResp.Allowed)

	_, err = svc.RemoveRole(ctx, &authzservicev1.RemoveRoleRequest{
		UserId: userID,
		RoleId: role.RoleId,
		Domain: domain,
	})
	require.NoError(t, err)

	postRemoveResp, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID,
		Domain: domain,
		Object: "svc:policy/list",
		Action: "GET",
	})
	require.NoError(t, err)
	require.False(t, postRemoveResp.Allowed)

	_, err = svc.DeleteRole(ctx, &authzservicev1.DeleteRoleRequest{RoleId: role.RoleId})
	require.NoError(t, err)

	listRolesResp, err := svc.ListRoles(ctx, &authzservicev1.ListRolesRequest{
		Portal:     authzentityv1.Portal_PORTAL_SYSTEM,
		ActiveOnly: false,
		PageSize:   0,
	})
	require.NoError(t, err)
	require.NotNil(t, listRolesResp.Roles)
}

func TestAuthZService_Live_PolicyPortalAndCache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLiveAuthzService(t, false)

	userID := uuid.New().String()
	domain := genSvcDomain("policy")
	insertAuthnUserForService(t, dbConn, userID)
	t.Cleanup(func() {
		cleanupDomainArtifacts(t, dbConn, domain)
		cleanupAuthnUserForService(t, dbConn, userID)
	})

	createdBy := uuid.New().String()
	createResp, err := svc.CreatePolicyRule(ctx, &authzservicev1.CreatePolicyRuleRequest{
		Subject:     "user:" + userID,
		Domain:      domain,
		Object:      "svc:claim/search",
		Action:      "GET",
		Effect:      authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		Description: "search allow",
		CreatedBy:   createdBy,
	})
	require.NoError(t, err)
	policyID := createResp.Policy.PolicyId

	_, err = svc.UpdatePolicyRule(ctx, &authzservicev1.UpdatePolicyRuleRequest{
		PolicyId: policyID,
		Action:   "POST",
		Effect:   authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		IsActive: true,
	})
	require.NoError(t, err)

	oldActionResp, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID, Domain: domain, Object: "svc:claim/search", Action: "GET",
	})
	require.NoError(t, err)
	require.False(t, oldActionResp.Allowed)

	newActionResp, err := svc.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
		UserId: userID, Domain: domain, Object: "svc:claim/search", Action: "POST",
	})
	require.NoError(t, err)
	require.True(t, newActionResp.Allowed)

	listPoliciesResp, err := svc.ListPolicyRules(ctx, &authzservicev1.ListPolicyRulesRequest{
		Domain:     domain,
		ActiveOnly: true,
		PageSize:   0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, listPoliciesResp.Policies)

	_, err = svc.DeletePolicyRule(ctx, &authzservicev1.DeletePolicyRuleRequest{PolicyId: policyID})
	require.NoError(t, err)

	updatePortalResp, err := svc.UpdatePortalConfig(ctx, &authzservicev1.UpdatePortalConfigRequest{
		Portal:                  authzentityv1.Portal_PORTAL_B2C,
		MfaRequired:             true,
		MfaMethods:              []string{"totp", "email_otp"},
		AccessTokenTtlSeconds:   1200,
		RefreshTokenTtlSeconds:  7200,
		SessionTtlSeconds:       3600,
		IdleTimeoutSeconds:      600,
		AllowConcurrentSessions: false,
		MaxConcurrentSessions:   1,
		UpdatedBy:               createdBy,
	})
	require.NoError(t, err)
	require.NotNil(t, updatePortalResp.Config)

	getPortalResp, err := svc.GetPortalConfig(ctx, &authzservicev1.GetPortalConfigRequest{
		Portal: authzentityv1.Portal_PORTAL_B2C,
	})
	require.NoError(t, err)
	require.True(t, getPortalResp.Config.MfaRequired)
	require.Equal(t, int32(1200), getPortalResp.Config.AccessTokenTtlSeconds)

	invalidated, err := svc.InvalidatePolicyCache(ctx, &authzservicev1.InvalidatePolicyCacheRequest{Domain: domain})
	require.NoError(t, err)
	require.True(t, invalidated.Invalidated)
}

func TestAuthZService_GetJWKS_RealKeyAndErrors(t *testing.T) {
	ctx := context.Background()
	svc := &AuthZService{}

	_, err := svc.GetJWKS(ctx, &authzservicev1.GetJWKSRequest{})
	require.Error(t, err, "default missing file should error")

	tmpBad, err := os.CreateTemp("", "jwks_bad_*.pem")
	require.NoError(t, err)
	defer os.Remove(tmpBad.Name())
	_, err = tmpBad.WriteString("not-a-pem")
	require.NoError(t, err)
	require.NoError(t, tmpBad.Close())
	t.Setenv("AUTHZ_JWKS_PUBLIC_KEY_PATH", tmpBad.Name())
	_, err = svc.GetJWKS(ctx, &authzservicev1.GetJWKSRequest{})
	require.Error(t, err)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pubASN1, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	require.NoError(t, err)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})
	tmpGood, err := os.CreateTemp("", "jwks_good_*.pem")
	require.NoError(t, err)
	defer os.Remove(tmpGood.Name())
	_, err = tmpGood.Write(pubPEM)
	require.NoError(t, err)
	require.NoError(t, tmpGood.Close())

	t.Setenv("AUTHZ_JWKS_PUBLIC_KEY_PATH", tmpGood.Name())
	t.Setenv("JWT_KEY_ID", "svc-live-kid")
	jwks, err := svc.GetJWKS(ctx, &authzservicev1.GetJWKSRequest{})
	require.NoError(t, err)
	require.Len(t, jwks.Keys, 1)
	require.Equal(t, "RSA", jwks.Keys[0].Kty)
	require.Equal(t, "RS256", jwks.Keys[0].Alg)
	require.Equal(t, "svc-live-kid", jwks.Keys[0].Kid)
	require.NotEmpty(t, jwks.Keys[0].N)
	require.NotEmpty(t, jwks.Keys[0].E)
}

func TestAuthZService_HelperFunctions(t *testing.T) {
	t.Setenv("AUTHZ_HELPER_TEST_ENV", "x")
	require.Equal(t, "x", getEnvOrDefault("AUTHZ_HELPER_TEST_ENV", "d"))
	require.Equal(t, "d", getEnvOrDefault("AUTHZ_HELPER_TEST_ENV_MISSING", "d"))

	ip, ua, sid, tid, did := extractContext(nil)
	require.Empty(t, ip)
	require.Empty(t, ua)
	require.Empty(t, sid)
	require.Empty(t, tid)
	require.Empty(t, did)

	ip, ua, sid, tid, did = extractContext(&authzservicev1.AccessContext{
		IpAddress: "1.1.1.1", UserAgent: "ua", SessionId: "sid", TokenId: "tid", DeviceId: "did",
	})
	require.Equal(t, "1.1.1.1", ip)
	require.Equal(t, "ua", ua)
	require.Equal(t, "sid", sid)
	require.Equal(t, "tid", tid)
	require.Equal(t, "did", did)
}
