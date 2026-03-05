package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func TestUserProfileRepository_LiveDB_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthnDB(t)
	repo := NewUserProfileRepository(dbConn)

	// UserProfile requires a valid user due to FK constraint
	userID := uuid.New().String()
	mobile := genValidMobile()

	// Pre-cleanup any stale data
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	_ = repo.DeleteByUserID(ctx, userID)

	// Insert the required user row first
	insertUserMinimal(t, dbConn, userID, mobile, "profile_test@example.com", "hash_profile", 0)

	t.Cleanup(func() {
		_ = repo.DeleteByUserID(ctx, userID)
		cleanupAuthnUser(ctx, t, dbConn, "", userID)
	})

	// --- Create ---
	profile := &authnentityv1.UserProfile{
		UserId:      userID,
		FullName:    "Integration Test User",
		DateOfBirth: timestamppb.New(time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)),
		Gender:      authnentityv1.Gender_GENDER_MALE,
		City:        "Dhaka",
		Country:     "BD",
		NidNumber:   genValidNID(),
		KycVerified: false,
	}
	require.NoError(t, repo.Create(ctx, profile))

	// --- GetByUserID ---
	got, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, got.UserId)
	require.Equal(t, "Integration Test User", got.FullName)
	require.Equal(t, "Dhaka", got.City)
	require.Equal(t, "BD", got.Country)
	require.False(t, got.KycVerified)

	// --- SetKYCVerified(true) ---
	verifiedAt := time.Now().UTC()
	require.NoError(t, repo.SetKYCVerified(ctx, userID, true, &verifiedAt))

	// --- GetByUserID again - verify kyc_verified=true ---
	updated, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.True(t, updated.KycVerified)

	// --- DeleteByUserID ---
	require.NoError(t, repo.DeleteByUserID(ctx, userID))

	_, err = repo.GetByUserID(ctx, userID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
