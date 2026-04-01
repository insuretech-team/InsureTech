package env

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Load loads the project .env file (searched upward from CWD), then .env.local
// if present, then normalizes variables.  After loading, it auto-detects whether
// the process is running inside Docker or directly on the host and rewrites any
// Docker Compose service hostnames (e.g. "authn", "kafka") to "localhost" so
// that services can bind and connect without a Docker DNS resolver.
func Load() error {
	paths := candidateEnvPaths()
	var lastErr error
	var loaded bool

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			// Use Overload (not Load) so .env values always take precedence over
			// any stale OS environment variables inherited from the parent shell.
			// This prevents issues where KYC_SERVICE_ENABLED=false or FLVE_HF_TOKEN=""
			// in the parent shell silently overrides the correct .env values.
			if err := godotenv.Overload(p); err == nil {
				loaded = true
				// .env.local can still override individual values if needed.
				localPath := filepath.Join(filepath.Dir(p), ".env.local")
				if _, err := os.Stat(localPath); err == nil {
					_ = godotenv.Overload(localPath)
				}
				normalize()
				return nil
			} else {
				lastErr = err
			}
		}
	}

	if err := godotenv.Overload(); err == nil {
		loaded = true
		if _, err := os.Stat(".env.local"); err == nil {
			_ = godotenv.Overload(".env.local")
		}
		normalize()
		return nil
	}

	normalize()

	if lastErr != nil && !loaded {
		return lastErr
	}
	if !loaded {
		return errors.New(".env not found")
	}
	return nil
}

// IsRunningInDocker reports whether the current process is running inside a
// Docker container (checks /.dockerenv and /proc/1/cgroup).
func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	data, err := os.ReadFile("/proc/1/cgroup")
	if err == nil {
		c := string(data)
		if strings.Contains(c, "docker") || strings.Contains(c, "containerd") {
			return true
		}
	}
	return false
}

// IsRunningInWSL reports whether the current process is running inside the
// Windows Subsystem for Linux.
func IsRunningInWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSL_INTEROP") != "" {
		return true
	}
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

// isDockerServiceHostname returns true for short, unqualified names that are
// Docker Compose service names.  Such names have no dots, are not IP addresses,
// and are not the standard loopback / wildcard identifiers.
func isDockerServiceHostname(host string) bool {
	h := strings.TrimSpace(strings.ToLower(host))
	switch h {
	case "", "localhost", "0.0.0.0", "127.0.0.1", "::1":
		return false
	}
	if strings.Contains(h, ".") {
		return false // FQDN or dotted IP
	}
	if net.ParseIP(h) != nil {
		return false // bare IPv6 or other IP literal
	}
	return true // single-label name → Docker Compose service
}

// normalize applies post-load fixups and, when running outside Docker,
// rewrites Docker Compose service hostnames to localhost equivalents.
func normalize() {
	// Strip accidental wrapping quotes from loaded .env values so services see
	// the actual host/URL/token rather than a literal quoted string.
	for _, kv := range os.Environ() {
		idx := strings.IndexByte(kv, '=')
		if idx < 0 {
			continue
		}
		key, val := kv[:idx], kv[idx+1:]
		if cleaned := trimWrappingQuotes(val); cleaned != val {
			_ = os.Setenv(key, cleaned)
		}
	}

	// PGPORT fallback from legacy alias.
	if os.Getenv("PGPORT") == "" {
		if v := os.Getenv("NEON_DB_PORT"); v != "" {
			_ = os.Setenv("PGPORT", v)
		}
	}

	if IsRunningInDocker() {
		// Inside Docker: the Docker DNS resolver handles service names. No rewrite needed.
		return
	}

	// ── Running on host (Windows / macOS / Linux bare-metal / WSL) ───────────
	// Rewrite Docker Compose service hostnames so services bind on 0.0.0.0 and
	// connect to each other via localhost instead of unresolvable service names.

	// Pass 1 — *_HOST bind vars  (e.g. AUTHN_HOST=authn → 0.0.0.0)
	for _, kv := range os.Environ() {
		idx := strings.IndexByte(kv, '=')
		if idx < 0 {
			continue
		}
		key, val := kv[:idx], kv[idx+1:]
		if strings.HasSuffix(key, "_HOST") && isDockerServiceHostname(val) {
			_ = os.Setenv(key, "0.0.0.0")
		}
	}

	// Pass 2 — *_GRPC_ADDR client endpoints  (e.g. AUTHN_GRPC_ADDR=authn:50060 → localhost:50060)
	for _, kv := range os.Environ() {
		idx := strings.IndexByte(kv, '=')
		if idx < 0 {
			continue
		}
		key, val := kv[:idx], kv[idx+1:]
		if strings.HasSuffix(key, "_GRPC_ADDR") {
			if rw := rewriteHostPort(val, "localhost"); rw != val {
				_ = os.Setenv(key, rw)
			}
		}
	}

	// Pass 3 — KAFKA_BROKERS comma-separated list  (e.g. kafka:9092 → localhost:9092)
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		parts := strings.Split(v, ",")
		changed := false
		for i, p := range parts {
			trimmed := strings.TrimSpace(p)
			if rw := rewriteHostPort(trimmed, "localhost"); rw != trimmed {
				parts[i] = rw
				changed = true
			}
		}
		if changed {
			_ = os.Setenv("KAFKA_BROKERS", strings.Join(parts, ","))
		}
	}

	// Pass 4 — REDIS_URL  (e.g. redis://redis:6379 → redis://localhost:6379)
	if v := os.Getenv("REDIS_URL"); v != "" {
		if rw := rewriteURLHostname(v, "localhost"); rw != v {
			_ = os.Setenv("REDIS_URL", rw)
		}
	}

	// Pass 5 — Misc URL / address vars
	for _, key := range []string{"GOTENBERG_URL", "NEWMAN_BASE_URL", "STORAGE_SERVICE_ADDRESS"} {
		v := os.Getenv(key)
		if v == "" {
			continue
		}
		rw := rewriteURLHostname(v, "localhost")
		if rw == v && isDockerServiceHostname(v) {
			rw = "localhost" // bare hostname (no scheme/port)
		}
		if rw != v {
			_ = os.Setenv(key, rw)
		}
	}
}

// rewriteHostPort replaces the host in a "host:port" string with replacement
// when the host is a Docker service name, otherwise returns addr unchanged.
func rewriteHostPort(addr, replacement string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	if isDockerServiceHostname(host) {
		return net.JoinHostPort(replacement, port)
	}
	return addr
}

// rewriteURLHostname replaces the hostname inside a URL string when it is a
// Docker service name.  Port (if present) is preserved.
// e.g. "http://redis:6379/0" → "http://localhost:6379/0"
func rewriteURLHostname(rawURL, replacement string) string {
	schemeEnd := strings.Index(rawURL, "://")
	if schemeEnd < 0 {
		return rawURL
	}
	rest := rawURL[schemeEnd+3:]
	slashIdx := strings.IndexAny(rest, "/")
	hostPart, pathPart := rest, ""
	if slashIdx >= 0 {
		hostPart = rest[:slashIdx]
		pathPart = rest[slashIdx:]
	}
	host, port, err := net.SplitHostPort(hostPart)
	if err != nil {
		// No port — hostPart is just the hostname.
		if isDockerServiceHostname(hostPart) {
			return rawURL[:schemeEnd+3] + replacement + pathPart
		}
		return rawURL
	}
	if isDockerServiceHostname(host) {
		return rawURL[:schemeEnd+3] + net.JoinHostPort(replacement, port) + pathPart
	}
	return rawURL
}

func trimWrappingQuotes(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			return strings.TrimSpace(value[1 : len(value)-1])
		}
	}
	return value
}

// candidateEnvPaths returns candidate .env locations, walking up from CWD.
func candidateEnvPaths() []string {
	var out []string
	wd, _ := os.Getwd()
	cur := wd
	for i := 0; i < 10 && cur != "" && cur != string(filepath.Separator); i++ {
		out = append(out, filepath.Join(cur, ".env"))
		cur = filepath.Dir(cur)
	}
	return out
}

// GetEnv retrieves environment variable with default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt retrieves environment variable as integer with default value
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsBool retrieves environment variable as boolean with default value
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsSlice parses comma-separated environment variable into slice
func GetEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		result := []string{}
		for _, item := range splitByComma(value) {
			if trimmed := trim(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

// splitByComma splits string by comma
func splitByComma(s string) []string {
	var result []string
	var current string
	for _, char := range s {
		if char == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// trim removes leading and trailing whitespace
func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}
