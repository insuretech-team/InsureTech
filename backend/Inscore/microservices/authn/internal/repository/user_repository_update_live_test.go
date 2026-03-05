package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
)

// TestUserRepository_LiveDB_UpdatePassword verifies that UpdatePassword persists the new hash.
func TestUserRepository_LiveDB_UpdatePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	email := "upd_pw_" + uuid.New().String()[:8] + "@example.com"

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, email, "old_hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn)

	err := repo.UpdatePassword(ctx, userID, "new_hash")
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, "new_hash", got.PasswordHash)
}

// TestUserRepository_LiveDB_UpdateEmailVerified verifies that UpdateEmailVerified sets EmailVerified=true
// and populates EmailVerifiedAt.
func TestUserRepository_LiveDB_UpdateEmailVerified(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	email := "upd_ev_" + uuid.New().String()[:8] + "@example.com"

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, email, "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn)

	err := repo.UpdateEmailVerified(ctx, userID)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.True(t, got.EmailVerified, "EmailVerified should be true after UpdateEmailVerified")
	require.NotNil(t, got.EmailVerifiedAt, "EmailVerifiedAt should be set after UpdateEmailVerified")
}

// TestUserRepository_LiveDB_IncrementEmailLoginAttempts verifies that each call increments the counter
// and that the returned value reflects the new count.
func TestUserRepository_LiveDB_IncrementEmailLoginAttempts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	email := "upd_ila_" + uuid.New().String()[:8] + "@example.com"

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, email, "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn)

	count1, err := repo.IncrementEmailLoginAttempts(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, int32(1), count1)

	count2, err := repo.IncrementEmailLoginAttempts(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, int32(2), count2)

	count3, err := repo.IncrementEmailLoginAttempts(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, int32(3), count3)

	got, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, int32(3), got.EmailLoginAttempts)
}

// TestUserRepository_LiveDB_LockEmailAuth verifies that LockEmailAuth sets EmailLockedUntil
// to a time in the future.
func TestUserRepository_LiveDB_LockEmailAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	email := "upd_lock_" + uuid.New().String()[:8] + "@example.com"

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, email, "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn)

	err := repo.LockEmailAuth(ctx, userID, 30*time.Minute)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, got.EmailLockedUntil, "EmailLockedUntil should be set after LockEmailAuth")
	require.True(t, got.EmailLockedUntil.AsTime().After(time.Now()), "EmailLockedUntil should be in the future")
}

// TestUserRepository_LiveDB_ResetEmailLoginAttempts verifies that ResetEmailLoginAttempts zeroes
// the attempt counter and clears the lock.
func TestUserRepository_LiveDB_ResetEmailLoginAttempts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	email := "upd_reset_" + uuid.New().String()[:8] + "@example.com"

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, email, "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn)

	// Set up some state to reset.
	_, err := repo.IncrementEmailLoginAttempts(ctx, userID)
	require.NoError(t, err)

	err = repo.LockEmailAuth(ctx, userID, 30*time.Minute)
	require.NoError(t, err)

	// Now reset.
	err = repo.ResetEmailLoginAttempts(ctx, userID)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, int32(0), got.EmailLoginAttempts, "EmailLoginAttempts should be 0 after reset")
	// EmailLockedUntil should be nil or zero after reset.
	if got.EmailLockedUntil != nil {
		require.True(t, got.EmailLockedUntil.AsTime().IsZero(), "EmailLockedUntil should be zero after reset")
	}
}
