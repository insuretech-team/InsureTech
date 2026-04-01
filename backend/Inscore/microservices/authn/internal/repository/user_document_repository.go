package repository

import (
	"context"
	"database/sql"
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
	// BUG FIX: Cannot use GORM First() with proto struct — *timestamppb.Timestamp fields cause errors.
	type row struct {
		UserDocumentId     string         `gorm:"column:user_document_id"`
		UserId             string         `gorm:"column:user_id"`
		DocumentTypeId     string         `gorm:"column:document_type_id"`
		FileUrl            string         `gorm:"column:file_url"`
		VerificationStatus string         `gorm:"column:verification_status"`
		PolicyId           sql.NullString `gorm:"column:policy_id"`
		VerifiedBy         sql.NullString `gorm:"column:verified_by"`
		RejectionReason    sql.NullString `gorm:"column:rejection_reason"`
		VerifiedAt         sql.NullTime   `gorm:"column:verified_at"`
		CreatedAt          sql.NullTime   `gorm:"column:created_at"`
		UpdatedAt          sql.NullTime   `gorm:"column:updated_at"`
	}
	var r2 row
	err := r.db.WithContext(ctx).
		Table("authn_schema.users_documents").
		Select("user_document_id, user_id, document_type_id, file_url, verification_status, policy_id, verified_by, rejection_reason, verified_at, created_at, updated_at").
		Where("user_document_id = ?", id).
		Limit(1).
		Scan(&r2).Error
	if err != nil {
		return nil, err
	}
	if r2.UserDocumentId == "" {
		return nil, gorm.ErrRecordNotFound
	}
	d := &authnv1.UserDocument{
		UserDocumentId:     r2.UserDocumentId,
		UserId:             r2.UserId,
		DocumentTypeId:     r2.DocumentTypeId,
		FileUrl:            r2.FileUrl,
		VerificationStatus: r2.VerificationStatus,
	}
	if r2.PolicyId.Valid {
		d.PolicyId = r2.PolicyId.String
	}
	if r2.VerifiedBy.Valid {
		d.VerifiedBy = r2.VerifiedBy.String
	}
	if r2.RejectionReason.Valid {
		d.RejectionReason = r2.RejectionReason.String
	}
	if r2.VerifiedAt.Valid {
		d.VerifiedAt = timestamppb.New(r2.VerifiedAt.Time)
	}
	if r2.CreatedAt.Valid {
		d.CreatedAt = timestamppb.New(r2.CreatedAt.Time)
	}
	if r2.UpdatedAt.Valid {
		d.UpdatedAt = timestamppb.New(r2.UpdatedAt.Time)
	}
	return d, nil
}

func (r *UserDocumentRepository) ListByUser(ctx context.Context, userID string) ([]*authnv1.UserDocument, error) {
	// BUG FIX: Cannot use GORM Find() with proto struct — *timestamppb.Timestamp fields
	// (VerifiedAt, CreatedAt, UpdatedAt) cause "invalid field" GORM errors.
	// Use raw SQL with intermediate scan variables instead.
	type row struct {
		UserDocumentId     string         `gorm:"column:user_document_id"`
		UserId             string         `gorm:"column:user_id"`
		DocumentTypeId     string         `gorm:"column:document_type_id"`
		FileUrl            string         `gorm:"column:file_url"`
		VerificationStatus string         `gorm:"column:verification_status"`
		PolicyId           sql.NullString `gorm:"column:policy_id"`
		VerifiedBy         sql.NullString `gorm:"column:verified_by"`
		RejectionReason    sql.NullString `gorm:"column:rejection_reason"`
		VerifiedAt         sql.NullTime   `gorm:"column:verified_at"`
		CreatedAt          sql.NullTime   `gorm:"column:created_at"`
		UpdatedAt          sql.NullTime   `gorm:"column:updated_at"`
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Table("authn_schema.users_documents").
		Select("user_document_id, user_id, document_type_id, file_url, verification_status, policy_id, verified_by, rejection_reason, verified_at, created_at, updated_at").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	docs := make([]*authnv1.UserDocument, 0, len(rows))
	for _, r := range rows {
		d := &authnv1.UserDocument{
			UserDocumentId:     r.UserDocumentId,
			UserId:             r.UserId,
			DocumentTypeId:     r.DocumentTypeId,
			FileUrl:            r.FileUrl,
			VerificationStatus: r.VerificationStatus,
		}
		if r.PolicyId.Valid {
			d.PolicyId = r.PolicyId.String
		}
		if r.VerifiedBy.Valid {
			d.VerifiedBy = r.VerifiedBy.String
		}
		if r.RejectionReason.Valid {
			d.RejectionReason = r.RejectionReason.String
		}
		if r.VerifiedAt.Valid {
			d.VerifiedAt = timestamppb.New(r.VerifiedAt.Time)
		}
		if r.CreatedAt.Valid {
			d.CreatedAt = timestamppb.New(r.CreatedAt.Time)
		}
		if r.UpdatedAt.Valid {
			d.UpdatedAt = timestamppb.New(r.UpdatedAt.Time)
		}
		docs = append(docs, d)
	}
	return docs, nil
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
