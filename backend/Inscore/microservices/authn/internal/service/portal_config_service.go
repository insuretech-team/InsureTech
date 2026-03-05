package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/metrics"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// PortalConfig holds portal-specific authentication configuration
type PortalConfig struct {
	PortalName            string // String name like "customer", "agent", etc.
	SessionTTLSeconds     int32
	RefreshTokenTTL       int32
	MFARequired           bool
	MFAMethod             string
	PasswordMinLength     int32
	PasswordRequireUpper  bool
	PasswordRequireLower  bool
	PasswordRequireDigit  bool
	PasswordRequireSymbol bool
	MaxSessionsPerUser    int32
	LoadedAt              time.Time
}

// PortalConfigCache caches portal configurations from AuthZ service
type PortalConfigCache struct {
	authzClient authzservicev1.AuthZServiceClient
	configs     map[string]*PortalConfig
	mu          sync.RWMutex
	ttl         time.Duration
}

// NewPortalConfigCache creates a new portal config cache
func NewPortalConfigCache(authzConn *grpc.ClientConn, ttl time.Duration) *PortalConfigCache {
	if authzConn == nil {
		logger.Warn("PortalConfigCache: authz connection is nil, using default configs")
		return &PortalConfigCache{
			authzClient: nil,
			configs:     make(map[string]*PortalConfig),
			ttl:         ttl,
		}
	}

	cache := &PortalConfigCache{
		authzClient: authzservicev1.NewAuthZServiceClient(authzConn),
		configs:     make(map[string]*PortalConfig),
		ttl:         ttl,
	}

	// Start background refresh
	go cache.backgroundRefresh()

	return cache
}

// Get retrieves portal configuration, loading from AuthZ if not cached
func (c *PortalConfigCache) Get(ctx context.Context, portal string) (*PortalConfig, error) {
	// Check cache first
	c.mu.RLock()
	config, exists := c.configs[portal]
	c.mu.RUnlock()

	if exists && time.Since(config.LoadedAt) < c.ttl {
		// Cache hit
		metrics.RecordPortalConfigCache(portal, true)
		return config, nil
	}

	// Cache miss or expired - load from AuthZ
	metrics.RecordPortalConfigCache(portal, false)
	return c.loadFromAuthZ(ctx, portal)
}

// loadFromAuthZ loads portal config from AuthZ service
func (c *PortalConfigCache) loadFromAuthZ(ctx context.Context, portal string) (*PortalConfig, error) {
	start := time.Now()
	
	if c.authzClient == nil {
		// Return default config if AuthZ is not available
		return c.getDefaultConfig(portal), nil
	}

	// Convert string portal to enum
	portalEnum := stringToPortalEnum(portal)
	
	resp, err := c.authzClient.GetPortalConfig(ctx, &authzservicev1.GetPortalConfigRequest{
		Portal: portalEnum,
	})
	
	// Record load duration
	duration := time.Since(start).Seconds()
	metrics.RecordPortalConfigLoad(portal, duration)

	if err != nil {
		logger.Warn("Failed to load portal config from AuthZ, using defaults",
			zap.String("portal", portal),
			zap.Error(err),
		)
		return c.getDefaultConfig(portal), nil
	}

	if resp.Error != nil && resp.Error.Code != "" {
		logger.Warn("AuthZ returned error for portal config, using defaults",
			zap.String("portal", portal),
			zap.String("error", resp.Error.Message),
		)
		return c.getDefaultConfig(portal), nil
	}

	if resp.Config == nil {
		logger.Warn("AuthZ returned empty config, using defaults",
			zap.String("portal", portal),
		)
		return c.getDefaultConfig(portal), nil
	}

	// Determine MFA method (take first from list or default to TOTP)
	mfaMethod := "TOTP"
	if len(resp.Config.MfaMethods) > 0 {
		mfaMethod = resp.Config.MfaMethods[0]
	}

	// Map AuthZ config to our internal format
	config := &PortalConfig{
		PortalName:            portal,
		SessionTTLSeconds:     resp.Config.SessionTtlSeconds,
		RefreshTokenTTL:       resp.Config.RefreshTokenTtlSeconds,
		MFARequired:           resp.Config.MfaRequired,
		MFAMethod:             mfaMethod,
		PasswordMinLength:     8,  // Default - not in proto yet
		PasswordRequireUpper:  true,
		PasswordRequireLower:  true,
		PasswordRequireDigit:  true,
		PasswordRequireSymbol: false,
		MaxSessionsPerUser:    resp.Config.MaxConcurrentSessions,
		LoadedAt:              time.Now(),
	}

	// Cache it
	c.mu.Lock()
	c.configs[portal] = config
	c.mu.Unlock()

	logger.Info("Loaded portal config from AuthZ",
		zap.String("portal", portal),
		zap.Int32("session_ttl", config.SessionTTLSeconds),
		zap.Bool("mfa_required", config.MFARequired),
	)

	return config, nil
}

// getDefaultConfig returns default configuration for a portal
func (c *PortalConfigCache) getDefaultConfig(portal string) *PortalConfig {
	// Default configurations per portal
	defaults := map[string]*PortalConfig{
		"customer": {
			PortalName:            "customer",
			SessionTTLSeconds:     3600,  // 1 hour
			RefreshTokenTTL:       86400, // 24 hours
			MFARequired:           false,
			MFAMethod:             "TOTP",
			PasswordMinLength:     8,
			PasswordRequireUpper:  true,
			PasswordRequireLower:  true,
			PasswordRequireDigit:  true,
			PasswordRequireSymbol: false,
			MaxSessionsPerUser:    5,
			LoadedAt:              time.Now(),
		},
		"agent": {
			PortalName:            "agent",
			SessionTTLSeconds:     7200,  // 2 hours
			RefreshTokenTTL:       86400, // 24 hours
			MFARequired:           true,
			MFAMethod:             "TOTP",
			PasswordMinLength:     10,
			PasswordRequireUpper:  true,
			PasswordRequireLower:  true,
			PasswordRequireDigit:  true,
			PasswordRequireSymbol: true,
			MaxSessionsPerUser:    3,
			LoadedAt:              time.Now(),
		},
		"business": {
			PortalName:            "business",
			SessionTTLSeconds:     7200,  // 2 hours
			RefreshTokenTTL:       86400, // 24 hours
			MFARequired:           false,
			MFAMethod:             "TOTP",
			PasswordMinLength:     8,
			PasswordRequireUpper:  true,
			PasswordRequireLower:  true,
			PasswordRequireDigit:  true,
			PasswordRequireSymbol: false,
			MaxSessionsPerUser:    10,
			LoadedAt:              time.Now(),
		},
		"system": {
			PortalName:            "system",
			SessionTTLSeconds:     14400, // 4 hours
			RefreshTokenTTL:       86400, // 24 hours
			MFARequired:           true,
			MFAMethod:             "TOTP",
			PasswordMinLength:     12,
			PasswordRequireUpper:  true,
			PasswordRequireLower:  true,
			PasswordRequireDigit:  true,
			PasswordRequireSymbol: true,
			MaxSessionsPerUser:    2,
			LoadedAt:              time.Now(),
		},
	}

	if config, exists := defaults[portal]; exists {
		return config
	}

	// Generic default for unknown portals
	return &PortalConfig{
		PortalName:            portal,
		SessionTTLSeconds:     3600,
		RefreshTokenTTL:       86400,
		MFARequired:           false,
		MFAMethod:             "TOTP",
		PasswordMinLength:     8,
		PasswordRequireUpper:  true,
		PasswordRequireLower:  true,
		PasswordRequireDigit:  true,
		PasswordRequireSymbol: false,
		MaxSessionsPerUser:    5,
		LoadedAt:              time.Now(),
	}
}

// Invalidate removes a portal config from cache (called when portal config is updated)
func (c *PortalConfigCache) Invalidate(portal string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.configs, portal)
	logger.Info("Invalidated portal config cache", zap.String("portal", portal))
}

// InvalidateAll clears all cached configs
func (c *PortalConfigCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.configs = make(map[string]*PortalConfig)
	logger.Info("Invalidated all portal config cache")
}

// backgroundRefresh periodically refreshes cached configs
func (c *PortalConfigCache) backgroundRefresh() {
	ticker := time.NewTicker(c.ttl / 2) // Refresh at half TTL
	defer ticker.Stop()

	for range ticker.C {
		c.mu.RLock()
		portals := make([]string, 0, len(c.configs))
		for portal := range c.configs {
			portals = append(portals, portal)
		}
		c.mu.RUnlock()

		// Refresh each portal config in background
		for _, portal := range portals {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := c.loadFromAuthZ(ctx, portal)
			if err != nil {
				logger.Warn("Background refresh failed for portal",
					zap.String("portal", portal),
					zap.Error(err),
				)
			}
			cancel()
		}
	}
}

// ValidatePassword validates a password against portal-specific policy
func (c *PortalConfigCache) ValidatePassword(ctx context.Context, portal, password string) error {
	config, err := c.Get(ctx, portal)
	if err != nil {
		return fmt.Errorf("failed to get portal config: %w", err)
	}

	if len(password) < int(config.PasswordMinLength) {
		metrics.RecordPasswordValidationFailure(portal, "too_short")
		return fmt.Errorf("password must be at least %d characters long", config.PasswordMinLength)
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSymbol := false

	for _, ch := range password {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		case (ch >= '!' && ch <= '/') || (ch >= ':' && ch <= '@') || (ch >= '[' && ch <= '`') || (ch >= '{' && ch <= '~'):
			hasSymbol = true
		}
	}

	if config.PasswordRequireUpper && !hasUpper {
		metrics.RecordPasswordValidationFailure(portal, "no_uppercase")
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if config.PasswordRequireLower && !hasLower {
		metrics.RecordPasswordValidationFailure(portal, "no_lowercase")
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if config.PasswordRequireDigit && !hasDigit {
		metrics.RecordPasswordValidationFailure(portal, "no_digit")
		return fmt.Errorf("password must contain at least one digit")
	}

	if config.PasswordRequireSymbol && !hasSymbol {
		metrics.RecordPasswordValidationFailure(portal, "no_symbol")
		return fmt.Errorf("password must contain at least one symbol")
	}

	return nil
}

// GetSessionTTL returns the session TTL for a portal
func (c *PortalConfigCache) GetSessionTTL(ctx context.Context, portal string) (time.Duration, error) {
	config, err := c.Get(ctx, portal)
	if err != nil {
		return 0, err
	}
	return time.Duration(config.SessionTTLSeconds) * time.Second, nil
}

// GetRefreshTokenTTL returns the refresh token TTL for a portal
func (c *PortalConfigCache) GetRefreshTokenTTL(ctx context.Context, portal string) (time.Duration, error) {
	config, err := c.Get(ctx, portal)
	if err != nil {
		return 0, err
	}
	return time.Duration(config.RefreshTokenTTL) * time.Second, nil
}

// IsMFARequired checks if MFA is required for a portal
func (c *PortalConfigCache) IsMFARequired(ctx context.Context, portal string) (bool, error) {
	config, err := c.Get(ctx, portal)
	if err != nil {
		return false, err
	}
	return config.MFARequired, nil
}

// stringToPortalEnum converts a portal string to Portal enum
func stringToPortalEnum(portal string) authzentityv1.Portal {
	switch portal {
	case "customer", "b2c":
		return authzentityv1.Portal_PORTAL_B2C
	case "agent":
		return authzentityv1.Portal_PORTAL_AGENT
	case "business":
		return authzentityv1.Portal_PORTAL_BUSINESS
	case "partner", "b2b":
		return authzentityv1.Portal_PORTAL_B2B
	case "system":
		return authzentityv1.Portal_PORTAL_SYSTEM
	case "regulator":
		return authzentityv1.Portal_PORTAL_REGULATOR
	default:
		return authzentityv1.Portal_PORTAL_UNSPECIFIED
	}
}
