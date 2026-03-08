package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// Config holds all configuration for the partner microservice
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Kafka       KafkaConfig
	Integration IntegrationConfig
	Security    SecurityConfig
}

// ServerConfig contains server settings
type ServerConfig struct {
	GRPCPort int
	HTTPPort int
	Host     string
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// RedisConfig contains Redis connection settings
type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

// KafkaConfig contains Kafka settings
type KafkaConfig struct {
	Brokers        []string
	Topic          string
	ConsumerTopics []string
	ConsumerGroup  string
	DLQTopic       string
}

// IntegrationConfig holds endpoints or keys for communicating with other services (like AuthN)
type IntegrationConfig struct {
	AuthNAddress string // For calling GetPartnerAPICredentials
	AuthZAddress string // For Casbin permission checks
}

// SecurityConfig contains security-related settings like encryption keys
type SecurityConfig struct {
	PIIEncryptionKey string // 32-byte key for AES-GCM
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: getEnvAsInt("PARTNER_GRPC_PORT", 50058),
			HTTPPort: getEnvAsInt("PARTNER_HTTP_PORT", 8088),
			Host:     getEnv("PARTNER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "insuretech_primary"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "insuretech_primary"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Kafka: KafkaConfig{
			Brokers:        getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:          getEnv("KAFKA_PARTNER_TOPIC", "partner-events"),
			ConsumerTopics: getEnvAsSlice("KAFKA_PARTNER_CONSUMER_TOPICS", []string{"policy.issued", "policy.renewed", "policy.events", "renewal.policy.renewed"}),
			ConsumerGroup:  getEnv("KAFKA_PARTNER_CONSUMER_GROUP", "partner-service-consumer"),
			DLQTopic:       getEnv("KAFKA_PARTNER_DLQ_TOPIC", "partner.dlq"),
		},
		Integration: IntegrationConfig{
			AuthNAddress: getEnv("AUTHN_GRPC_ADDR", getEnv("AUTHN_SERVICE_ADDRESS", "")),
			AuthZAddress: getEnv("AUTHZ_GRPC_ADDR", getEnv("AUTHZ_SERVICE_ADDRESS", "")),
		},
		Security: SecurityConfig{
			PIIEncryptionKey: getEnv("PII_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef"), // Default 32-byte key for local dev
		},
	}

	// Validation
	if err := cfg.Validate(); err != nil {
		logger.Errorf("config validation failed: %v", err)
		return nil, errors.New("config validation failed")
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Password == "" {
		// Log warning instead of failing strict locally for now, typically want to enforce it
		logger.Warnf("DB_PASSWORD is empty")
	}
	return nil
}

// Helper functions for environment variable parsing

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
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return defaultValue
	}
	return out
}
