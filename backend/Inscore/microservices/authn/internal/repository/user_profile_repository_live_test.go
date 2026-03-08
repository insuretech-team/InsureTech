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

	// --- Update (raw map path) ---
	profile.FullName = "Updated Name"
	profile.City = "Chittagong"
	profile.DateOfBirth = timestamppb.New(time.Date(1992, 3, 20, 0, 0, 0, 0, time.UTC))
	require.NoError(t, repo.Update(ctx, profile))

	afterUpdate, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, "Updated Name", afterUpdate.FullName)
	require.Equal(t, "Chittagong", afterUpdate.City)

	// --- DeleteByUserID ---
	require.NoError(t, repo.DeleteByUserID(ctx, userID))

	_, err = repo.GetByUserID(ctx, userID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

// TestUserProfileRepository_AutoCreate_NilDOB verifies that Create succeeds when
// DateOfBirth is nil (auto-create path from GetUserProfile/UpdateUserProfile).
// This tests the NOT NULL sentinel fallback introduced to fix the profile 404 bug.
func TestUserProfileRepository_AutoCreate_NilDOB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthnDB(t)
	repo := NewUserProfileRepository(dbConn)

	userID := uuid.New().String()
	mobile := genValidMobile()

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	_ = repo.DeleteByUserID(ctx, userID)
	insertUserMinimal(t, dbConn, userID, mobile, "nilbirth_test@example.com", "hash_nilbirth", 0)

	t.Cleanup(func() {
		_ = repo.DeleteByUserID(ctx, userID)
		cleanupAuthnUser(ctx, t, dbConn, "", userID)
	})

	// Simulate auto-create with nil DateOfBirth (as done by GetUserProfile on first visit).
	profile := &authnentityv1.UserProfile{
		UserId:      userID,
		DateOfBirth: nil, // intentionally nil — must use 1900-01-01 sentinel
	}
	err := repo.Create(ctx, profile)
	require.NoError(t, err, "Create with nil DateOfBirth must not violate NOT NULL constraint")

	got, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, got.UserId)
	// date_of_birth should be the 1900-01-01 sentinel — year must be ≤ 1970 so
	// the BFF renders the date field as blank.
	if got.DateOfBirth != nil && got.DateOfBirth.IsValid() {
		require.LessOrEqual(t, got.DateOfBirth.AsTime().Year(), 1970,
			"sentinel DOB year should be ≤ 1970 so BFF renders date field as blank")
	}
}

// TestUserProfileRepository_Update_NilDOBSkipped verifies that Update does NOT
// overwrite a real date_of_birth when the incoming proto has a nil/sentinel DOB.
func TestUserProfileRepository_Update_NilDOBSkipped(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthnDB(t)
	repo := NewUserProfileRepository(dbConn)

	userID := uuid.New().String()
	mobile := genValidMobile()

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	_ = repo.DeleteByUserID(ctx, userID)
	insertUserMinimal(t, dbConn, userID, mobile, "skipbirth_test@example.com", "hash_skipbirth", 0)

	t.Cleanup(func() {
		_ = repo.DeleteByUserID(ctx, userID)
		cleanupAuthnUser(ctx, t, dbConn, "", userID)
	})

	realDOB := time.Date(1988, 7, 4, 0, 0, 0, 0, time.UTC)
	profile := &authnentityv1.UserProfile{
		UserId:      userID,
		FullName:    "DOB Skip Test",
		DateOfBirth: timestamppb.New(realDOB),
		City:        "Sylhet",
		Country:     "BD",
	}
	require.NoError(t, repo.Create(ctx, profile))

	// Now update with nil DateOfBirth — real DOB must be preserved in DB.
	profile.FullName = "DOB Skip Test Updated"
	profile.DateOfBirth = nil
	require.NoError(t, repo.Update(ctx, profile))

	got, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, "DOB Skip Test Updated", got.FullName)
	// DateOfBirth should still be the original real value.
	if got.DateOfBirth != nil && got.DateOfBirth.IsValid() {
		require.Equal(t, realDOB.Year(), got.DateOfBirth.AsTime().Year(),
			"DateOfBirth must not be overwritten when update sends nil")
		require.Equal(t, int(realDOB.Month()), int(got.DateOfBirth.AsTime().Month()))
		require.Equal(t, realDOB.Day(), got.DateOfBirth.AsTime().Day())
	}
}
