package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	KafkaBrokers                []string
	DefaultPayeeID              string
	PublicBaseURL               string
	HTTPTimeout                 time.Duration
	SSLCommerzStoreID           string
	SSLCommerzStorePassword     string
	SSLCommerzHostedBaseURL     string
	SSLCommerzAPIBaseURL        string
	SSLCommerzValidationBaseURL string
	SSLCommerzRefundBaseURL     string
}

func Load() (*Config, error) {
	return &Config{
		KafkaBrokers:                splitCSV(os.Getenv("KAFKA_BROKERS")),
		DefaultPayeeID:              strings.TrimSpace(os.Getenv("PAYMENT_DEFAULT_PAYEE_ID")),
		PublicBaseURL:               strings.TrimRight(strings.TrimSpace(os.Getenv("PAYMENT_PUBLIC_BASE_URL")), "/"),
		HTTPTimeout:                 durationSecondsEnv("PAYMENT_HTTP_TIMEOUT_SECONDS", 15),
		SSLCommerzStoreID:           strings.TrimSpace(os.Getenv("PAYMENT_SSLCOMMERZ_STORE_ID")),
		SSLCommerzStorePassword:     strings.TrimSpace(os.Getenv("PAYMENT_SSLCOMMERZ_STORE_PASSWORD")),
		SSLCommerzHostedBaseURL:     firstNonEmpty(os.Getenv("PAYMENT_SSLCOMMERZ_HOSTED_BASE_URL"), "https://sandbox.sslcommerz.com"),
		SSLCommerzAPIBaseURL:        firstNonEmpty(os.Getenv("PAYMENT_SSLCOMMERZ_API_BASE_URL"), "https://sandbox.sslcommerz.com"),
		SSLCommerzValidationBaseURL: firstNonEmpty(os.Getenv("PAYMENT_SSLCOMMERZ_VALIDATION_BASE_URL"), "https://sandbox.sslcommerz.com"),
		SSLCommerzRefundBaseURL:     firstNonEmpty(os.Getenv("PAYMENT_SSLCOMMERZ_REFUND_BASE_URL"), "https://sandbox.sslcommerz.com"),
	}, nil
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func durationSecondsEnv(key string, fallback int) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return time.Duration(fallback) * time.Second
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return time.Duration(fallback) * time.Second
	}
	return time.Duration(seconds) * time.Second
}
