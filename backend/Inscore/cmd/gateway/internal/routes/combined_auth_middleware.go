package routes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/metrics"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CombinedAuthResult holds the result of AuthN + AuthZ validation
type CombinedAuthResult struct {
	// AuthN fields
	UserID      string
	SessionID   string
	SessionType string
	UserType    string
	Portal      string
	TenantID    string
	TokenID     string
	DeviceID    string
	Valid       bool
	
	// AuthZ fields
	Allowed      bool
	MatchedRule  string
	Reason       string
	
	// Cache metadata
	CachedAt time.Time
}

// CombinedAuthCache caches auth validation results
type CombinedAuthCache struct {
	cache map[string]*CombinedAuthResult
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewCombinedAuthCache creates a new combined auth cache
func NewCombinedAuthCache(ttl time.Duration) *CombinedAuthCache {
	c := &CombinedAuthCache{
		cache: make(map[string]*CombinedAuthResult),
		ttl:   ttl,
	}
	
	// Start cleanup goroutine
	go c.cleanupExpired()
	
	return c
}

func (c *CombinedAuthCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		for key, result := range c.cache {
			if time.Since(result.CachedAt) > c.ttl {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *CombinedAuthCache) Get(key string) (*CombinedAuthResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result, exists := c.cache[key]
	if !exists {
		return nil, false
	}
	
	// Check if expired
	if time.Since(result.CachedAt) > c.ttl {
		return nil, false
	}
	
	return result, true
}

func (c *CombinedAuthCache) Set(key string, result *CombinedAuthResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	result.CachedAt = time.Now()
	c.cache[key] = result
}

func (c *CombinedAuthCache) Invalidate(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Remove all entries for this user
	for key := range c.cache {
		if strings.Contains(key, userID) {
			delete(c.cache, key)
		}
	}
}

// CombinedAuthMiddleware performs both AuthN and AuthZ in a single optimized flow
// with caching and circuit breaker support
type CombinedAuthMiddleware struct {
	authnClient      authnservicev1.AuthServiceClient
	authzClient      authzservicev1.AuthZServiceClient
	cache            *CombinedAuthCache
	authnCircuit     *middleware.CircuitBreaker
	authzCircuit     *middleware.CircuitBreaker
	servicePrefix    string
	extractResource  ResourceExtractorFn
}

// NewCombinedAuthMiddleware creates an optimized combined auth middleware
func NewCombinedAuthMiddleware(
	authnConn *grpc.ClientConn,
	authzConn *grpc.ClientConn,
	servicePrefix string,
	extractResource ResourceExtractorFn,
) *CombinedAuthMiddleware {
	return &CombinedAuthMiddleware{
		authnClient:     authnservicev1.NewAuthServiceClient(authnConn),
		authzClient:     authzservicev1.NewAuthZServiceClient(authzConn),
		cache:           NewCombinedAuthCache(30 * time.Second), // 30s TTL as per plan
		authnCircuit:    middleware.NewCircuitBreaker("authn", 5, 10*time.Second),
		authzCircuit:    middleware.NewCircuitBreaker("authz", 5, 10*time.Second),
		servicePrefix:   servicePrefix,
		extractResource: extractResource,
	}
}

// Middleware returns the HTTP middleware handler
func (m *CombinedAuthMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			
			// Extract token/session
			jwt := bearerToken(r.Header.Get("Authorization"))
			sessionToken := ""
			if c, err := r.Cookie(SessionCookieName); err == nil {
				sessionToken = c.Value
			}
			
			// Build resource info for AuthZ
			resource := ""
			if m.extractResource != nil {
				resource = m.extractResource(r)
			}
			object := buildObject(m.servicePrefix, resource)
			action := r.Method
			
			// Generate cache key
			cacheKey := m.buildCacheKey(jwt, sessionToken, object, action)
			
			// Check cache first
			if cached, found := m.cache.Get(cacheKey); found {
				if cached.Valid && cached.Allowed {
					logger.Debug("Combined auth cache HIT",
						zap.String("user_id", cached.UserID),
						zap.String("object", object),
						zap.String("action", action),
					)
					
					// Record cache hit
					metrics.RecordCombinedAuthCacheHit(true)
					duration := time.Since(time.Now()).Seconds() // Minimal duration for cache hit
					metrics.RecordCombinedAuthRequest(cached.Portal, "success", duration, true)
					
					m.populateContext(r, cached)
					next.ServeHTTP(w, r)
					return
				}
			}
			
			// Cache miss
			metrics.RecordCombinedAuthCacheHit(false)
			
			// Cache miss - perform full validation
			validationStart := time.Now()
			result, err := m.validateAuth(ctx, r, jwt, sessionToken, object, action)
			validationDuration := time.Since(validationStart).Seconds()
			
			if err != nil {
				logger.Error("Combined auth validation failed", zap.Error(err))
				metrics.RecordCombinedAuthRequest("unknown", "error", validationDuration, false)
				http.Error(w, "Authentication/Authorization failed", http.StatusUnauthorized)
				return
			}
			
			// Cache the result
			m.cache.Set(cacheKey, result)
			
			if !result.Valid {
				metrics.RecordCombinedAuthRequest(result.Portal, "authn_failed", validationDuration, false)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			
			if !result.Allowed {
				reason := result.Reason
				if reason == "" {
					reason = "no matching policy"
				}
				logger.Warn("AuthZ DENY",
					zap.String("user_id", result.UserID),
					zap.String("object", object),
					zap.String("action", action),
					zap.String("reason", reason),
				)
				metrics.RecordCombinedAuthRequest(result.Portal, "authz_failed", validationDuration, false)
				http.Error(w, "Forbidden: "+reason, http.StatusForbidden)
				return
			}
			
			// Success
			metrics.RecordCombinedAuthRequest(result.Portal, "success", validationDuration, false)
			
			// Populate context and headers
			m.populateContext(r, result)
			
			logger.Debug("Combined auth SUCCESS",
				zap.String("user_id", result.UserID),
				zap.String("object", object),
				zap.String("action", action),
			)
			
			next.ServeHTTP(w, r)
		})
	}
}

// validateAuth performs the actual AuthN + AuthZ validation
func (m *CombinedAuthMiddleware) validateAuth(
	ctx context.Context,
	r *http.Request,
	jwt, sessionToken, object, action string,
) (*CombinedAuthResult, error) {
	result := &CombinedAuthResult{}
	
	// Step 1: AuthN validation with circuit breaker
	var authnResp *authnservicev1.ValidateTokenResponse
	err := m.authnCircuit.Execute(func() error {
		// Prepare metadata
		md := metadata.New(map[string]string{
			"authorization":   r.Header.Get("Authorization"),
			"cookie":          r.Header.Get("Cookie"),
			"x-csrf-token":    r.Header.Get("X-CSRF-Token"),
			"x-device-id":     r.Header.Get("X-Device-Id"),
			"x-forwarded-for": r.Header.Get("X-Forwarded-For"),
			"x-real-ip":       r.Header.Get("X-Real-Ip"),
			"user-agent":      r.UserAgent(),
		})
		ctx = metadata.NewOutgoingContext(ctx, md)
		
		resp, err := m.authnClient.ValidateToken(ctx, &authnservicev1.ValidateTokenRequest{
			AccessToken: jwt,
			SessionId:   sessionToken,
		})
		authnResp = resp
		return err
	})
	
	if err != nil {
		return nil, fmt.Errorf("authn validation failed: %w", err)
	}
	
	if authnResp == nil || !authnResp.Valid {
		result.Valid = false
		return result, nil
	}
	
	// Populate AuthN fields
	result.Valid = true
	result.UserID = authnResp.UserId
	result.SessionID = authnResp.SessionId
	result.SessionType = authnResp.SessionType
	result.UserType = authnResp.UserType
	result.Portal = authnResp.Portal
	result.TenantID = authnResp.TenantId
	result.TokenID = authnResp.TokenId
	result.DeviceID = authnResp.DeviceId
	
	// Device binding check
	requestDeviceID := strings.TrimSpace(r.Header.Get("X-Device-Id"))
	if requestDeviceID != "" && authnResp.SessionType == "JWT" && authnResp.DeviceId != "" {
		if requestDeviceID != authnResp.DeviceId {
			result.Valid = false
			result.Reason = "device mismatch"
			return result, nil
		}
	}
	
	// Step 2: AuthZ validation with circuit breaker
	domain := buildDomain(result.Portal, result.TenantID)
	
	// Build access context with API key scopes if present
	accessCtx := &authzservicev1.AccessContext{
		SessionId: result.SessionID,
		TokenId:   result.TokenID,
		DeviceId:  result.DeviceID,
		IpAddress: realIP(r),
		UserAgent: r.UserAgent(),
	}
	
	// Check if this is an API key authentication and pass scopes to AuthZ
	if authnResp.SessionType == "API_KEY" && len(authnResp.ApiKeyScopes) > 0 {
		if accessCtx.Attributes == nil {
			accessCtx.Attributes = make(map[string]string)
		}
		
		// Pass API key scopes as comma-separated string
		// The AuthZ service will validate these scopes before checking Casbin policies
		accessCtx.Attributes["api_key_scopes"] = strings.Join(authnResp.ApiKeyScopes, ",")
		accessCtx.Attributes["auth_type"] = "api_key"
		
		logger.Debug("API key authentication detected",
			zap.String("user_id", result.UserID),
			zap.Int("scope_count", len(authnResp.ApiKeyScopes)),
		)
	}
	
	var authzResp *authzservicev1.CheckAccessResponse
	err = m.authzCircuit.Execute(func() error {
		resp, err := m.authzClient.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
			UserId: result.UserID,
			Domain: domain,
			Object: object,
			Action: action,
			Context: accessCtx,
		})
		authzResp = resp
		return err
	})
	
	if err != nil {
		return nil, fmt.Errorf("authz validation failed: %w", err)
	}
	
	if authzResp != nil {
		result.Allowed = authzResp.Allowed
		result.MatchedRule = authzResp.MatchedRule
		result.Reason = authzResp.Reason
	}
	
	return result, nil
}

// populateContext adds auth data to request context and headers
func (m *CombinedAuthMiddleware) populateContext(r *http.Request, result *CombinedAuthResult) {
	// Set headers for downstream services
	r.Header.Set("X-User-ID", result.UserID)
	r.Header.Set("X-Session-ID", result.SessionID)
	r.Header.Set("X-Session-Type", result.SessionType)
	r.Header.Set("X-User-Type", result.UserType)
	r.Header.Set("X-Portal", result.Portal)
	r.Header.Set("X-Tenant-ID", result.TenantID)
	r.Header.Set("X-Token-ID", result.TokenID)
	r.Header.Set("X-Device-ID", result.DeviceID)
	
	// Store in context
	ctx := r.Context()
	ctx = context.WithValue(ctx, "user_id", result.UserID)
	ctx = context.WithValue(ctx, "session_id", result.SessionID)
	ctx = context.WithValue(ctx, "session_type", result.SessionType)
	ctx = context.WithValue(ctx, "user_type", result.UserType)
	ctx = context.WithValue(ctx, "portal", result.Portal)
	ctx = context.WithValue(ctx, "tenant_id", result.TenantID)
	ctx = context.WithValue(ctx, "token_id", result.TokenID)
	ctx = context.WithValue(ctx, "device_id", result.DeviceID)
	
	*r = *r.WithContext(ctx)
}

// buildCacheKey generates a cache key from auth parameters
func (m *CombinedAuthMiddleware) buildCacheKey(jwt, sessionToken, object, action string) string {
	// Use token/session + object + action as cache key
	key := fmt.Sprintf("%s:%s:%s:%s", jwt, sessionToken, object, action)
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// GetCircuitBreakerStates returns the current state of circuit breakers
func (m *CombinedAuthMiddleware) GetCircuitBreakerStates() map[string]string {
	return map[string]string{
		"authn": circuitStateToString(m.authnCircuit.State()),
		"authz": circuitStateToString(m.authzCircuit.State()),
	}
}

// circuitStateToString converts circuit state to string
func circuitStateToString(s middleware.CircuitState) string {
	switch s {
	case middleware.StateClosed:
		return "CLOSED"
	case middleware.StateOpen:
		return "OPEN"
	case middleware.StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}
