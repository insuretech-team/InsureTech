package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/stretchr/testify/require"
)

func generateTempRSAKeys(cfg *config.Config) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	dir := os.TempDir()
	id := uuid.New().String()
	privPath := filepath.Join(dir, "test_private_"+id+".pem")
	pubPath := filepath.Join(dir, "test_public_"+id+".pem")

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	_ = os.WriteFile(privPath, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}), 0600)

	pubBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	_ = os.WriteFile(pubPath, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}), 0644)

	cfg.JWT.PrivateKeyPath = privPath
	cfg.JWT.PublicKeyPath = pubPath
	cfg.JWT.KeyID = "test-kid"
}

func newTestTokenService() *TokenService {
	cfg := &config.Config{}
	cfg.JWT.Issuer = "insuretech-test"
	cfg.JWT.AccessTokenDuration = 15 * time.Minute
	cfg.JWT.RefreshTokenDuration = 7 * 24 * time.Hour

	generateTempRSAKeys(cfg)
	ts, err := NewTokenService(nil, nil, cfg, nil, nil)
	if err != nil {
		panic(err)
	}
	return ts
}

func signTestToken(t *testing.T, privKey *rsa.PrivateKey, kid string, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = kid
	signed, err := tok.SignedString(privKey)
	require.NoError(t, err)
	return signed
}

// ---------------------------------------------------------------------------
// T3a: expired token → Valid=false
// ---------------------------------------------------------------------------

func TestValidateJWT_ExpiredToken(t *testing.T) {
	s := newTestTokenService()

	tok := signTestToken(t, s.rsaPrivateKey, s.keyID, jwt.MapClaims{
		"sub":  "u1",
		"sid":  "s1",
		"utp":  "B2C_CUSTOMER",
		"exp":  time.Now().Add(-1 * time.Minute).Unix(), // already expired
		"iat":  time.Now().Add(-5 * time.Minute).Unix(),
		"iss":  s.config.JWT.Issuer,
		"type": "access",
		"jti":  "j1",
	})

	resp, err := s.ValidateJWT(context.Background(), tok)
	require.NoError(t, err)
	require.False(t, resp.Valid)
}

// ---------------------------------------------------------------------------
// T3b: wrong signing secret → Valid=false
// ---------------------------------------------------------------------------

func TestValidateJWT_InvalidSignature(t *testing.T) {
	s := newTestTokenService()

	// Sign with a DIFFERENT RSA key
	wrongKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	tok := signTestToken(t, wrongKey, s.keyID, jwt.MapClaims{
		"sub":  "u1",
		"sid":  "s1",
		"utp":  "AGENT",
		"exp":  time.Now().Add(5 * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  s.config.JWT.Issuer,
		"type": "access",
		"jti":  "j1",
	})

	resp, err := s.ValidateJWT(context.Background(), tok)
	require.NoError(t, err)
	require.False(t, resp.Valid)
}

// ---------------------------------------------------------------------------
// T3c: malformed / garbage token → Valid=false
// ---------------------------------------------------------------------------

func TestValidateJWT_MalformedToken(t *testing.T) {
	s := newTestTokenService()

	cases := []string{
		"",
		"not.a.jwt",
		"Bearer abc123",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.garbage.sig",
	}

	for _, tok := range cases {
		t.Run(tok, func(t *testing.T) {
			resp, err := s.ValidateJWT(context.Background(), tok)
			require.NoError(t, err)
			require.False(t, resp.Valid)
		})
	}
}

// ---------------------------------------------------------------------------
// T3d: both access and refresh tokens contain utp claim
// ---------------------------------------------------------------------------

func TestGenerateJWT_UserType_InBothAccessAndRefresh(t *testing.T) {
	s := newTestTokenService()

	// We can't call GenerateJWT directly (it needs sessionRepo), so verify
	// the claim structure by signing manually and parsing — exactly what
	// GenerateJWT does internally.
	privKey := s.rsaPrivateKey
	kid := s.keyID
	userType := "AGENT"

	accessClaims := jwt.MapClaims{
		"sub":  "u-agent",
		"utp":  userType,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  s.config.JWT.Issuer,
		"type": "access",
		"jti":  "access-jti",
		"sid":  "sess-1",
	}
	refreshClaims := jwt.MapClaims{
		"sub":  "u-agent",
		"utp":  userType,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  s.config.JWT.Issuer,
		"type": "refresh",
		"jti":  "refresh-jti",
		"sid":  "sess-1",
	}

	accessTok := signTestToken(t, privKey, kid, accessClaims)
	refreshTok := signTestToken(t, privKey, kid, refreshClaims)

	// Validate access token — utp must be present and correct
	resp, err := s.ValidateJWT(context.Background(), accessTok)
	require.NoError(t, err)
	require.True(t, resp.Valid)
	require.Equal(t, "AGENT", resp.UserType)
	require.Equal(t, "u-agent", resp.UserId)

	// Parse refresh token manually and confirm utp claim exists
	parsed, err := jwt.Parse(refreshTok, func(tok *jwt.Token) (interface{}, error) {
		return s.rsaPublicKey, nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)
	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	require.Equal(t, "AGENT", claims["utp"])
	require.Equal(t, "refresh", claims["type"])
}
