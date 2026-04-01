package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	policyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/policy/entity/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type QuotationRepository struct {
	db *gorm.DB
}

func NewQuotationRepository(db *gorm.DB) *QuotationRepository {
	return &QuotationRepository{db: db}
}

var quotationProtoJSONMarshaler = protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: false}
var quotationProtoJSONUnmarshaler = protojson.UnmarshalOptions{DiscardUnknown: true}
// nullUUID converts an empty string to sql.NullString{Valid:false} so Postgres
// stores NULL instead of an empty string in uuid columns.
func nullUUID(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func scanQuotationMoney(raw sql.NullString) *commonv1.Money {
	if !raw.Valid || strings.TrimSpace(raw.String) == "" || raw.String == "null" {
		return nil
	}

	var money commonv1.Money
	if err := quotationProtoJSONUnmarshaler.Unmarshal([]byte(raw.String), &money); err != nil {
		_ = json.Unmarshal([]byte(raw.String), &money)
	}

	return &money
}

func marshalQuotationMoney(m *commonv1.Money) (string, error) {
	if m == nil {
		return "null", nil
	}

	b, err := quotationProtoJSONMarshaler.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (r *QuotationRepository) Create(ctx context.Context, quotation *policyv1.Quotation) (*policyv1.Quotation, error) {
	if quotation.QuotationId == "" {
		return nil, fmt.Errorf("quotation_id is required")
	}

	estimatedPremiumJSON, err := marshalQuotationMoney(quotation.EstimatedPremium)
	if err != nil {
		return nil, fmt.Errorf("marshal estimated_premium: %w", err)
	}

	quotedAmountJSON, err := marshalQuotationMoney(quotation.QuotedAmount)
	if err != nil {
		return nil, fmt.Errorf("marshal quoted_amount: %w", err)
	}

	var submissionDate, validUntil, approvedAt sql.NullTime
	if quotation.SubmissionDate != nil {
		submissionDate = sql.NullTime{Time: quotation.SubmissionDate.AsTime(), Valid: true}
	}
	if quotation.ValidUntil != nil {
		validUntil = sql.NullTime{Time: quotation.ValidUntil.AsTime(), Valid: true}
	}
	if quotation.ApprovedAt != nil {
		approvedAt = sql.NullTime{Time: quotation.ApprovedAt.AsTime(), Valid: true}
	}

	// Use the proto enum string as-is (DB stores full proto enum names like "INSURANCE_TYPE_LIFE")
	// Default to INSURANCE_TYPE_LIFE for B2C quotations when unspecified
	categoryStr := strings.ToUpper(quotation.InsuranceCategory.String())
	if categoryStr == "" || strings.HasSuffix(categoryStr, "_UNSPECIFIED") {
		categoryStr = "INSURANCE_TYPE_LIFE" // Default for B2C quotations
	}

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.quotations
			(quotation_id, business_id, insurer_name, plan_id, insurance_category, department_id,
			 employee_no, estimated_premium, quoted_amount,
			 status, submission_date, valid_until, quotation_number, plan_name,
			 created_by_user_id, approved_by_user_id, approved_at, rejection_reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb, $9::jsonb, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		quotation.QuotationId,
		nullUUID(quotation.BusinessId),
		quotation.InsurerName,
		nullUUID(quotation.PlanId),
		categoryStr,
		nullUUID(quotation.DepartmentId),
		quotation.EmployeeNo,
		estimatedPremiumJSON,
		quotedAmountJSON,
		strings.ToUpper(quotation.Status.String()),
		submissionDate,
		validUntil,
		quotation.QuotationNumber,
		quotation.PlanName,
		nullUUID(quotation.CreatedByUserId),
		nullUUID(quotation.ApprovedByUserId),
		approvedAt,
		quotation.RejectionReason,
	).Error
	if err != nil {
		return nil, fmt.Errorf("failed to insert quotation: %w", err)
	}

	return r.GetByID(ctx, quotation.QuotationId)
}

func (r *QuotationRepository) GetByID(ctx context.Context, quotationID string) (*policyv1.Quotation, error) {
	var (
		quot                 policyv1.Quotation
		insuranceCategoryStr sql.NullString
		statusStr            sql.NullString
		businessID           sql.NullString
		insurerName          sql.NullString
		planID               sql.NullString
		departmentID         sql.NullString
		estimatedPremiumJSON sql.NullString
		quotedAmountJSON     sql.NullString
		submissionDate       sql.NullTime
		validUntil           sql.NullTime
		quotationNumber      sql.NullString
		planName             sql.NullString
		createdByUserID      sql.NullString
		approvedByUserID     sql.NullString
		approvedAt           sql.NullTime
		rejectionReason      sql.NullString
		createdAt            time.Time
		updatedAt            time.Time
		deletedAt            sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT quotation_id, business_id, insurer_name, plan_id, insurance_category, department_id,
		       employee_no, COALESCE(estimated_premium::text, 'null') AS estimated_premium,
		       COALESCE(quoted_amount::text, 'null') AS quoted_amount,
		       status, submission_date, valid_until, quotation_number, plan_name,
		       created_by_user_id, approved_by_user_id, approved_at, rejection_reason,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.quotations
		WHERE quotation_id = $1 AND deleted_at IS NULL`,
		quotationID,
	).Row().Scan(
		&quot.QuotationId,
		&businessID,
		&insurerName,
		&planID,
		&insuranceCategoryStr,
		&departmentID,
		&quot.EmployeeNo,
		&estimatedPremiumJSON,
		&quotedAmountJSON,
		&statusStr,
		&submissionDate,
		&validUntil,
		&quotationNumber,
		&planName,
		&createdByUserID,
		&approvedByUserID,
		&approvedAt,
		&rejectionReason,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get quotation: %w", err)
	}

	if businessID.Valid {
		quot.BusinessId = businessID.String
	}
	if insurerName.Valid {
		quot.InsurerName = insurerName.String
	}
	if planID.Valid {
		quot.PlanId = planID.String
	}
	if departmentID.Valid {
		quot.DepartmentId = departmentID.String
	}
	if quotationNumber.Valid {
		quot.QuotationNumber = quotationNumber.String
	}
	if planName.Valid {
		quot.PlanName = planName.String
	}
	if createdByUserID.Valid {
		quot.CreatedByUserId = createdByUserID.String
	}
	if approvedByUserID.Valid {
		quot.ApprovedByUserId = approvedByUserID.String
	}
	if rejectionReason.Valid {
		quot.RejectionReason = rejectionReason.String
	}

	if insuranceCategoryStr.Valid {
		k := strings.ToUpper(insuranceCategoryStr.String)
		if v, ok := commonv1.InsuranceType_value[k]; ok {
			quot.InsuranceCategory = commonv1.InsuranceType(v)
		}
	}

	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := policyv1.QuotationStatus_value[k]; ok {
			quot.Status = policyv1.QuotationStatus(v)
		}
	}

	quot.EstimatedPremium = scanQuotationMoney(estimatedPremiumJSON)
	quot.QuotedAmount = scanQuotationMoney(quotedAmountJSON)

	if submissionDate.Valid {
		quot.SubmissionDate = timestamppb.New(submissionDate.Time)
	}
	if validUntil.Valid {
		quot.ValidUntil = timestamppb.New(validUntil.Time)
	}
	if approvedAt.Valid {
		quot.ApprovedAt = timestamppb.New(approvedAt.Time)
	}

	quot.CreatedAt = timestamppb.New(createdAt)
	quot.UpdatedAt = timestamppb.New(updatedAt)

	return &quot, nil
}

func (r *QuotationRepository) Update(ctx context.Context, quotation *policyv1.Quotation) (*policyv1.Quotation, error) {
	estimatedPremiumJSON, err := marshalQuotationMoney(quotation.EstimatedPremium)
	if err != nil {
		return nil, fmt.Errorf("marshal estimated_premium: %w", err)
	}

	quotedAmountJSON, err := marshalQuotationMoney(quotation.QuotedAmount)
	if err != nil {
		return nil, fmt.Errorf("marshal quoted_amount: %w", err)
	}

	var submissionDate, validUntil, approvedAt sql.NullTime
	if quotation.SubmissionDate != nil {
		submissionDate = sql.NullTime{Time: quotation.SubmissionDate.AsTime(), Valid: true}
	}
	if quotation.ValidUntil != nil {
		validUntil = sql.NullTime{Time: quotation.ValidUntil.AsTime(), Valid: true}
	}
	if quotation.ApprovedAt != nil {
		approvedAt = sql.NullTime{Time: quotation.ApprovedAt.AsTime(), Valid: true}
	}

	err = r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.quotations
		SET business_id = $2,
		    insurer_name = $3,
		    plan_id = $4,
		    insurance_category = $5,
		    department_id = $6,
		    employee_no = $7,
		    estimated_premium = $8::jsonb,
		    quoted_amount = $9::jsonb,
		    status = $10,
		    submission_date = $11,
		    valid_until = $12,
		    quotation_number = $13,
		    plan_name = $14,
		    created_by_user_id = $15,
		    approved_by_user_id = $16,
		    approved_at = $17,
		    rejection_reason = $18,
		    updated_at = NOW()
		WHERE quotation_id = $1 AND deleted_at IS NULL`,
		quotation.QuotationId,
		nullUUID(quotation.BusinessId),
		quotation.InsurerName,
		nullUUID(quotation.PlanId),
		strings.ToUpper(quotation.InsuranceCategory.String()),
		nullUUID(quotation.DepartmentId),
		quotation.EmployeeNo,
		estimatedPremiumJSON,
		quotedAmountJSON,
		strings.ToUpper(quotation.Status.String()),
		submissionDate,
		validUntil,
		quotation.QuotationNumber,
		quotation.PlanName,
		nullUUID(quotation.CreatedByUserId),
		nullUUID(quotation.ApprovedByUserId),
		approvedAt,
		quotation.RejectionReason,
	).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update quotation: %w", err)
	}

	return r.GetByID(ctx, quotation.QuotationId)
}

func (r *QuotationRepository) Delete(ctx context.Context, quotationID string) error {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.quotations
		SET deleted_at = NOW()
		WHERE quotation_id = $1 AND deleted_at IS NULL`,
		quotationID,
	).Error
	if err != nil {
		return fmt.Errorf("failed to delete quotation: %w", err)
	}

	return nil
}

func (r *QuotationRepository) List(ctx context.Context, businessID string, page, pageSize int) ([]*policyv1.Quotation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	var total int64
	countQuery := `SELECT COUNT(*) FROM insurance_schema.quotations WHERE deleted_at IS NULL`
	if businessID != "" {
		countQuery += ` AND business_id = $1`
		err := r.db.WithContext(ctx).Raw(countQuery, businessID).Scan(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count quotations: %w", err)
		}
	} else {
		err := r.db.WithContext(ctx).Raw(countQuery).Scan(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count quotations: %w", err)
		}
	}

	query := `
		SELECT quotation_id, business_id, insurer_name, plan_id, insurance_category, department_id,
		       employee_no, COALESCE(estimated_premium::text, 'null') AS estimated_premium,
		       COALESCE(quoted_amount::text, 'null') AS quoted_amount,
		       status, submission_date, valid_until, quotation_number, plan_name,
		       created_by_user_id, approved_by_user_id, approved_at, rejection_reason,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.quotations
		WHERE deleted_at IS NULL`
	if businessID != "" {
		query += ` AND business_id = $1`
		query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, pageSize, offset)
	} else {
		query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, pageSize, offset)
	}

	var (
		rows *sql.Rows
		err  error
	)
	if businessID != "" {
		rows, err = r.db.WithContext(ctx).Raw(query, businessID).Rows()
	} else {
		rows, err = r.db.WithContext(ctx).Raw(query).Rows()
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list quotations: %w", err)
	}
	defer rows.Close()

	quotations := make([]*policyv1.Quotation, 0)
	for rows.Next() {
		var (
			quot                 policyv1.Quotation
			insuranceCategoryStr sql.NullString
			statusStr            sql.NullString
			rowBusinessID        sql.NullString
			insurerName          sql.NullString
			planID               sql.NullString
			departmentID         sql.NullString
			estimatedPremiumJSON sql.NullString
			quotedAmountJSON     sql.NullString
			submissionDate       sql.NullTime
			validUntil           sql.NullTime
			quotationNumber      sql.NullString
			planName             sql.NullString
			createdByUserID      sql.NullString
			approvedByUserID     sql.NullString
			approvedAt           sql.NullTime
			rejectionReason      sql.NullString
			createdAt            time.Time
			updatedAt            time.Time
			deletedAt            sql.NullTime
		)

		err := rows.Scan(
			&quot.QuotationId,
			&rowBusinessID,
			&insurerName,
			&planID,
			&insuranceCategoryStr,
			&departmentID,
			&quot.EmployeeNo,
			&estimatedPremiumJSON,
			&quotedAmountJSON,
			&statusStr,
			&submissionDate,
			&validUntil,
			&quotationNumber,
			&planName,
			&createdByUserID,
			&approvedByUserID,
			&approvedAt,
			&rejectionReason,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan quotation: %w", err)
		}

		if rowBusinessID.Valid {
			quot.BusinessId = rowBusinessID.String
		}
		if insurerName.Valid {
			quot.InsurerName = insurerName.String
		}
		if planID.Valid {
			quot.PlanId = planID.String
		}
		if departmentID.Valid {
			quot.DepartmentId = departmentID.String
		}
		if quotationNumber.Valid {
			quot.QuotationNumber = quotationNumber.String
		}
		if planName.Valid {
			quot.PlanName = planName.String
		}
		if createdByUserID.Valid {
			quot.CreatedByUserId = createdByUserID.String
		}
		if approvedByUserID.Valid {
			quot.ApprovedByUserId = approvedByUserID.String
		}
		if rejectionReason.Valid {
			quot.RejectionReason = rejectionReason.String
		}

		if insuranceCategoryStr.Valid {
			k := strings.ToUpper(insuranceCategoryStr.String)
			if v, ok := commonv1.InsuranceType_value[k]; ok {
				quot.InsuranceCategory = commonv1.InsuranceType(v)
			}
		}

		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := policyv1.QuotationStatus_value[k]; ok {
				quot.Status = policyv1.QuotationStatus(v)
			}
		}

		quot.EstimatedPremium = scanQuotationMoney(estimatedPremiumJSON)
		quot.QuotedAmount = scanQuotationMoney(quotedAmountJSON)

		if submissionDate.Valid {
			quot.SubmissionDate = timestamppb.New(submissionDate.Time)
		}
		if validUntil.Valid {
			quot.ValidUntil = timestamppb.New(validUntil.Time)
		}
		if approvedAt.Valid {
			quot.ApprovedAt = timestamppb.New(approvedAt.Time)
		}

		quot.CreatedAt = timestamppb.New(createdAt)
		quot.UpdatedAt = timestamppb.New(updatedAt)

		quotations = append(quotations, &quot)
	}

	return quotations, total, nil
}


