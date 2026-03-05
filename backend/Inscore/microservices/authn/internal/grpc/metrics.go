package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// MetricsUnaryInterceptor records per-RPC latency, call count, and error codes
// using structured logging. Swap out for Prometheus counters/histograms when
// a metrics registry is wired.
func MetricsUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}

		logger := appLogger.GetLogger()
		logger.Info("grpc.rpc",
			zap.String("method", info.FullMethod),
			zap.String("code", code.String()),
			zap.Duration("duration_ms", duration),
			zap.Bool("error", err != nil),
		)

		return resp, err
	}
}

// MetricsStreamInterceptor records per-stream metrics.
func MetricsStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		err := handler(srv, ss)
		duration := time.Since(start)

		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}

		logger := appLogger.GetLogger()
		logger.Info("grpc.stream",
			zap.String("method", info.FullMethod),
			zap.String("code", code.String()),
			zap.Duration("duration_ms", duration),
			zap.Bool("error", err != nil),
		)

		return err
	}
}
