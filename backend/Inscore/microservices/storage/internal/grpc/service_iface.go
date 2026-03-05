package server

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/service"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
)

// storageServiceIface decouples transport from concrete service implementation.
type storageServiceIface interface {
	UploadFile(ctx context.Context, tenantID string, content []byte, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error)
	UploadFiles(ctx context.Context, tenantID string, files []service.UploadFileInput, uploadedBy string) ([]*storageentityv1.StoredFile, error)
	GetFile(ctx context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error)
	UpdateFileMetadata(ctx context.Context, in *service.UpdateFileInput) (*storageentityv1.StoredFile, error)
	GetPresignedUploadURL(ctx context.Context, tenantID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresInMinutes int32, uploadedBy string) (string, string, string, error)
	FinalizeDirectUpload(ctx context.Context, tenantID, fileID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error)
	GetPresignedDownloadURL(ctx context.Context, tenantID string, fileID string, expiresInMinutes int32) (string, *timestamppb.Timestamp, error)
	DeleteFile(ctx context.Context, tenantID string, fileID string, deletedBy string) error
	ListFiles(ctx context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error)
}
