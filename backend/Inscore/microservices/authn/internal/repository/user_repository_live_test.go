package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_LiveDB_Updates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "user_repo_update@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn.Table("authn_schema.users"))

	require.NoError(t, repo.UpdatePassword(ctx, userID, "hash2"))
	u, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, u.UserId)

	require.NoError(t, repo.UpdateEmailVerified(ctx, userID))
	u2, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.True(t, u2.EmailVerified)

	_, _ = repo.IncrementEmailLoginAttempts(ctx, userID)
	require.NoError(t, repo.LockEmailAuth(ctx, userID, 5*time.Minute))
	require.NoError(t, repo.ResetEmailLoginAttempts(ctx, userID))
}
