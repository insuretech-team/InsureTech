package config

import (
	"os"
	"strconv"
)

type Config struct {
	GRPCPort        int
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	KafkaBrokers    []string
	AuthZServiceURL string
}

func Load() (*Config, error) {
	grpcPort, _ := strconv.Atoi(getEnv("GRPC_PORT", "50112"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return &Config{
		GRPCPort:        grpcPort,
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          dbPort,
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "insuretech"),
		KafkaBrokers:    []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		AuthZServiceURL: getEnv("AUTHZ_GRPC_ADDR", getEnv("AUTHZ_SERVICE_URL", "")),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
