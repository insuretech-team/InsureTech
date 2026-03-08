package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	claimsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/claims/entity/v1"
)

type ClaimDocumentRepository struct {
	db *gorm.DB
}

func NewClaimDocumentRepository(db *gorm.DB) *ClaimDocumentRepository {
	return &ClaimDocumentRepository{db: db}
}

func (r *ClaimDocumentRepository) Create(ctx context.Context, doc *claimsv1.ClaimDocument) (*claimsv1.ClaimDocument, error) {
	if doc.DocumentId == "" {
		return nil, fmt.Errorf("document_id is required")
	}
	
	var uploadedAt interface{}
	if doc.UploadedAt != nil {
		uploadedAt = doc.UploadedAt.AsTime()
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.claim_documents
			(document_id, claim_id, document_type, file_url, file_hash,
			 uploaded_at, verified, verified_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`,
		doc.DocumentId, doc.ClaimId, doc.DocumentType, doc.FileUrl, doc.FileHash,
		uploadedAt, doc.Verified, doc.VerifiedBy,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert claim document: %w", err)
	}

	return r.GetByID(ctx, doc.DocumentId)
}

func (r *ClaimDocumentRepository) GetByID(ctx context.Context, documentID string) (*claimsv1.ClaimDocument, error) {
	var (
		d           claimsv1.ClaimDocument
		uploadedAt  time.Time
		verifiedBy  sql.NullString
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT document_id, claim_id, document_type, file_url, file_hash,
		       uploaded_at, verified, verified_by, created_at, updated_at
		FROM insurance_schema.claim_documents
		WHERE document_id = $1`,
		documentID,
	).Row().Scan(
		&d.DocumentId, &d.ClaimId, &d.DocumentType, &d.FileUrl, &d.FileHash,
		&uploadedAt, &d.Verified, &verifiedBy, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get claim document: %w", err)
	}

	if !uploadedAt.IsZero() {
		d.UploadedAt = timestamppb.New(uploadedAt)
	}
	if verifiedBy.Valid {
		d.VerifiedBy = verifiedBy.String
	}
	if !createdAt.IsZero() {
		d.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		d.UpdatedAt = timestamppb.New(updatedAt)
	}

	return &d, nil
}

func (r *ClaimDocumentRepository) Update(ctx context.Context, doc *claimsv1.ClaimDocument) (*claimsv1.ClaimDocument, error) {
	var uploadedAt interface{}
	if doc.UploadedAt != nil {
		uploadedAt = doc.UploadedAt.AsTime()
	}
	
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.claim_documents
		SET claim_id = $2, document_type = $3, file_url = $4, file_hash = $5,
		    uploaded_at = $6, verified = $7, verified_by = $8, updated_at = NOW()
		WHERE document_id = $1`,
		doc.DocumentId, doc.ClaimId, doc.DocumentType, doc.FileUrl, doc.FileHash,
		uploadedAt, doc.Verified, doc.VerifiedBy,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update claim document: %w", err)
	}

	return r.GetByID(ctx, doc.DocumentId)
}

func (r *ClaimDocumentRepository) Delete(ctx context.Context, documentID string) error {
	return r.db.WithContext(ctx).Exec(`DELETE FROM insurance_schema.claim_documents WHERE document_id = $1`, documentID).Error
}

func (r *ClaimDocumentRepository) ListByClaimID(ctx context.Context, claimID string) ([]*claimsv1.ClaimDocument, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT document_id, claim_id, document_type, file_url, file_hash,
		       uploaded_at, verified, verified_by, created_at, updated_at
		FROM insurance_schema.claim_documents
		WHERE claim_id = $1
		ORDER BY uploaded_at DESC`,
		claimID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list claim documents: %w", err)
	}
	defer rows.Close()

	documents := make([]*claimsv1.ClaimDocument, 0)
	for rows.Next() {
		var (
			d           claimsv1.ClaimDocument
			uploadedAt  time.Time
			verifiedBy  sql.NullString
			createdAt   time.Time
			updatedAt   time.Time
		)

		err := rows.Scan(
			&d.DocumentId, &d.ClaimId, &d.DocumentType, &d.FileUrl, &d.FileHash,
			&uploadedAt, &d.Verified, &verifiedBy, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan claim document: %w", err)
		}

		if !uploadedAt.IsZero() {
			d.UploadedAt = timestamppb.New(uploadedAt)
		}
		if verifiedBy.Valid {
			d.VerifiedBy = verifiedBy.String
		}
		if !createdAt.IsZero() {
			d.CreatedAt = timestamppb.New(createdAt)
		}
		if !updatedAt.IsZero() {
			d.UpdatedAt = timestamppb.New(updatedAt)
		}

		documents = append(documents, &d)
	}

	return documents, nil
}
