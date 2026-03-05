package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestRequestIDHelpers(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-request-id", "rid-1"))
	require.Equal(t, "rid-1", requestIDFromIncomingMD(ctx))

	ctx2, rid := withRequestID(ctx)
	require.Equal(t, "rid-1", rid)
	require.Equal(t, "rid-1", getRequestID(ctx2))
	require.Equal(t, "unknown", formatFullMethod(""))
	require.Equal(t, "/svc/m", formatFullMethod("/svc/m"))
}

func TestRecoveryAndRequestIDUnaryInterceptors(t *testing.T) {
	rec := recoveryUnaryInterceptor()
	_, err := rec(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/x"}, func(ctx context.Context, req any) (any, error) {
		panic("boom")
	})
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))

	reqIDInt := requestIDUnaryInterceptor()
	_, err = reqIDInt(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/x"}, func(ctx context.Context, req any) (any, error) {
		require.NotEmpty(t, getRequestID(ctx))
		return "ok", nil
	})
	require.NoError(t, err)

	logInt := loggingUnaryInterceptor()
	_, err = logInt(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/x"}, func(ctx context.Context, req any) (any, error) {
		return nil, errors.New("bad")
	})
	require.Error(t, err)
}

func TestRateLimitAndMetricsInterceptors_NoRedis(t *testing.T) {
	rl := rateLimitUnaryInterceptor(nil)
	_, err := rl(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/insuretech.authn.services.v1.AuthService/Login"}, func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	})
	require.NoError(t, err)

	m := MetricsUnaryInterceptor()
	_, err = m(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/x"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.Error(codes.InvalidArgument, "bad")
	})
	require.Error(t, err)
}

func TestDefaultInterceptorsBuilders(t *testing.T) {
	require.NotEmpty(t, defaultUnaryInterceptors())
	require.NotEmpty(t, unaryInterceptorsWithRateLimit(nil))
	require.NotEmpty(t, defaultStreamInterceptors())
}
