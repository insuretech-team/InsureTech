package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPolicyRepo_LiveDB_CRUDAndList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewPolicyRepo(dbConn)

	policyID := uuid.New().String()
	domain := "system:" + newLiveID("tenant")

	t.Cleanup(func() { cleanupPolicyByID(t, dbConn, policyID) })

	created, err := repo.Create(ctx, &authzentityv1.PolicyRule{
		PolicyId:    policyID,
		Subject:     "role:" + newLiveID("role"),
		Domain:      domain,
		Object:      "svc:policy/create",
		Action:      "POST",
		Effect:      authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW,
		Condition:   "",
		Description: "initial policy",
		IsActive:    true,
		CreatedBy:   uuid.New().String(),
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	})
	require.NoError(t, err)
	require.Equal(t, policyID, created.PolicyId)

	byID, err := repo.GetByID(ctx, policyID)
	require.NoError(t, err)
	require.Equal(t, created.Subject, byID.Subject)
	require.Equal(t, domain, byID.Domain)
	require.Equal(t, created.Object, byID.Object)
	require.Equal(t, "POST", byID.Action)
	require.Equal(t, authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW, byID.Effect)
	require.Equal(t, created.Condition, byID.Condition)
	require.Equal(t, "initial policy", byID.Description)
	require.True(t, byID.IsActive)
	require.Equal(t, created.CreatedBy, byID.CreatedBy)
	require.NotNil(t, byID.CreatedAt)
	require.NotNil(t, byID.UpdatedAt)

	byID.Action = "PUT"
	byID.Description = "updated policy"
	updated, err := repo.Update(ctx, byID)
	require.NoError(t, err)
	require.Equal(t, "PUT", updated.Action)
	require.Equal(t, "updated policy", updated.Description)

	listActive, err := repo.List(ctx, domain, true, 50, 0)
	require.NoError(t, err)
	require.NotEmpty(t, listActive)
	for _, p := range listActive {
		require.True(t, p.IsActive)
	}

	require.NoError(t, repo.SoftDelete(ctx, policyID))
	afterDelete, err := repo.GetByID(ctx, policyID)
	require.NoError(t, err)
	require.False(t, afterDelete.IsActive)
}
