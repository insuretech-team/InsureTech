package cache

// permission_cache.go — Redis-backed permission cache for AuthZ service
// Reduces database queries for repeated permission checks.
//
// Cache key format: "authz:perm:{user_id}:{domain}:{object}:{action}"
// Cache value: "allow" or "deny"
// TTL: configurable (default 5-15 minutes)

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// PermissionCache caches authorization decisions in Redis.
type PermissionCache struct {
	redis redis.UniversalClient
	ttl   time.Duration
}

// NewPermissionCache creates a new permission cache.
// ttl: cache time-to-live (recommended: 5-15 minutes)
func NewPermissionCache(redisClient redis.UniversalClient, ttl time.Duration) *PermissionCache {
	if ttl == 0 {
		ttl = 5 * time.Minute // default
	}
	return &PermissionCache{
		redis: redisClient,
		ttl:   ttl,
	}
}

// Get retrieves a cached permission decision.
// Returns (decision, found). Decision is true=allow, false=deny.
func (c *PermissionCache) Get(ctx context.Context, userID, domain, object, action string) (bool, bool) {
	if c.redis == nil {
		return false, false
	}

	key := c.buildKey(userID, domain, object, action)
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return false, false // cache miss or error
	}

	return val == "allow", true
}

// Set stores a permission decision in cache.
func (c *PermissionCache) Set(ctx context.Context, userID, domain, object, action string, allowed bool) error {
	if c.redis == nil {
		return nil // no-op if Redis not available
	}

	key := c.buildKey(userID, domain, object, action)
	val := "deny"
	if allowed {
		val = "allow"
	}

	return c.redis.Set(ctx, key, val, c.ttl).Err()
}

// InvalidateUser invalidates all cached permissions for a specific user.
// Called when user roles change.
func (c *PermissionCache) InvalidateUser(ctx context.Context, userID string) error {
	if c.redis == nil {
		return nil
	}

	pattern := fmt.Sprintf("authz:perm:%s:*", userID)
	return c.deleteByPattern(ctx, pattern)
}

// InvalidateDomain invalidates all cached permissions for a domain.
// Called when domain-level policies change.
func (c *PermissionCache) InvalidateDomain(ctx context.Context, domain string) error {
	if c.redis == nil {
		return nil
	}

	pattern := fmt.Sprintf("authz:perm:*:%s:*", domain)
	return c.deleteByPattern(ctx, pattern)
}

// InvalidateAll clears the entire permission cache.
// Called on major policy changes.
func (c *PermissionCache) InvalidateAll(ctx context.Context) error {
	if c.redis == nil {
		return nil
	}

	pattern := "authz:perm:*"
	return c.deleteByPattern(ctx, pattern)
}

// buildKey constructs the cache key.
func (c *PermissionCache) buildKey(userID, domain, object, action string) string {
	return fmt.Sprintf("authz:perm:%s:%s:%s:%s", userID, domain, object, action)
}

// deleteByPattern deletes all keys matching a pattern (using SCAN).
func (c *PermissionCache) deleteByPattern(ctx context.Context, pattern string) error {
	iter := c.redis.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := c.redis.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
