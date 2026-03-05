package repository

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

func TestAuditRepo_LiveDB_CreateAndListFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewAuditRepo(dbConn)

	userID := uuid.New().String()
	insertAuthnUserMinimal(t, dbConn, userID)
	t.Cleanup(func() { cleanupAuthnUserByID(t, dbConn, userID) })

	domain := "agent:" + newLiveID("audit")
	auditID1 := uuid.New().String()
	auditID2 := uuid.New().String()
	t.Cleanup(func() {
		cleanupAuditByID(t, dbConn, auditID1)
		cleanupAuditByID(t, dbConn, auditID2)
	})

	err := repo.Create(ctx, &authzentityv1.AccessDecisionAudit{
		AuditId:     auditID1,
		UserId:      userID,
		SessionId:   uuid.New().String(),
		Domain:      domain,
		Subject:     "user:" + userID,
		Object:      "svc:policy/list",
		Action:      "GET",
		Decision:    authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		MatchedRule: "p, role:admin, " + domain + ", svc:policy/list, GET, allow",
		IpAddress:   "127.0.0.1",
		UserAgent:   "live-test",
		// DecidedAt intentionally omitted to verify defaulting in repository.
	})
	require.NoError(t, err)

	decidedAt2 := timestamppb.New(time.Now().Add(1 * time.Second))
	err = repo.Create(ctx, &authzentityv1.AccessDecisionAudit{
		AuditId:     auditID2,
		UserId:      userID,
		SessionId:   uuid.New().String(),
		Domain:      domain,
		Subject:     "user:" + userID,
		Object:      "svc:policy/create",
		Action:      "POST",
		Decision:    authzentityv1.PolicyEffect_POLICY_EFFECT_DENY,
		MatchedRule: "p, role:viewer, " + domain + ", svc:policy/create, POST, deny",
		IpAddress:   "127.0.0.1",
		UserAgent:   "live-test",
		DecidedAt:   decidedAt2,
	})
	require.NoError(t, err)

	allRows, totalAll, err := repo.List(ctx, &authzservicev1.ListAccessDecisionAuditsRequest{
		UserId: userID,
		Domain: domain,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, totalAll, int64(2))
	require.NotEmpty(t, allRows)

	denyRows, totalDeny, err := repo.List(ctx, &authzservicev1.ListAccessDecisionAuditsRequest{
		UserId:   userID,
		Domain:   domain,
		Decision: authzentityv1.PolicyEffect_POLICY_EFFECT_DENY,
		From:     timestamppb.New(time.Now().Add(-1 * time.Hour)),
		To:       timestamppb.New(time.Now().Add(1 * time.Hour)),
		PageSize: 10,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, totalDeny, int64(1))
	require.NotEmpty(t, denyRows)
	for _, row := range denyRows {
		require.Equal(t, authzentityv1.PolicyEffect_POLICY_EFFECT_DENY, row.Decision)
		require.NotNil(t, row.DecidedAt)
	}
}

