package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/service"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
)

type fakeStorageService struct {
	uploadFileFn              func(ctx context.Context, tenantID string, content []byte, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error)
	uploadFilesFn             func(ctx context.Context, tenantID string, files []service.UploadFileInput, uploadedBy string) ([]*storageentityv1.StoredFile, error)
	getFileFn                 func(ctx context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error)
	updateFileMetadataFn      func(ctx context.Context, in *service.UpdateFileInput) (*storageentityv1.StoredFile, error)
	getPresignedUploadURLFn   func(ctx context.Context, tenantID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresInMinutes int32, uploadedBy string) (string, string, string, error)
	finalizeDirectUploadFn    func(ctx context.Context, tenantID, fileID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error)
	getPresignedDownloadURLFn func(ctx context.Context, tenantID string, fileID string, expiresInMinutes int32) (string, *timestamppb.Timestamp, error)
	deleteFileFn              func(ctx context.Context, tenantID string, fileID string, deletedBy string) error
	listFilesFn               func(ctx context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error)
}

func (f *fakeStorageService) UploadFile(ctx context.Context, tenantID string, content []byte, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error) {
	return f.uploadFileFn(ctx, tenantID, content, filename, contentType, fileType, referenceID, referenceType, isPublic, expiresAt, uploadedBy)
}

func (f *fakeStorageService) UploadFiles(ctx context.Context, tenantID string, files []service.UploadFileInput, uploadedBy string) ([]*storageentityv1.StoredFile, error) {
	return f.uploadFilesFn(ctx, tenantID, files, uploadedBy)
}

func (f *fakeStorageService) GetFile(ctx context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error) {
	return f.getFileFn(ctx, tenantID, fileID)
}

func (f *fakeStorageService) UpdateFileMetadata(ctx context.Context, in *service.UpdateFileInput) (*storageentityv1.StoredFile, error) {
	return f.updateFileMetadataFn(ctx, in)
}

func (f *fakeStorageService) GetPresignedUploadURL(ctx context.Context, tenantID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresInMinutes int32, uploadedBy string) (string, string, string, error) {
	return f.getPresignedUploadURLFn(ctx, tenantID, filename, contentType, fileType, referenceID, referenceType, isPublic, expiresInMinutes, uploadedBy)
}

func (f *fakeStorageService) FinalizeDirectUpload(ctx context.Context, tenantID, fileID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error) {
	return f.finalizeDirectUploadFn(ctx, tenantID, fileID, filename, contentType, fileType, referenceID, referenceType, isPublic, expiresAt, uploadedBy)
}

func (f *fakeStorageService) GetPresignedDownloadURL(ctx context.Context, tenantID string, fileID string, expiresInMinutes int32) (string, *timestamppb.Timestamp, error) {
	return f.getPresignedDownloadURLFn(ctx, tenantID, fileID, expiresInMinutes)
}

func (f *fakeStorageService) DeleteFile(ctx context.Context, tenantID string, fileID string, deletedBy string) error {
	return f.deleteFileFn(ctx, tenantID, fileID, deletedBy)
}

func (f *fakeStorageService) ListFiles(ctx context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error) {
	return f.listFilesFn(ctx, tenantID, fileType, referenceID, referenceType, uploadedBy, limit, offset)
}

func TestActorFromContext(t *testing.T) {
	assert.Equal(t, "fallback", actorFromContext(context.Background(), "fallback"))

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-actor-id", " actor-1 "))
	assert.Equal(t, "actor-1", actorFromContext(ctx, "fallback"))

	ctx = metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-user-id", "", "user-id", "user-2"))
	assert.Equal(t, "user-2", actorFromContext(ctx, "fallback"))
}

func TestUploadFileHandler(t *testing.T) {
	var gotContentType string
	var gotUploadedBy string
	handler := NewStorageHandler(&fakeStorageService{
		uploadFileFn: func(_ context.Context, tenantID string, content []byte, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error) {
			gotContentType = contentType
			gotUploadedBy = uploadedBy
			return &storageentityv1.StoredFile{FileId: "file-1", TenantId: tenantID, Filename: filename}, nil
		},
	})

	_, err := handler.UploadFile(context.Background(), &storageservicev1.UploadFileRequest{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-user-id", "user-1"))
	resp, err := handler.UploadFile(ctx, &storageservicev1.UploadFileRequest{
		TenantId: "tenant-1",
		Content:  []byte("hello"),
		Filename: "doc.pdf",
	})
	require.NoError(t, err)
	assert.Equal(t, "file-1", resp.GetFile().GetFileId())
	assert.Equal(t, "application/octet-stream", gotContentType)
	assert.Equal(t, "user-1", gotUploadedBy)

	handler = NewStorageHandler(&fakeStorageService{
		uploadFileFn: func(context.Context, string, []byte, string, string, storageentityv1.FileType, string, string, bool, *timestamppb.Timestamp, string) (*storageentityv1.StoredFile, error) {
			return nil, service.ErrInvalidInput
		},
	})
	_, err = handler.UploadFile(ctx, &storageservicev1.UploadFileRequest{TenantId: "tenant-1", Content: []byte("hello"), Filename: "doc.pdf"})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	handler = NewStorageHandler(&fakeStorageService{
		uploadFileFn: func(context.Context, string, []byte, string, string, storageentityv1.FileType, string, string, bool, *timestamppb.Timestamp, string) (*storageentityv1.StoredFile, error) {
			return nil, errors.New("boom")
		},
	})
	_, err = handler.UploadFile(ctx, &storageservicev1.UploadFileRequest{TenantId: "tenant-1", Content: []byte("hello"), Filename: "doc.pdf"})
	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestUploadFilesHandler(t *testing.T) {
	var gotUploadedBy string
	handler := NewStorageHandler(&fakeStorageService{
		uploadFilesFn: func(_ context.Context, _ string, files []service.UploadFileInput, uploadedBy string) ([]*storageentityv1.StoredFile, error) {
			gotUploadedBy = uploadedBy
			return []*storageentityv1.StoredFile{{FileId: files[0].Filename}}, nil
		},
	})

	_, err := handler.UploadFiles(context.Background(), &storageservicev1.UploadFilesRequest{TenantId: "tenant-1"})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-sub", "subject-1"))
	resp, err := handler.UploadFiles(ctx, &storageservicev1.UploadFilesRequest{
		TenantId: "tenant-1",
		Files: []*storageservicev1.FileUpload{{
			Content:  []byte("hello"),
			Filename: "doc.pdf",
		}},
	})
	require.NoError(t, err)
	assert.Len(t, resp.GetFiles(), 1)
	assert.Equal(t, "subject-1", gotUploadedBy)

	_, err = handler.UploadFiles(context.Background(), &storageservicev1.UploadFilesRequest{
		TenantId: "tenant-1",
		Files:    []*storageservicev1.FileUpload{{Filename: "missing-content"}},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	handler = NewStorageHandler(&fakeStorageService{
		uploadFilesFn: func(context.Context, string, []service.UploadFileInput, string) ([]*storageentityv1.StoredFile, error) {
			return nil, errors.New("boom")
		},
	})
	_, err = handler.UploadFiles(context.Background(), &storageservicev1.UploadFilesRequest{
		TenantId: "tenant-1",
		Files:    []*storageservicev1.FileUpload{{Content: []byte("hello"), Filename: "doc.pdf"}},
	})
	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestGetAndUpdateFileHandlers(t *testing.T) {
	handler := NewStorageHandler(&fakeStorageService{
		getFileFn: func(_ context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error) {
			if fileID == "missing" {
				return nil, service.ErrFileNotFound
			}
			return &storageentityv1.StoredFile{FileId: fileID, TenantId: tenantID}, nil
		},
		updateFileMetadataFn: func(_ context.Context, in *service.UpdateFileInput) (*storageentityv1.StoredFile, error) {
			if in.FileID == "missing" {
				return nil, service.ErrFileNotFound
			}
			if in.Filename == nil && in.ContentType == nil && in.FileType == nil && in.ReferenceID == nil && in.ReferenceType == nil && in.IsPublic == nil && in.ExpiresAt == nil && !in.ClearExpires {
				return nil, service.ErrNoMetadataUpdates
			}
			return &storageentityv1.StoredFile{FileId: in.FileID, TenantId: in.TenantID}, nil
		},
	})

	_, err := handler.GetFile(context.Background(), &storageservicev1.GetFileRequest{TenantId: "tenant-1", FileId: "missing"})
	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	resp, err := handler.GetFile(context.Background(), &storageservicev1.GetFileRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.NoError(t, err)
	assert.Equal(t, "file-1", resp.GetFile().GetFileId())

	handler = NewStorageHandler(&fakeStorageService{
		getFileFn: func(context.Context, string, string) (*storageentityv1.StoredFile, error) {
			return nil, service.ErrInvalidInput
		},
		updateFileMetadataFn: func(_ context.Context, in *service.UpdateFileInput) (*storageentityv1.StoredFile, error) {
			return &storageentityv1.StoredFile{FileId: in.FileID, TenantId: in.TenantID}, nil
		},
	})
	_, err = handler.GetFile(context.Background(), &storageservicev1.GetFileRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	handler = NewStorageHandler(&fakeStorageService{
		getFileFn: func(_ context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error) {
			if fileID == "missing" {
				return nil, service.ErrFileNotFound
			}
			return &storageentityv1.StoredFile{FileId: fileID, TenantId: tenantID}, nil
		},
		updateFileMetadataFn: func(_ context.Context, in *service.UpdateFileInput) (*storageentityv1.StoredFile, error) {
			if in.FileID == "missing" {
				return nil, service.ErrFileNotFound
			}
			if in.Filename == nil && in.ContentType == nil && in.FileType == nil && in.ReferenceID == nil && in.ReferenceType == nil && in.IsPublic == nil && in.ExpiresAt == nil && !in.ClearExpires {
				return nil, service.ErrNoMetadataUpdates
			}
			return &storageentityv1.StoredFile{FileId: in.FileID, TenantId: in.TenantID}, nil
		},
	})
	_, err = handler.UpdateFile(context.Background(), &storageservicev1.UpdateFileRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	filename := "renamed.pdf"
	isPublic := true
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("user-id", "user-99"))
	updateResp, err := handler.UpdateFile(ctx, &storageservicev1.UpdateFileRequest{
		TenantId: "tenant-1",
		FileId:   "file-1",
		Filename: &filename,
		IsPublic: &isPublic,
	})
	require.NoError(t, err)
	assert.Equal(t, "file-1", updateResp.GetFile().GetFileId())

	handler = NewStorageHandler(&fakeStorageService{
		getFileFn: func(_ context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error) {
			return &storageentityv1.StoredFile{FileId: fileID, TenantId: tenantID}, nil
		},
		updateFileMetadataFn: func(context.Context, *service.UpdateFileInput) (*storageentityv1.StoredFile, error) {
			return nil, service.ErrInvalidInput
		},
	})
	_, err = handler.UpdateFile(context.Background(), &storageservicev1.UpdateFileRequest{
		TenantId: "tenant-1",
		FileId:   "file-1",
		Filename: &filename,
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestUploadAndFinalizeURLHandlers(t *testing.T) {
	handler := NewStorageHandler(&fakeStorageService{
		getPresignedUploadURLFn: func(_ context.Context, tenantID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresInMinutes int32, uploadedBy string) (string, string, string, error) {
			return "https://upload", "file-1", "key-1", nil
		},
		finalizeDirectUploadFn: func(_ context.Context, tenantID, fileID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error) {
			if fileID == "missing" {
				return nil, service.ErrFileNotFound
			}
			return &storageentityv1.StoredFile{FileId: fileID, TenantId: tenantID}, nil
		},
	})

	_, err := handler.GetUploadURL(context.Background(), &storageservicev1.GetUploadURLRequest{TenantId: "tenant-1"})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	resp, err := handler.GetUploadURL(context.Background(), &storageservicev1.GetUploadURLRequest{
		TenantId: "tenant-1",
		Filename: "doc.pdf",
	})
	require.NoError(t, err)
	assert.Equal(t, "https://upload", resp.GetUploadUrl())

	handler = NewStorageHandler(&fakeStorageService{
		getPresignedUploadURLFn: func(context.Context, string, string, string, storageentityv1.FileType, string, string, bool, int32, string) (string, string, string, error) {
			return "", "", "", errors.New("boom")
		},
		finalizeDirectUploadFn: func(_ context.Context, tenantID, fileID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error) {
			return &storageentityv1.StoredFile{FileId: fileID, TenantId: tenantID}, nil
		},
	})
	_, err = handler.GetUploadURL(context.Background(), &storageservicev1.GetUploadURLRequest{
		TenantId: "tenant-1",
		Filename: "doc.pdf",
	})
	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))

	handler = NewStorageHandler(&fakeStorageService{
		getPresignedUploadURLFn: func(_ context.Context, tenantID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresInMinutes int32, uploadedBy string) (string, string, string, error) {
			return "https://upload", "file-1", "key-1", nil
		},
		finalizeDirectUploadFn: func(_ context.Context, tenantID, fileID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresAt *timestamppb.Timestamp, uploadedBy string) (*storageentityv1.StoredFile, error) {
			if fileID == "missing" {
				return nil, service.ErrFileNotFound
			}
			return &storageentityv1.StoredFile{FileId: fileID, TenantId: tenantID}, nil
		},
	})
	_, err = handler.FinalizeUpload(context.Background(), &storageservicev1.FinalizeUploadRequest{TenantId: "tenant-1", FileId: "missing"})
	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	finalized, err := handler.FinalizeUpload(context.Background(), &storageservicev1.FinalizeUploadRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.NoError(t, err)
	assert.Equal(t, "file-1", finalized.GetFile().GetFileId())

	handler = NewStorageHandler(&fakeStorageService{
		getPresignedUploadURLFn: func(_ context.Context, tenantID, filename, contentType string, fileType storageentityv1.FileType, referenceID, referenceType string, isPublic bool, expiresInMinutes int32, uploadedBy string) (string, string, string, error) {
			return "https://upload", "file-1", "key-1", nil
		},
		finalizeDirectUploadFn: func(context.Context, string, string, string, string, storageentityv1.FileType, string, string, bool, *timestamppb.Timestamp, string) (*storageentityv1.StoredFile, error) {
			return nil, errors.New("boom")
		},
	})
	_, err = handler.FinalizeUpload(context.Background(), &storageservicev1.FinalizeUploadRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestDownloadDeleteAndListHandlers(t *testing.T) {
	expiresAt := timestamppb.New(time.Now().UTC())
	handler := NewStorageHandler(&fakeStorageService{
		getPresignedDownloadURLFn: func(_ context.Context, tenantID string, fileID string, expiresInMinutes int32) (string, *timestamppb.Timestamp, error) {
			if fileID == "missing" {
				return "", nil, service.ErrFileNotFound
			}
			return "https://download", expiresAt, nil
		},
		deleteFileFn: func(_ context.Context, tenantID string, fileID string, deletedBy string) error {
			switch fileID {
			case "missing":
				return service.ErrFileNotFound
			case "storage-down":
				return service.ErrStorageUnavailable
			default:
				return nil
			}
		},
		listFilesFn: func(_ context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error) {
			if tenantID == "bad-tenant" {
				return nil, 0, service.ErrInvalidInput
			}
			return []*storageentityv1.StoredFile{{FileId: "file-1", TenantId: tenantID}}, 5, nil
		},
	})

	_, err := handler.GetDownloadURL(context.Background(), &storageservicev1.GetDownloadURLRequest{TenantId: "tenant-1", FileId: "missing"})
	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	downloadResp, err := handler.GetDownloadURL(context.Background(), &storageservicev1.GetDownloadURLRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.NoError(t, err)
	assert.Equal(t, "https://download", downloadResp.GetDownloadUrl())
	assert.Equal(t, expiresAt, downloadResp.GetExpiresAt())

	handler = NewStorageHandler(&fakeStorageService{
		getPresignedDownloadURLFn: func(context.Context, string, string, int32) (string, *timestamppb.Timestamp, error) {
			return "", nil, errors.New("boom")
		},
		deleteFileFn: func(_ context.Context, tenantID string, fileID string, deletedBy string) error {
			return nil
		},
		listFilesFn: func(_ context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error) {
			return []*storageentityv1.StoredFile{{FileId: "file-1", TenantId: tenantID}}, 5, nil
		},
	})
	_, err = handler.GetDownloadURL(context.Background(), &storageservicev1.GetDownloadURLRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))

	handler = NewStorageHandler(&fakeStorageService{
		getPresignedDownloadURLFn: func(_ context.Context, tenantID string, fileID string, expiresInMinutes int32) (string, *timestamppb.Timestamp, error) {
			if fileID == "missing" {
				return "", nil, service.ErrFileNotFound
			}
			return "https://download", expiresAt, nil
		},
		deleteFileFn: func(_ context.Context, tenantID string, fileID string, deletedBy string) error {
			switch fileID {
			case "missing":
				return service.ErrFileNotFound
			case "storage-down":
				return service.ErrStorageUnavailable
			default:
				return nil
			}
		},
		listFilesFn: func(_ context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error) {
			if tenantID == "bad-tenant" {
				return nil, 0, service.ErrInvalidInput
			}
			return []*storageentityv1.StoredFile{{FileId: "file-1", TenantId: tenantID}}, 5, nil
		},
	})
	_, err = handler.DeleteFile(context.Background(), &storageservicev1.DeleteFileRequest{TenantId: "tenant-1", FileId: "missing"})
	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	_, err = handler.DeleteFile(context.Background(), &storageservicev1.DeleteFileRequest{TenantId: "tenant-1", FileId: "storage-down"})
	require.Error(t, err)
	assert.Equal(t, codes.Unavailable, status.Code(err))

	deleteResp, err := handler.DeleteFile(context.Background(), &storageservicev1.DeleteFileRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.NoError(t, err)
	assert.True(t, deleteResp.GetSuccess())

	handler = NewStorageHandler(&fakeStorageService{
		getPresignedDownloadURLFn: func(_ context.Context, tenantID string, fileID string, expiresInMinutes int32) (string, *timestamppb.Timestamp, error) {
			return "https://download", expiresAt, nil
		},
		deleteFileFn: func(_ context.Context, tenantID string, fileID string, deletedBy string) error {
			return service.ErrInvalidInput
		},
		listFilesFn: func(_ context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error) {
			return nil, 0, errors.New("boom")
		},
	})
	_, err = handler.DeleteFile(context.Background(), &storageservicev1.DeleteFileRequest{TenantId: "tenant-1", FileId: "file-1"})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	handler = NewStorageHandler(&fakeStorageService{
		getPresignedDownloadURLFn: func(_ context.Context, tenantID string, fileID string, expiresInMinutes int32) (string, *timestamppb.Timestamp, error) {
			return "https://download", expiresAt, nil
		},
		deleteFileFn: func(_ context.Context, tenantID string, fileID string, deletedBy string) error {
			return nil
		},
		listFilesFn: func(_ context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error) {
			if tenantID == "bad-tenant" {
				return nil, 0, service.ErrInvalidInput
			}
			return []*storageentityv1.StoredFile{{FileId: "file-1", TenantId: tenantID}}, 5, nil
		},
	})
	_, err = handler.ListFiles(context.Background(), &storageservicev1.ListFilesRequest{
		TenantId: "tenant-1",
		Page:     &commonv1.PaginationRequest{Page: -1},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	listResp, err := handler.ListFiles(context.Background(), &storageservicev1.ListFilesRequest{
		TenantId: "tenant-1",
		Page:     &commonv1.PaginationRequest{Page: 2, PageSize: 2},
	})
	require.NoError(t, err)
	assert.Len(t, listResp.GetFiles(), 1)
	assert.Equal(t, int32(5), listResp.GetPage().GetTotalItems())
	assert.Equal(t, int32(3), listResp.GetPage().GetTotalPages())
	assert.True(t, listResp.GetPage().GetHasNext())
	assert.True(t, listResp.GetPage().GetHasPrevious())

	_, err = handler.ListFiles(context.Background(), &storageservicev1.ListFilesRequest{TenantId: "bad-tenant"})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	_, err = handler.ListFiles(context.Background(), &storageservicev1.ListFilesRequest{
		TenantId: "tenant-1",
		Page:     &commonv1.PaginationRequest{PageSize: -1},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	handler = NewStorageHandler(&fakeStorageService{
		getPresignedDownloadURLFn: func(_ context.Context, tenantID string, fileID string, expiresInMinutes int32) (string, *timestamppb.Timestamp, error) {
			return "https://download", expiresAt, nil
		},
		deleteFileFn: func(_ context.Context, tenantID string, fileID string, deletedBy string) error {
			return nil
		},
		listFilesFn: func(_ context.Context, tenantID string, fileType storageentityv1.FileType, referenceID, referenceType, uploadedBy string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error) {
			return nil, 0, errors.New("boom")
		},
	})
	_, err = handler.ListFiles(context.Background(), &storageservicev1.ListFilesRequest{TenantId: "tenant-1"})
	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}
