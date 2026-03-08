package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/events"
	fraudgrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/service"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/interceptors"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	opsconfig "github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"
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
	if err := appLogger.Initialize(appLogger.NoFileConfig()); err != nil {
		appLogger.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.GetLogger().Sync() //nolint:errcheck
	_ = env.Load()

	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatal("failed to load fraud config", zap.Error(err))
	}

	servicesConfigPath, err := opsconfig.ResolveConfigPath("services.yaml")
	if err != nil {
		appLogger.Fatal("failed to resolve services.yaml", zap.Error(err))
	}
	raw, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		appLogger.Fatal("failed to read services.yaml", zap.Error(err))
	}
	var svcConfig ServicesConfig
	if err := yaml.Unmarshal(raw, &svcConfig); err != nil {
		appLogger.Fatal("failed to parse services.yaml", zap.Error(err))
	}

	fraudSvc, ok := svcConfig.Services["fraud"]
	if !ok {
		appLogger.Fatal("configuration for 'fraud' service not found in services.yaml")
	}
	if os.Getenv("FRAUD_PORT") != "" || os.Getenv("FRAUD_GRPC_PORT") != "" || os.Getenv("FRAUD_HTTP_PORT") != "" {
		appLogger.Warn("FRAUD_PORT/FRAUD_GRPC_PORT/FRAUD_HTTP_PORT env values are ignored; using backend/inscore/configs/services.yaml")
	}
	grpcPort := fraudSvc.Ports.Grpc
	appLogger.Info("service configured from services.yaml",
		zap.String("service", fraudSvc.Name),
		zap.Int("grpc_port", fraudSvc.Ports.Grpc),
		zap.Int("http_port", fraudSvc.Ports.Http),
	)

	dbConfigPath, err := opsconfig.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatal("failed to resolve database config path", zap.Error(err))
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Fatal("failed to initialize database manager", zap.Error(err))
	}
	defer db.Manager.Close()

	database := db.GetDB()
	ruleRepo := repository.NewFraudRuleRepository(database)
	alertRepo := repository.NewFraudAlertRepository(database)
	caseRepo := repository.NewFraudCaseRepository(database)

	kafkaProducer, err := producer.NewEventProducerWithRetry(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		"fraud-service",
		5,
		3*time.Second,
	)
	if err != nil {
		appLogger.Warn("Kafka producer init failed — fraud events will be dropped", zap.Error(err))
		kafkaProducer = nil
	} else {
		defer kafkaProducer.Close()
	}
	eventPublisher := events.NewPublisher(kafkaProducer, cfg.Kafka.Topic)

	fraudSvcLogic := service.NewFraudService(ruleRepo, alertRepo, caseRepo, eventPublisher)
	consumer := events.NewConsumer(fraudSvcLogic)
	consumerGroup, consumerErr := kafkaconsumer.NewConsumerGroup(kafkaconsumer.Config{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.ConsumerGroup,
		Topics:  cfg.Kafka.ConsumerTopics,
		Handler: func(ctx context.Context, msg *kafkaconsumer.Message) error {
			return consumer.HandleMessage(ctx, msg.Topic, string(msg.Key), msg.Value)
		},
		DLQTopic: cfg.Kafka.DLQTopic,
		ClientID: "fraud-consumer",
	})
	if consumerErr != nil {
		appLogger.Warn("Kafka consumer group failed to start — fraud async checks disabled", zap.Error(consumerErr))
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
	handler := fraudgrpc.NewFraudHandler(fraudSvcLogic)

	authzAddr := resolveServiceAddr(cfg.Integration.AuthZAddress, svcConfig.Services, "authz")
	if authzAddr == "" {
		appLogger.Fatal("authz address is empty; configure AUTHZ_SERVICE_ADDRESS or services.yaml entry")
	}
	appLogger.Info("connecting to authz", zap.String("address", authzAddr))
	authzConn, err := grpc.Dial(authzAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		appLogger.Fatal("failed to connect to authz", zap.Error(err), zap.String("address", authzAddr))
	}
	defer authzConn.Close()
	authzClient := authzservicev1.NewAuthZServiceClient(authzConn)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.NewAuthZInterceptor(authzClient, interceptors.DefaultSkipMethods),
			loggingInterceptor(),
			recoveryInterceptor(),
		),
	)

	fraudservicev1.RegisterFraudServiceServer(grpcServer, handler)
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("fraud", grpc_health_v1.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)

	addr := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(grpcPort))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		appLogger.Fatal("failed to listen", zap.String("addr", addr), zap.Error(err))
	}

	appLogger.Info("fraud gRPC server listening", zap.String("addr", addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			appLogger.Error("gRPC serve error", zap.Error(err))
		}
	}()

	<-quit
	appLogger.Info("shutting down fraud service...")
	healthSrv.SetServingStatus("fraud", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
		appLogger.Info("fraud service stopped cleanly")
	case <-time.After(15 * time.Second):
		appLogger.Warn("graceful stop timed out — forcing")
		grpcServer.Stop()
	}
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

func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		appLogger.Info("grpc request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
		return resp, err
	}
}

func recoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				appLogger.Error("panic recovered",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
