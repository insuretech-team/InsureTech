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

func TestOTPRepository_LiveDB_Flow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	// Create user (minimal + cleanup)
	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "otp_repo_test@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	otpID := uuid.New().String()
	otp := &authnentityv1.OTP{
		OtpId:      otpID,
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
	require.NoError(t, repo.Create(ctx, otp))
	t.Cleanup(func() {
		// best-effort cleanup
		_ = dbConn.Table("authn_schema.otps").Where("otp_id = ?", otpID).Delete(map[string]any{}).Error
	})

	// Read
	got, err := repo.GetByID(ctx, otpID)
	require.NoError(t, err)
	require.Equal(t, otpID, got.OtpId)
	require.Equal(t, mobile, got.Recipient)

	// Update attempts
	require.NoError(t, repo.IncrementAttempts(ctx, otpID))
	got2, err := repo.GetByID(ctx, otpID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, got2.Attempts, int32(1))

	// Mark verified
	require.NoError(t, repo.MarkVerified(ctx, otpID))
	got3, err := repo.GetByID(ctx, otpID)
	require.NoError(t, err)
	require.True(t, got3.Verified)
}
