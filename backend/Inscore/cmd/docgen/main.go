package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	docgenpkg "github.com/newage-saint/insuretech/backend/inscore/microservices/docgen"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	documentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if err := logger.Initialize(logger.NoFileConfig()); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.GetLogger().Sync() //nolint:errcheck

	logger.Info("Starting DocGen Service...")
	_ = env.Load()

	// Use ops/config for proper path resolution from project root
	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		logger.Fatal("Failed to resolve database config path", zap.Error(err))
	}

	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		logger.Fatal("Failed to initialize database manager", zap.Error(err))
	}
	defer db.Manager.Close()

	// Get database connection
	gormDB := db.Manager.GetDB()
	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Fatal("Failed to get sql.DB from gorm", zap.Error(err))
	}

	logger.Info("DocGen service database ready",
		zap.String("active_db", string(db.Manager.GetCurrentType())),
		zap.Bool("failover_enabled", db.Manager.GetPrimaryDB() != nil && db.Manager.GetBackupDB() != nil))

	var storageConn *grpc.ClientConn
	if storageAddr := strings.TrimSpace(os.Getenv("STORAGE_GRPC_ADDR")); storageAddr != "" {
		dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		storageConn, err = grpc.DialContext(dialCtx, storageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		cancel()
		if err != nil {
			logger.Warn("Failed to connect storage service; docgen will use inline-only mode", zap.String("addr", storageAddr), zap.Error(err))
			storageConn = nil
		} else {
			defer storageConn.Close()
			logger.Info("Storage service connected", zap.String("addr", storageAddr))
		}
	}

	grpcPort := getEnvOrDefault("DOCGEN_GRPC_PORT", "50280")
	if grpcPort != "" && grpcPort[0] != ':' {
		grpcPort = ":" + grpcPort
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	docServer, err := docgenpkg.NewDocumentServer(sqlDB, storageConn)
	if err != nil {
		logger.Fatal("Failed to create docgen server", zap.Error(err))
	}
	documentservicev1.RegisterDocumentServiceServer(grpcServer, docServer.Handler())

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	go func() {
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			logger.Error("DocGen gRPC server error", zap.Error(serveErr))
		}
	}()

	logger.Info("DocGen service running", zap.String("grpc_port", grpcPort))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down docgen service...")
	grpcServer.GracefulStop()
	logger.Info("DocGen service stopped")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
