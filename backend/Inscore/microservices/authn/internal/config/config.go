package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	opsconfig "github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"gopkg.in/yaml.v3"
)

const (
	defaultAuthnGRPCPort = 50060
	defaultAuthnHTTPPort = 50061
)

// Config holds all configuration for the authn microservice
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	KYC      KYCConfig
	SMS      SMSConfig
	Email    EmailConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	Security SecurityConfig
	FLVE     FLVEConfig
}

// FLVEConfig contains Face Liveness & Verification Engine settings
type FLVEConfig struct {
	Port              int
	Backend           string
	HFEndpoint        string
	HFToken           string
	MainCDN           string
	CDNURL            string
	AccessKey         string
	SecretKey         string
	ModelsPath        string
	ModelDetector     string
	ModelEmbedding    string
	ModelLiveness     string
	UseGPU            bool
	LivenessThreshold float64
	MatchThreshold    float64
	SessionType       string
	RedisURL          string
	SessionTTL        int
	JWTSecret         string
}

// ServerConfig contains server settings
type ServerConfig struct {
	GRPCPort int
	HTTPPort int
	Host     string
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// JWTConfig contains JWT token settings.
// RS256 is mandatory — HS256 is removed.
// Private key is loaded from file (path) or Vault (ref) at startup.
// Public key PEM is served via JWKS endpoint (/.well-known/jwks.json).
type JWTConfig struct {
	// RS256 signing — private key PEM file path (loaded at startup, never stored in DB)
	PrivateKeyPath string // JWT_PRIVATE_KEY_PATH  e.g. /secrets/jwt_rsa_private.pem
	PublicKeyPath  string // JWT_PUBLIC_KEY_PATH   e.g. /secrets/jwt_rsa_public.pem
	KeyID          string // JWT_KEY_ID (kid header) e.g. "insuretech-2025-01"
	// Token expiry — configurable per environment; defaults: access=15m, refresh=7d
	AccessTokenDuration  time.Duration // JWT_ACCESS_TOKEN_DURATION  (default 15m)
	RefreshTokenDuration time.Duration // JWT_REFRESH_TOKEN_DURATION (default 7d)
	Issuer               string        // JWT_ISSUER (default "insuretech-authn")
	Audience             string        // JWT_AUDIENCE (default "insuretech-api")
}

// KYCConfig contains downstream KYC service integration settings.
type KYCConfig struct {
	Enabled bool
	Address string
	Timeout time.Duration
	Token   string
}

// SMSConfig contains SMS provider settings (SSL Wireless)
type SMSConfig struct {
	Provider string
	APIBase  string
	SID      string
	APIKey   string
	// BTRC Masking
	MaskingEnabled  bool
	MaskingSenderID string
	// Non-masking fallback
	NonMaskingEnabled bool
	NonMaskingSender  string
	// DLR webhook
	DLRWebhookURL string
}

// EmailConfig contains email settings
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	From     string
	Username string
	Password string
	TLS      bool
}

// RedisConfig contains Redis connection settings
type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

// KafkaConfig contains Kafka settings
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// SecurityConfig contains security-related settings
type SecurityConfig struct {
	ServerSessionDuration time.Duration // Web portal server-side session duration (default 12h)
	OTPLength             int
	OTPExpiry             time.Duration
	OTPMaxAttempts        int
	OTPCooldown           time.Duration
	BCryptCost            int
	RateLimitPerMinute    int
	RateLimitPerDay       int
	IdleTimeoutDuration   time.Duration // Idle session timeout (default 0 = disabled)
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file first
	if err := env.Load(); err != nil {
		logger.Warnf("Failed to load .env file: %v (using system environment variables)", err)
	}

	// Parse flve.yaml
	var flveYAML struct {
		FLVE struct {
			Port              int     `yaml:"port"`
			Backend           string  `yaml:"backend"`
			HFEndpoint        string  `yaml:"hf_endpoint"`
			HFToken           string  `yaml:"hf_token"`
			MainCDN           string  `yaml:"main_cdn"`
			CDNURL            string  `yaml:"cdn_url"`
			AccessKey         string  `yaml:"access_key"`
			SecretKey         string  `yaml:"secret_key"`
			ModelsPath        string  `yaml:"models_path"`
			ModelDetector     string  `yaml:"model_detector"`
			ModelEmbedding    string  `yaml:"model_embedding"`
			ModelLiveness     string  `yaml:"model_liveness"`
			UseGPU            bool    `yaml:"use_gpu"`
			LivenessThreshold float64 `yaml:"liveness_threshold"`
			MatchThreshold    float64 `yaml:"match_threshold"`
			SessionType       string  `yaml:"session_type"`
			RedisURL          string  `yaml:"redis_url"`
			SessionTTL        int     `yaml:"session_ttl"`
			JWTSecret         string  `yaml:"jwt_secret"`
		} `yaml:"flve"`
	}

	flveConfigPath, err := opsconfig.ResolveConfigPath("flve.yaml")
	if err == nil {
		data, readErr := os.ReadFile(flveConfigPath)
		if readErr == nil {
			_ = yaml.Unmarshal(data, &flveYAML)
		} else {
			logger.Warnf("Failed to read flve.yaml: %v", readErr)
		}
	} else {
		logger.Warnf("Failed to resolve flve.yaml: %v", err)
	}

	authnGRPCPort, authnHTTPPort := loadAuthnServicePorts()
	if os.Getenv("AUTHN_PORT") != "" || os.Getenv("AUTHN_GRPC_PORT") != "" || os.Getenv("AUTHN_HTTP_PORT") != "" {
		logger.Warn("AUTHN_PORT/AUTHN_GRPC_PORT/AUTHN_HTTP_PORT are ignored; authn ports are loaded from backend/inscore/configs/services.yaml")
	}

	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: authnGRPCPort,
			HTTPPort: authnHTTPPort,
			Host:     getEnv("AUTHN_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "insuretech_primary"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "insuretech_primary"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			PrivateKeyPath:       getEnv("JWT_PRIVATE_KEY_PATH", "/secrets/jwt_rsa_private.pem"),
			PublicKeyPath:        getEnv("JWT_PUBLIC_KEY_PATH", "/secrets/jwt_rsa_public.pem"),
			KeyID:                getEnv("JWT_KEY_ID", "insuretech-2025-01"),
			AccessTokenDuration:  getEnvAsDuration("JWT_ACCESS_TOKEN_DURATION", 15*time.Minute),
			RefreshTokenDuration: getEnvAsDuration("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour),
			Issuer:               getEnv("JWT_ISSUER", "insuretech-authn"),
			Audience:             getEnv("JWT_AUDIENCE", "insuretech-api"),
		},
		KYC: KYCConfig{
			Enabled: getEnvAsBool("KYC_SERVICE_ENABLED", false),
			Address: getEnv("KYC_SERVICE_ADDRESS", ""),
			Timeout: getEnvAsDuration("KYC_SERVICE_TIMEOUT", 5*time.Second),
			Token:   getEnv("KYC_FLVE_TOKEN", ""),
		},
		SMS: SMSConfig{
			Provider:          getEnv("SMS_PROVIDER", "sslwireless"),
			APIBase:           getEnv("SSLWIRELESS_API_BASE", ""),
			SID:               getEnv("SSLWIRELESS_SID", ""),
			APIKey:            getEnv("SSLWIRELESS_API_KEY", ""),
			MaskingEnabled:    getEnvAsBool("SSLWIRELESS_MASKING_ENABLED", true),
			MaskingSenderID:   getEnv("SSLWIRELESS_SENDER_ID", "LABAIDINS"),
			NonMaskingEnabled: getEnvAsBool("SSLWIRELESS_NONMASKING_ENABLED", true),
			NonMaskingSender:  getEnv("SSLWIRELESS_NONMASKING_SENDER", ""),
			DLRWebhookURL:     getEnv("SSLWIRELESS_DLR_WEBHOOK_URL", ""),
		},
		Email: EmailConfig{
			SMTPHost: getEnv("EMAIL_SMTP_HOST", "smtp.labaidinsuretech.com"),
			SMTPPort: getEnvAsInt("EMAIL_SMTP_PORT", 587),
			From:     getEnv("EMAIL_FROM", "verification@labaidinsuretech.com"),
			Username: getEnv("EMAIL_USERNAME", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
			TLS:      getEnvAsBool("EMAIL_TLS", true),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Kafka: KafkaConfig{
			Brokers: getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:   getEnv("KAFKA_AUTHN_TOPIC", "authn-events"),
		},
		Security: SecurityConfig{
			ServerSessionDuration: getEnvAsDuration("SERVER_SESSION_DURATION", 12*time.Hour),
			OTPLength:             getEnvAsInt("OTP_LENGTH", 6),
			OTPExpiry:             getEnvAsDuration("OTP_EXPIRY", 5*time.Minute),
			OTPMaxAttempts:        getEnvAsInt("OTP_MAX_ATTEMPTS", 3),
			OTPCooldown:           getEnvAsDuration("OTP_COOLDOWN", 60*time.Second),
			BCryptCost:            getEnvAsInt("BCRYPT_COST", 10),
			RateLimitPerMinute:    getEnvAsInt("RATE_LIMIT_PER_MINUTE", 3),
			RateLimitPerDay:       getEnvAsInt("RATE_LIMIT_PER_DAY", 10),
			IdleTimeoutDuration:   time.Duration(getEnvAsInt("IDLE_TIMEOUT_SECONDS", 0)) * time.Second,
		},
		FLVE: FLVEConfig{
			Port: getEnvAsInt("FLVE_PORT", func() int {
				if flveYAML.FLVE.Port != 0 {
					return flveYAML.FLVE.Port
				} else {
					return 50051
				}
			}()),
			Backend: getEnv("FLVE_BACKEND", func() string {
				if flveYAML.FLVE.Backend != "" {
					return flveYAML.FLVE.Backend
				} else {
					return "hybrid"
				}
			}()),
			HFEndpoint: getEnv("FLVE_HF_ENDPOINT", flveYAML.FLVE.HFEndpoint),
			HFToken:    getEnv("FLVE_HF_TOKEN", flveYAML.FLVE.HFToken),
			MainCDN:    getEnv("MAIN_CDN", flveYAML.FLVE.MainCDN),
			CDNURL:     getEnv("CDN_URL", flveYAML.FLVE.CDNURL),
			AccessKey:  getEnv("ACCESS_KEY", flveYAML.FLVE.AccessKey),
			SecretKey:  getEnv("SECRET_KEY", flveYAML.FLVE.SecretKey),
			ModelsPath: getEnv("FLVE_MODELS_PATH", func() string {
				if flveYAML.FLVE.ModelsPath != "" {
					return flveYAML.FLVE.ModelsPath
				} else {
					return "./models"
				}
			}()),
			ModelDetector: getEnv("FLVE_MODEL_DETECTOR", func() string {
				if flveYAML.FLVE.ModelDetector != "" {
					return flveYAML.FLVE.ModelDetector
				} else {
					return "yolo-face.onnx"
				}
			}()),
			ModelEmbedding: getEnv("FLVE_MODEL_EMBEDDING", func() string {
				if flveYAML.FLVE.ModelEmbedding != "" {
					return flveYAML.FLVE.ModelEmbedding
				} else {
					return "arcface.onnx"
				}
			}()),
			ModelLiveness: getEnv("FLVE_MODEL_LIVENESS", func() string {
				if flveYAML.FLVE.ModelLiveness != "" {
					return flveYAML.FLVE.ModelLiveness
				} else {
					return "liveness.onnx"
				}
			}()),
			UseGPU: getEnvAsBool("FLVE_USE_GPU", flveYAML.FLVE.UseGPU),
			LivenessThreshold: getEnvAsFloat64("FLVE_LIVENESS_THRESHOLD", func() float64 {
				if flveYAML.FLVE.LivenessThreshold != 0 {
					return flveYAML.FLVE.LivenessThreshold
				} else {
					return 0.5
				}
			}()),
			MatchThreshold: getEnvAsFloat64("FLVE_MATCH_THRESHOLD", func() float64 {
				if flveYAML.FLVE.MatchThreshold != 0 {
					return flveYAML.FLVE.MatchThreshold
				} else {
					return 0.6
				}
			}()),
			SessionType: getEnv("FLVE_SESSION_TYPE", func() string {
				if flveYAML.FLVE.SessionType != "" {
					return flveYAML.FLVE.SessionType
				} else {
					return "memory"
				}
			}()),
			RedisURL: getEnv("FLVE_REDIS_URL", func() string {
				if flveYAML.FLVE.RedisURL != "" {
					return flveYAML.FLVE.RedisURL
				} else {
					return "localhost:6379"
				}
			}()),
			SessionTTL: getEnvAsInt("FLVE_SESSION_TTL", func() int {
				if flveYAML.FLVE.SessionTTL != 0 {
					return flveYAML.FLVE.SessionTTL
				} else {
					return 3600
				}
			}()),
			JWTSecret: getEnv("FLVE_JWT_SECRET", flveYAML.FLVE.JWTSecret),
		},
	}

	// Validation
	if err := cfg.Validate(); err != nil {
		logger.Errorf("config validation failed: %v", err)
		return nil, errors.New("config validation failed")
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JWT.PrivateKeyPath == "" {
		logger.Errorf("JWT_PRIVATE_KEY_PATH is required (RS256 private key PEM file)")
		return errors.New("JWT_PRIVATE_KEY_PATH is required (RS256 private key PEM file)")
	}
	if c.JWT.PublicKeyPath == "" {
		logger.Errorf("JWT_PUBLIC_KEY_PATH is required (RS256 public key PEM file)")
		return errors.New("JWT_PUBLIC_KEY_PATH is required (RS256 public key PEM file)")
	}
	if c.JWT.KeyID == "" {
		logger.Errorf("JWT_KEY_ID is required (kid header for JWKS)")
		return errors.New("JWT_KEY_ID is required (kid header for JWKS)")
	}
	if c.KYC.Enabled && c.KYC.Address == "" {
		logger.Errorf("KYC_SERVICE_ADDRESS is required when KYC_SERVICE_ENABLED=true")
		return errors.New("KYC_SERVICE_ADDRESS is required when KYC_SERVICE_ENABLED=true")
	}
	if c.SMS.APIBase == "" && c.SMS.MaskingEnabled {
		logger.Errorf("SSLWIRELESS_API_BASE is required when SMS is enabled")
		return errors.New("SSLWIRELESS_API_BASE is required when SMS is enabled")
	}
	if c.SMS.SID == "" && c.SMS.MaskingEnabled {
		logger.Errorf("SSLWIRELESS_SID is required when SMS is enabled")
		return errors.New("SSLWIRELESS_SID is required when SMS is enabled")
	}
	if c.SMS.APIKey == "" && c.SMS.MaskingEnabled {
		logger.Errorf("SSLWIRELESS_API_KEY is required when SMS is enabled")
		return errors.New("SSLWIRELESS_API_KEY is required when SMS is enabled")
	}
	if c.Database.Password == "" {
		logger.Errorf("DB_PASSWORD is required")
		return errors.New("DB_PASSWORD is required")
	}
	return nil
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	parts := strings.Split(valueStr, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	if len(result) == 0 {
		return defaultValue
	}
	return result
}

func getEnvAsFloat64(key string, defaultValue float64) float64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func loadAuthnServicePorts() (grpcPort int, httpPort int) {
	grpcPort = defaultAuthnGRPCPort
	httpPort = defaultAuthnHTTPPort

	type servicesConfig struct {
		Services map[string]struct {
			Ports struct {
				Grpc int `yaml:"grpc"`
				Http int `yaml:"http"`
			} `yaml:"ports"`
		} `yaml:"services"`
	}

	servicesConfigPath, err := opsconfig.ResolveConfigPath("services.yaml")
	if err != nil {
		logger.Warnf("Failed to resolve services.yaml for authn ports: %v (using defaults)", err)
		return grpcPort, httpPort
	}

	data, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		logger.Warnf("Failed to read services.yaml for authn ports: %v (using defaults)", err)
		return grpcPort, httpPort
	}

	var cfg servicesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Warnf("Failed to parse services.yaml for authn ports: %v (using defaults)", err)
		return grpcPort, httpPort
	}

	authnService, ok := cfg.Services["authn"]
	if !ok {
		logger.Warnf("Authn service not found in services.yaml (using defaults)")
		return grpcPort, httpPort
	}
	if authnService.Ports.Grpc > 0 {
		grpcPort = authnService.Ports.Grpc
	}
	if authnService.Ports.Http > 0 {
		httpPort = authnService.Ports.Http
	}
	return grpcPort, httpPort
}
