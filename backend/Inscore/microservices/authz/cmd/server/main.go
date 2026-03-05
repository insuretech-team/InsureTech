package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/enforcer"
	authzEvents "github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/events"
	authzgrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/seeder"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/service"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	opsconfig "github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
)

// ServicesConfig matches backend/inscore/configs/services.yaml structure.
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
	// ── Logger ────────────────────────────────────────────────────────────────
	if err := appLogger.Initialize(appLogger.NoFileConfig()); err != nil {
		appLogger.Fatalf("failed to initialize logger: %v", err)
	}
	layerLogger := appLogger.GetLogger().WithOptions(zap.AddCallerSkip(-1))
	defer layerLogger.Sync() //nolint:errcheck

	// ── Load Environment Variables ───────────────────────────────────────────
	if err := env.Load(); err != nil {
		appLogger.Warn("No .env file found, using system environment variables", zap.Error(err))
	} else {
		appLogger.Info(".env file loaded successfully")
	}

	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatal("failed to load config", zap.Error(err))
	}

	// ── Service Port Config (strictly from services.yaml) ────────────────────
	servicesConfigPath, err := opsconfig.ResolveConfigPath("services.yaml")
	if err != nil {
		appLogger.Fatal("failed to resolve services.yaml path", zap.Error(err))
	}
	servicesData, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		appLogger.Fatal("failed to read services.yaml", zap.Error(err))
	}
	var svcConfig ServicesConfig
	if err := yaml.Unmarshal(servicesData, &svcConfig); err != nil {
		appLogger.Fatal("failed to parse services.yaml", zap.Error(err))
	}
	authzServiceConfig, exists := svcConfig.Services["authz"]
	if !exists {
		appLogger.Fatal("configuration for 'authz' service not found in services.yaml")
	}
	if os.Getenv("AUTHZ_PORT") != "" || os.Getenv("AUTHZ_GRPC_PORT") != "" || os.Getenv("AUTHZ_HTTP_PORT") != "" {
		appLogger.Warn("AUTHZ_PORT/AUTHZ_GRPC_PORT/AUTHZ_HTTP_PORT env values are ignored; using backend/inscore/configs/services.yaml")
	}
	grpcPort := authzServiceConfig.Ports.Grpc
	appLogger.Info("service configured from services.yaml",
		zap.String("service", authzServiceConfig.Name),
		zap.Int("grpc_port", authzServiceConfig.Ports.Grpc),
		zap.Int("http_port", authzServiceConfig.Ports.Http),
	)

	// ── Database ──────────────────────────────────────────────────────────────
	dbConfigPath, err := opsconfig.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatal("failed to resolve database config path", zap.Error(err))
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Fatal("failed to initialize database manager for service", zap.Error(err))
	}
	defer db.Manager.Close()
	database := db.GetDB()

	// ── Casbin Enforcer ───────────────────────────────────────────────────────
	enf, err := enforcer.New(database, cfg.Casbin.ModelPath)
	if err != nil {
		appLogger.Fatal("failed to init casbin enforcer", zap.Error(err))
	}
	if cfg.Casbin.AutoReloadInterval > 0 {
		enf.StartAutoReload(int(cfg.Casbin.AutoReloadInterval.Seconds()))
	}

	// ── Repositories ──────────────────────────────────────────────────────────
	roleRepo := repository.NewRoleRepo(database)
	userRoleRepo := repository.NewUserRoleRepo(database)
	policyRepo := repository.NewPolicyRepo(database)
	portalRepo := repository.NewPortalRepo(database)
	auditRepo := repository.NewAuditRepo(database)
	tokenConfigRepo := repository.NewTokenConfigRepo(database)

	// ── Seeder ────────────────────────────────────────────────────────────────
	sd := seeder.New(roleRepo, policyRepo, enf, portalRepo, tokenConfigRepo, database, layerLogger)

	// 📡 Kafka Producer ────────────────────────────────────────────────────────
	kafkaProducer, err := producer.NewEventProducerWithRetry(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		"authz-service",
		5,
		3*time.Second,
	)
	if err != nil {
		appLogger.Warn("Kafka producer init failed — events will be dropped", zap.Error(err))
		kafkaProducer = nil
	} else {
		defer kafkaProducer.Close()
		appLogger.Info("Kafka producer initialized")
	}
	eventPublisher := authzEvents.NewPublisher(kafkaProducer)

	// 🔧 Service ───────────────────────────────────────────────────────────────
	svc := service.New(
		enf,
		roleRepo,
		userRoleRepo,
		policyRepo,
		portalRepo,
		auditRepo,
		cfg.Casbin.AuditAllDecisions,
		eventPublisher,
	)

	// 📥 Kafka Consumer Group ──────────────────────────────────────────────────
	consumerTopics := []string{
		"authn.user.registered",
		authzEvents.TopicAuthZEvents,
	}
	fanOut := authzEvents.FanOutHandler(authzEvents.TopicHandlers{
		"authn.user.registered":      authzEvents.NewUserRegisteredHandler(enf),
		authzEvents.TopicAuthZEvents: authzEvents.NewPolicyCacheInvalidatedHandler(enf),
	})
	consumerGroup, consumerErr := kafkaconsumer.NewConsumerGroup(kafkaconsumer.Config{
		Brokers:  cfg.Kafka.Brokers,
		GroupID:  "authz-service-consumer",
		Topics:   consumerTopics,
		Handler:  fanOut,
		DLQTopic: "authz.dlq",
		ClientID: "authz-consumer",
	})
	if consumerErr != nil {
		appLogger.Warn("Kafka consumer group failed to start — events will not be consumed", zap.Error(consumerErr))
	} else {
		consumerCtx, consumerCancel := context.WithCancel(context.Background())
		defer consumerCancel()
		go consumerGroup.Start(consumerCtx)
		defer func() {
			consumerCancel()
			_ = consumerGroup.Close()
		}()
		appLogger.Info("Kafka consumer group started", zap.Strings("topics", consumerTopics))
	}

	// ── Auto-seed on startup ──────────────────────────────────────────────────
	seedCtx, seedCancel := context.WithTimeout(context.Background(), 60*time.Second)
	if err := sd.SeedAllPortals(seedCtx); err != nil {
		appLogger.Warn("portal seeding had errors (non-fatal)", zap.Error(err))
	}
	seedCancel()

	// ── JWT Interceptor ───────────────────────────────────────────────────────
	// Parse RSA public key from config. If the PEM is empty (dev/test), the
	// interceptor operates in no-op mode (all requests pass through).
	publicKey, err := middleware.ParseRSAPublicKeyFromPEM(cfg.Auth.PublicKeyPEM)
	if err != nil {
		appLogger.Fatal("failed to parse RSA public key for JWT interceptor", zap.Error(err))
	}
	if publicKey == nil {
		appLogger.Warn("AUTHZ_JWT_PUBLIC_KEY_PEM not set — JWT validation disabled (no-op mode)")
	}

	// Methods that are publicly accessible without a JWT.
	publicMethods := []string{
		"/insuretech.authz.services.v1.AuthZService/CheckPermission",
		"/insuretech.authz.services.v1.AuthZService/CheckAccess",
		"/insuretech.authz.services.v1.AuthZService/GetJWKS",
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
	}

	jwtInterceptor := middleware.NewJWTInterceptor(publicKey, publicMethods)
	rateLimiter := middleware.NewRateLimiter(100, 200) // 100 rps steady, burst 200

	// ── gRPC Server ───────────────────────────────────────────────────────────
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			jwtInterceptor.UnaryInterceptor(),
			rateLimiter.UnaryInterceptor(),
			loggingInterceptor(),
			recoveryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			jwtInterceptor.StreamInterceptor(),
			rateLimiter.StreamInterceptor(),
		),
	)

	authzservicev1.RegisterAuthZServiceServer(grpcServer, authzgrpc.NewAuthZHandler(svc))

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("authz", grpc_health_v1.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)

	// ── Listen ────────────────────────────────────────────────────────────────
	addr := cfg.Server.Host + ":" + strconv.Itoa(grpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		appLogger.Fatal("failed to listen", zap.String("addr", addr), zap.Error(err))
	}
	appLogger.Info("authz gRPC server listening", zap.String("addr", addr))

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			appLogger.Error("gRPC serve error", zap.Error(err))
		}
	}()

	<-quit
	appLogger.Info("shutting down authz service...")
	healthSrv.SetServingStatus("authz", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
		appLogger.Info("authz service stopped cleanly")
	case <-time.After(15 * time.Second):
		appLogger.Warn("graceful stop timed out — forcing")
		grpcServer.Stop()
	}
}

func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		appLogger.Info("grpc", zap.String("method", info.FullMethod), zap.Duration("dur", time.Since(start)), zap.Error(err))
		return resp, err
	}
}

func recoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				appLogger.Error("panic in handler", zap.String("method", info.FullMethod), zap.Any("panic", r))
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
