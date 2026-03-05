package repository

import (
	"context"
	"errors"
	"time"

	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// UserDocumentRepository provides access to authn_schema.users_documents.
type UserDocumentRepository struct{ db *gorm.DB }

func NewUserDocumentRepository(db *gorm.DB) *UserDocumentRepository {
	return &UserDocumentRepository{db: db}
}

func (r *UserDocumentRepository) Create(ctx context.Context, d *authnv1.UserDocument) error {
	if d.CreatedAt == nil {
		d.CreatedAt = timestamppb.Now()
	}
	d.UpdatedAt = timestamppb.Now()

	// IMPORTANT: some columns are UUID in Postgres (e.g. policy_id, verified_by).
	// Proto fields are strings; empty string would be written as '' which Postgres rejects for uuid.
	// Insert via explicit column map to ensure optional UUID columns become NULL (not '').
	values := map[string]any{
		"user_document_id":    d.UserDocumentId,
		"user_id":             d.UserId,
		"document_type_id":    d.DocumentTypeId,
		"file_url":            d.FileUrl,
		"verification_status": d.VerificationStatus,
		"created_at":          d.CreatedAt.AsTime(),
		"updated_at":          d.UpdatedAt.AsTime(),
	}
	// Only include optional UUID/timestamp columns when they have values.
	if d.PolicyId != "" {
		values["policy_id"] = d.PolicyId
	}
	if d.VerifiedBy != "" {
		values["verified_by"] = d.VerifiedBy
	}
	if d.VerifiedAt != nil {
		values["verified_at"] = d.VerifiedAt.AsTime()
	}

	return r.db.WithContext(ctx).Table("authn_schema.users_documents").Create(values).Error
}

func (r *UserDocumentRepository) GetByID(ctx context.Context, id string) (*authnv1.UserDocument, error) {
	var d authnv1.UserDocument
	err := r.db.WithContext(ctx).
		Table("authn_schema.users_documents").
		Where("user_document_id = ?", id).
		First(&d).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *UserDocumentRepository) ListByUser(ctx context.Context, userID string) ([]*authnv1.UserDocument, error) {
	var docs []*authnv1.UserDocument
	err := r.db.WithContext(ctx).
		Table("authn_schema.users_documents").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&docs).Error
	return docs, err
}

func (r *UserDocumentRepository) UpdateVerification(ctx context.Context, id, status string, verifiedBy *string, verifiedAt *time.Time) error {
	upd := map[string]any{"verification_status": status, "updated_at": time.Now()}
	if verifiedBy != nil {
		upd["verified_by"] = *verifiedBy
	}
	if verifiedAt != nil {
		upd["verified_at"] = *verifiedAt
	}
	return r.db.WithContext(ctx).
		Table("authn_schema.users_documents").
		Where("user_document_id = ?", id).
		Updates(upd).Error
}

func (r *UserDocumentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("authn_schema.users_documents").Where("user_document_id = ?", id).Delete(map[string]any{}).Error
}

// Update applies partial updates to a user document.
func (r *UserDocumentRepository) Update(ctx context.Context, id string, updates map[string]any) error {
	if len(updates) == 0 {
		return errors.New("no updates provided")
	}
	updates["updated_at"] = time.Now()
	tx := r.db.WithContext(ctx).Table("authn_schema.users_documents").Where("user_document_id = ?", id).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// MarkVerified updates the verification status of a user document.
func (r *UserDocumentRepository) MarkVerified(ctx context.Context, docID, verifiedBy, status, rejectionReason string) error {
	upd := map[string]any{
		"verification_status": status,
		"verified_by":         verifiedBy,
		"verified_at":         "NOW()",
	}
	if rejectionReason != "" {
		upd["rejection_reason"] = rejectionReason
	}
	return r.db.WithContext(ctx).Table("authn_schema.users_documents").Where("user_document_id = ?", docID).Updates(upd).Error
}
