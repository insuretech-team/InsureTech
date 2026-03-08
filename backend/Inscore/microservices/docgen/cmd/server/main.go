package main

import (
	"context"
	"database/sql"
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
)

func main() {
	if err := logger.Initialize(logger.NoFileConfig()); err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer logger.GetLogger().Sync() //nolint:errcheck

	logger.Info("Starting DocGen gRPC server...")
	_ = env.Load()

	sqlDB, cleanupDB, err := setupDatabase()
	if err != nil {
		logger.Fatal("failed to setup database", zap.Error(err))
	}
	defer cleanupDB()

	var storageConn *grpc.ClientConn
	if storageAddr := strings.TrimSpace(getFirstEnv("STORAGE_GRPC_ADDR", "STORAGE_SERVICE_ADDR")); storageAddr != "" {
		dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		storageConn, err = grpc.DialContext(
			dialCtx,
			storageAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()
		if err != nil {
			logger.Warn("failed to connect storage service; running in inline-only mode", zap.String("addr", storageAddr), zap.Error(err))
			storageConn = nil
		} else {
			defer func() { _ = storageConn.Close() }()
			logger.Info("Storage service connected", zap.String("addr", storageAddr))
		}
	}

	docgenServer, err := docgenpkg.NewDocumentServer(sqlDB, storageConn)
	if err != nil {
		logger.Fatal("failed to create docgen server", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	documentservicev1.RegisterDocumentServiceServer(grpcServer, docgenServer.Handler())

	grpcPort := getEnvOrDefault("DOCGEN_GRPC_PORT", "50280")
	if grpcPort != "" && grpcPort[0] != ':' {
		grpcPort = ":" + grpcPort
	}

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err), zap.String("grpc_port", grpcPort))
	}

	go func() {
		logger.Info("DocGen gRPC server listening", zap.String("grpc_port", grpcPort))
		if serveErr := grpcServer.Serve(listener); serveErr != nil {
			logger.Error("gRPC server error", zap.Error(serveErr))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigChan
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	if err := docgenServer.Close(); err != nil {
		logger.Warn("DocGen server resource shutdown reported errors", zap.Error(err))
	}
	grpcServer.GracefulStop()
	logger.Info("DocGen server stopped")
}

// setupDatabase initializes the database manager and returns the underlying sql.DB plus cleanup callback.
func setupDatabase() (*sql.DB, func(), error) {
	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		return nil, nil, fmt.Errorf("resolve database config path: %w", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		return nil, nil, fmt.Errorf("initialize database manager: %w", err)
	}

	gormDB := db.Manager.GetDB()
	sqlDB, err := gormDB.DB()
	if err != nil {
		db.Manager.Close()
		return nil, nil, fmt.Errorf("get sql.DB from gorm: %w", err)
	}

	cleanup := func() {
		db.Manager.Close()
	}
	return sqlDB, cleanup, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}

func getFirstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}
