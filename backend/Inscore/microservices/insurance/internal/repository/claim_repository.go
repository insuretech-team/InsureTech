package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	claimsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/claims/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type ClaimRepository struct {
	db *gorm.DB
}

func NewClaimRepository(db *gorm.DB) *ClaimRepository {
	return &ClaimRepository{db: db}
}

func (r *ClaimRepository) Create(ctx context.Context, claim *claimsv1.Claim) (*claimsv1.Claim, error) {
	if claim.ClaimId == "" {
		return nil, fmt.Errorf("claim_id is required")
	}
	
	// Extract Money values
	claimedAmount := int64(0)
	claimedCurrency := "BDT"
	if claim.ClaimedAmount != nil {
		claimedAmount = claim.ClaimedAmount.Amount
		claimedCurrency = claim.ClaimedAmount.Currency
	}
	
	approvedAmount := int64(0)
	approvedCurrency := "BDT"
	if claim.ApprovedAmount != nil {
		approvedAmount = claim.ApprovedAmount.Amount
		approvedCurrency = claim.ApprovedAmount.Currency
	}
	
	settledAmount := int64(0)
	settledCurrency := "BDT"
	if claim.SettledAmount != nil {
		settledAmount = claim.SettledAmount.Amount
		settledCurrency = claim.SettledAmount.Currency
	}
	
	deductibleAmount := int64(0)
	if claim.DeductibleAmount != nil {
		deductibleAmount = claim.DeductibleAmount.Amount
	}
	
	coPayAmount := int64(0)
	if claim.CoPayAmount != nil {
		coPayAmount = claim.CoPayAmount.Amount
	}
	
	// Handle timestamps
	var incidentDate, submittedAt, approvedAt, settledAt interface{}
	if claim.IncidentDate != nil {
		incidentDate = claim.IncidentDate.AsTime()
	}
	if claim.SubmittedAt != nil {
		submittedAt = claim.SubmittedAt.AsTime()
	}
	if claim.ApprovedAt != nil {
		approvedAt = claim.ApprovedAt.AsTime()
	}
	if claim.SettledAt != nil {
		settledAt = claim.SettledAt.AsTime()
	}
	
	// Handle JSONB in_app_messages
	var inAppMessages interface{}
	if claim.InAppMessages != "" {
		inAppMessages = claim.InAppMessages
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.claims
			(claim_id, claim_number, policy_id, customer_id, status, type,
			 claimed_amount, approved_amount, settled_amount,
			 claimed_currency, approved_currency, settled_currency,
			 incident_date, incident_description, submitted_at, approved_at, settled_at,
			 rejection_reason, place_of_incident, bank_details_for_payout,
			 appeal_option_available, in_app_messages, processing_type,
			 deductible_amount, co_pay_amount, processor_notes,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, NOW(), NOW())`,
		claim.ClaimId,
		claim.ClaimNumber,
		claim.PolicyId,
		claim.CustomerId,
		strings.ToUpper(claim.Status.String()),
		strings.ToUpper(claim.Type.String()),
		claimedAmount,
		approvedAmount,
		settledAmount,
		claimedCurrency,
		approvedCurrency,
		settledCurrency,
		incidentDate,
		claim.IncidentDescription,
		submittedAt,
		approvedAt,
		settledAt,
		claim.RejectionReason,
		claim.PlaceOfIncident,
		claim.BankDetailsForPayout,
		claim.AppealOptionAvailable,
		inAppMessages,
		strings.ToUpper(claim.ProcessingType.String()),
		deductibleAmount,
		coPayAmount,
		claim.ProcessorNotes,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert claim: %w", err)
	}

	return r.GetByID(ctx, claim.ClaimId)
}

func (r *ClaimRepository) GetByID(ctx context.Context, claimID string) (*claimsv1.Claim, error) {
	var (
		c                  claimsv1.Claim
		statusStr          sql.NullString
		typeStr            sql.NullString
		processingTypeStr  sql.NullString
		claimedAmount      int64
		approvedAmount     sql.NullInt64
		settledAmount      sql.NullInt64
		claimedCurrency    string
		approvedCurrency   string
		settledCurrency    string
		deductibleAmount   sql.NullInt64
		coPayAmount        sql.NullInt64
		incidentDate       time.Time
		submittedAt        time.Time
		approvedAt         sql.NullTime
		settledAt          sql.NullTime
		createdAt          time.Time
		updatedAt          time.Time
		deletedAt          sql.NullTime
		inAppMessages      sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT claim_id, claim_number, policy_id, customer_id, status, type,
		       claimed_amount, approved_amount, settled_amount,
		       claimed_currency, approved_currency, settled_currency,
		       incident_date, incident_description, submitted_at, approved_at, settled_at,
		       COALESCE(rejection_reason, '') as rejection_reason,
		       COALESCE(place_of_incident, '') as place_of_incident,
		       COALESCE(bank_details_for_payout, '') as bank_details_for_payout,
		       COALESCE(appeal_option_available, false) as appeal_option_available,
		       in_app_messages, processing_type,
		       deductible_amount, co_pay_amount,
		       COALESCE(processor_notes, '') as processor_notes,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.claims
		WHERE claim_id = $1 AND deleted_at IS NULL`,
		claimID,
	).Row().Scan(
		&c.ClaimId,
		&c.ClaimNumber,
		&c.PolicyId,
		&c.CustomerId,
		&statusStr,
		&typeStr,
		&claimedAmount,
		&approvedAmount,
		&settledAmount,
		&claimedCurrency,
		&approvedCurrency,
		&settledCurrency,
		&incidentDate,
		&c.IncidentDescription,
		&submittedAt,
		&approvedAt,
		&settledAt,
		&c.RejectionReason,
		&c.PlaceOfIncident,
		&c.BankDetailsForPayout,
		&c.AppealOptionAvailable,
		&inAppMessages,
		&processingTypeStr,
		&deductibleAmount,
		&coPayAmount,
		&c.ProcessorNotes,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get claim: %w", err)
	}

	// Set Money fields
	c.ClaimedAmount = &commonv1.Money{
		Amount:   claimedAmount,
		Currency: claimedCurrency,
	}
	
	// Set currency companion fields
	c.ClaimedCurrency = claimedCurrency
	c.ApprovedCurrency = approvedCurrency
	c.SettledCurrency = settledCurrency
	
	if approvedAmount.Valid {
		c.ApprovedAmount = &commonv1.Money{
			Amount:   approvedAmount.Int64,
			Currency: approvedCurrency,
		}
	}
	if settledAmount.Valid {
		c.SettledAmount = &commonv1.Money{
			Amount:   settledAmount.Int64,
			Currency: settledCurrency,
		}
	}
	if deductibleAmount.Valid {
		c.DeductibleAmount = &commonv1.Money{
			Amount:   deductibleAmount.Int64,
			Currency: "BDT",
		}
	}
	if coPayAmount.Valid {
		c.CoPayAmount = &commonv1.Money{
			Amount:   coPayAmount.Int64,
			Currency: "BDT",
		}
	}
	
	// Set in_app_messages
	if inAppMessages.Valid {
		c.InAppMessages = inAppMessages.String
	}

	// Parse enums
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := claimsv1.ClaimStatus_value[k]; ok {
			c.Status = claimsv1.ClaimStatus(v)
		}
	}
	if typeStr.Valid {
		k := strings.ToUpper(typeStr.String)
		if v, ok := claimsv1.ClaimType_value[k]; ok {
			c.Type = claimsv1.ClaimType(v)
		}
	}
	if processingTypeStr.Valid {
		k := strings.ToUpper(processingTypeStr.String)
		if v, ok := claimsv1.ClaimProcessingType_value[k]; ok {
			c.ProcessingType = claimsv1.ClaimProcessingType(v)
		}
	}

	if !incidentDate.IsZero() {
		c.IncidentDate = timestamppb.New(incidentDate)
	}
	if !submittedAt.IsZero() {
		c.SubmittedAt = timestamppb.New(submittedAt)
	}
	if approvedAt.Valid {
		c.ApprovedAt = timestamppb.New(approvedAt.Time)
	}
	if settledAt.Valid {
		c.SettledAt = timestamppb.New(settledAt.Time)
	}
	if !createdAt.IsZero() {
		c.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		c.UpdatedAt = timestamppb.New(updatedAt)
	}
	if deletedAt.Valid {
		c.DeletedAt = timestamppb.New(deletedAt.Time)
	}

	return &c, nil
}

func (r *ClaimRepository) Update(ctx context.Context, claim *claimsv1.Claim) (*claimsv1.Claim, error) {
	// Extract Money values
	claimedAmount := int64(0)
	claimedCurrency := "BDT"
	if claim.ClaimedAmount != nil {
		claimedAmount = claim.ClaimedAmount.Amount
		claimedCurrency = claim.ClaimedAmount.Currency
	}
	
	approvedAmount := int64(0)
	approvedCurrency := "BDT"
	if claim.ApprovedAmount != nil {
		approvedAmount = claim.ApprovedAmount.Amount
		approvedCurrency = claim.ApprovedAmount.Currency
	}
	
	settledAmount := int64(0)
	settledCurrency := "BDT"
	if claim.SettledAmount != nil {
		settledAmount = claim.SettledAmount.Amount
		settledCurrency = claim.SettledAmount.Currency
	}
	
	deductibleAmount := int64(0)
	if claim.DeductibleAmount != nil {
		deductibleAmount = claim.DeductibleAmount.Amount
	}
	
	coPayAmount := int64(0)
	if claim.CoPayAmount != nil {
		coPayAmount = claim.CoPayAmount.Amount
	}
	
	// Handle timestamps
	var incidentDate, submittedAt, approvedAt, settledAt interface{}
	if claim.IncidentDate != nil {
		incidentDate = claim.IncidentDate.AsTime()
	}
	if claim.SubmittedAt != nil {
		submittedAt = claim.SubmittedAt.AsTime()
	}
	if claim.ApprovedAt != nil {
		approvedAt = claim.ApprovedAt.AsTime()
	}
	if claim.SettledAt != nil {
		settledAt = claim.SettledAt.AsTime()
	}
	
	// Handle JSONB in_app_messages
	var inAppMessages interface{}
	if claim.InAppMessages != "" {
		inAppMessages = claim.InAppMessages
	}
	
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.claims
		SET claim_number = $2,
		    policy_id = $3,
		    customer_id = $4,
		    status = $5,
		    type = $6,
		    claimed_amount = $7,
		    approved_amount = $8,
		    settled_amount = $9,
		    claimed_currency = $10,
		    approved_currency = $11,
		    settled_currency = $12,
		    incident_date = $13,
		    incident_description = $14,
		    submitted_at = $15,
		    approved_at = $16,
		    settled_at = $17,
		    rejection_reason = $18,
		    place_of_incident = $19,
		    bank_details_for_payout = $20,
		    appeal_option_available = $21,
		    in_app_messages = $22,
		    processing_type = $23,
		    deductible_amount = $24,
		    co_pay_amount = $25,
		    processor_notes = $26,
		    updated_at = NOW()
		WHERE claim_id = $1 AND deleted_at IS NULL`,
		claim.ClaimId,
		claim.ClaimNumber,
		claim.PolicyId,
		claim.CustomerId,
		strings.ToUpper(claim.Status.String()),
		strings.ToUpper(claim.Type.String()),
		claimedAmount,
		approvedAmount,
		settledAmount,
		claimedCurrency,
		approvedCurrency,
		settledCurrency,
		incidentDate,
		claim.IncidentDescription,
		submittedAt,
		approvedAt,
		settledAt,
		claim.RejectionReason,
		claim.PlaceOfIncident,
		claim.BankDetailsForPayout,
		claim.AppealOptionAvailable,
		inAppMessages,
		strings.ToUpper(claim.ProcessingType.String()),
		deductibleAmount,
		coPayAmount,
		claim.ProcessorNotes,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update claim: %w", err)
	}

	return r.GetByID(ctx, claim.ClaimId)
}

func (r *ClaimRepository) List(ctx context.Context, policyID, customerID string, page, pageSize int) ([]*claimsv1.Claim, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// Build count query
	countQuery := `SELECT COUNT(*) FROM insurance_schema.claims WHERE deleted_at IS NULL`
	var countArgs []interface{}
	argIdx := 1
	
	if policyID != "" {
		countQuery += fmt.Sprintf(` AND policy_id = $%d`, argIdx)
		countArgs = append(countArgs, policyID)
		argIdx++
	}
	if customerID != "" {
		countQuery += fmt.Sprintf(` AND customer_id = $%d`, argIdx)
		countArgs = append(countArgs, customerID)
		argIdx++
	}

	// Get total count
	var total int64
	err := r.db.WithContext(ctx).Raw(countQuery, countArgs...).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count claims: %w", err)
	}

	// Build list query
	query := `
		SELECT claim_id, claim_number, policy_id, customer_id, status, type,
		       claimed_amount, approved_amount, settled_amount,
		       claimed_currency, approved_currency, settled_currency,
		       incident_date, incident_description, submitted_at, approved_at, settled_at,
		       COALESCE(rejection_reason, '') as rejection_reason,
		       COALESCE(place_of_incident, '') as place_of_incident,
		       COALESCE(bank_details_for_payout, '') as bank_details_for_payout,
		       COALESCE(appeal_option_available, false) as appeal_option_available,
		       in_app_messages, processing_type,
		       deductible_amount, co_pay_amount,
		       COALESCE(processor_notes, '') as processor_notes,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.claims
		WHERE deleted_at IS NULL`
	
	var queryArgs []interface{}
	argIdx = 1
	
	if policyID != "" {
		query += fmt.Sprintf(` AND policy_id = $%d`, argIdx)
		queryArgs = append(queryArgs, policyID)
		argIdx++
	}
	if customerID != "" {
		query += fmt.Sprintf(` AND customer_id = $%d`, argIdx)
		queryArgs = append(queryArgs, customerID)
		argIdx++
	}
	
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query, queryArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list claims: %w", err)
	}
	defer rows.Close()

	claims := make([]*claimsv1.Claim, 0)
	for rows.Next() {
		var (
			c                  claimsv1.Claim
			statusStr          sql.NullString
			typeStr            sql.NullString
			processingTypeStr  sql.NullString
			claimedAmount      int64
			approvedAmount     sql.NullInt64
			settledAmount      sql.NullInt64
			claimedCurrency    string
			approvedCurrency   string
			settledCurrency    string
			deductibleAmount   sql.NullInt64
			coPayAmount        sql.NullInt64
			incidentDate       time.Time
			submittedAt        time.Time
			approvedAt         sql.NullTime
			settledAt          sql.NullTime
			createdAt          time.Time
			updatedAt          time.Time
			deletedAt          sql.NullTime
			inAppMessages      sql.NullString
		)

		err := rows.Scan(
			&c.ClaimId,
			&c.ClaimNumber,
			&c.PolicyId,
			&c.CustomerId,
			&statusStr,
			&typeStr,
			&claimedAmount,
			&approvedAmount,
			&settledAmount,
			&claimedCurrency,
			&approvedCurrency,
			&settledCurrency,
			&incidentDate,
			&c.IncidentDescription,
			&submittedAt,
			&approvedAt,
			&settledAt,
			&c.RejectionReason,
			&c.PlaceOfIncident,
			&c.BankDetailsForPayout,
			&c.AppealOptionAvailable,
			&inAppMessages,
			&processingTypeStr,
			&deductibleAmount,
			&coPayAmount,
			&c.ProcessorNotes,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan claim: %w", err)
		}

		// Set Money fields
		c.ClaimedAmount = &commonv1.Money{
			Amount:   claimedAmount,
			Currency: claimedCurrency,
		}
		
		// Set currency companion fields
		c.ClaimedCurrency = claimedCurrency
		c.ApprovedCurrency = approvedCurrency
		c.SettledCurrency = settledCurrency
		
		if approvedAmount.Valid {
			c.ApprovedAmount = &commonv1.Money{
				Amount:   approvedAmount.Int64,
				Currency: approvedCurrency,
			}
		}
		if settledAmount.Valid {
			c.SettledAmount = &commonv1.Money{
				Amount:   settledAmount.Int64,
				Currency: settledCurrency,
			}
		}
		if deductibleAmount.Valid {
			c.DeductibleAmount = &commonv1.Money{
				Amount:   deductibleAmount.Int64,
				Currency: "BDT",
			}
		}
		if coPayAmount.Valid {
			c.CoPayAmount = &commonv1.Money{
				Amount:   coPayAmount.Int64,
				Currency: "BDT",
			}
		}
		
		// Set in_app_messages
		if inAppMessages.Valid {
			c.InAppMessages = inAppMessages.String
		}

		// Parse enums
		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := claimsv1.ClaimStatus_value[k]; ok {
				c.Status = claimsv1.ClaimStatus(v)
			}
		}
		if typeStr.Valid {
			k := strings.ToUpper(typeStr.String)
			if v, ok := claimsv1.ClaimType_value[k]; ok {
				c.Type = claimsv1.ClaimType(v)
			}
		}
		if processingTypeStr.Valid {
			k := strings.ToUpper(processingTypeStr.String)
			if v, ok := claimsv1.ClaimProcessingType_value[k]; ok {
				c.ProcessingType = claimsv1.ClaimProcessingType(v)
			}
		}

		if !incidentDate.IsZero() {
			c.IncidentDate = timestamppb.New(incidentDate)
		}
		if !submittedAt.IsZero() {
			c.SubmittedAt = timestamppb.New(submittedAt)
		}
		if approvedAt.Valid {
			c.ApprovedAt = timestamppb.New(approvedAt.Time)
		}
		if settledAt.Valid {
			c.SettledAt = timestamppb.New(settledAt.Time)
		}
		if !createdAt.IsZero() {
			c.CreatedAt = timestamppb.New(createdAt)
		}
		if !updatedAt.IsZero() {
			c.UpdatedAt = timestamppb.New(updatedAt)
		}
		if deletedAt.Valid {
			c.DeletedAt = timestamppb.New(deletedAt.Time)
		}

		claims = append(claims, &c)
	}

	return claims, total, nil
}

func (r *ClaimRepository) Delete(ctx context.Context, claimID string) error {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.claims
		SET deleted_at = NOW()
		WHERE claim_id = $1 AND deleted_at IS NULL`,
		claimID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete claim: %w", err)
	}

	return nil
}
