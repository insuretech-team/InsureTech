package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	kycv1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDocumentVerificationRepository_LiveDB_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthnDB(t)
	kycRepo := NewKYCVerificationRepository(dbConn)
	docRepo := NewDocumentVerificationRepository(dbConn)

	// Create a parent KYC verification record (FK requirement)
	kycID := uuid.New().String()
	entityID := uuid.New().String()

	_ = kycRepo.Delete(ctx, kycID)
	t.Cleanup(func() {
		_, _ = docRepo.DeleteByKYC(ctx, kycID)
		_ = kycRepo.Delete(ctx, kycID)
	})

	require.NoError(t, kycRepo.Create(ctx, &kycv1.KYCVerification{
		Id:         kycID,
		Type:       kycv1.VerificationType_VERIFICATION_TYPE_KYC,
		EntityType: "USER",
		EntityId:   entityID,
		Method:     kycv1.VerificationMethod_VERIFICATION_METHOD_MANUAL,
		Status:     kycv1.VerificationStatus_VERIFICATION_STATUS_PENDING,
	}))

	// --- Create DocumentVerification ---
	docID := uuid.New().String()
	_ = dbConn.Table("authn_schema.document_verifications").Where("doc_verification_id = ?", docID).Delete(map[string]any{}).Error
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.document_verifications").Where("doc_verification_id = ?", docID).Delete(map[string]any{}).Error
	})

	doc := &kycv1.DocumentVerification{
		Id:                docID,
		KycVerificationId: kycID,
		DocumentType:      kycv1.DocumentType_DOCUMENT_TYPE_NID,
		DocumentNumber:    "DOC-" + uuid.New().String()[:8],
		Status:            kycv1.DocumentStatus_DOCUMENT_STATUS_PENDING,
		ConfidenceScore:   0.75,
	}
	require.NoError(t, docRepo.Create(ctx, doc))

	// --- GetByID ---
	got, err := docRepo.GetByID(ctx, docID)
	require.NoError(t, err)
	require.Equal(t, docID, got.Id)
	require.Equal(t, kycID, got.KycVerificationId)
	require.Equal(t, kycv1.DocumentStatus_DOCUMENT_STATUS_PENDING, got.Status)

	// --- ListByKYC ---
	list, err := docRepo.ListByKYC(ctx, kycID, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, list)
	found := false
	for _, item := range list {
		if item.Id == docID {
			found = true
			break
		}
	}
	require.True(t, found, "created document verification should appear in ListByKYC")

	// --- UpdateStatus (VERIFIED with confidence) ---
	confidence := float32(0.98)
	require.NoError(t, docRepo.UpdateStatus(ctx, docID, kycv1.DocumentStatus_DOCUMENT_STATUS_VERIFIED, &confidence))

	updated, err := docRepo.GetByID(ctx, docID)
	require.NoError(t, err)
	require.Equal(t, kycv1.DocumentStatus_DOCUMENT_STATUS_VERIFIED, updated.Status)

	// --- Delete ---
	require.NoError(t, dbConn.Table("authn_schema.document_verifications").Where("doc_verification_id = ?", docID).Delete(map[string]any{}).Error)

	_, err = docRepo.GetByID(ctx, docID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
