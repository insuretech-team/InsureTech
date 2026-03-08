package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	documentv1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
)

// DocumentGenerationRepository handles storage_schema.document_generations.
type DocumentGenerationRepository struct {
	db *sqlx.DB
}

func NewDocumentGenerationRepository(db *sqlx.DB) *DocumentGenerationRepository {
	return &DocumentGenerationRepository{db: db}
}

func (r *DocumentGenerationRepository) Create(ctx context.Context, doc *documentv1.DocumentGeneration) (*documentv1.DocumentGeneration, error) {
	query := `
		INSERT INTO storage_schema.document_generations (
			generation_id, document_template_id, entity_type, entity_id, data,
			status, file_url, file_size_bytes, qr_code_data, generated_by,
			generated_at, audit_info, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
		RETURNING created_at, updated_at
	`

	generatedAt := sql.NullTime{Valid: false}
	if doc.GeneratedAt != nil {
		generatedAt = sql.NullTime{Time: doc.GeneratedAt.AsTime(), Valid: true}
	} else {
		now := time.Now().UTC()
		generatedAt = sql.NullTime{Time: now, Valid: true}
		doc.GeneratedAt = timestamppb.New(now)
	}

	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query,
		doc.Id,
		doc.DocumentTemplateId,
		doc.EntityType,
		doc.EntityId,
		doc.Data,
		generationStatusDBValue(doc.Status),
		nullString(doc.FileUrl),
		doc.FileSizeBytes,
		nullString(doc.QrCodeData),
		nullString(doc.GeneratedBy),
		generatedAt,
		"{}",
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		if isUndefinedColumnErr(err) {
			legacyQuery := `
				INSERT INTO storage_schema.document_generations (
					generation_id, document_template_id, entity_type, entity_id, data,
					status, file_url, file_size_bytes, qr_code_data, generated_by,
					generated_at, audit_info
				) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
			`
			_, legacyErr := r.db.ExecContext(ctx, legacyQuery,
				doc.Id,
				doc.DocumentTemplateId,
				doc.EntityType,
				doc.EntityId,
				doc.Data,
				generationStatusDBValue(doc.Status),
				nullString(doc.FileUrl),
				doc.FileSizeBytes,
				nullString(doc.QrCodeData),
				nullString(doc.GeneratedBy),
				generatedAt,
				"{}",
			)
			if legacyErr != nil {
				return nil, fmt.Errorf("failed to create document generation: %w", legacyErr)
			}
		} else {
			return nil, fmt.Errorf("failed to create document generation: %w", err)
		}
	}

	if doc.AuditInfo == nil {
		doc.AuditInfo = &commonv1.AuditInfo{}
	}
	if createdAt.Valid {
		doc.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		doc.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	return doc, nil
}

func (r *DocumentGenerationRepository) GetByID(ctx context.Context, documentID string) (*documentv1.DocumentGeneration, error) {
	query := `
		SELECT generation_id, document_template_id, entity_type, entity_id, data,
			status, file_url, file_size_bytes, qr_code_data, generated_by,
			generated_at, created_at, updated_at
		FROM storage_schema.document_generations
		WHERE generation_id = $1
	`
	doc, err := scanDocumentGeneration(func(dest ...any) error {
		return r.db.QueryRowContext(ctx, query, documentID).Scan(dest...)
	})
	if err != nil && isUndefinedColumnErr(err) {
		legacyQuery := `
			SELECT generation_id, document_template_id, entity_type, entity_id, data,
				status, file_url, file_size_bytes, qr_code_data, generated_by,
				generated_at
			FROM storage_schema.document_generations
			WHERE generation_id = $1
		`
		doc, err = scanDocumentGenerationLegacy(func(dest ...any) error {
			return r.db.QueryRowContext(ctx, legacyQuery, documentID).Scan(dest...)
		})
	}
	if err == sql.ErrNoRows {
		return nil, ErrDocumentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return doc, nil
}

func (r *DocumentGenerationRepository) ListByEntity(
	ctx context.Context,
	entityType, entityID string,
	status *documentv1.GenerationStatus,
	limit, offset int,
) ([]*documentv1.DocumentGeneration, int, error) {
	conditions := []string{"entity_type = $1", "entity_id = $2"}
	args := []any{entityType, entityID}
	if status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, generationStatusDBValue(*status))
	}
	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM storage_schema.document_generations WHERE %s`, where)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT generation_id, document_template_id, entity_type, entity_id, data,
			status, file_url, file_size_bytes, qr_code_data, generated_by,
			generated_at, created_at, updated_at
		FROM storage_schema.document_generations
		WHERE %s
		ORDER BY generated_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, append(args, limit, offset)...)
	useLegacyScan := false
	if err != nil {
		if isUndefinedColumnErr(err) {
			legacyQuery := fmt.Sprintf(`
				SELECT generation_id, document_template_id, entity_type, entity_id, data,
					status, file_url, file_size_bytes, qr_code_data, generated_by,
					generated_at
				FROM storage_schema.document_generations
				WHERE %s
				ORDER BY generated_at DESC
				LIMIT $%d OFFSET $%d
			`, where, len(args)+1, len(args)+2)
			rows, err = r.db.QueryContext(ctx, legacyQuery, append(args, limit, offset)...)
			useLegacyScan = true
		}
		if err != nil {
			return nil, 0, fmt.Errorf("failed to list documents: %w", err)
		}
	}
	defer rows.Close()

	docs := make([]*documentv1.DocumentGeneration, 0)
	for rows.Next() {
		var (
			doc     *documentv1.DocumentGeneration
			scanErr error
		)
		if useLegacyScan {
			doc, scanErr = scanDocumentGenerationLegacy(rows.Scan)
		} else {
			doc, scanErr = scanDocumentGeneration(rows.Scan)
		}
		if scanErr != nil {
			return nil, 0, fmt.Errorf("failed to scan document: %w", scanErr)
		}
		docs = append(docs, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate documents: %w", err)
	}
	return docs, total, nil
}

func (r *DocumentGenerationRepository) Delete(ctx context.Context, documentID string) error {
	query := `DELETE FROM storage_schema.document_generations WHERE generation_id = $1`
	res, err := r.db.ExecContext(ctx, query, documentID)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if ra == 0 {
		return ErrDocumentNotFound
	}
	return nil
}

func scanDocumentGeneration(scan func(dest ...any) error) (*documentv1.DocumentGeneration, error) {
	var doc documentv1.DocumentGeneration
	var statusDB string
	var fileURL, qrCodeData, generatedBy sql.NullString
	var generatedAt, createdAt, updatedAt sql.NullTime

	err := scan(
		&doc.Id,
		&doc.DocumentTemplateId,
		&doc.EntityType,
		&doc.EntityId,
		&doc.Data,
		&statusDB,
		&fileURL,
		&doc.FileSizeBytes,
		&qrCodeData,
		&generatedBy,
		&generatedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	doc.Status = parseGenerationStatus(statusDB)
	if fileURL.Valid {
		doc.FileUrl = fileURL.String
	}
	if qrCodeData.Valid {
		doc.QrCodeData = qrCodeData.String
	}
	if generatedBy.Valid {
		doc.GeneratedBy = generatedBy.String
	}
	if generatedAt.Valid {
		doc.GeneratedAt = timestamppb.New(generatedAt.Time)
	}
	doc.AuditInfo = &commonv1.AuditInfo{}
	if createdAt.Valid {
		doc.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		doc.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &doc, nil
}

func scanDocumentGenerationLegacy(scan func(dest ...any) error) (*documentv1.DocumentGeneration, error) {
	var doc documentv1.DocumentGeneration
	var statusDB string
	var fileURL, qrCodeData, generatedBy sql.NullString
	var generatedAt sql.NullTime

	err := scan(
		&doc.Id,
		&doc.DocumentTemplateId,
		&doc.EntityType,
		&doc.EntityId,
		&doc.Data,
		&statusDB,
		&fileURL,
		&doc.FileSizeBytes,
		&qrCodeData,
		&generatedBy,
		&generatedAt,
	)
	if err != nil {
		return nil, err
	}

	doc.Status = parseGenerationStatus(statusDB)
	if fileURL.Valid {
		doc.FileUrl = fileURL.String
	}
	if qrCodeData.Valid {
		doc.QrCodeData = qrCodeData.String
	}
	if generatedBy.Valid {
		doc.GeneratedBy = generatedBy.String
	}
	if generatedAt.Valid {
		doc.GeneratedAt = timestamppb.New(generatedAt.Time)
	}
	doc.AuditInfo = &commonv1.AuditInfo{}
	return &doc, nil
}

func parseGenerationStatus(v string) documentv1.GenerationStatus {
	s := strings.TrimSpace(v)
	if s == "" {
		return documentv1.GenerationStatus_GENERATION_STATUS_UNSPECIFIED
	}
	if i, err := strconv.Atoi(s); err == nil {
		st := documentv1.GenerationStatus(i)
		if _, ok := documentv1.GenerationStatus_name[int32(st)]; ok {
			return st
		}
	}
	if i, ok := documentv1.GenerationStatus_value[s]; ok {
		return documentv1.GenerationStatus(i)
	}
	if i, ok := documentv1.GenerationStatus_value["GENERATION_STATUS_"+s]; ok {
		return documentv1.GenerationStatus(i)
	}
	return documentv1.GenerationStatus_GENERATION_STATUS_UNSPECIFIED
}

func generationStatusDBValue(status documentv1.GenerationStatus) string {
	switch status {
	case documentv1.GenerationStatus_GENERATION_STATUS_PENDING:
		return "PENDING"
	case documentv1.GenerationStatus_GENERATION_STATUS_PROCESSING:
		return "PROCESSING"
	case documentv1.GenerationStatus_GENERATION_STATUS_COMPLETED:
		return "COMPLETED"
	case documentv1.GenerationStatus_GENERATION_STATUS_FAILED:
		return "FAILED"
	default:
		return "UNSPECIFIED"
	}
}
