package service

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	refreshRateLimit  = 10               // max calls per window per user
	refreshRateWindow = 15 * time.Minute // sliding window duration
)

// refreshRateLimiter enforces per-user refresh token rate limits.
// Uses Redis INCR+EXPIRE when a Redis client is provided; falls back to an
// in-memory sliding-window counter otherwise.
type refreshRateLimiter struct {
	rdb     redis.UniversalClient // nil → in-memory fallback
	mu      sync.Mutex
	buckets map[string]*refreshBucket
}

type refreshBucket struct {
	count     int
	windowEnd time.Time
}

// newRefreshRateLimiter creates a rate limiter using the provided Redis client.
// Pass nil to use the in-memory fallback (single-instance only).
func newRefreshRateLimiter(rdb redis.UniversalClient) *refreshRateLimiter {
	return &refreshRateLimiter{
		rdb:     rdb,
		buckets: make(map[string]*refreshBucket),
	}
}

// Allow returns true if the refresh call for userID is within the rate limit.
func (r *refreshRateLimiter) Allow(userID string) bool {
	if r.rdb != nil {
		return r.redisAllow(userID)
	}
	return r.memAllow(userID)
}

func (r *refreshRateLimiter) redisAllow(userID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := "refresh_rl:" + userID
	pipe := r.rdb.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, refreshRateWindow)
	if _, err := pipe.Exec(ctx); err != nil {
		// Redis error → fail open.
		return true
	}
	return incrCmd.Val() <= int64(refreshRateLimit)
}

func (r *refreshRateLimiter) memAllow(userID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	b, ok := r.buckets[userID]
	if !ok || now.After(b.windowEnd) {
		r.buckets[userID] = &refreshBucket{
			count:     1,
			windowEnd: now.Add(refreshRateWindow),
		}
		return true
	}
	if b.count >= refreshRateLimit {
		return false
	}
	b.count++
	return true
}
