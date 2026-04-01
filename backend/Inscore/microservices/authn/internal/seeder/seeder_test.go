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

func TestSeedB2bAdminUser_SkipAndValidation(t *testing.T) {
	err := SeedB2bAdminUser(context.Background(), nil)
	require.NoError(t, err)

	t.Setenv("B2B_ADMIN", "")
	t.Setenv("B2B_ADMIN_MOBILE", "")
	t.Setenv("B2B_ADMIN_PASSWORD", "")
	t.Setenv("B2B_ADMIN_PASSWARD", "")
	err = SeedB2bAdminUser(context.Background(), nil)
	require.NoError(t, err)

	t.Setenv("B2B_ADMIN", "b2b_admin@example.com")
	t.Setenv("B2B_ADMIN_MOBILE", "bad-mobile")
	t.Setenv("B2B_ADMIN_PASSWARD", "SeedPass!1")
	err = SeedB2bAdminUser(context.Background(), nil)
	require.Error(t, err)
}

func TestNormalizeAdminEmail(t *testing.T) {
	got, err := normalizeAdminEmail(` "Faruk.Hannan@LifePlusBD.com" `)
	require.NoError(t, err)
	require.Equal(t, "faruk.hannan@lifeplusbd.com", got)

	_, err = normalizeAdminEmail("not-an-email")
	require.Error(t, err)
}

func TestSeedB2bAdminUser_AcceptsPreferredPasswordEnv(t *testing.T) {
	t.Setenv("B2B_ADMIN", ` "b2b_admin@example.com" `)
	t.Setenv("B2B_ADMIN_MOBILE", "01712345678")
	t.Setenv("B2B_ADMIN_PASSWORD", ` "SeedPass!1" `)
	t.Setenv("B2B_ADMIN_PASSWARD", "")

	err := SeedB2bAdminUser(context.Background(), nil)
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

func TestSeedAdminUser_LiveDB_CreateAndPreserveExistingCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	dbConn := testSeederLiveDB(t)
	ctx := context.Background()
	userRepo := repository.NewUserRepository(dbConn)

	email := "seed_admin_" + fmt.Sprintf("%d", time.Now().UnixNano()) + "@example.com"
	rawMobile := fmt.Sprintf("017%08d", time.Now().UnixNano()%100000000)
	mobile := "+880" + rawMobile[1:]
	t.Setenv("ADMIN_EMAIL", email)
	t.Setenv("ADMIN_MOBILE", rawMobile)
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
	require.Equal(t, oldHash, u2.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(u2.PasswordHash), []byte("SeedPass!1")))
}

func TestSeedB2bAdminUser_LiveDB_CreateAndPreserveExistingCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	dbConn := testSeederLiveDB(t)
	ctx := context.Background()
	userRepo := repository.NewUserRepository(dbConn)

	email := "seed_b2b_admin_" + fmt.Sprintf("%d", time.Now().UnixNano()) + "@example.com"
	rawMobile := fmt.Sprintf("017%08d", time.Now().UnixNano()%100000000)
	mobile := "+880" + rawMobile[1:]
	t.Setenv("B2B_ADMIN", email)
	t.Setenv("B2B_ADMIN_MOBILE", rawMobile)
	t.Setenv("B2B_ADMIN_PASSWARD", "SeedPass!1")
	t.Cleanup(func() { cleanupSeededAdmin(t, dbConn, email) })

	require.NoError(t, SeedB2bAdminUser(ctx, dbConn))
	u, err := userRepo.GetByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, mobile, u.MobileNumber)
	require.True(t, u.EmailVerified)
	require.NotEmpty(t, u.PasswordHash)

	oldHash := u.PasswordHash
	t.Setenv("B2B_ADMIN_PASSWARD", "SeedPass!2")
	require.NoError(t, SeedB2bAdminUser(ctx, dbConn))
	u2, err := userRepo.GetByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, oldHash, u2.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(u2.PasswordHash), []byte("SeedPass!1")))
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
