package repository

import (
	"context"
	"strings"
	"time"

	kycv1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// KYCVerificationRepository provides access to authn_schema.kyc_verifications.
type KYCVerificationRepository struct{ db *gorm.DB }

func NewKYCVerificationRepository(db *gorm.DB) *KYCVerificationRepository {
	return &KYCVerificationRepository{db: db}
}

func (r *KYCVerificationRepository) Create(ctx context.Context, k *kycv1.KYCVerification) error {
	return r.db.WithContext(ctx).Exec(
		`insert into authn_schema.kyc_verifications
			(verification_id, type, entity_type, entity_id, method, provider, provider_reference, documents, status, verification_result, rejection_reason, verified_by, verified_at, expires_at, audit_info)
		 values (?, ?, ?, ?, ?, ?, ?, ?::jsonb, ?, ?::jsonb, ?, ?, ?, ?, '{}'::jsonb)`,
		k.Id,
		strings.TrimPrefix(k.Type.String(), "VERIFICATION_TYPE_"),
		k.EntityType,
		k.EntityId,
		strings.TrimPrefix(k.Method.String(), "VERIFICATION_METHOD_"),
		nullableString(k.Provider),
		nullableString(k.ProviderReference),
		nullableJSON(k.Documents),
		verificationStatusToString(k.Status),
		nullableJSON(k.VerificationResult),
		nullableString(k.RejectionReason),
		nullableString(k.VerifiedBy),
		nilOrTime(k.VerifiedAt),
		nilOrTime(k.ExpiresAt),
	).Error
}

func (r *KYCVerificationRepository) GetByID(ctx context.Context, id string) (*kycv1.KYCVerification, error) {
	return r.getOne(ctx, "verification_id = ?", id)
}

func (r *KYCVerificationRepository) GetByEntity(ctx context.Context, entityType, entityID string) (*kycv1.KYCVerification, error) {
	return r.getOne(ctx, "entity_type = ? AND entity_id = ?", entityType, entityID)
}

const kycCols = `verification_id, type, entity_type, entity_id, method, provider, provider_reference, status, rejection_reason, verified_by, verified_at, expires_at`

func scanKYCVerification(row interface{ Scan(...any) error }) (*kycv1.KYCVerification, error) {
	var k kycv1.KYCVerification
	var typeStr, methodStr, statusStr string
	var provider, providerRef, rejectionReason, verifiedBy *string
	var verifiedAt, expiresAt *time.Time
	if err := row.Scan(&k.Id, &typeStr, &k.EntityType, &k.EntityId, &methodStr, &provider, &providerRef, &statusStr, &rejectionReason, &verifiedBy, &verifiedAt, &expiresAt); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if provider != nil {
		k.Provider = *provider
	}
	if providerRef != nil {
		k.ProviderReference = *providerRef
	}
	if rejectionReason != nil {
		k.RejectionReason = *rejectionReason
	}
	if verifiedBy != nil {
		k.VerifiedBy = *verifiedBy
	}
	k.Type = verificationTypeFromString(typeStr)
	k.Method = verificationMethodFromString(methodStr)
	k.Status = verificationStatusFromString(statusStr)
	if verifiedAt != nil {
		k.VerifiedAt = timestamppb.New(*verifiedAt)
	}
	if expiresAt != nil {
		k.ExpiresAt = timestamppb.New(*expiresAt)
	}
	return &k, nil
}

func (r *KYCVerificationRepository) ListByStatus(ctx context.Context, status kycv1.VerificationStatus, limit, offset int) ([]*kycv1.KYCVerification, error) {
	q := `select ` + kycCols + ` from authn_schema.kyc_verifications where status = ? order by verification_id desc`
	args := []any{verificationStatusToString(status)}
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
	var out []*kycv1.KYCVerification
	for rows.Next() {
		k, err := scanKYCVerification(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

func (r *KYCVerificationRepository) UpdateStatus(ctx context.Context, id string, status kycv1.VerificationStatus, rejectionReason *string) error {
	upd := map[string]any{"status": verificationStatusToString(status)}
	if rejectionReason != nil {
		upd["rejection_reason"] = *rejectionReason
	}
	return r.db.WithContext(ctx).Table("authn_schema.kyc_verifications").Where("verification_id = ?", id).Updates(upd).Error
}

func (r *KYCVerificationRepository) MarkVerified(ctx context.Context, id, verifiedBy string, verifiedAt time.Time, expiresAt *time.Time) error {
	upd := map[string]any{
		"status":      verificationStatusToString(kycv1.VerificationStatus_VERIFICATION_STATUS_VERIFIED),
		"verified_by": verifiedBy,
		"verified_at": verifiedAt,
	}
	if expiresAt != nil {
		upd["expires_at"] = *expiresAt
	}
	return r.db.WithContext(ctx).Table("authn_schema.kyc_verifications").Where("verification_id = ?", id).Updates(upd).Error
}

func (r *KYCVerificationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("authn_schema.kyc_verifications").Where("verification_id = ?", id).Delete(map[string]any{}).Error
}

func (r *KYCVerificationRepository) getOne(ctx context.Context, where string, args ...any) (*kycv1.KYCVerification, error) {
	q := `select ` + kycCols + ` from authn_schema.kyc_verifications where ` + where + ` limit 1`
	row := r.db.WithContext(ctx).Raw(q, args...).Row()
	if err := row.Err(); err != nil {
		return nil, err
	}
	k, err := scanKYCVerification(row)
	if err != nil {
		return nil, err
	}
	if k.Id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return k, nil
}

func verificationTypeFromString(s string) kycv1.VerificationType {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "VERIFICATION_TYPE_")
	if v, ok := kycv1.VerificationType_value["VERIFICATION_TYPE_"+s]; ok {
		return kycv1.VerificationType(v)
	}
	return kycv1.VerificationType_VERIFICATION_TYPE_UNSPECIFIED
}

func verificationMethodFromString(s string) kycv1.VerificationMethod {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "VERIFICATION_METHOD_")
	if v, ok := kycv1.VerificationMethod_value["VERIFICATION_METHOD_"+s]; ok {
		return kycv1.VerificationMethod(v)
	}
	return kycv1.VerificationMethod_VERIFICATION_METHOD_UNSPECIFIED
}

func verificationStatusToString(s kycv1.VerificationStatus) string {
	out := strings.TrimPrefix(s.String(), "VERIFICATION_STATUS_")
	return out
}

func verificationStatusFromString(s string) kycv1.VerificationStatus {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "VERIFICATION_STATUS_")
	if v, ok := kycv1.VerificationStatus_value["VERIFICATION_STATUS_"+s]; ok {
		return kycv1.VerificationStatus(v)
	}
	return kycv1.VerificationStatus_VERIFICATION_STATUS_UNSPECIFIED
}
