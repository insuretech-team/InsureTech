package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("PARTNER_GRPC_PORT", "50100")
	t.Setenv("PARTNER_HTTP_PORT", "50101")
	t.Setenv("PARTNER_HOST", "127.0.0.1")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("KAFKA_BROKERS", "broker-1:9092, broker-2:9092")
	t.Setenv("AUTHN_GRPC_ADDR", "authn:50051")
	t.Setenv("AUTHZ_SERVICE_ADDRESS", "authz:50052")
	t.Setenv("PII_ENCRYPTION_KEY", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, 50100, cfg.Server.GRPCPort)
	assert.Equal(t, 50101, cfg.Server.HTTPPort)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, []string{"broker-1:9092", "broker-2:9092"}, cfg.Kafka.Brokers)
	assert.Equal(t, "authn:50051", cfg.Integration.AuthNAddress)
	assert.Equal(t, "authz:50052", cfg.Integration.AuthZAddress)
	assert.Equal(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", cfg.Security.PIIEncryptionKey)
}

func TestConfigHelpers(t *testing.T) {
	t.Setenv("INT_OK", "42")
	t.Setenv("INT_BAD", "nope")
	t.Setenv("BOOL_OK", "true")
	t.Setenv("BOOL_BAD", "nope")
	t.Setenv("DURATION_OK", "30s")
	t.Setenv("DURATION_BAD", "later")
	t.Setenv("SLICE_OK", " a, b ,, c ")
	t.Setenv("SLICE_EMPTY", " , , ")

	assert.Equal(t, "fallback", getEnv("MISSING_ENV", "fallback"))
	assert.Equal(t, 42, getEnvAsInt("INT_OK", 5))
	assert.Equal(t, 5, getEnvAsInt("INT_BAD", 5))
	assert.True(t, getEnvAsBool("BOOL_OK", false))
	assert.True(t, getEnvAsBool("BOOL_BAD", true))
	assert.Equal(t, 30*time.Second, getEnvAsDuration("DURATION_OK", time.Minute))
	assert.Equal(t, time.Minute, getEnvAsDuration("DURATION_BAD", time.Minute))
	assert.Equal(t, []string{"a", "b", "c"}, getEnvAsSlice("SLICE_OK", []string{"x"}))
	assert.Equal(t, []string{"x"}, getEnvAsSlice("SLICE_EMPTY", []string{"x"}))
}
