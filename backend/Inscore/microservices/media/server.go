package media

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/config"
	grpcserver "github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/grpc"
	mediakafka "github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/kafka"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/processor"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/service"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/worker"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	mediaservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/services/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// jobRepoAdapter adapts *repository.ProcessingJobRepository to worker.JobUpdater.
// The repository returns []*repository.ProcessingJobRecord but the worker interface
// requires []*worker.ProcessingJobRecord — this adapter bridges the two.
type jobRepoAdapter struct {
	repo *repository.ProcessingJobRepository
}

type mediaDownloaderAdapter struct {
	repo          *repository.MediaRepository
	storageClient service.StorageDownloadClient
	httpClient    *http.Client
}

func (a *jobRepoAdapter) MarkJobStarted(ctx context.Context, jobID string) error {
	return a.repo.MarkJobStarted(ctx, jobID)
}

func (a *jobRepoAdapter) MarkJobCompleted(ctx context.Context, jobID string, result string) error {
	return a.repo.MarkJobCompleted(ctx, jobID, result)
}

func (a *jobRepoAdapter) MarkJobFailed(ctx context.Context, jobID string, errMsg string) error {
	return a.repo.MarkJobFailed(ctx, jobID, errMsg)
}

func (a *jobRepoAdapter) GetPendingJobs(ctx context.Context, limit int) ([]*worker.ProcessingJobRecord, error) {
	repoJobs, err := a.repo.GetPendingJobs(ctx, limit)
	if err != nil {
		return nil, err
	}
	out := make([]*worker.ProcessingJobRecord, 0, len(repoJobs))
	for _, j := range repoJobs {
		if j == nil {
			continue
		}
		out = append(out, &worker.ProcessingJobRecord{
			ID:       j.ID,
			MediaID:  j.MediaID,
			TenantID: j.TenantID,
			JobType:  j.JobType,
			Priority: j.Priority,
		})
	}
	return out, nil
}

func (a *mediaDownloaderAdapter) DownloadFile(ctx context.Context, tenantID, mediaID string) ([]byte, string, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, "", fmt.Errorf("tenant_id is required for media download")
	}
	if strings.TrimSpace(mediaID) == "" {
		return nil, "", fmt.Errorf("media_id is required for media download")
	}

	media, err := a.repo.GetByID(ctx, tenantID, mediaID)
	if err != nil {
		return nil, "", fmt.Errorf("load media metadata: %w", err)
	}
	if a.storageClient == nil {
		return nil, "", fmt.Errorf("storage download client is not configured")
	}

	resp, err := a.storageClient.GetDownloadURL(ctx, &storageservicev1.GetDownloadURLRequest{
		TenantId: tenantID,
		FileId:   media.FileId,
	})
	if err != nil {
		return nil, "", fmt.Errorf("resolve download url: %w", err)
	}
	if strings.TrimSpace(resp.GetDownloadUrl()) == "" {
		return nil, "", fmt.Errorf("storage service returned an empty download URL")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, resp.GetDownloadUrl(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("build download request: %w", err)
	}

	httpResp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("download media bytes: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < http.StatusOK || httpResp.StatusCode >= http.StatusMultipleChoices {
		return nil, "", fmt.Errorf("download media bytes: unexpected status %s", httpResp.Status)
	}

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read media bytes: %w", err)
	}

	mimeType := strings.TrimSpace(httpResp.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = media.MimeType
	}

	return data, mimeType, nil
}

// MediaServer wraps the gRPC service with supporting infrastructure
type MediaServer struct {
	svc                mediaservicev1.MediaServiceServer
	kafkaPublisher     *mediakafka.Publisher
	processingWorker   *worker.ProcessingWorker
	workerCtx          context.Context
	workerCtxCancel    context.CancelFunc
}

// NewMediaServer creates a media gRPC server with optional storage integration and Kafka/worker integration.
// It accepts a *gorm.DB (returned by db.GetDB()) and extracts the underlying *sql.DB for sqlx.
func NewMediaServer(gormDB *gorm.DB, storageConn *grpc.ClientConn) (*MediaServer, error) {
	// Extract underlying *sql.DB from gorm so sqlx can wrap it
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB from gorm: %w", err)
	}
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Warnf("Failed to load media config: %v (using defaults)", err)
		cfg = &config.Config{
			WorkerCount:     5,
			ThumbnailWidth:  300,
			ThumbnailHeight: 300,
		}
	}

	sqlxDB := sqlx.NewDb(sqlDB, "postgres")

	mediaRepo := repository.NewMediaRepository(sqlxDB)
	jobRepo := repository.NewProcessingJobRepository(sqlxDB)

	var storageClient service.StorageDownloadClient
	if storageConn != nil {
		storageClient = storageservicev1.NewStorageServiceClient(storageConn)
	}

	mediaService := service.NewMediaServiceWithStorage(mediaRepo, jobRepo, storageClient)

	// Initialize Kafka Publisher (if brokers are configured)
	var kafkaPublisher *mediakafka.Publisher
	if len(cfg.KafkaBrokers) > 0 {
		pub, err := mediakafka.NewPublisher(cfg.KafkaBrokers, cfg.KafkaMediaTopic)
		if err != nil {
			logger.Warnf("Failed to initialize Kafka publisher: %v (events will not be published)", err)
		} else {
			kafkaPublisher = pub
			mediaService.SetKafkaPublisher(kafkaPublisher)
		}
	}

	// Initialize Processors
	imageProcessor := processor.NewImageProcessor()
	ocrProcessor := processor.NewOCRProcessor(cfg.OCREnabled)
	virusScanner := processor.NewVirusScanner(cfg.VirusScanEnabled, cfg.ClamAVAddr)

	// Initialize Processing Worker (jobRepo adapted to worker.JobUpdater interface)
	processingWorker := worker.NewProcessingWorker(
		cfg.WorkerCount,
		cfg.ThumbnailWidth,
		cfg.ThumbnailHeight,
		&mediaDownloaderAdapter{
			repo:          mediaRepo,
			storageClient: storageClient,
			httpClient:    &http.Client{},
		},
		mediaRepo,                      // implements worker.MediaUpdater
		&jobRepoAdapter{repo: jobRepo}, // adapts ProcessingJobRepository to worker.JobUpdater
		kafkaPublisher,
		imageProcessor,
		ocrProcessor,
		virusScanner,
	)

	// Create context for worker
	workerCtx, workerCtxCancel := context.WithCancel(context.Background())

	// Start processing worker
	processingWorker.Start(workerCtx)
	logger.Infof("Media processing worker pool started (workers=%d)", cfg.WorkerCount)

	return &MediaServer{
		svc:              grpcserver.NewMediaHandler(mediaService),
		kafkaPublisher:   kafkaPublisher,
		processingWorker: processingWorker,
		workerCtx:        workerCtx,
		workerCtxCancel:  workerCtxCancel,
	}, nil
}

// GetGRPCServer returns the gRPC service server
func (ms *MediaServer) GetGRPCServer() mediaservicev1.MediaServiceServer {
	return ms.svc
}

// Close gracefully shuts down the media server and all supporting infrastructure
func (ms *MediaServer) Close() error {
	logger.Info("Shutting down media server...")

	// Stop the processing worker
	if ms.processingWorker != nil {
		ms.workerCtxCancel()
		ms.processingWorker.Stop()
		logger.Info("Processing worker stopped")
	}

	// Close Kafka publisher
	if ms.kafkaPublisher != nil {
		if err := ms.kafkaPublisher.Close(); err != nil {
			logger.Errorf("Failed to close Kafka publisher: %v", err)
		} else {
			logger.Info("Kafka publisher closed")
		}
	}

	return nil
}
