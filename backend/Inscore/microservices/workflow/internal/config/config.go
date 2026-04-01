package config

import (
	"os"
	"strings"
)

type Config struct {
	GRPCPort        int      // set from services.yaml by main.go
	KafkaBrokers    []string
	AuthzServiceURL string
}

func Load() (*Config, error) {
	brokersStr := getEnv("KAFKA_BROKERS", "localhost:9092")
	var brokers []string
	for _, b := range strings.Split(brokersStr, ",") {
		if b = strings.TrimSpace(b); b != "" {
			brokers = append(brokers, b)
		}
	}

	return &Config{
		GRPCPort:        50180, // overridden by services.yaml in main.go
		KafkaBrokers:    brokers,
		AuthzServiceURL: getEnv("AUTHZ_GRPC_ADDR", getEnv("AUTHZ_SERVICE_URL", "")),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
