package repository

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRoleRepo_LiveDB_DeleteAndMissingPaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewRoleRepo(dbConn)

	roleID := uuid.New().String()
	roleName := newLiveID("edge_role")
	t.Cleanup(func() { cleanupRoleByID(t, dbConn, roleID) })

	_, err := repo.Create(ctx, &authzentityv1.Role{
		RoleId:      roleID,
		Name:        roleName,
		Portal:      authzentityv1.Portal_PORTAL_SYSTEM,
		Description: "edge",
		IsActive:    true,
		CreatedBy:   uuid.New().String(),
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	})
	require.NoError(t, err)

	require.NoError(t, repo.Delete(ctx, roleID))

	_, err = repo.GetByID(ctx, roleID)
	require.Error(t, err)
	_, err = repo.GetByNameAndPortal(ctx, roleName, authzentityv1.Portal_PORTAL_SYSTEM)
	require.Error(t, err)
	_, err = repo.GetByName(ctx, authzentityv1.Portal_PORTAL_SYSTEM.String(), roleName)
	require.Error(t, err)
}

func TestRoleRepo_LiveDB_DuplicateAndContextErrorPaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewRoleRepo(dbConn)

	roleID := uuid.New().String()
	name := newLiveID("dup_role")
	role := &authzentityv1.Role{
		RoleId:      roleID,
		Name:        name,
		Portal:      authzentityv1.Portal_PORTAL_SYSTEM,
		Description: "dup-test",
		IsActive:    true,
		CreatedBy:   uuid.New().String(),
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	}
	t.Cleanup(func() { cleanupRoleByID(t, dbConn, roleID) })

	_, err := repo.Create(ctx, role)
	require.NoError(t, err)
	_, err = repo.Create(ctx, role)
	require.Error(t, err)

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	_, err = repo.Update(canceled, role)
	require.Error(t, err)
	err = repo.Delete(canceled, roleID)
	require.Error(t, err)

	all, err := repo.List(ctx, authzentityv1.Portal_PORTAL_UNSPECIFIED, false, 10, 0)
	require.NoError(t, err)
	require.NotNil(t, all)
}

func TestPolicyRepo_LiveDB_DeleteAndMissingPaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewPolicyRepo(dbConn)

	policyID := uuid.New().String()
	domain := "system:" + newLiveID("edge")
	t.Cleanup(func() { cleanupPolicyByID(t, dbConn, policyID) })

	_, err := repo.Create(ctx, &authzentityv1.PolicyRule{
		PolicyId:    policyID,
		Subject:     "role:" + newLiveID("edge"),
		Domain:      domain,
		Object:      "svc:edge/test",
		Action:      "GET",
		Effect:      authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		Description: "edge",
		IsActive:    true,
		CreatedBy:   uuid.New().String(),
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	})
	require.NoError(t, err)

	require.NoError(t, repo.Delete(ctx, policyID))
	_, err = repo.GetByID(ctx, policyID)
	require.Error(t, err)

	listAll, err := repo.List(ctx, "", false, 10, 0)
	require.NoError(t, err)
	require.NotNil(t, listAll)
}

func TestPolicyRepo_LiveDB_DuplicateAndContextErrorPaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewPolicyRepo(dbConn)

	policyID := uuid.New().String()
	domain := "system:" + newLiveID("dup")
	p := &authzentityv1.PolicyRule{
		PolicyId:    policyID,
		Subject:     "role:" + newLiveID("dup"),
		Domain:      domain,
		Object:      "svc:dup/test",
		Action:      "GET",
		Effect:      authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		Description: "dup",
		IsActive:    true,
		CreatedBy:   uuid.New().String(),
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	}
	t.Cleanup(func() { cleanupPolicyByID(t, dbConn, policyID) })

	_, err := repo.Create(ctx, p)
	require.NoError(t, err)
	_, err = repo.Create(ctx, p)
	require.Error(t, err)

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	_, err = repo.Update(canceled, p)
	require.Error(t, err)
	err = repo.Delete(canceled, policyID)
	require.Error(t, err)
	_, err = repo.List(canceled, domain, true, 10, 0)
	require.Error(t, err)
}

func TestPortalRepo_LiveDB_EdgePaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewPortalRepo(dbConn)

	_, err := repo.Upsert(ctx, nil)
	require.Error(t, err)

	// existing portal upsert with valid UUID updated_by branch
	_, err = repo.Upsert(ctx, &authzentityv1.PortalConfig{
		Portal:                  authzentityv1.Portal_PORTAL_B2C,
		MfaRequired:             false,
		MfaMethods:              []string{"sms_otp"},
		AccessTokenTtlSeconds:   900,
		RefreshTokenTtlSeconds:  604800,
		SessionTtlSeconds:       28800,
		IdleTimeoutSeconds:      1800,
		AllowConcurrentSessions: true,
		MaxConcurrentSessions:   3,
		UpdatedBy:               uuid.New().String(),
	})
	require.NoError(t, err)

	_, err = repo.GetByPortal(ctx, authzentityv1.Portal_PORTAL_UNSPECIFIED)
	require.NoError(t, err)

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	listCanceled, err := repo.List(canceled)
	if err == nil {
		require.NotNil(t, listCanceled)
	}
}

func TestAuditRepo_LiveDB_PageSizeBranchesAndCreateError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewAuditRepo(dbConn)

	userID := uuid.New().String()
	insertAuthnUserMinimal(t, dbConn, userID)
	t.Cleanup(func() { cleanupAuthnUserByID(t, dbConn, userID) })

	auditID := uuid.New().String()
	t.Cleanup(func() { cleanupAuditByID(t, dbConn, auditID) })
	require.NoError(t, repo.Create(ctx, &authzentityv1.AccessDecisionAudit{
		AuditId:   auditID,
		UserId:    userID,
		SessionId: uuid.New().String(),
		Domain:    "system:" + newLiveID("audit"),
		Subject:   "user:" + userID,
		Object:    "svc:audit/list",
		Action:    "GET",
		Decision:  authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		IpAddress: "127.0.0.1",
		UserAgent: "edge-test",
	}))

	rows1, _, err := repo.List(ctx, &authzservicev1.ListAccessDecisionAuditsRequest{
		UserId:   userID,
		PageSize: 0, // default branch
	})
	require.NoError(t, err)
	require.NotEmpty(t, rows1)

	rows2, _, err := repo.List(ctx, &authzservicev1.ListAccessDecisionAuditsRequest{
		UserId:   userID,
		PageSize: 1000, // >500 cap branch
	})
	require.NoError(t, err)
	require.NotEmpty(t, rows2)

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	err = repo.Create(canceled, &authzentityv1.AccessDecisionAudit{
		AuditId:   uuid.New().String(),
		UserId:    userID,
		SessionId: uuid.New().String(),
		Domain:    "system:" + newLiveID("audit"),
		Subject:   "user:" + userID,
		Object:    "svc:audit/create",
		Action:    "POST",
		Decision:  authzentityv1.PolicyEffect_POLICY_EFFECT_DENY,
		IpAddress: "127.0.0.1",
		UserAgent: "edge-test",
	})
	require.Error(t, err)
}

func TestTokenConfigRepo_LiveDB_DuplicateCreateAndNoActiveInTx(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewTokenConfigRepo(dbConn)

	kid := newLiveID("dup_kid")
	t.Cleanup(func() { cleanupTokenConfigByKID(t, dbConn, kid) })

	_, err := repo.Create(ctx, &authzentityv1.TokenConfig{
		Kid:           kid,
		Algorithm:     "RS256",
		PublicKeyPem:  "-----BEGIN PUBLIC KEY-----\nDUP\n-----END PUBLIC KEY-----",
		PrivateKeyRef: "secret/authz/" + kid,
		IsActive:      false,
		CreatedAt:     timestamppb.Now(),
	})
	require.NoError(t, err)

	_, err = repo.Create(ctx, &authzentityv1.TokenConfig{
		Kid:           kid,
		Algorithm:     "RS256",
		PublicKeyPem:  "-----BEGIN PUBLIC KEY-----\nDUP2\n-----END PUBLIC KEY-----",
		PrivateKeyRef: "secret/authz/" + kid,
		IsActive:      false,
		CreatedAt:     timestamppb.Now(),
	})
	require.Error(t, err)

	tx := dbConn.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { _ = tx.Rollback().Error })
	txRepo := NewTokenConfigRepo(tx)
	require.NoError(t, tx.Exec(`UPDATE authz_schema.token_configs SET is_active = false`).Error)
	_, err = txRepo.GetActive(ctx)
	require.Error(t, err)

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	_, err = repo.List(canceled)
	require.Error(t, err)
}

func TestUserRoleRepo_LiveDB_ConcurrencyAndErrorBranches(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewUserRoleRepo(dbConn)
	roleRepo := NewRoleRepo(dbConn)

	userID := uuid.New().String()
	roleID := uuid.New().String()
	domain := "agent:" + newLiveID("concurrent")

	insertAuthnUserMinimal(t, dbConn, userID)
	t.Cleanup(func() { cleanupAuthnUserByID(t, dbConn, userID) })
	_, err := roleRepo.Create(ctx, &authzentityv1.Role{
		RoleId:      roleID,
		Name:        newLiveID("con_role"),
		Portal:      authzentityv1.Portal_PORTAL_AGENT,
		Description: "concurrency role",
		IsActive:    true,
		CreatedBy:   uuid.New().String(),
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	})
	require.NoError(t, err)
	t.Cleanup(func() { cleanupRoleByID(t, dbConn, roleID) })

	const n = 8
	var wg sync.WaitGroup
	errCh := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, e := repo.Assign(ctx, &authzentityv1.UserRole{
				UserRoleId: uuid.New().String(),
				UserId:     userID,
				RoleId:     roleID,
				Domain:     domain,
				AssignedBy: uuid.New().String(),
				AssignedAt: timestamppb.Now(),
				ExpiresAt:  timestamppb.New(time.Now().Add(time.Duration(i+1) * time.Hour)),
			})
			errCh <- e
		}(i)
	}
	wg.Wait()
	close(errCh)
	for e := range errCh {
		require.NoError(t, e)
	}

	listByUserAllDomains, err := repo.ListByUser(ctx, userID, "")
	require.NoError(t, err)
	require.NotEmpty(t, listByUserAllDomains)

	listByRoleAllDomains, err := repo.ListByRole(ctx, roleID, "")
	require.NoError(t, err)
	require.NotEmpty(t, listByRoleAllDomains)

	// FK error branch
	_, err = repo.Assign(ctx, &authzentityv1.UserRole{
		UserRoleId: uuid.New().String(),
		UserId:     uuid.New().String(),
		RoleId:     uuid.New().String(),
		Domain:     "agent:" + newLiveID("bad"),
		AssignedBy: uuid.New().String(),
		AssignedAt: timestamppb.Now(),
	})
	require.Error(t, err)

	// context canceled error branch in revoke
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	err = repo.Revoke(canceled, userID, roleID, domain)
	require.Error(t, err)
}

func TestCasbinRuleRepo_LiveDB_ErrorBranches(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewCasbinRuleRepo(dbConn)
	domain := "agent:" + newLiveID("edgecasbin")
	sub := "role:" + newLiveID("edge")

	t.Cleanup(func() { cleanupCasbinByDomainPrefix(t, dbConn, domain) })

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	_, err := repo.Upsert(canceled, &authzentityv1.CasbinRule{
		Ptype: "p", V0: sub, V1: domain, V2: "svc:x", V3: "GET", V4: "allow",
	})
	require.Error(t, err)

	err = repo.Delete(canceled, &authzentityv1.CasbinRule{
		Ptype: "p", V0: sub, V1: domain, V2: "svc:x", V3: "GET",
	})
	require.Error(t, err)

	_, err = repo.ListByDomain(canceled, domain)
	require.Error(t, err)

	err = repo.DeleteByDomainAndSubject(canceled, domain, sub)
	require.Error(t, err)
}
