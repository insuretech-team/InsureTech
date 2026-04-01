package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("FRAUD_HOST", "127.0.0.1")
	t.Setenv("KAFKA_BROKERS", "broker-1:9092, broker-2:9092")
	t.Setenv("KAFKA_FRAUD_TOPIC", "fraud.topic")
	t.Setenv("KAFKA_FRAUD_CONSUMER_TOPICS", "claim.submitted, policy.issued")
	t.Setenv("KAFKA_FRAUD_CONSUMER_GROUP", "fraud-group")
	t.Setenv("KAFKA_FRAUD_DLQ_TOPIC", "fraud.dlq.custom")
	t.Setenv("AUTHZ_GRPC_ADDR", "authz:50052")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", cfg.Server.Host)
	require.Equal(t, []string{"broker-1:9092", "broker-2:9092"}, cfg.Kafka.Brokers)
	require.Equal(t, "fraud.topic", cfg.Kafka.Topic)
	require.Equal(t, []string{"claim.submitted", "policy.issued"}, cfg.Kafka.ConsumerTopics)
	require.Equal(t, "fraud-group", cfg.Kafka.ConsumerGroup)
	require.Equal(t, "fraud.dlq.custom", cfg.Kafka.DLQTopic)
	require.Equal(t, "authz:50052", cfg.Integration.AuthZAddress)
}

func TestValidateAndHelpers(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Host: ""},
		Kafka:  KafkaConfig{Brokers: []string{"broker:9092"}},
	}
	require.Error(t, cfg.Validate())

	cfg.Server.Host = "0.0.0.0"
	cfg.Kafka.Brokers = nil
	require.Error(t, cfg.Validate())

	t.Setenv("INT_OK", "42")
	t.Setenv("INT_BAD", "bad")
	t.Setenv("SLICE_OK", " a, b ,, c ")
	t.Setenv("SLICE_EMPTY", " , ")
	require.Equal(t, "fallback", getEnv("MISSING_VALUE", "fallback"))
	require.Equal(t, 42, getEnvAsInt("INT_OK", 7))
	require.Equal(t, 7, getEnvAsInt("INT_BAD", 7))
	require.Equal(t, []string{"a", "b", "c"}, getEnvAsSlice("SLICE_OK", []string{"x"}))
	require.Equal(t, []string{"x"}, getEnvAsSlice("SLICE_EMPTY", []string{"x"}))
}
