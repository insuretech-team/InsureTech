package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/middleware"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestArgon2id_HashAndVerify(t *testing.T) {
	hash, err := HashPassword("P@ssw0rd!Strong", nil)
	require.NoError(t, err)
	require.True(t, IsArgon2idHash(hash))

	ok, err := VerifyPassword("P@ssw0rd!Strong", hash)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = VerifyPassword("wrong", hash)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestPasswordHash_DualVerify(t *testing.T) {
	argHash, err := hashPassword("Abcdef1!")
	require.NoError(t, err)

	valid, needsRehash, err := verifyPassword("Abcdef1!", argHash)
	require.NoError(t, err)
	require.True(t, valid)
	require.False(t, needsRehash)

	bcHash, err := bcrypt.GenerateFromPassword([]byte("Abcdef1!"), bcrypt.DefaultCost)
	require.NoError(t, err)
	valid, needsRehash, err = verifyPassword("Abcdef1!", string(bcHash))
	require.NoError(t, err)
	require.True(t, valid)
	require.True(t, needsRehash)
}

func TestHelpers_ParseAndMapping(t *testing.T) {
	require.Equal(t, authnentityv1.DeviceType_DEVICE_TYPE_WEB, parseDeviceType("WEB"))
	require.Equal(t, authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_ANDROID, parseDeviceType("ANDROID"))
	require.Equal(t, authnentityv1.DeviceType_DEVICE_TYPE_API, parseDeviceType("unknown"))

	require.Equal(t, authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE, mapDeviceTypeToSessionType(authnentityv1.DeviceType_DEVICE_TYPE_WEB))
	require.Equal(t, authnentityv1.SessionType_SESSION_TYPE_JWT, mapDeviceTypeToSessionType(authnentityv1.DeviceType_DEVICE_TYPE_API))

	require.Equal(t, "PORTAL_B2C", portalConfigKeyForUserType("B2C_CUSTOMER"))
	require.Equal(t, "PORTAL_SYSTEM", portalConfigKeyForUserType("USER_TYPE_SYSTEM_USER"))
	require.Equal(t, "", portalConfigKeyForUserType("UNKNOWN"))

	require.Equal(t, authnentityv1.Gender_GENDER_MALE, parseGender("MALE"))
	require.Equal(t, authnentityv1.Gender_GENDER_FEMALE, parseGender("GENDER_FEMALE"))
	require.Equal(t, authnentityv1.Gender_GENDER_UNSPECIFIED, parseGender("nope"))

	require.True(t, isEmailAuthUser(authnentityv1.UserType_USER_TYPE_AGENT))
	require.False(t, isEmailAuthUser(authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER))
	require.Equal(t, "u***@example.com", maskEmail("user@example.com"))
	require.Equal(t, "***", maskEmail("invalid"))
	require.Equal(t, authnentityv1.UserType_USER_TYPE_SYSTEM_USER, parseUserType("SYSTEM_USER"))
}

func TestSessionTokenLookup_Deterministic(t *testing.T) {
	a := sessionTokenLookup("token-1")
	b := sessionTokenLookup("token-1")
	c := sessionTokenLookup("token-2")
	require.Equal(t, a, b)
	require.NotEqual(t, a, c)
	require.Len(t, a, 64)
}

func TestJTIBlocklist_Noop(t *testing.T) {
	bl := NewJTIBlocklist(nil)
	err := bl.Block(context.Background(), "jti-1", time.Now().Add(1*time.Minute))
	require.NoError(t, err)

	blocked, err := bl.IsBlocked(context.Background(), "jti-1")
	require.NoError(t, err)
	require.False(t, blocked)

	opt := WithJTIBlocklist(bl)
	ts := &TokenService{}
	require.NotNil(t, opt(ts))
}

func TestMFASession_AndTrustedDevice_NoRedis(t *testing.T) {
	svc := &AuthService{tokenService: &TokenService{}}
	t.Setenv("MFA_SESSION_TTL_SECONDS", "180")
	require.Equal(t, 180*time.Second, mfaSessionTTL())
	require.Equal(t, "mfa:session:abc", mfaSessionKey("abc"))

	token, err := svc.StoreMFASessionToken(context.Background(), "u1", "d1", "ANDROID", "127.0.0.1")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	_, _, _, _, err = svc.ConsumeMFASessionToken(context.Background(), token)
	require.Error(t, err)

	t.Setenv("TRUSTED_DEVICE_TTL_DAYS", "7")
	require.Equal(t, 7*24*time.Hour, trustedDeviceTTL())
	require.Equal(t, "trusted:device:u1:d1", trustedDeviceKey("u1", "d1"))
	require.False(t, svc.isTrustedDevice(context.Background(), "u1", "d1"))
	require.NotPanics(t, func() { svc.markTrustedDevice(context.Background(), "u1", "d1") })
}

func TestSessionLimiter_NoRedis(t *testing.T) {
	sl := NewSessionLimiter(nil, 0)
	require.Equal(t, "sessions:active:u1", sl.key("u1"))

	evicted, err := sl.TrackSession(context.Background(), "u1", "s1", time.Now().Add(1*time.Hour))
	require.NoError(t, err)
	require.Empty(t, evicted)
	require.NoError(t, sl.RemoveSession(context.Background(), "u1", "s1"))
	count, err := sl.ActiveCount(context.Background(), "u1")
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func TestJWKSService_GetJWKS(t *testing.T) {
	nilSvc := NewJWKSService(nil, "kid")
	resp, err := nilSvc.GetJWKS(context.Background(), &authnservicev1.GetJWKSRequest{})
	require.NoError(t, err)
	require.Empty(t, resp.Keys)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	svc := NewJWKSService(&priv.PublicKey, "kid-1")
	resp, err = svc.GetJWKS(context.Background(), &authnservicev1.GetJWKSRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Keys, 1)
	require.Equal(t, "kid-1", resp.Keys[0].Kid)
	require.Equal(t, "RSA", resp.Keys[0].Kty)
	require.Equal(t, "RS256", resp.Keys[0].Alg)
}

func TestTOTPService_GenerateAndValidate(t *testing.T) {
	ts := NewTOTPService()
	uri, secret, err := ts.GenerateKey("InsureTech", "user@example.com")
	require.NoError(t, err)
	require.Contains(t, uri, "otpauth://")
	require.NotEmpty(t, secret)

	code, err := totp.GenerateCode(secret, time.Now().UTC())
	require.NoError(t, err)
	valid, err := ts.Validate(code, secret)
	require.NoError(t, err)
	require.True(t, valid)
}

func TestStubCryptoHelpers(t *testing.T) {
	t.Setenv("TOTP_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString(make([]byte, 32)))
	key := totpEncryptionKey()
	require.Len(t, key, 32)

	ciphertext, err := aesGCMEncrypt("secret", key)
	require.NoError(t, err)
	plain, err := aesGCMDecrypt(ciphertext, key)
	require.NoError(t, err)
	require.Equal(t, "secret", plain)

	_, err = aesGCMDecrypt("not-base64", key)
	require.Error(t, err)
}

func TestServiceMethods_NilRepoErrors(t *testing.T) {
	svc := &AuthService{}
	ctx := context.Background()

	_, err := svc.InitiateKYC(ctx, &authnservicev1.InitiateKYCRequest{UserId: "u1"})
	require.Error(t, err)
	_, err = svc.GetKYCStatus(ctx, &authnservicev1.GetKYCStatusRequest{UserId: "u1"})
	require.Error(t, err)
	_, err = svc.ApproveKYC(ctx, &authnservicev1.ApproveKYCRequest{KycId: "k1"})
	require.Error(t, err)
	_, err = svc.RejectKYC(ctx, &authnservicev1.RejectKYCRequest{KycId: "k1"})
	require.Error(t, err)
	_, err = svc.VerifyDocument(ctx, &authnservicev1.VerifyDocumentRequest{UserDocumentId: "d1"})
	require.Error(t, err)

	_, err = svc.CreateVoiceSession(ctx, &authnservicev1.CreateVoiceSessionRequest{UserId: "u1"})
	require.Error(t, err)
	_, err = svc.GetVoiceSession(ctx, &authnservicev1.GetVoiceSessionRequest{VoiceSessionId: "v1"})
	require.Error(t, err)
	_, err = svc.EndVoiceSession(ctx, &authnservicev1.EndVoiceSessionRequest{VoiceSessionId: "v1"})
	require.Error(t, err)

	_, err = svc.UploadUserDocument(ctx, &authnservicev1.UploadUserDocumentRequest{})
	require.Error(t, err)
	_, err = svc.ListUserDocuments(ctx, &authnservicev1.ListUserDocumentsRequest{UserId: "u1"})
	require.Error(t, err)
	_, err = svc.GetUserDocument(ctx, &authnservicev1.GetUserDocumentRequest{UserDocumentId: "d1"})
	require.Error(t, err)
	_, err = svc.DeleteUserDocument(ctx, &authnservicev1.DeleteUserDocumentRequest{UserDocumentId: "d1"})
	require.Error(t, err)
	_, err = svc.ListDocumentTypes(ctx, &authnservicev1.ListDocumentTypesRequest{})
	require.Error(t, err)

	_, err = svc.CreateUserProfile(ctx, &authnservicev1.CreateUserProfileRequest{UserId: "u1"})
	require.Error(t, err)
	_, err = svc.GetUserProfile(ctx, &authnservicev1.GetUserProfileRequest{UserId: "u1"})
	require.Error(t, err)
	_, err = svc.UpdateUserProfile(ctx, &authnservicev1.UpdateUserProfileRequest{UserId: "u1"})
	require.Error(t, err)

	_, err = svc.CreateAPIKey(ctx, &authnservicev1.CreateAPIKeyRequest{Name: "k"})
	require.Error(t, err)
	_, err = svc.ListAPIKeys(ctx, &authnservicev1.ListAPIKeysRequest{OwnerId: "o1"})
	require.Error(t, err)
	_, err = svc.RevokeAPIKey(ctx, &authnservicev1.RevokeAPIKeyRequest{KeyId: "k1"})
	require.Error(t, err)
}

func TestConstructorsAndSimpleHelpers(t *testing.T) {
	cfg := &config.Config{}
	cfg.JWT.Issuer = "test"
	cfg.JWT.KeyID = "kid-1"
	cfg.JWT.AccessTokenDuration = time.Minute
	cfg.JWT.RefreshTokenDuration = time.Hour
	generateTempRSAKeys(cfg)

	meta := middleware.NewMetadataExtractor()
	pub := events.NewPublisher((events.EventProducer)(nil))

	svc := NewAuthServiceWithAPIKey(nil, nil, nil, nil, nil, nil, pub, cfg, meta)
	require.NotNil(t, svc)
	require.Equal(t, cfg, svc.config)

	ts1, err := NewTokenServiceWithRedis(nil, nil, cfg, pub, meta, nil)
	require.NoError(t, err)
	require.NotNil(t, ts1)
	ts2, err := NewTokenServiceWithSessionLimiter(nil, nil, cfg, pub, meta, nil, 2)
	require.NoError(t, err)
	require.NotNil(t, ts2)
	require.NotNil(t, ts2.sessionLimiter)

	require.NoError(t, ts1.BlockJTI(context.Background(), "j1", time.Minute))
	require.False(t, ts1.isJTIBlocked(context.Background(), "j1"))
	require.NotNil(t, ts1.PublicKey())
	require.Equal(t, cfg.JWT.KeyID, ts1.KeyID())
}

func TestOTPHelpersAndResendErrors(t *testing.T) {
	code, err := generateOTPCode(6)
	require.NoError(t, err)
	require.Len(t, code, 6)
	_, err = generateOTPCode(3)
	require.Error(t, err)

	otpSvc := &OTPService{}
	require.Error(t, otpSvc.HandleDLR(context.Background(), []byte(`{}`)))

	authSvc := &AuthService{
		config: &config.Config{},
	}
	authSvc.config.Security.OTPCooldown = time.Minute
	_, err = authSvc.ResendOTP(context.Background(), nil)
	require.Error(t, err)
	_, err = authSvc.ResendOTP(context.Background(), &authnservicev1.ResendOTPRequest{OriginalOtpId: "x"})
	require.Error(t, err)
}

func TestTOTPAliasMethodsAndProfileURL(t *testing.T) {
	svc := &AuthService{}
	require.Panics(t, func() {
		_, _ = svc.EnrollTOTP(context.Background(), &authnservicev1.EnableTOTPRequest{UserId: "u1"})
	})
	require.Panics(t, func() {
		_, _ = svc.ConfirmTOTP(context.Background(), &authnservicev1.VerifyTOTPRequest{UserId: "u1", TotpCode: "123456"})
	})

	t.Setenv("AWS_ACCESS_KEY_ID", "dummy")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "dummy")
	t.Setenv("AWS_REGION", "ap-southeast-1")
	t.Setenv("S3_BUCKET", "bucket-test")
	resp, err := svc.GetProfilePhotoUploadURL(context.Background(), &authnservicev1.GetProfilePhotoUploadURLRequest{
		UserId:      "u1",
		ContentType: "image/jpeg",
	})
	require.NoError(t, err)
	require.Contains(t, resp.UploadUrl, "bucket-test")
	require.Contains(t, resp.FileUrl, "bucket-test")
	require.Equal(t, int32(900), resp.ExpiresInSeconds)
}

func TestAdditionalCoverageHelpers(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	t.Cleanup(func() { _ = rdb.Close() })

	otpSvc := NewOTPServiceWithRedis(nil, nil, nil, &config.Config{}, nil, rdb)
	require.NotNil(t, otpSvc)
	require.NotNil(t, otpSvc.redisClient)
	otpSvc.incrementRateLimit(context.Background(), "u1", "login")

	rl := newRefreshRateLimiter(rdb)
	require.True(t, rl.redisAllow("user-x")) // fail-open on connection error still executes redis path

	uri, secret, err := EnrollTOTPForUser("InsureTech", "user@example.com")
	require.NoError(t, err)
	require.Contains(t, uri, "otpauth://")
	code, err := totp.GenerateCode(secret, time.Now().UTC())
	require.NoError(t, err)
	ok, err := ValidateTOTPCode(code, secret)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestUpdateNotificationPreferences_LiveDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)
	ctx := context.Background()

	numStr := strconv.FormatInt(time.Now().UnixNano()%1_000_000_000, 10)
	mobile := "+8801" + strings.Repeat("0", 9-len(numStr)) + numStr
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "notify_"+strconv.FormatInt(time.Now().UnixNano(), 10)+"@example.com", "Str0ng!Notify1")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	resp, err := svc.UpdateNotificationPreferences(ctx, &authnservicev1.UpdateNotificationPreferencesRequest{
		UserId:                 userID,
		NotificationPreference: `{"email":true}`,
		PreferredLanguage:      "en",
	})
	require.NoError(t, err)
	require.Contains(t, resp.Message, "updated")
}
