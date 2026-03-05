package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// Initialize logger
	if err := logger.Initialize(logger.NoFileConfig()); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.GetLogger().Sync() //nolint:errcheck

	logger.Info("Starting Tenant Service...")
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

	logger.Info("Tenant service database ready",
		zap.String("active_db", string(db.Manager.GetCurrentType())),
		zap.Bool("failover_enabled", db.Manager.GetPrimaryDB() != nil && db.Manager.GetBackupDB() != nil))

	// TODO: Initialize tenant server when implemented
	// tenantServer, err := tenantpkg.NewTenantServer(sqlDB)
	// if err != nil {
	// 	logger.Fatal("Failed to create tenant server", zap.Error(err))
	// }
	_ = sqlDB

	// Setup gRPC server
	grpcPort := getEnvOrDefault("TENANT_GRPC_PORT", "50050")
	if grpcPort != "" && grpcPort[0] != ':' {
		grpcPort = ":" + grpcPort
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	// TODO: Register tenant service when implemented
	// tenantservicev1.RegisterTenantServiceServer(grpcServer, tenantServer)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Start gRPC server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("Tenant gRPC server error", zap.Error(err))
		}
	}()

	logger.Info("Tenant service running",
		zap.String("grpc_port", grpcPort))

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down tenant service...")
	grpcServer.GracefulStop()
	logger.Info("Tenant service stopped")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
