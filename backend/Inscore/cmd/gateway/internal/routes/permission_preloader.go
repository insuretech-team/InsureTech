package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/metrics"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// PermissionSet holds a user's permissions for a specific portal
type PermissionSet struct {
	UserID      string          `json:"user_id"`
	Portal      string          `json:"portal"`
	TenantID    string          `json:"tenant_id"`
	Permissions map[string]bool `json:"permissions"` // object:action -> allowed
	Roles       []string        `json:"roles"`
	LoadedAt    time.Time       `json:"loaded_at"`
	ExpiresAt   time.Time       `json:"expires_at"`
}

// PermissionPreloader batches and caches user permissions for UI applications
type PermissionPreloader struct {
	authzClient authzservicev1.AuthZServiceClient
	cache       map[string]*PermissionSet
	mu          sync.RWMutex
	ttl         time.Duration
}

// NewPermissionPreloader creates a new permission pre-loader
func NewPermissionPreloader(authzConn *grpc.ClientConn, ttl time.Duration) *PermissionPreloader {
	p := &PermissionPreloader{
		authzClient: authzservicev1.NewAuthZServiceClient(authzConn),
		cache:       make(map[string]*PermissionSet),
		ttl:         ttl,
	}

	// Start cleanup goroutine
	go p.cleanupExpired()

	return p
}

func (p *PermissionPreloader) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		for key, perm := range p.cache {
			if time.Now().After(perm.ExpiresAt) {
				delete(p.cache, key)
			}
		}
		p.mu.Unlock()
	}
}

// PreloadPermissions loads all permissions for a user on login
// This is called after successful authentication to populate the permission cache
func (p *PermissionPreloader) PreloadPermissions(ctx context.Context, userID, portal, tenantID string) (*PermissionSet, error) {
	start := time.Now()

	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s:%s", userID, portal, tenantID)

	p.mu.RLock()
	cached, exists := p.cache[cacheKey]
	p.mu.RUnlock()

	if exists && time.Now().Before(cached.ExpiresAt) {
		// Cache hit
		metrics.RecordPermissionPreloadCacheHit(portal, true)
		return cached, nil
	}

	// Cache miss - load from AuthZ service
	metrics.RecordPermissionPreloadCacheHit(portal, false)
	domain := buildDomain(portal, tenantID)

	// Get user roles
	authzCtx := metadata.AppendToOutgoingContext(ctx, "x-internal-service", "gateway")

	rolesResp, err := p.authzClient.ListUserRoles(authzCtx, &authzservicev1.ListUserRolesRequest{
		UserId: userID,
		Domain: domain,
	})
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordPermissionPreload(portal, "error", duration, 0)
		return nil, fmt.Errorf("failed to list user roles: %w", err)
	}

	roles := make([]string, 0)
	if rolesResp != nil && rolesResp.Roles != nil {
		for _, role := range rolesResp.Roles {
			if role != nil {
				roles = append(roles, role.RoleId)
			}
		}
	}

	// Get user permissions
	permResp, err := p.authzClient.GetUserPermissions(authzCtx, &authzservicev1.GetUserPermissionsRequest{
		UserId: userID,
		Domain: domain,
	})
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordPermissionPreload(portal, "error", duration, 0)
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Build permission map
	permissions := make(map[string]bool)
	if permResp != nil && permResp.Permissions != nil {
		for _, perm := range permResp.Permissions {
			if perm != nil {
				// Format: object:action -> true
				key := fmt.Sprintf("%s:%s", perm.Object, perm.Action)
				permissions[key] = true
			}
		}
	}

	// Create permission set
	now := time.Now()
	permSet := &PermissionSet{
		UserID:      userID,
		Portal:      portal,
		TenantID:    tenantID,
		Permissions: permissions,
		Roles:       roles,
		LoadedAt:    now,
		ExpiresAt:   now.Add(p.ttl),
	}

	// Cache it
	p.mu.Lock()
	p.cache[cacheKey] = permSet
	p.mu.Unlock()

	logger.Info("Preloaded permissions for user",
		zap.String("user_id", userID),
		zap.String("portal", portal),
		zap.Int("permission_count", len(permissions)),
		zap.Int("role_count", len(roles)),
	)

	// Record metrics
	duration := time.Since(start).Seconds()
	metrics.RecordPermissionPreload(portal, "success", duration, len(permissions))

	return permSet, nil
}

// CheckPermission checks if a user has a specific permission (uses cache)
func (p *PermissionPreloader) CheckPermission(userID, portal, tenantID, object, action string) bool {
	cacheKey := fmt.Sprintf("%s:%s:%s", userID, portal, tenantID)

	p.mu.RLock()
	defer p.mu.RUnlock()

	permSet, exists := p.cache[cacheKey]
	if !exists || time.Now().After(permSet.ExpiresAt) {
		return false
	}

	permKey := fmt.Sprintf("%s:%s", object, action)
	return permSet.Permissions[permKey]
}

// GetPermissions returns the cached permission set for a user
func (p *PermissionPreloader) GetPermissions(userID, portal, tenantID string) (*PermissionSet, bool) {
	cacheKey := fmt.Sprintf("%s:%s:%s", userID, portal, tenantID)

	p.mu.RLock()
	defer p.mu.RUnlock()

	permSet, exists := p.cache[cacheKey]
	if !exists || time.Now().After(permSet.ExpiresAt) {
		return nil, false
	}

	return permSet, true
}

// InvalidateUser removes cached permissions for a user
func (p *PermissionPreloader) InvalidateUser(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Remove all entries for this user
	for key := range p.cache {
		if len(key) > 0 && key[:len(userID)] == userID {
			delete(p.cache, key)
		}
	}
}

// PermissionsHandler returns an HTTP handler that returns preloaded permissions
// This is typically called after login to send the UI all user permissions at once
func (p *PermissionPreloader) PermissionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user info from context (set by auth middleware)
		userID := r.Header.Get("X-User-ID")
		portal := r.Header.Get("X-Portal")
		tenantID := r.Header.Get("X-Tenant-ID")

		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if we need to refresh permissions
		refresh := r.URL.Query().Get("refresh") == "true"

		var permSet *PermissionSet
		var err error

		if refresh {
			// Force reload
			permSet, err = p.PreloadPermissions(r.Context(), userID, portal, tenantID)
		} else {
			// Try cache first
			var found bool
			permSet, found = p.GetPermissions(userID, portal, tenantID)
			if !found {
				permSet, err = p.PreloadPermissions(r.Context(), userID, portal, tenantID)
			}
		}

		if err != nil {
			logger.Error("Failed to load permissions", zap.Error(err), zap.String("user_id", userID))
			http.Error(w, "Failed to load permissions", http.StatusInternalServerError)
			return
		}

		// Return as JSON
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "private, max-age=300") // 5 min browser cache

		if err := json.NewEncoder(w).Encode(permSet); err != nil {
			logger.Error("Failed to encode permissions", zap.Error(err))
			http.Error(w, "Failed to encode permissions", http.StatusInternalServerError)
			return
		}
	}
}

// UIPermissionsList returns common UI permissions to check on initial load
// This helps reduce the number of individual permission checks
var UIPermissionsList = map[string][]string{
	"customer": {
		"svc:policy/read",
		"svc:policy/create",
		"svc:policy/update",
		"svc:claim/read",
		"svc:claim/create",
		"svc:payment/read",
		"svc:payment/create",
		"svc:profile/read",
		"svc:profile/update",
	},
	"agent": {
		"svc:policy/read",
		"svc:policy/create",
		"svc:policy/update",
		"svc:policy/delete",
		"svc:claim/read",
		"svc:claim/create",
		"svc:claim/update",
		"svc:customer/read",
		"svc:customer/create",
		"svc:commission/read",
	},
	"business": {
		"svc:policy/read",
		"svc:policy/create",
		"svc:policy/update",
		"svc:claim/read",
		"svc:claim/create",
		"svc:payment/read",
		"svc:employee/read",
		"svc:employee/create",
		"svc:employee/update",
	},
	"system": {
		"svc:user/read",
		"svc:user/create",
		"svc:user/update",
		"svc:user/delete",
		"svc:role/read",
		"svc:role/create",
		"svc:role/update",
		"svc:role/delete",
		"svc:policy/read",
		"svc:policy/create",
		"svc:policy/update",
		"svc:policy/delete",
		"svc:claim/read",
		"svc:claim/update",
		"svc:claim/approve",
		"svc:audit/read",
		"svc:report/read",
	},
}

// BatchCheckPermissions checks multiple permissions at once
func (p *PermissionPreloader) BatchCheckPermissions(ctx context.Context, userID, portal, tenantID string, checks []PermissionCheck) (map[string]bool, error) {
	domain := buildDomain(portal, tenantID)

	// Try to use cached permissions first
	if permSet, found := p.GetPermissions(userID, portal, tenantID); found {
		results := make(map[string]bool)
		for _, check := range checks {
			permKey := fmt.Sprintf("%s:%s", check.Object, check.Action)
			results[permKey] = permSet.Permissions[permKey]
		}
		return results, nil
	}

	// Cache miss - perform individual checks
	// Note: BatchCheckAccess is not available in the proto, so we perform individual checks
	// In production, this method can be optimized by adding BatchCheckAccess to the authz proto
	results := make(map[string]bool)
	for _, check := range checks {
		authzCtx := metadata.AppendToOutgoingContext(ctx, "x-internal-service", "gateway")
		resp, err := p.authzClient.CheckAccess(authzCtx, &authzservicev1.CheckAccessRequest{
			UserId: userID,
			Domain: domain,
			Object: check.Object,
			Action: check.Action,
		})

		permKey := fmt.Sprintf("%s:%s", check.Object, check.Action)
		if err != nil {
			logger.Warnf("Failed to check permission %s: %v", permKey, err)
			results[permKey] = false
		} else if resp != nil {
			results[permKey] = resp.Allowed
		} else {
			results[permKey] = false
		}
	}

	return results, nil
}

// PermissionCheck represents a single permission to check
type PermissionCheck struct {
	Object string `json:"object"`
	Action string `json:"action"`
}
