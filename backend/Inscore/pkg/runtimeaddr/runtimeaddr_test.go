package runtimeaddr

import "testing"

func TestNormalizeKafkaBrokersDedupes(t *testing.T) {
	got := NormalizeKafkaBrokers([]string{"localhost:9092", "localhost:9092"})
	if len(got) == 0 || got[0] != "localhost:9092" {
		t.Fatalf("unexpected brokers: %v", got)
	}
}

func TestNormalizeRedisURLInvalid(t *testing.T) {
	raw := "not-a-url"
	if got := NormalizeRedisURL(raw); got != raw {
		t.Fatalf("NormalizeRedisURL() = %q", got)
	}
}
