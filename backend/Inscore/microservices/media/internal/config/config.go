package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	opsconfig "github.com/newage-saint/insuretech/ops/config"
	"gopkg.in/yaml.v3"
)

const (
	defaultMediaGRPCPort = 50260
	defaultMediaHTTPPort = 50261
)

// Config holds all configuration for the media microservice
type Config struct {
	Port                int
	StorageServiceAddr  string
	KafkaBrokers        []string
	KafkaMediaTopic     string
	MaxFileSizeMB       int
	AllowedMIMETypes    []string
	WorkerCount         int
	VirusScanEnabled    bool
	OCREnabled          bool
	ClamAVAddr          string
	ThumbnailWidth      int
	ThumbnailHeight     int
}

// Load loads configuration from environment variables and services.yaml
func Load() (*Config, error) {
	mediaPort := loadMediaServicePort()

	cfg := &Config{
		Port:               mediaPort,
		StorageServiceAddr: getEnv("STORAGE_SERVICE_ADDR", ""),
		KafkaBrokers:       getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		KafkaMediaTopic:    getEnv("KAFKA_MEDIA_TOPIC", "media-events"),
		MaxFileSizeMB:      getEnvAsInt("MEDIA_MAX_FILE_SIZE_MB", 50),
		AllowedMIMETypes:   getEnvAsSlice("MEDIA_ALLOWED_TYPES", []string{"image/jpeg", "image/png", "image/webp", "application/pdf", "video/mp4"}),
		WorkerCount:        getEnvAsInt("MEDIA_WORKER_COUNT", 5),
		VirusScanEnabled:   getEnvAsBool("MEDIA_VIRUS_SCAN_ENABLED", false),
		OCREnabled:         getEnvAsBool("MEDIA_OCR_ENABLED", false),
		ClamAVAddr:         getEnv("CLAMAV_ADDR", "localhost:3310"),
		ThumbnailWidth:     getEnvAsInt("MEDIA_THUMBNAIL_WIDTH", 300),
		ThumbnailHeight:    getEnvAsInt("MEDIA_THUMBNAIL_HEIGHT", 300),
	}

	logger.Infof("Media service config loaded: port=%d, workers=%d, virusScan=%v, ocr=%v",
		cfg.Port, cfg.WorkerCount, cfg.VirusScanEnabled, cfg.OCREnabled)

	return cfg, nil
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

func loadMediaServicePort() int {
	grpcPort := defaultMediaGRPCPort

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
		logger.Warnf("Failed to resolve services.yaml for media port: %v (using default %d)", err, grpcPort)
		return grpcPort
	}

	data, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		logger.Warnf("Failed to read services.yaml for media port: %v (using default %d)", err, grpcPort)
		return grpcPort
	}

	var cfg servicesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Warnf("Failed to parse services.yaml for media port: %v (using default %d)", err, grpcPort)
		return grpcPort
	}

	mediaService, ok := cfg.Services["media"]
	if !ok {
		logger.Warnf("Media service not found in services.yaml (using default %d)", grpcPort)
		return grpcPort
	}
	if mediaService.Ports.Grpc > 0 {
		grpcPort = mediaService.Ports.Grpc
	}
	return grpcPort
}
