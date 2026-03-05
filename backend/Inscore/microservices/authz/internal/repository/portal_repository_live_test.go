package repository

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	_ "github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/seeder"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	authzLiveDBOnce sync.Once
	authzLiveDB     *gorm.DB
	authzLiveDBErr  error
)

func testAuthzDB(t *testing.T) *gorm.DB {
	t.Helper()

	authzLiveDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()

		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}

		authzLiveDBErr = db.InitializeManagerForService(configPath)
		if authzLiveDBErr != nil {
			return
		}
		authzLiveDB = db.GetDB()
	})

	if authzLiveDBErr != nil {
		t.Skipf("skipping live DB test: %v", authzLiveDBErr)
	}
	if authzLiveDB == nil {
		t.Skip("skipping live DB test: db is nil")
	}
	return authzLiveDB
}

// TestPortalRepo_LiveDB_UpsertAndGet verifies portal config upsert/read path
// against the live DB for a concrete portal.
func TestPortalRepo_LiveDB_UpsertAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)
	tx := dbConn.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { _ = tx.Rollback().Error })

	portalRepo := NewPortalRepo(tx)

	in := &authzentityv1.PortalConfig{
		Portal:                  authzentityv1.Portal_PORTAL_SYSTEM,
		MfaRequired:             true,
		MfaMethods:              []string{"totp"},
		AccessTokenTtlSeconds:   900,
		RefreshTokenTtlSeconds:  604800,
		SessionTtlSeconds:       28800,
		IdleTimeoutSeconds:      1800,
		AllowConcurrentSessions: false,
		MaxConcurrentSessions:   1,
		UpdatedBy:               "live-test",
	}
	_, err := portalRepo.Upsert(ctx, in)
	require.NoError(t, err)

	out, err := portalRepo.GetByPortal(ctx, authzentityv1.Portal_PORTAL_SYSTEM)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, in.Portal, out.Portal)
	require.Equal(t, in.MfaRequired, out.MfaRequired)
	require.ElementsMatch(t, in.MfaMethods, out.MfaMethods)
	require.Equal(t, in.AccessTokenTtlSeconds, out.AccessTokenTtlSeconds)
	require.Equal(t, in.RefreshTokenTtlSeconds, out.RefreshTokenTtlSeconds)
	require.Equal(t, in.SessionTtlSeconds, out.SessionTtlSeconds)
	require.Equal(t, in.IdleTimeoutSeconds, out.IdleTimeoutSeconds)
	require.Equal(t, in.AllowConcurrentSessions, out.AllowConcurrentSessions)
	require.Equal(t, in.MaxConcurrentSessions, out.MaxConcurrentSessions)

	list, err := portalRepo.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, list)
}

// TestPortalSeeder_LiveDB_SeedPortalConfigs verifies all 6 portal defaults are
// upserted through seeder->repository and readable from DB.
func TestPortalSeeder_LiveDB_SeedPortalConfigs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)
	tx := dbConn.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { _ = tx.Rollback().Error })

	portalRepo := NewPortalRepo(tx)
	sd := seeder.New(nil, nil, nil, portalRepo, nil, nil, zap.NewNop())
	require.NoError(t, sd.SeedPortalConfigs(ctx))

	expected := map[authzentityv1.Portal]struct {
		required bool
		methods  []string
	}{
		authzentityv1.Portal_PORTAL_SYSTEM:    {required: true, methods: []string{"totp"}},
		authzentityv1.Portal_PORTAL_BUSINESS:  {required: true, methods: []string{"email_otp"}},
		authzentityv1.Portal_PORTAL_B2B:       {required: true, methods: []string{"totp", "email_otp"}},
		authzentityv1.Portal_PORTAL_AGENT:     {required: false, methods: []string{"sms_otp"}},
		authzentityv1.Portal_PORTAL_REGULATOR: {required: true, methods: []string{"totp"}},
		authzentityv1.Portal_PORTAL_B2C:       {required: false, methods: []string{"sms_otp"}},
	}

	for portal, exp := range expected {
		cfg, err := portalRepo.GetByPortal(ctx, portal)
		require.NoError(t, err, "portal lookup should succeed for %s", portal.String())
		require.NotNil(t, cfg)
		require.Equal(t, exp.required, cfg.MfaRequired, "mfa_required mismatch for %s", portal.String())
		require.ElementsMatch(t, exp.methods, cfg.MfaMethods, "mfa_methods mismatch for %s", portal.String())
	}
}
