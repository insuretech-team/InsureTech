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

func TestUserRoleRepo_LiveDB_AssignUpsertListRemove(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)

	userID := uuid.New().String()
	roleID := uuid.New().String()
	domain := "agent:" + newLiveID("tenant")
	assignedBy1 := uuid.New().String()
	assignedBy2 := uuid.New().String()
	roleName := newLiveID("live_ur_role")

	insertAuthnUserMinimal(t, dbConn, userID)
	t.Cleanup(func() { cleanupAuthnUserByID(t, dbConn, userID) })

	roleRepo := NewRoleRepo(dbConn)
	_, err := roleRepo.Create(ctx, &authzentityv1.Role{
		RoleId:      roleID,
		Name:        roleName,
		Portal:      authzentityv1.Portal_PORTAL_AGENT,
		Description: "role for user-role tests",
		IsActive:    true,
		CreatedBy:   assignedBy1,
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	})
	require.NoError(t, err)
	t.Cleanup(func() { cleanupRoleByID(t, dbConn, roleID) })

	repo := NewUserRoleRepo(dbConn)
	urID := uuid.New().String()
	firstExpiry := timestamppb.New(time.Now().Add(24 * time.Hour))
	_, err = repo.Assign(ctx, &authzentityv1.UserRole{
		UserRoleId: urID,
		UserId:     userID,
		RoleId:     roleID,
		Domain:     domain,
		AssignedBy: assignedBy1,
		AssignedAt: timestamppb.Now(),
		ExpiresAt:  firstExpiry,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = dbConn.Exec(`DELETE FROM authz_schema.user_roles WHERE user_role_id = ?`, urID).Error
	})

	secondExpiry := timestamppb.New(time.Now().Add(48 * time.Hour))
	_, err = repo.Assign(ctx, &authzentityv1.UserRole{
		UserRoleId: uuid.New().String(),
		UserId:     userID,
		RoleId:     roleID,
		Domain:     domain,
		AssignedBy: assignedBy2,
		AssignedAt: timestamppb.Now(),
		ExpiresAt:  secondExpiry,
	})
	require.NoError(t, err)

	list, err := repo.ListByUser(ctx, userID, domain)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, assignedBy2, list[0].AssignedBy, "upsert should update assigned_by")
	require.NotNil(t, list[0].ExpiresAt)
	require.WithinDuration(t, secondExpiry.AsTime(), list[0].ExpiresAt.AsTime(), 2*time.Second)

	require.NoError(t, repo.Remove(ctx, userID, roleID, domain))
	listAfterRemove, err := repo.ListByUser(ctx, userID, domain)
	require.NoError(t, err)
	require.Len(t, listAfterRemove, 0)
}

func TestUserRoleRepo_LiveDB_ListByRoleAndRevoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewUserRoleRepo(dbConn)
	roleRepo := NewRoleRepo(dbConn)

	userID := uuid.New().String()
	roleID := uuid.New().String()
	roleName := newLiveID("live_revoke_role")
	domain := "business:" + newLiveID("tenant")
	assignedBy := uuid.New().String()

	insertAuthnUserMinimal(t, dbConn, userID)
	t.Cleanup(func() { cleanupAuthnUserByID(t, dbConn, userID) })
	_, err := roleRepo.Create(ctx, &authzentityv1.Role{
		RoleId:      roleID,
		Name:        roleName,
		Portal:      authzentityv1.Portal_PORTAL_BUSINESS,
		Description: "role for revoke tests",
		IsActive:    true,
		CreatedBy:   assignedBy,
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	})
	require.NoError(t, err)
	t.Cleanup(func() { cleanupRoleByID(t, dbConn, roleID) })

	assigned, err := repo.Assign(ctx, &authzentityv1.UserRole{
		UserRoleId: uuid.New().String(),
		UserId:     userID,
		RoleId:     roleID,
		Domain:     domain,
		AssignedBy: assignedBy,
		AssignedAt: timestamppb.Now(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = dbConn.Exec(`DELETE FROM authz_schema.user_roles WHERE user_role_id = ?`, assigned.UserRoleId).Error
	})

	listByRole, err := repo.ListByRole(ctx, roleID, domain)
	require.NoError(t, err)
	require.Len(t, listByRole, 1)
	require.Equal(t, userID, listByRole[0].UserId)

	require.NoError(t, repo.Revoke(ctx, userID, roleID, domain))
	listByRoleAfter, err := repo.ListByRole(ctx, roleID, domain)
	require.NoError(t, err)
	require.Len(t, listByRoleAfter, 0)
}

