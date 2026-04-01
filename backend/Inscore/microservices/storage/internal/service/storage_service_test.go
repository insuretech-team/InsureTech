package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/index"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/repository"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
)

type fakeFileRepo struct {
	files                map[string]*storageentityv1.StoredFile
	createErr            error
	createErrAt          int
	createCalls          int
	listResult           []*storageentityv1.StoredFile
	listTotal            int
	listErr              error
	listAllResult        []*storageentityv1.StoredFile
	listAllErr           error
	listAllCalls         int
	getErr               error
	deleteErr            error
	deletedIDs           []string
	updateAfterResult    *storageentityv1.StoredFile
	updateAfterErr       error
	updateMetadataResult *storageentityv1.StoredFile
	updateMetadataErr    error
	lastPatch            *repository.FileMetadataPatch
}

func newFakeFileRepo() *fakeFileRepo {
	return &fakeFileRepo{files: map[string]*storageentityv1.StoredFile{}}
}

func (f *fakeFileRepo) Create(_ context.Context, tenantID string, file *storageentityv1.StoredFile) (*storageentityv1.StoredFile, error) {
	f.createCalls++
	if f.createErr != nil && (f.createErrAt == 0 || f.createCalls == f.createErrAt) {
		return nil, f.createErr
	}
	cloned := cloneStoredFile(file)
	cloned.TenantId = tenantID
	f.files[cloned.FileId] = cloned
	return cloneStoredFile(cloned), nil
}

func (f *fakeFileRepo) GetByID(_ context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	file, ok := f.files[fileID]
	if !ok || file.GetTenantId() != tenantID {
		return nil, repository.ErrFileNotFound
	}
	return cloneStoredFile(file), nil
}

func (f *fakeFileRepo) List(_ context.Context, _ string, _ storageentityv1.FileType, _, _ string, _, _ int32) ([]*storageentityv1.StoredFile, int, error) {
	if f.listErr != nil {
		return nil, 0, f.listErr
	}
	return cloneStoredFiles(f.listResult), f.listTotal, nil
}

func (f *fakeFileRepo) ListAllByUploadedBy(_ context.Context, _ string, _ string) ([]*storageentityv1.StoredFile, error) {
	f.listAllCalls++
	if f.listAllErr != nil {
		return nil, f.listAllErr
	}
	return cloneStoredFiles(f.listAllResult), nil
}

func (f *fakeFileRepo) Delete(_ context.Context, _ string, fileID string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	if _, ok := f.files[fileID]; !ok {
		return repository.ErrFileNotFound
	}
	delete(f.files, fileID)
	f.deletedIDs = append(f.deletedIDs, fileID)
	return nil
}

func (f *fakeFileRepo) UpdateAfterDirectUpload(_ context.Context, tenantID string, file *storageentityv1.StoredFile) (*storageentityv1.StoredFile, error) {
	if f.updateAfterErr != nil {
		return nil, f.updateAfterErr
	}
	if f.updateAfterResult != nil {
		f.files[f.updateAfterResult.FileId] = cloneStoredFile(f.updateAfterResult)
		return cloneStoredFile(f.updateAfterResult), nil
	}
	cloned := cloneStoredFile(file)
	cloned.TenantId = tenantID
	f.files[cloned.FileId] = cloned
	return cloneStoredFile(cloned), nil
}

func (f *fakeFileRepo) UpdateMetadata(_ context.Context, tenantID, fileID string, patch *repository.FileMetadataPatch) (*storageentityv1.StoredFile, error) {
	f.lastPatch = patch
	if f.updateMetadataErr != nil {
		return nil, f.updateMetadataErr
	}
	if f.updateMetadataResult != nil {
		f.files[fileID] = cloneStoredFile(f.updateMetadataResult)
		return cloneStoredFile(f.updateMetadataResult), nil
	}
	file, ok := f.files[fileID]
	if !ok {
		return nil, repository.ErrFileNotFound
	}
	cloned := cloneStoredFile(file)
	cloned.TenantId = tenantID
	return cloned, nil
}

type uploadCall struct {
	key         string
	contentType string
	isPublic    bool
}

type fakeObjectStore struct {
	bucket                string
	generatedKeys         []string
	uploadCalls           []uploadCall
	uploadURL             string
	uploadCDNURL          string
	uploadErr             error
	deletedKeys           []string
	deleteErr             error
	headSize              int64
	headContentType       string
	headErr               error
	presignedUploadURL    string
	presignedUploadErr    error
	presignedUploadExpiry time.Duration
	presignedDownloadURL  string
	presignedDownloadErr  error
	presignedDownloadExp  time.Duration
	objectURL             string
	cdnURL                string
}

func (f *fakeObjectStore) GenerateInsuranceKey(tenantID string, fileID string, referenceType string, referenceID string, filename string) string {
	key := fmt.Sprintf("%s/%s/%s/%s", tenantID, referenceType, referenceID, filename)
	f.generatedKeys = append(f.generatedKeys, key)
	return key
}

func (f *fakeObjectStore) UploadFile(_ context.Context, key string, _ []byte, contentType string, isPublic bool) (string, string, error) {
	f.uploadCalls = append(f.uploadCalls, uploadCall{key: key, contentType: contentType, isPublic: isPublic})
	if f.uploadErr != nil {
		return "", "", f.uploadErr
	}
	return f.uploadURL, f.uploadCDNURL, nil
}

func (f *fakeObjectStore) DeleteFile(_ context.Context, key string) error {
	f.deletedKeys = append(f.deletedKeys, key)
	if f.deleteErr != nil {
		return f.deleteErr
	}
	return nil
}

func (f *fakeObjectStore) HeadObject(_ context.Context, _ string) (int64, string, error) {
	if f.headErr != nil {
		return 0, "", f.headErr
	}
	return f.headSize, f.headContentType, nil
}

func (f *fakeObjectStore) GetPresignedUploadURL(_ context.Context, _ string, expiresIn time.Duration) (string, error) {
	f.presignedUploadExpiry = expiresIn
	if f.presignedUploadErr != nil {
		return "", f.presignedUploadErr
	}
	return f.presignedUploadURL, nil
}

func (f *fakeObjectStore) GetPresignedDownloadURL(_ context.Context, _ string, expiresIn time.Duration) (string, error) {
	f.presignedDownloadExp = expiresIn
	if f.presignedDownloadErr != nil {
		return "", f.presignedDownloadErr
	}
	return f.presignedDownloadURL, nil
}

func (f *fakeObjectStore) GetBucket() string {
	return f.bucket
}

func (f *fakeObjectStore) BuildObjectURLs(_ string) (string, string) {
	return f.objectURL, f.cdnURL
}

type fakeEventPublisher struct {
	uploadedSources   []string
	uploadedBy        []string
	metadataUpdated   [][]string
	metadataUpdatedBy []string
	uploadURLIssuedBy []string
	finalizedBy       []string
	deletedBy         []string
}

func (f *fakeEventPublisher) PublishFileUploaded(_ context.Context, _ *storageentityv1.StoredFile, source, uploadedBy string) error {
	f.uploadedSources = append(f.uploadedSources, source)
	f.uploadedBy = append(f.uploadedBy, uploadedBy)
	return nil
}

func (f *fakeEventPublisher) PublishUploadURLIssued(_ context.Context, _, _, _, _, _, _ string, _ bool, _ time.Time, requestedBy string) error {
	f.uploadURLIssuedBy = append(f.uploadURLIssuedBy, requestedBy)
	return nil
}

func (f *fakeEventPublisher) PublishFileUploadFinalized(_ context.Context, _ *storageentityv1.StoredFile, finalizedBy string) error {
	f.finalizedBy = append(f.finalizedBy, finalizedBy)
	return nil
}

func (f *fakeEventPublisher) PublishFileMetadataUpdated(_ context.Context, _, _ string, updatedFields []string, updatedBy string) error {
	f.metadataUpdated = append(f.metadataUpdated, append([]string(nil), updatedFields...))
	f.metadataUpdatedBy = append(f.metadataUpdatedBy, updatedBy)
	return nil
}

func (f *fakeEventPublisher) PublishFileDeleted(_ context.Context, _, _, _, deletedBy string) error {
	f.deletedBy = append(f.deletedBy, deletedBy)
	return nil
}

func TestUploadFile(t *testing.T) {
	t.Run("validates required fields", func(t *testing.T) {
		svc := newStorageService(newFakeFileRepo(), &fakeObjectStore{}, nil, index.NewUserFileIndex())

		_, err := svc.UploadFile(context.Background(), "", []byte("x"), "a.txt", "", storageentityv1.FileType_FILE_TYPE_DOCUMENT, "", "", false, nil, "")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInput)

		_, err = svc.UploadFile(context.Background(), "tenant-1", nil, "a.txt", "", storageentityv1.FileType_FILE_TYPE_DOCUMENT, "", "", false, nil, "")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInput)

		_, err = svc.UploadFile(context.Background(), "tenant-1", []byte("x"), "   ", "", storageentityv1.FileType_FILE_TYPE_DOCUMENT, "", "", false, nil, "")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("defaults content type and uploader on success", func(t *testing.T) {
		repo := newFakeFileRepo()
		store := &fakeObjectStore{
			bucket:       "bucket-1",
			uploadURL:    "https://objects/uploaded",
			uploadCDNURL: "https://cdn/uploaded",
		}
		publisher := &fakeEventPublisher{}
		svc := newStorageService(repo, store, publisher, index.NewUserFileIndex())

		file, err := svc.UploadFile(
			context.Background(),
			"tenant-1",
			[]byte("hello"),
			" doc.pdf ",
			"",
			storageentityv1.FileType_FILE_TYPE_DOCUMENT,
			"ref-1",
			"claim",
			true,
			nil,
			"",
		)
		require.NoError(t, err)
		require.NotNil(t, file)
		assert.Equal(t, "application/octet-stream", file.ContentType)
		assert.Equal(t, "tenant-1", file.UploadedBy)
		assert.Equal(t, "bucket-1", file.Bucket)
		assert.Equal(t, "DIRECT", publisher.uploadedSources[0])
		assert.Equal(t, "tenant-1", publisher.uploadedBy[0])
		require.Len(t, store.uploadCalls, 1)
		assert.Equal(t, "application/octet-stream", store.uploadCalls[0].contentType)
		assert.True(t, store.uploadCalls[0].isPublic)
	})

	t.Run("cleans up storage when metadata create fails", func(t *testing.T) {
		repo := newFakeFileRepo()
		repo.createErr = errors.New("insert failed")
		store := &fakeObjectStore{uploadURL: "u", uploadCDNURL: "c"}
		svc := newStorageService(repo, store, nil, index.NewUserFileIndex())

		_, err := svc.UploadFile(
			context.Background(),
			"tenant-1",
			[]byte("hello"),
			"doc.pdf",
			"application/pdf",
			storageentityv1.FileType_FILE_TYPE_DOCUMENT,
			"",
			"",
			false,
			nil,
			"user-1",
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to store file metadata")
		require.Len(t, store.deletedKeys, 1)
		assert.True(t, strings.Contains(store.deletedKeys[0], "tenant-1"))
	})
}

func TestUploadFilesRollbackOnFailure(t *testing.T) {
	repo := newFakeFileRepo()
	repo.createErr = errors.New("boom")
	repo.createErrAt = 2
	store := &fakeObjectStore{uploadURL: "u", uploadCDNURL: "c"}
	svc := newStorageService(repo, store, nil, index.NewUserFileIndex())

	_, err := svc.UploadFiles(context.Background(), "tenant-1", []UploadFileInput{
		{Content: []byte("a"), Filename: "a.pdf", ContentType: "application/pdf"},
		{Content: []byte("b"), Filename: "b.pdf", ContentType: "application/pdf"},
	}, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "batch upload failed at index 1")
	assert.Len(t, repo.deletedIDs, 1)
	assert.Len(t, store.deletedKeys, 2)
}

func TestGetFileAndDeleteFileErrorMapping(t *testing.T) {
	t.Run("maps not found from repository", func(t *testing.T) {
		svc := newStorageService(newFakeFileRepo(), &fakeObjectStore{}, nil, index.NewUserFileIndex())

		_, err := svc.GetFile(context.Background(), "tenant-1", "missing")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrFileNotFound)
	})

	t.Run("maps storage delete failures", func(t *testing.T) {
		repo := newFakeFileRepo()
		repo.files["file-1"] = sampleStoredFile("tenant-1", "file-1", "user-1")
		store := &fakeObjectStore{deleteErr: errors.New("s3 down")}
		svc := newStorageService(repo, store, nil, index.NewUserFileIndex())

		err := svc.DeleteFile(context.Background(), "tenant-1", "file-1", "")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrStorageUnavailable)
	})

	t.Run("deletes successfully and defaults actor", func(t *testing.T) {
		repo := newFakeFileRepo()
		repo.files["file-1"] = sampleStoredFile("tenant-1", "file-1", "user-1")
		store := &fakeObjectStore{}
		publisher := &fakeEventPublisher{}
		svc := newStorageService(repo, store, publisher, index.NewUserFileIndex())

		err := svc.DeleteFile(context.Background(), "tenant-1", "file-1", "")
		require.NoError(t, err)
		assert.Equal(t, []string{"file-1"}, repo.deletedIDs)
		assert.Equal(t, []string{"tenant-1"}, publisher.deletedBy)
	})
}

func TestUpdateFileMetadata(t *testing.T) {
	t.Run("validates and trims inputs", func(t *testing.T) {
		repo := newFakeFileRepo()
		updated := sampleStoredFile("tenant-1", "file-1", "actor-1")
		repo.updateMetadataResult = updated
		publisher := &fakeEventPublisher{}
		svc := newStorageService(repo, &fakeObjectStore{}, publisher, index.NewUserFileIndex())

		filename := "  renamed.pdf  "
		contentType := " text/plain "
		referenceID := "   "
		referenceType := " claim "
		isPublic := true
		fileType := storageentityv1.FileType_FILE_TYPE_IMAGE

		file, err := svc.UpdateFileMetadata(context.Background(), &UpdateFileInput{
			TenantID:      "tenant-1",
			FileID:        "file-1",
			Filename:      &filename,
			ContentType:   &contentType,
			FileType:      &fileType,
			ReferenceID:   &referenceID,
			ReferenceType: &referenceType,
			IsPublic:      &isPublic,
			ClearExpires:  true,
			UpdatedBy:     " actor-1 ",
		})
		require.NoError(t, err)
		require.NotNil(t, file)
		require.NotNil(t, repo.lastPatch)
		require.NotNil(t, repo.lastPatch.Filename)
		assert.Equal(t, "renamed.pdf", *repo.lastPatch.Filename)
		require.NotNil(t, repo.lastPatch.ContentType)
		assert.Equal(t, "text/plain", *repo.lastPatch.ContentType)
		require.NotNil(t, repo.lastPatch.ReferenceID)
		assert.Equal(t, "", *repo.lastPatch.ReferenceID)
		require.NotNil(t, repo.lastPatch.ReferenceType)
		assert.Equal(t, "claim", *repo.lastPatch.ReferenceType)
		require.NotNil(t, repo.lastPatch.UploadedBy)
		assert.Equal(t, "actor-1", *repo.lastPatch.UploadedBy)
		assert.True(t, repo.lastPatch.ClearExpires)
		assert.Equal(t, []string{"actor-1"}, publisher.metadataUpdatedBy)
		assert.Equal(t, []string{"filename", "content_type", "file_type", "reference_id", "reference_type", "is_public", "expires_at"}, publisher.metadataUpdated[0])
	})

	t.Run("maps repository sentinels", func(t *testing.T) {
		repo := newFakeFileRepo()
		svc := newStorageService(repo, &fakeObjectStore{}, nil, index.NewUserFileIndex())

		_, err := svc.UpdateFileMetadata(context.Background(), nil)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNoMetadataUpdates)

		blank := "   "
		_, err = svc.UpdateFileMetadata(context.Background(), &UpdateFileInput{
			TenantID: "tenant-1",
			FileID:   "file-1",
			Filename: &blank,
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInput)

		repo.updateMetadataErr = repository.ErrFileNotFound
		name := "ok.pdf"
		_, err = svc.UpdateFileMetadata(context.Background(), &UpdateFileInput{
			TenantID: "tenant-1",
			FileID:   "file-1",
			Filename: &name,
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrFileNotFound)

		repo.updateMetadataErr = repository.ErrNoMetadataUpdates
		_, err = svc.UpdateFileMetadata(context.Background(), &UpdateFileInput{
			TenantID: "tenant-1",
			FileID:   "file-1",
			Filename: &name,
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNoMetadataUpdates)
	})
}

func TestPresignedUploadFinalizeAndDownload(t *testing.T) {
	t.Run("issues upload URL with defaults", func(t *testing.T) {
		repo := newFakeFileRepo()
		store := &fakeObjectStore{
			bucket:             "bucket-1",
			presignedUploadURL: "https://upload-url",
			objectURL:          "https://objects/file",
			cdnURL:             "https://cdn/file",
		}
		publisher := &fakeEventPublisher{}
		svc := newStorageService(repo, store, publisher, index.NewUserFileIndex())

		uploadURL, fileID, key, err := svc.GetPresignedUploadURL(
			context.Background(),
			"tenant-1",
			"doc.pdf",
			"",
			storageentityv1.FileType_FILE_TYPE_DOCUMENT,
			"ref-1",
			"claim",
			true,
			0,
			"",
		)
		require.NoError(t, err)
		assert.Equal(t, "https://upload-url", uploadURL)
		assert.NotEmpty(t, fileID)
		assert.NotEmpty(t, key)
		assert.Equal(t, 15*time.Minute, store.presignedUploadExpiry)
		assert.Equal(t, []string{"tenant-1"}, publisher.uploadURLIssuedBy)
	})

	t.Run("finalizes direct upload and falls back to detected content type", func(t *testing.T) {
		repo := newFakeFileRepo()
		repo.files["file-1"] = sampleStoredFile("tenant-1", "file-1", "user-1")
		store := &fakeObjectStore{headSize: 123, headContentType: "image/png"}
		publisher := &fakeEventPublisher{}
		svc := newStorageService(repo, store, publisher, index.NewUserFileIndex())

		file, err := svc.FinalizeDirectUpload(
			context.Background(),
			"tenant-1",
			"file-1",
			"photo.png",
			"",
			storageentityv1.FileType_FILE_TYPE_IMAGE,
			"ref-2",
			"claim",
			false,
			nil,
			"",
		)
		require.NoError(t, err)
		assert.Equal(t, int64(123), file.SizeBytes)
		assert.Equal(t, "image/png", file.ContentType)
		assert.Equal(t, "tenant-1", publisher.finalizedBy[0])
		assert.Equal(t, "FINALIZE", publisher.uploadedSources[0])
	})

	t.Run("generates download URL with default expiry", func(t *testing.T) {
		repo := newFakeFileRepo()
		repo.files["file-1"] = sampleStoredFile("tenant-1", "file-1", "user-1")
		store := &fakeObjectStore{presignedDownloadURL: "https://download-url"}
		svc := newStorageService(repo, store, nil, index.NewUserFileIndex())

		url, expiresAt, err := svc.GetPresignedDownloadURL(context.Background(), "tenant-1", "file-1", 0)
		require.NoError(t, err)
		assert.Equal(t, "https://download-url", url)
		require.NotNil(t, expiresAt)
		assert.Equal(t, 60*time.Minute, store.presignedDownloadExp)
	})
}

func TestListFilesAndFilterPagination(t *testing.T) {
	t.Run("uses repository listing when uploader filter is absent", func(t *testing.T) {
		repo := newFakeFileRepo()
		repo.listResult = []*storageentityv1.StoredFile{sampleStoredFile("tenant-1", "file-1", "user-1")}
		repo.listTotal = 1
		svc := newStorageService(repo, &fakeObjectStore{}, nil, index.NewUserFileIndex())

		files, total, err := svc.ListFiles(context.Background(), "tenant-1", storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "", "", "", 10, 0)
		require.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, 1, total)
	})

	t.Run("hits warmed user index before repository", func(t *testing.T) {
		repo := newFakeFileRepo()
		userIndex := index.NewUserFileIndex()
		userIndex.WarmUser("tenant-1", "user-1", []*storageentityv1.StoredFile{
			sampleStoredFile("tenant-1", "file-1", "user-1"),
			sampleStoredFile("tenant-1", "file-2", "user-1"),
		})
		svc := newStorageService(repo, &fakeObjectStore{}, nil, userIndex)

		files, total, err := svc.ListFiles(context.Background(), "tenant-1", storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "", "", "user-1", 1, 0)
		require.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, 2, total)
		assert.Zero(t, repo.listAllCalls)
	})

	t.Run("warms index from repository when needed", func(t *testing.T) {
		repo := newFakeFileRepo()
		repo.listAllResult = []*storageentityv1.StoredFile{
			sampleStoredFile("tenant-1", "file-1", "user-2"),
			sampleStoredFile("tenant-1", "file-2", "user-2"),
			sampleStoredFile("tenant-1", "file-3", "user-2"),
		}
		repo.listAllResult[0].ReferenceId = "ref-1"
		repo.listAllResult[1].ReferenceId = "ref-1"
		repo.listAllResult[2].ReferenceId = "ref-2"
		svc := newStorageService(repo, &fakeObjectStore{}, nil, index.NewUserFileIndex())

		files, total, err := svc.ListFiles(context.Background(), "tenant-1", storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "ref-1", "", "user-2", 1, 1)
		require.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, 2, total)
		assert.Equal(t, 1, repo.listAllCalls)
	})

	t.Run("filters and paginates robustly", func(t *testing.T) {
		files, total := filterAndPaginateFiles([]*storageentityv1.StoredFile{
			nil,
			sampleStoredFile("tenant-1", "file-1", "user-1"),
			sampleStoredFile("tenant-1", "file-2", "user-1"),
		}, storageentityv1.FileType_FILE_TYPE_DOCUMENT, "", "", 5, -2)
		assert.Len(t, files, 2)
		assert.Equal(t, 2, total)
	})
}

func sampleStoredFile(tenantID, fileID, uploadedBy string) *storageentityv1.StoredFile {
	return &storageentityv1.StoredFile{
		FileId:        fileID,
		TenantId:      tenantID,
		Filename:      fileID + ".pdf",
		ContentType:   "application/pdf",
		SizeBytes:     64,
		StorageKey:    "keys/" + fileID,
		Bucket:        "bucket-1",
		Url:           "https://objects/" + fileID,
		CdnUrl:        "https://cdn/" + fileID,
		FileType:      storageentityv1.FileType_FILE_TYPE_DOCUMENT,
		ReferenceId:   "ref-1",
		ReferenceType: "claim",
		UploadedBy:    uploadedBy,
		CreatedAt:     timestamppb.New(time.Now().UTC()),
	}
}

func cloneStoredFile(file *storageentityv1.StoredFile) *storageentityv1.StoredFile {
	if file == nil {
		return nil
	}
	return proto.Clone(file).(*storageentityv1.StoredFile)
}

func cloneStoredFiles(files []*storageentityv1.StoredFile) []*storageentityv1.StoredFile {
	out := make([]*storageentityv1.StoredFile, 0, len(files))
	for _, file := range files {
		out = append(out, cloneStoredFile(file))
	}
	return out
}
