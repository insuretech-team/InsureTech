package repository

import (
	"context"

	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"gorm.io/gorm"
)

// DocumentTypeRepository provides access to authn_schema.document_types.
// Uses proto-generated DocumentType struct directly with GORM.
type DocumentTypeRepository struct{ db *gorm.DB }

func NewDocumentTypeRepository(db *gorm.DB) *DocumentTypeRepository {
	return &DocumentTypeRepository{db: db}
}

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
func (r *DocumentTypeRepository) ListActive(ctx context.Context) ([]*authnv1.DocumentType, error) {
	var out []*authnv1.DocumentType
	if err := r.db.WithContext(ctx).
		Table("authn_schema.document_types").
		Where("is_active = true").
		Order("code asc").
		Scan(&out).Error; err != nil {
		return nil, err
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

func (r *DocumentTypeRepository) getOne(ctx context.Context, where string, args ...any) (*authnv1.DocumentType, error) {
	d := &authnv1.DocumentType{}
	if err := r.db.WithContext(ctx).
		Table("authn_schema.document_types").
		Where(where, args...).
		Limit(1).
		Scan(d).Error; err != nil {
		return nil, err
	}
	if d.DocumentTypeId == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return d, nil
}
