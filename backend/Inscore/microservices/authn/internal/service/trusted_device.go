package service

import (
	"context"
	"os"
	"strconv"
	"time"
)

func trustedDeviceTTL() time.Duration {
	days := 30
	if v := os.Getenv("TRUSTED_DEVICE_TTL_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			days = n
		}
	}
	return time.Duration(days) * 24 * time.Hour
}

func trustedDeviceKey(userID, deviceID string) string {
	return "trusted:device:" + userID + ":" + deviceID
}

func (s *AuthService) isTrustedDevice(ctx context.Context, userID, deviceID string) bool {
	if s == nil || s.tokenService == nil || s.tokenService.rdb == nil || userID == "" || deviceID == "" {
		return false
	}
	v, err := s.tokenService.rdb.Get(ctx, trustedDeviceKey(userID, deviceID)).Result()
	return err == nil && v == "1"
}

func (s *AuthService) markTrustedDevice(ctx context.Context, userID, deviceID string) {
	if s == nil || s.tokenService == nil || s.tokenService.rdb == nil || userID == "" || deviceID == "" {
		return
	}
	_ = s.tokenService.rdb.Set(ctx, trustedDeviceKey(userID, deviceID), "1", trustedDeviceTTL()).Err()
}
