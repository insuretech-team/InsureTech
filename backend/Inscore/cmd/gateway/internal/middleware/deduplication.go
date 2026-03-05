package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"sync"
	"time"
)

// DedupEntry represents a deduplicated request
type DedupEntry struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Timestamp  time.Time
}

// Deduplicator prevents duplicate request processing using idempotency keys
type Deduplicator struct {
	cache   map[string]*DedupEntry
	mu      sync.RWMutex
	ttl     time.Duration
	maxSize int
	hits    int64
	total   int64
}

// NewDeduplicator creates a new request deduplicator
func NewDeduplicator(ttl time.Duration, maxSize int) *Deduplicator {
	d := &Deduplicator{
		cache:   make(map[string]*DedupEntry),
		ttl:     ttl,
		maxSize: maxSize,
	}

	// Start cleanup goroutine
	go d.cleanupExpired()

	return d
}

// cleanupExpired periodically removes expired entries
func (d *Deduplicator) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		d.mu.Lock()
		now := time.Now()
		for key, entry := range d.cache {
			if now.Sub(entry.Timestamp) > d.ttl {
				delete(d.cache, key)
			}
		}
		d.mu.Unlock()
	}
}

// Get retrieves a deduplicated entry
func (d *Deduplicator) Get(key string) (*DedupEntry, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entry, exists := d.cache[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Since(entry.Timestamp) > d.ttl {
		return nil, false
	}

	return entry, true
}

// Set stores a deduplicated entry
func (d *Deduplicator) Set(key string, entry *DedupEntry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// If cache is full, remove oldest entry
	if len(d.cache) >= d.maxSize {
		var oldestKey string
		var oldestTime time.Time

		for k, v := range d.cache {
			if oldestKey == "" || v.Timestamp.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.Timestamp
			}
		}

		if oldestKey != "" {
			delete(d.cache, oldestKey)
		}
	}

	d.cache[key] = entry
}

// generateIdempotencyKey generates a key from request if not provided
func generateIdempotencyKey(r *http.Request, body []byte) string {
	// Combine method, path, body hash, and user ID
	hash := sha256.New()
	hash.Write([]byte(r.Method))
	hash.Write([]byte(r.URL.Path))
	hash.Write([]byte(r.URL.RawQuery))
	hash.Write(body)

	// Include user identifier if available
	if auth := r.Header.Get("Authorization"); auth != "" {
		hash.Write([]byte(auth))
	}
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		hash.Write([]byte(userID))
	}

	return hex.EncodeToString(hash.Sum(nil))
}

// shouldDeduplicate checks if request should be deduplicated
func shouldDeduplicate(r *http.Request) bool {
	// Only deduplicate mutation operations
	if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
		return false
	}

	// Deduplicate POST/PUT/PATCH/DELETE
	return true
}

// dedupResponseWriter captures response for deduplication
type dedupResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newDedupResponseWriter(w http.ResponseWriter) *dedupResponseWriter {
	return &dedupResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           &bytes.Buffer{},
	}
}

func (w *dedupResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *dedupResponseWriter) Write(b []byte) (int, error) {
	// Write to both response and buffer
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Flush implements http.Flusher
func (w *dedupResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Middleware implements request deduplication
func (d *Deduplicator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d.total++

		// Check if request should be deduplicated
		if !shouldDeduplicate(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Get idempotency key from header
		idempotencyKey := r.Header.Get("Idempotency-Key")

		// If no key provided, generate one from request content
		var requestBody []byte
		if idempotencyKey == "" {
			// Read body to generate key
			if r.Body != nil {
				var err error
				requestBody, err = io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Failed to read request body", http.StatusBadRequest)
					return
				}
				// Restore body for next handler
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			idempotencyKey = generateIdempotencyKey(r, requestBody)
		}

		// Check if request was already processed
		if entry, found := d.Get(idempotencyKey); found {
			d.hits++

			// Return cached response
			w.Header().Set("X-Idempotent-Replay", "true")
			w.Header().Set("X-Original-Timestamp", entry.Timestamp.Format(time.RFC3339))

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

		// Process new request - capture response
		w.Header().Set("X-Idempotent-Replay", "false")
		drw := newDedupResponseWriter(w)

		next.ServeHTTP(drw, r)

		// Cache successful responses (2xx)
		if drw.statusCode >= 200 && drw.statusCode < 300 {
			entry := &DedupEntry{
				StatusCode: drw.statusCode,
				Headers:    w.Header().Clone(),
				Body:       drw.body.Bytes(),
				Timestamp:  time.Now(),
			}

			d.Set(idempotencyKey, entry)
		}
	})
}

// Stats returns deduplication statistics
func (d *Deduplicator) Stats() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var hitRate float64
	if d.total > 0 {
		hitRate = float64(d.hits) / float64(d.total) * 100
	}

	return map[string]interface{}{
		"entries":  len(d.cache),
		"max_size": d.maxSize,
		"ttl":      d.ttl.String(),
		"total":    d.total,
		"hits":     d.hits,
		"hit_rate": hitRate,
	}
}

// Clear removes all entries
func (d *Deduplicator) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cache = make(map[string]*DedupEntry)
}
