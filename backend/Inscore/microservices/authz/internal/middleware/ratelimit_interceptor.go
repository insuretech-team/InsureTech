package middleware

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// tokenBucket holds token-bucket state for a single client IP.
type tokenBucket struct {
	tokens   float64
	lastSeen time.Time
}

// RateLimiter is a token-bucket rate limiter keyed by client IP.
// It is safe for concurrent use.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64       // tokens added per second
	capacity float64       // maximum burst capacity
	cleanup  time.Duration // how often to evict idle buckets
	stopCh   chan struct{}
}

// NewRateLimiter creates a RateLimiter with the given steady-state rate (rps)
// and burst capacity. A background goroutine cleans up idle buckets every minute.
func NewRateLimiter(rate, capacity int) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     float64(rate),
		capacity: float64(capacity),
		cleanup:  time.Minute,
		stopCh:   make(chan struct{}),
	}
	go rl.cleanupLoop()
	return rl
}

// Stop terminates the background cleanup goroutine.
func (r *RateLimiter) Stop() {
	close(r.stopCh)
}

// allow implements the token-bucket algorithm for the given IP.
// Returns true if the request should be allowed, false if rate-limited.
func (r *RateLimiter) allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	b, exists := r.buckets[ip]
	if !exists {
		// First request from this IP — start with a full bucket.
		r.buckets[ip] = &tokenBucket{
			tokens:   r.capacity - 1, // consume one immediately
			lastSeen: now,
		}
		return true
	}

	// Refill tokens based on elapsed time.
	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens += elapsed * r.rate
	if b.tokens > r.capacity {
		b.tokens = r.capacity
	}
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// cleanupLoop removes buckets that have been idle for longer than the cleanup interval.
func (r *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(r.cleanup)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			cutoff := time.Now().Add(-r.cleanup)
			for ip, b := range r.buckets {
				if b.lastSeen.Before(cutoff) {
					delete(r.buckets, ip)
				}
			}
			r.mu.Unlock()
		case <-r.stopCh:
			return
		}
	}
}

// UnaryInterceptor returns a gRPC unary server interceptor that enforces rate limits.
func (r *RateLimiter) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ip := peerIP(ctx)
		if !r.allow(ip) {
			return nil, status.Errorf(codes.ResourceExhausted,
				"rate limit exceeded — too many requests from %s", ip)
		}
		return handler(ctx, req)
	}
}

// StreamInterceptor returns a gRPC stream server interceptor that enforces rate limits.
func (r *RateLimiter) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ip := peerIP(ss.Context())
		if !r.allow(ip) {
			return status.Errorf(codes.ResourceExhausted,
				"rate limit exceeded — too many requests from %s", ip)
		}
		return handler(srv, ss)
	}
}

// peerIP extracts the client IP address from the gRPC peer info in the context.
// Falls back to "unknown" if unavailable.
func peerIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p == nil || p.Addr == nil {
		return "unknown"
	}
	addr := p.Addr.String()
	// Strip port suffix if present (host:port → host).
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i]
		}
	}
	return addr
}
