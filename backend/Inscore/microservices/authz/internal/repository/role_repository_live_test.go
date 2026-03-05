package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRoleRepo_LiveDB_CRUDAndList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewRoleRepo(dbConn)

	roleID := uuid.New().String()
	roleName := newLiveID("live_role")
	createdBy := uuid.New().String()
	now := time.Now()

	t.Cleanup(func() { cleanupRoleByID(t, dbConn, roleID) })

	in := &authzentityv1.Role{
		RoleId:      roleID,
		Name:        roleName,
		Portal:      authzentityv1.Portal_PORTAL_AGENT,
		Description: "initial role",
		IsSystem:    false,
		IsActive:    true,
		CreatedBy:   createdBy,
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}

	created, err := repo.Create(ctx, in)
	require.NoError(t, err)
	require.Equal(t, roleID, created.RoleId)

	byID, err := repo.GetByID(ctx, roleID)
	require.NoError(t, err)
	require.Equal(t, roleName, byID.Name)
	require.Equal(t, authzentityv1.Portal_PORTAL_AGENT, byID.Portal)
	require.Equal(t, "initial role", byID.Description)
	require.False(t, byID.IsSystem)
	require.True(t, byID.IsActive)
	require.Equal(t, createdBy, byID.CreatedBy)
	require.NotNil(t, byID.CreatedAt)
	require.NotNil(t, byID.UpdatedAt)

	byName, err := repo.GetByNameAndPortal(ctx, roleName, authzentityv1.Portal_PORTAL_AGENT)
	require.NoError(t, err)
	require.Equal(t, roleID, byName.RoleId)

	byNameRaw, err := repo.GetByName(ctx, authzentityv1.Portal_PORTAL_AGENT.String(), roleName)
	require.NoError(t, err)
	require.Equal(t, roleID, byNameRaw.RoleId)

	byID.Description = "updated role"
	byID.IsActive = true
	updated, err := repo.Update(ctx, byID)
	require.NoError(t, err)
	require.Equal(t, "updated role", updated.Description)

	listAll, err := repo.List(ctx, authzentityv1.Portal_PORTAL_AGENT, false, 50, 0)
	require.NoError(t, err)
	require.NotEmpty(t, listAll)

	require.NoError(t, repo.SoftDelete(ctx, roleID))

	afterDelete, err := repo.GetByID(ctx, roleID)
	require.NoError(t, err)
	require.False(t, afterDelete.IsActive)

	listActive, err := repo.List(ctx, authzentityv1.Portal_PORTAL_AGENT, true, 200, 0)
	require.NoError(t, err)
	for _, r := range listActive {
		require.True(t, r.IsActive)
	}
}
