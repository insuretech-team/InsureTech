package env

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Load attempts to load the project's .env file from the repository root.
// It searches upward from the current working directory for a .env file.
// After loading .env, it also loads .env.local if it exists (overriding values).
// After loading, it also normalizes common variables (e.g., PGPORT from NEON_DB_PORT).
func Load() error {
	paths := candidateEnvPaths()
	var lastErr error
	var loaded bool
	
	// First, load .env file
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			// Use Overload so .env always wins over any inherited env vars
			// (e.g. stale PGHOST set by a previous shell session).
			if err := godotenv.Overload(p); err == nil {
				loaded = true
				// Now try to load .env.local from the same directory (overrides .env)
				localPath := filepath.Join(filepath.Dir(p), ".env.local")
				if _, err := os.Stat(localPath); err == nil {
					_ = godotenv.Overload(localPath) // Overload to override existing values
				}
				normalize()
				return nil
			} else {
				lastErr = err
			}
		}
	}
	
	// Fallback: try default godotenv.Overload in CWD
	if err := godotenv.Overload(); err == nil {
		loaded = true
		// Try .env.local in CWD
		if _, err := os.Stat(".env.local"); err == nil {
			_ = godotenv.Overload(".env.local")
		}
		normalize()
		return nil
	}

	// If no .env file found, just normalize existing env vars and return success
	// This allows using system environment variables
	normalize()

	if lastErr != nil && !loaded {
		return lastErr
	}
	if !loaded {
		return errors.New(".env not found")
	}
	return nil
}

// candidateEnvPaths returns a list of potential .env locations walking up the tree.
func candidateEnvPaths() []string {
	var out []string
	// Start from CWD and go up more levels to handle deep test directories
	wd, _ := os.Getwd()
	cur := wd
	for i := 0; i < 10 && cur != "" && cur != string(filepath.Separator); i++ {
		out = append(out, filepath.Join(cur, ".env"))
		cur = filepath.Dir(cur)
	}
	return out
}

// normalize maps or sets env variables required by the app if not present.
func normalize() {
	// If PGPORT is not set but NEON_DB_PORT is, set PGPORT accordingly
	if os.Getenv("PGPORT") == "" {
		if v := os.Getenv("NEON_DB_PORT"); v != "" {
			_ = os.Setenv("PGPORT", v)
		}
	}
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
