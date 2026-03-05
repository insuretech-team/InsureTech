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
	"google.golang.org/grpc/metadata"
)

// TestTokenService_ValidateJWT_DeviceBinding_LiveDB verifies x-device-id metadata
// must match ins_device claim during JWT validation.
func TestTokenService_ValidateJWT_DeviceBinding_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testServiceLiveDB(t)
	tx := dbConn.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { _ = tx.Rollback().Error })

	userID := uuid.New().String()
	numStr := strconv.FormatInt(time.Now().UnixNano()%1_000_000_000, 10)
	mobile := "+8801" + strings.Repeat("0", 9-len(numStr)) + numStr

	err := tx.Exec(
		`INSERT INTO authn_schema.users
		   (user_id, mobile_number, password_hash, status, user_type, created_at, updated_at)
		 VALUES (?, ?, 'test-hash', 'USER_STATUS_ACTIVE', 'USER_TYPE_B2C_CUSTOMER', NOW(), NOW())`,
		userID, mobile,
	).Error
	require.NoError(t, err)

	sessionRepo := repository.NewSessionRepository(tx.Table("authn_schema.sessions"))
	cfg := &config.Config{}
	cfg.JWT.Issuer = "insuretech-test"
	cfg.JWT.AccessTokenDuration = 15 * time.Minute
	cfg.JWT.RefreshTokenDuration = 7 * 24 * time.Hour
	generateTempRSAKeys(cfg)

	svc, err := NewTokenService(sessionRepo, nil, cfg, nil, nil)
	require.NoError(t, err)

	tokenPair, err := svc.GenerateJWT(
		ctx,
		userID,
		authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER.String(),
		"",
		"dev-abc-1",
		authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_ANDROID,
		"127.0.0.1",
		"test-agent",
	)
	require.NoError(t, err)

	matchCtx := metadata.NewIncomingContext(ctx, metadata.Pairs("x-device-id", "dev-abc-1"))
	matchResp, err := svc.ValidateJWT(matchCtx, tokenPair.AccessToken)
	require.NoError(t, err)
	require.True(t, matchResp.Valid, "expected token to validate with matching device id")

	mismatchCtx := metadata.NewIncomingContext(ctx, metadata.Pairs("x-device-id", "dev-other-9"))
	mismatchResp, err := svc.ValidateJWT(mismatchCtx, tokenPair.AccessToken)
	require.NoError(t, err)
	require.False(t, mismatchResp.Valid, "expected token validation to fail on device-id mismatch")
}
