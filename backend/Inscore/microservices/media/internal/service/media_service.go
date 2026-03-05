package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/repository"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	mediav1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/entity/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	"google.golang.org/grpc"
)

// StorageDownloadClient is the subset of storage RPCs media needs for secure downloads.
type StorageDownloadClient interface {
	GetDownloadURL(ctx context.Context, in *storageservicev1.GetDownloadURLRequest, opts ...grpc.CallOption) (*storageservicev1.GetDownloadURLResponse, error)
}

// MediaService handles media file operations.
type MediaService struct {
	mediaRepo *repository.MediaRepository
	jobRepo   *repository.ProcessingJobRepository
	kafkaPublisher         interface{} // *mediakafka.Publisher, nil-safe
	storageDownloadClient  StorageDownloadClient
	fallbackDownloadBase   string
	defaultURLExpiryMinute int32
}

var (
	ErrMediaNotFound = errors.New("media file not found")
	ErrJobNotFound   = errors.New("processing job not found")
	ErrInvalidInput  = errors.New("invalid input")
)

// UploadMediaInput contains input for uploading media.
type UploadMediaInput struct {
	FileID        string
	TenantID      string
	MediaType     mediav1.MediaType
	MimeType      string
	FileSizeBytes int64
	EntityType    string
	EntityID      string
	UploadedBy    string
	Width         int32
	Height        int32
	DPI           int32
}

// NewMediaService creates a new media service.
func NewMediaService(mediaRepo *repository.MediaRepository, jobRepo *repository.ProcessingJobRepository) *MediaService {
	return NewMediaServiceWithStorage(mediaRepo, jobRepo, nil)
}

// NewMediaServiceWithStorage creates a media service with optional storage URL resolver.
func NewMediaServiceWithStorage(
	mediaRepo *repository.MediaRepository,
	jobRepo *repository.ProcessingJobRepository,
	storageClient StorageDownloadClient,
) *MediaService {
	fallbackBase := strings.TrimRight(strings.TrimSpace(os.Getenv("INSURETECH_CDN_URL")), "/")
	if fallbackBase == "" {
		mainCDN := strings.TrimRight(strings.TrimSpace(os.Getenv("MAIN_CDN")), "/")
		if mainCDN != "" {
			fallbackBase = mainCDN + "/insuretech"
		}
	}
	expiryMinutes := int32(60)
	if raw := strings.TrimSpace(os.Getenv("MEDIA_DOWNLOAD_URL_EXPIRES_MINUTES")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			expiryMinutes = int32(parsed)
		}
	}

	return &MediaService{
		mediaRepo:              mediaRepo,
		jobRepo:                jobRepo,
		kafkaPublisher:         nil,
		storageDownloadClient:  storageClient,
		fallbackDownloadBase:   fallbackBase,
		defaultURLExpiryMinute: expiryMinutes,
	}
}

// SetKafkaPublisher sets the Kafka publisher for event publishing
func (s *MediaService) SetKafkaPublisher(pub interface{}) {
	s.kafkaPublisher = pub
}

// UploadMedia creates a new media file record.
func (s *MediaService) UploadMedia(ctx context.Context, input *UploadMediaInput) (*mediav1.MediaFile, error) {
	if strings.TrimSpace(input.TenantID) == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if input.FileID == "" {
		return nil, fmt.Errorf("%w: file_id is required", ErrInvalidInput)
	}
	if input.MimeType == "" {
		return nil, fmt.Errorf("%w: mime_type is required", ErrInvalidInput)
	}
	if input.UploadedBy == "" {
		return nil, fmt.Errorf("%w: uploaded_by is required", ErrInvalidInput)
	}

	media := &mediav1.MediaFile{
		Id:               uuid.New().String(),
		FileId:           input.FileID,
		TenantId:         input.TenantID,
		MediaType:        input.MediaType,
		MimeType:         input.MimeType,
		FileSizeBytes:    input.FileSizeBytes,
		Width:            input.Width,
		Height:           input.Height,
		Dpi:              input.DPI,
		EntityType:       input.EntityType,
		EntityId:         input.EntityID,
		ValidationStatus: mediav1.ValidationStatus_VALIDATION_STATUS_PENDING,
		VirusScanStatus:  mediav1.VirusScanStatus_VIRUS_SCAN_STATUS_PENDING,
		UploadedBy:       input.UploadedBy,
		AuditInfo: &commonv1.AuditInfo{
			CreatedBy: input.UploadedBy,
			UpdatedBy: input.UploadedBy,
		},
	}

	createdMedia, err := s.mediaRepo.Create(ctx, media)
	if err != nil {
		return nil, fmt.Errorf("failed to create media file: %w", err)
	}

	// Publish file uploaded event (fire-and-forget in goroutine)
	if s.kafkaPublisher != nil {
		go s.publishFileUploadedAsync(ctx, createdMedia.Id, input.TenantID, input.FileID, input.MimeType, input.FileSizeBytes)
	}

	// Enqueue processing jobs based on media type
	go s.enqueueProcessingJobs(context.Background(), createdMedia.Id, input.TenantID, input.MimeType)

	return createdMedia, nil
}

// GetMedia retrieves a media file by ID.
func (s *MediaService) GetMedia(ctx context.Context, tenantID, mediaID string) (*mediav1.MediaFile, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if mediaID == "" {
		return nil, fmt.Errorf("%w: media_id is required", ErrInvalidInput)
	}

	media, err := s.mediaRepo.GetByID(ctx, tenantID, mediaID)
	if err != nil {
		if errors.Is(err, repository.ErrMediaNotFound) {
			return nil, ErrMediaNotFound
		}
		return nil, fmt.Errorf("failed to get media file: %w", err)
	}

	return media, nil
}

// ListMediaByEntity retrieves media files for an entity with optional filters.
func (s *MediaService) ListMediaByEntity(
	ctx context.Context,
	tenantID, entityType, entityID string,
	mediaType *mediav1.MediaType,
	validationStatus *mediav1.ValidationStatus,
	page, pageSize int,
) ([]*mediav1.MediaFile, int, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, 0, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if entityType == "" || entityID == "" {
		return nil, 0, fmt.Errorf("%w: entity_type and entity_id are required", ErrInvalidInput)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	mediaFiles, total, err := s.mediaRepo.ListByEntity(ctx, tenantID, entityType, entityID, mediaType, validationStatus, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list media files: %w", err)
	}

	return mediaFiles, total, nil
}

// ValidateMedia validates a media file.
func (s *MediaService) ValidateMedia(ctx context.Context, tenantID, mediaID string, validationRules []string) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if mediaID == "" {
		return fmt.Errorf("%w: media_id is required", ErrInvalidInput)
	}

	media, err := s.mediaRepo.GetByID(ctx, tenantID, mediaID)
	if err != nil {
		if errors.Is(err, repository.ErrMediaNotFound) {
			return ErrMediaNotFound
		}
		return fmt.Errorf("failed to get media file: %w", err)
	}

	_ = validationRules // TODO: Implement validation rules processing.
	validationStatus := mediav1.ValidationStatus_VALIDATION_STATUS_VALIDATED
	validationErrors := ""
	if media.FileSizeBytes > 10*1024*1024 {
		validationStatus = mediav1.ValidationStatus_VALIDATION_STATUS_REJECTED
		validationErrors = "File size exceeds 10MB limit"
	}

	if err := s.mediaRepo.UpdateValidationStatus(ctx, tenantID, mediaID, validationStatus, validationErrors); err != nil {
		return fmt.Errorf("failed to update validation status: %w", err)
	}

	// TODO: Publish MediaValidationCompletedEvent.
	return nil
}

// RequestProcessing creates a processing job for a media file.
func (s *MediaService) RequestProcessing(ctx context.Context, tenantID, mediaID string, processingType mediav1.ProcessingType, priority int32) (string, error) {
	if strings.TrimSpace(tenantID) == "" {
		return "", fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if mediaID == "" {
		return "", fmt.Errorf("%w: media_id is required", ErrInvalidInput)
	}

	_, err := s.mediaRepo.GetByID(ctx, tenantID, mediaID)
	if err != nil {
		if errors.Is(err, repository.ErrMediaNotFound) {
			return "", ErrMediaNotFound
		}
		return "", fmt.Errorf("failed to get media file: %w", err)
	}

	if priority < 1 || priority > 10 {
		priority = 5
	}

	job := &mediav1.ProcessingJob{
		Id:             uuid.New().String(),
		MediaId:        mediaID,
		ProcessingType: processingType,
		Status:         mediav1.ProcessingStatus_PROCESSING_STATUS_PENDING,
		Priority:       priority,
		RetryCount:     0,
		MaxRetries:     3,
		AuditInfo:      &commonv1.AuditInfo{},
	}

	createdJob, err := s.jobRepo.Create(ctx, job)
	if err != nil {
		return "", fmt.Errorf("failed to create processing job: %w", err)
	}

	// TODO: Publish MediaProcessingJobCreatedEvent.
	return createdJob.Id, nil
}

// GetProcessingJob retrieves a processing job by ID.
func (s *MediaService) GetProcessingJob(ctx context.Context, tenantID, jobID string) (*mediav1.ProcessingJob, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if jobID == "" {
		return nil, fmt.Errorf("%w: job_id is required", ErrInvalidInput)
	}

	job, err := s.jobRepo.GetByID(ctx, tenantID, jobID)
	if err != nil {
		if errors.Is(err, repository.ErrJobNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, fmt.Errorf("failed to get processing job: %w", err)
	}

	return job, nil
}

// ListProcessingJobs retrieves processing jobs with optional filters.
func (s *MediaService) ListProcessingJobs(
	ctx context.Context,
	tenantID, mediaID string,
	processingType *mediav1.ProcessingType,
	status *mediav1.ProcessingStatus,
	page, pageSize int,
) ([]*mediav1.ProcessingJob, int, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, 0, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	jobs, total, err := s.jobRepo.List(ctx, tenantID, mediaID, processingType, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list processing jobs: %w", err)
	}
	return jobs, total, nil
}

// DeleteMedia hard deletes a media file.
func (s *MediaService) DeleteMedia(ctx context.Context, tenantID, mediaID string) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if mediaID == "" {
		return fmt.Errorf("%w: media_id is required", ErrInvalidInput)
	}

	_, err := s.mediaRepo.GetByID(ctx, tenantID, mediaID)
	if err != nil {
		if errors.Is(err, repository.ErrMediaNotFound) {
			return ErrMediaNotFound
		}
		return fmt.Errorf("failed to get media file: %w", err)
	}

	if err := s.mediaRepo.Delete(ctx, tenantID, mediaID); err != nil {
		return fmt.Errorf("failed to delete media file: %w", err)
	}

	// Publish file deleted event (fire-and-forget in goroutine)
	if s.kafkaPublisher != nil {
		go s.publishFileDeletedAsync(ctx, mediaID, tenantID)
	}

	return nil
}

// UpdateOCRText updates the OCR extracted text for a media file.
func (s *MediaService) UpdateOCRText(ctx context.Context, tenantID, mediaID, ocrText string) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if mediaID == "" {
		return fmt.Errorf("%w: media_id is required", ErrInvalidInput)
	}

	if err := s.mediaRepo.UpdateOCRText(ctx, tenantID, mediaID, ocrText); err != nil {
		if errors.Is(err, repository.ErrMediaNotFound) {
			return ErrMediaNotFound
		}
		return fmt.Errorf("failed to update OCR text: %w", err)
	}

	// TODO: Publish MediaOCRCompletedEvent.
	return nil
}

// UpdateVirusScanStatus updates the virus scan status.
func (s *MediaService) UpdateVirusScanStatus(ctx context.Context, tenantID, mediaID string, status mediav1.VirusScanStatus) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if mediaID == "" {
		return fmt.Errorf("%w: media_id is required", ErrInvalidInput)
	}

	if err := s.mediaRepo.UpdateVirusScanStatus(ctx, tenantID, mediaID, status); err != nil {
		if errors.Is(err, repository.ErrMediaNotFound) {
			return ErrMediaNotFound
		}
		return fmt.Errorf("failed to update virus scan status: %w", err)
	}

	// TODO: Publish MediaVirusScanCompletedEvent.
	return nil
}

// UpdateProcessedFiles updates the optimized and thumbnail file references.
func (s *MediaService) UpdateProcessedFiles(ctx context.Context, tenantID, mediaID, optimizedFileID, thumbnailFileID string) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if mediaID == "" {
		return fmt.Errorf("%w: media_id is required", ErrInvalidInput)
	}

	if err := s.mediaRepo.UpdateProcessedFiles(ctx, tenantID, mediaID, optimizedFileID, thumbnailFileID); err != nil {
		if errors.Is(err, repository.ErrMediaNotFound) {
			return ErrMediaNotFound
		}
		return fmt.Errorf("failed to update processed files: %w", err)
	}

	return nil
}

// ResolveFileDownloadURL resolves a secure download URL via storage service, with CDN fallback.
func (s *MediaService) ResolveFileDownloadURL(ctx context.Context, tenantID, fileID string) (string, int64, error) {
	if strings.TrimSpace(fileID) == "" {
		return "", 0, fmt.Errorf("%w: file_id is required", ErrInvalidInput)
	}

	expiresInSeconds := int64(s.defaultURLExpiryMinute) * 60
	if s.storageDownloadClient != nil && strings.TrimSpace(tenantID) != "" {
		resp, err := s.storageDownloadClient.GetDownloadURL(ctx, &storageservicev1.GetDownloadURLRequest{
			TenantId:         tenantID,
			FileId:           fileID,
			ExpiresInMinutes: s.defaultURLExpiryMinute,
		})
		if err == nil && strings.TrimSpace(resp.GetDownloadUrl()) != "" {
			if resp.ExpiresAt != nil {
				if secs := int64(time.Until(resp.ExpiresAt.AsTime()).Seconds()); secs > 0 {
					expiresInSeconds = secs
				}
			}
			return resp.DownloadUrl, expiresInSeconds, nil
		}
	}

	if s.fallbackDownloadBase == "" {
		return "", 0, fmt.Errorf("storage download URL resolution unavailable")
	}
	return fmt.Sprintf("%s/%s", s.fallbackDownloadBase, url.PathEscape(fileID)), expiresInSeconds, nil
}

// kafkaPublisherInterface defines the Kafka publisher methods needed by MediaService
type kafkaPublisherInterface interface {
	PublishFileUploaded(ctx context.Context, mediaID, tenantID, filename, mimeType string, sizeBytes int64) error
	PublishFileDeleted(ctx context.Context, mediaID, tenantID string) error
}

// publishFileUploadedAsync publishes a file uploaded event asynchronously
func (s *MediaService) publishFileUploadedAsync(ctx context.Context, mediaID, tenantID, fileID, mimeType string, sizeBytes int64) {
	pub, ok := s.kafkaPublisher.(kafkaPublisherInterface)
	if !ok || pub == nil {
		return
	}
	_ = pub.PublishFileUploaded(ctx, mediaID, tenantID, fileID, mimeType, sizeBytes)
}

// publishFileDeletedAsync publishes a file deleted event asynchronously
func (s *MediaService) publishFileDeletedAsync(ctx context.Context, mediaID, tenantID string) {
	pub, ok := s.kafkaPublisher.(kafkaPublisherInterface)
	if !ok || pub == nil {
		return
	}
	_ = pub.PublishFileDeleted(ctx, mediaID, tenantID)
}

// enqueueProcessingJobs enqueues processing jobs based on media MIME type
func (s *MediaService) enqueueProcessingJobs(ctx context.Context, mediaID, tenantID, mimeType string) {
	// Check if MIME type is an image
	isImage := strings.HasPrefix(mimeType, "image/")
	isPDF := mimeType == "application/pdf"

	// Enqueue THUMBNAIL job for images
	if isImage {
		thumbnailJob := &mediav1.ProcessingJob{
			Id:             uuid.New().String(),
			MediaId:        mediaID,
			ProcessingType: mediav1.ProcessingType_PROCESSING_TYPE_THUMBNAIL,
			Status:         mediav1.ProcessingStatus_PROCESSING_STATUS_PENDING,
			Priority:       5,
			RetryCount:     0,
			MaxRetries:     3,
			AuditInfo:      &commonv1.AuditInfo{},
		}
		_, _ = s.jobRepo.Create(ctx, thumbnailJob)
	}

	// Enqueue OPTIMIZATION job for images
	if isImage {
		optimizationJob := &mediav1.ProcessingJob{
			Id:             uuid.New().String(),
			MediaId:        mediaID,
			ProcessingType: mediav1.ProcessingType_PROCESSING_TYPE_OPTIMIZATION,
			Status:         mediav1.ProcessingStatus_PROCESSING_STATUS_PENDING,
			Priority:       5,
			RetryCount:     0,
			MaxRetries:     3,
			AuditInfo:      &commonv1.AuditInfo{},
		}
		_, _ = s.jobRepo.Create(ctx, optimizationJob)
	}

	// Enqueue VIRUS_SCAN job for all files
	virusScanJob := &mediav1.ProcessingJob{
		Id:             uuid.New().String(),
		MediaId:        mediaID,
		ProcessingType: mediav1.ProcessingType_PROCESSING_TYPE_VIRUS_SCAN,
		Status:         mediav1.ProcessingStatus_PROCESSING_STATUS_PENDING,
		Priority:       5,
		RetryCount:     0,
		MaxRetries:     3,
		AuditInfo:      &commonv1.AuditInfo{},
	}
	_, _ = s.jobRepo.Create(ctx, virusScanJob)

	// Enqueue OCR job for PDFs and images
	if isPDF || isImage {
		ocrJob := &mediav1.ProcessingJob{
			Id:             uuid.New().String(),
			MediaId:        mediaID,
			ProcessingType: mediav1.ProcessingType_PROCESSING_TYPE_OCR,
			Status:         mediav1.ProcessingStatus_PROCESSING_STATUS_PENDING,
			Priority:       5,
			RetryCount:     0,
			MaxRetries:     3,
			AuditInfo:      &commonv1.AuditInfo{},
		}
		_, _ = s.jobRepo.Create(ctx, ocrJob)
	}
}
