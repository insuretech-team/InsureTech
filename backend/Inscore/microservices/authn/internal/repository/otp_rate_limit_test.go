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

// TestOTPRepository_CountRecentOTPs_RateLimit_LiveDB verifies that
// CountRecentOTPs correctly counts only OTPs within the time window,
// which is the core behaviour relied upon by the rate limiter in resend_otp.go.
func TestOTPRepository_CountRecentOTPs_RateLimit_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthnDB(t)

	// ── Setup user ────────────────────────────────────────────────────────────
	userID := uuid.New().String()
	mobile := genValidMobile()

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	insertUserMinimal(t, dbConn, userID, mobile, userID+"@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	// Cleanup all OTPs by recipient at the end.
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.otps").Where("recipient = ?", mobile).Delete(map[string]any{}).Error
	})

	repo := NewOTPRepository(dbConn.Table("authn_schema.otps"))

	// ── Insert 3 recent OTPs ──────────────────────────────────────────────────
	for i := 0; i < 3; i++ {
		otp := &authnentityv1.OTP{
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
		require.NoError(t, repo.Create(ctx, otp))
	}

	// ── Count within the 10-minute window — expect 3 ─────────────────────────
	count, err := repo.CountRecentOTPs(ctx, mobile, time.Now().Add(-10*time.Minute))
	require.NoError(t, err)
	require.Equal(t, int64(3), count)

	// ── Insert 1 OTP with created_at and expires_at 30 minutes ago (outside window) ──
	oldOtpID := uuid.New().String()
	oldTime := time.Now().Add(-30 * time.Minute)

	// Insert directly via raw SQL so we can set created_at in the past,
	// bypassing the repo.Create() which stamps created_at = NOW().
	err = dbConn.Exec(
		`INSERT INTO authn_schema.otps
		 (otp_id, user_id, otp_hash, purpose, channel, recipient, device_type, ip_address, attempts, verified, created_at, expires_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		oldOtpID, userID, "hash", "login", "sms", mobile, "WEB", "127.0.0.1", 0, false,
		oldTime, oldTime,
	).Error
	require.NoError(t, err)

	// ── Count again — old OTP must NOT be included ────────────────────────────
	count2, err := repo.CountRecentOTPs(ctx, mobile, time.Now().Add(-10*time.Minute))
	require.NoError(t, err)
	require.Equal(t, int64(3), count2, "old OTP outside the window should not be counted")
}
