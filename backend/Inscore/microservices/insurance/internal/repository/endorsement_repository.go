package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	endorsementv1 "github.com/newage-saint/insuretech/gen/go/insuretech/endorsement/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type EndorsementRepository struct {
	db *gorm.DB
}

func NewEndorsementRepository(db *gorm.DB) *EndorsementRepository {
	return &EndorsementRepository{db: db}
}

func (r *EndorsementRepository) Create(ctx context.Context, endorsement *endorsementv1.Endorsement) (*endorsementv1.Endorsement, error) {
	if endorsement.Id == "" {
		return nil, fmt.Errorf("endorsement_id is required")
	}

	var createdBy string
	err := r.db.WithContext(ctx).Raw(`SELECT user_id FROM authn_schema.users LIMIT 1`).Scan(&createdBy).Error
	if err != nil || createdBy == "" {
		return nil, fmt.Errorf("failed to get valid user for created_by: %w", err)
	}

	auditInfoJSON := fmt.Sprintf(`{"created_by":"%s","created_at":"%s"}`, createdBy, time.Now().Format(time.RFC3339))

	// Handle Money type
	premiumAdjustment := int64(0)
	premiumAdjustmentCurrency := "BDT"
	if endorsement.PremiumAdjustment != nil {
		premiumAdjustment = endorsement.PremiumAdjustment.Amount
		premiumAdjustmentCurrency = endorsement.PremiumAdjustment.Currency
	}

	var effectiveDate, approvedAt sql.NullTime
	if endorsement.EffectiveDate != nil {
		effectiveDate = sql.NullTime{Time: endorsement.EffectiveDate.AsTime(), Valid: true}
	}
	if endorsement.ApprovedAt != nil {
		approvedAt = sql.NullTime{Time: endorsement.ApprovedAt.AsTime(), Valid: true}
	}

	var approvedBy sql.NullString
	if endorsement.ApprovedBy != "" {
		approvedBy = sql.NullString{String: endorsement.ApprovedBy, Valid: true}
	}

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.endorsements
			(endorsement_id, endorsement_number, policy_id, type, reason, changes,
			 premium_adjustment, premium_adjustment_currency, premium_refund_required,
			 status, requested_by, approved_by, effective_date, approved_at, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		endorsement.Id,
		endorsement.EndorsementNumber,
		endorsement.PolicyId,
		strings.ToUpper(endorsement.Type.String()),
		endorsement.Reason,
		endorsement.Changes,
		premiumAdjustment,
		premiumAdjustmentCurrency,
		endorsement.PremiumRefundRequired,
		strings.ToUpper(endorsement.Status.String()),
		endorsement.RequestedBy,
		approvedBy,
		effectiveDate,
		approvedAt,
		auditInfoJSON,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert endorsement: %w", err)
	}

	return r.GetByID(ctx, endorsement.Id)
}

func (r *EndorsementRepository) GetByID(ctx context.Context, endorsementID string) (*endorsementv1.Endorsement, error) {
	var (
		end                       endorsementv1.Endorsement
		typeStr                   sql.NullString
		statusStr                 sql.NullString
		changes                   sql.NullString
		premiumAdjustment         int64
		premiumAdjustmentCurrency string
		approvedBy                sql.NullString
		effectiveDate             sql.NullTime
		approvedAt                sql.NullTime
		auditInfo                 sql.NullString
		createdAt                 time.Time
		updatedAt                 time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT endorsement_id, endorsement_number, policy_id, type, reason, changes,
		       premium_adjustment, premium_adjustment_currency, premium_refund_required,
		       status, requested_by, approved_by, effective_date, approved_at, audit_info,
		       created_at, updated_at
		FROM insurance_schema.endorsements
		WHERE endorsement_id = $1`,
		endorsementID,
	).Row().Scan(
		&end.Id,
		&end.EndorsementNumber,
		&end.PolicyId,
		&typeStr,
		&end.Reason,
		&changes,
		&premiumAdjustment,
		&premiumAdjustmentCurrency,
		&end.PremiumRefundRequired,
		&statusStr,
		&end.RequestedBy,
		&approvedBy,
		&effectiveDate,
		&approvedAt,
		&auditInfo,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get endorsement: %w", err)
	}

	if typeStr.Valid {
		k := strings.ToUpper(typeStr.String)
		if v, ok := endorsementv1.EndorsementType_value[k]; ok {
			end.Type = endorsementv1.EndorsementType(v)
		}
	}

	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := endorsementv1.EndorsementStatus_value[k]; ok {
			end.Status = endorsementv1.EndorsementStatus(v)
		}
	}

	if changes.Valid {
		end.Changes = changes.String
	}

	end.PremiumAdjustment = &commonv1.Money{
		Amount:   premiumAdjustment,
		Currency: premiumAdjustmentCurrency,
	}

	if approvedBy.Valid {
		end.ApprovedBy = approvedBy.String
	}

	if effectiveDate.Valid {
		end.EffectiveDate = timestamppb.New(effectiveDate.Time)
	}

	if approvedAt.Valid {
		end.ApprovedAt = timestamppb.New(approvedAt.Time)
	}

	if auditInfo.Valid {
		end.AuditInfo = &commonv1.AuditInfo{}
	}

	return &end, nil
}

func (r *EndorsementRepository) Update(ctx context.Context, endorsement *endorsementv1.Endorsement) (*endorsementv1.Endorsement, error) {
	premiumAdjustment := int64(0)
	premiumAdjustmentCurrency := "BDT"
	if endorsement.PremiumAdjustment != nil {
		premiumAdjustment = endorsement.PremiumAdjustment.Amount
		premiumAdjustmentCurrency = endorsement.PremiumAdjustment.Currency
	}

	var effectiveDate, approvedAt sql.NullTime
	if endorsement.EffectiveDate != nil {
		effectiveDate = sql.NullTime{Time: endorsement.EffectiveDate.AsTime(), Valid: true}
	}
	if endorsement.ApprovedAt != nil {
		approvedAt = sql.NullTime{Time: endorsement.ApprovedAt.AsTime(), Valid: true}
	}

	var approvedBy sql.NullString
	if endorsement.ApprovedBy != "" {
		approvedBy = sql.NullString{String: endorsement.ApprovedBy, Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.endorsements
		SET endorsement_number = $2,
		    policy_id = $3,
		    type = $4,
		    reason = $5,
		    changes = $6,
		    premium_adjustment = $7,
		    premium_adjustment_currency = $8,
		    premium_refund_required = $9,
		    status = $10,
		    requested_by = $11,
		    approved_by = $12,
		    effective_date = $13,
		    approved_at = $14,
		    updated_at = NOW()
		WHERE endorsement_id = $1`,
		endorsement.Id,
		endorsement.EndorsementNumber,
		endorsement.PolicyId,
		strings.ToUpper(endorsement.Type.String()),
		endorsement.Reason,
		endorsement.Changes,
		premiumAdjustment,
		premiumAdjustmentCurrency,
		endorsement.PremiumRefundRequired,
		strings.ToUpper(endorsement.Status.String()),
		endorsement.RequestedBy,
		approvedBy,
		effectiveDate,
		approvedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update endorsement: %w", err)
	}

	return r.GetByID(ctx, endorsement.Id)
}

func (r *EndorsementRepository) Delete(ctx context.Context, endorsementID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.endorsements
		WHERE endorsement_id = $1`,
		endorsementID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete endorsement: %w", err)
	}

	return nil
}

func (r *EndorsementRepository) ListByPolicyID(ctx context.Context, policyID string) ([]*endorsementv1.Endorsement, error) {
	query := `
		SELECT endorsement_id, endorsement_number, policy_id, type, reason, changes,
		       premium_adjustment, premium_adjustment_currency, premium_refund_required,
		       status, requested_by, approved_by, effective_date, approved_at, audit_info,
		       created_at, updated_at
		FROM insurance_schema.endorsements
		WHERE policy_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.WithContext(ctx).Raw(query, policyID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list endorsements: %w", err)
	}
	defer rows.Close()

	endorsements := make([]*endorsementv1.Endorsement, 0)
	for rows.Next() {
		var (
			end                       endorsementv1.Endorsement
			typeStr                   sql.NullString
			statusStr                 sql.NullString
			changes                   sql.NullString
			premiumAdjustment         int64
			premiumAdjustmentCurrency string
			approvedBy                sql.NullString
			effectiveDate             sql.NullTime
			approvedAt                sql.NullTime
			auditInfo                 sql.NullString
			createdAt                 time.Time
			updatedAt                 time.Time
		)

		err := rows.Scan(
			&end.Id,
			&end.EndorsementNumber,
			&end.PolicyId,
			&typeStr,
			&end.Reason,
			&changes,
			&premiumAdjustment,
			&premiumAdjustmentCurrency,
			&end.PremiumRefundRequired,
			&statusStr,
			&end.RequestedBy,
			&approvedBy,
			&effectiveDate,
			&approvedAt,
			&auditInfo,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan endorsement: %w", err)
		}

		if typeStr.Valid {
			k := strings.ToUpper(typeStr.String)
			if v, ok := endorsementv1.EndorsementType_value[k]; ok {
				end.Type = endorsementv1.EndorsementType(v)
			}
		}

		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := endorsementv1.EndorsementStatus_value[k]; ok {
				end.Status = endorsementv1.EndorsementStatus(v)
			}
		}

		if changes.Valid {
			end.Changes = changes.String
		}

		end.PremiumAdjustment = &commonv1.Money{
			Amount:   premiumAdjustment,
			Currency: premiumAdjustmentCurrency,
		}

		if approvedBy.Valid {
			end.ApprovedBy = approvedBy.String
		}

		if effectiveDate.Valid {
			end.EffectiveDate = timestamppb.New(effectiveDate.Time)
		}

		if approvedAt.Valid {
			end.ApprovedAt = timestamppb.New(approvedAt.Time)
		}

		if auditInfo.Valid {
			end.AuditInfo = &commonv1.AuditInfo{}
		}

		endorsements = append(endorsements, &end)
	}

	return endorsements, nil
}
