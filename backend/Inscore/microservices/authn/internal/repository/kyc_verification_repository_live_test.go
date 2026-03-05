package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	kycv1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestKYCVerificationRepository_LiveDB_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthnDB(t)
	repo := NewKYCVerificationRepository(dbConn)

	// Create a real user so verified_by FK is satisfied.
	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	insertUserMinimal(t, dbConn, userID, mobile, "kyc_test@example.com", "hash", 1)
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })

	kycID := uuid.New().String()
	entityID := userID // use the real user as entity
	entityType := "USER"

	// Pre-cleanup and deferred cleanup
	_ = repo.Delete(ctx, kycID)
	t.Cleanup(func() { _ = repo.Delete(ctx, kycID) })

	// --- Create ---
	kyc := &kycv1.KYCVerification{
		Id:         kycID,
		Type:       kycv1.VerificationType_VERIFICATION_TYPE_KYC,
		EntityType: entityType,
		EntityId:   entityID,
		Method:     kycv1.VerificationMethod_VERIFICATION_METHOD_MANUAL,
		Status:     kycv1.VerificationStatus_VERIFICATION_STATUS_PENDING,
	}
	require.NoError(t, repo.Create(ctx, kyc))

	// --- GetByID ---
	got, err := repo.GetByID(ctx, kycID)
	require.NoError(t, err)
	require.Equal(t, kycID, got.Id)
	require.Equal(t, entityType, got.EntityType)
	require.Equal(t, entityID, got.EntityId)

	// --- GetByEntity ---
	byEntity, err := repo.GetByEntity(ctx, entityType, entityID)
	require.NoError(t, err)
	require.Equal(t, kycID, byEntity.Id)

	// --- ListByStatus (PENDING) ---
	list, err := repo.ListByStatus(ctx, kycv1.VerificationStatus_VERIFICATION_STATUS_PENDING, 50, 0)
	require.NoError(t, err)
	found := false
	for _, item := range list {
		if item.Id == kycID {
			found = true
			break
		}
	}
	require.True(t, found, "created KYC record should appear in ListByStatus(PENDING)")

	// --- UpdateStatus (REJECTED with reason) ---
	reason := "document unclear"
	require.NoError(t, repo.UpdateStatus(ctx, kycID, kycv1.VerificationStatus_VERIFICATION_STATUS_REJECTED, &reason))

	rejected, err := repo.GetByID(ctx, kycID)
	require.NoError(t, err)
	require.Equal(t, kycv1.VerificationStatus_VERIFICATION_STATUS_REJECTED, rejected.Status)

	// --- MarkVerified ---
	verifiedAt := time.Now().UTC().Truncate(time.Second)
	expiresAt := verifiedAt.Add(365 * 24 * time.Hour)
	verifiedBy := userID // must be a real user_id due to FK on verified_by
	require.NoError(t, repo.MarkVerified(ctx, kycID, verifiedBy, verifiedAt, &expiresAt))

	verified, err := repo.GetByID(ctx, kycID)
	require.NoError(t, err)
	require.Equal(t, kycv1.VerificationStatus_VERIFICATION_STATUS_VERIFIED, verified.Status)

	// --- Delete ---
	require.NoError(t, repo.Delete(ctx, kycID))

	_, err = repo.GetByID(ctx, kycID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
