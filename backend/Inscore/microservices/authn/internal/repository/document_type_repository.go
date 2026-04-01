package repository

import (
	"context"
	"database/sql"

	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// DocumentTypeRepository provides access to authn_schema.document_types.
// Uses raw SQL scan rows (not GORM struct scan) to avoid *timestamppb.Timestamp field errors.
type DocumentTypeRepository struct{ db *gorm.DB }

func NewDocumentTypeRepository(db *gorm.DB) *DocumentTypeRepository {
	return &DocumentTypeRepository{db: db}
}

// docTypeRow is an intermediate struct for GORM scan that uses sql.NullTime
// instead of *timestamppb.Timestamp to avoid "invalid field" GORM errors.
type docTypeRow struct {
	DocumentTypeId string         `gorm:"column:document_type_id"`
	Code           string         `gorm:"column:code"`
	Name           string         `gorm:"column:name"`
	Description    sql.NullString `gorm:"column:description"`
	IsActive       bool           `gorm:"column:is_active"`
	CreatedAt      sql.NullTime   `gorm:"column:created_at"`
	UpdatedAt      sql.NullTime   `gorm:"column:updated_at"`
}

func (r docTypeRow) toProto() *authnv1.DocumentType {
	d := &authnv1.DocumentType{
		DocumentTypeId: r.DocumentTypeId,
		Code:           r.Code,
		Name:           r.Name,
		IsActive:       r.IsActive,
	}
	if r.Description.Valid {
		d.Description = r.Description.String
	}
	if r.CreatedAt.Valid {
		d.CreatedAt = timestamppb.New(r.CreatedAt.Time)
	}
	if r.UpdatedAt.Valid {
		d.UpdatedAt = timestamppb.New(r.UpdatedAt.Time)
	}
	return d
}

const docTypeSelectCols = "document_type_id, code, name, description, is_active, created_at, updated_at"

// Create inserts a new document type using raw SQL.
func (r *DocumentTypeRepository) Create(ctx context.Context, d *authnv1.DocumentType) error {
	return r.db.WithContext(ctx).Exec(
		`insert into authn_schema.document_types (document_type_id, code, name, description, is_active, created_at, updated_at)
		 values (?, ?, ?, ?, ?, now(), now())`,
		d.DocumentTypeId, d.Code, d.Name, nullableString(d.Description), d.IsActive,
	).Error
}

// GetByID returns a document type by primary key.
func (r *DocumentTypeRepository) GetByID(ctx context.Context, id string) (*authnv1.DocumentType, error) {
	return r.getOne(ctx, "document_type_id = ?", id)
}

// GetByCode returns a document type by code.
func (r *DocumentTypeRepository) GetByCode(ctx context.Context, code string) (*authnv1.DocumentType, error) {
	return r.getOne(ctx, "code = ?", code)
}

// ListActive returns all active document types ordered by code.
// BUG FIX: Use intermediate docTypeRow struct to avoid GORM "invalid field CreatedAt" error
// when scanning into proto struct with *timestamppb.Timestamp fields.
func (r *DocumentTypeRepository) ListActive(ctx context.Context) ([]*authnv1.DocumentType, error) {
	var rows []docTypeRow
	if err := r.db.WithContext(ctx).
		Table("authn_schema.document_types").
		Select(docTypeSelectCols).
		Where("is_active = true").
		Order("code asc").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*authnv1.DocumentType, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.toProto())
	}
	return out, nil
}

// SetActive updates the is_active flag for a document type.
func (r *DocumentTypeRepository) SetActive(ctx context.Context, id string, active bool) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.document_types").
		Where("document_type_id = ?", id).
		Update("is_active", active).Error
}

// Delete hard-deletes a document type.
func (r *DocumentTypeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.document_types").
		Where("document_type_id = ?", id).
		Delete(map[string]any{}).Error
}

// getOne is a helper that scans a single document type row using the intermediate struct.
func (r *DocumentTypeRepository) getOne(ctx context.Context, where string, args ...any) (*authnv1.DocumentType, error) {
	var row docTypeRow
	if err := r.db.WithContext(ctx).
		Table("authn_schema.document_types").
		Select(docTypeSelectCols).
		Where(where, args...).
		Limit(1).
		Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.DocumentTypeId == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return row.toProto(), nil
}
