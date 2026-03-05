package docgen

import (
	"context"
	"database/sql"

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
)

// DocGenServer holds the document generation server with all dependencies
type DocGenServer struct {
	handler        documentservicev1.DocumentServiceServer
	kafkaPublisher *kafka.Publisher
	genWorker      *worker.GenerationWorker
}

// NewDocumentServer creates a document gRPC server with optional storage integration, Kafka publisher, and async worker.
func NewDocumentServer(db *sql.DB, storageConn *grpc.ClientConn) (*DocGenServer, error) {
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

	// Initialize Kafka publisher if brokers are configured
	var kafkaPublisher *kafka.Publisher
	if cfg != nil && len(cfg.KafkaBrokers) > 0 {
		kafkaPublisher, err = kafka.NewPublisher(cfg.KafkaBrokers, cfg.KafkaDocgenTopic)
		if err != nil {
			logger.Warnf("failed to initialize Kafka publisher: %v (continuing without Kafka)", err)
		}
	} else {
		logger.Warn("Kafka brokers not configured, Kafka publisher disabled")
	}

	// Inject Kafka publisher into the service
	docService.SetKafkaPublisher(kafkaPublisher)

	// Initialize async generation worker if enabled
	var genWorker *worker.GenerationWorker
	if cfg != nil && cfg.AsyncGeneration {
		genWorker = worker.NewGenerationWorker(
			cfg.AsyncWorkerCount,
			100, // queueSize
			func(ctx context.Context, req worker.GenerationRequest) error {
				// Generator function implementation would go here
				logger.Infof("Async generation request: %s", req.GenerationID)
				return nil
			},
		)
	}

	server := &DocGenServer{
		handler:        grpcserver.NewDocumentHandler(docService),
		kafkaPublisher: kafkaPublisher,
		genWorker:      genWorker,
	}

	// Start the generation worker if initialized
	if genWorker != nil {
		genWorker.Start(context.Background())
		logger.Info("Async generation worker started")
	}

	return server, nil
}

// Handler returns the gRPC DocumentServiceServer handler
func (s *DocGenServer) Handler() documentservicev1.DocumentServiceServer {
	return s.handler
}

// Close gracefully shuts down the server's resources
func (s *DocGenServer) Close() error {
	var errs []error

	if s.genWorker != nil {
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
