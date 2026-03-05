package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/stretchr/testify/require"
)

func TestTokenService_ValidateJWT_PopulatesUserTypeFromClaim(t *testing.T) {
	cfg := &config.Config{}
	cfg.JWT.Issuer = "test"

	generateTempRSAKeys(cfg)
	s, err := NewTokenService(nil, nil, cfg, nil, nil)
	require.NoError(t, err)

	claims := jwt.MapClaims{
		"sub":  "u1",
		"sid":  "s1",
		"utp":  "AGENT",
		"exp":  time.Now().Add(5 * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  cfg.JWT.Issuer,
		"type": "access",
		"jti":  "j1",
	}
	tok := signTestToken(t, s.rsaPrivateKey, s.keyID, claims)

	resp, err := s.ValidateJWT(context.Background(), tok)
	require.NoError(t, err)
	require.True(t, resp.Valid)
	require.Equal(t, "u1", resp.UserId)
	require.Equal(t, "AGENT", resp.UserType)
	require.Equal(t, "s1", resp.SessionId)
}

func TestTokenService_ValidateJWT_MissingUserTypeClaim_DoesNotFail(t *testing.T) {
	cfg := &config.Config{}
	cfg.JWT.Issuer = "test"

	generateTempRSAKeys(cfg)
	s, err := NewTokenService(nil, nil, cfg, nil, nil)
	require.NoError(t, err)

	claims := jwt.MapClaims{
		"sub":  "u1",
		"sid":  "s1",
		"exp":  time.Now().Add(5 * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  cfg.JWT.Issuer,
		"type": "access",
		"jti":  "j1",
	}
	tok := signTestToken(t, s.rsaPrivateKey, s.keyID, claims)

	resp, err := s.ValidateJWT(context.Background(), tok)
	require.NoError(t, err)
	require.True(t, resp.Valid)
	require.Equal(t, "u1", resp.UserId)
	require.Equal(t, "", resp.UserType)
}
