package config

import (
	"errors"
	"fmt"
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
	defaultDocgenHTTPPort    = 50281
	defaultGotenbergHTTPPort = 50330
)

// Config holds all configuration for the docgen microservice
type Config struct {
	Port                 int
	GotenbergURL         string
	StorageServiceAddr   string
	KafkaBrokers         []string
	KafkaDocgenTopic     string
	AsyncGeneration      bool
	AsyncWorkerCount     int
	MaxGenerationTimeout time.Duration
}

// Load loads configuration from environment variables and services.yaml
func Load() (*Config, error) {
	// Load .env file first
	if err := env.Load(); err != nil {
		logger.Warnf("Failed to load .env file: %v (using system environment variables)", err)
	}

	docgenHTTPPort := loadDocgenServicePort()
	gotenbergURLDefault := fmt.Sprintf("http://localhost:%d", loadServiceHTTPPort("gotenberg", defaultGotenbergHTTPPort))

	cfg := &Config{
		Port:                 docgenHTTPPort,
		GotenbergURL:         getEnv("GOTENBERG_URL", gotenbergURLDefault),
		StorageServiceAddr:   getEnv("STORAGE_SERVICE_ADDR", ""),
		KafkaBrokers:         getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		KafkaDocgenTopic:     getEnv("KAFKA_DOCGEN_TOPIC", "docgen-events"),
		AsyncGeneration:      getEnvAsBool("DOCGEN_ASYNC_GENERATION", false),
		AsyncWorkerCount:     getEnvAsInt("DOCGEN_WORKER_COUNT", 3),
		MaxGenerationTimeout: getEnvAsDuration("DOCGEN_TIMEOUT", 30*time.Second),
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
	if c.Port <= 0 {
		logger.Errorf("Invalid port: %d", c.Port)
		return errors.New("invalid port")
	}
	if c.GotenbergURL == "" {
		logger.Errorf("GOTENBERG_URL is required")
		return errors.New("GOTENBERG_URL is required")
	}
	if len(c.KafkaBrokers) == 0 {
		logger.Errorf("KAFKA_BROKERS is required")
		return errors.New("KAFKA_BROKERS is required")
	}
	if c.AsyncWorkerCount <= 0 {
		logger.Errorf("DOCGEN_WORKER_COUNT must be greater than 0")
		return errors.New("DOCGEN_WORKER_COUNT must be greater than 0")
	}
	if c.MaxGenerationTimeout <= 0 {
		logger.Errorf("DOCGEN_TIMEOUT must be greater than 0")
		return errors.New("DOCGEN_TIMEOUT must be greater than 0")
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
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	if len(result) == 0 {
		return defaultValue
	}
	return result
}

func loadDocgenServicePort() int {
	return loadServiceHTTPPort("docgen", defaultDocgenHTTPPort)
}

func loadServiceHTTPPort(serviceName string, defaultPort int) int {
	port := defaultPort

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
		logger.Warnf("Failed to resolve services.yaml for docgen port: %v (using default %d)", err, port)
		return port
	}

	data, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		logger.Warnf("Failed to read services.yaml for docgen port: %v (using default %d)", err, port)
		return port
	}

	var cfg servicesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Warnf("Failed to parse services.yaml for %s port: %v (using default %d)", serviceName, err, port)
		return port
	}

	svc, ok := cfg.Services[serviceName]
	if !ok {
		logger.Warnf("%s service not found in services.yaml (using default %d)", serviceName, port)
		return port
	}
	if svc.Ports.Http > 0 {
		port = svc.Ports.Http
	}
	return port
}
