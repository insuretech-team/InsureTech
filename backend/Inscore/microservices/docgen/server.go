package docgen

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/config"
	grpcserver "github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/grpc"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/kafka"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/service"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/worker"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	documentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

// DocGenServer holds the document generation server with all dependencies.
type DocGenServer struct {
	handler        documentservicev1.DocumentServiceServer
	kafkaPublisher *kafka.Publisher
	genWorker      *worker.GenerationWorker
	workerCancel   context.CancelFunc
}

// NewDocumentServer creates a document gRPC server with optional storage integration,
// Kafka publisher, and async generation worker.
func NewDocumentServer(db *sql.DB, storageConn *grpc.ClientConn) (*DocGenServer, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Warnf("failed to load config: %v (continuing with defaults)", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")
	templateRepo := repository.NewDocumentTemplateRepository(sqlxDB)
	generationRepo := repository.NewDocumentGenerationRepository(sqlxDB)

	var storageClient service.StorageClient
	if storageConn != nil {
		storageClient = storageservicev1.NewStorageServiceClient(storageConn)
	}

	docService, err := service.NewDocumentService(templateRepo, generationRepo, storageClient)
	if err != nil {
		return nil, err
	}

	// Initialize Kafka publisher if brokers are configured.
	var kafkaPublisher *kafka.Publisher
	if cfg != nil && len(cfg.KafkaBrokers) > 0 {
		kafkaPublisher, err = kafka.NewPublisher(cfg.KafkaBrokers, cfg.KafkaDocgenTopic)
		if err != nil {
			logger.Warnf("failed to initialize Kafka publisher: %v (continuing without Kafka)", err)
		}
	} else {
		logger.Warn("Kafka brokers not configured, Kafka publisher disabled")
	}

	docService.SetKafkaPublisher(kafkaPublisher)
	if cfg != nil {
		docService.SetPDFRenderer(cfg.GotenbergURL, cfg.MaxGenerationTimeout)
	}

	// Initialize async generation worker if enabled.
	var genWorker *worker.GenerationWorker
	var workerCancel context.CancelFunc
	if cfg != nil && cfg.AsyncGeneration {
		workerCtx, cancel := context.WithCancel(context.Background())
		workerCancel = cancel

		genWorker = worker.NewGenerationWorker(
			cfg.AsyncWorkerCount,
			100,
			func(ctx context.Context, req worker.GenerationRequest) error {
				payload := map[string]any{}
				for k, v := range req.Data {
					payload[k] = v
				}
				if strings.TrimSpace(req.GenerationID) != "" {
					payload["_generation_id"] = req.GenerationID
				}
				data, err := structpb.NewStruct(payload)
				if err != nil {
					return fmt.Errorf("invalid async generation payload: %w", err)
				}
				_, genErr := docService.GenerateDocument(
					ctx,
					req.TemplateID,
					req.EntityType,
					req.EntityID,
					data,
					true,
					req.TenantID,
					req.ActorID,
				)
				return genErr
			},
		)
		genWorker.Start(workerCtx)
		logger.Infof("Async generation worker started with %d workers", cfg.AsyncWorkerCount)
	}

	return &DocGenServer{
		handler:        grpcserver.NewDocumentHandler(docService),
		kafkaPublisher: kafkaPublisher,
		genWorker:      genWorker,
		workerCancel:   workerCancel,
	}, nil
}

// Handler returns the gRPC DocumentServiceServer handler.
func (s *DocGenServer) Handler() documentservicev1.DocumentServiceServer {
	return s.handler
}

// EnqueueGeneration adds a generation request to the async queue.
func (s *DocGenServer) EnqueueGeneration(req worker.GenerationRequest) error {
	if s.genWorker == nil {
		return errors.New("async generation worker is not enabled")
	}
	return s.genWorker.Enqueue(req)
}

// Close gracefully shuts down the server's resources.
func (s *DocGenServer) Close() error {
	var errs []error

	if s.genWorker != nil {
		if s.workerCancel != nil {
			s.workerCancel()
		}
		s.genWorker.Stop()
		logger.Info("Async generation worker stopped")
	}

	if s.kafkaPublisher != nil {
		if err := s.kafkaPublisher.Close(); err != nil {
			logger.Errorf("failed to close Kafka publisher: %v", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
