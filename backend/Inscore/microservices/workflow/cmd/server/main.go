package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/events"
	workflowgrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/service"
	kafkaproducer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
	configpkg "github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

// ServicesConfig mirrors the structure of backend/inscore/configs/services.yaml
type ServicesConfig = serviceaddr.ServicesConfig

func main() {
	// ── 1. Logger ─────────────────────────────────────────────────────────────
	if err := appLogger.Initialize(appLogger.NoFileConfig()); err != nil {
		appLogger.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger.Info("Starting workflow-engine microservice...")
	if err := env.Load(); err != nil {
		appLogger.Warn("No .env file found, using system environment variables", zap.Error(err))
	}

	// ── 2. Config ─────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatalf("Failed to load config: %v", err)
	}

	servicesConfigPath, err := configpkg.ResolveConfigPath("services.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve services.yaml path: %v", err)
	}
	servicesData, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		appLogger.Fatalf("Failed to read services.yaml: %v", err)
	}
	var svcConfig ServicesConfig
	if err := yaml.Unmarshal(servicesData, &svcConfig); err != nil {
		appLogger.Fatalf("Failed to parse services.yaml: %v", err)
	}

	workflowServiceConfig, exists := svcConfig.Services["workflow"]
	if !exists {
		appLogger.Fatal("configuration for 'workflow' service not found in services.yaml")
	}
	cfg.GRPCPort = workflowServiceConfig.Ports.Grpc
	appLogger.Info("service configured from services.yaml",
		zap.String("service", workflowServiceConfig.Name),
		zap.Int("grpc_port", workflowServiceConfig.Ports.Grpc),
		zap.Int("http_port", workflowServiceConfig.Ports.Http),
	)

	// ── 3. Database ───────────────────────────────────────────────────────────
	dbConfigPath, err := configpkg.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve database config path: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Manager.Close()

	database := db.GetDB()

	// ── 4. Repository ─────────────────────────────────────────────────────────
	repo := repository.New(database)

	// ── 5. Kafka producer (optional — service starts without Kafka) ───────────
	var producer *kafkaproducer.EventProducer
	if len(cfg.KafkaBrokers) > 0 {
		kp, err := kafkaproducer.NewEventProducer(
			cfg.KafkaBrokers,
			events.TopicWorkflowStarted,
			"workflow-engine-producer",
		)
		if err != nil {
			appLogger.Warn("Kafka producer unavailable — events will be skipped", zap.Error(err))
		} else {
			producer = kp
			appLogger.Info("Kafka producer connected")
		}
	}

	// ── 6. Publisher + Service ────────────────────────────────────────────────
	publisher := events.NewPublisher(producer)
	svc := service.New(repo, publisher)

	// ── 7. gRPC server ────────────────────────────────────────────────────────
	serverCfg := &workflowgrpc.Config{
		Port: fmt.Sprintf("%d", cfg.GRPCPort),
		DB:   database,
	}

	grpcServer, err := workflowgrpc.NewServer(
		serverCfg, svc,
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(),
		),
	)
	if err != nil {
		appLogger.Fatalf("failed to create gRPC server: %v", err)
	}

	// ── 8. Graceful shutdown ──────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Start(); err != nil {
			appLogger.Fatalf("gRPC server error: %v", err)
		}
	}()

	appLogger.Info("workflow-engine ready", zap.Int("grpc_port", cfg.GRPCPort))

	<-quit
	appLogger.Info("Shutting down workflow-engine...")

	grpcServer.Stop()

	if producer != nil {
		producer.Close()
	}

	// Allow in-flight requests to drain
	time.Sleep(500 * time.Millisecond)
	appLogger.Info("workflow-engine stopped gracefully")
}
