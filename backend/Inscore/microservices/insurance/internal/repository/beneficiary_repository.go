package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	beneficiaryv1 "github.com/newage-saint/insuretech/gen/go/insuretech/beneficiary/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type BeneficiaryRepository struct {
	db *gorm.DB
}

func NewBeneficiaryRepository(db *gorm.DB) *BeneficiaryRepository {
	return &BeneficiaryRepository{db: db}
}

func (r *BeneficiaryRepository) Create(ctx context.Context, beneficiary *beneficiaryv1.Beneficiary) (*beneficiaryv1.Beneficiary, error) {
	if beneficiary.BeneficiaryId == "" {
		return nil, fmt.Errorf("beneficiary_id is required")
	}

	var createdBy string
	err := r.db.WithContext(ctx).Raw(`SELECT user_id FROM authn_schema.users LIMIT 1`).Scan(&createdBy).Error
	if err != nil || createdBy == "" {
		return nil, fmt.Errorf("failed to get valid user for created_by: %w", err)
	}

	auditInfo := &commonv1.AuditInfo{
		CreatedBy: createdBy,
		CreatedAt: timestamppb.Now(),
	}
	auditInfoJSON := fmt.Sprintf(`{"created_by":"%s","created_at":"%s"}`, auditInfo.CreatedBy, auditInfo.CreatedAt.AsTime().Format(time.RFC3339))

	var kycCompletedAt sql.NullTime
	if beneficiary.KycCompletedAt != nil {
		kycCompletedAt = sql.NullTime{Time: beneficiary.KycCompletedAt.AsTime(), Valid: true}
	}

	var referredBy, partnerID sql.NullString
	if beneficiary.ReferredBy != "" {
		referredBy = sql.NullString{String: beneficiary.ReferredBy, Valid: true}
	}
	if beneficiary.PartnerId != "" {
		partnerID = sql.NullString{String: beneficiary.PartnerId, Valid: true}
	}

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.beneficiaries
			(beneficiary_id, user_id, type, code, status, kyc_status, kyc_completed_at, 
			 risk_score, referral_code, referred_by, partner_id, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		beneficiary.BeneficiaryId,
		beneficiary.UserId,
		strings.ToUpper(beneficiary.Type.String()),
		beneficiary.Code,
		strings.ToUpper(beneficiary.Status.String()),
		strings.ToUpper(beneficiary.KycStatus.String()),
		kycCompletedAt,
		beneficiary.RiskScore,
		beneficiary.ReferralCode,
		referredBy,
		partnerID,
		auditInfoJSON,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert beneficiary: %w", err)
	}

	return r.GetByID(ctx, beneficiary.BeneficiaryId)
}

func (r *BeneficiaryRepository) GetByID(ctx context.Context, beneficiaryID string) (*beneficiaryv1.Beneficiary, error) {
	var (
		b              beneficiaryv1.Beneficiary
		typeStr        sql.NullString
		statusStr      sql.NullString
		kycStatusStr   sql.NullString
		kycCompletedAt sql.NullTime
		riskScore      sql.NullString
		referralCode   sql.NullString
		referredBy     sql.NullString
		partnerID      sql.NullString
		auditInfo      sql.NullString
		createdAt      time.Time
		updatedAt      time.Time
		deletedAt      sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT beneficiary_id, user_id, type, code, status, kyc_status, kyc_completed_at,
		       risk_score, referral_code, referred_by, partner_id, audit_info,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.beneficiaries
		WHERE beneficiary_id = $1 AND deleted_at IS NULL`,
		beneficiaryID,
	).Row().Scan(
		&b.BeneficiaryId,
		&b.UserId,
		&typeStr,
		&b.Code,
		&statusStr,
		&kycStatusStr,
		&kycCompletedAt,
		&riskScore,
		&referralCode,
		&referredBy,
		&partnerID,
		&auditInfo,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get beneficiary: %w", err)
	}

	if typeStr.Valid {
		k := strings.ToUpper(typeStr.String)
		if v, ok := beneficiaryv1.BeneficiaryType_value[k]; ok {
			b.Type = beneficiaryv1.BeneficiaryType(v)
		}
	}

	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := beneficiaryv1.BeneficiaryStatus_value[k]; ok {
			b.Status = beneficiaryv1.BeneficiaryStatus(v)
		}
	}

	if kycStatusStr.Valid {
		k := strings.ToUpper(kycStatusStr.String)
		if v, ok := beneficiaryv1.KYCStatus_value[k]; ok {
			b.KycStatus = beneficiaryv1.KYCStatus(v)
		}
	}

	if kycCompletedAt.Valid {
		b.KycCompletedAt = timestamppb.New(kycCompletedAt.Time)
	}

	if riskScore.Valid {
		b.RiskScore = riskScore.String
	}

	if referralCode.Valid {
		b.ReferralCode = referralCode.String
	}

	if referredBy.Valid {
		b.ReferredBy = referredBy.String
	}

	if partnerID.Valid {
		b.PartnerId = partnerID.String
	}

	if auditInfo.Valid {
		b.AuditInfo = &commonv1.AuditInfo{}
	}

	return &b, nil
}

func (r *BeneficiaryRepository) Update(ctx context.Context, beneficiary *beneficiaryv1.Beneficiary) (*beneficiaryv1.Beneficiary, error) {
	var kycCompletedAt sql.NullTime
	if beneficiary.KycCompletedAt != nil {
		kycCompletedAt = sql.NullTime{Time: beneficiary.KycCompletedAt.AsTime(), Valid: true}
	}

	var referredBy, partnerID sql.NullString
	if beneficiary.ReferredBy != "" {
		referredBy = sql.NullString{String: beneficiary.ReferredBy, Valid: true}
	}
	if beneficiary.PartnerId != "" {
		partnerID = sql.NullString{String: beneficiary.PartnerId, Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.beneficiaries
		SET user_id = $2,
		    type = $3,
		    code = $4,
		    status = $5,
		    kyc_status = $6,
		    kyc_completed_at = $7,
		    risk_score = $8,
		    referral_code = $9,
		    referred_by = $10,
		    partner_id = $11,
		    updated_at = NOW()
		WHERE beneficiary_id = $1 AND deleted_at IS NULL`,
		beneficiary.BeneficiaryId,
		beneficiary.UserId,
		strings.ToUpper(beneficiary.Type.String()),
		beneficiary.Code,
		strings.ToUpper(beneficiary.Status.String()),
		strings.ToUpper(beneficiary.KycStatus.String()),
		kycCompletedAt,
		beneficiary.RiskScore,
		beneficiary.ReferralCode,
		referredBy,
		partnerID,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update beneficiary: %w", err)
	}

	return r.GetByID(ctx, beneficiary.BeneficiaryId)
}

func (r *BeneficiaryRepository) Delete(ctx context.Context, beneficiaryID string) error {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.beneficiaries
		SET deleted_at = NOW()
		WHERE beneficiary_id = $1 AND deleted_at IS NULL`,
		beneficiaryID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete beneficiary: %w", err)
	}

	return nil
}

func (r *BeneficiaryRepository) List(ctx context.Context, page, pageSize int) ([]*beneficiaryv1.Beneficiary, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	var total int64
	err := r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM insurance_schema.beneficiaries WHERE deleted_at IS NULL`).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count beneficiaries: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT beneficiary_id, user_id, type, code, status, kyc_status, kyc_completed_at,
		       risk_score, referral_code, referred_by, partner_id, audit_info,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.beneficiaries
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC LIMIT %d OFFSET %d`, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list beneficiaries: %w", err)
	}
	defer rows.Close()

	beneficiaries := make([]*beneficiaryv1.Beneficiary, 0)
	for rows.Next() {
		var (
			b              beneficiaryv1.Beneficiary
			typeStr        sql.NullString
			statusStr      sql.NullString
			kycStatusStr   sql.NullString
			kycCompletedAt sql.NullTime
			riskScore      sql.NullString
			referralCode   sql.NullString
			referredBy     sql.NullString
			partnerID      sql.NullString
			auditInfo      sql.NullString
			createdAt      time.Time
			updatedAt      time.Time
			deletedAt      sql.NullTime
		)

		err := rows.Scan(
			&b.BeneficiaryId,
			&b.UserId,
			&typeStr,
			&b.Code,
			&statusStr,
			&kycStatusStr,
			&kycCompletedAt,
			&riskScore,
			&referralCode,
			&referredBy,
			&partnerID,
			&auditInfo,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan beneficiary: %w", err)
		}

		if typeStr.Valid {
			k := strings.ToUpper(typeStr.String)
			if v, ok := beneficiaryv1.BeneficiaryType_value[k]; ok {
				b.Type = beneficiaryv1.BeneficiaryType(v)
			}
		}

		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := beneficiaryv1.BeneficiaryStatus_value[k]; ok {
				b.Status = beneficiaryv1.BeneficiaryStatus(v)
			}
		}

		if kycStatusStr.Valid {
			k := strings.ToUpper(kycStatusStr.String)
			if v, ok := beneficiaryv1.KYCStatus_value[k]; ok {
				b.KycStatus = beneficiaryv1.KYCStatus(v)
			}
		}

		if kycCompletedAt.Valid {
			b.KycCompletedAt = timestamppb.New(kycCompletedAt.Time)
		}

		if riskScore.Valid {
			b.RiskScore = riskScore.String
		}

		if referralCode.Valid {
			b.ReferralCode = referralCode.String
		}

		if referredBy.Valid {
			b.ReferredBy = referredBy.String
		}

		if partnerID.Valid {
			b.PartnerId = partnerID.String
		}

		if auditInfo.Valid {
			b.AuditInfo = &commonv1.AuditInfo{}
		}

		beneficiaries = append(beneficiaries, &b)
	}

	return beneficiaries, total, nil
}
