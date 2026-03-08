package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// Config contains fraud microservice runtime settings.
type Config struct {
	Server      ServerConfig
	Kafka       KafkaConfig
	Integration IntegrationConfig
}

type ServerConfig struct {
	Host string
}

type KafkaConfig struct {
	Brokers        []string
	Topic          string
	ConsumerTopics []string
	ConsumerGroup  string
	DLQTopic       string
}

type IntegrationConfig struct {
	AuthZAddress string
}

// Load reads fraud configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("FRAUD_HOST", "0.0.0.0"),
		},
		Kafka: KafkaConfig{
			Brokers:        getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:          getEnv("KAFKA_FRAUD_TOPIC", "fraud.events"),
			ConsumerTopics: getEnvAsSlice("KAFKA_FRAUD_CONSUMER_TOPICS", []string{"policy.issued", "policy.renewed", "claim.created", "claim.submitted", "customer.risk.updated"}),
			ConsumerGroup:  getEnv("KAFKA_FRAUD_CONSUMER_GROUP", "fraud-service-consumer"),
			DLQTopic:       getEnv("KAFKA_FRAUD_DLQ_TOPIC", "fraud.dlq"),
		},
		Integration: IntegrationConfig{
			AuthZAddress: getEnv("AUTHZ_GRPC_ADDR", getEnv("AUTHZ_SERVICE_ADDRESS", "")),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate performs basic config validation.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Server.Host) == "" {
		return errors.New("FRAUD_HOST cannot be empty")
	}
	if len(c.Kafka.Brokers) == 0 {
		return errors.New("KAFKA_BROKERS cannot be empty")
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
