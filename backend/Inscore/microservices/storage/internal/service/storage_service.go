package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/index"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/s3"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
)

// StorageService handles file storage operations
type StorageService struct {
	fileRepo       *repository.FileRepository
	s3Client       *s3.Client
	eventPublisher *events.Publisher
	userFileIndex  *index.UserFileIndex
}

var (
	ErrFileNotFound       = errors.New("file not found")
	ErrNoMetadataUpdates  = errors.New("no metadata fields provided")
	ErrInvalidInput       = errors.New("invalid input")
	ErrStorageUnavailable = errors.New("storage backend unavailable")
)

// UploadFileInput is a transport-agnostic input for uploads.
type UploadFileInput struct {
	Content       []byte
	Filename      string
	ContentType   string
	FileType      storageentityv1.FileType
	ReferenceID   string
	ReferenceType string
	IsPublic      bool
	ExpiresAt     *timestamppb.Timestamp
}

// UpdateFileInput is a transport-agnostic input for metadata updates.
type UpdateFileInput struct {
	TenantID      string
	FileID        string
	Filename      *string
	ContentType   *string
	FileType      *storageentityv1.FileType
	ReferenceID   *string
	ReferenceType *string
	IsPublic      *bool
	ExpiresAt     *timestamppb.Timestamp
	ClearExpires  bool
	UpdatedBy     string
}

// NewStorageService creates a new storage service
func NewStorageService(fileRepo *repository.FileRepository, s3Client *s3.Client) *StorageService {
	return NewStorageServiceWithPublisher(fileRepo, s3Client, nil)
}

// NewStorageServiceWithPublisher creates a new storage service with optional event publisher.
func NewStorageServiceWithPublisher(
	fileRepo *repository.FileRepository,
	s3Client *s3.Client,
	eventPublisher *events.Publisher,
) *StorageService {
	return &StorageService{
		fileRepo:       fileRepo,
		s3Client:       s3Client,
		eventPublisher: eventPublisher,
		userFileIndex:  index.NewUserFileIndex(),
	}
}

// UploadFile uploads a file to S3 and stores metadata in database
func (s *StorageService) UploadFile(
	ctx context.Context,
	tenantID string,
	content []byte,
	filename string,
	contentType string,
	fileType storageentityv1.FileType,
	referenceID string,
	referenceType string,
	isPublic bool,
	expiresAt *timestamppb.Timestamp,
	uploadedBy string,
) (*storageentityv1.StoredFile, error) {
	return s.uploadFileInternal(
		ctx,
		tenantID,
		content,
		filename,
		contentType,
		fileType,
		referenceID,
		referenceType,
		isPublic,
		expiresAt,
		uploadedBy,
		true,
		"DIRECT",
	)
}

func (s *StorageService) uploadFileInternal(
	ctx context.Context,
	tenantID string,
	content []byte,
	filename string,
	contentType string,
	fileType storageentityv1.FileType,
	referenceID string,
	referenceType string,
	isPublic bool,
	expiresAt *timestamppb.Timestamp,
	uploadedBy string,
	emitUploadedEvent bool,
	uploadedEventSource string,
) (*storageentityv1.StoredFile, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if len(content) == 0 {
		return nil, fmt.Errorf("%w: content is required", ErrInvalidInput)
	}
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, fmt.Errorf("%w: filename is required", ErrInvalidInput)
	}
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if strings.TrimSpace(uploadedBy) == "" {
		uploadedBy = tenantID
	}

	// Generate unique file ID
	fileID := uuid.New().String()

	// Generate insurance-domain key.
	s3Key := s.s3Client.GenerateInsuranceKey(tenantID, fileID, referenceType, referenceID, filename)

	// Upload to S3
	url, cdnURL, err := s.s3Client.UploadFile(ctx, s3Key, content, contentType, isPublic)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Create database record
	file := &storageentityv1.StoredFile{
		FileId:        fileID,
		TenantId:      tenantID,
		Filename:      filename,
		ContentType:   contentType,
		SizeBytes:     int64(len(content)),
		StorageKey:    s3Key,
		Bucket:        s.s3Client.GetBucket(),
		Url:           url,
		CdnUrl:        cdnURL,
		FileType:      fileType,
		ReferenceId:   referenceID,
		ReferenceType: referenceType,
		IsPublic:      isPublic,
		ExpiresAt:     expiresAt,
		UploadedBy:    uploadedBy,
	}

	created, err := s.fileRepo.Create(ctx, tenantID, file)
	if err != nil {
		// Try to cleanup S3 file if database insert fails
		_ = s.s3Client.DeleteFile(ctx, s3Key)
		return nil, fmt.Errorf("failed to store file metadata: %w", err)
	}
	if s.userFileIndex != nil {
		s.userFileIndex.Upsert(created)
	}
	if emitUploadedEvent && s.eventPublisher != nil {
		_ = s.eventPublisher.PublishFileUploaded(ctx, created, uploadedEventSource, uploadedBy)
	}

	return created, nil
}

// UploadFiles uploads multiple files atomically. If any upload fails, created files are rolled back.
func (s *StorageService) UploadFiles(
	ctx context.Context,
	tenantID string,
	files []UploadFileInput,
	uploadedBy string,
) ([]*storageentityv1.StoredFile, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("%w: files are required", ErrInvalidInput)
	}
	uploadedBy = strings.TrimSpace(uploadedBy)
	if uploadedBy == "" {
		uploadedBy = tenantID
	}

	created := make([]*storageentityv1.StoredFile, 0, len(files))
	for i, f := range files {
		file, err := s.uploadFileInternal(
			ctx,
			tenantID,
			f.Content,
			f.Filename,
			f.ContentType,
			f.FileType,
			f.ReferenceID,
			f.ReferenceType,
			f.IsPublic,
			f.ExpiresAt,
			uploadedBy,
			false,
			"",
		)
		if err != nil {
			rollbackCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			for _, c := range created {
				_ = s.DeleteFile(rollbackCtx, tenantID, c.FileId, uploadedBy)
			}
			cancel()
			return nil, fmt.Errorf("batch upload failed at index %d: %w", i, err)
		}
		created = append(created, file)
	}
	if s.eventPublisher != nil {
		for _, file := range created {
			_ = s.eventPublisher.PublishFileUploaded(ctx, file, "BATCH", uploadedBy)
		}
	}

	return created, nil
}

// GetFile retrieves file metadata
func (s *StorageService) GetFile(ctx context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if strings.TrimSpace(fileID) == "" {
		return nil, fmt.Errorf("%w: file_id is required", ErrInvalidInput)
	}

	file, err := s.fileRepo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		if errors.Is(err, repository.ErrFileNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return file, nil
}

// UpdateFileMetadata partially updates mutable metadata fields for a file.
func (s *StorageService) UpdateFileMetadata(ctx context.Context, in *UpdateFileInput) (*storageentityv1.StoredFile, error) {
	if in == nil {
		return nil, ErrNoMetadataUpdates
	}

	tenantID := strings.TrimSpace(in.TenantID)
	fileID := strings.TrimSpace(in.FileID)
	if tenantID == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if fileID == "" {
		return nil, fmt.Errorf("%w: file_id is required", ErrInvalidInput)
	}

	if in.Filename != nil {
		trimmed := strings.TrimSpace(*in.Filename)
		if trimmed == "" {
			return nil, fmt.Errorf("%w: filename cannot be empty", ErrInvalidInput)
		}
		in.Filename = &trimmed
	}
	if in.ContentType != nil {
		trimmed := strings.TrimSpace(*in.ContentType)
		if trimmed == "" {
			return nil, fmt.Errorf("%w: content_type cannot be empty", ErrInvalidInput)
		}
		in.ContentType = &trimmed
	}
	if in.ReferenceType != nil {
		trimmed := strings.TrimSpace(*in.ReferenceType)
		in.ReferenceType = &trimmed
	}
	if in.ReferenceID != nil {
		trimmed := strings.TrimSpace(*in.ReferenceID)
		in.ReferenceID = &trimmed
	}

	updatedBy := strings.TrimSpace(in.UpdatedBy)
	var updatedByPtr *string
	if updatedBy != "" {
		updatedByPtr = &updatedBy
	}

	patch := &repository.FileMetadataPatch{
		Filename:      in.Filename,
		ContentType:   in.ContentType,
		FileType:      in.FileType,
		ReferenceID:   in.ReferenceID,
		ReferenceType: in.ReferenceType,
		IsPublic:      in.IsPublic,
		ExpiresAt:     in.ExpiresAt,
		ClearExpires:  in.ClearExpires,
		UploadedBy:    updatedByPtr,
	}

	file, err := s.fileRepo.UpdateMetadata(ctx, tenantID, fileID, patch)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrFileNotFound):
			return nil, ErrFileNotFound
		case errors.Is(err, repository.ErrNoMetadataUpdates):
			return nil, ErrNoMetadataUpdates
		default:
			return nil, err
		}
	}
	if s.eventPublisher != nil {
		updatedFields := make([]string, 0, 8)
		if in.Filename != nil {
			updatedFields = append(updatedFields, "filename")
		}
		if in.ContentType != nil {
			updatedFields = append(updatedFields, "content_type")
		}
		if in.FileType != nil {
			updatedFields = append(updatedFields, "file_type")
		}
		if in.ReferenceID != nil {
			updatedFields = append(updatedFields, "reference_id")
		}
		if in.ReferenceType != nil {
			updatedFields = append(updatedFields, "reference_type")
		}
		if in.IsPublic != nil {
			updatedFields = append(updatedFields, "is_public")
		}
		if in.ExpiresAt != nil || in.ClearExpires {
			updatedFields = append(updatedFields, "expires_at")
		}
		_ = s.eventPublisher.PublishFileMetadataUpdated(ctx, tenantID, fileID, updatedFields, updatedBy)
	}
	if s.userFileIndex != nil {
		s.userFileIndex.Upsert(file)
	}
	return file, nil
}

// GetPresignedUploadURL generates a presigned URL for direct upload
func (s *StorageService) GetPresignedUploadURL(
	ctx context.Context,
	tenantID string,
	filename string,
	contentType string,
	fileType storageentityv1.FileType,
	referenceID string,
	referenceType string,
	isPublic bool,
	expiresInMinutes int32,
	uploadedBy string,
) (string, string, string, error) {
	if strings.TrimSpace(uploadedBy) == "" {
		uploadedBy = tenantID
	}
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}

	// Generate file ID and S3 key
	fileID := uuid.New().String()
	s3Key := s.s3Client.GenerateInsuranceKey(tenantID, fileID, referenceType, referenceID, filename)

	// Generate presigned URL
	expiresIn := time.Duration(expiresInMinutes) * time.Minute
	if expiresInMinutes == 0 {
		expiresIn = 15 * time.Minute // Default 15 minutes
	}

	uploadURL, err := s.s3Client.GetPresignedUploadURL(ctx, s3Key, expiresIn)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Persist placeholder metadata for finalize step.
	url, cdnURL := s.s3Client.BuildObjectURLs(s3Key)
	file := &storageentityv1.StoredFile{
		FileId:        fileID,
		TenantId:      tenantID,
		Filename:      filename,
		ContentType:   strings.TrimSpace(contentType),
		SizeBytes:     0,
		StorageKey:    s3Key,
		Bucket:        s.s3Client.GetBucket(),
		Url:           url,
		CdnUrl:        cdnURL,
		FileType:      fileType,
		ReferenceId:   referenceID,
		ReferenceType: referenceType,
		IsPublic:      isPublic,
		UploadedBy:    uploadedBy,
	}
	created, err := s.fileRepo.Create(ctx, tenantID, file)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create pending file metadata: %w", err)
	}
	if s.userFileIndex != nil {
		s.userFileIndex.Upsert(created)
	}
	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishUploadURLIssued(
			ctx,
			tenantID,
			fileID,
			filename,
			s3Key,
			referenceID,
			referenceType,
			isPublic,
			time.Now().Add(expiresIn).UTC(),
			uploadedBy,
		)
	}

	return uploadURL, fileID, s3Key, nil
}

// FinalizeDirectUpload verifies uploaded object and updates metadata.
func (s *StorageService) FinalizeDirectUpload(
	ctx context.Context,
	tenantID string,
	fileID string,
	filename string,
	contentType string,
	fileType storageentityv1.FileType,
	referenceID string,
	referenceType string,
	isPublic bool,
	expiresAt *timestamppb.Timestamp,
	uploadedBy string,
) (*storageentityv1.StoredFile, error) {
	if strings.TrimSpace(uploadedBy) == "" {
		uploadedBy = tenantID
	}

	existing, err := s.fileRepo.GetByID(ctx, tenantID, fileID)
	if err != nil {
		if errors.Is(err, repository.ErrFileNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get pending file metadata: %w", err)
	}

	sizeBytes, detectedContentType, err := s.s3Client.HeadObject(ctx, existing.StorageKey)
	if err != nil {
		return nil, fmt.Errorf("uploaded object not found in storage: %w", err)
	}

	finalContentType := strings.TrimSpace(contentType)
	if finalContentType == "" {
		finalContentType = detectedContentType
	}
	if finalContentType == "" {
		finalContentType = existing.ContentType
	}

	updated := &storageentityv1.StoredFile{
		FileId:        fileID,
		TenantId:      tenantID,
		Filename:      filename,
		ContentType:   finalContentType,
		SizeBytes:     sizeBytes,
		StorageKey:    existing.StorageKey,
		Bucket:        existing.Bucket,
		Url:           existing.Url,
		CdnUrl:        existing.CdnUrl,
		FileType:      fileType,
		ReferenceId:   referenceID,
		ReferenceType: referenceType,
		IsPublic:      isPublic,
		ExpiresAt:     expiresAt,
		UploadedBy:    uploadedBy,
	}

	finalFile, err := s.fileRepo.UpdateAfterDirectUpload(ctx, tenantID, updated)
	if err != nil {
		return nil, err
	}
	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishFileUploadFinalized(ctx, finalFile, uploadedBy)
		_ = s.eventPublisher.PublishFileUploaded(ctx, finalFile, "FINALIZE", uploadedBy)
	}
	if s.userFileIndex != nil {
		s.userFileIndex.Upsert(finalFile)
	}
	return finalFile, nil
}

// GetPresignedDownloadURL generates a presigned URL for downloading
func (s *StorageService) GetPresignedDownloadURL(
	ctx context.Context,
	tenantID string,
	fileID string,
	expiresInMinutes int32,
) (string, *timestamppb.Timestamp, error) {
	// Get file metadata
	file, err := s.GetFile(ctx, tenantID, fileID)
	if err != nil {
		return "", nil, err
	}

	// Generate presigned URL
	expiresIn := time.Duration(expiresInMinutes) * time.Minute
	if expiresInMinutes == 0 {
		expiresIn = 60 * time.Minute // Default 60 minutes
	}

	downloadURL, err := s.s3Client.GetPresignedDownloadURL(ctx, file.StorageKey, expiresIn)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	expiresAt := timestamppb.New(time.Now().Add(expiresIn))
	return downloadURL, expiresAt, nil
}

// DeleteFile deletes a file from S3 and database
func (s *StorageService) DeleteFile(ctx context.Context, tenantID string, fileID string, deletedBy string) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if strings.TrimSpace(fileID) == "" {
		return fmt.Errorf("%w: file_id is required", ErrInvalidInput)
	}

	deletedBy = strings.TrimSpace(deletedBy)
	if deletedBy == "" {
		deletedBy = tenantID
	}

	// Get file metadata
	file, err := s.GetFile(ctx, tenantID, fileID)
	if err != nil {
		return err
	}

	// Delete from S3
	if err := s.s3Client.DeleteFile(ctx, file.StorageKey); err != nil {
		return fmt.Errorf("%w: failed to delete from S3: %v", ErrStorageUnavailable, err)
	}

	// Delete from database
	if err := s.fileRepo.Delete(ctx, tenantID, fileID); err != nil {
		if errors.Is(err, repository.ErrFileNotFound) {
			return ErrFileNotFound
		}
		return fmt.Errorf("failed to delete file metadata: %w", err)
	}
	if s.userFileIndex != nil {
		s.userFileIndex.Delete(tenantID, fileID)
	}
	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishFileDeleted(ctx, tenantID, fileID, file.StorageKey, deletedBy)
	}

	return nil
}

// ListFiles lists files with filters
func (s *StorageService) ListFiles(
	ctx context.Context,
	tenantID string,
	fileType storageentityv1.FileType,
	referenceID string,
	referenceType string,
	uploadedBy string,
	limit, offset int32,
) ([]*storageentityv1.StoredFile, int, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, 0, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	uploadedBy = strings.TrimSpace(uploadedBy)
	if uploadedBy != "" {
		if s.userFileIndex != nil {
			if files, total, ok := s.userFileIndex.List(tenantID, uploadedBy, fileType, referenceID, referenceType, limit, offset); ok {
				return files, total, nil
			}
		}

		allUserFiles, err := s.fileRepo.ListAllByUploadedBy(ctx, tenantID, uploadedBy)
		if err != nil {
			return nil, 0, err
		}
		if s.userFileIndex != nil {
			s.userFileIndex.WarmUser(tenantID, uploadedBy, allUserFiles)
			if files, total, ok := s.userFileIndex.List(tenantID, uploadedBy, fileType, referenceID, referenceType, limit, offset); ok {
				return files, total, nil
			}
		}

		files, total := filterAndPaginateFiles(allUserFiles, fileType, referenceID, referenceType, limit, offset)
		return files, total, nil
	}

	return s.fileRepo.List(ctx, tenantID, fileType, referenceID, referenceType, limit, offset)
}

func filterAndPaginateFiles(
	files []*storageentityv1.StoredFile,
	fileType storageentityv1.FileType,
	referenceID string,
	referenceType string,
	limit int32,
	offset int32,
) ([]*storageentityv1.StoredFile, int) {
	referenceID = strings.TrimSpace(referenceID)
	referenceType = strings.TrimSpace(referenceType)

	filtered := make([]*storageentityv1.StoredFile, 0, len(files))
	for _, file := range files {
		if file == nil {
			continue
		}
		if fileType != storageentityv1.FileType_FILE_TYPE_UNSPECIFIED && file.GetFileType() != fileType {
			continue
		}
		if referenceID != "" && file.GetReferenceId() != referenceID {
			continue
		}
		if referenceType != "" && file.GetReferenceType() != referenceType {
			continue
		}
		filtered = append(filtered, file)
	}

	total := len(filtered)
	start := int(offset)
	if start < 0 {
		start = 0
	}
	if start > total {
		start = total
	}
	end := total
	if limit > 0 {
		end = start + int(limit)
		if end > total {
			end = total
		}
	}

	return filtered[start:end], total
}
