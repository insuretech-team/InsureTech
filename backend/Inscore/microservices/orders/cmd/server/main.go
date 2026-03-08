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
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/consumers"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/events"
	ordersgrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/service"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	kafkaproducer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	documentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
	orderservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/services/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	configpkg "github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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
	// ── 1. Logger ────────────────────────────────────────────────────────────
	if err := appLogger.Initialize(appLogger.NoFileConfig()); err != nil {
		appLogger.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger.Info("Starting Orders microservice...")
	if err := env.Load(); err != nil {
		appLogger.Warn("No .env file found, using system environment variables", zap.Error(err))
	}

	// ── 2. Config ────────────────────────────────────────────────────────────
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

	ordersServiceConfig, exists := svcConfig.Services["orders"]
	if !exists {
		appLogger.Fatal("configuration for 'orders' service not found in services.yaml")
	}
	cfg.GRPCPort = ordersServiceConfig.Ports.Grpc
	appLogger.Info("service configured from services.yaml",
		zap.String("service", ordersServiceConfig.Name),
		zap.Int("grpc_port", ordersServiceConfig.Ports.Grpc),
		zap.Int("http_port", ordersServiceConfig.Ports.Http),
	)

	// Resolve downstream payment service address
	paymentAddr := resolveServiceAddr(cfg.PaymentServiceURL, svcConfig.Services, "payment")
	if paymentAddr != "" {
		cfg.PaymentServiceURL = paymentAddr
		appLogger.Info("Payment service address resolved", zap.String("addr", paymentAddr))
	} else {
		appLogger.Warn("Payment service address not resolved — payment initiation will fail until payment-service is reachable")
	}

	// ── 3. Database ──────────────────────────────────────────────────────────
	dbConfigPath, err := configpkg.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve database config path: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Manager.Close()

	database := db.GetDB()

	// ── 4. Repository ────────────────────────────────────────────────────────
	repo := repository.NewOrderRepository(database)

	// ── 5. Kafka producer (optional — service starts without Kafka) ──────────
	var kafkaProducer *kafkaproducer.EventProducer
	kafkaBrokers := cfg.KafkaBrokers
	if len(kafkaBrokers) > 0 {
		kp, err := kafkaproducer.NewEventProducer(kafkaBrokers, events.TopicOrderCreated, "orders-service-producer")
		if err != nil {
			appLogger.Warnf("Failed to create Kafka producer (continuing without events): %v", err)
		} else {
			kafkaProducer = kp
			appLogger.Info("Kafka producer connected", zap.Strings("brokers", kafkaBrokers))
		}
	}

	publisher := events.NewPublisher(kafkaProducer)

	// ── 6. AuthZ gRPC client (optional — interceptor fails open if unavailable) ──
	var authzClient authzservicev1.AuthZServiceClient
	authzAddr := resolveServiceAddr(cfg.AuthzServiceURL, svcConfig.Services, "authz")
	if authzAddr == "" {
		authzAddr = cfg.AuthzServiceURL
	}
	authzConn, err := grpc.NewClient(authzAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		appLogger.Warnf("Failed to dial authz-service at %s — AuthZ interceptor will be disabled: %v", authzAddr, err)
	} else {
		authzClient = authzservicev1.NewAuthZServiceClient(authzConn)
		appLogger.Info("AuthZ service dialed", zap.String("addr", authzAddr))
		defer authzConn.Close()
	}

	// ── 7. Payment gRPC client (optional — InitiatePayment fails if unavailable) ─
	var paymentClient paymentservicev1.PaymentServiceClient
	if cfg.PaymentServiceURL != "" {
		payConn, err := grpc.NewClient(cfg.PaymentServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			appLogger.Warnf("Failed to dial payment-service at %s — InitiatePayment will fail: %v", cfg.PaymentServiceURL, err)
		} else {
			paymentClient = paymentservicev1.NewPaymentServiceClient(payConn)
			appLogger.Info("Payment service dialed", zap.String("addr", cfg.PaymentServiceURL))
			defer payConn.Close()
		}
	}

	// ── 8. Service ───────────────────────────────────────────────────────────
	svc := service.NewOrderService(repo, publisher, paymentClient)

	// ── 9. Kafka consumer (optional) ─────────────────────────────────────────
	if kafkaProducer != nil {
		// ── Docgen gRPC client (receipt PDF generation after payment) ───────
		docgenAddr := envOrDefault("DOCGEN_GRPC_ADDR", "docgen-service:50170")
		docgenConn, _ := grpc.NewClient(docgenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		var docgenClient documentservicev1.DocumentServiceClient
		if docgenConn != nil {
			docgenClient = documentservicev1.NewDocumentServiceClient(docgenConn)
		}

		// ── Storage gRPC client (manual-proof file validation) ────────────────
		storageAddr := envOrDefault("STORAGE_GRPC_ADDR", "storage-service:50175")
		storageConn, _ := grpc.NewClient(storageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		var storageClient storageservicev1.StorageServiceClient
		if storageConn != nil {
			storageClient = storageservicev1.NewStorageServiceClient(storageConn)
		}

		consumer := consumers.NewEventConsumer(repo, svc, docgenClient, storageClient)
		go startEventConsumer(cfg, consumer)
	}

	// ── 10. gRPC server with AuthZ interceptor ────────────────────────────────
	var serverOpts []grpc.ServerOption
	if authzClient != nil {
		authzInterceptor := middleware.NewOrderAuthZInterceptor(authzClient)
		serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(
			authzInterceptor.UnaryServerInterceptor(),
		))
		appLogger.Info("AuthZ interceptor wired into gRPC server")
	} else {
		appLogger.Warn("AuthZ interceptor NOT active — running without authorization enforcement")
	}

	grpcServer := grpc.NewServer(serverOpts...)

	handler := ordersgrpc.NewOrderHandler(svc)
	orderservicev1.RegisterOrderServiceServer(grpcServer, handler)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("orders", grpc_health_v1.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		appLogger.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	go func() {
		appLogger.Infof("Orders gRPC server listening on %s", addr)
		if err := grpcServer.Serve(lis); err != nil {
			appLogger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// ── 9. Graceful shutdown ──────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down Orders microservice...")
	healthSrv.SetServingStatus("orders", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
		appLogger.Info("Orders service stopped cleanly")
	case <-time.After(15 * time.Second):
		appLogger.Warn("Graceful stop timed out — forcing")
		grpcServer.Stop()
	}

	if kafkaProducer != nil {
		kafkaProducer.Close()
	}
}

// startEventConsumer starts the Kafka consumer for payment, policy, and B2B events.
func startEventConsumer(cfg *config.Config, consumer *consumers.EventConsumer) {
	topics := []string{
		events.TopicPaymentCompleted,       // payment fully confirmed
		events.TopicPaymentFailed,          // payment gateway failure
		events.TopicPaymentVerified,        // manual proof approved → re-trigger
		events.TopicManualReviewRequested,  // manual proof submitted → hold
		events.TopicManualPaymentReviewed,  // manual review decision (approve/reject)
		events.TopicPolicyIssued,           // policy issued → fulfillment complete
		events.TopicB2BPurchaseOrderApproved, // B2B PO approved → create order
	}

	consumerCfg := kafkaconsumer.Config{
		Brokers:        cfg.KafkaBrokers,
		GroupID:        "orders-service-group",
		Topics:         topics,
		ClientID:       "orders-service-consumer",
		InitialOffset:  -1,
		SessionTimeout: 10 * time.Second,
		Handler: func(ctx context.Context, msg *kafkaconsumer.Message) error {
			switch msg.Topic {
			case events.TopicPaymentCompleted:
				return consumer.HandlePaymentCompleted(ctx, msg.Value)
			case events.TopicPaymentFailed:
				return consumer.HandlePaymentFailed(ctx, msg.Value)
			case events.TopicPaymentVerified:
				return consumer.HandlePaymentVerified(ctx, msg.Value)
			case events.TopicManualReviewRequested:
				return consumer.HandleManualReviewRequested(ctx, msg.Value)
			case events.TopicManualPaymentReviewed:
				return consumer.HandleManualPaymentReviewed(ctx, msg.Value)
			case events.TopicPolicyIssued:
				return consumer.HandlePolicyIssued(ctx, msg.Value)
			case events.TopicB2BPurchaseOrderApproved:
				return consumer.HandleB2BPurchaseOrderApproved(ctx, msg.Value)
			default:
				appLogger.Warnf("orders consumer: unknown topic %s", msg.Topic)
				return nil
			}
		},
		DLQTopic: "orders-service-dlq",
	}

	kafkaConsumer, err := kafkaconsumer.NewConsumerGroup(consumerCfg)
	if err != nil {
		appLogger.Errorf("Failed to create Kafka consumer: %v", err)
		return
	}

	appLogger.Infof("Orders event consumer started, listening to topics: %v", topics)
	kafkaConsumer.Start(context.Background())
}

// envOrDefault returns the environment variable value or falls back to defaultVal.
func envOrDefault(key, defaultVal string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return defaultVal
}

// resolveServiceAddr resolves a service gRPC address from services.yaml,
// falling back to the explicit override if provided.
func resolveServiceAddr(explicit string, services map[string]struct {
	Name  string `yaml:"name"`
	Ports struct {
		Grpc int `yaml:"grpc"`
		Http int `yaml:"http"`
	} `yaml:"ports"`
}, key string) string {
	if v := strings.TrimSpace(explicit); v != "" {
		return v
	}
	svc, ok := services[key]
	if !ok || svc.Ports.Grpc <= 0 {
		return ""
	}
	return key + ":" + strconv.Itoa(svc.Ports.Grpc)
}
