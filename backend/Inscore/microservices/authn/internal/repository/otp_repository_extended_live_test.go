package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// buildOTP is a local helper to construct a minimal valid OTP for testing.
func buildOTP(userID, mobile string) *authnentityv1.OTP {
	return &authnentityv1.OTP{
		OtpId:      uuid.New().String(),
		UserId:     userID,
		OtpHash:    "hash",
		Purpose:    "login",
		Channel:    "sms",
		Recipient:  mobile,
		DeviceType: "WEB",
		IpAddress:  "127.0.0.1",
		Attempts:   0,
		Verified:   false,
		CreatedAt:  timestamppb.Now(),
		ExpiresAt:  timestamppb.New(time.Now().Add(5 * time.Minute)),
	}
}

// setupOTPTestUser creates a user and registers cleanup, returning userID and mobile.
func setupOTPTestUser(ctx context.Context, t *testing.T) (userID, mobile string) {
	t.Helper()
	dbConn := testAuthnDB(t)
	userID = uuid.New().String()
	mobile = genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	insertUserMinimal(t, dbConn, userID, mobile, userID+"@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })
	return userID, mobile
}

// registerOTPCleanup registers a t.Cleanup to delete the OTP row by ID.
func registerOTPCleanup(t *testing.T, otpID string) {
	t.Helper()
	dbConn := testAuthnDB(t)
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.otps").Where("otp_id = ?", otpID).Delete(map[string]any{}).Error
	})
}

// TestOTPRepository_GetByProviderMessageID_LiveDB verifies that an OTP can be
// retrieved by its provider_message_id after creation.
func TestOTPRepository_GetByProviderMessageID_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	providerMsgID := "msg-" + uuid.New().String()
	otp := buildOTP(userID, mobile)
	otp.ProviderMessageId = providerMsgID

	require.NoError(t, repo.Create(ctx, otp))
	registerOTPCleanup(t, otp.OtpId)

	got, err := repo.GetByProviderMessageID(ctx, providerMsgID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, otp.OtpId, got.OtpId)
	require.Equal(t, providerMsgID, got.ProviderMessageId)
}

// TestOTPRepository_ExpireOTP_LiveDB verifies that ExpireOTP sets expires_at
// to a time in the past (or at most now), making the OTP effectively expired.
func TestOTPRepository_ExpireOTP_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otp := buildOTP(userID, mobile)
	// Ensure it starts with a future expiry.
	otp.ExpiresAt = timestamppb.New(time.Now().Add(10 * time.Minute))

	require.NoError(t, repo.Create(ctx, otp))
	registerOTPCleanup(t, otp.OtpId)

	require.NoError(t, repo.ExpireOTP(ctx, otp.OtpId))

	got, err := repo.GetByID(ctx, otp.OtpId)
	require.NoError(t, err)
	require.NotNil(t, got)
	// expires_at must now be <= now (with 1-second tolerance for clock skew).
	require.True(t, got.ExpiresAt.AsTime().Before(time.Now().Add(time.Second)),
		"expected expires_at to be set to now or past, got %v", got.ExpiresAt.AsTime())
}

// TestOTPRepository_UpdateDLRStatus_LiveDB verifies that the DLR status and
// dlr_received_at are persisted correctly after UpdateDLRStatus.
func TestOTPRepository_UpdateDLRStatus_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otp := buildOTP(userID, mobile)
	otp.ProviderMessageId = "dlr-test-" + uuid.New().String()
	require.NoError(t, repo.Create(ctx, otp))
	registerOTPCleanup(t, otp.OtpId)

	require.NoError(t, repo.UpdateDLRStatus(ctx, otp.ProviderMessageId, "DELIVERED", ""))

	got, err := repo.GetByID(ctx, otp.OtpId)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "DELIVERED", got.DlrStatus)
}

// TestOTPRepository_CountRecentOTPs_LiveDB verifies that two OTPs created for
// the same recipient are counted correctly.
func TestOTPRepository_CountRecentOTPs_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otp1 := buildOTP(userID, mobile)
	otp2 := buildOTP(userID, mobile)

	require.NoError(t, repo.Create(ctx, otp1))
	registerOTPCleanup(t, otp1.OtpId)

	require.NoError(t, repo.Create(ctx, otp2))
	registerOTPCleanup(t, otp2.OtpId)

	since := time.Now().Add(-1 * time.Hour)
	count, err := repo.CountRecentOTPs(ctx, mobile, since)
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, int64(2))
}

// TestOTPRepository_GetLastOTP_LiveDB verifies that GetLastOTP returns the
// most recent OTP (one of the two created OTPs) for a recipient.
func TestOTPRepository_GetLastOTP_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otp1 := buildOTP(userID, mobile)
	otp2 := buildOTP(userID, mobile)

	require.NoError(t, repo.Create(ctx, otp1))
	registerOTPCleanup(t, otp1.OtpId)

	// Small sleep to ensure distinct created_at values.
	time.Sleep(5 * time.Millisecond)

	require.NoError(t, repo.Create(ctx, otp2))
	registerOTPCleanup(t, otp2.OtpId)

	got, err := repo.GetLastOTP(ctx, mobile)
	require.NoError(t, err)
	require.NotNil(t, got)

	validIDs := map[string]bool{otp1.OtpId: true, otp2.OtpId: true}
	require.True(t, validIDs[got.OtpId], "expected OtpId to be one of the two created OTPs, got %s", got.OtpId)
}

// TestOTPRepository_ListByRecipient_LiveDB creates two OTPs (one verified,
// one not) and verifies that the verified filter works correctly.
func TestOTPRepository_ListByRecipient_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otpVerified := buildOTP(userID, mobile)
	otpUnverified := buildOTP(userID, mobile)

	require.NoError(t, repo.Create(ctx, otpVerified))
	registerOTPCleanup(t, otpVerified.OtpId)

	require.NoError(t, repo.Create(ctx, otpUnverified))
	registerOTPCleanup(t, otpUnverified.OtpId)

	// Mark the first OTP as verified.
	require.NoError(t, repo.MarkVerified(ctx, otpVerified.OtpId))

	// Confirm verified flag is set.
	confirmed, err := repo.GetByID(ctx, otpVerified.OtpId)
	require.NoError(t, err)
	require.True(t, confirmed.Verified)

	// List only verified OTPs.
	trueVal := true
	verifiedList, err := repo.ListByRecipient(ctx, mobile, &trueVal, "")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(verifiedList), 1)
	for _, o := range verifiedList {
		require.True(t, o.Verified, "expected only verified OTPs in result")
	}

	// List only unverified OTPs.
	falseVal := false
	unverifiedList, err := repo.ListByRecipient(ctx, mobile, &falseVal, "")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(unverifiedList), 1)
	for _, o := range unverifiedList {
		require.False(t, o.Verified, "expected only unverified OTPs in result")
	}
}

// TestOTPRepository_CleanupExpiredOTPs_LiveDB creates an OTP with expires_at
// in the past and verifies it is removed by CleanupExpiredOTPs.
func TestOTPRepository_CleanupExpiredOTPs_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otp := buildOTP(userID, mobile)
	// Set expires_at in the past so it qualifies for cleanup.
	otp.ExpiresAt = timestamppb.New(time.Now().Add(-1 * time.Hour))

	require.NoError(t, repo.Create(ctx, otp))
	// Register cleanup in case the test fails before CleanupExpiredOTPs runs.
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.otps").Where("otp_id = ?", otp.OtpId).Delete(map[string]any{}).Error
	})

	rowsAffected, err := repo.CleanupExpiredOTPs(ctx, time.Now())
	require.NoError(t, err)
	require.GreaterOrEqual(t, rowsAffected, int64(1))
}

// TestOTPRepository_GetPendingDLRs_LiveDB creates an OTP with DlrStatus="PENDING"
// and channel="sms", then verifies it appears in GetPendingDLRs results.
func TestOTPRepository_GetPendingDLRs_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otp := buildOTP(userID, mobile)
	otp.DlrStatus = "PENDING"
	otp.Channel = "sms"

	require.NoError(t, repo.Create(ctx, otp))
	registerOTPCleanup(t, otp.OtpId)

	results, err := repo.GetPendingDLRs(ctx, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 1)

	found := false
	for _, r := range results {
		if r.OtpId == otp.OtpId {
			found = true
			break
		}
	}
	require.True(t, found, "expected created OTP with PENDING DLR status to appear in results")
}

// TestOTPRepository_GetStatsByCarrier_LiveDB creates an OTP with Carrier="GP"
// and verifies the stats map contains at least 1 for that carrier.
func TestOTPRepository_GetStatsByCarrier_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID, mobile := setupOTPTestUser(ctx, t)
	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otp := buildOTP(userID, mobile)
	otp.Carrier = "GP"
	otp.Channel = "sms"

	require.NoError(t, repo.Create(ctx, otp))
	registerOTPCleanup(t, otp.OtpId)

	since := time.Now().Add(-1 * time.Hour)
	stats, err := repo.GetStatsByCarrier(ctx, since)
	require.NoError(t, err)
	require.NotNil(t, stats)
	require.GreaterOrEqual(t, stats["GP"], int64(1))
}

// TestOTPRepository_GetDeliveryRate_LiveDB verifies that GetDeliveryRate
// returns a valid rate between 0 and 100 without error.
func TestOTPRepository_GetDeliveryRate_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	since := time.Now().Add(-24 * time.Hour)
	rate, err := repo.GetDeliveryRate(ctx, since)
	require.NoError(t, err)
	require.GreaterOrEqual(t, rate, float64(0))
	require.LessOrEqual(t, rate, float64(100))
}
