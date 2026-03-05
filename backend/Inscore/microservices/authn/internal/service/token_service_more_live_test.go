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
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/stretchr/testify/require"
)

func newLiveTokenService(t *testing.T) (*TokenService, *repository.SessionRepository, *repository.UserRepository) {
	t.Helper()
	dbConn := testServiceLiveDB(t)
	sessionRepo := repository.NewSessionRepository(dbConn.Table("authn_schema.sessions"))
	userRepo := repository.NewUserRepository(dbConn.Table("authn_schema.users"))

	cfg := &config.Config{}
	cfg.JWT.Issuer = "insuretech-test"
	cfg.JWT.AccessTokenDuration = 15 * time.Minute
	cfg.JWT.RefreshTokenDuration = 7 * 24 * time.Hour
	cfg.Security.ServerSessionDuration = 2 * time.Hour
	cfg.Security.BCryptCost = 10

	generateTempRSAKeys(cfg)
	svc, err := NewTokenService(sessionRepo, userRepo, cfg, nil, nil)
	require.NoError(t, err)
	return svc, sessionRepo, userRepo
}

func createLiveUserForTokenTests(t *testing.T, userRepo *repository.UserRepository) string {
	t.Helper()
	ctx := context.Background()
	numStr := strconv.FormatInt(time.Now().UnixNano()%1_000_000_000, 10)
	mobile := "+8801" + strings.Repeat("0", 9-len(numStr)) + numStr
	u, err := userRepo.Create(ctx, mobile, "hash", "tok_live_"+uuid.NewString()[:8]+"@example.com", authnentityv1.UserStatus_USER_STATUS_ACTIVE)
	require.NoError(t, err)
	return u.UserId
}

func cleanupTokenUser(t *testing.T, userID string) {
	t.Helper()
	dbConn := testServiceLiveDB(t)
	_ = dbConn.Exec(`DELETE FROM authn_schema.sessions WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.otps WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.users WHERE user_id = ?`, userID).Error
}

func TestTokenService_LiveDB_ServerSideSessionAndCSRF(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	svc, _, userRepo := newLiveTokenService(t)
	userID := createLiveUserForTokenTests(t, userRepo)
	t.Cleanup(func() { cleanupTokenUser(t, userID) })

	ss, err := svc.GenerateServerSideSession(ctx, userID, "dev-web", authnentityv1.DeviceType_DEVICE_TYPE_WEB, "127.0.0.1", "ua")
	require.NoError(t, err)
	require.NotEmpty(t, ss.SessionID)
	require.NotEmpty(t, ss.SessionToken)
	require.NotEmpty(t, ss.CSRFToken)

	v, err := svc.ValidateServerSideSession(ctx, ss.SessionToken)
	require.NoError(t, err)
	require.True(t, v.Valid)
	require.Equal(t, userID, v.UserId)
	require.Equal(t, "SERVER_SIDE", v.SessionType)

	ok, err := svc.ValidateCSRFToken(ctx, ss.SessionID, ss.CSRFToken)
	require.NoError(t, err)
	require.True(t, ok)
	ok, err = svc.ValidateCSRFToken(ctx, ss.SessionID, "wrong")
	require.NoError(t, err)
	require.False(t, ok)
}

func TestTokenService_LiveDB_RefreshAndRevoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	svc, sessionRepo, userRepo := newLiveTokenService(t)
	userID := createLiveUserForTokenTests(t, userRepo)
	t.Cleanup(func() { cleanupTokenUser(t, userID) })

	pair, err := svc.GenerateJWT(ctx, userID, authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER.String(), "", "dev-api", authnentityv1.DeviceType_DEVICE_TYPE_API, "127.0.0.1", "ua")
	require.NoError(t, err)
	require.NotEmpty(t, pair.AccessToken)
	require.NotEmpty(t, pair.RefreshToken)

	ref, err := svc.RefreshJWT(ctx, pair.RefreshToken)
	require.NoError(t, err)
	require.NotEqual(t, pair.AccessToken, ref.AccessToken)
	require.NotEqual(t, pair.RefreshToken, ref.RefreshToken)
	require.NotEqual(t, pair.SessionID, ref.SessionId)

	_, err = sessionRepo.GetByID(ctx, pair.SessionID)
	require.Error(t, err) // old session should be revoked

	require.NoError(t, svc.RevokeSession(ctx, ref.SessionId))
	_, err = sessionRepo.GetByID(ctx, ref.SessionId)
	require.Error(t, err)
}

func TestTokenService_JWKSAndRandom(t *testing.T) {
	svc, _, _ := newLiveTokenService(t)
	resp, err := svc.GetJWKS(context.Background(), &authnservicev1.GetJWKSRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Keys, 1)
	require.Equal(t, "RSA", resp.Keys[0].Kty)
	require.NotEmpty(t, resp.Keys[0].N)
	require.NotEmpty(t, resp.Keys[0].E)

	s, err := generateSecureRandomString(16)
	require.NoError(t, err)
	require.Len(t, s, 32)
}
