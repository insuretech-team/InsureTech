package middleware

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type rlStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *rlStream) Context() context.Context { return s.ctx }

func TestRateLimiter_AllowAndInterceptors(t *testing.T) {
	rl := NewRateLimiter(1, 2)
	defer rl.Stop()
	rl.cleanup = 20 * time.Millisecond

	require.True(t, rl.allow("127.0.0.1"))
	require.True(t, rl.allow("127.0.0.1"))
	require.False(t, rl.allow("127.0.0.1"))

	ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("10.0.0.1"), Port: 9090}})
	_, err := rl.UnaryInterceptor()(ctx, "req", &grpc.UnaryServerInfo{FullMethod: "/m"}, func(ctx context.Context, req interface{}) (interface{}, error) { return "ok", nil })
	require.NoError(t, err)
	_, err = rl.UnaryInterceptor()(ctx, "req", &grpc.UnaryServerInfo{FullMethod: "/m"}, func(ctx context.Context, req interface{}) (interface{}, error) { return "ok", nil })
	require.NoError(t, err)
	_, err = rl.UnaryInterceptor()(ctx, "req", &grpc.UnaryServerInfo{FullMethod: "/m"}, func(ctx context.Context, req interface{}) (interface{}, error) { return "ok", nil })
	require.Equal(t, codes.ResourceExhausted, status.Code(err))

	err = rl.StreamInterceptor()("srv", &rlStream{ctx: ctx}, &grpc.StreamServerInfo{FullMethod: "/m"}, func(srv interface{}, ss grpc.ServerStream) error { return nil })
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestRateLimiter_CleanupLoop(t *testing.T) {
	rl := NewRateLimiter(1, 1)
	rl.cleanup = 10 * time.Millisecond
	rl.mu.Lock()
	rl.buckets["stale"] = &tokenBucket{tokens: 0, lastSeen: time.Now().Add(-time.Second)}
	rl.mu.Unlock()
	time.Sleep(40 * time.Millisecond)
	rl.mu.Lock()
	_, ok := rl.buckets["stale"]
	rl.mu.Unlock()
	rl.Stop()
	require.False(t, ok)
}
