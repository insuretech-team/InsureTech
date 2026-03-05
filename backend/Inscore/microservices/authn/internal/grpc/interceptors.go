package grpc

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

// requestIDFromIncomingMD extracts request id from gRPC metadata.
func requestIDFromIncomingMD(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	// common keys
	for _, k := range []string{"x-request-id", "x-correlation-id", "request-id"} {
		if v := md.Get(k); len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

func withRequestID(ctx context.Context) (context.Context, string) {
	rid := requestIDFromIncomingMD(ctx)
	if rid == "" {
		rid = uuid.NewString()
	}
	return context.WithValue(ctx, requestIDKey, rid), rid
}

func getRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// recoveryUnaryInterceptor converts panics into gRPC Internal errors.
func recoveryUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				appLogger.Error("grpc panic recovered",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.ByteString("stack", debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

// requestIDUnaryInterceptor ensures every request has a request id in context.
func requestIDUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, _ = withRequestID(ctx)
		return handler(ctx, req)
	}
}

// loggingUnaryInterceptor logs method, duration, status code and request-id.
func loggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		ctx, rid := withRequestID(ctx)

		resp, err := handler(ctx, req)
		st := status.Convert(err)

		fields := []zap.Field{
			zap.String("grpc.method", info.FullMethod),
			zap.String("request_id", rid),
			zap.String("grpc.code", st.Code().String()),
			zap.Duration("duration", time.Since(start)),
		}
		if err != nil {
			appLogger.Warn("grpc request", append(fields, zap.String("error", st.Message()))...)
		} else {
			appLogger.Info("grpc request", fields...)
		}
		return resp, err
	}
}

// recoveryStreamInterceptor converts panics into gRPC Internal errors for streams.
func recoveryStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				appLogger.Error("grpc stream panic recovered",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.ByteString("stack", debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(srv, ss)
	}
}

// requestIDStreamInterceptor attaches request-id to stream context.
func requestIDStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, _ := withRequestID(ss.Context())
		wrapped := &wrappedServerStream{ServerStream: ss, ctx: ctx}
		return handler(srv, wrapped)
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context { return w.ctx }

// loggingStreamInterceptor logs stream start/end.
func loggingStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		ctx, rid := withRequestID(ss.Context())
		wrapped := &wrappedServerStream{ServerStream: ss, ctx: ctx}

		err := handler(srv, wrapped)
		st := status.Convert(err)

		fields := []zap.Field{
			zap.String("grpc.method", info.FullMethod),
			zap.String("request_id", rid),
			zap.String("grpc.code", st.Code().String()),
			zap.Duration("duration", time.Since(start)),
			zap.Bool("grpc.is_client_stream", info.IsClientStream),
			zap.Bool("grpc.is_server_stream", info.IsServerStream),
		}
		if err != nil {
			appLogger.Warn("grpc stream", append(fields, zap.String("error", st.Message()))...)
		} else {
			appLogger.Info("grpc stream", fields...)
		}
		return err
	}
}

// rateLimitedMethods maps gRPC full method names that should be rate-limited
// at the server interceptor level, with their per-minute window limits.
var rateLimitedMethods = map[string]int64{
	"/insuretech.authn.services.v1.AuthService/SendOTP":      10,
	"/insuretech.authn.services.v1.AuthService/ResendOTP":    10,
	"/insuretech.authn.services.v1.AuthService/SendEmailOTP": 10,
	"/insuretech.authn.services.v1.AuthService/Login":        20,
	"/insuretech.authn.services.v1.AuthService/EmailLogin":   20,
	"/insuretech.authn.services.v1.AuthService/RefreshToken": 60,
}

// rateLimitUnaryInterceptor applies Redis-backed per-IP rate limiting for
// sensitive methods. Falls back gracefully (fail-open) when Redis is nil or
// unavailable.
func rateLimitUnaryInterceptor(rdb redis.UniversalClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		limit, ok := rateLimitedMethods[info.FullMethod]
		if !ok || rdb == nil {
			return handler(ctx, req)
		}

		// Derive rate-limit key from client IP via gRPC metadata.
		clientIP := "unknown"
		if md, ok2 := metadata.FromIncomingContext(ctx); ok2 {
			for _, hdr := range []string{"x-forwarded-for", "x-real-ip"} {
				if v := md.Get(hdr); len(v) > 0 {
					clientIP = v[0]
					break
				}
			}
		}

		rlKey := fmt.Sprintf("grpc_rl:%s:%s", info.FullMethod, clientIP)
		pipe := rdb.Pipeline()
		incrCmd := pipe.Incr(ctx, rlKey)
		pipe.Expire(ctx, rlKey, time.Minute)
		if _, err := pipe.Exec(ctx); err != nil {
			// Redis unavailable — fail open.
			appLogger.Warn("grpc rate limit: redis unavailable, failing open",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return handler(ctx, req)
		}

		count := incrCmd.Val()
		if count > limit {
			retryAfter := 60 - time.Now().Second()
			_ = grpc.SetHeader(ctx, metadata.Pairs(
				"x-ratelimit-limit", strconv.FormatInt(limit, 10),
				"x-ratelimit-remaining", "0",
				"retry-after", strconv.Itoa(retryAfter),
			))
			return nil, status.Errorf(codes.ResourceExhausted,
				"rate limit exceeded for %s: %d/%d requests per minute", info.FullMethod, count, limit)
		}

		remaining := limit - count
		_ = grpc.SetHeader(ctx, metadata.Pairs(
			"x-ratelimit-limit", strconv.FormatInt(limit, 10),
			"x-ratelimit-remaining", strconv.FormatInt(remaining, 10),
		))

		return handler(ctx, req)
	}
}

// defaultUnaryInterceptors returns the base interceptor chain (no rate limiting).
func defaultUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		recoveryUnaryInterceptor(),
		requestIDUnaryInterceptor(),
		loggingUnaryInterceptor(),
	}
}

// unaryInterceptorsWithRateLimit returns the full interceptor chain with Redis rate limiting.
func unaryInterceptorsWithRateLimit(rdb redis.UniversalClient) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		recoveryUnaryInterceptor(),
		requestIDUnaryInterceptor(),
		rateLimitUnaryInterceptor(rdb),
		loggingUnaryInterceptor(),
	}
}

func defaultStreamInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		recoveryStreamInterceptor(),
		requestIDStreamInterceptor(),
		loggingStreamInterceptor(),
	}
}

func formatFullMethod(m string) string {
	if m == "" {
		return "unknown"
	}
	return m
}

// keep helpers referenced to avoid unused-import errors.
var _ = getRequestID
var _ = formatFullMethod
