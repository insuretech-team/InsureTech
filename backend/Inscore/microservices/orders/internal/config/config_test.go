package config

import "testing"

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("GRPC_PORT", "50221")
	t.Setenv("DB_HOST", "db.internal")
	t.Setenv("DB_PORT", "5439")
	t.Setenv("DB_USER", "orders")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "ordersdb")
	t.Setenv("KAFKA_BROKERS", "broker-1:9092")
	t.Setenv("PAYMENT_SERVICE_URL", "payment:50190")
	t.Setenv("AUTHZ_GRPC_ADDR", "authz:50052")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.GRPCPort != 50221 || cfg.DBPort != 5439 {
		t.Fatalf("unexpected numeric config: %+v", cfg)
	}
	if cfg.DBHost != "db.internal" || cfg.DBUser != "orders" || cfg.DBName != "ordersdb" {
		t.Fatalf("unexpected db config: %+v", cfg)
	}
	if len(cfg.KafkaBrokers) != 1 || cfg.KafkaBrokers[0] != "broker-1:9092" {
		t.Fatalf("unexpected brokers: %+v", cfg.KafkaBrokers)
	}
	if cfg.PaymentServiceURL != "payment:50190" || cfg.AuthzServiceURL != "authz:50052" {
		t.Fatalf("unexpected service urls: %+v", cfg)
	}
}

func TestGetEnvFallback(t *testing.T) {
	if got := getEnv("ORDERS_TEST_MISSING", "fallback"); got != "fallback" {
		t.Fatalf("getEnv() = %q", got)
	}
}
