package payment

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	paymentcfg "github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/config"
	paymentevents "github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/events"
	paymentgrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/grpc"
	paymentmw "github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/providers/sslcommerz"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/repository"
	paymentservice "github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/service"
	kafkaproducer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Run() {
	if err := logger.Initialize(logger.NoFileConfig()); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.GetLogger().Sync() //nolint:errcheck

	logger.Info("Starting Payment Service...")
	_ = env.Load()

	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		logger.Fatal("Failed to resolve database config path", zap.Error(err))
	}

	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		logger.Fatal("Failed to initialize database manager", zap.Error(err))
	}
	defer db.Manager.Close()

	gormDB := db.Manager.GetDB()
	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Fatal("Failed to get sql.DB from gorm", zap.Error(err))
	}

	logger.Info("Payment service database ready",
		zap.String("active_db", string(db.Manager.GetCurrentType())),
		zap.Bool("failover_enabled", db.Manager.GetPrimaryDB() != nil && db.Manager.GetBackupDB() != nil))
	_ = sqlDB

	cfg, err := paymentcfg.Load()
	if err != nil {
		logger.Fatal("Failed to load payment config", zap.Error(err))
	}

	repo := repository.NewPaymentRepository(gormDB)

	var kafkaProducer *kafkaproducer.EventProducer
	if len(cfg.KafkaBrokers) > 0 {
		kafkaProducer, err = kafkaproducer.NewEventProducer(cfg.KafkaBrokers, paymentevents.TopicPaymentInitiated, "payment-service-producer")
		if err != nil {
			logger.Warn("Failed to create Kafka producer, continuing without events", zap.Error(err))
		}
	}

	publisher := paymentevents.NewPublisher(kafkaProducer)
	gatewayClient := sslcommerz.NewClient(cfg)
	svc := paymentservice.NewPaymentService(repo, publisher, cfg, gatewayClient)
	handler := paymentgrpc.NewPaymentHandler(svc)

	grpcPort := getEnvOrDefault("PAYMENT_GRPC_PORT", "50190")
	if grpcPort != "" && grpcPort[0] != ':' {
		grpcPort = ":" + grpcPort
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// ── AuthZ gRPC client ──────────────────────────────────────────────────────
	// Connect to the authz-service for Casbin policy checks. Dial is non-blocking;
	// the interceptor handles unavailability gracefully (fail-open with a warning).
	authzAddr := getEnvOrDefault("AUTHZ_GRPC_ADDR", "authz-service:50153")
	authzConn, err := grpc.NewClient(authzAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// Non-fatal: log and continue with fail-open behavior in the interceptor.
		logger.Warn("Failed to connect to authz-service, payment AuthZ will fail open", zap.String("addr", authzAddr), zap.Error(err))
	}
	var authzInterceptor *paymentmw.PaymentAuthZInterceptor
	if authzConn != nil {
		authzClient := authzservicev1.NewAuthZServiceClient(authzConn)
		authzInterceptor = paymentmw.NewPaymentAuthZInterceptor(authzClient)
	}

	// ── Build gRPC server with interceptors ───────────────────────────────────
	var serverOpts []grpc.ServerOption
	if authzInterceptor != nil {
		serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(
			authzInterceptor.UnaryServerInterceptor(),
		))
	}

	grpcServer := grpc.NewServer(serverOpts...)
	paymentservicev1.RegisterPaymentServiceServer(grpcServer, handler)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("payment", grpc_health_v1.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("Payment gRPC server error", zap.Error(err))
		}
	}()

	logger.Info("Payment service running", zap.String("grpc_port", grpcPort))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down payment service...")
	healthServer.SetServingStatus("payment", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()
	if kafkaProducer != nil {
		_ = kafkaProducer.Close()
	}
	if authzConn != nil {
		_ = authzConn.Close()
	}
	logger.Info("Payment service stopped")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}
