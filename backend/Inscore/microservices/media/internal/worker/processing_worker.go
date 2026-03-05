package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// MediaDownloader provides media file download functionality.
type MediaDownloader interface {
	DownloadFile(ctx context.Context, mediaID string) ([]byte, string, error)
}

// MediaUpdater provides media file update functionality.
type MediaUpdater interface {
	UpdateProcessingResult(ctx context.Context, mediaID string, ocrText string, width, height int) error
	UpdateVirusScanResult(ctx context.Context, mediaID string, clean bool, virusName string) error
}

// JobUpdater provides processing job update functionality.
type JobUpdater interface {
	MarkJobStarted(ctx context.Context, jobID string) error
	MarkJobCompleted(ctx context.Context, jobID string, result string) error
	MarkJobFailed(ctx context.Context, jobID string, errMsg string) error
	GetPendingJobs(ctx context.Context, limit int) ([]*ProcessingJobRecord, error)
}

// ProcessingJobRecord represents a media processing job.
type ProcessingJobRecord struct {
	ID       string
	MediaID  string
	TenantID string
	JobType  string // THUMBNAIL, OPTIMIZATION, OCR, VIRUS_SCAN
	Priority int
}

// EventPublisher provides event publishing functionality.
type EventPublisher interface {
	PublishProcessingStarted(ctx context.Context, jobID, mediaID, tenantID, jobType string) error
	PublishProcessingCompleted(ctx context.Context, jobID, mediaID, tenantID, jobType string) error
	PublishProcessingFailed(ctx context.Context, jobID, mediaID, tenantID, jobType, reason string) error
	PublishVirusDetected(ctx context.Context, mediaID, tenantID, virusName string) error
}

// ImageProcessorI provides image processing functionality.
type ImageProcessorI interface {
	GenerateThumbnail(data []byte, width, height int) ([]byte, error)
	CompressImage(data []byte, quality int, mimeType string) ([]byte, error)
	GetImageDimensions(data []byte) (width, height int, err error)
}

// OCRProcessorI provides OCR processing functionality.
type OCRProcessorI interface {
	ExtractText(ctx context.Context, data []byte, mimeType string) (string, error)
	IsEnabled() bool
}

// VirusScannerI provides virus scanning functionality.
type VirusScannerI interface {
	Scan(ctx context.Context, data []byte) (clean bool, virusName string, err error)
	IsEnabled() bool
}

// ProcessingWorker manages media file processing with a worker pool.
type ProcessingWorker struct {
	workerCount  int
	pollInterval time.Duration
	thumbnailW   int
	thumbnailH   int
	downloader   MediaDownloader
	mediaUpdater MediaUpdater
	jobUpdater   JobUpdater
	publisher    EventPublisher
	imageProc    ImageProcessorI
	ocrProc      OCRProcessorI
	virusScanner VirusScannerI
	wg           sync.WaitGroup
}

// NewProcessingWorker creates a new media processing worker pool.
func NewProcessingWorker(
	workerCount, thumbnailW, thumbnailH int,
	downloader MediaDownloader,
	mediaUpdater MediaUpdater,
	jobUpdater JobUpdater,
	publisher EventPublisher,
	imageProc ImageProcessorI,
	ocrProc OCRProcessorI,
	virusScanner VirusScannerI,
) *ProcessingWorker {
	if workerCount <= 0 {
		workerCount = 1
	}
	if thumbnailW <= 0 {
		thumbnailW = 200
	}
	if thumbnailH <= 0 {
		thumbnailH = 200
	}
	return &ProcessingWorker{
		workerCount:  workerCount,
		pollInterval: 5 * time.Second,
		thumbnailW:   thumbnailW,
		thumbnailH:   thumbnailH,
		downloader:   downloader,
		mediaUpdater: mediaUpdater,
		jobUpdater:   jobUpdater,
		publisher:    publisher,
		imageProc:    imageProc,
		ocrProc:      ocrProc,
		virusScanner: virusScanner,
	}
}

// Start initiates the worker goroutines that poll for pending jobs.
func (w *ProcessingWorker) Start(ctx context.Context) {
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.pollingWorker(ctx, i)
	}
}

// Stop waits for all workers to finish.
func (w *ProcessingWorker) Stop() {
	w.wg.Wait()
}

// pollingWorker is the main worker loop that polls for pending jobs.
func (w *ProcessingWorker) pollingWorker(ctx context.Context, workerID int) {
	defer w.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[media-worker-%d] panicked: %v", workerID, r)
		}
	}()

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("[media-worker-%d] shutting down due to context cancellation", workerID)
			return
		case <-ticker.C:
			w.pollAndProcess(ctx, workerID)
		}
	}
}

// pollAndProcess polls for pending jobs and processes them.
func (w *ProcessingWorker) pollAndProcess(ctx context.Context, workerID int) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[media-worker-%d] panic in pollAndProcess: %v", workerID, r)
		}
	}()

	jobs, err := w.jobUpdater.GetPendingJobs(ctx, 10)
	if err != nil {
		logger.Errorf("[media-worker-%d] failed to get pending jobs: %v", workerID, err)
		return
	}

	for _, job := range jobs {
		if job == nil {
			continue
		}
		w.processJob(ctx, job, workerID)
	}
}

// processJob dispatches job processing based on job type.
func (w *ProcessingWorker) processJob(ctx context.Context, job *ProcessingJobRecord, workerID int) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[media-worker-%d] panic processing job %s: %v", workerID, job.ID, r)
		}
	}()

	logger.Infof("[media-worker-%d] processing job id=%s type=%s media=%s", workerID, job.ID, job.JobType, job.MediaID)

	// Publish processing started event
	if w.publisher != nil {
		if err := w.publisher.PublishProcessingStarted(ctx, job.ID, job.MediaID, job.TenantID, job.JobType); err != nil {
			logger.Warnf("[media-worker-%d] failed to publish processing started event: %v", workerID, err)
		}
	}

	// Mark job as started
	if err := w.jobUpdater.MarkJobStarted(ctx, job.ID); err != nil {
		logger.Errorf("[media-worker-%d] failed to mark job as started: %v", workerID, err)
		return
	}

	var procErr error
	switch job.JobType {
	case "THUMBNAIL":
		procErr = w.processThumbnail(ctx, job, workerID)
	case "OPTIMIZATION":
		procErr = w.processOptimization(ctx, job, workerID)
	case "OCR":
		procErr = w.processOCR(ctx, job, workerID)
	case "VIRUS_SCAN":
		procErr = w.processVirusScan(ctx, job, workerID)
	default:
		procErr = fmt.Errorf("unknown job type: %s", job.JobType)
	}

	if procErr != nil {
		logger.Errorf("[media-worker-%d] job %s failed: %v", workerID, job.ID, procErr)
		if err := w.jobUpdater.MarkJobFailed(ctx, job.ID, procErr.Error()); err != nil {
			logger.Errorf("[media-worker-%d] failed to mark job as failed: %v", workerID, err)
		}
		if w.publisher != nil {
			if err := w.publisher.PublishProcessingFailed(ctx, job.ID, job.MediaID, job.TenantID, job.JobType, procErr.Error()); err != nil {
				logger.Warnf("[media-worker-%d] failed to publish processing failed event: %v", workerID, err)
			}
		}
		return
	}

	if err := w.jobUpdater.MarkJobCompleted(ctx, job.ID, ""); err != nil {
		logger.Errorf("[media-worker-%d] failed to mark job as completed: %v", workerID, err)
		return
	}

	if w.publisher != nil {
		if err := w.publisher.PublishProcessingCompleted(ctx, job.ID, job.MediaID, job.TenantID, job.JobType); err != nil {
			logger.Warnf("[media-worker-%d] failed to publish processing completed event: %v", workerID, err)
		}
	}

	logger.Infof("[media-worker-%d] job %s completed successfully", workerID, job.ID)
}

// processThumbnail handles thumbnail generation.
func (w *ProcessingWorker) processThumbnail(ctx context.Context, job *ProcessingJobRecord, workerID int) error {
	logger.Infof("[media-worker-%d] generating thumbnail for media=%s", workerID, job.MediaID)

	data, mimeType, err := w.downloader.DownloadFile(ctx, job.MediaID)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	logger.Infof("[media-worker-%d] downloaded file mime=%s size=%d", workerID, mimeType, len(data))

	if w.imageProc == nil {
		return fmt.Errorf("image processor not available")
	}

	width, height, err := w.imageProc.GetImageDimensions(data)
	if err != nil {
		return fmt.Errorf("failed to get image dimensions: %w", err)
	}

	_, err = w.imageProc.GenerateThumbnail(data, w.thumbnailW, w.thumbnailH)
	if err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	if err := w.mediaUpdater.UpdateProcessingResult(ctx, job.MediaID, "", width, height); err != nil {
		return fmt.Errorf("failed to update processing result: %w", err)
	}

	logger.Infof("[media-worker-%d] thumbnail completed media=%s dims=%dx%d", workerID, job.MediaID, width, height)
	return nil
}

// processOptimization handles image optimization.
func (w *ProcessingWorker) processOptimization(ctx context.Context, job *ProcessingJobRecord, workerID int) error {
	logger.Infof("[media-worker-%d] optimizing image for media=%s", workerID, job.MediaID)

	data, mimeType, err := w.downloader.DownloadFile(ctx, job.MediaID)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	if w.imageProc == nil {
		return fmt.Errorf("image processor not available")
	}

	compressedData, err := w.imageProc.CompressImage(data, 80, mimeType)
	if err != nil {
		return fmt.Errorf("failed to compress image: %w", err)
	}

	logger.Infof("[media-worker-%d] image optimized media=%s original=%d compressed=%d bytes",
		workerID, job.MediaID, len(data), len(compressedData))
	// TODO: Upload compressed image to storage
	return nil
}

// processOCR handles OCR text extraction.
func (w *ProcessingWorker) processOCR(ctx context.Context, job *ProcessingJobRecord, workerID int) error {
	if w.ocrProc == nil || !w.ocrProc.IsEnabled() {
		logger.Warnf("[media-worker-%d] OCR not enabled, skipping job %s", workerID, job.ID)
		return nil // non-fatal: skip gracefully
	}

	logger.Infof("[media-worker-%d] extracting text via OCR for media=%s", workerID, job.MediaID)

	data, mimeType, err := w.downloader.DownloadFile(ctx, job.MediaID)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	ocrText, err := w.ocrProc.ExtractText(ctx, data, mimeType)
	if err != nil {
		return fmt.Errorf("failed to extract text: %w", err)
	}

	if err := w.mediaUpdater.UpdateProcessingResult(ctx, job.MediaID, ocrText, 0, 0); err != nil {
		return fmt.Errorf("failed to update processing result: %w", err)
	}

	logger.Infof("[media-worker-%d] OCR completed media=%s text_length=%d", workerID, job.MediaID, len(ocrText))
	return nil
}

// processVirusScan handles virus scanning.
func (w *ProcessingWorker) processVirusScan(ctx context.Context, job *ProcessingJobRecord, workerID int) error {
	if w.virusScanner == nil || !w.virusScanner.IsEnabled() {
		logger.Warnf("[media-worker-%d] virus scanner not enabled, skipping job %s", workerID, job.ID)
		return nil // non-fatal: skip gracefully
	}

	logger.Infof("[media-worker-%d] scanning file for viruses media=%s", workerID, job.MediaID)

	data, mimeType, err := w.downloader.DownloadFile(ctx, job.MediaID)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	_ = mimeType

	clean, virusName, err := w.virusScanner.Scan(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	if err := w.mediaUpdater.UpdateVirusScanResult(ctx, job.MediaID, clean, virusName); err != nil {
		return fmt.Errorf("failed to update virus scan result: %w", err)
	}

	if !clean {
		logger.Warnf("[media-worker-%d] virus detected media=%s virus=%s", workerID, job.MediaID, virusName)
		if w.publisher != nil {
			if err := w.publisher.PublishVirusDetected(ctx, job.MediaID, job.TenantID, virusName); err != nil {
				logger.Warnf("[media-worker-%d] failed to publish virus detected event: %v", workerID, err)
			}
		}
		return fmt.Errorf("virus detected: %s", virusName)
	}

	logger.Infof("[media-worker-%d] virus scan clean media=%s", workerID, job.MediaID)
	return nil
}
