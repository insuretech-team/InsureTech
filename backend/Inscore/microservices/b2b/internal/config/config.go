package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	GRPCPort int
	HTTPPort int
	Host     string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: getEnvAsInt("B2B_GRPC_PORT", 50112),
			HTTPPort: getEnvAsInt("B2B_HTTP_PORT", 50113),
			Host:     getEnv("B2B_HOST", "0.0.0.0"),
		},
	}
	return cfg, nil
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
