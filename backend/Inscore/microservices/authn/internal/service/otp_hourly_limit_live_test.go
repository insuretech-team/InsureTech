package service

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestOTPService_CheckRateLimit_Hourly_LiveDB verifies the 3/hour cap for a recipient.
func TestOTPService_CheckRateLimit_Hourly_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testServiceLiveDB(t)
	tx := dbConn.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { _ = tx.Rollback().Error })

	otpRepo := repository.NewOTPRepository(tx.Table("authn_schema.otps"))
	cfg := &config.Config{}
	cfg.Security.RateLimitPerMinute = 100
	cfg.Security.RateLimitPerDay = 1000

	svc := NewOTPService(otpRepo, nil, nil, cfg, nil)

	recipient := "8801999999999"
	userID := uuid.New().String()
	numStr := strconv.FormatInt(time.Now().UnixNano()%1_000_000_000, 10)
	mobile := "+8801" + strings.Repeat("0", 9-len(numStr)) + numStr
	require.NoError(t, tx.Exec(
		`INSERT INTO authn_schema.users
		   (user_id, mobile_number, password_hash, status, user_type, created_at, updated_at)
		 VALUES (?, ?, 'test-hash', 'USER_STATUS_ACTIVE', 'USER_TYPE_B2C_CUSTOMER', NOW(), NOW())`,
		userID, mobile,
	).Error)

	for i := 0; i < 3; i++ {
		otp := &authnentityv1.OTP{
			OtpId:      uuid.New().String(),
			UserId:     userID,
			Purpose:    "login",
			DeviceType: "MOBILE_ANDROID",
			IpAddress:  "127.0.0.1",
			Recipient:  recipient,
			Channel:    "sms",
			OtpHash:    "hash",
			ExpiresAt:  timestamppb.New(time.Now().Add(5 * time.Minute)),
			Verified:   false,
			Attempts:   0,
			DlrStatus:  "PENDING",
		}
		require.NoError(t, otpRepo.Create(ctx, otp))
	}

	err := svc.checkRateLimit(ctx, recipient, "login")
	require.Error(t, err)
	require.Contains(t, err.Error(), "last hour")
}
