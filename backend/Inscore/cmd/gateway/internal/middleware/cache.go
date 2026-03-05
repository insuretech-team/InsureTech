package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CacheEntry represents a cached response
type CacheEntry struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	CachedAt   time.Time
	ExpiresAt  time.Time
}

// IsExpired checks if cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// ResponseCache manages response caching
type ResponseCache struct {
	cache      map[string]*CacheEntry
	mu         sync.RWMutex
	defaultTTL time.Duration
	maxSize    int // Maximum cache entries
}

// NewResponseCache creates a new response cache
func NewResponseCache(defaultTTL time.Duration, maxSize int) *ResponseCache {
	rc := &ResponseCache{
		cache:      make(map[string]*CacheEntry),
		defaultTTL: defaultTTL,
		maxSize:    maxSize,
	}

	// Start cleanup goroutine
	go rc.cleanupExpired()

	return rc
}

// cleanupExpired periodically removes expired entries
func (rc *ResponseCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rc.mu.Lock()
		for key, entry := range rc.cache {
			if entry.IsExpired() {
				delete(rc.cache, key)
			}
		}
		rc.mu.Unlock()
	}
}

// Get retrieves a cached entry
func (rc *ResponseCache) Get(key string) (*CacheEntry, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	entry, exists := rc.cache[key]
	if !exists || entry.IsExpired() {
		return nil, false
	}

	return entry, true
}

// Set stores a cache entry
func (rc *ResponseCache) Set(key string, entry *CacheEntry) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// If cache is full, remove oldest entry (simple LRU)
	if len(rc.cache) >= rc.maxSize {
		var oldestKey string
		var oldestTime time.Time

		for k, v := range rc.cache {
			if oldestKey == "" || v.CachedAt.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.CachedAt
			}
		}

		if oldestKey != "" {
			delete(rc.cache, oldestKey)
		}
	}

	rc.cache[key] = entry
}

// Clear removes all cache entries
func (rc *ResponseCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.cache = make(map[string]*CacheEntry)
}

// Stats returns cache statistics
func (rc *ResponseCache) Stats() map[string]interface{} {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return map[string]interface{}{
		"entries":     len(rc.cache),
		"max_size":    rc.maxSize,
		"default_ttl": rc.defaultTTL.String(),
	}
}

// generateCacheKey creates a unique cache key from request
func generateCacheKey(r *http.Request) string {
	// Include method, path, query, and authorization header
	key := r.Method + ":" + r.URL.Path + ":" + r.URL.RawQuery

	// Include auth header for user-specific caching
	if auth := r.Header.Get("Authorization"); auth != "" {
		hash := sha256.Sum256([]byte(auth))
		key += ":" + hex.EncodeToString(hash[:8])
	}

	return key
}

// shouldCache determines if a request should be cached
func shouldCache(r *http.Request) bool {
	// Only cache GET requests
	if r.Method != http.MethodGet {
		return false
	}

	// Don't cache if client explicitly requests fresh data
	if r.Header.Get("Cache-Control") == "no-cache" {
		return false
	}

	// Cache these paths
	cachePaths := []string{
		"/v1/products",            // Product listings
		"/v1/products/",           // Individual products
		"/v1/products/categories", // Categories
		"/v1/products/brands",     // Brands
		"/v1/shipping/carriers",   // Carriers list
	}

	for _, path := range cachePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return true
		}
	}

	return false
}

// getTTLForPath returns cache TTL based on path
func getTTLForPath(path string) time.Duration {
	// Different TTLs for different endpoints
	switch {
	case strings.HasPrefix(path, "/v1/products/categories"):
		return 10 * time.Minute // Categories rarely change
	case strings.HasPrefix(path, "/v1/products/brands"):
		return 10 * time.Minute // Brands rarely change
	case strings.HasPrefix(path, "/v1/shipping/carriers"):
		return 15 * time.Minute // Carriers rarely change
	case strings.Contains(path, "/v1/products/") && !strings.Contains(path, "?"):
		return 5 * time.Minute // Individual products
	case strings.HasPrefix(path, "/v1/products"):
		return 2 * time.Minute // Product listings (may have inventory changes)
	default:
		return 1 * time.Minute // Default TTL
	}
}

// cachedResponseWriter captures response for caching
type cachedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
	tee        *bytes.Buffer
}

func newCachedResponseWriter(w http.ResponseWriter) *cachedResponseWriter {
	buf := &bytes.Buffer{}
	return &cachedResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           buf,
		tee:            buf,
	}
}

func (w *cachedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *cachedResponseWriter) Write(b []byte) (int, error) {
	// Write to both response and buffer
	w.tee.Write(b)
	return w.ResponseWriter.Write(b)
}

// Flush implements http.Flusher
func (w *cachedResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// CacheMiddleware provides response caching
func (rc *ResponseCache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request should be cached
		if !shouldCache(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Generate cache key
		cacheKey := generateCacheKey(r)

		// Try to get from cache
		if entry, found := rc.Get(cacheKey); found {
			// Cache hit!
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("X-Cache-Expires", entry.ExpiresAt.Format(time.RFC3339))

			// Copy cached headers
			for key, values := range entry.Headers {
				for _, value := range values {
					w.Header().Set(key, value)
				}
			}

			w.WriteHeader(entry.StatusCode)
			w.Write(entry.Body)
			return
		}

		// Cache miss - capture response
		w.Header().Set("X-Cache", "MISS")
		crw := newCachedResponseWriter(w)

		next.ServeHTTP(crw, r)

		// Cache successful responses
		if crw.statusCode == http.StatusOK {
			ttl := getTTLForPath(r.URL.Path)

			entry := &CacheEntry{
				StatusCode: crw.statusCode,
				Headers:    w.Header().Clone(),
				Body:       crw.body.Bytes(),
				CachedAt:   time.Now(),
				ExpiresAt:  time.Now().Add(ttl),
			}

			rc.Set(cacheKey, entry)
		}
	})
}

// CacheInvalidationMiddleware provides cache invalidation on mutations
func (rc *ResponseCache) InvalidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clear cache on POST/PUT/DELETE/PATCH (mutations)
		if r.Method == http.MethodPost || r.Method == http.MethodPut ||
			r.Method == http.MethodDelete || r.Method == http.MethodPatch {

			// Clear relevant cache entries based on path
			if strings.HasPrefix(r.URL.Path, "/v1/products") {
				// In a real implementation, you'd selectively invalidate
				// For now, clear all cache on any product mutation
				rc.Clear()
				w.Header().Set("X-Cache-Invalidated", "true")
			}
		}

		next.ServeHTTP(w, r)
	})
}
