package repository

import (
	"context"
	"strings"

	kycv1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	"gorm.io/gorm"
)

// DocumentVerificationRepository provides access to authn_schema.document_verifications.
type DocumentVerificationRepository struct{ db *gorm.DB }

func NewDocumentVerificationRepository(db *gorm.DB) *DocumentVerificationRepository {
	return &DocumentVerificationRepository{db: db}
}

func (r *DocumentVerificationRepository) Create(ctx context.Context, d *kycv1.DocumentVerification) error {
	return r.db.WithContext(ctx).Exec(
		`insert into authn_schema.document_verifications
			(doc_verification_id, kyc_verification_id, document_type, document_number, document_url, extracted_data, status, confidence_score, audit_info)
		 values (?, ?, ?, ?, ?, ?::jsonb, ?, ?, '{}'::jsonb)`,
		d.Id,
		d.KycVerificationId,
		strings.TrimPrefix(d.DocumentType.String(), "DOCUMENT_TYPE_"),
		d.DocumentNumber,
		nullableString(d.DocumentUrl),
		nullableJSON(d.ExtractedData),
		strings.TrimPrefix(d.Status.String(), "DOCUMENT_STATUS_"),
		d.ConfidenceScore,
	).Error
}

const docVerificationCols = `doc_verification_id, kyc_verification_id, document_type, document_number, document_url, extracted_data, status, confidence_score`

func scanDocVerification(row interface{ Scan(...any) error }) (*kycv1.DocumentVerification, error) {
	var d kycv1.DocumentVerification
	var docTypeStr, statusStr string
	var extractedData []byte
	var docURL *string
	if err := row.Scan(&d.Id, &d.KycVerificationId, &docTypeStr, &d.DocumentNumber, &docURL, &extractedData, &statusStr, &d.ConfidenceScore); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if docURL != nil {
		d.DocumentUrl = *docURL
	}
	d.DocumentType = docTypeFromString(docTypeStr)
	d.Status = docStatusFromString(statusStr)
	return &d, nil
}

func (r *DocumentVerificationRepository) GetByID(ctx context.Context, id string) (*kycv1.DocumentVerification, error) {
	row := r.db.WithContext(ctx).Raw(`select `+docVerificationCols+` from authn_schema.document_verifications where doc_verification_id = ? limit 1`, id).Row()
	if err := row.Err(); err != nil {
		return nil, err
	}
	d, err := scanDocVerification(row)
	if err != nil {
		return nil, err
	}
	if d.Id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return d, nil
}

func (r *DocumentVerificationRepository) ListByKYC(ctx context.Context, kycID string, limit, offset int) ([]*kycv1.DocumentVerification, error) {
	q := `select ` + docVerificationCols + ` from authn_schema.document_verifications where kyc_verification_id = ? order by doc_verification_id desc`
	args := []any{kycID}
	if limit > 0 {
		q += " limit ?"
		args = append(args, limit)
	}
	if offset > 0 {
		q += " offset ?"
		args = append(args, offset)
	}
	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*kycv1.DocumentVerification
	for rows.Next() {
		d, err := scanDocVerification(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *DocumentVerificationRepository) UpdateStatus(ctx context.Context, id string, st kycv1.DocumentStatus, confidence *float32) error {
	upd := map[string]any{"status": docStatusToString(st)}
	if confidence != nil {
		upd["confidence_score"] = *confidence
	}
	return r.db.WithContext(ctx).Table("authn_schema.document_verifications").Where("doc_verification_id = ?", id).Updates(upd).Error
}

func (r *DocumentVerificationRepository) DeleteByKYC(ctx context.Context, kycID string) (int64, error) {
	res := r.db.WithContext(ctx).Table("authn_schema.document_verifications").Where("kyc_verification_id = ?", kycID).Delete(map[string]any{})
	return res.RowsAffected, res.Error
}

func docTypeFromString(s string) kycv1.DocumentType {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "DOCUMENT_TYPE_")
	if v, ok := kycv1.DocumentType_value["DOCUMENT_TYPE_"+s]; ok {
		return kycv1.DocumentType(v)
	}
	return kycv1.DocumentType_DOCUMENT_TYPE_UNSPECIFIED
}

func docStatusToString(s kycv1.DocumentStatus) string {
	return strings.TrimPrefix(s.String(), "DOCUMENT_STATUS_")
}

func docStatusFromString(s string) kycv1.DocumentStatus {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "DOCUMENT_STATUS_")
	if v, ok := kycv1.DocumentStatus_value["DOCUMENT_STATUS_"+s]; ok {
		return kycv1.DocumentStatus(v)
	}
	return kycv1.DocumentStatus_DOCUMENT_STATUS_UNSPECIFIED
}
