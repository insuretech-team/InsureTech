package config

import (
	"os"
	"strconv"
)

type Config struct {
	GRPCPort          int
	DBHost            string
	DBPort            int
	DBUser            string
	DBPassword        string
	DBName            string
	KafkaBrokers      []string
	PaymentServiceURL string
	AuthzServiceURL   string
}

func Load() (*Config, error) {
	grpcPort, _ := strconv.Atoi(getEnv("GRPC_PORT", "50142"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return &Config{
		GRPCPort:          grpcPort,
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            dbPort,
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "insuretech"),
		KafkaBrokers:      []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "localhost:50190"),
		AuthzServiceURL:   getEnv("AUTHZ_GRPC_ADDR", getEnv("AUTHZ_SERVICE_URL", "")),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
