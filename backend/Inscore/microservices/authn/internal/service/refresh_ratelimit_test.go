package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestRefreshRateLimiter_TenCallsAllowed verifies that exactly 10 calls within
// the window all return true.
func TestRefreshRateLimiter_TenCallsAllowed(t *testing.T) {
	rl := newRefreshRateLimiter(nil)
	userID := "user-abc"

	for i := 0; i < refreshRateLimit; i++ {
		require.True(t, rl.Allow(userID), "call %d should be allowed", i+1)
	}
}

// TestRefreshRateLimiter_EleventhCallDenied verifies that the 11th call within
// the same window returns false.
func TestRefreshRateLimiter_EleventhCallDenied(t *testing.T) {
	rl := newRefreshRateLimiter(nil)
	userID := "user-xyz"

	for i := 0; i < refreshRateLimit; i++ {
		rl.Allow(userID)
	}

	require.False(t, rl.Allow(userID), "11th call should be denied")
}

// TestRefreshRateLimiter_WindowExpiry verifies that after the window expires
// the counter resets and calls are allowed again.
func TestRefreshRateLimiter_WindowExpiry(t *testing.T) {
	rl := newRefreshRateLimiter(nil)
	userID := "user-reset"

	// Exhaust the limit.
	for i := 0; i < refreshRateLimit; i++ {
		rl.Allow(userID)
	}
	require.False(t, rl.Allow(userID), "should be rate-limited before window expires")

	// Manually expire the window by backdating windowEnd.
	rl.mu.Lock()
	rl.buckets[userID].windowEnd = time.Now().Add(-1 * time.Second)
	rl.mu.Unlock()

	// After the window has expired, the next call should create a fresh bucket.
	require.True(t, rl.Allow(userID), "should be allowed after window expiry")

	// And subsequent calls within the new window should still be allowed.
	for i := 1; i < refreshRateLimit; i++ {
		require.True(t, rl.Allow(userID), "call %d in new window should be allowed", i+1)
	}

	// 11th call in the new window should be denied.
	require.False(t, rl.Allow(userID), "11th call in new window should be denied")
}

// TestRefreshRateLimiter_IndependentUsers verifies that rate limits are
// tracked independently per user.
func TestRefreshRateLimiter_IndependentUsers(t *testing.T) {
	rl := newRefreshRateLimiter(nil)

	// Exhaust user A.
	for i := 0; i < refreshRateLimit; i++ {
		rl.Allow("user-A")
	}
	require.False(t, rl.Allow("user-A"), "user-A should be rate-limited")

	// user-B should be unaffected.
	require.True(t, rl.Allow("user-B"), "user-B should not be affected by user-A limit")
}
