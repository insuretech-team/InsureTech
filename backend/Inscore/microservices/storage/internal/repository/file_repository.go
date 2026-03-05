package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"

	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
)

// FileRepository handles database operations for stored files
type FileRepository struct {
	db *sqlx.DB
}

var (
	ErrFileNotFound      = errors.New("file not found")
	ErrNoMetadataUpdates = errors.New("no metadata fields provided")
)

// FileMetadataPatch contains mutable metadata fields for partial updates.
type FileMetadataPatch struct {
	Filename      *string
	ContentType   *string
	FileType      *storageentityv1.FileType
	ReferenceID   *string
	ReferenceType *string
	IsPublic      *bool
	ExpiresAt     *timestamppb.Timestamp
	ClearExpires  bool
	UploadedBy    *string
}

func enumToDBValue(fileType storageentityv1.FileType) string {
	return fileType.String()
}

func dbValueToEnum(v string) storageentityv1.FileType {
	s := strings.TrimSpace(v)
	if s == "" {
		return storageentityv1.FileType_FILE_TYPE_UNSPECIFIED
	}
	if n, err := strconv.Atoi(s); err == nil {
		ft := storageentityv1.FileType(n)
		if _, ok := storageentityv1.FileType_name[int32(ft)]; ok {
			return ft
		}
	}
	if n, ok := storageentityv1.FileType_value[s]; ok {
		return storageentityv1.FileType(n)
	}
	return storageentityv1.FileType_FILE_TYPE_UNSPECIFIED
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *sqlx.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create stores a new file record
func (r *FileRepository) Create(ctx context.Context, tenantID string, file *storageentityv1.StoredFile) (*storageentityv1.StoredFile, error) {
	query := `
		INSERT INTO storage_schema.files (
			file_id, tenant_id, filename, content_type, size_bytes, storage_key, bucket,
			url, cdn_url, file_type, reference_id, reference_type, is_public, 
			expires_at, uploaded_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt sql.NullTime
	var expiresAt sql.NullTime
	if file.ExpiresAt != nil {
		expiresAt = sql.NullTime{Time: file.ExpiresAt.AsTime(), Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query,
		file.FileId,
		tenantID,
		file.Filename,
		file.ContentType,
		file.SizeBytes,
		file.StorageKey,
		file.Bucket,
		file.Url,
		file.CdnUrl,
		enumToDBValue(file.FileType),
		file.ReferenceId,
		file.ReferenceType,
		file.IsPublic,
		expiresAt,
		file.UploadedBy,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	if createdAt.Valid {
		file.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		file.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return file, nil
}

// GetByID retrieves a file by ID
func (r *FileRepository) GetByID(ctx context.Context, tenantID string, fileID string) (*storageentityv1.StoredFile, error) {
	query := `
		SELECT file_id, tenant_id, filename, content_type, size_bytes, storage_key, bucket,
			   url, cdn_url, file_type, reference_id, reference_type, is_public,
			   expires_at, uploaded_by, created_at, updated_at
		FROM storage_schema.files
		WHERE tenant_id = $1 AND file_id = $2
	`

	var file storageentityv1.StoredFile
	var fileID_, tenantID_, uploadedBy_ string
	var fileType_ string
	var createdAt, updatedAt, expiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tenantID, fileID).Scan(
		&fileID_,
		&tenantID_,
		&file.Filename,
		&file.ContentType,
		&file.SizeBytes,
		&file.StorageKey,
		&file.Bucket,
		&file.Url,
		&file.CdnUrl,
		&fileType_,
		&file.ReferenceId,
		&file.ReferenceType,
		&file.IsPublic,
		&expiresAt,
		&uploadedBy_,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	file.FileId = fileID_
	file.TenantId = tenantID_
	file.FileType = dbValueToEnum(fileType_)
	file.UploadedBy = uploadedBy_

	if createdAt.Valid {
		file.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		file.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	if expiresAt.Valid {
		file.ExpiresAt = timestamppb.New(expiresAt.Time)
	}

	return &file, nil
}

// List retrieves files with filters
func (r *FileRepository) List(ctx context.Context, tenantID string, fileType storageentityv1.FileType, referenceID string, referenceType string, limit, offset int32) ([]*storageentityv1.StoredFile, int, error) {
	// Build query with filters
	baseQuery := `
		SELECT file_id, tenant_id, filename, content_type, size_bytes, storage_key, bucket,
			   url, cdn_url, file_type, reference_id, reference_type, is_public,
			   expires_at, uploaded_by, created_at, updated_at
		FROM storage_schema.files
		WHERE tenant_id = $1
	`
	countQuery := `SELECT COUNT(*) FROM storage_schema.files WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	argCount := 2

	// Add filters
	if fileType != storageentityv1.FileType_FILE_TYPE_UNSPECIFIED {
		baseQuery += fmt.Sprintf(" AND file_type = $%d", argCount)
		countQuery += fmt.Sprintf(" AND file_type = $%d", argCount)
		args = append(args, enumToDBValue(fileType))
		argCount++
	}

	if referenceID != "" {
		baseQuery += fmt.Sprintf(" AND reference_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND reference_id = $%d", argCount)
		args = append(args, referenceID)
		argCount++
	}

	if referenceType != "" {
		baseQuery += fmt.Sprintf(" AND reference_type = $%d", argCount)
		countQuery += fmt.Sprintf(" AND reference_type = $%d", argCount)
		args = append(args, referenceType)
		argCount++
	}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}

	// Add pagination
	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}
	defer rows.Close()

	var files []*storageentityv1.StoredFile
	for rows.Next() {
		var file storageentityv1.StoredFile
		var fileID_, tenantID_, uploadedBy_ string
		var fileType_ string
		var createdAt, updatedAt, expiresAt sql.NullTime

		err := rows.Scan(
			&fileID_,
			&tenantID_,
			&file.Filename,
			&file.ContentType,
			&file.SizeBytes,
			&file.StorageKey,
			&file.Bucket,
			&file.Url,
			&file.CdnUrl,
			&fileType_,
			&file.ReferenceId,
			&file.ReferenceType,
			&file.IsPublic,
			&expiresAt,
			&uploadedBy_,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan file: %w", err)
		}

		file.FileId = fileID_
		file.TenantId = tenantID_
		file.FileType = dbValueToEnum(fileType_)
		file.UploadedBy = uploadedBy_

		if createdAt.Valid {
			file.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			file.UpdatedAt = timestamppb.New(updatedAt.Time)
		}
		if expiresAt.Valid {
			file.ExpiresAt = timestamppb.New(expiresAt.Time)
		}

		files = append(files, &file)
	}

	return files, total, nil
}

// ListAllByUploadedBy retrieves all files for a tenant and uploader, ordered by newest first.
func (r *FileRepository) ListAllByUploadedBy(ctx context.Context, tenantID string, uploadedBy string) ([]*storageentityv1.StoredFile, error) {
	query := `
		SELECT file_id, tenant_id, filename, content_type, size_bytes, storage_key, bucket,
			   url, cdn_url, file_type, reference_id, reference_type, is_public,
			   expires_at, uploaded_by, created_at, updated_at
		FROM storage_schema.files
		WHERE tenant_id = $1 AND uploaded_by = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, uploadedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to list files by uploader: %w", err)
	}
	defer rows.Close()

	files := make([]*storageentityv1.StoredFile, 0, 64)
	for rows.Next() {
		var file storageentityv1.StoredFile
		var fileID_, tenantID_, uploadedBy_ string
		var fileType_ string
		var createdAt, updatedAt, expiresAt sql.NullTime

		err := rows.Scan(
			&fileID_,
			&tenantID_,
			&file.Filename,
			&file.ContentType,
			&file.SizeBytes,
			&file.StorageKey,
			&file.Bucket,
			&file.Url,
			&file.CdnUrl,
			&fileType_,
			&file.ReferenceId,
			&file.ReferenceType,
			&file.IsPublic,
			&expiresAt,
			&uploadedBy_,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		file.FileId = fileID_
		file.TenantId = tenantID_
		file.FileType = dbValueToEnum(fileType_)
		file.UploadedBy = uploadedBy_

		if createdAt.Valid {
			file.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			file.UpdatedAt = timestamppb.New(updatedAt.Time)
		}
		if expiresAt.Valid {
			file.ExpiresAt = timestamppb.New(expiresAt.Time)
		}

		files = append(files, &file)
	}

	return files, nil
}

// Delete removes a file record
func (r *FileRepository) Delete(ctx context.Context, tenantID string, fileID string) error {
	query := `DELETE FROM storage_schema.files WHERE tenant_id = $1 AND file_id = $2`

	result, err := r.db.ExecContext(ctx, query, tenantID, fileID)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrFileNotFound
	}

	return nil
}

// UpdateAfterDirectUpload updates metadata for a file previously created via presigned upload.
func (r *FileRepository) UpdateAfterDirectUpload(ctx context.Context, tenantID string, file *storageentityv1.StoredFile) (*storageentityv1.StoredFile, error) {
	query := `
		UPDATE storage_schema.files
		SET
			filename = $3,
			content_type = $4,
			size_bytes = $5,
			file_type = $6,
			reference_id = $7,
			reference_type = $8,
			is_public = $9,
			expires_at = $10,
			uploaded_by = $11,
			updated_at = NOW()
		WHERE tenant_id = $1 AND file_id = $2
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt sql.NullTime
	var expiresAt sql.NullTime
	if file.ExpiresAt != nil {
		expiresAt = sql.NullTime{Time: file.ExpiresAt.AsTime(), Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query,
		tenantID,
		file.FileId,
		file.Filename,
		file.ContentType,
		file.SizeBytes,
		enumToDBValue(file.FileType),
		file.ReferenceId,
		file.ReferenceType,
		file.IsPublic,
		expiresAt,
		file.UploadedBy,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to finalize file metadata: %w", err)
	}

	if createdAt.Valid {
		file.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		file.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return file, nil
}

// UpdateMetadata updates mutable file metadata fields and returns the updated file.
func (r *FileRepository) UpdateMetadata(ctx context.Context, tenantID, fileID string, patch *FileMetadataPatch) (*storageentityv1.StoredFile, error) {
	if patch == nil {
		return nil, ErrNoMetadataUpdates
	}

	setClauses := make([]string, 0, 10)
	args := []any{tenantID, fileID}
	argPos := 3

	if patch.Filename != nil {
		setClauses = append(setClauses, fmt.Sprintf("filename = $%d", argPos))
		args = append(args, *patch.Filename)
		argPos++
	}
	if patch.ContentType != nil {
		setClauses = append(setClauses, fmt.Sprintf("content_type = $%d", argPos))
		args = append(args, *patch.ContentType)
		argPos++
	}
	if patch.FileType != nil {
		setClauses = append(setClauses, fmt.Sprintf("file_type = $%d", argPos))
		args = append(args, enumToDBValue(*patch.FileType))
		argPos++
	}
	if patch.ReferenceID != nil {
		setClauses = append(setClauses, fmt.Sprintf("reference_id = $%d", argPos))
		args = append(args, *patch.ReferenceID)
		argPos++
	}
	if patch.ReferenceType != nil {
		setClauses = append(setClauses, fmt.Sprintf("reference_type = $%d", argPos))
		args = append(args, *patch.ReferenceType)
		argPos++
	}
	if patch.IsPublic != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_public = $%d", argPos))
		args = append(args, *patch.IsPublic)
		argPos++
	}
	if patch.ClearExpires {
		setClauses = append(setClauses, "expires_at = NULL")
	} else if patch.ExpiresAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("expires_at = $%d", argPos))
		args = append(args, sql.NullTime{Time: patch.ExpiresAt.AsTime(), Valid: true})
		argPos++
	}
	if patch.UploadedBy != nil {
		setClauses = append(setClauses, fmt.Sprintf("uploaded_by = $%d", argPos))
		args = append(args, *patch.UploadedBy)
		argPos++
	}

	if len(setClauses) == 0 {
		return nil, ErrNoMetadataUpdates
	}

	query := fmt.Sprintf(`
		UPDATE storage_schema.files
		SET %s, updated_at = NOW()
		WHERE tenant_id = $1 AND file_id = $2
	`, strings.Join(setClauses, ", "))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update file metadata: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrFileNotFound
	}

	return r.GetByID(ctx, tenantID, fileID)
}
