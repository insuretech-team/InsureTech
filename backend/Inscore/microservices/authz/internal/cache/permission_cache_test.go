package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMiniRedis(t *testing.T) (*miniredis.Miniredis, redis.UniversalClient) {
	t.Helper()
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestPermissionCache_GetSet(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	cache := NewPermissionCache(client, 5*time.Minute)
	ctx := context.Background()

	// Cache miss
	decision, found := cache.Get(ctx, "user:u1", "system:root", "svc:policy/read", "GET")
	assert.False(t, found)
	assert.False(t, decision)

	// Set allow
	err := cache.Set(ctx, "user:u1", "system:root", "svc:policy/read", "GET", true)
	require.NoError(t, err)

	// Cache hit - allow
	decision, found = cache.Get(ctx, "user:u1", "system:root", "svc:policy/read", "GET")
	assert.True(t, found)
	assert.True(t, decision)

	// Set deny
	err = cache.Set(ctx, "user:u1", "system:root", "svc:policy/create", "POST", false)
	require.NoError(t, err)

	// Cache hit - deny
	decision, found = cache.Get(ctx, "user:u1", "system:root", "svc:policy/create", "POST")
	assert.True(t, found)
	assert.False(t, decision)
}

func TestPermissionCache_InvalidateUser(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	cache := NewPermissionCache(client, 5*time.Minute)
	ctx := context.Background()

	// Set multiple permissions for user
	_ = cache.Set(ctx, "user:u1", "system:root", "svc:policy/read", "GET", true)
	_ = cache.Set(ctx, "user:u1", "customer:t1", "svc:claim/read", "GET", true)
	_ = cache.Set(ctx, "user:u2", "system:root", "svc:policy/read", "GET", true)

	// Invalidate user:u1
	err := cache.InvalidateUser(ctx, "user:u1")
	require.NoError(t, err)

	// user:u1 permissions should be gone
	_, found := cache.Get(ctx, "user:u1", "system:root", "svc:policy/read", "GET")
	assert.False(t, found)

	_, found = cache.Get(ctx, "user:u1", "customer:t1", "svc:claim/read", "GET")
	assert.False(t, found)

	// user:u2 permissions should still exist
	_, found = cache.Get(ctx, "user:u2", "system:root", "svc:policy/read", "GET")
	assert.True(t, found)
}

func TestPermissionCache_InvalidateDomain(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	cache := NewPermissionCache(client, 5*time.Minute)
	ctx := context.Background()

	// Set permissions in different domains
	_ = cache.Set(ctx, "user:u1", "system:root", "svc:policy/read", "GET", true)
	_ = cache.Set(ctx, "user:u2", "system:root", "svc:claim/read", "GET", true)
	_ = cache.Set(ctx, "user:u1", "customer:t1", "svc:policy/read", "GET", true)

	// Invalidate system:root domain
	err := cache.InvalidateDomain(ctx, "system:root")
	require.NoError(t, err)

	// system:root permissions should be gone
	_, found := cache.Get(ctx, "user:u1", "system:root", "svc:policy/read", "GET")
	assert.False(t, found)

	_, found = cache.Get(ctx, "user:u2", "system:root", "svc:claim/read", "GET")
	assert.False(t, found)

	// customer:t1 permissions should still exist
	_, found = cache.Get(ctx, "user:u1", "customer:t1", "svc:policy/read", "GET")
	assert.True(t, found)
}

func TestPermissionCache_InvalidateAll(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	cache := NewPermissionCache(client, 5*time.Minute)
	ctx := context.Background()

	// Set multiple permissions
	_ = cache.Set(ctx, "user:u1", "system:root", "svc:policy/read", "GET", true)
	_ = cache.Set(ctx, "user:u2", "customer:t1", "svc:claim/read", "GET", true)

	// Invalidate all
	err := cache.InvalidateAll(ctx)
	require.NoError(t, err)

	// All permissions should be gone
	_, found := cache.Get(ctx, "user:u1", "system:root", "svc:policy/read", "GET")
	assert.False(t, found)

	_, found = cache.Get(ctx, "user:u2", "customer:t1", "svc:claim/read", "GET")
	assert.False(t, found)
}

func TestPermissionCache_NilRedis(t *testing.T) {
	cache := NewPermissionCache(nil, 5*time.Minute)
	ctx := context.Background()

	// Should not panic with nil Redis
	decision, found := cache.Get(ctx, "user:u1", "system:root", "svc:policy/read", "GET")
	assert.False(t, found)
	assert.False(t, decision)

	err := cache.Set(ctx, "user:u1", "system:root", "svc:policy/read", "GET", true)
	assert.NoError(t, err)

	err = cache.InvalidateUser(ctx, "user:u1")
	assert.NoError(t, err)
}
