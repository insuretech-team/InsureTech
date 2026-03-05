package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// otpBucket is the in-memory fallback bucket when Redis is unavailable.
type otpBucket struct {
	mu        sync.Mutex
	count     int
	lastReset time.Time
}

type ipWindowEntry struct {
	count      int
	windowFrom time.Time
}

var (
	otpRedisOnce   sync.Once
	otpRedisClient *redis.Client
)

func getOTPRedisClient() *redis.Client {
	otpRedisOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			return
		}
		otpRedisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	})
	return otpRedisClient
}

// OTPRateLimit returns middleware that limits OTP requests to maxRequests per
// composite key (IP + mobile number) per window duration.
//
// Uses Redis INCR+EXPIRE when REDIS_ADDR is set; falls back to in-memory
// sliding-window counter otherwise.
//
// Returns HTTP 429 with X-RateLimit-* and Retry-After headers on breach.
func OTPRateLimit(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	var buckets sync.Map // in-memory fallback

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)
			mobile := r.Header.Get("X-Mobile-Number")
			key := fmt.Sprintf("otp_rl:%s:%s", ip, mobile)

			var count int
			var retryAfterSec int
			allowed := true

			rdb := getOTPRedisClient()
			if rdb != nil {
				// Redis path: INCR + EXPIRE (set TTL only on first increment).
				ctx := context.Background()
				pipe := rdb.Pipeline()
				incrCmd := pipe.Incr(ctx, key)
				pipe.Expire(ctx, key, window)
				_, err := pipe.Exec(ctx)
				if err == nil {
					count = int(incrCmd.Val())
					ttl, _ := rdb.TTL(ctx, key).Result()
					retryAfterSec = int(math.Ceil(ttl.Seconds()))
					if retryAfterSec < 0 {
						retryAfterSec = int(window.Seconds())
					}
					allowed = count <= maxRequests
				} else {
					// Redis error → fail open (use in-memory fallback below).
					rdb = nil
				}
			}

			if rdb == nil {
				// In-memory fallback.
				val, _ := buckets.LoadOrStore(key, &otpBucket{lastReset: time.Now()})
				bucket := val.(*otpBucket)
				bucket.mu.Lock()
				now := time.Now()
				if now.Sub(bucket.lastReset) >= window {
					bucket.count = 0
					bucket.lastReset = now
				}
				bucket.count++
				count = bucket.count
				elapsed := now.Sub(bucket.lastReset)
				retryAfterSec = int(math.Ceil((window - elapsed).Seconds()))
				allowed = count <= maxRequests
				bucket.mu.Unlock()
			}

			remaining := maxRequests - count
			if remaining < 0 {
				remaining = 0
			}

			// Always set rate-limit headers.
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(maxRequests))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", strconv.Itoa(retryAfterSec))
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error":               "rate limit exceeded",
					"retry_after_seconds": retryAfterSec,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// IPWindowRateLimit enforces a fixed-window per-IP rate limit.
// Used for login/register endpoint throttling.
func IPWindowRateLimit(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	entries := map[string]*ipWindowEntry{}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)
			now := time.Now()

			mu.Lock()
			entry, ok := entries[ip]
			if !ok || now.Sub(entry.windowFrom) >= window {
				entry = &ipWindowEntry{count: 0, windowFrom: now}
				entries[ip] = entry
			}
			entry.count++
			currentCount := entry.count
			windowStart := entry.windowFrom
			mu.Unlock()

			if currentCount > maxRequests {
				retryAfter := int((window - now.Sub(windowStart)).Seconds())
				if retryAfter < 1 {
					retryAfter = 1
				}
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
