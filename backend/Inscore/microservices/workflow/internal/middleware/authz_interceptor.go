package middleware

import (
	"context"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"google.golang.org/grpc"
)

// MetadataExtractInterceptor is retained as a no-op for backwards compatibility.
// Workflow now reads metadata directly via pkg/grpcmeta.
func MetadataExtractInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		return handler(ctx, req)
	}
}

// LoggingInterceptor logs all incoming gRPC calls for observability.
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		userID := grpcmeta.ActorID(ctx, "")
		method := info.FullMethod
		if idx := strings.LastIndex(method, "/"); idx >= 0 {
			method = method[idx+1:]
		}
		appLogger.Infof("workflow-engine: %s user=%s", method, userID)
		resp, err := handler(ctx, req)
		if err != nil {
			// AlreadyExists is expected during idempotent seeding — log at DEBUG not WARN
			if strings.Contains(err.Error(), "AlreadyExists") || strings.Contains(err.Error(), "already exists") {
				appLogger.Debug("workflow-engine: " + method + " idempotent (already exists)")
			} else {
				appLogger.Warnf("workflow-engine: %s error: %v", method, err)
			}
		}
		return resp, err
	}
}
