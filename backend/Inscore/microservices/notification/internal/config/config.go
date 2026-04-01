package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	opsconfig "github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"gopkg.in/yaml.v3"
)

const (
	defaultNotificationGRPCPort = 50230
	defaultNotificationHTTPPort = 50231
)

type Config struct {
	Server            ServerConfig
	Database          DatabaseConfig
	Kafka             KafkaConfig
	Topics            TopicConfig
	Email             EmailConfig
	SMS               SMSConfig
	Push              PushConfig
	Webhook           WebhookConfig
	Delivery          DeliveryConfig
	NotificationRetry NotificationRetryConfig
}

type ServerConfig struct {
	GRPCPort int
	HTTPPort int
	Host     string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type TopicConfig struct {
	SubscriptionProfile string
	EnabledGroups       []string
	DisabledGroups      []string
	AllowTopics         []string
	DenyTopics          []string
	ExtraTopics         []string
	IncludeReserved     bool
	ConsumerGroupID     string
	DLQTopic            string
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	From     string
	Username string
	Password string
	TLS      bool
}

type SMSConfig struct {
	Provider          string
	APIBase           string
	SID               string
	APIKey            string
	MaskingEnabled    bool
	MaskingSenderID   string
	NonMaskingEnabled bool
	NonMaskingSender  string
}

type PushConfig struct {
	Provider    string
	Endpoint    string
	ServerKey   string
	Timeout     time.Duration
	MockSuccess bool
}

type WebhookConfig struct {
	Enabled     bool
	Timeout     time.Duration
	BatchSize   int
	MaxAttempts int
	Backoff     []time.Duration
	UserAgent   string
}

type DeliveryConfig struct {
	DispatchInterval time.Duration
	BatchSize        int
}

type NotificationRetryConfig struct {
	MaxAttempts int
	Backoff     []time.Duration
}

func Load() (*Config, error) {
	if err := env.Load(); err != nil {
		logger.Warnf("Failed to load .env file: %v (using system environment variables)", err)
	}

	grpcPort, httpPort := loadNotificationServicePorts()
	if os.Getenv("NOTIFICATION_PORT") != "" || os.Getenv("NOTIFICATION_GRPC_PORT") != "" || os.Getenv("NOTIFICATION_HTTP_PORT") != "" {
		logger.Warn("NOTIFICATION_PORT/NOTIFICATION_GRPC_PORT/NOTIFICATION_HTTP_PORT are ignored; notification ports are loaded from backend/inscore/configs/services.yaml")
	}

	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: grpcPort,
			HTTPPort: httpPort,
			Host:     getEnv("NOTIFICATION_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "insuretech_primary"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "insuretech_primary"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers: getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:   getEnv("KAFKA_NOTIFICATION_TOPIC", "notification-events"),
		},
		Topics: TopicConfig{
			SubscriptionProfile: getEnv("NOTIFICATION_SUBSCRIPTION_PROFILE", "customer_core"),
			EnabledGroups:       getEnvAsSlice("NOTIFICATION_TOPIC_GROUPS", nil),
			DisabledGroups:      getEnvAsSlice("NOTIFICATION_DISABLED_TOPIC_GROUPS", nil),
			AllowTopics:         getEnvAsSlice("NOTIFICATION_TOPIC_ALLOWLIST", nil),
			DenyTopics:          getEnvAsSlice("NOTIFICATION_TOPIC_DENYLIST", nil),
			ExtraTopics:         getEnvAsSlice("NOTIFICATION_EXTRA_TOPICS", nil),
			IncludeReserved:     getEnvAsBool("NOTIFICATION_INCLUDE_RESERVED_TOPICS", false),
			ConsumerGroupID:     getEnv("NOTIFICATION_CONSUMER_GROUP_ID", "notification-service-consumer"),
			DLQTopic:            getEnv("NOTIFICATION_CONSUMER_DLQ_TOPIC", "notification.dlq"),
		},
		Email: EmailConfig{
			SMTPHost: getEnv("EMAIL_INFO_SMTP_HOST", getEnv("EMAIL_SMTP_HOST", "smtp.gmail.com")),
			SMTPPort: getEnvAsInt("EMAIL_INFO_SMTP_PORT", getEnvAsInt("EMAIL_SMTP_PORT", 587)),
			From:     getEnv("EMAIL_INFO_FROM", getEnv("EMAIL_FROM", "info@labaidinsuretech.com")),
			Username: getEnv("EMAIL_INFO_USERNAME", getEnv("EMAIL_USERNAME", "")),
			Password: getEnv("EMAIL_INFO_PASSWORD", getEnv("EMAIL_PASSWORD", "")),
			TLS:      getEnvAsBool("EMAIL_INFO_TLS", getEnvAsBool("EMAIL_TLS", false)),
		},
		SMS: SMSConfig{
			Provider:          getEnv("SMS_PROVIDER", "sslwireless"),
			APIBase:           getEnv("SSLWIRELESS_API_BASE", ""),
			SID:               getEnv("SSLWIRELESS_SID", ""),
			APIKey:            getEnv("SSLWIRELESS_API_KEY", ""),
			MaskingEnabled:    getEnvAsBool("SSLWIRELESS_MASKING_ENABLED", true),
			MaskingSenderID:   getEnv("SSLWIRELESS_SENDER_ID", "LABAIDINS"),
			NonMaskingEnabled: getEnvAsBool("SSLWIRELESS_NONMASKING_ENABLED", true),
			NonMaskingSender:  getEnv("SSLWIRELESS_NONMASKING_SENDER", ""),
		},
		Push: PushConfig{
			Provider:    strings.ToLower(getEnv("NOTIFICATION_PUSH_PROVIDER", "")),
			Endpoint:    getEnv("NOTIFICATION_PUSH_ENDPOINT", "https://fcm.googleapis.com/fcm/send"),
			ServerKey:   getEnv("NOTIFICATION_PUSH_SERVER_KEY", ""),
			Timeout:     getEnvAsDuration("NOTIFICATION_PUSH_TIMEOUT", 10*time.Second),
			MockSuccess: getEnvAsBool("NOTIFICATION_PUSH_MOCK_SUCCESS", false),
		},
		Webhook: WebhookConfig{
			Enabled:     getEnvAsBool("NOTIFICATION_WEBHOOK_ENABLED", true),
			Timeout:     getEnvAsDuration("NOTIFICATION_WEBHOOK_TIMEOUT", 10*time.Second),
			BatchSize:   getEnvAsInt("NOTIFICATION_WEBHOOK_BATCH_SIZE", 50),
			MaxAttempts: getEnvAsInt("NOTIFICATION_WEBHOOK_MAX_ATTEMPTS", 5),
			Backoff: []time.Duration{
				getEnvAsDuration("NOTIFICATION_WEBHOOK_BACKOFF_1", 30*time.Second),
				getEnvAsDuration("NOTIFICATION_WEBHOOK_BACKOFF_2", 2*time.Minute),
				getEnvAsDuration("NOTIFICATION_WEBHOOK_BACKOFF_3", 10*time.Minute),
			},
			UserAgent: getEnv("NOTIFICATION_WEBHOOK_USER_AGENT", "insuretech-notification-webhook/1.0"),
		},
		Delivery: DeliveryConfig{
			DispatchInterval: getEnvAsDuration("NOTIFICATION_DISPATCH_INTERVAL", 15*time.Second),
			BatchSize:        getEnvAsInt("NOTIFICATION_DISPATCH_BATCH_SIZE", 50),
		},
		NotificationRetry: NotificationRetryConfig{
			MaxAttempts: getEnvAsInt("NOTIFICATION_MAX_RETRY_ATTEMPTS", 3),
			Backoff: []time.Duration{
				getEnvAsDuration("NOTIFICATION_RETRY_BACKOFF_1", time.Minute),
				getEnvAsDuration("NOTIFICATION_RETRY_BACKOFF_2", 5*time.Minute),
				getEnvAsDuration("NOTIFICATION_RETRY_BACKOFF_3", 15*time.Minute),
			},
		},
	}

	if err := cfg.Validate(); err != nil {
		logger.Errorf("config validation failed: %v", err)
		return nil, errors.New("config validation failed")
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return errors.New("DB_PASSWORD is required")
	}
	if c.Delivery.BatchSize <= 0 {
		return errors.New("NOTIFICATION_DISPATCH_BATCH_SIZE must be greater than 0")
	}
	if c.NotificationRetry.MaxAttempts <= 0 {
		return errors.New("NOTIFICATION_MAX_RETRY_ATTEMPTS must be greater than 0")
	}
	if c.Webhook.BatchSize <= 0 {
		return errors.New("NOTIFICATION_WEBHOOK_BATCH_SIZE must be greater than 0")
	}
	if c.Webhook.MaxAttempts <= 0 {
		return errors.New("NOTIFICATION_WEBHOOK_MAX_ATTEMPTS must be greater than 0")
	}
	if strings.TrimSpace(c.Topics.ConsumerGroupID) == "" {
		return errors.New("NOTIFICATION_CONSUMER_GROUP_ID is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	parts := strings.Split(valueStr, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		values = append(values, trimmed)
	}
	if len(values) == 0 {
		return defaultValue
	}
	return values
}

func loadNotificationServicePorts() (grpcPort int, httpPort int) {
	grpcPort = defaultNotificationGRPCPort
	httpPort = defaultNotificationHTTPPort

	type servicesConfig struct {
		Services map[string]struct {
			Ports struct {
				Grpc int `yaml:"grpc"`
				Http int `yaml:"http"`
			} `yaml:"ports"`
		} `yaml:"services"`
	}

	servicesConfigPath, err := opsconfig.ResolveConfigPath("services.yaml")
	if err != nil {
		logger.Warnf("Failed to resolve services.yaml for notification ports: %v (using defaults)", err)
		return grpcPort, httpPort
	}

	data, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		logger.Warnf("Failed to read services.yaml for notification ports: %v (using defaults)", err)
		return grpcPort, httpPort
	}

	var cfg servicesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Warnf("Failed to parse services.yaml for notification ports: %v (using defaults)", err)
		return grpcPort, httpPort
	}

	serviceCfg, ok := cfg.Services["notification"]
	if !ok {
		logger.Warnf("Notification service not found in services.yaml (using defaults)")
		return grpcPort, httpPort
	}
	if serviceCfg.Ports.Grpc > 0 {
		grpcPort = serviceCfg.Ports.Grpc
	}
	if serviceCfg.Ports.Http > 0 {
		httpPort = serviceCfg.Ports.Http
	}
	return grpcPort, httpPort
}
