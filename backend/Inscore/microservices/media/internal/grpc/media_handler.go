package server

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/media/internal/service"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	mediav1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/entity/v1"
	mediaservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/services/v1"
)

// MediaHandler implements the MediaService gRPC server.
type MediaHandler struct {
	mediaservicev1.UnimplementedMediaServiceServer
	mediaService *service.MediaService
}

// NewMediaHandler creates a new media handler.
func NewMediaHandler(mediaService *service.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

// actorFromContext extracts the user/actor ID from gRPC metadata.
func actorFromContext(ctx context.Context, fallback string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fallback
	}
	keys := []string{"x-user-id", "x-actor-id", "x-sub", "user-id"}
	for _, k := range keys {
		vals := md.Get(k)
		if len(vals) == 0 {
			continue
		}
		v := strings.TrimSpace(vals[0])
		if v != "" {
			return v
		}
	}
	return fallback
}

func tenantFromContext(ctx context.Context, fallback string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fallback
	}
	keys := []string{"x-tenant-id", "tenant-id", "x-tenant"}
	for _, k := range keys {
		vals := md.Get(k)
		if len(vals) == 0 {
			continue
		}
		v := strings.TrimSpace(vals[0])
		if v != "" {
			return v
		}
	}
	return fallback
}

func invalidArg(msg string) *commonv1.Error {
	return &commonv1.Error{Code: "INVALID_ARGUMENT", Message: msg}
}

func internalErr(msg string) *commonv1.Error {
	return &commonv1.Error{Code: "INTERNAL", Message: msg}
}

func mapError(err error, notFoundMsg, internalMsg string) *commonv1.Error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return invalidArg(err.Error())
	case errors.Is(err, service.ErrMediaNotFound), errors.Is(err, service.ErrJobNotFound):
		return &commonv1.Error{Code: "NOT_FOUND", Message: notFoundMsg}
	default:
		return internalErr(internalMsg)
	}
}

func parseMediaType(input string) (*mediav1.MediaType, error) {
	t := parseEnumWithPrefix(input, "MEDIA_TYPE_")
	if t == "" {
		return nil, nil
	}
	v, ok := mediav1.MediaType_value[t]
	if !ok {
		return nil, errors.New("invalid media_type")
	}
	mt := mediav1.MediaType(v)
	return &mt, nil
}

func parseValidationStatus(input string) (*mediav1.ValidationStatus, error) {
	t := parseEnumWithPrefix(input, "VALIDATION_STATUS_")
	if t == "" {
		return nil, nil
	}
	v, ok := mediav1.ValidationStatus_value[t]
	if !ok {
		return nil, errors.New("invalid validation_status")
	}
	vs := mediav1.ValidationStatus(v)
	return &vs, nil
}

func parseProcessingType(input string) (*mediav1.ProcessingType, error) {
	t := parseEnumWithPrefix(input, "PROCESSING_TYPE_")
	if t == "" {
		return nil, nil
	}
	v, ok := mediav1.ProcessingType_value[t]
	if !ok {
		return nil, errors.New("invalid processing_type")
	}
	pt := mediav1.ProcessingType(v)
	return &pt, nil
}

func parseProcessingStatus(input string) (*mediav1.ProcessingStatus, error) {
	t := parseEnumWithPrefix(input, "PROCESSING_STATUS_")
	if t == "" {
		return nil, nil
	}
	v, ok := mediav1.ProcessingStatus_value[t]
	if !ok {
		return nil, errors.New("invalid status")
	}
	ps := mediav1.ProcessingStatus(v)
	return &ps, nil
}

func parseEnumWithPrefix(input, prefix string) string {
	t := strings.TrimSpace(strings.ToUpper(input))
	if t == "" {
		return ""
	}
	if strings.HasPrefix(t, prefix) {
		return t
	}
	return prefix + t
}

// UploadMedia handles media file upload.
func (h *MediaHandler) UploadMedia(ctx context.Context, req *mediaservicev1.UploadMediaRequest) (*mediaservicev1.UploadMediaResponse, error) {
	if req.FileId == "" {
		return &mediaservicev1.UploadMediaResponse{Error: invalidArg("file_id is required")}, nil
	}
	if req.MimeType == "" {
		return &mediaservicev1.UploadMediaResponse{Error: invalidArg("mime_type is required")}, nil
	}

	tenantID := strings.TrimSpace(req.TenantId)
	if tenantID == "" {
		tenantID = tenantFromContext(ctx, "")
	}
	if tenantID == "" {
		return &mediaservicev1.UploadMediaResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	uploadedBy := strings.TrimSpace(req.UploadedBy)
	if uploadedBy == "" {
		uploadedBy = actorFromContext(ctx, tenantID)
	}

	mediaType := mediav1.MediaType_MEDIA_TYPE_UNSPECIFIED
	if strings.TrimSpace(req.MediaType) != "" {
		parsed, err := parseMediaType(req.MediaType)
		if err != nil || parsed == nil {
			return &mediaservicev1.UploadMediaResponse{Error: invalidArg("invalid media_type")}, nil
		}
		mediaType = *parsed
	}

	media, err := h.mediaService.UploadMedia(ctx, &service.UploadMediaInput{
		FileID:        req.FileId,
		TenantID:      tenantID,
		MediaType:     mediaType,
		MimeType:      req.MimeType,
		FileSizeBytes: req.FileSizeBytes,
		EntityType:    req.EntityType,
		EntityID:      req.EntityId,
		UploadedBy:    uploadedBy,
	})
	if err != nil {
		return &mediaservicev1.UploadMediaResponse{Error: mapError(err, "media file not found", "failed to upload media")}, nil
	}

	if req.AutoValidate {
		_ = h.mediaService.ValidateMedia(ctx, tenantID, media.Id, nil)
	}
	if req.AutoOptimize {
		_, _ = h.mediaService.RequestProcessing(ctx, tenantID, media.Id, mediav1.ProcessingType_PROCESSING_TYPE_OPTIMIZATION, 5)
	}
	if req.AutoThumbnail {
		_, _ = h.mediaService.RequestProcessing(ctx, tenantID, media.Id, mediav1.ProcessingType_PROCESSING_TYPE_THUMBNAIL, 5)
	}

	return &mediaservicev1.UploadMediaResponse{MediaId: media.Id, Message: "Media uploaded successfully"}, nil
}

// GetMedia retrieves media file metadata.
func (h *MediaHandler) GetMedia(ctx context.Context, req *mediaservicev1.GetMediaRequest) (*mediaservicev1.GetMediaResponse, error) {
	if req.MediaId == "" {
		return &mediaservicev1.GetMediaResponse{Error: invalidArg("media_id is required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.GetMediaResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	media, err := h.mediaService.GetMedia(ctx, tenantID, req.MediaId)
	if err != nil {
		return &mediaservicev1.GetMediaResponse{Error: mapError(err, "media file not found", "failed to get media file")}, nil
	}

	return &mediaservicev1.GetMediaResponse{Media: media}, nil
}

// ListMedia lists media files for an entity.
func (h *MediaHandler) ListMedia(ctx context.Context, req *mediaservicev1.ListMediaRequest) (*mediaservicev1.ListMediaResponse, error) {
	if req.EntityType == "" || req.EntityId == "" {
		return &mediaservicev1.ListMediaResponse{Error: invalidArg("entity_type and entity_id are required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.ListMediaResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	mediaType, err := parseMediaType(req.MediaType)
	if err != nil {
		return &mediaservicev1.ListMediaResponse{Error: invalidArg("invalid media_type")}, nil
	}
	validationStatus, err := parseValidationStatus(req.ValidationStatus)
	if err != nil {
		return &mediaservicev1.ListMediaResponse{Error: invalidArg("invalid validation_status")}, nil
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	mediaFiles, total, err := h.mediaService.ListMediaByEntity(
		ctx,
		tenantID,
		req.EntityType,
		req.EntityId,
		mediaType,
		validationStatus,
		int(page),
		int(pageSize),
	)
	if err != nil {
		return &mediaservicev1.ListMediaResponse{Error: mapError(err, "media files not found", "failed to list media files")}, nil
	}

	return &mediaservicev1.ListMediaResponse{MediaFiles: mediaFiles, TotalCount: int32(total)}, nil
}

// DownloadMedia generates a download URL for a media file.
func (h *MediaHandler) DownloadMedia(ctx context.Context, req *mediaservicev1.DownloadMediaRequest) (*mediaservicev1.DownloadMediaResponse, error) {
	if req.MediaId == "" {
		return &mediaservicev1.DownloadMediaResponse{Error: invalidArg("media_id is required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.DownloadMediaResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	media, err := h.mediaService.GetMedia(ctx, tenantID, req.MediaId)
	if err != nil {
		return &mediaservicev1.DownloadMediaResponse{Error: mapError(err, "media file not found", "failed to get media file")}, nil
	}

	downloadURL, expiresIn, err := h.mediaService.ResolveFileDownloadURL(ctx, tenantID, media.FileId)
	if err != nil {
		return &mediaservicev1.DownloadMediaResponse{Error: internalErr("failed to generate download URL")}, nil
	}
	return &mediaservicev1.DownloadMediaResponse{DownloadUrl: downloadURL, ExpiresInSeconds: expiresIn}, nil
}

// DownloadOptimized generates a download URL for optimized version.
func (h *MediaHandler) DownloadOptimized(ctx context.Context, req *mediaservicev1.DownloadOptimizedRequest) (*mediaservicev1.DownloadOptimizedResponse, error) {
	if req.MediaId == "" {
		return &mediaservicev1.DownloadOptimizedResponse{Error: invalidArg("media_id is required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.DownloadOptimizedResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	media, err := h.mediaService.GetMedia(ctx, tenantID, req.MediaId)
	if err != nil {
		return &mediaservicev1.DownloadOptimizedResponse{Error: mapError(err, "media file not found", "failed to get media file")}, nil
	}
	if strings.TrimSpace(media.OptimizedFileId) == "" {
		return &mediaservicev1.DownloadOptimizedResponse{Error: &commonv1.Error{Code: "NOT_FOUND", Message: "optimized version not available"}}, nil
	}

	downloadURL, expiresIn, err := h.mediaService.ResolveFileDownloadURL(ctx, tenantID, media.OptimizedFileId)
	if err != nil {
		return &mediaservicev1.DownloadOptimizedResponse{Error: internalErr("failed to generate download URL")}, nil
	}
	return &mediaservicev1.DownloadOptimizedResponse{DownloadUrl: downloadURL, ExpiresInSeconds: expiresIn}, nil
}

// DownloadThumbnail generates a download URL for thumbnail.
func (h *MediaHandler) DownloadThumbnail(ctx context.Context, req *mediaservicev1.DownloadThumbnailRequest) (*mediaservicev1.DownloadThumbnailResponse, error) {
	if req.MediaId == "" {
		return &mediaservicev1.DownloadThumbnailResponse{Error: invalidArg("media_id is required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.DownloadThumbnailResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	media, err := h.mediaService.GetMedia(ctx, tenantID, req.MediaId)
	if err != nil {
		return &mediaservicev1.DownloadThumbnailResponse{Error: mapError(err, "media file not found", "failed to get media file")}, nil
	}
	if strings.TrimSpace(media.ThumbnailFileId) == "" {
		return &mediaservicev1.DownloadThumbnailResponse{Error: &commonv1.Error{Code: "NOT_FOUND", Message: "thumbnail not available"}}, nil
	}

	downloadURL, expiresIn, err := h.mediaService.ResolveFileDownloadURL(ctx, tenantID, media.ThumbnailFileId)
	if err != nil {
		return &mediaservicev1.DownloadThumbnailResponse{Error: internalErr("failed to generate download URL")}, nil
	}
	return &mediaservicev1.DownloadThumbnailResponse{DownloadUrl: downloadURL, ExpiresInSeconds: expiresIn}, nil
}

// DeleteMedia deletes a media file.
func (h *MediaHandler) DeleteMedia(ctx context.Context, req *mediaservicev1.DeleteMediaRequest) (*mediaservicev1.DeleteMediaResponse, error) {
	if req.MediaId == "" {
		return &mediaservicev1.DeleteMediaResponse{Error: invalidArg("media_id is required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.DeleteMediaResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	if err := h.mediaService.DeleteMedia(ctx, tenantID, req.MediaId); err != nil {
		return &mediaservicev1.DeleteMediaResponse{Error: mapError(err, "media file not found", "failed to delete media file")}, nil
	}

	return &mediaservicev1.DeleteMediaResponse{Message: "Media file deleted successfully"}, nil
}

// ValidateMedia validates a media file.
func (h *MediaHandler) ValidateMedia(ctx context.Context, req *mediaservicev1.ValidateMediaRequest) (*mediaservicev1.ValidateMediaResponse, error) {
	if req.MediaId == "" {
		return &mediaservicev1.ValidateMediaResponse{Error: invalidArg("media_id is required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.ValidateMediaResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	if err := h.mediaService.ValidateMedia(ctx, tenantID, req.MediaId, req.ValidationRules); err != nil {
		return &mediaservicev1.ValidateMediaResponse{Error: mapError(err, "media file not found", "failed to validate media")}, nil
	}

	return &mediaservicev1.ValidateMediaResponse{ValidationStatus: "VALIDATED"}, nil
}

// RequestProcessing creates a processing job.
func (h *MediaHandler) RequestProcessing(ctx context.Context, req *mediaservicev1.RequestProcessingRequest) (*mediaservicev1.RequestProcessingResponse, error) {
	if req.MediaId == "" {
		return &mediaservicev1.RequestProcessingResponse{Error: invalidArg("media_id is required")}, nil
	}
	if req.ProcessingType == "" {
		return &mediaservicev1.RequestProcessingResponse{Error: invalidArg("processing_type is required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.RequestProcessingResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	parsedProcessingType, err := parseProcessingType(req.ProcessingType)
	if err != nil || parsedProcessingType == nil {
		return &mediaservicev1.RequestProcessingResponse{Error: invalidArg("invalid processing_type")}, nil
	}

	priority := req.Priority
	if priority < 1 || priority > 10 {
		priority = 5
	}

	jobID, err := h.mediaService.RequestProcessing(ctx, tenantID, req.MediaId, *parsedProcessingType, priority)
	if err != nil {
		return &mediaservicev1.RequestProcessingResponse{Error: mapError(err, "media file not found", "failed to create processing job")}, nil
	}

	return &mediaservicev1.RequestProcessingResponse{JobId: jobID, Status: "PENDING", Message: "Processing job created successfully"}, nil
}

// GetProcessingJob retrieves a processing job.
func (h *MediaHandler) GetProcessingJob(ctx context.Context, req *mediaservicev1.GetProcessingJobRequest) (*mediaservicev1.GetProcessingJobResponse, error) {
	if req.JobId == "" {
		return &mediaservicev1.GetProcessingJobResponse{Error: invalidArg("job_id is required")}, nil
	}
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.GetProcessingJobResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	job, err := h.mediaService.GetProcessingJob(ctx, tenantID, req.JobId)
	if err != nil {
		return &mediaservicev1.GetProcessingJobResponse{Error: mapError(err, "processing job not found", "failed to get processing job")}, nil
	}

	return &mediaservicev1.GetProcessingJobResponse{Job: job}, nil
}

// ListProcessingJobs lists processing jobs.
func (h *MediaHandler) ListProcessingJobs(ctx context.Context, req *mediaservicev1.ListProcessingJobsRequest) (*mediaservicev1.ListProcessingJobsResponse, error) {
	tenantID := tenantFromContext(ctx, "")
	if tenantID == "" {
		return &mediaservicev1.ListProcessingJobsResponse{Error: invalidArg("tenant_id is required")}, nil
	}

	processingType, err := parseProcessingType(req.ProcessingType)
	if err != nil {
		return &mediaservicev1.ListProcessingJobsResponse{Error: invalidArg("invalid processing_type")}, nil
	}
	processingStatus, err := parseProcessingStatus(req.Status)
	if err != nil {
		return &mediaservicev1.ListProcessingJobsResponse{Error: invalidArg("invalid status")}, nil
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	jobs, total, err := h.mediaService.ListProcessingJobs(
		ctx,
		tenantID,
		req.MediaId,
		processingType,
		processingStatus,
		int(page),
		int(pageSize),
	)
	if err != nil {
		return &mediaservicev1.ListProcessingJobsResponse{Error: mapError(err, "processing jobs not found", "failed to list processing jobs")}, nil
	}

	return &mediaservicev1.ListProcessingJobsResponse{Jobs: jobs, TotalCount: int32(total)}, nil
}
