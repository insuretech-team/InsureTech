package service

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	_ "github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	svcTestDBOnce sync.Once
	svcTestDB     *gorm.DB
	svcTestDBErr  error
)

func testServiceLiveDB(t *testing.T) *gorm.DB {
	t.Helper()

	svcTestDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())

		if err := env.Load(); err != nil {
			logger.Warnf("Warning: couldn't load .env: %v", err)
		}

		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			// Relative to this package: backend/inscore/microservices/authn/internal/service
			configPath = "../../../../database.yaml"
		}

		svcTestDBErr = db.InitializeManagerForService(configPath)
		if svcTestDBErr != nil {
			return
		}

		schema.RegisterSerializer("proto_timestamp", db.ProtoTimestampSerializer{})
		svcTestDB = db.GetDB()
		if svcTestDB != nil {
			svcTestDB = svcTestDB.Debug()
		}
	})

	if svcTestDBErr != nil {
		t.Skipf("skipping live DB test: failed to init DB: %v", svcTestDBErr)
	}
	if svcTestDB == nil {
		t.Skip("skipping live DB test: DB is nil")
	}

	return svcTestDB
}

// TestTokenService_ValidateJWT_AfterRevocation_LiveDB verifies that after a
// session is revoked via sessionRepo.Revoke(), ValidateJWT (stateful with
// revocation check via ValidateJWTStrict) correctly returns Valid: false.
func TestTokenService_ValidateJWT_AfterRevocation_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testServiceLiveDB(t)

	// ── Setup user ────────────────────────────────────────────────────────────
	userID := uuid.New().String()

	// Insert a minimal user row so the FK on sessions is satisfied.
	// Live DB enforces chk_users_mobile_number (+8801XXXXXXXXX, digits only).
	mobileNumber := "+8801" + time.Now().Format("150405000")
	err := dbConn.Exec(
		`INSERT INTO authn_schema.users
		   (user_id, mobile_number, password_hash, status, user_type, created_at, updated_at)
		 VALUES (?, ?, 'test-hash', 'USER_STATUS_ACTIVE', 'USER_TYPE_B2C_CUSTOMER', NOW(), NOW())
		 ON CONFLICT DO NOTHING`,
		userID, mobileNumber,
	).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = dbConn.Exec(`DELETE FROM authn_schema.sessions WHERE user_id = ?`, userID).Error
		_ = dbConn.Exec(`DELETE FROM authn_schema.users   WHERE user_id = ?`, userID).Error
	})

	// ── Build repos + service ─────────────────────────────────────────────────
	sessionRepo := repository.NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	cfg := &config.Config{}
	cfg.JWT.Issuer = "insuretech-test"
	cfg.JWT.AccessTokenDuration = 15 * time.Minute
	cfg.JWT.RefreshTokenDuration = 7 * 24 * time.Hour

	generateTempRSAKeys(cfg)
	svc, err := NewTokenService(sessionRepo, nil, cfg, nil, nil)
	require.NoError(t, err)

	// ── Generate token pair ───────────────────────────────────────────────────
	tokenPair, err := svc.GenerateJWT(
		ctx,
		userID,
		"B2C_CUSTOMER",
		"",
		"",
		authnentityv1.DeviceType_DEVICE_TYPE_API,
		"127.0.0.1",
		"test",
	)
	require.NoError(t, err)
	require.NotEmpty(t, tokenPair.AccessToken)
	require.NotEmpty(t, tokenPair.SessionID)

	// ── Verify token is valid before revocation ───────────────────────────────
	resp, err := svc.ValidateJWT(ctx, tokenPair.AccessToken)
	require.NoError(t, err)
	require.True(t, resp.Valid, "expected token to be valid before revocation")

	// Strict validation should also pass before revocation.
	strictResp, err := svc.ValidateJWTStrict(ctx, tokenPair.AccessToken)
	require.NoError(t, err)
	require.True(t, strictResp.Valid, "expected strict validation to pass before revocation")

	// ── Revoke the session ────────────────────────────────────────────────────
	require.NoError(t, sessionRepo.Revoke(ctx, tokenPair.SessionID))

	// ── ValidateJWT (stateless) still returns Valid=true (expected: stateless) ─
	// The stateless ValidateJWT only checks signature + expiry; revocation is
	// handled by ValidateJWTStrict (which checks the DB). Confirm strict path
	// returns Valid=false after revocation.
	strictAfter, err := svc.ValidateJWTStrict(ctx, tokenPair.AccessToken)
	require.NoError(t, err)
	require.False(t, strictAfter.Valid, "expected ValidateJWTStrict to return Valid=false after revocation")
}
