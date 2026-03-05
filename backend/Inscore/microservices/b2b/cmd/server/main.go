package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	b2bconfig "github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/config"
	b2bgrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/service"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type ServicesConfig struct {
	Services map[string]struct {
		Name  string `yaml:"name"`
		Ports struct {
			Grpc int `yaml:"grpc"`
			Http int `yaml:"http"`
		} `yaml:"ports"`
	} `yaml:"services"`
}

func main() {
	if err := appLogger.Initialize(appLogger.NoFileConfig()); err != nil {
		appLogger.Fatalf("failed to initialize logger: %v", err)
	}
	_ = env.Load()
	appLogger.Info("Starting B2B microservice...")

	servicesConfigPath, err := config.ResolveConfigPath("services.yaml")
	if err != nil {
		appLogger.Fatalf("failed to resolve services.yaml path: %v", err)
	}
	servicesData, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		appLogger.Fatalf("failed to read services.yaml: %v", err)
	}

	var svcConfig ServicesConfig
	if err := yaml.Unmarshal(servicesData, &svcConfig); err != nil {
		appLogger.Fatalf("failed to parse services.yaml: %v", err)
	}

	b2bSvc, exists := svcConfig.Services["b2b"]
	if !exists {
		appLogger.Fatal("configuration for 'b2b' service not found in services.yaml")
	}
	if os.Getenv("B2B_PORT") != "" || os.Getenv("B2B_GRPC_PORT") != "" || os.Getenv("B2B_HTTP_PORT") != "" {
		appLogger.Warn("B2B_PORT/B2B_GRPC_PORT/B2B_HTTP_PORT env values are ignored; using backend/inscore/configs/services.yaml")
	}
	grpcPort := strconv.Itoa(b2bSvc.Ports.Grpc)
	appLogger.Info("service configured from services.yaml",
		zap.String("service", b2bSvc.Name),
		zap.Int("grpc_port", b2bSvc.Ports.Grpc),
		zap.Int("http_port", b2bSvc.Ports.Http),
	)

	if _, err := b2bconfig.Load(); err != nil {
		appLogger.Fatalf("failed to load b2b config: %v", err)
	}

	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("failed to resolve database config path: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Fatal("database initialization failed", zap.Error(err))
	}
	defer db.Manager.Close()
	database := db.GetDB()

	portalRepo := repository.NewPortalRepository(database)
	b2bService := service.NewB2BService(portalRepo)

	serverConfig := b2bgrpc.DefaultServerConfig()
	serverConfig.Port = grpcPort
	serverConfig.DB = database

	server, err := b2bgrpc.NewServer(serverConfig, b2bService)
	if err != nil {
		appLogger.Fatalf("failed to create gRPC server: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.HealthCheck(ctx); err != nil {
		appLogger.Fatal("Server health check failed", zap.Error(err))
	}

	go func() {
		if err := server.Start(); err != nil {
			appLogger.Fatalf("server crashed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down b2b service...")
	server.Stop()
	appLogger.Info("B2B service stopped")
}
