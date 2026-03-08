package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/consumers"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/events"
	b2bgrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/service"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	kafkaproducer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	configpkg "github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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
	// 1. Initialize Logger
	if err := appLogger.Initialize(appLogger.NoFileConfig()); err != nil {
		appLogger.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger.Info("Starting B2B microservice (VSA)...")
	if err := env.Load(); err != nil {
		appLogger.Warn("No .env file found, using system environment variables", zap.Error(err))
	}

	// Load configuration
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

	b2bServiceConfig, exists := svcConfig.Services["b2b"]
	if !exists {
		appLogger.Fatal("configuration for 'b2b' service not found in services.yaml")
	}
	if os.Getenv("B2B_PORT") != "" || os.Getenv("B2B_GRPC_PORT") != "" || os.Getenv("B2B_HTTP_PORT") != "" || os.Getenv("GRPC_PORT") != "" {
		appLogger.Warn("B2B_PORT/B2B_GRPC_PORT/B2B_HTTP_PORT/GRPC_PORT env values are ignored; using backend/inscore/configs/services.yaml")
	}
	cfg.GRPCPort = b2bServiceConfig.Ports.Grpc
	appLogger.Info("service configured from services.yaml",
		zap.String("service", b2bServiceConfig.Name),
		zap.Int("grpc_port", b2bServiceConfig.Ports.Grpc),
		zap.Int("http_port", b2bServiceConfig.Ports.Http),
	)

	authzAddr := resolveServiceAddr(cfg.AuthZServiceURL, svcConfig.Services, "authz")
	if authzAddr == "" {
		appLogger.Fatal("authz address is empty; configure AUTHZ_SERVICE_URL or services.yaml entry")
	}
	cfg.AuthZServiceURL = authzAddr
	appLogger.Info("Connecting to AuthZ service", zap.String("addr", cfg.AuthZServiceURL))

	// Initialize database using shared database manager
	dbConfigPath, err := configpkg.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve database config path: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Manager.Close()

	database := db.GetDB()

	// Initialize Kafka producer with retry
	kafkaProducer, err := kafkaproducer.NewEventProducerWithRetry(
		cfg.KafkaBrokers,
		"", // topic will be specified per message
		"b2b-service",
		5,             // max retries
		2*time.Second, // retry delay
	)
	if err != nil {
		appLogger.Warnf("Failed to initialize Kafka producer: %v (continuing without events)", err)
		kafkaProducer = nil
	}

	// Initialize event publisher.
	// Always create a Publisher — even when Kafka is unavailable.
	// events.NewPublisher accepts a nil producer and the internal publish()
	// method already no-ops gracefully when producer == nil.
	// Passing a typed-nil *events.Publisher to NewB2BService causes a panic
	// because the interface is non-nil but the receiver is nil.
	publisher := events.NewPublisher(kafkaProducer)

	// Initialize repository
	repo := repository.NewPortalRepository(database)

	// Initialize service
	svc := service.NewB2BService(repo, publisher)

	// Connect to AuthZ service
	authzConn, err := grpc.Dial(cfg.AuthZServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		appLogger.Fatalf("Failed to connect to AuthZ service: %v", err)
	}
	defer authzConn.Close()

	authzClient := authzservicev1.NewAuthZServiceClient(authzConn)

	// Create wrapper for middleware
	authzClientWrapper := &authzClientAdapter{client: authzClient}

	// Initialize authorization interceptor
	authzInterceptor := middleware.NewAuthZInterceptor(authzClientWrapper)

	// Initialize event consumer
	if kafkaProducer != nil {
		consumer := consumers.NewEventConsumer(authzClientWrapper)
		go startEventConsumer(cfg, consumer)
	}

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authzInterceptor.UnaryServerInterceptor()),
	)

	// Register B2B service
	handler := b2bgrpc.NewB2BHandler(svc)
	b2bservicev1.RegisterB2BServiceServer(grpcServer, handler)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		appLogger.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		appLogger.Infof("B2B microservice listening on 0.0.0.0:%d", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			appLogger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down B2B microservice...")
	grpcServer.GracefulStop()
	if kafkaProducer != nil {
		kafkaProducer.Close()
	}
}

func startEventConsumer(cfg *config.Config, consumer *consumers.EventConsumer) {
	// Subscribe to B2B events
	b2bTopics := []string{
		events.TopicOrganisationCreated,
		events.TopicB2BAdminAssigned,
		events.TopicOrgMemberAdded,
		events.TopicOrganisationApproved,
	}

	// Subscribe to AuthN events (to react to user registrations)
	authnTopics := []string{
		"authn.user.registered",
	}

	// Subscribe to AuthZ events (to react to role assignments)
	authzTopics := []string{
		"authz.role.assigned",
	}

	// Combine all topics
	allTopics := append(b2bTopics, authnTopics...)
	allTopics = append(allTopics, authzTopics...)

	// Create Kafka consumer configuration
	consumerCfg := kafkaconsumer.Config{
		Brokers:        cfg.KafkaBrokers,
		GroupID:        "b2b-service-group",
		Topics:         allTopics,
		ClientID:       "b2b-service-consumer",
		InitialOffset:  -1, // Start from latest
		SessionTimeout: 10 * time.Second,
		Handler: func(ctx context.Context, msg *kafkaconsumer.Message) error {
			// Route message to appropriate handler based on topic
			switch msg.Topic {
			case events.TopicOrganisationCreated:
				return consumer.HandleOrganisationCreated(ctx, msg.Value)
			case events.TopicB2BAdminAssigned:
				return consumer.HandleB2BAdminAssigned(ctx, msg.Value)
			case events.TopicOrgMemberAdded:
				return consumer.HandleOrgMemberAdded(ctx, msg.Value)
			case events.TopicOrganisationApproved:
				return consumer.HandleOrganisationApproved(ctx, msg.Value)
			case "authn.user.registered":
				return consumer.HandleUserRegistered(ctx, msg.Value)
			case "authz.role.assigned":
				return consumer.HandleRoleAssigned(ctx, msg.Value)
			default:
				appLogger.Warnf("Unknown topic: %s", msg.Topic)
				return nil
			}
		},
		DLQTopic: "b2b-service-dlq", // Dead letter queue for failed messages
	}

	kafkaConsumer, err := kafkaconsumer.NewConsumerGroup(consumerCfg)
	if err != nil {
		appLogger.Errorf("Failed to create Kafka consumer: %v", err)
		return
	}

	appLogger.Infof("B2B event consumer started, listening to topics: %v", allTopics)

	// Start consuming in blocking mode
	kafkaConsumer.Start(context.Background())
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
	return key + ":" + strconv.Itoa(svc.Ports.Grpc)
}

// authzClientAdapter adapts the gRPC client to the interface expected by middleware and consumers
type authzClientAdapter struct {
	client authzservicev1.AuthZServiceClient
}

func (a *authzClientAdapter) CheckAccess(ctx context.Context, req *authzservicev1.CheckAccessRequest) (*authzservicev1.CheckAccessResponse, error) {
	return a.client.CheckAccess(withInternalServiceContext(ctx, "b2b-service"), req)
}

func (a *authzClientAdapter) AssignRole(ctx context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error) {
	return a.client.AssignRole(withInternalServiceContext(ctx, "b2b-service"), req)
}

func (a *authzClientAdapter) CreatePolicyRule(ctx context.Context, req *authzservicev1.CreatePolicyRuleRequest) (*authzservicev1.CreatePolicyRuleResponse, error) {
	return a.client.CreatePolicyRule(withInternalServiceContext(ctx, "b2b-service"), req)
}

func (a *authzClientAdapter) ListRoles(ctx context.Context, req *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error) {
	return a.client.ListRoles(withInternalServiceContext(ctx, "b2b-service"), req)
}

func withInternalServiceContext(ctx context.Context, serviceName string) context.Context {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		cloned := md.Copy()
		cloned.Set("x-internal-service", serviceName)
		return metadata.NewOutgoingContext(ctx, cloned)
	}

	return metadata.NewOutgoingContext(ctx, metadata.Pairs("x-internal-service", serviceName))
}
