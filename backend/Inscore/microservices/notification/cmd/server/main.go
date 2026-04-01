package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	notificationconfig "github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/delivery"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/events"
	notificationgrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/notification/internal/service"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcclient"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafkaapp"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/runtimeaddr"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type servicesConfig = serviceaddr.ServicesConfig

func main() {
	_ = appLogger.Initialize(appLogger.Config{
		Level:  "info",
		Format: "text",
		Output: "console",
	})
	appLogger.Info("Starting notification microservice...")

	cfg, err := notificationconfig.Load()
	if err != nil {
		appLogger.Fatalf("Failed to load notification config: %v", err)
	}

	servicesConfigPath, err := config.ResolveConfigPath("services.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve services.yaml path: %v", err)
	}
	servicesData, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		appLogger.Fatalf("Failed to read services.yaml: %v", err)
	}
	var svcConfig servicesConfig
	if err := yaml.Unmarshal(servicesData, &svcConfig); err != nil {
		appLogger.Fatalf("Failed to parse services.yaml: %v", err)
	}

	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve database config path: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Fatalf("Failed to initialize database manager: %v", err)
	}
	defer db.Manager.Close()

	database := db.GetDB()
	notificationRepo := repository.NewNotificationRepository(database)
	templateRepo := repository.NewTemplateRepository(database)
	userRepo := repository.NewUserRepository(database)
	lookupRepo := repository.NewLookupRepository(database)
	pushTokenRepo := repository.NewPushTokenRepository(database)
	webhookRepo := repository.NewWebhookRepository(database)

	kafkaBrokers := runtimeaddr.NormalizeKafkaBrokers(cfg.Kafka.Brokers)
	kafkaProducer, err := producer.NewEventProducerWithRetry(kafkaBrokers, cfg.Kafka.Topic, "notification-service", 5, 3*time.Second)
	if err != nil {
		appLogger.Warnf("Kafka producer initialization failed (notification state events will be dropped): %v", err)
		kafkaProducer = nil
	} else {
		defer kafkaProducer.Close()
	}

	eventPublisher := events.NewPublisher(kafkaProducer)
	notificationService := service.NewService(
		notificationRepo,
		templateRepo,
		userRepo,
		lookupRepo,
		pushTokenRepo,
		webhookRepo,
		delivery.NewEmailClient(delivery.EmailConfig(cfg.Email)),
		delivery.NewSMSClient(delivery.SMSConfig(cfg.SMS)),
		delivery.NewPushClient(delivery.PushConfig(cfg.Push)),
		delivery.NewWebhookClient(delivery.WebhookConfig{
			Enabled:   cfg.Webhook.Enabled,
			Timeout:   cfg.Webhook.Timeout,
			UserAgent: cfg.Webhook.UserAgent,
		}),
		eventPublisher,
		cfg,
	)

	authnAddr := serviceaddr.ResolveFromServicesMap("", svcConfig.Services, os.Getenv("NOTIFICATION_SERVICE_DISCOVERY_HOST"), "authn")
	if authnAddr == "" {
		authnAddr = serviceaddr.ResolveFromServicesMap("", svcConfig.Services, os.Getenv("SERVICE_DISCOVERY_HOST"), "authn")
	}
	if authnAddr != "" {
		authnConn, err := grpcclient.NewClient(authnAddr)
		if err != nil {
			appLogger.Warnf("Failed to dial authn-service at %s — preference updates will be unavailable: %v", authnAddr, err)
		} else {
			defer authnConn.Close()
			notificationService.WithAuthNPreferenceClient(authnservicev1.NewAuthServiceClient(authnConn))
			appLogger.Info("AuthN preference client wired", zap.String("addr", authnAddr))
		}
	} else {
		appLogger.Warn("AuthN service address not resolved — preference updates will be unavailable")
	}

	dispatchCtx, dispatchCancel := context.WithCancel(context.Background())
	defer dispatchCancel()
	notificationService.StartDispatcher(dispatchCtx)

	plan := events.SubscriptionPlan{
		Profile:               events.SubscriptionProfile(cfg.Topics.SubscriptionProfile),
		EnabledGroups:         cfg.Topics.EnabledGroups,
		DisabledGroups:        cfg.Topics.DisabledGroups,
		AllowTopics:           cfg.Topics.AllowTopics,
		DenyTopics:            cfg.Topics.DenyTopics,
		ExtraTopics:           cfg.Topics.ExtraTopics,
		IncludeReservedTopics: cfg.Topics.IncludeReserved,
	}
	consumerTopics := events.ConsumerTopicsForPlan(plan)

	var consumerGroup *kafkaapp.ManagedConsumer
	if len(consumerTopics) > 0 {
		eventConsumer := events.NewConsumer(notificationService)
		consumerGroup, err = kafkaapp.StartConsumerGroup(kafkaconsumer.Config{
			Brokers:  kafkaBrokers,
			GroupID:  cfg.Topics.ConsumerGroupID,
			Topics:   consumerTopics,
			DLQTopic: cfg.Topics.DLQTopic,
			ClientID: "notification-consumer",
			Handler: func(ctx context.Context, msg *kafkaconsumer.Message) error {
				return eventConsumer.HandleMessage(ctx, msg.Topic, msg.Value)
			},
		})
		if err != nil {
			appLogger.Warnf("Kafka consumer initialization failed: %v", err)
		} else {
			defer func() { _ = consumerGroup.Close() }()
			appLogger.Infof("Notification Kafka consumer started for topics: %v", consumerTopics)
		}
	}

	serverCfg := notificationgrpc.DefaultServerConfig()
	serverCfg.Host = cfg.Server.Host
	serverCfg.Port = cfg.Server.GRPCPort
	serverCfg.DB = database

	server, err := notificationgrpc.NewServer(serverCfg, notificationService)
	if err != nil {
		appLogger.Fatalf("Failed to create notification gRPC server: %v", err)
	}

	healthCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.HealthCheck(healthCtx); err != nil {
		appLogger.Fatalf("Notification server health check failed: %v", err)
	}

	go func() {
		if err := server.Start(); err != nil {
			appLogger.Fatalf("Notification server crashed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down notification service...")
	dispatchCancel()
	server.Stop()
	notificationService.WaitForDispatcher()
	appLogger.Info("Notification service stopped.")
}
