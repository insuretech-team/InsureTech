package routes

import (
	"testing"
	"time"
)

func TestCombinedAuthCache_SetAndGet(t *testing.T) {
	cache := NewCombinedAuthCache(5 * time.Minute)
	
	result := &CombinedAuthResult{
		UserID:      "user123",
		SessionID:   "session456",
		Valid:       true,
		Allowed:     true,
		MatchedRule: "policy1",
	}
	
	key := "test_key_123"
	
	// Set the result
	cache.Set(key, result)
	
	// Get it back
	cached, found := cache.Get(key)
	
	if !found {
		t.Error("Expected to find cached result")
	}
	
	if cached.UserID != result.UserID {
		t.Errorf("UserID = %v, want %v", cached.UserID, result.UserID)
	}
	
	if cached.SessionID != result.SessionID {
		t.Errorf("SessionID = %v, want %v", cached.SessionID, result.SessionID)
	}
	
	if cached.Valid != result.Valid {
		t.Errorf("Valid = %v, want %v", cached.Valid, result.Valid)
	}
	
	if cached.Allowed != result.Allowed {
		t.Errorf("Allowed = %v, want %v", cached.Allowed, result.Allowed)
	}
}

func TestCombinedAuthCache_Expiry(t *testing.T) {
	cache := NewCombinedAuthCache(100 * time.Millisecond) // Very short TTL
	
	result := &CombinedAuthResult{
		UserID:  "user123",
		Valid:   true,
		Allowed: true,
	}
	
	key := "test_key_expiry"
	
	// Set the result
	cache.Set(key, result)
	
	// Should be found immediately
	_, found := cache.Get(key)
	if !found {
		t.Error("Expected to find cached result immediately")
	}
	
	// Wait for expiry
	time.Sleep(150 * time.Millisecond)
	
	// Should not be found after expiry
	_, found = cache.Get(key)
	if found {
		t.Error("Expected cached result to be expired")
	}
}

func TestCombinedAuthCache_Invalidate(t *testing.T) {
	cache := NewCombinedAuthCache(5 * time.Minute)
	
	// Set multiple results for same user
	cache.Set("user123:policy:read", &CombinedAuthResult{UserID: "user123", Valid: true, Allowed: true})
	cache.Set("user123:claim:write", &CombinedAuthResult{UserID: "user123", Valid: true, Allowed: true})
	cache.Set("user456:policy:read", &CombinedAuthResult{UserID: "user456", Valid: true, Allowed: true})
	
	// Verify all are cached
	if _, found := cache.Get("user123:policy:read"); !found {
		t.Error("Expected to find user123:policy:read")
	}
	if _, found := cache.Get("user123:claim:write"); !found {
		t.Error("Expected to find user123:claim:write")
	}
	if _, found := cache.Get("user456:policy:read"); !found {
		t.Error("Expected to find user456:policy:read")
	}
	
	// Invalidate user123
	cache.Invalidate("user123")
	
	// user123 entries should be gone
	if _, found := cache.Get("user123:policy:read"); found {
		t.Error("Expected user123:policy:read to be invalidated")
	}
	if _, found := cache.Get("user123:claim:write"); found {
		t.Error("Expected user123:claim:write to be invalidated")
	}
	
	// user456 should still be there
	if _, found := cache.Get("user456:policy:read"); !found {
		t.Error("Expected user456:policy:read to still exist")
	}
}

func TestPermissionPreloader_CheckPermission(t *testing.T) {
	preloader := &PermissionPreloader{
		cache: make(map[string]*PermissionSet),
	}
	
	// Add a permission set
	now := time.Now()
	permSet := &PermissionSet{
		UserID:   "user123",
		Portal:   "customer",
		TenantID: "tenant456",
		Permissions: map[string]bool{
			"svc:policy:read":  true,  // Format: object:action
			"svc:policy:write": true,
			"svc:claim:read":   true,
		},
		LoadedAt:  now,
		ExpiresAt: now.Add(10 * time.Minute),
	}
	
	cacheKey := "user123:customer:tenant456"
	preloader.cache[cacheKey] = permSet
	
	tests := []struct {
		name     string
		object   string
		action   string
		expected bool
	}{
		{
			name:     "has read permission",
			object:   "svc:policy",
			action:   "read",
			expected: true,
		},
		{
			name:     "has write permission",
			object:   "svc:policy",
			action:   "write",
			expected: true,
		},
		{
			name:     "no delete permission",
			object:   "svc:policy",
			action:   "delete",
			expected: false,
		},
		{
			name:     "has claim read permission",
			object:   "svc:claim",
			action:   "read",
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := preloader.CheckPermission("user123", "customer", "tenant456", tt.object, tt.action)
			
			if allowed != tt.expected {
				t.Errorf("CheckPermission() = %v, want %v", allowed, tt.expected)
			}
		})
	}
}

func TestPermissionPreloader_InvalidateUser(t *testing.T) {
	preloader := &PermissionPreloader{
		cache: make(map[string]*PermissionSet),
	}
	
	// Add permission sets for multiple users
	preloader.cache["user123:customer:tenant1"] = &PermissionSet{UserID: "user123"}
	preloader.cache["user123:agent:tenant1"] = &PermissionSet{UserID: "user123"}
	preloader.cache["user456:customer:tenant1"] = &PermissionSet{UserID: "user456"}
	
	if len(preloader.cache) != 3 {
		t.Fatalf("Expected 3 cached items, got %d", len(preloader.cache))
	}
	
	// Invalidate user123
	preloader.InvalidateUser("user123")
	
	// user123 entries should be removed
	if _, exists := preloader.cache["user123:customer:tenant1"]; exists {
		t.Error("Expected user123:customer:tenant1 to be invalidated")
	}
	if _, exists := preloader.cache["user123:agent:tenant1"]; exists {
		t.Error("Expected user123:agent:tenant1 to be invalidated")
	}
	
	// user456 should still exist
	if _, exists := preloader.cache["user456:customer:tenant1"]; !exists {
		t.Error("Expected user456:customer:tenant1 to still exist")
	}
}
