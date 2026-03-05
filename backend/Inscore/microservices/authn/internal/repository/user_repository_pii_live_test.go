package repository

import (
	"context"
	"testing"
	"time"

	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
)

const (
	testPIIAESKey  = "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
	testPIIHMACKey = "ffeeddccbbaa99887766554433221100ffeeddccbbaa99887766554433221100"
)

func TestPIIUserRepository_LiveDB_EncryptedCRUDAndPassThrough(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	t.Setenv("PII_AES_KEY", testPIIAESKey)
	t.Setenv("PII_HMAC_KEY", testPIIHMACKey)

	ctx := context.Background()
	dbConn := testAuthnDB(t)
	inner := NewUserRepository(dbConn)
	repo := NewPIIUserRepository(inner)
	require.NotNil(t, repo.enc)

	mobile := genValidMobile()
	email := "pii_live_" + time.Now().Format("150405000") + "@example.com"
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })

	created, err := repo.Create(ctx, mobile, "hash_pii", email, authnentityv1.UserStatus_USER_STATUS_ACTIVE)
	require.NoError(t, err)
	require.NotEmpty(t, created.UserId)
	require.Equal(t, mobile, created.MobileNumber)
	require.Equal(t, email, created.Email)
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", created.UserId) })

	var storedMobile, storedEmail, mobileIdx, emailIdx, bioEnc, bioIdx string
	err = dbConn.Raw(
		`SELECT mobile_number, email, COALESCE(mobile_number_idx,''), COALESCE(email_idx,''), COALESCE(biometric_token_enc,''), COALESCE(biometric_token_idx,'')
		   FROM authn_schema.users WHERE user_id = ?`,
		created.UserId,
	).Scan(&struct {
		A string
	}{}).Error
	require.NoError(t, err)
	row := dbConn.Raw(
		`SELECT mobile_number, email, COALESCE(mobile_number_idx,''), COALESCE(email_idx,''), COALESCE(biometric_token_enc,''), COALESCE(biometric_token_idx,'')
		   FROM authn_schema.users WHERE user_id = ?`,
		created.UserId,
	).Row()
	require.NoError(t, row.Scan(&storedMobile, &storedEmail, &mobileIdx, &emailIdx, &bioEnc, &bioIdx))
	require.NotEmpty(t, storedMobile)
	require.NotEmpty(t, storedEmail)
	require.NotEmpty(t, mobileIdx)
	require.NotEmpty(t, emailIdx)
	require.NotEmpty(t, bioEnc)
	require.NotEmpty(t, bioIdx)

	gotByID, err := repo.GetByID(ctx, created.UserId)
	require.NoError(t, err)
	require.Equal(t, mobile, gotByID.MobileNumber)
	require.Equal(t, email, gotByID.Email)

	gotByMobile, err := repo.GetByMobileNumber(ctx, mobile)
	require.NoError(t, err)
	require.Equal(t, created.UserId, gotByMobile.UserId)
	require.Equal(t, mobile, gotByMobile.MobileNumber)

	gotByEmail, err := repo.GetByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, created.UserId, gotByEmail.UserId)
	require.Equal(t, email, gotByEmail.Email)

	require.NoError(t, repo.UpdatePassword(ctx, created.UserId, "hash_new"))
	require.NoError(t, repo.UpdateEmailVerified(ctx, created.UserId))
	require.NoError(t, repo.UpdateStatus(ctx, created.UserId, authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	require.NoError(t, repo.UpdateLastLogin(ctx, created.UserId, "WEB"))
	count, err := repo.IncrementEmailLoginAttempts(ctx, created.UserId)
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, int32(1))
	require.NoError(t, repo.LockEmailAuth(ctx, created.UserId, 5*time.Minute))
	require.NoError(t, repo.ResetEmailLoginAttempts(ctx, created.UserId))
}

func TestPIIUserRepository_LiveDB_CreateFullAndFallbackMode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	// Fallback mode (keys unset) should keep plaintext behavior.
	t.Setenv("PII_AES_KEY", "")
	t.Setenv("PII_HMAC_KEY", "")
	innerFallback := NewUserRepository(dbConn)
	repoFallback := NewPIIUserRepository(innerFallback)
	require.Nil(t, repoFallback.enc)
	encVal, err := repoFallback.encryptField("abc")
	require.NoError(t, err)
	require.Equal(t, "abc", encVal)
	decVal, err := repoFallback.decryptField("xyz")
	require.NoError(t, err)
	require.Equal(t, "xyz", decVal)
	require.Equal(t, "plain", repoFallback.blindIndex("plain"))

	mobile1 := genValidMobile()
	email1 := "pii_fb_" + time.Now().Format("150405001") + "@example.com"
	cleanupAuthnUser(ctx, t, dbConn, mobile1, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile1, "") })
	u1, err := repoFallback.Create(ctx, mobile1, "hash_fb", email1, authnentityv1.UserStatus_USER_STATUS_ACTIVE)
	require.NoError(t, err)
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", u1.UserId) })
	gotFB, err := repoFallback.GetByMobileNumber(ctx, mobile1)
	require.NoError(t, err)
	require.Equal(t, u1.UserId, gotFB.UserId)

	// Encrypted mode with CreateFull + index-based lookup.
	t.Setenv("PII_AES_KEY", testPIIAESKey)
	t.Setenv("PII_HMAC_KEY", testPIIHMACKey)
	innerEnc := NewUserRepository(dbConn)
	repoEnc := NewPIIUserRepository(innerEnc)
	require.NotNil(t, repoEnc.enc)

	mobile2 := genValidMobile()
	email2 := "pii_full_" + time.Now().Format("150405002") + "@example.com"
	cleanupAuthnUser(ctx, t, dbConn, mobile2, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile2, "") })

	user := &authnentityv1.User{
		MobileNumber: mobile2,
		Email:        email2,
		PasswordHash: "hash_full",
		Status:       authnentityv1.UserStatus_USER_STATUS_ACTIVE,
		UserType:     authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER,
	}
	require.NoError(t, repoEnc.CreateFull(ctx, user))
	require.NotEmpty(t, user.UserId)
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", user.UserId) })
	require.NotEmpty(t, user.BiometricTokenEnc)

	got2ByMobile, err := repoEnc.GetByMobileNumber(ctx, mobile2)
	require.NoError(t, err)
	require.Equal(t, user.UserId, got2ByMobile.UserId)
	require.Equal(t, mobile2, got2ByMobile.MobileNumber)

	got2ByEmail, err := repoEnc.GetByEmail(ctx, email2)
	require.NoError(t, err)
	require.Equal(t, user.UserId, got2ByEmail.UserId)
	require.Equal(t, email2, got2ByEmail.Email)
}
