package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	mediav1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MediaRepository handles database operations for media files.
type MediaRepository struct {
	db *sqlx.DB
}

var (
	ErrMediaNotFound = errors.New("media file not found")
	ErrNoUpdates     = errors.New("no fields to update")
)

// NewMediaRepository creates a new media repository.
func NewMediaRepository(db *sqlx.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

// Create stores a new media file record.
func (r *MediaRepository) Create(ctx context.Context, media *mediav1.MediaFile) (*mediav1.MediaFile, error) {
	query := `
		INSERT INTO media_schema.media_files (
			media_id, file_id, tenant_id, media_type, mime_type, file_size_bytes,
			width, height, dpi, optimized_file_id, thumbnail_file_id,
			entity_type, entity_id, ocr_text, validation_status, validation_errors,
			virus_scan_status, uploaded_by, audit_info, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query,
		media.Id,
		media.FileId,
		media.TenantId,
		media.MediaType.String(),
		media.MimeType,
		media.FileSizeBytes,
		nullInt32(media.Width),
		nullInt32(media.Height),
		nullInt32(media.Dpi),
		nullString(media.OptimizedFileId),
		nullString(media.ThumbnailFileId),
		nullString(media.EntityType),
		nullString(media.EntityId),
		nullString(media.OcrText),
		media.ValidationStatus.String(),
		nullString(media.ValidationErrors),
		media.VirusScanStatus.String(),
		media.UploadedBy,
		"{}",
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create media file: %w", err)
	}

	if media.AuditInfo == nil {
		media.AuditInfo = &commonv1.AuditInfo{}
	}
	if createdAt.Valid {
		media.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		media.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return media, nil
}

// GetByID retrieves a media file by ID, scoped by tenant when tenantID is provided.
func (r *MediaRepository) GetByID(ctx context.Context, tenantID, mediaID string) (*mediav1.MediaFile, error) {
	query := `
		SELECT
			media_id, file_id, tenant_id, media_type, mime_type, file_size_bytes,
			width, height, dpi, optimized_file_id, thumbnail_file_id,
			entity_type, entity_id, ocr_text, validation_status, validation_errors,
			virus_scan_status, uploaded_by, audit_info, created_at, updated_at
		FROM media_schema.media_files
		WHERE media_id = $1`
	args := []any{mediaID}
	if strings.TrimSpace(tenantID) != "" {
		query += " AND tenant_id = $2"
		args = append(args, tenantID)
	}

	var media mediav1.MediaFile
	var mediaType, validationStatus, virusScanStatus string
	var width, height, dpi sql.NullInt32
	var optimizedFileID, thumbnailFileID, entityType, entityID, ocrText, validationErrors sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&media.Id,
		&media.FileId,
		&media.TenantId,
		&mediaType,
		&media.MimeType,
		&media.FileSizeBytes,
		&width,
		&height,
		&dpi,
		&optimizedFileID,
		&thumbnailFileID,
		&entityType,
		&entityID,
		&ocrText,
		&validationStatus,
		&validationErrors,
		&virusScanStatus,
		&media.UploadedBy,
		&media.AuditInfo,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrMediaNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get media file: %w", err)
	}

	media.MediaType = parseMediaType(mediaType)
	media.ValidationStatus = parseValidationStatus(validationStatus)
	media.VirusScanStatus = parseVirusScanStatus(virusScanStatus)

	if width.Valid {
		media.Width = width.Int32
	}
	if height.Valid {
		media.Height = height.Int32
	}
	if dpi.Valid {
		media.Dpi = dpi.Int32
	}
	if optimizedFileID.Valid {
		media.OptimizedFileId = optimizedFileID.String
	}
	if thumbnailFileID.Valid {
		media.ThumbnailFileId = thumbnailFileID.String
	}
	if entityType.Valid {
		media.EntityType = entityType.String
	}
	if entityID.Valid {
		media.EntityId = entityID.String
	}
	if ocrText.Valid {
		media.OcrText = ocrText.String
	}
	if validationErrors.Valid {
		media.ValidationErrors = validationErrors.String
	}

	if media.AuditInfo == nil {
		media.AuditInfo = &commonv1.AuditInfo{}
	}
	if createdAt.Valid {
		media.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		media.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &media, nil
}

// ListByEntity retrieves media files for a specific entity with optional filters.
func (r *MediaRepository) ListByEntity(
	ctx context.Context,
	tenantID, entityType, entityID string,
	mediaType *mediav1.MediaType,
	validationStatus *mediav1.ValidationStatus,
	limit, offset int,
) ([]*mediav1.MediaFile, int, error) {
	conditions := []string{"entity_type = $1", "entity_id = $2"}
	args := []any{entityType, entityID}

	if strings.TrimSpace(tenantID) != "" {
		conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", len(args)+1))
		args = append(args, tenantID)
	}
	if mediaType != nil {
		conditions = append(conditions, fmt.Sprintf("media_type = $%d", len(args)+1))
		args = append(args, mediaType.String())
	}
	if validationStatus != nil {
		conditions = append(conditions, fmt.Sprintf("validation_status = $%d", len(args)+1))
		args = append(args, validationStatus.String())
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM media_schema.media_files
		WHERE %s
	`, whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count media files: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT
			media_id, file_id, tenant_id, media_type, mime_type, file_size_bytes,
			width, height, dpi, optimized_file_id, thumbnail_file_id,
			entity_type, entity_id, ocr_text, validation_status, validation_errors,
			virus_scan_status, uploaded_by, created_at, updated_at
		FROM media_schema.media_files
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(args)+1, len(args)+2)

	queryArgs := append(append([]any{}, args...), limit, offset)
	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list media files: %w", err)
	}
	defer rows.Close()

	var mediaFiles []*mediav1.MediaFile
	for rows.Next() {
		var media mediav1.MediaFile
		var mediaTypeStr, validationStatusStr, virusScanStatusStr string
		var width, height, dpi sql.NullInt32
		var optimizedFileID, thumbnailFileID, entityTypeVal, entityIDVal, ocrText, validationErrors sql.NullString
		var createdAt, updatedAt sql.NullTime

		if err := rows.Scan(
			&media.Id,
			&media.FileId,
			&media.TenantId,
			&mediaTypeStr,
			&media.MimeType,
			&media.FileSizeBytes,
			&width,
			&height,
			&dpi,
			&optimizedFileID,
			&thumbnailFileID,
			&entityTypeVal,
			&entityIDVal,
			&ocrText,
			&validationStatusStr,
			&validationErrors,
			&virusScanStatusStr,
			&media.UploadedBy,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan media file: %w", err)
		}

		media.MediaType = parseMediaType(mediaTypeStr)
		media.ValidationStatus = parseValidationStatus(validationStatusStr)
		media.VirusScanStatus = parseVirusScanStatus(virusScanStatusStr)

		if width.Valid {
			media.Width = width.Int32
		}
		if height.Valid {
			media.Height = height.Int32
		}
		if dpi.Valid {
			media.Dpi = dpi.Int32
		}
		if optimizedFileID.Valid {
			media.OptimizedFileId = optimizedFileID.String
		}
		if thumbnailFileID.Valid {
			media.ThumbnailFileId = thumbnailFileID.String
		}
		if entityTypeVal.Valid {
			media.EntityType = entityTypeVal.String
		}
		if entityIDVal.Valid {
			media.EntityId = entityIDVal.String
		}
		if ocrText.Valid {
			media.OcrText = ocrText.String
		}
		if validationErrors.Valid {
			media.ValidationErrors = validationErrors.String
		}

		media.AuditInfo = &commonv1.AuditInfo{}
		if createdAt.Valid {
			media.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			media.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		mediaFiles = append(mediaFiles, &media)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating media files: %w", err)
	}

	return mediaFiles, total, nil
}

// UpdateValidationStatus updates the validation status of a media file.
func (r *MediaRepository) UpdateValidationStatus(ctx context.Context, tenantID, mediaID string, status mediav1.ValidationStatus, errors string) error {
	query := `
		UPDATE media_schema.media_files
		SET validation_status = $1, validation_errors = $2, updated_at = NOW()
		WHERE media_id = $3`
	args := []any{status.String(), nullString(errors), mediaID}
	if strings.TrimSpace(tenantID) != "" {
		query += " AND tenant_id = $4"
		args = append(args, tenantID)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update validation status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrMediaNotFound
	}

	return nil
}

// UpdateVirusScanStatus updates the virus scan status.
func (r *MediaRepository) UpdateVirusScanStatus(ctx context.Context, tenantID, mediaID string, status mediav1.VirusScanStatus) error {
	query := `
		UPDATE media_schema.media_files
		SET virus_scan_status = $1, updated_at = NOW()
		WHERE media_id = $2`
	args := []any{status.String(), mediaID}
	if strings.TrimSpace(tenantID) != "" {
		query += " AND tenant_id = $3"
		args = append(args, tenantID)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update virus scan status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrMediaNotFound
	}

	return nil
}

// UpdateOCRText updates the OCR extracted text.
func (r *MediaRepository) UpdateOCRText(ctx context.Context, tenantID, mediaID, ocrText string) error {
	query := `
		UPDATE media_schema.media_files
		SET ocr_text = $1, updated_at = NOW()
		WHERE media_id = $2`
	args := []any{ocrText, mediaID}
	if strings.TrimSpace(tenantID) != "" {
		query += " AND tenant_id = $3"
		args = append(args, tenantID)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update OCR text: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrMediaNotFound
	}

	return nil
}

// UpdateProcessedFiles updates optimized and thumbnail file references.
func (r *MediaRepository) UpdateProcessedFiles(ctx context.Context, tenantID, mediaID, optimizedFileID, thumbnailFileID string) error {
	query := `
		UPDATE media_schema.media_files
		SET optimized_file_id = $1, thumbnail_file_id = $2, updated_at = NOW()
		WHERE media_id = $3`
	args := []any{nullString(optimizedFileID), nullString(thumbnailFileID), mediaID}
	if strings.TrimSpace(tenantID) != "" {
		query += " AND tenant_id = $4"
		args = append(args, tenantID)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update processed files: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrMediaNotFound
	}

	return nil
}

// Delete hard deletes a media file.
func (r *MediaRepository) Delete(ctx context.Context, tenantID, mediaID string) error {
	query := `DELETE FROM media_schema.media_files WHERE media_id = $1`
	args := []any{mediaID}
	if strings.TrimSpace(tenantID) != "" {
		query += " AND tenant_id = $2"
		args = append(args, tenantID)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete media file: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrMediaNotFound
	}

	return nil
}

// Helper functions.
func nullInt32(val int32) sql.NullInt32 {
	if val == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: val, Valid: true}
}

func nullString(val string) sql.NullString {
	if val == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: val, Valid: true}
}

func parseMediaType(s string) mediav1.MediaType {
	if val, ok := mediav1.MediaType_value[s]; ok {
		return mediav1.MediaType(val)
	}
	return mediav1.MediaType_MEDIA_TYPE_UNSPECIFIED
}

func parseValidationStatus(s string) mediav1.ValidationStatus {
	if val, ok := mediav1.ValidationStatus_value[s]; ok {
		return mediav1.ValidationStatus(val)
	}
	return mediav1.ValidationStatus_VALIDATION_STATUS_UNSPECIFIED
}

func parseVirusScanStatus(s string) mediav1.VirusScanStatus {
	if val, ok := mediav1.VirusScanStatus_value[s]; ok {
		return mediav1.VirusScanStatus(val)
	}
	return mediav1.VirusScanStatus_VIRUS_SCAN_STATUS_UNSPECIFIED
}

// DownloadFile retrieves the raw file data for a media file (for worker processing)
// This is a stub implementation that returns error since actual file data is in storage service
func (r *MediaRepository) DownloadFile(ctx context.Context, mediaID string) ([]byte, string, error) {
	// TODO: Implement integration with storage service to download actual file data
	return nil, "", fmt.Errorf("file download not yet implemented - requires storage service integration")
}

// UpdateProcessingResult updates processing results after a job completes (OCR text and/or dimensions)
func (r *MediaRepository) UpdateProcessingResult(ctx context.Context, mediaID string, ocrText string, width, height int) error {
	var updates []string
	var args []any

	if ocrText != "" {
		updates = append(updates, fmt.Sprintf("ocr_text = $%d", len(args)+1))
		args = append(args, ocrText)
	}
	if width > 0 {
		updates = append(updates, fmt.Sprintf("width = $%d", len(args)+1))
		args = append(args, width)
	}
	if height > 0 {
		updates = append(updates, fmt.Sprintf("height = $%d", len(args)+1))
		args = append(args, height)
	}

	if len(updates) == 0 {
		return nil
	}

	updates = append(updates, fmt.Sprintf("updated_at = NOW()"))

	query := fmt.Sprintf(`
		UPDATE media_schema.media_files
		SET %s
		WHERE media_id = $%d
	`, strings.Join(updates, ", "), len(args)+1)

	args = append(args, mediaID)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update processing result: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrMediaNotFound
	}

	return nil
}

// UpdateVirusScanResult updates virus scan results after scanning completes
func (r *MediaRepository) UpdateVirusScanResult(ctx context.Context, mediaID string, clean bool, virusName string) error {
	status := mediav1.VirusScanStatus_VIRUS_SCAN_STATUS_CLEAN
	if !clean {
		status = mediav1.VirusScanStatus_VIRUS_SCAN_STATUS_INFECTED
	}

	query := `
		UPDATE media_schema.media_files
		SET virus_scan_status = $1, updated_at = NOW()
		WHERE media_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status.String(), mediaID)
	if err != nil {
		return fmt.Errorf("failed to update virus scan result: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrMediaNotFound
	}

	return nil
}
