package seeder

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestSeedAdminUser_SkipAndNilDB(t *testing.T) {
	err := SeedAdminUser(context.Background(), nil)
	require.NoError(t, err)

	t.Setenv("ADMIN_EMAIL", "")
	t.Setenv("ADMIN_MOBILE", "")
	t.Setenv("ADMIN_PASSWORD", "")
	err = SeedAdminUser(context.Background(), nil)
	require.NoError(t, err)
}

func TestNormalizeAdminMobile(t *testing.T) {
	got, err := normalizeAdminMobile("01347-210751")
	require.NoError(t, err)
	require.Equal(t, "+8801347210751", got)

	_, err = normalizeAdminMobile("abc")
	require.Error(t, err)
}

func TestSeedDocumentTypes_NilDB(t *testing.T) {
	err := SeedDocumentTypes(context.Background(), nil)
	require.NoError(t, err)
}

var (
	seederDBOnce sync.Once
	seederDB     *gorm.DB
	seederDBErr  error
)

func testSeederLiveDB(t *testing.T) *gorm.DB {
	t.Helper()
	seederDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}
		seederDBErr = db.InitializeManagerForService(configPath)
		if seederDBErr != nil {
			return
		}
		schema.RegisterSerializer("proto_timestamp", db.ProtoTimestampSerializer{})
		seederDB = db.GetDB()
	})
	if seederDBErr != nil || seederDB == nil {
		t.Skipf("skipping live DB test: %v", seederDBErr)
	}
	return seederDB
}

func cleanupSeededAdmin(t *testing.T, dbConn *gorm.DB, email string) {
	t.Helper()
	_ = dbConn.Exec(`DELETE FROM authn_schema.sessions WHERE user_id IN (SELECT user_id FROM authn_schema.users WHERE email = ?)`, email).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.otps WHERE user_id IN (SELECT user_id FROM authn_schema.users WHERE email = ?)`, email).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.users WHERE email = ?`, email).Error
}

func TestSeedAdminUser_LiveDB_CreateAndUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	dbConn := testSeederLiveDB(t)
	ctx := context.Background()
	userRepo := repository.NewUserRepository(dbConn)

	email := "seed_admin_" + fmt.Sprintf("%d", time.Now().UnixNano()) + "@example.com"
	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	t.Setenv("ADMIN_EMAIL", email)
	t.Setenv("ADMIN_MOBILE", mobile)
	t.Setenv("ADMIN_PASSWORD", "SeedPass!1")
	t.Cleanup(func() { cleanupSeededAdmin(t, dbConn, email) })

	require.NoError(t, SeedAdminUser(ctx, dbConn))
	u, err := userRepo.GetByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, mobile, u.MobileNumber)
	require.True(t, u.EmailVerified)
	require.NotEmpty(t, u.PasswordHash)

	oldHash := u.PasswordHash
	t.Setenv("ADMIN_PASSWORD", "SeedPass!2")
	require.NoError(t, SeedAdminUser(ctx, dbConn))
	u2, err := userRepo.GetByEmail(ctx, email)
	require.NoError(t, err)
	require.NotEqual(t, oldHash, u2.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(u2.PasswordHash), []byte("SeedPass!2")))
}

func TestSeedDocumentTypes_LiveDB_Idempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	dbConn := testSeederLiveDB(t)
	ctx := context.Background()
	repo := repository.NewDocumentTypeRepository(dbConn)

	require.NoError(t, SeedDocumentTypes(ctx, dbConn))
	require.NoError(t, SeedDocumentTypes(ctx, dbConn))

	for _, code := range []string{"NID", "PASSPORT", "BIRTH_CERTIFICATE", "DRIVING_LICENSE", "TIN_CERTIFICATE"} {
		dt, err := repo.GetByCode(ctx, code)
		require.NoError(t, err)
		require.NotNil(t, dt)
		require.Equal(t, code, dt.Code)
	}
}
