package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ─── Create (returns *User) ───────────────────────────────────────────────────

func TestUserRepository_LiveDB_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	mobile := genValidMobile()
	email := "user_create_" + uuid.New().String()[:8] + "@example.com"

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })

	repo := NewUserRepository(dbConn.Table("authn_schema.users"))

	created, err := repo.Create(ctx, mobile, "hashed_password", email, authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION)
	require.NoError(t, err)
	require.NotEmpty(t, created.UserId, "Create should return a user with a generated ID")
	require.Equal(t, mobile, created.MobileNumber)
	require.Equal(t, email, created.Email)
	require.Equal(t, authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION, created.Status)
	require.NotNil(t, created.CreatedAt)
	require.NotNil(t, created.UpdatedAt)
	require.NotNil(t, created.WalletBalance)

	// Register cleanup by ID now that we know it.
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", created.UserId) })

	// Round-trip: confirm the row is actually in the DB.
	got, err := repo.GetByID(ctx, created.UserId)
	require.NoError(t, err)
	require.Equal(t, created.UserId, got.UserId)
	require.Equal(t, mobile, got.MobileNumber)
}

// ─── CreateFull ───────────────────────────────────────────────────────────────

func TestUserRepository_LiveDB_CreateFull(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	mobile := genValidMobile()
	email := "user_createfull_" + uuid.New().String()[:8] + "@example.com"
	userID := uuid.New().String()

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn.Table("authn_schema.users"))

	user := &authnentityv1.User{
		UserId:       userID,
		MobileNumber: mobile,
		Email:        email,
		PasswordHash: "full_hashed_pw",
		Status:       authnentityv1.UserStatus_USER_STATUS_ACTIVE,
		UserType:     authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER,
		WalletBalance: &commonv1.Money{
			Amount:   0,
			Currency: "BDT",
		},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	err := repo.CreateFull(ctx, user)
	require.NoError(t, err)
	require.NotEmpty(t, user.UserId, "CreateFull should keep or assign a UserId")

	// Confirm the row is retrievable.
	got, err := repo.GetByID(ctx, user.UserId)
	require.NoError(t, err)
	require.Equal(t, user.UserId, got.UserId)
	require.Equal(t, mobile, got.MobileNumber)
	require.Equal(t, email, got.Email)
	require.Equal(t, authnentityv1.UserStatus_USER_STATUS_ACTIVE, got.Status)
}

// ─── GetByMobileNumber ────────────────────────────────────────────────────────

func TestUserRepository_LiveDB_GetByMobileNumber(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "mobile_lookup@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn.Table("authn_schema.users"))

	got, err := repo.GetByMobileNumber(ctx, mobile)
	require.NoError(t, err)
	require.Equal(t, userID, got.UserId)
	require.Equal(t, mobile, got.MobileNumber)
}

// ─── GetByEmail ───────────────────────────────────────────────────────────────

func TestUserRepository_LiveDB_GetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	email := "email_lookup_" + uuid.New().String()[:8] + "@example.com"

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, email, "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn.Table("authn_schema.users"))

	got, err := repo.GetByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, userID, got.UserId)
	require.Equal(t, email, got.Email)
}

// ─── UpdateStatus ─────────────────────────────────────────────────────────────

func TestUserRepository_LiveDB_UpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "status_update@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn.Table("authn_schema.users"))

	// Change status to SUSPENDED.
	err := repo.UpdateStatus(ctx, userID, authnentityv1.UserStatus_USER_STATUS_SUSPENDED)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, authnentityv1.UserStatus_USER_STATUS_SUSPENDED, got.Status)

	// Change status back to ACTIVE.
	err = repo.UpdateStatus(ctx, userID, authnentityv1.UserStatus_USER_STATUS_ACTIVE)
	require.NoError(t, err)

	got2, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, authnentityv1.UserStatus_USER_STATUS_ACTIVE, got2.Status)
}

// ─── UpdateLastLogin ──────────────────────────────────────────────────────────

func TestUserRepository_LiveDB_UpdateLastLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "last_login@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	repo := NewUserRepository(dbConn.Table("authn_schema.users"))

	sessionType := "SESSION_TYPE_JWT"
	err := repo.UpdateLastLogin(ctx, userID, sessionType)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, got.LastLoginAt, "last_login_at should be set after UpdateLastLogin")
	require.Equal(t, sessionType, got.LastLoginSessionType)
}
