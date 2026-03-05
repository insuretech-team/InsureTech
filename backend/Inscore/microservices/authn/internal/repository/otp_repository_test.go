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

func TestOTPRepository_CreateAndGet_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	// Use a constraint-safe mobile for live DB.
	mobile := "+8801" + time.Now().Format("150405") + "000" // +8801 + 9 digits
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "otp_test@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))
	otp := &authnentityv1.OTP{
		OtpId:      uuid.New().String(),
		UserId:     userID,
		OtpHash:    "hash", // repo stores hash, not raw code
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

	err := repo.Create(ctx, otp)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, otp.OtpId)
	require.NoError(t, err)
	require.Equal(t, otp.OtpId, got.OtpId)
	require.Equal(t, userID, got.UserId)
}
