package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// JTIBlocklist uses Redis to track revoked JWT IDs.
// Key pattern: jti:blocked:<jti>  value: "1"  TTL = remaining token lifetime.
//
// NOTE: TokenService already provides inline BlockJTI / isJTIBlocked methods
// that operate on the same Redis client. JTIBlocklist is a standalone helper
// for callers that need to manage blocklist operations outside of TokenService
// (e.g. integration tests, admin tooling, or future microservice boundaries).
type JTIBlocklist struct {
	rdb redis.UniversalClient
}

// NewJTIBlocklist creates a new JTIBlocklist backed by the given Redis client.
// Pass nil to get a no-op blocklist (all Block calls silently succeed, all
// IsBlocked calls return false).
func NewJTIBlocklist(rdb redis.UniversalClient) *JTIBlocklist {
	return &JTIBlocklist{rdb: rdb}
}

// Block adds a JTI to the blocklist with a TTL equal to the token's remaining
// lifetime. expiry is the token's exp claim time. The minimum TTL is 1 second;
// already-expired tokens are skipped (zero TTL).
func (b *JTIBlocklist) Block(ctx context.Context, jti string, expiry time.Time) error {
	if b == nil || b.rdb == nil || jti == "" {
		return nil
	}
	ttl := time.Until(expiry)
	if ttl <= 0 {
		return nil // token already expired — no point blocking
	}
	if ttl < time.Second {
		ttl = time.Second
	}
	key := "jti:blocked:" + jti
	return b.rdb.Set(ctx, key, "1", ttl).Err()
}

// IsBlocked returns true if the JTI is currently on the blocklist.
func (b *JTIBlocklist) IsBlocked(ctx context.Context, jti string) (bool, error) {
	if b == nil || b.rdb == nil || jti == "" {
		return false, nil
	}
	key := "jti:blocked:" + jti
	count, err := b.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// WithJTIBlocklist returns the TokenService unchanged (JTI blocking is handled
// natively by TokenService via its embedded Redis client). This method exists
// so callers can wire a *JTIBlocklist into the service chain using the
// functional-options pattern without needing to change call sites.
//
// If you need the standalone JTIBlocklist to share the same Redis instance as
// the TokenService, obtain the client from the TokenService's constructor
// (NewTokenServiceWithRedis) and pass it to NewJTIBlocklist.
func WithJTIBlocklist(bl *JTIBlocklist) func(*TokenService) *TokenService {
	return func(ts *TokenService) *TokenService {
		// TokenService already maintains its own rdb field and calls
		// BlockJTI / isJTIBlocked directly. If the caller supplied a
		// JTIBlocklist with a different client, prefer that client so
		// both operate on the same keyspace.
		if bl != nil && bl.rdb != nil && ts.rdb == nil {
			ts.rdb = bl.rdb
		}
		return ts
	}
}
