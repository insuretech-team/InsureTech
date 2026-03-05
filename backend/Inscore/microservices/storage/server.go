package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	storageevents "github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/s3"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/service"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"gopkg.in/yaml.v3"
)

// NewStorageServer creates a new storage gRPC server with all dependencies
func NewStorageServer(db *sql.DB) (storageservicev1.StorageServiceServer, error) {
	return NewStorageServerWithProducer(db, nil)
}

// NewStorageServerWithProducer creates a new storage gRPC server with optional Kafka producer.
func NewStorageServerWithProducer(db *sql.DB, producer storageevents.EventProducer) (storageservicev1.StorageServiceServer, error) {
	// Wrap sql.DB with sqlx
	sqlxDB := sqlx.NewDb(db, "postgres")

	// Create file repository
	fileRepo := repository.NewFileRepository(sqlxDB)

	// Load S3 configuration
	configPath, err := config.ResolveConfigPath("s3.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve s3 config path: %w", err)
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open s3 config: %w", err)
	}
	defer file.Close()

	var cfg struct {
		Storage struct {
			S3 struct {
				Bucket      string `yaml:"bucket"`
				Region      string `yaml:"region"`
				Endpoint    string `yaml:"endpoint"`
				CDNEndpoint string `yaml:"cdn_endpoint"`
				RootFolder  string `yaml:"root_folder"`
			} `yaml:"s3"`
		} `yaml:"storage"`
	}

	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode s3 config: %w", err)
	}

	// Expand environment variables in config values
	cfg.Storage.S3.Bucket = os.ExpandEnv(cfg.Storage.S3.Bucket)
	cfg.Storage.S3.Region = os.ExpandEnv(cfg.Storage.S3.Region)
	cfg.Storage.S3.Endpoint = os.ExpandEnv(cfg.Storage.S3.Endpoint)
	cfg.Storage.S3.CDNEndpoint = os.ExpandEnv(cfg.Storage.S3.CDNEndpoint)
	cfg.Storage.S3.RootFolder = os.ExpandEnv(cfg.Storage.S3.RootFolder)

	// Create S3 client from config and environment variables
	// Secrets (ACCESS_KEY_ID, SECRET_ACCESS_KEY) are only from environment variables
	s3Config := s3.Config{
		AccessKeyID:     os.Getenv("SPACES_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("SPACES_SECRET_ACCESS_KEY"),
		Bucket:          cfg.Storage.S3.Bucket,
		Region:          cfg.Storage.S3.Region,
		Endpoint:        cfg.Storage.S3.Endpoint,
		CDNEndpoint:     cfg.Storage.S3.CDNEndpoint,
		RootFolder:      cfg.Storage.S3.RootFolder,
	}

	s3Client, err := s3.NewClient(s3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	// Load optional storage layout templates.
	layoutConfigPath, err := config.ResolveConfigPath("storage_layout.yaml")
	if err == nil {
		layoutFile, openErr := os.Open(layoutConfigPath)
		if openErr == nil {
			defer layoutFile.Close()

			var layoutCfg struct {
				StorageLayout struct {
					BasePrefix string `yaml:"base_prefix"`
					Categories map[string]struct {
						ReferenceType    string `yaml:"reference_type"`
						FolderTemplate   string `yaml:"folder_template"`
						FilenameTemplate string `yaml:"filename_template"`
					} `yaml:"categories"`
				} `yaml:"storage_layout"`
			}
			if decodeErr := yaml.NewDecoder(layoutFile).Decode(&layoutCfg); decodeErr == nil {
				layoutTemplates := make(map[string]s3.LayoutTemplate, len(layoutCfg.StorageLayout.Categories))
				for _, cat := range layoutCfg.StorageLayout.Categories {
					if cat.ReferenceType == "" {
						continue
					}
					layoutTemplates[cat.ReferenceType] = s3.LayoutTemplate{
						FolderTemplate:   cat.FolderTemplate,
						FilenameTemplate: cat.FilenameTemplate,
					}
				}
				basePrefix := os.ExpandEnv(layoutCfg.StorageLayout.BasePrefix)
				s3Client.SetLayout(basePrefix, layoutTemplates)
			}
		}
	}

	// Create storage service
	eventPublisher := storageevents.NewPublisher(producer)
	storageService := service.NewStorageServiceWithPublisher(fileRepo, s3Client, eventPublisher)

	// Create and return gRPC handler
	return server.NewStorageHandler(storageService), nil
}
