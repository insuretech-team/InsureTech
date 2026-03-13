package client

import (
	"time"
)

// Config represents the SDK configuration
type Config struct {
	// API Configuration
	BaseURL string
	APIKey  string
	
	// HTTP Client Configuration
	Timeout       time.Duration
	MaxRetries    int
	RetryDelay    time.Duration
	
	// TLS Configuration
	InsecureSkipVerify bool
	
	// Logging
	Debug bool
	
	// Rate Limiting
	RateLimit         int
	RateLimitInterval time.Duration
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		BaseURL:           "https://api.labaidinsuretech.com",
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		RetryDelay:        time.Second,
		InsecureSkipVerify: false,
		Debug:             false,
		RateLimit:         100,
		RateLimitInterval: time.Minute,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return ErrInvalidConfig{Field: "BaseURL", Message: "base URL is required"}
	}
	
	if c.APIKey == "" {
		return ErrInvalidConfig{Field: "APIKey", Message: "API key is required"}
	}
	
	if c.Timeout <= 0 {
		return ErrInvalidConfig{Field: "Timeout", Message: "timeout must be positive"}
	}
	
	if c.MaxRetries < 0 {
		return ErrInvalidConfig{Field: "MaxRetries", Message: "max retries cannot be negative"}
	}
	
	return nil
}

// ErrInvalidConfig represents a configuration validation error
type ErrInvalidConfig struct {
	Field   string
	Message string
}

func (e ErrInvalidConfig) Error() string {
	return "invalid config field '" + e.Field + "': " + e.Message
}

// Environment represents the API environment
type Environment string

const (
	// EnvironmentProduction is the production environment
	EnvironmentProduction Environment = "production"
	
	// EnvironmentStaging is the staging environment
	EnvironmentStaging Environment = "staging"
	
	// EnvironmentDevelopment is the development environment
	EnvironmentDevelopment Environment = "development"
)

// GetBaseURL returns the base URL for the given environment
func (e Environment) GetBaseURL() string {
	switch e {
	case EnvironmentProduction:
		return "https://api.insuretech.com"
	case EnvironmentStaging:
		return "https://api.staging.insuretech.com"
	case EnvironmentDevelopment:
		return "https://api.dev.insuretech.com"
	default:
		return "https://api.insuretech.com"
	}
}

// NewConfigForEnvironment creates a new Config for the specified environment
func NewConfigForEnvironment(env Environment, apiKey string) *Config {
	config := DefaultConfig()
	config.BaseURL = env.GetBaseURL()
	config.APIKey = apiKey
	return config
}
