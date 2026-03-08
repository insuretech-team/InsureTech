package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	policyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/policy/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type QuotationRepository struct {
	db *gorm.DB
}

func NewQuotationRepository(db *gorm.DB) *QuotationRepository {
	return &QuotationRepository{db: db}
}

func (r *QuotationRepository) Create(ctx context.Context, quotation *policyv1.Quotation) (*policyv1.Quotation, error) {
	if quotation.QuotationId == "" {
		return nil, fmt.Errorf("quotation_id is required")
	}

	// Handle Money types
	estimatedPremium := int64(0)
	estimatedPremiumCurrency := "BDT"
	if quotation.EstimatedPremium != nil {
		estimatedPremium = quotation.EstimatedPremium.Amount
		estimatedPremiumCurrency = quotation.EstimatedPremium.Currency
	}

	quotedAmount := int64(0)
	quotedAmountCurrency := "BDT"
	if quotation.QuotedAmount != nil {
		quotedAmount = quotation.QuotedAmount.Amount
		quotedAmountCurrency = quotation.QuotedAmount.Currency
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

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.quotations
			(quotation_id, business_id, insurer_name, plan_id, insurance_category, department_id,
			 employee_no, estimated_premium, estimated_premium_currency, quoted_amount, quoted_amount_currency,
			 status, submission_date, valid_until, quotation_number, plan_name,
			 created_by_user_id, approved_by_user_id, approved_at, rejection_reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
		quotation.QuotationId,
		quotation.BusinessId,
		quotation.InsurerName,
		quotation.PlanId,
		strings.ToUpper(quotation.InsuranceCategory.String()),
		quotation.DepartmentId,
		quotation.EmployeeNo,
		estimatedPremium,
		estimatedPremiumCurrency,
		quotedAmount,
		quotedAmountCurrency,
		strings.ToUpper(quotation.Status.String()),
		submissionDate,
		validUntil,
		quotation.QuotationNumber,
		quotation.PlanName,
		quotation.CreatedByUserId,
		quotation.ApprovedByUserId,
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
		quot                      policyv1.Quotation
		insuranceCategoryStr      sql.NullString
		statusStr                 sql.NullString
		businessID                sql.NullString
		insurerName               sql.NullString
		planID                    sql.NullString
		departmentID              sql.NullString
		estimatedPremium          int64
		estimatedPremiumCurrency  string
		quotedAmount              int64
		quotedAmountCurrency      string
		submissionDate            sql.NullTime
		validUntil                sql.NullTime
		quotationNumber           sql.NullString
		planName                  sql.NullString
		createdByUserID           sql.NullString
		approvedByUserID          sql.NullString
		approvedAt                sql.NullTime
		rejectionReason           sql.NullString
		createdAt                 time.Time
		updatedAt                 time.Time
		deletedAt                 sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT quotation_id, business_id, insurer_name, plan_id, insurance_category, department_id,
		       employee_no, estimated_premium, estimated_premium_currency, quoted_amount, quoted_amount_currency,
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
		&estimatedPremium,
		&estimatedPremiumCurrency,
		&quotedAmount,
		&quotedAmountCurrency,
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

	quot.EstimatedPremium = &commonv1.Money{
		Amount:   estimatedPremium,
		Currency: estimatedPremiumCurrency,
	}

	quot.QuotedAmount = &commonv1.Money{
		Amount:   quotedAmount,
		Currency: quotedAmountCurrency,
	}

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
	estimatedPremium := int64(0)
	estimatedPremiumCurrency := "BDT"
	if quotation.EstimatedPremium != nil {
		estimatedPremium = quotation.EstimatedPremium.Amount
		estimatedPremiumCurrency = quotation.EstimatedPremium.Currency
	}

	quotedAmount := int64(0)
	quotedAmountCurrency := "BDT"
	if quotation.QuotedAmount != nil {
		quotedAmount = quotation.QuotedAmount.Amount
		quotedAmountCurrency = quotation.QuotedAmount.Currency
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

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.quotations
		SET business_id = $2,
		    insurer_name = $3,
		    plan_id = $4,
		    insurance_category = $5,
		    department_id = $6,
		    employee_no = $7,
		    estimated_premium = $8,
		    estimated_premium_currency = $9,
		    quoted_amount = $10,
		    quoted_amount_currency = $11,
		    status = $12,
		    submission_date = $13,
		    valid_until = $14,
		    quotation_number = $15,
		    plan_name = $16,
		    created_by_user_id = $17,
		    approved_by_user_id = $18,
		    approved_at = $19,
		    rejection_reason = $20,
		    updated_at = NOW()
		WHERE quotation_id = $1 AND deleted_at IS NULL`,
		quotation.QuotationId,
		quotation.BusinessId,
		quotation.InsurerName,
		quotation.PlanId,
		strings.ToUpper(quotation.InsuranceCategory.String()),
		quotation.DepartmentId,
		quotation.EmployeeNo,
		estimatedPremium,
		estimatedPremiumCurrency,
		quotedAmount,
		quotedAmountCurrency,
		strings.ToUpper(quotation.Status.String()),
		submissionDate,
		validUntil,
		quotation.QuotationNumber,
		quotation.PlanName,
		quotation.CreatedByUserId,
		quotation.ApprovedByUserId,
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

	// Get total count
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

	// Get quotations
	query := `
		SELECT quotation_id, business_id, insurer_name, plan_id, insurance_category, department_id,
		       employee_no, estimated_premium, estimated_premium_currency, quoted_amount, quoted_amount_currency,
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

	var rows *sql.Rows
	var err error
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
			quot                     policyv1.Quotation
			insuranceCategoryStr     sql.NullString
			statusStr                sql.NullString
			businessID               sql.NullString
			insurerName              sql.NullString
			planID                   sql.NullString
			departmentID             sql.NullString
			estimatedPremium         int64
			estimatedPremiumCurrency string
			quotedAmount             int64
			quotedAmountCurrency     string
			submissionDate           sql.NullTime
			validUntil               sql.NullTime
			quotationNumber          sql.NullString
			planName                 sql.NullString
			createdByUserID          sql.NullString
			approvedByUserID         sql.NullString
			approvedAt               sql.NullTime
			rejectionReason          sql.NullString
			createdAt                time.Time
			updatedAt                time.Time
			deletedAt                sql.NullTime
		)

		err := rows.Scan(
			&quot.QuotationId,
			&businessID,
			&insurerName,
			&planID,
			&insuranceCategoryStr,
			&departmentID,
			&quot.EmployeeNo,
			&estimatedPremium,
			&estimatedPremiumCurrency,
			&quotedAmount,
			&quotedAmountCurrency,
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

		quot.EstimatedPremium = &commonv1.Money{
			Amount:   estimatedPremium,
			Currency: estimatedPremiumCurrency,
		}

		quot.QuotedAmount = &commonv1.Money{
			Amount:   quotedAmount,
			Currency: quotedAmountCurrency,
		}

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
