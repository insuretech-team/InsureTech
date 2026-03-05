package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	partnerconfig "github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/events"
	partnergrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/service"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

// ServicesConfig structure matches services.yaml.
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
	appLogger.Info("Starting Partner microservice...")

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

	partnerSvc, exists := svcConfig.Services["partner"]
	if !exists {
		appLogger.Fatal("configuration for 'partner' service not found in services.yaml")
	}
	if os.Getenv("PARTNER_PORT") != "" || os.Getenv("PARTNER_GRPC_PORT") != "" || os.Getenv("PARTNER_HTTP_PORT") != "" {
		appLogger.Warn("PARTNER_PORT/PARTNER_GRPC_PORT/PARTNER_HTTP_PORT env values are ignored; using backend/inscore/configs/services.yaml")
	}
	grpcPort := strconv.Itoa(partnerSvc.Ports.Grpc)
	appLogger.Info("service configured from services.yaml",
		zap.String("service", partnerSvc.Name),
		zap.Int("grpc_port", partnerSvc.Ports.Grpc),
		zap.Int("http_port", partnerSvc.Ports.Http),
	)

	cfg, err := partnerconfig.Load()
	if err != nil {
		appLogger.Fatalf("failed to load partner config: %v", err)
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

	partnerRepo := repository.NewPartnerRepository(database, cfg.Security.PIIEncryptionKey)
	agentRepo := repository.NewAgentRepository(database, cfg.Security.PIIEncryptionKey)
	commissionRepo := repository.NewCommissionRepository(database)

	kafkaProducer, err := producer.NewEventProducerWithRetry(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		"partner-service",
		5,
		2*time.Second,
	)
	var eventPublisher *events.Publisher
	if err != nil {
		appLogger.Warn("Kafka producer init failed; partner events will be dropped", zap.Error(err))
		eventPublisher = events.NewPublisher(nil, cfg.Kafka.Topic)
	} else {
		defer kafkaProducer.Close()
		eventPublisher = events.NewPublisher(kafkaProducer, cfg.Kafka.Topic)
	}

	authnAddr := resolveServiceAddr(cfg.Integration.AuthNAddress, svcConfig.Services, "authn")
	if authnAddr == "" {
		appLogger.Fatal("authn address is empty; configure AUTHN_SERVICE_ADDRESS or services.yaml entry")
	}
	appLogger.Info("Connecting to AuthN service", zap.String("addr", authnAddr))
	authnConn, err := grpc.Dial(authnAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		appLogger.Fatal("Failed to connect to AuthN service", zap.Error(err), zap.String("addr", authnAddr))
	}
	defer authnConn.Close()
	authnClient := authnservicev1.NewAuthServiceClient(authnConn)

	partnerService := service.NewPartnerService(partnerRepo, agentRepo, commissionRepo, eventPublisher, authnClient)

	consumerGroup, consumerErr := kafkaconsumer.NewConsumerGroup(kafkaconsumer.Config{
		Brokers:  cfg.Kafka.Brokers,
		GroupID:  cfg.Kafka.ConsumerGroup,
		Topics:   cfg.Kafka.ConsumerTopics,
		Handler:  events.NewPolicyLifecycleHandler(partnerService),
		DLQTopic: cfg.Kafka.DLQTopic,
		ClientID: "partner-consumer",
	})
	if consumerErr != nil {
		appLogger.Warn("Kafka consumer group failed to start; policy commission events will not be consumed", zap.Error(consumerErr))
	} else {
		consumerCtx, consumerCancel := context.WithCancel(context.Background())
		defer consumerCancel()
		go consumerGroup.Start(consumerCtx)
		defer func() {
			consumerCancel()
			_ = consumerGroup.Close()
		}()
		appLogger.Info("Kafka consumer group started", zap.Strings("topics", cfg.Kafka.ConsumerTopics))
	}

	authzAddr := resolveServiceAddr(cfg.Integration.AuthZAddress, svcConfig.Services, "authz")
	if authzAddr == "" {
		appLogger.Fatal("authz address is empty; configure AUTHZ_SERVICE_ADDRESS or services.yaml entry")
	}
	appLogger.Info("Connecting to AuthZ service", zap.String("addr", authzAddr))
	authzConn, err := grpc.Dial(authzAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		appLogger.Fatal("Failed to connect to AuthZ service", zap.Error(err), zap.String("addr", authzAddr))
	}
	defer authzConn.Close()
	authzClient := authzservicev1.NewAuthZServiceClient(authzConn)

	serverConfig := partnergrpc.DefaultServerConfig()
	serverConfig.Port = grpcPort
	serverConfig.DB = database

	server, err := partnergrpc.NewServer(serverConfig, partnerService, authzClient)
	if err != nil {
		appLogger.Fatalf("Failed to create gRPC server: %v", err)
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

	appLogger.Info("Shutting down partner service...")
	server.Stop()
	appLogger.Info("Partner service stopped")
}

func resolveServiceAddr(explicit string, services map[string]struct {
	Name  string `yaml:"name"`
	Ports struct {
		Grpc int `yaml:"grpc"`
		Http int `yaml:"http"`
	} `yaml:"ports"`
}, key string) string {
	v := strings.TrimSpace(explicit)
	if v != "" {
		return v
	}
	svc, ok := services[key]
	if !ok || svc.Ports.Grpc <= 0 {
		return ""
	}
	return "localhost:" + strconv.Itoa(svc.Ports.Grpc)
}
