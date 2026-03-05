package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	documentv1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrTemplateNotFound = errors.New("document template not found")
)

// DocumentTemplateRepository handles storage_schema.document_templates.
type DocumentTemplateRepository struct {
	db *sqlx.DB
}

func NewDocumentTemplateRepository(db *sqlx.DB) *DocumentTemplateRepository {
	return &DocumentTemplateRepository{db: db}
}

func (r *DocumentTemplateRepository) Create(ctx context.Context, tpl *documentv1.DocumentTemplate) (*documentv1.DocumentTemplate, error) {
	query := `
		INSERT INTO storage_schema.document_templates (
			template_id, name, type, description, template_content, output_format,
			variables, version, is_active, audit_info, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query,
		tpl.Id,
		tpl.Name,
		tpl.Type.String(),
		nullString(tpl.Description),
		tpl.TemplateContent,
		tpl.OutputFormat.String(),
		nullString(tpl.Variables),
		tpl.Version,
		tpl.IsActive,
		"{}",
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	if tpl.AuditInfo == nil {
		tpl.AuditInfo = &commonv1.AuditInfo{}
	}
	if createdAt.Valid {
		tpl.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		tpl.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	return tpl, nil
}

func (r *DocumentTemplateRepository) UpsertByName(ctx context.Context, tpl *documentv1.DocumentTemplate) (*documentv1.DocumentTemplate, error) {
	query := `
		INSERT INTO storage_schema.document_templates (
			template_id, name, type, description, template_content, output_format,
			variables, version, is_active, audit_info, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
		ON CONFLICT (name)
		DO UPDATE SET
			type = EXCLUDED.type,
			description = EXCLUDED.description,
			template_content = EXCLUDED.template_content,
			output_format = EXCLUDED.output_format,
			variables = EXCLUDED.variables,
			is_active = EXCLUDED.is_active,
			version = storage_schema.document_templates.version + 1,
			updated_at = NOW()
		RETURNING template_id, version, created_at, updated_at
	`

	var createdAt, updatedAt sql.NullTime
	var id string
	var version int32
	err := r.db.QueryRowContext(ctx, query,
		tpl.Id,
		tpl.Name,
		tpl.Type.String(),
		nullString(tpl.Description),
		tpl.TemplateContent,
		tpl.OutputFormat.String(),
		nullString(tpl.Variables),
		maxInt32(1, tpl.Version),
		tpl.IsActive,
		"{}",
	).Scan(&id, &version, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert template: %w", err)
	}

	tpl.Id = id
	tpl.Version = version
	if tpl.AuditInfo == nil {
		tpl.AuditInfo = &commonv1.AuditInfo{}
	}
	if createdAt.Valid {
		tpl.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		tpl.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	return tpl, nil
}

func (r *DocumentTemplateRepository) GetByID(ctx context.Context, templateID string) (*documentv1.DocumentTemplate, error) {
	query := `
		SELECT template_id, name, type, description, template_content, output_format,
			variables, version, is_active, created_at, updated_at
		FROM storage_schema.document_templates
		WHERE template_id = $1
	`
	return r.scanOne(ctx, query, templateID)
}

func (r *DocumentTemplateRepository) GetByName(ctx context.Context, name string) (*documentv1.DocumentTemplate, error) {
	query := `
		SELECT template_id, name, type, description, template_content, output_format,
			variables, version, is_active, created_at, updated_at
		FROM storage_schema.document_templates
		WHERE name = $1
	`
	return r.scanOne(ctx, query, name)
}

func (r *DocumentTemplateRepository) List(ctx context.Context, docType *documentv1.DocumentType, activeOnly bool, limit, offset int) ([]*documentv1.DocumentTemplate, int, error) {
	conditions := []string{"1=1"}
	args := make([]any, 0, 3)

	if docType != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, docType.String())
	}
	if activeOnly {
		conditions = append(conditions, "is_active = true")
	}

	where := strings.Join(conditions, " AND ")
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM storage_schema.document_templates WHERE %s`, where)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count templates: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT template_id, name, type, description, template_content, output_format,
			variables, version, is_active, created_at, updated_at
		FROM storage_schema.document_templates
		WHERE %s
		ORDER BY updated_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)
	rows, err := r.db.QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list templates: %w", err)
	}
	defer rows.Close()

	templates := make([]*documentv1.DocumentTemplate, 0)
	for rows.Next() {
		tpl, scanErr := scanTemplate(rows.Scan)
		if scanErr != nil {
			return nil, 0, fmt.Errorf("failed to scan template: %w", scanErr)
		}
		templates = append(templates, tpl)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate templates: %w", err)
	}
	return templates, total, nil
}

func (r *DocumentTemplateRepository) Update(ctx context.Context, templateID string, tpl *documentv1.DocumentTemplate) error {
	query := `
		UPDATE storage_schema.document_templates
		SET name = $1,
			type = $2,
			description = $3,
			template_content = $4,
			output_format = $5,
			variables = $6,
			is_active = $7,
			version = $8,
			updated_at = NOW()
		WHERE template_id = $9
	`
	res, err := r.db.ExecContext(ctx, query,
		tpl.Name,
		tpl.Type.String(),
		nullString(tpl.Description),
		tpl.TemplateContent,
		tpl.OutputFormat.String(),
		nullString(tpl.Variables),
		tpl.IsActive,
		maxInt32(1, tpl.Version),
		templateID,
	)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if ra == 0 {
		return ErrTemplateNotFound
	}
	return nil
}

func (r *DocumentTemplateRepository) Deactivate(ctx context.Context, templateID string) error {
	query := `UPDATE storage_schema.document_templates SET is_active = false, updated_at = NOW() WHERE template_id = $1`
	res, err := r.db.ExecContext(ctx, query, templateID)
	if err != nil {
		return fmt.Errorf("failed to deactivate template: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if ra == 0 {
		return ErrTemplateNotFound
	}
	return nil
}

func (r *DocumentTemplateRepository) Delete(ctx context.Context, templateID string) error {
	query := `DELETE FROM storage_schema.document_templates WHERE template_id = $1`
	res, err := r.db.ExecContext(ctx, query, templateID)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if ra == 0 {
		return ErrTemplateNotFound
	}
	return nil
}

func (r *DocumentTemplateRepository) scanOne(ctx context.Context, query string, arg any) (*documentv1.DocumentTemplate, error) {
	tpl, err := scanTemplate(func(dest ...any) error {
		return r.db.QueryRowContext(ctx, query, arg).Scan(dest...)
	})
	if err == sql.ErrNoRows {
		return nil, ErrTemplateNotFound
	}
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func scanTemplate(scan func(dest ...any) error) (*documentv1.DocumentTemplate, error) {
	var tpl documentv1.DocumentTemplate
	var typeDB, formatDB string
	var description, variables sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := scan(
		&tpl.Id,
		&tpl.Name,
		&typeDB,
		&description,
		&tpl.TemplateContent,
		&formatDB,
		&variables,
		&tpl.Version,
		&tpl.IsActive,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	tpl.Type = parseDocumentType(typeDB)
	tpl.OutputFormat = parseOutputFormat(formatDB)
	if description.Valid {
		tpl.Description = description.String
	}
	if variables.Valid {
		tpl.Variables = variables.String
	}
	tpl.AuditInfo = &commonv1.AuditInfo{}
	if createdAt.Valid {
		tpl.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		tpl.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &tpl, nil
}

func parseDocumentType(v string) documentv1.DocumentType {
	s := strings.TrimSpace(v)
	if s == "" {
		return documentv1.DocumentType_DOCUMENT_TYPE_UNSPECIFIED
	}
	if i, err := strconv.Atoi(s); err == nil {
		t := documentv1.DocumentType(i)
		if _, ok := documentv1.DocumentType_name[int32(t)]; ok {
			return t
		}
	}
	if i, ok := documentv1.DocumentType_value[s]; ok {
		return documentv1.DocumentType(i)
	}
	if i, ok := documentv1.DocumentType_value["DOCUMENT_TYPE_"+s]; ok {
		return documentv1.DocumentType(i)
	}
	return documentv1.DocumentType_DOCUMENT_TYPE_UNSPECIFIED
}

func parseOutputFormat(v string) documentv1.OutputFormat {
	s := strings.TrimSpace(v)
	if s == "" {
		return documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED
	}
	if i, err := strconv.Atoi(s); err == nil {
		o := documentv1.OutputFormat(i)
		if _, ok := documentv1.OutputFormat_name[int32(o)]; ok {
			return o
		}
	}
	if i, ok := documentv1.OutputFormat_value[s]; ok {
		return documentv1.OutputFormat(i)
	}
	if i, ok := documentv1.OutputFormat_value["OUTPUT_FORMAT_"+s]; ok {
		return documentv1.OutputFormat(i)
	}
	return documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED
}

func parseTemplateVariables(raw string) []string {
	out := make([]string, 0)
	if strings.TrimSpace(raw) == "" {
		return out
	}
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}

func nullString(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

func maxInt32(min, v int32) int32 {
	if v < min {
		return min
	}
	return v
}
