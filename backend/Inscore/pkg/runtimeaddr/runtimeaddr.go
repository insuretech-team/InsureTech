package runtimeaddr

import (
	"net"
	"net/url"
	"strings"

	"github.com/newage-saint/insuretech/ops/env"
)

// NormalizeKafkaBrokers expands runtime-friendly Kafka broker fallbacks for host, Docker, and WSL.
func NormalizeKafkaBrokers(brokers []string) []string {
	if len(brokers) == 0 {
		return brokers
	}

	inDocker := env.IsRunningInDocker()
	inWSL := env.IsRunningInWSL()
	normalized := make([]string, 0, len(brokers)*3)

	for _, broker := range brokers {
		broker = strings.TrimSpace(broker)
		if broker == "" {
			continue
		}

		host, port, err := net.SplitHostPort(broker)
		if err != nil {
			switch broker {
			case "localhost", "127.0.0.1":
				if inDocker {
					// Inside Docker: use kafka service name only — never add host.docker.internal
				normalized = append(normalized, "kafka:9092")
				} else if inWSL {
					normalized = append(normalized, "host.docker.internal:9092", "localhost:9092")
				} else {
					normalized = append(normalized, "localhost:9092", "host.docker.internal:9092")
				}
			case "kafka":
				if inDocker {
					// Already the correct Docker service name — no fallbacks needed
				normalized = append(normalized, "kafka:9092")
				} else {
					normalized = append(normalized, "host.docker.internal:9092", "localhost:9092")
				}
			default:
				normalized = append(normalized, broker)
			}
			continue
		}

		switch host {
		case "localhost", "127.0.0.1", "::1":
			if inDocker {
				normalized = append(normalized, net.JoinHostPort("kafka", port), net.JoinHostPort("host.docker.internal", port), broker)
			} else if inWSL {
				normalized = append(normalized, net.JoinHostPort("host.docker.internal", port), broker)
			} else {
				normalized = append(normalized, broker, net.JoinHostPort("host.docker.internal", port))
			}
		case "kafka":
			if inDocker {
				// Inside Docker: kafka service name is directly reachable via Docker network.
				// Do NOT add host.docker.internal — on Linux Docker it resolves to the host
				// IP which may map to [::1]:9092 causing consumer connection refused errors.
				normalized = append(normalized, broker)
			} else {
				normalized = append(normalized, net.JoinHostPort("host.docker.internal", port), net.JoinHostPort("localhost", port))
			}
		case "host.docker.internal":
			if inDocker {
				normalized = append(normalized, broker, net.JoinHostPort("kafka", port))
			} else {
				normalized = append(normalized, broker, net.JoinHostPort("localhost", port))
			}
		default:
			normalized = append(normalized, broker)
		}
	}

	return dedupeStrings(normalized)
}

// NormalizeRedisURL rewrites loopback Redis hosts for Docker and WSL execution.
func NormalizeRedisURL(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	host := parsed.Hostname()
	port := parsed.Port()
	if port == "" {
		port = "6379"
	}

	switch host {
	case "localhost", "127.0.0.1", "::1":
		if env.IsRunningInDocker() {
			parsed.Host = net.JoinHostPort("redis", port)
		} else if env.IsRunningInWSL() {
			parsed.Host = net.JoinHostPort("host.docker.internal", port)
		}
	}
	return parsed.String()
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
