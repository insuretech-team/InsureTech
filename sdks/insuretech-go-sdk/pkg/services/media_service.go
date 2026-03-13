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

// GetMedia Get media file
func (s *MediaService) GetMedia(ctx context.Context, mediaId string) (*models.MediaRetrievalResponse, error) {
	path := "/v1/media/{media_id}"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	var result models.MediaRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteMedia Delete media file
func (s *MediaService) DeleteMedia(ctx context.Context, mediaId string) error {
	path := "/v1/media/{media_id}"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// ListProcessingJobs List processing jobs
func (s *MediaService) ListProcessingJobs(ctx context.Context) (*models.ProcessingJobsListingResponse, error) {
	path := "/v1/processing-jobs"
	var result models.ProcessingJobsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadMedia Download media file
func (s *MediaService) DownloadMedia(ctx context.Context, mediaId string) (*models.MediaDownloadResponse, error) {
	path := "/v1/media/{media_id}/download"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	var result models.MediaDownloadResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RequestProcessing Request processing (OCR, optimization, etc
func (s *MediaService) RequestProcessing(ctx context.Context, mediaId string, req *models.RequestProcessingRequest) (*models.RequestProcessingResponse, error) {
	path := "/v1/media/{media_id}/process"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	var result models.RequestProcessingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProcessingJob Get processing job status
func (s *MediaService) GetProcessingJob(ctx context.Context, jobId string) (*models.ProcessingJobRetrievalResponse, error) {
	path := "/v1/processing-jobs/{job_id}"
	path = strings.ReplaceAll(path, "{job_id}", jobId)
	var result models.ProcessingJobRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadThumbnail Download thumbnail
func (s *MediaService) DownloadThumbnail(ctx context.Context, mediaId string) (*models.ThumbnailDownloadResponse, error) {
	path := "/v1/media/{media_id}/thumbnail"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	var result models.ThumbnailDownloadResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ValidateMedia Validate media file
func (s *MediaService) ValidateMedia(ctx context.Context, mediaId string, req *models.MediaValidationRequest) (*models.MediaValidationResponse, error) {
	path := "/v1/media/{media_id}:validate"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	var result models.MediaValidationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMedia List media files for entity
func (s *MediaService) ListMedia(ctx context.Context, entityType string, entityId string) (*models.MediaListingResponse, error) {
	path := "/v1/entities/{entity_type}/{entity_id}/media"
	path = strings.ReplaceAll(path, "{entity_type}", entityType)
	path = strings.ReplaceAll(path, "{entity_id}", entityId)
	var result models.MediaListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadMedia Upload media file
func (s *MediaService) UploadMedia(ctx context.Context, req *models.MediaUploadRequest) (*models.MediaUploadResponse, error) {
	path := "/v1/media"
	var result models.MediaUploadResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadOptimized Download optimized version
func (s *MediaService) DownloadOptimized(ctx context.Context, mediaId string) (*models.OptimizedDownloadResponse, error) {
	path := "/v1/media/{media_id}/optimized"
	path = strings.ReplaceAll(path, "{media_id}", mediaId)
	var result models.OptimizedDownloadResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

