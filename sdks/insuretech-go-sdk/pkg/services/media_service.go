package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// MediaService handles media-related API calls
type MediaService struct {
	Client Client
}

// ListMedia List media files for entity
func (s *MediaService) ListMedia(ctx context.Context, entityType string, entityId string) error {
	path := "/v1/entities/{entity_type}/{entity_id}/media"
	path = strings.ReplaceAll(path, "{entity_type}", entityType)
	path = strings.ReplaceAll(path, "{entity_id}", entityId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UploadMedia Upload media file
func (s *MediaService) UploadMedia(ctx context.Context, req *models.MediaUploadRequest) error {
	path := "/v1/media"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetMedia Get media file
func (s *MediaService) GetMedia(ctx context.Context, mediaId string) error {
	path := "/v1/media/{media_id}"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// DeleteMedia Delete media file
func (s *MediaService) DeleteMedia(ctx context.Context, mediaId string) error {
	path := "/v1/media/{media_id}"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// DownloadMedia Download media file
func (s *MediaService) DownloadMedia(ctx context.Context, mediaId string) error {
	path := "/v1/media/{media_id}/download"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// DownloadOptimized Download optimized version
func (s *MediaService) DownloadOptimized(ctx context.Context, mediaId string) error {
	path := "/v1/media/{media_id}/optimized"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RequestProcessing Request processing (OCR, optimization, etc
func (s *MediaService) RequestProcessing(ctx context.Context, mediaId string, req *models.RequestProcessingRequest) error {
	path := "/v1/media/{media_id}/process"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// DownloadThumbnail Download thumbnail
func (s *MediaService) DownloadThumbnail(ctx context.Context, mediaId string) error {
	path := "/v1/media/{media_id}/thumbnail"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ValidateMedia Validate media file
func (s *MediaService) ValidateMedia(ctx context.Context, mediaId string, req *models.MediaValidationRequest) error {
	path := "/v1/media/{media_id}:validate"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListProcessingJobs List processing jobs
func (s *MediaService) ListProcessingJobs(ctx context.Context) error {
	path := "/v1/processing-jobs"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetProcessingJob Get processing job status
func (s *MediaService) GetProcessingJob(ctx context.Context, jobId string) error {
	path := "/v1/processing-jobs/{job_id}"
	path = strings.ReplaceAll(path, "{job_id}", jobId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

