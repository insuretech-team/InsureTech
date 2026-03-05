package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
)

var testTenantID = getTestTenantID()

func getTestTenantID() string {
	if tenantID := os.Getenv("DEFAULT_TENANT_ID"); tenantID != "" {
		return tenantID
	}
	return "00000000-0000-0000-0000-000000000001" // Fallback UUID for testing
}

func TestFileRepository_Create(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx, testTenantID)

	// Setup
	sqlDB, err := getTestDB().DB()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")
	repo := NewFileRepository(sqlxDB)

	// Create test file
	fileID := uuid.New().String()
	uploadedBy := uuid.New().String()

	file := &storageentityv1.StoredFile{
		FileId:        fileID,
		TenantId:      testTenantID,
		Filename:      "test-image.jpg",
		ContentType:   "image/jpeg",
		SizeBytes:     12345,
		StorageKey:    "lpc/assets/test-tenant-123/2025/01/test.jpg",
		Bucket:        "merchbd",
		Url:           "https://merchbd.sgp1.digitaloceanspaces.com/lpc/assets/test-tenant-123/2025/01/test.jpg",
		CdnUrl:        "https://merchbd.sgp1.cdn.digitaloceanspaces.com/lpc/assets/test-tenant-123/2025/01/test.jpg",
		FileType:      storageentityv1.FileType_FILE_TYPE_IMAGE,
		ReferenceId:   "",
		ReferenceType: "",
		IsPublic:      true,
		UploadedBy:    uploadedBy,
	}

	// Test Create
	created, err := repo.Create(ctx, testTenantID, file)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotNil(t, created.CreatedAt)
	assert.NotNil(t, created.UpdatedAt)
	assert.Equal(t, fileID, created.FileId)
	assert.Equal(t, "test-image.jpg", created.Filename)
}

func TestFileRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx, testTenantID)

	// Setup
	sqlDB, err := getTestDB().DB()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")
	repo := NewFileRepository(sqlxDB)

	// Create test file first
	fileID := uuid.New().String()
	uploadedBy := uuid.New().String()

	file := &storageentityv1.StoredFile{
		FileId:        fileID,
		TenantId:      testTenantID,
		Filename:      "test-document.pdf",
		ContentType:   "application/pdf",
		SizeBytes:     54321,
		StorageKey:    "lpc/orders/test-tenant-123/2025/01/invoice.pdf",
		Bucket:        "merchbd",
		Url:           "https://merchbd.sgp1.digitaloceanspaces.com/lpc/orders/test-tenant-123/2025/01/invoice.pdf",
		FileType:      storageentityv1.FileType_FILE_TYPE_INVOICE,
		ReferenceId:   uuid.New().String(),
		ReferenceType: "order",
		IsPublic:      false,
		UploadedBy:    uploadedBy,
	}

	_, err = repo.Create(ctx, testTenantID, file)
	require.NoError(t, err)

	// Test GetByID
	retrieved, err := repo.GetByID(ctx, testTenantID, fileID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, fileID, retrieved.FileId)
	assert.Equal(t, "test-document.pdf", retrieved.Filename)
	assert.Equal(t, storageentityv1.FileType_FILE_TYPE_INVOICE, retrieved.FileType)
	assert.NotEmpty(t, retrieved.ReferenceId)
}

func TestFileRepository_List(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx, testTenantID)

	// Setup
	sqlDB, err := getTestDB().DB()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")
	repo := NewFileRepository(sqlxDB)

	// Create multiple test files
	uploadedBy := uuid.New().String()
	orderRefID := uuid.New().String()

	files := []*storageentityv1.StoredFile{
		{
			FileId:      uuid.New().String(),
			TenantId:    testTenantID,
			Filename:    "image1.jpg",
			ContentType: "image/jpeg",
			SizeBytes:   1000,
			StorageKey:  "lpc/assets/test-tenant-123/2025/01/image1.jpg",
			Bucket:      "merchbd",
			Url:         "https://merchbd.sgp1.digitaloceanspaces.com/lpc/assets/test-tenant-123/2025/01/image1.jpg",
			FileType:    storageentityv1.FileType_FILE_TYPE_IMAGE,
			IsPublic:    true,
			UploadedBy:  uploadedBy,
		},
		{
			FileId:      uuid.New().String(),
			TenantId:    testTenantID,
			Filename:    "image2.jpg",
			ContentType: "image/jpeg",
			SizeBytes:   2000,
			StorageKey:  "lpc/assets/test-tenant-123/2025/01/image2.jpg",
			Bucket:      "merchbd",
			Url:         "https://merchbd.sgp1.digitaloceanspaces.com/lpc/assets/test-tenant-123/2025/01/image2.jpg",
			FileType:    storageentityv1.FileType_FILE_TYPE_IMAGE,
			IsPublic:    true,
			UploadedBy:  uploadedBy,
		},
		{
			FileId:        uuid.New().String(),
			TenantId:      testTenantID,
			Filename:      "invoice.pdf",
			ContentType:   "application/pdf",
			SizeBytes:     3000,
			StorageKey:    "lpc/orders/test-tenant-123/2025/01/invoice.pdf",
			Bucket:        "merchbd",
			Url:           "https://merchbd.sgp1.digitaloceanspaces.com/lpc/orders/test-tenant-123/2025/01/invoice.pdf",
			FileType:      storageentityv1.FileType_FILE_TYPE_INVOICE,
			ReferenceId:   orderRefID,
			ReferenceType: "order",
			IsPublic:      false,
			UploadedBy:    uploadedBy,
		},
	}

	for _, file := range files {
		_, err := repo.Create(ctx, testTenantID, file)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // Ensure different created_at times
	}

	// Test List all files
	allFiles, total, err := repo.List(ctx, testTenantID, storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "", "", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, allFiles, 3)

	// Test List with file type filter
	imageFiles, total, err := repo.List(ctx, testTenantID, storageentityv1.FileType_FILE_TYPE_IMAGE, "", "", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, imageFiles, 2)

	// Test List with reference ID filter
	orderFiles, total, err := repo.List(ctx, testTenantID, storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, orderRefID, "", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, orderFiles, 1)
	assert.Equal(t, "invoice.pdf", orderFiles[0].Filename)

	// Test pagination
	page1, _, err := repo.List(ctx, testTenantID, storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "", "", 2, 0)
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, _, err := repo.List(ctx, testTenantID, storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "", "", 2, 2)
	require.NoError(t, err)
	assert.Len(t, page2, 1)
}

func TestFileRepository_Delete(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx, testTenantID)

	// Setup
	sqlDB, err := getTestDB().DB()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")
	repo := NewFileRepository(sqlxDB)

	// Create test file
	fileID := uuid.New().String()
	uploadedBy := uuid.New().String()

	file := &storageentityv1.StoredFile{
		FileId:      fileID,
		TenantId:    testTenantID,
		Filename:    "to-delete.jpg",
		ContentType: "image/jpeg",
		SizeBytes:   1000,
		StorageKey:  "lpc/assets/test-tenant-123/2025/01/to-delete.jpg",
		Bucket:      "merchbd",
		Url:         "https://merchbd.sgp1.digitaloceanspaces.com/lpc/assets/test-tenant-123/2025/01/to-delete.jpg",
		FileType:    storageentityv1.FileType_FILE_TYPE_IMAGE,
		IsPublic:    true,
		UploadedBy:  uploadedBy,
	}

	_, err = repo.Create(ctx, testTenantID, file)
	require.NoError(t, err)

	// Test Delete
	err = repo.Delete(ctx, testTenantID, fileID)
	require.NoError(t, err)

	// Verify deleted
	_, err = repo.GetByID(ctx, testTenantID, fileID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestFileRepository_WithExpiration(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx, testTenantID)

	// Setup
	sqlDB, err := getTestDB().DB()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")
	repo := NewFileRepository(sqlxDB)

	// Create test file with expiration
	fileID := uuid.New().String()
	uploadedBy := uuid.New().String()
	expiresAt := timestamppb.New(time.Now().Add(24 * time.Hour))

	file := &storageentityv1.StoredFile{
		FileId:      fileID,
		TenantId:    testTenantID,
		Filename:    "temp-file.jpg",
		ContentType: "image/jpeg",
		SizeBytes:   1000,
		StorageKey:  "lpc/assets/test-tenant-123/2025/01/temp.jpg",
		Bucket:      "merchbd",
		Url:         "https://merchbd.sgp1.digitaloceanspaces.com/lpc/assets/test-tenant-123/2025/01/temp.jpg",
		FileType:    storageentityv1.FileType_FILE_TYPE_IMAGE,
		IsPublic:    false,
		ExpiresAt:   expiresAt,
		UploadedBy:  uploadedBy,
	}

	created, err := repo.Create(ctx, testTenantID, file)
	require.NoError(t, err)
	assert.NotNil(t, created.ExpiresAt)

	// Retrieve and verify expiration is preserved
	retrieved, err := repo.GetByID(ctx, testTenantID, fileID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved.ExpiresAt)
	assert.WithinDuration(t, expiresAt.AsTime(), retrieved.ExpiresAt.AsTime(), time.Second)
}
