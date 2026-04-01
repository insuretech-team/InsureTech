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
	authnconfig "github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/consumers"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/email"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/events"
	authnGrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/seeder"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/service"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/sms"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcclient"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafkaapp"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/runtimeaddr"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

// ServicesConfig structure matches services.yaml
type ServicesConfig = serviceaddr.ServicesConfig

func main() {
	// 1. Initialize Logger
	_ = appLogger.Initialize(appLogger.Config{
		Level:  "info",
		Format: "text",
		Output: "console",
	})
	appLogger.Info("Starting AuthN microservice (VSA)...")

	// 2. Load Configuration (Port)
	servicesConfigPath, err := config.ResolveConfigPath("services.yaml")
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

	authnConfig, exists := svcConfig.Services["authn"]
	if !exists {
		appLogger.Fatal("Configuration for 'authn' service not found in services.yaml")
	}
	port := strconv.Itoa(authnConfig.Ports.Grpc)
	if os.Getenv("AUTHN_PORT") != "" || os.Getenv("AUTHN_GRPC_PORT") != "" || os.Getenv("AUTHN_HTTP_PORT") != "" {
		appLogger.Warn("AUTHN_PORT/AUTHN_GRPC_PORT/AUTHN_HTTP_PORT env values are ignored; using backend/inscore/configs/services.yaml")
	}

	appLogger.Infof("Service '%s' configured on port %s", authnConfig.Name, port)

	// 3. Load AuthN Configuration
	cfg, err := authnconfig.Load()
	if err != nil {
		appLogger.Fatalf("Failed to load authn config: %v", err)
	}

	// 4. Initialize Infrastructure (DB)
	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve database config path: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Errorf("Failed to initialize database: %v", err)
		appLogger.Fatal("Database initialization failed")
	}
	defer db.Manager.Close()

	database := db.GetDB()

	// 5. Initialize Repositories
	sessionRepo := repository.NewSessionRepository(database)
	userRepo := repository.NewUserRepository(database)
	otpRepo := repository.NewOTPRepository(database)
	apiKeyRepo := repository.NewApiKeyRepository(database)
	userProfileRepo := repository.NewUserProfileRepository(database)
	userDocumentRepo := repository.NewUserDocumentRepository(database)
	documentTypeRepo := repository.NewDocumentTypeRepository(database)
	kycRepo := repository.NewKYCVerificationRepository(database)
	voiceRepo := repository.NewVoiceSessionRepository(database)

	// 6. Initialize Kafka Event Producer (with retry for startup resilience)
	kafkaBrokers := runtimeaddr.NormalizeKafkaBrokers(cfg.Kafka.Brokers)
	appLogger.Infof("Connecting to Kafka brokers: %v", kafkaBrokers)
	kafkaProducer, err := producer.NewEventProducerWithRetry(
		kafkaBrokers,
		cfg.Kafka.Topic,
		"authn-service",
		5,             // max retries
		3*time.Second, // retry delay
	)
	if err != nil {
		// Non-fatal: authn can run without Kafka but events won't be published
		appLogger.Errorf("Kafka producer initialization failed (events will be dropped): %v", err)
		kafkaProducer = nil
	} else {
		defer kafkaProducer.Close()
		appLogger.Info("Kafka producer initialized successfully")
	}
	eventPublisher := events.NewPublisher(kafkaProducer)

	// 7. Initialize SMS Client
	smsClient := sms.NewSSLWirelessClient(cfg)

	// 7b. Initialize Email Client (for Business Beneficiary + System User email OTP)
	// OTP / verification email client (noreply@labaidinsuretech.com)
	emailClient := email.NewClient(email.Config{
		SMTPHost: cfg.Email.SMTPHost,
		SMTPPort: cfg.Email.SMTPPort,
		From:     cfg.Email.From,
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
		TLS:      cfg.Email.TLS,
	})
	appLogger.Infof("Email client initialized (SMTP: %s:%d, from: %s)", cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.From)

	// Info / transactional email client (info@labaidinsuretech.com)
	emailInfoClient := email.NewClient(email.Config{
		SMTPHost: cfg.EmailInfo.SMTPHost,
		SMTPPort: cfg.EmailInfo.SMTPPort,
		From:     cfg.EmailInfo.From,
		Username: cfg.EmailInfo.Username,
		Password: cfg.EmailInfo.Password,
		TLS:      cfg.EmailInfo.TLS,
	})
	appLogger.Infof("Email info client initialized (SMTP: %s:%d, from: %s)", cfg.EmailInfo.SMTPHost, cfg.EmailInfo.SMTPPort, cfg.EmailInfo.From)
	_ = emailInfoClient // available for future transactional email use

	// 7d. Initialize Redis client (optional — used for JTI blocklist + session limiter)
	var redisClient redis.UniversalClient
	redisURL := runtimeaddr.NormalizeRedisURL(cfg.Redis.URL)
	if redisURL != cfg.Redis.URL {
		appLogger.Warnf("Redis URL normalized for runtime: %s -> %s", cfg.Redis.URL, redisURL)
	}
	if redisURL != "" {
		opt, parseErr := redis.ParseURL(redisURL)
		if parseErr != nil {
			appLogger.Warnf("Redis URL parse failed (%s): %v — running without Redis", redisURL, parseErr)
		} else {
			if cfg.Redis.Password != "" {
				opt.Password = cfg.Redis.Password
			}
			opt.DB = cfg.Redis.DB
			rdb := redis.NewClient(opt)
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := rdb.Ping(pingCtx).Err(); pingErr != nil {
				appLogger.Warnf("Redis ping failed (%s): %v — running without Redis", redisURL, pingErr)
			} else {
				redisClient = rdb
				appLogger.Infof("Redis connected: %s (db=%d)", redisURL, cfg.Redis.DB)
			}
			pingCancel()
		}
	} else {
		appLogger.Warn("REDIS_URL not set — JTI blocklist and session limiter disabled")
	}

	// 8. Initialize Services
	metadataExtractor := middleware.NewMetadataExtractor()
	// TokenService: use Redis-backed session limiter when available (JTI blocklist + concurrent session enforcement).
	// maxSessions ≤ 0 defaults to 5 in NewTokenServiceWithSessionLimiter.
	tokenService, err := service.NewTokenServiceWithSessionLimiter(sessionRepo, userRepo, cfg, eventPublisher, metadataExtractor, redisClient, 0)
	if err != nil {
		appLogger.Fatalf("failed to initialize token service: %v", err)
	}
	var otpService *service.OTPService
	if redisClient != nil {
		otpService = service.NewOTPServiceWithRedis(otpRepo, smsClient, emailClient, cfg, eventPublisher, redisClient)
	} else {
		otpService = service.NewOTPService(otpRepo, smsClient, emailClient, cfg, eventPublisher)
	}
	authService := service.NewAuthService(
		tokenService,
		otpService,
		userRepo,
		sessionRepo,
		otpRepo,
		apiKeyRepo,
		userProfileRepo,
		userDocumentRepo,
		documentTypeRepo,
		kycRepo,
		voiceRepo,
		eventPublisher,
		cfg,
		metadataExtractor,
	)

	// Wire purpose-built FLVEAdapter when HF endpoint is configured.
	if cfg.FLVE.HFEndpoint != "" {
		timeout := cfg.KYC.Timeout
		if timeout <= 0 {
			timeout = 30 * time.Second
		}
		appLogger.Infof("FLVE eKYC adapter enabled: endpoint=%s token_len=%d", cfg.FLVE.HFEndpoint, len(cfg.FLVE.HFToken))
		authService.SetFLVEAdapter(service.NewFLVEAdapter(cfg.FLVE.HFEndpoint, cfg.FLVE.HFToken, timeout))
	}

	// Downstream KYC microservice wiring.
	// Reads KYC_SERVICE_ADDRESS directly from OS env (not via config struct) to avoid
	// godotenv not-overwrite behaviour when KYC_SERVICE_ENABLED was previously false
	// in the shell environment. Address presence is the only gate needed.
	var kycConn *grpc.ClientConn
	kycAddress := os.Getenv("KYC_SERVICE_ADDRESS")
	if kycAddress == "" {
		kycAddress = cfg.KYC.Address
	}
	if kycAddress != "" {
		addressLower := strings.ToLower(kycAddress)
		if strings.HasPrefix(addressLower, "http://") || strings.HasPrefix(addressLower, "https://") {
			authService.SetExternalKYCClient(service.NewFLVEExternalKYCClient(kycAddress, cfg.KYC.Token, cfg.KYC.Timeout))
			appLogger.Infof("Downstream FLVE KYC client (legacy) enabled: %s", kycAddress)
		} else {
			// Non-blocking dial — gRPC will reconnect automatically when the KYC service is ready.
			// WithBlock() caused silent failures when KYC service wasn't up yet during authn startup.
			conn, dialErr := grpcclient.NewClient(kycAddress)
			if dialErr != nil {
				appLogger.Warnf("Downstream KYC client setup failed (%s): %v — using local KYC repository path", kycAddress, dialErr)
			} else {
				kycConn = conn
				authService.SetExternalKYCClient(kycservicev1.NewKYCServiceClient(conn))
				appLogger.Infof("Downstream KYC gRPC client enabled: %s", kycAddress)
			}
		}
	} else {
		appLogger.Warn("KYC_SERVICE_ADDRESS not set — using local KYC repository (no FLVE eKYC)")
	}
	if kycConn != nil {
		defer func() { _ = kycConn.Close() }()
	}

	// 7c. Initialize Kafka Consumer Group (authn domain topics + authz.events)
	const topicAuthzEvents = "authz.events"
	fanOut := consumers.FanOutHandler(consumers.TopicHandlers{
		events.TopicSMSDeliveryReport:      consumers.NewSMSDLRHandler(otpRepo),
		events.TopicAccountLocked:          consumers.NewAccountLockedHandler(userRepo, smsClient),
		events.TopicUserRegistered:         consumers.NewUserRegisteredHandler(emailClient, smsClient),
		events.TopicPasswordChanged:        consumers.NewPasswordChangedHandler(userRepo, smsClient),
		events.TopicPasswordResetRequested: consumers.NewPasswordResetRequestedHandler(userRepo, smsClient),
		events.TopicSessionRevoked:         consumers.NewSessionRevokedAllHandler(userRepo, smsClient),
		// Sprint 1.9: consume PortalConfigUpdatedEvent from AuthZ to keep local
		// portal config cache (MFA requirements, session limits, TTLs) up-to-date
		// without synchronous gRPC calls on every login.
		topicAuthzEvents: consumers.NewPortalConfigUpdatedHandler(),
	})
	consumerTopics := []string{
		events.TopicSMSDeliveryReport,
		events.TopicAccountLocked,
		events.TopicUserRegistered,
		events.TopicPasswordChanged,
		events.TopicPasswordResetRequested,
		events.TopicSessionRevoked,
		topicAuthzEvents,
	}
	consumerGroup, consumerErr := kafkaapp.StartConsumerGroup(kafkaconsumer.Config{
		Brokers:  kafkaBrokers,
		GroupID:  "authn-service-consumer",
		Topics:   consumerTopics,
		Handler:  fanOut,
		DLQTopic: "authn.dlq",
		ClientID: "authn-consumer",
	})
	if consumerErr != nil {
		appLogger.Warnf("Kafka consumer group failed to start (events will not be consumed): %v", consumerErr)
	} else {
		defer func() {
			_ = consumerGroup.Close()
		}()
		appLogger.Infof("Kafka consumer group started (topics=%v)", consumerTopics)
	}

	// 8a. Seed default admin user (idempotent)
	if err := seeder.SeedAdminUser(context.Background(), database); err != nil {
		appLogger.Warnf("Admin seeder: %v", err)
	}

	if err := seeder.SeedB2bAdminUser(context.Background(), database); err != nil {
		appLogger.Warnf("B2B Admin seeder: %v", err)
	}

	// 8b. Background cleanup jobs (sessions + OTPs)
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	defer cleanupCancel()
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// Expired sessions
				count, err := sessionRepo.CleanupExpiredSessions(cleanupCtx)
				if err != nil {
					appLogger.Errorf("Session cleanup error: %v", err)
				} else if count > 0 {
					appLogger.Infof("Cleaned up %d expired sessions", count)
				}
				// Expired OTPs (older than 24h)
				otpCount, err := otpRepo.CleanupExpiredOTPs(cleanupCtx, time.Now().Add(-24*time.Hour))
				if err != nil {
					appLogger.Errorf("OTP cleanup error: %v", err)
				} else if otpCount > 0 {
					appLogger.Infof("Cleaned up %d expired OTPs", otpCount)
				}
			case <-cleanupCtx.Done():
				return
			}
		}
	}()

	// 8. Initialize gRPC Server
	serverConfig := authnGrpc.DefaultServerConfig()
	serverConfig.Host = cfg.Server.Host
	serverConfig.Port = port
	serverConfig.DB = database

	server, err := authnGrpc.NewServer(serverConfig, authService)
	if err != nil {
		appLogger.Fatalf("Failed to create gRPC server: %v", err)
	}

	// 8. Health Check (retry for transient cloud DB latency)
	var healthErr error
	for attempt := 1; attempt <= 3; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		healthErr = server.HealthCheck(ctx)
		cancel()
		if healthErr == nil {
			break
		}
		appLogger.Warnf("Server health check attempt %d/3 failed: %v", attempt, healthErr)
		if attempt < 3 {
			time.Sleep(2 * time.Second)
		}
	}
	if healthErr != nil {
		appLogger.Fatalf("Server health check failed after retries: %v", healthErr)
	}

	// 9. Start Server
	go func() {
		if err := server.Start(); err != nil {
			appLogger.Fatalf("Server crashed: %v", err)
		}
	}()

	// 10. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down...")
	server.Stop()
	appLogger.Info("Stopped.")
}
