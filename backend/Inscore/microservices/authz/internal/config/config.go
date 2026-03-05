package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	opsconfig "github.com/newage-saint/insuretech/ops/config"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the authz microservice.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	Casbin   CasbinConfig
	Auth     AuthConfig
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

// AuthConfig holds JWT validation settings for the gRPC interceptor.
type AuthConfig struct {
	// PublicKeyPEM is the RS256 public key in PEM format used to verify JWTs.
	// Read from AUTHZ_JWT_PUBLIC_KEY_PEM env var.
	// If empty, JWT validation is disabled (no-op interceptor — dev/test only).
	PublicKeyPEM string
}

type ServerConfig struct {
	GRPCPort int
	HTTPPort int
	Host     string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
	// PolicyCacheTTL: how long Casbin policy decisions are cached in Redis.
	PolicyCacheTTL time.Duration
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// CasbinConfig controls Casbin enforcer behaviour.
type CasbinConfig struct {
	// ModelPath: path to casbin PERM model .conf file.
	// If empty, the built-in model is used (recommended).
	ModelPath string
	// AutoReloadInterval: how often to reload policies from DB (0 = disabled).
	AutoReloadInterval time.Duration
	// AuditAllDecisions: if true, every CheckAccess call writes to access_decision_audits.
	// For production set to false (only DENY events are always audited).
	AuditAllDecisions bool
	// DenyByDefault: if true, missing policy = DENY (always true — enforced in code).
	DenyByDefault bool
}

func Load() (*Config, error) {
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

	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: getEnvAsInt("AUTHZ_GRPC_PORT", 50052),
			HTTPPort: getEnvAsInt("AUTHZ_HTTP_PORT", 8081),
			Host:     getEnv("AUTHZ_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "insuretech_primary"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "insuretech_primary"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			URL:            getEnv("REDIS_URL", "redis://localhost:6379"),
			Password:       getEnv("REDIS_PASSWORD", ""),
			DB:             getEnvAsInt("REDIS_DB", 1), // separate DB from authn
			PolicyCacheTTL: getEnvAsDuration("AUTHZ_POLICY_CACHE_TTL", 5*time.Minute),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			Topic:   getEnv("KAFKA_AUTHZ_TOPIC", "authz-events"),
		},
		Casbin: CasbinConfig{
			ModelPath:          getEnv("CASBIN_MODEL_PATH", ""),
			AutoReloadInterval: getEnvAsDuration("CASBIN_RELOAD_INTERVAL", 30*time.Second),
			AuditAllDecisions:  getEnvAsBool("CASBIN_AUDIT_ALL", false),
			DenyByDefault:      true, // always deny-by-default — non-configurable
		},
		Auth: AuthConfig{
			PublicKeyPEM: resolveAuthZPublicKeyPEM(),
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
	if err := cfg.Validate(); err != nil {
		return nil, errors.New("authz config validation failed: " + err.Error())
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return errors.New("DB_PASSWORD is required")
	}
	return nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func getEnvAsInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
func getEnvAsBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}
func getEnvAsDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
func getEnvAsFloat64(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

func resolveAuthZPublicKeyPEM() string {
	// Preferred explicit authz-specific key.
	if pem := getEnv("AUTHZ_JWT_PUBLIC_KEY_PEM", ""); pem != "" {
		return pem
	}

	// Backward-compatible shared key env used by authn + seeder.
	if pem := getEnv("JWT_PUBLIC_KEY_PEM", ""); pem != "" {
		return pem
	}

	// File-path fallback (authz-specific first, then shared path).
	for _, keyPath := range []string{
		getEnv("AUTHZ_JWT_PUBLIC_KEY_PATH", ""),
		getEnv("JWT_PUBLIC_KEY_PATH", ""),
	} {
		if keyPath == "" {
			continue
		}
		data, err := os.ReadFile(keyPath)
		if err != nil {
			logger.Warnf("Failed to read JWT public key file from %s: %v", keyPath, err)
			continue
		}
		return string(data)
	}

	return ""
}
