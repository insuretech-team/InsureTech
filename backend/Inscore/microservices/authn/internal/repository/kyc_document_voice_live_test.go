package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	kycv1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	voicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/voice/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDocumentType_UserDocument_UserProfile_Voice_KYC_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("live DB")
	}
	ctx := context.Background()
	db := testAuthnDB(t)

	// create a user for FK tables
	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, db, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, db, mobile, "") })
	insertUserMinimal(t, db, userID, mobile, "combo_test@example.com", "hash", 0)
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, db, "", userID) })

	// document type
	dtRepo := NewDocumentTypeRepository(db)
	dtID := uuid.New().String()
	_ = dtRepo.Delete(ctx, dtID)
	t.Cleanup(func() { _ = dtRepo.Delete(ctx, dtID) })
	require.NoError(t, dtRepo.Create(ctx, &authnv1.DocumentType{DocumentTypeId: dtID, Code: "TEST_" + uuid.New().String()[:6], Name: "Test", Description: "", IsActive: true}))
	_, err := dtRepo.GetByID(ctx, dtID)
	require.NoError(t, err)

	// user document
	udRepo := NewUserDocumentRepository(db)
	udID := uuid.New().String()
	_ = udRepo.Delete(ctx, udID)
	t.Cleanup(func() { _ = udRepo.Delete(ctx, udID) })
	require.NoError(t, udRepo.Create(ctx, &authnv1.UserDocument{UserDocumentId: udID, UserId: userID, DocumentTypeId: dtID, FileUrl: "s3://x", VerificationStatus: "PENDING"}))
	list, err := udRepo.ListByUser(ctx, userID)
	require.NoError(t, err)
	require.NotEmpty(t, list)

	// user profile
	upRepo := NewUserProfileRepository(db)
	_ = upRepo.DeleteByUserID(ctx, userID)
	t.Cleanup(func() { _ = upRepo.DeleteByUserID(ctx, userID) })
	require.NoError(t, upRepo.Create(ctx, &authnv1.UserProfile{UserId: userID, FullName: "Test User", Gender: authnv1.Gender_GENDER_MALE, DateOfBirth: timestamppb.New(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)), AddressLine1: "addr", City: "city", District: "dist", Division: "div", Country: "BD", NidNumber: genValidNID(), KycVerified: false, ConsentPrivacyAcceptance: true}))
	_, err = upRepo.GetByUserID(ctx, userID)
	require.NoError(t, err)

	// voice session
	vsRepo := NewVoiceSessionRepository(db)
	vsID := uuid.New().String()
	_ = vsRepo.Delete(ctx, vsID)
	t.Cleanup(func() { _ = vsRepo.Delete(ctx, vsID) })
	require.NoError(t, vsRepo.Create(ctx, &voicev1.VoiceSession{Id: vsID, SessionId: "ext-" + uuid.New().String(), UserId: userID, PhoneNumber: mobile, Language: "bn", Status: voicev1.SessionStatus_SESSION_STATUS_ACTIVE, Context: "{}"}))
	_, err = vsRepo.GetByID(ctx, vsID)
	require.NoError(t, err)

	// kyc + document verification
	kycRepo := NewKYCVerificationRepository(db)
	docRepo := NewDocumentVerificationRepository(db)
	kycID := uuid.New().String()
	_ = kycRepo.Delete(ctx, kycID)
	t.Cleanup(func() { _ = kycRepo.Delete(ctx, kycID) })
	// cleanup doc verifications via FK cascade, but be explicit
	_, _ = docRepo.DeleteByKYC(ctx, kycID)
	t.Cleanup(func() { _, _ = docRepo.DeleteByKYC(ctx, kycID) })

	require.NoError(t, kycRepo.Create(ctx, &kycv1.KYCVerification{Id: kycID, Type: kycv1.VerificationType_VERIFICATION_TYPE_KYC, EntityType: "USER", EntityId: userID, Method: kycv1.VerificationMethod_VERIFICATION_METHOD_MANUAL, Status: kycv1.VerificationStatus_VERIFICATION_STATUS_PENDING}))
	_, err = kycRepo.GetByID(ctx, kycID)
	require.NoError(t, err)

	docID := uuid.New().String()
	require.NoError(t, docRepo.Create(ctx, &kycv1.DocumentVerification{Id: docID, KycVerificationId: kycID, DocumentType: kycv1.DocumentType_DOCUMENT_TYPE_NID, DocumentNumber: "N123", Status: kycv1.DocumentStatus_DOCUMENT_STATUS_PENDING, ConfidenceScore: 0.5}))
	items, err := docRepo.ListByKYC(ctx, kycID, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	// mark verified
	now := time.Now()
	require.NoError(t, kycRepo.MarkVerified(ctx, kycID, userID, now, nil))
}
