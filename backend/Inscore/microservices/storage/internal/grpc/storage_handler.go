package server

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/service"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
)

// StorageHandler implements the StorageService gRPC server
type StorageHandler struct {
	storageservicev1.UnimplementedStorageServiceServer
	storageService storageServiceIface
}

// NewStorageHandler creates a new storage handler
func NewStorageHandler(storageService storageServiceIface) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
	}
}

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

// UploadFile uploads a file to storage
func (h *StorageHandler) UploadFile(ctx context.Context, req *storageservicev1.UploadFileRequest) (*storageservicev1.UploadFileResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	tenantID := req.TenantId

	// Validate request
	if len(req.Content) == 0 {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	if req.Filename == "" {
		return nil, status.Error(codes.InvalidArgument, "filename is required")
	}
	if req.ContentType == "" {
		req.ContentType = "application/octet-stream"
	}

	// Get uploaded_by from context (from auth service)
	uploadedBy := actorFromContext(ctx, tenantID)

	// Upload file
	file, err := h.storageService.UploadFile(
		ctx,
		tenantID,
		req.Content,
		req.Filename,
		req.ContentType,
		req.FileType,
		req.ReferenceId,
		req.ReferenceType,
		req.IsPublic,
		req.ExpiresAt,
		uploadedBy,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to upload file: %v", err)
	}

	return &storageservicev1.UploadFileResponse{
		File: file,
	}, nil
}

// UploadFiles uploads multiple files
func (h *StorageHandler) UploadFiles(ctx context.Context, req *storageservicev1.UploadFilesRequest) (*storageservicev1.UploadFilesResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	if len(req.Files) == 0 {
		return nil, status.Error(codes.InvalidArgument, "files are required")
	}
	tenantID := req.TenantId

	uploadedBy := actorFromContext(ctx, tenantID)
	input := make([]service.UploadFileInput, 0, len(req.Files))
	for i, fileUpload := range req.Files {
		if len(fileUpload.Content) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "files[%d].content is required", i)
		}
		if strings.TrimSpace(fileUpload.Filename) == "" {
			return nil, status.Errorf(codes.InvalidArgument, "files[%d].filename is required", i)
		}
		input = append(input, service.UploadFileInput{
			Content:       fileUpload.Content,
			Filename:      fileUpload.Filename,
			ContentType:   fileUpload.ContentType,
			FileType:      fileUpload.FileType,
			ReferenceID:   fileUpload.ReferenceId,
			ReferenceType: fileUpload.ReferenceType,
			IsPublic:      fileUpload.IsPublic,
		})
	}

	uploadedFiles, err := h.storageService.UploadFiles(ctx, tenantID, input, uploadedBy)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to upload files: %v", err)
	}

	return &storageservicev1.UploadFilesResponse{
		Files: uploadedFiles,
	}, nil
}

// GetFile retrieves file metadata
func (h *StorageHandler) GetFile(ctx context.Context, req *storageservicev1.GetFileRequest) (*storageservicev1.GetFileResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	if strings.TrimSpace(req.FileId) == "" {
		return nil, status.Error(codes.InvalidArgument, "file_id is required")
	}
	tenantID := req.TenantId

	file, err := h.storageService.GetFile(ctx, tenantID, req.FileId)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, service.ErrFileNotFound):
			return nil, status.Error(codes.NotFound, "file not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get file: %v", err)
		}
	}

	return &storageservicev1.GetFileResponse{
		File: file,
	}, nil
}

// UpdateFile updates mutable metadata fields for an existing file.
func (h *StorageHandler) UpdateFile(ctx context.Context, req *storageservicev1.UpdateFileRequest) (*storageservicev1.UpdateFileResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	if req.FileId == "" {
		return nil, status.Error(codes.InvalidArgument, "file_id is required")
	}

	tenantID := req.TenantId
	uploadedBy := actorFromContext(ctx, tenantID)
	input := &service.UpdateFileInput{
		TenantID:     tenantID,
		FileID:       req.FileId,
		ClearExpires: req.ClearExpiresAt,
		UpdatedBy:    uploadedBy,
	}

	if req.Filename != nil {
		filename := req.GetFilename()
		input.Filename = &filename
	}
	if req.ContentType != nil {
		contentType := req.GetContentType()
		input.ContentType = &contentType
	}
	if req.FileType != nil {
		ft := req.GetFileType()
		input.FileType = &ft
	}
	if req.ReferenceId != nil {
		referenceID := req.GetReferenceId()
		input.ReferenceID = &referenceID
	}
	if req.ReferenceType != nil {
		referenceType := req.GetReferenceType()
		input.ReferenceType = &referenceType
	}
	if req.IsPublic != nil {
		isPublic := req.GetIsPublic()
		input.IsPublic = &isPublic
	}
	if req.ExpiresAt != nil {
		input.ExpiresAt = req.ExpiresAt
	}

	file, err := h.storageService.UpdateFileMetadata(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, service.ErrNoMetadataUpdates):
			return nil, status.Error(codes.InvalidArgument, "no metadata fields provided")
		case errors.Is(err, service.ErrFileNotFound):
			return nil, status.Error(codes.NotFound, "file not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to update file metadata: %v", err)
		}
	}

	return &storageservicev1.UpdateFileResponse{File: file}, nil
}

// GetUploadURL generates a presigned URL for direct upload
func (h *StorageHandler) GetUploadURL(ctx context.Context, req *storageservicev1.GetUploadURLRequest) (*storageservicev1.GetUploadURLResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	if strings.TrimSpace(req.Filename) == "" {
		return nil, status.Error(codes.InvalidArgument, "filename is required")
	}
	if req.ContentType == "" {
		req.ContentType = "application/octet-stream"
	}
	tenantID := req.TenantId

	uploadedBy := actorFromContext(ctx, tenantID)

	uploadURL, fileID, s3Key, err := h.storageService.GetPresignedUploadURL(
		ctx,
		tenantID,
		req.Filename,
		req.ContentType,
		req.FileType,
		req.ReferenceId,
		req.ReferenceType,
		req.IsPublic,
		req.ExpiresInMinutes,
		uploadedBy,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate upload URL: %v", err)
	}

	return &storageservicev1.GetUploadURLResponse{
		UploadUrl:  uploadURL,
		FileId:     fileID,
		StorageKey: s3Key, // Fix: struct field is StorageKey
	}, nil
}

// FinalizeUpload verifies direct upload and persists final metadata.
func (h *StorageHandler) FinalizeUpload(ctx context.Context, req *storageservicev1.FinalizeUploadRequest) (*storageservicev1.FinalizeUploadResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	if req.FileId == "" {
		return nil, status.Error(codes.InvalidArgument, "file_id is required")
	}
	tenantID := req.TenantId

	uploadedBy := actorFromContext(ctx, tenantID)
	file, err := h.storageService.FinalizeDirectUpload(
		ctx,
		tenantID,
		req.FileId,
		req.Filename,
		req.ContentType,
		req.FileType,
		req.ReferenceId,
		req.ReferenceType,
		req.IsPublic,
		req.ExpiresAt,
		uploadedBy,
	)
	if err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			return nil, status.Error(codes.NotFound, "file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to finalize upload: %v", err)
	}

	return &storageservicev1.FinalizeUploadResponse{File: file}, nil
}

// GetDownloadURL generates a presigned URL for download
func (h *StorageHandler) GetDownloadURL(ctx context.Context, req *storageservicev1.GetDownloadURLRequest) (*storageservicev1.GetDownloadURLResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	tenantID := req.TenantId

	// Fix: FileId is string
	downloadURL, expiresAt, err := h.storageService.GetPresignedDownloadURL(
		ctx,
		tenantID,
		req.FileId,
		req.ExpiresInMinutes,
	)
	if err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			return nil, status.Error(codes.NotFound, "file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to generate download URL: %v", err)
	}

	return &storageservicev1.GetDownloadURLResponse{
		DownloadUrl: downloadURL,
		ExpiresAt:   expiresAt,
	}, nil
}

// DeleteFile deletes a file
func (h *StorageHandler) DeleteFile(ctx context.Context, req *storageservicev1.DeleteFileRequest) (*storageservicev1.DeleteFileResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	if strings.TrimSpace(req.FileId) == "" {
		return nil, status.Error(codes.InvalidArgument, "file_id is required")
	}
	tenantID := req.TenantId

	deletedBy := actorFromContext(ctx, tenantID)
	if err := h.storageService.DeleteFile(ctx, tenantID, req.FileId, deletedBy); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, service.ErrFileNotFound):
			return nil, status.Error(codes.NotFound, "file not found")
		case errors.Is(err, service.ErrStorageUnavailable):
			return nil, status.Error(codes.Unavailable, "storage backend unavailable")
		default:
			return nil, status.Errorf(codes.Internal, "failed to delete file: %v", err)
		}
	}

	return &storageservicev1.DeleteFileResponse{
		Success: true,
	}, nil
}

// ListFiles lists files with filters
func (h *StorageHandler) ListFiles(ctx context.Context, req *storageservicev1.ListFilesRequest) (*storageservicev1.ListFilesResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}
	tenantID := req.TenantId

	// Parse pagination
	limit := int32(50)
	offset := int32(0)
	currentPage := int32(1)
	if req.Page != nil {
		if req.Page.Page < 0 {
			return nil, status.Error(codes.InvalidArgument, "page must be >= 0")
		}
		if req.Page.PageSize < 0 {
			return nil, status.Error(codes.InvalidArgument, "page_size must be >= 0")
		}
		if req.Page.Page > 0 {
			currentPage = req.Page.Page
		}
		if req.Page.PageSize > 0 {
			limit = req.Page.PageSize
		}
		offset = (currentPage - 1) * limit
	}

	files, total, err := h.storageService.ListFiles(
		ctx,
		tenantID,
		req.FileType,
		req.ReferenceId,
		req.ReferenceType,
		req.UploadedBy,
		limit,
		offset,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to list files: %v", err)
	}
	totalPages := int32(0)
	if total > 0 && limit > 0 {
		totalPages = int32((total + int(limit) - 1) / int(limit))
	}
	hasNext := totalPages > 0 && currentPage < totalPages
	hasPrevious := totalPages > 0 && currentPage > 1

	return &storageservicev1.ListFilesResponse{
		Files: files,
		Page: &commonv1.PaginationResponse{
			TotalItems:  int32(total),
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			PageSize:    limit,
			HasNext:     hasNext,
			HasPrevious: hasPrevious,
		},
	}, nil
}
