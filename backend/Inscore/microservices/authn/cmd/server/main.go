package main

import (
	"context"
	"net"
	"net/url"
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
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/producer"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

// ServicesConfig structure matches services.yaml
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
	kafkaBrokers := normalizeKafkaBrokers(cfg.Kafka.Brokers)
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
	emailClient := email.NewClient(email.Config{
		SMTPHost: cfg.Email.SMTPHost,
		SMTPPort: cfg.Email.SMTPPort,
		From:     cfg.Email.From,
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
		TLS:      cfg.Email.TLS,
	})
	appLogger.Infof("Email client initialized (SMTP: %s:%d, from: %s)", cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.From)

	// 7d. Initialize Redis client (optional — used for JTI blocklist + session limiter)
	var redisClient redis.UniversalClient
	redisURL := normalizeRedisURL(cfg.Redis.URL)
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
	otpService := service.NewOTPService(otpRepo, smsClient, emailClient, cfg, eventPublisher)
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

	// Optional downstream KYC client wiring (Phase B).
	var kycConn *grpc.ClientConn
	if cfg.KYC.Enabled && cfg.KYC.Address != "" {
		addressLower := strings.ToLower(cfg.KYC.Address)
		if strings.HasPrefix(addressLower, "http://") || strings.HasPrefix(addressLower, "https://") {
			authService.SetExternalKYCClient(service.NewFLVEExternalKYCClient(cfg.KYC.Address, cfg.KYC.Token, cfg.KYC.Timeout))
			appLogger.Infof("Downstream FLVE KYC client enabled: %s", cfg.KYC.Address)
		} else {
			dialCtx, dialCancel := context.WithTimeout(context.Background(), cfg.KYC.Timeout)
			defer dialCancel()
			conn, dialErr := grpc.DialContext(
				dialCtx,
				cfg.KYC.Address,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			)
			if dialErr != nil {
				appLogger.Warnf("Downstream KYC dial failed (%s): %v — using local KYC repository path", cfg.KYC.Address, dialErr)
			} else {
				kycConn = conn
				authService.SetExternalKYCClient(kycservicev1.NewKYCServiceClient(conn))
				appLogger.Infof("Downstream KYC client enabled: %s", cfg.KYC.Address)
			}
		}
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
	consumerGroup, consumerErr := kafkaconsumer.NewConsumerGroup(kafkaconsumer.Config{
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
		consumerCtx, consumerCancel := context.WithCancel(context.Background())
		defer consumerCancel()
		go consumerGroup.Start(consumerCtx)
		defer func() {
			consumerCancel()
			_ = consumerGroup.Close()
		}()
		appLogger.Infof("Kafka consumer group started (topics=%v)", consumerTopics)
	}

	// 8a. Seed default admin user (idempotent)
	if err := seeder.SeedAdminUser(context.Background(), database); err != nil {
		appLogger.Warnf("Admin seeder: %v", err)
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

func normalizeKafkaBrokers(brokers []string) []string {
	if len(brokers) == 0 {
		return brokers
	}
	inDocker := isRunningInDocker()
	inWSL := isRunningInWSL()

	normalized := make([]string, 0, len(brokers)*2)
	for _, broker := range brokers {
		if broker == "" {
			continue
		}

		host, port, err := net.SplitHostPort(broker)
		if err == nil {
			if host == "localhost" || host == "127.0.0.1" || host == "::1" {
				if inDocker {
					normalized = append(normalized, net.JoinHostPort("kafka", port))
					normalized = append(normalized, net.JoinHostPort("host.docker.internal", port))
					normalized = append(normalized, broker)
				} else if inWSL {
					normalized = append(normalized, net.JoinHostPort("host.docker.internal", port))
					normalized = append(normalized, broker)
				} else {
					normalized = append(normalized, broker)
					normalized = append(normalized, net.JoinHostPort("host.docker.internal", port))
				}
				continue
			}
			if host == "kafka" {
				if inDocker {
					normalized = append(normalized, broker)
					normalized = append(normalized, net.JoinHostPort("host.docker.internal", port))
				} else {
					normalized = append(normalized, net.JoinHostPort("host.docker.internal", port))
					normalized = append(normalized, net.JoinHostPort("localhost", port))
				}
				continue
			}
			if host == "host.docker.internal" {
				if inDocker {
					normalized = append(normalized, broker)
					normalized = append(normalized, net.JoinHostPort("kafka", port))
				} else {
					normalized = append(normalized, broker)
					normalized = append(normalized, net.JoinHostPort("localhost", port))
				}
				continue
			}
			normalized = append(normalized, broker)
		} else {
			if broker == "localhost" || broker == "127.0.0.1" {
				if inDocker {
					normalized = append(normalized, "kafka:9092", "host.docker.internal:9092", "localhost:9092")
				} else if inWSL {
					normalized = append(normalized, "host.docker.internal:9092", "localhost:9092")
				} else {
					normalized = append(normalized, "localhost:9092", "host.docker.internal:9092")
				}
				continue
			}
			if broker == "kafka" {
				if inDocker {
					normalized = append(normalized, "kafka:9092", "host.docker.internal:9092")
				} else {
					normalized = append(normalized, "host.docker.internal:9092", "localhost:9092")
				}
				continue
			}
			normalized = append(normalized, broker)
		}
	}

	deduped := dedupeStrings(normalized)
	if !equalStringSlices(brokers, deduped) {
		appLogger.Warnf("Kafka brokers normalized with runtime fallbacks: %v -> %v", brokers, deduped)
	}
	return deduped
}

func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func isRunningInWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSL_INTEROP") != "" {
		return true
	}
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(data))
	return strings.Contains(lower, "microsoft")
}

func dedupeStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func normalizeRedisURL(raw string) string {
	if raw == "" {
		return raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	host := parsed.Hostname()
	port := parsed.Port()
	if port == "" {
		port = "6379"
	}

	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		if isRunningInDocker() {
			parsed.Host = net.JoinHostPort("redis", port)
		} else if isRunningInWSL() {
			parsed.Host = net.JoinHostPort("host.docker.internal", port)
		}
	}
	return parsed.String()
}
