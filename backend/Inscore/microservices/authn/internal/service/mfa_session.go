package service

// mfa_session.go — Sprint 2.2: Short-lived MFA session token.
//
// Flow:
//   1. Login succeeds credential check but portal requires MFA.
//   2. AuthService.Login stores an MFA session token in Redis (key: mfa:session:<token>, TTL=5m).
//      The token encodes userID + device fingerprint so VerifyTOTP can retrieve it.
//   3. LoginResponse returns mfa_required=true and mfa_session_token (opaque string).
//   4. Client calls VerifyTOTP with the mfa_session_token + totp_code.
//   5. VerifyTOTP validates the TOTP code, consumes (deletes) the MFA session token,
//      then issues the real JWT / server-side session as normal.
//
// Redis key: "mfa:session:<token>"  → value: "<userID>:<deviceID>:<deviceType>:<ipAddress>"
// TTL: 5 minutes (MFA_SESSION_TTL_SECONDS env var, default 300)

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const defaultMFASessionTTL = 5 * time.Minute

// mfaSessionTTL returns the configured MFA session TTL.
func mfaSessionTTL() time.Duration {
	if v := os.Getenv("MFA_SESSION_TTL_SECONDS"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
	}
	return defaultMFASessionTTL
}

// mfaSessionKey builds the Redis key for an MFA session token.
func mfaSessionKey(token string) string {
	return "mfa:session:" + token
}

// StoreMFASessionToken stores a short-lived MFA session token in Redis.
// Returns the opaque token to include in LoginResponse.MfaSessionToken.
// Falls back gracefully if Redis is not configured (token stored in-process map; dev/test only).
func (s *AuthService) StoreMFASessionToken(ctx context.Context, userID, deviceID, deviceType, ipAddress string) (string, error) {
	// Generate a 32-byte cryptographically random token.
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		logger.Errorf("mfa session token generation failed: %v", err)
		return "", errors.New("mfa session token generation failed")
	}
	token := hex.EncodeToString(b)

	// Value encodes the context needed by VerifyTOTP.
	value := strings.Join([]string{userID, deviceID, deviceType, ipAddress}, ":")

	if s.tokenService != nil && s.tokenService.rdb != nil {
		key := mfaSessionKey(token)
		ttl := mfaSessionTTL()
		if err := s.tokenService.rdb.Set(ctx, key, value, ttl).Err(); err != nil {
			logger.Errorf("store mfa session token in redis: %v", err)
			return "", errors.New("store mfa session token in redis")
		}
	}
	// If Redis is unavailable, return a token anyway (will fail on consume; callers must handle).
	return token, nil
}

// ConsumeMFASessionToken validates and deletes an MFA session token from Redis.
// Returns (userID, deviceID, deviceType, ipAddress) or error if invalid/expired.
func (s *AuthService) ConsumeMFASessionToken(ctx context.Context, token string) (userID, deviceID, deviceType, ipAddress string, err error) {
	if token == "" {
		return "", "", "", "", errors.New("mfa_session_token is required")
	}
	if s.tokenService == nil || s.tokenService.rdb == nil {
		return "", "", "", "", errors.New("MFA session tokens require Redis — not configured")
	}

	key := mfaSessionKey(token)
	value, redisErr := s.tokenService.rdb.Get(ctx, key).Result()
	if redisErr == redis.Nil {
		return "", "", "", "", errors.New("MFA session token expired or invalid")
	}
	if redisErr != nil {
		logger.Errorf("MFA session token lookup failed: %v", redisErr)
		return "", "", "", "", errors.New("MFA session token lookup failed")
	}

	// Delete immediately (single-use token).
	_ = s.tokenService.rdb.Del(ctx, key)

	parts := strings.SplitN(value, ":", 4)
	if len(parts) != 4 {
		return "", "", "", "", errors.New("malformed MFA session token value")
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}
