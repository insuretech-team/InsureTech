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

type PolicyRepository struct {
	db *gorm.DB
}

func NewPolicyRepository(db *gorm.DB) *PolicyRepository {
	return &PolicyRepository{db: db}
}

func (r *PolicyRepository) Create(ctx context.Context, policy *policyv1.Policy) (*policyv1.Policy, error) {
	if policy.PolicyId == "" {
		return nil, fmt.Errorf("policy_id is required")
	}
	
	// Extract Money values
	premiumAmount := int64(0)
	premiumCurrency := "BDT"
	if policy.PremiumAmount != nil {
		premiumAmount = policy.PremiumAmount.Amount
		premiumCurrency = policy.PremiumAmount.Currency
	}
	
	sumInsured := int64(0)
	sumInsuredCurrency := "BDT"
	if policy.SumInsured != nil {
		sumInsured = policy.SumInsured.Amount
		sumInsuredCurrency = policy.SumInsured.Currency
	}
	
	vatTax := int64(0)
	if policy.VatTax != nil {
		vatTax = policy.VatTax.Amount
	}
	
	serviceFee := int64(0)
	if policy.ServiceFee != nil {
		serviceFee = policy.ServiceFee.Amount
	}
	
	totalPayable := int64(0)
	if policy.TotalPayable != nil {
		totalPayable = policy.TotalPayable.Amount
	}
	
	// Handle timestamps
	var startDate, endDate, issuedAt interface{}
	if policy.StartDate != nil {
		startDate = policy.StartDate.AsTime()
	}
	if policy.EndDate != nil {
		endDate = policy.EndDate.AsTime()
	}
	if policy.IssuedAt != nil {
		issuedAt = policy.IssuedAt.AsTime()
	}
	
	// Handle JSONB underwriting_data
	var underwritingData interface{}
	if policy.UnderwritingData != "" {
		underwritingData = policy.UnderwritingData
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.policies
			(policy_id, policy_number, product_id, customer_id, partner_id, agent_id,
			 quote_id, underwriting_decision_id, status,
			 premium_amount, sum_insured, tenure_months,
			 start_date, end_date, issued_at,
			 premium_currency, sum_insured_currency,
			 vat_tax, service_fee, total_payable,
			 payment_frequency, payment_gateway_reference, receipt_number,
			 occupation_risk_class, has_existing_policies, claims_history_summary,
			 provider_name, underwriting_data,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, NOW(), NOW())`,
		policy.PolicyId,
		policy.PolicyNumber,
		policy.ProductId,
		policy.CustomerId,
		policy.PartnerId,
		policy.AgentId,
		policy.QuoteId,
		policy.UnderwritingDecisionId,
		strings.ToUpper(policy.Status.String()),
		premiumAmount,
		sumInsured,
		policy.TenureMonths,
		startDate,
		endDate,
		issuedAt,
		premiumCurrency,
		sumInsuredCurrency,
		vatTax,
		serviceFee,
		totalPayable,
		policy.PaymentFrequency,
		policy.PaymentGatewayReference,
		policy.ReceiptNumber,
		policy.OccupationRiskClass,
		policy.HasExistingPolicies,
		policy.ClaimsHistorySummary,
		policy.ProviderName,
		underwritingData,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert policy: %w", err)
	}

	return r.GetByID(ctx, policy.PolicyId)
}

func (r *PolicyRepository) GetByID(ctx context.Context, policyID string) (*policyv1.Policy, error) {
	var (
		p                  policyv1.Policy
		statusStr          sql.NullString
		premiumAmount      int64
		sumInsured         int64
		premiumCurrency    string
		sumInsuredCurrency string
		vatTax             sql.NullInt64
		serviceFee         sql.NullInt64
		totalPayable       sql.NullInt64
		startDate          time.Time
		endDate            time.Time
		issuedAt           sql.NullTime
		createdAt          time.Time
		updatedAt          time.Time
		deletedAt          sql.NullTime
		underwritingData   sql.NullString
		partnerId          sql.NullString
		agentId            sql.NullString
		quoteId            sql.NullString
		underwritingDecId  sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT policy_id, policy_number, product_id, customer_id, 
		       partner_id, agent_id, quote_id, underwriting_decision_id,
		       status, premium_amount, sum_insured, tenure_months,
		       start_date, end_date, issued_at,
		       premium_currency, sum_insured_currency,
		       vat_tax, service_fee, total_payable,
		       COALESCE(payment_frequency, '') as payment_frequency,
		       COALESCE(payment_gateway_reference, '') as payment_gateway_reference,
		       COALESCE(receipt_number, '') as receipt_number,
		       COALESCE(occupation_risk_class, '') as occupation_risk_class,
		       COALESCE(has_existing_policies, false) as has_existing_policies,
		       COALESCE(claims_history_summary, '') as claims_history_summary,
		       COALESCE(provider_name, '') as provider_name,
		       COALESCE(policy_document_url, '') as policy_document_url,
		       underwriting_data,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.policies
		WHERE policy_id = $1 AND deleted_at IS NULL`,
		policyID,
	).Row().Scan(
		&p.PolicyId,
		&p.PolicyNumber,
		&p.ProductId,
		&p.CustomerId,
		&partnerId,
		&agentId,
		&quoteId,
		&underwritingDecId,
		&statusStr,
		&premiumAmount,
		&sumInsured,
		&p.TenureMonths,
		&startDate,
		&endDate,
		&issuedAt,
		&premiumCurrency,
		&sumInsuredCurrency,
		&vatTax,
		&serviceFee,
		&totalPayable,
		&p.PaymentFrequency,
		&p.PaymentGatewayReference,
		&p.ReceiptNumber,
		&p.OccupationRiskClass,
		&p.HasExistingPolicies,
		&p.ClaimsHistorySummary,
		&p.ProviderName,
		&p.PolicyDocumentUrl,
		&underwritingData,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	// Set Money fields
	p.PremiumAmount = &commonv1.Money{
		Amount:   premiumAmount,
		Currency: premiumCurrency,
	}
	p.SumInsured = &commonv1.Money{
		Amount:   sumInsured,
		Currency: sumInsuredCurrency,
	}
	
	// Set currency companion fields
	p.PremiumCurrency = premiumCurrency
	p.SumInsuredCurrency = sumInsuredCurrency
	
	if vatTax.Valid {
		p.VatTax = &commonv1.Money{
			Amount:   vatTax.Int64,
			Currency: "BDT",
		}
	}
	if serviceFee.Valid {
		p.ServiceFee = &commonv1.Money{
			Amount:   serviceFee.Int64,
			Currency: "BDT",
		}
	}
	if totalPayable.Valid {
		p.TotalPayable = &commonv1.Money{
			Amount:   totalPayable.Int64,
			Currency: "BDT",
		}
	}
	
	// Set nullable foreign keys
	if partnerId.Valid {
		p.PartnerId = partnerId.String
	}
	if agentId.Valid {
		p.AgentId = agentId.String
	}
	if quoteId.Valid {
		p.QuoteId = quoteId.String
	}
	if underwritingDecId.Valid {
		p.UnderwritingDecisionId = underwritingDecId.String
	}
	
	// Set underwriting_data
	if underwritingData.Valid {
		p.UnderwritingData = underwritingData.String
	}

	// Parse status enum
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := policyv1.PolicyStatus_value[k]; ok {
			p.Status = policyv1.PolicyStatus(v)
		}
	}

	if !startDate.IsZero() {
		p.StartDate = timestamppb.New(startDate)
	}
	if !endDate.IsZero() {
		p.EndDate = timestamppb.New(endDate)
	}
	if issuedAt.Valid {
		p.IssuedAt = timestamppb.New(issuedAt.Time)
	}
	if !createdAt.IsZero() {
		p.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		p.UpdatedAt = timestamppb.New(updatedAt)
	}
	if deletedAt.Valid {
		p.DeletedAt = timestamppb.New(deletedAt.Time)
	}

	return &p, nil
}

func (r *PolicyRepository) Update(ctx context.Context, policy *policyv1.Policy) (*policyv1.Policy, error) {
	// Extract Money values
	premiumAmount := int64(0)
	premiumCurrency := "BDT"
	if policy.PremiumAmount != nil {
		premiumAmount = policy.PremiumAmount.Amount
		premiumCurrency = policy.PremiumAmount.Currency
	}
	
	sumInsured := int64(0)
	sumInsuredCurrency := "BDT"
	if policy.SumInsured != nil {
		sumInsured = policy.SumInsured.Amount
		sumInsuredCurrency = policy.SumInsured.Currency
	}
	
	vatTax := int64(0)
	if policy.VatTax != nil {
		vatTax = policy.VatTax.Amount
	}
	
	serviceFee := int64(0)
	if policy.ServiceFee != nil {
		serviceFee = policy.ServiceFee.Amount
	}
	
	totalPayable := int64(0)
	if policy.TotalPayable != nil {
		totalPayable = policy.TotalPayable.Amount
	}
	
	// Handle timestamps
	var startDate, endDate, issuedAt interface{}
	if policy.StartDate != nil {
		startDate = policy.StartDate.AsTime()
	}
	if policy.EndDate != nil {
		endDate = policy.EndDate.AsTime()
	}
	if policy.IssuedAt != nil {
		issuedAt = policy.IssuedAt.AsTime()
	}
	
	// Handle JSONB underwriting_data
	var underwritingData interface{}
	if policy.UnderwritingData != "" {
		underwritingData = policy.UnderwritingData
	}
	
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.policies
		SET policy_number = $2,
		    product_id = $3,
		    customer_id = $4,
		    partner_id = $5,
		    agent_id = $6,
		    quote_id = $7,
		    underwriting_decision_id = $8,
		    status = $9,
		    premium_amount = $10,
		    sum_insured = $11,
		    tenure_months = $12,
		    start_date = $13,
		    end_date = $14,
		    issued_at = $15,
		    premium_currency = $16,
		    sum_insured_currency = $17,
		    vat_tax = $18,
		    service_fee = $19,
		    total_payable = $20,
		    payment_frequency = $21,
		    payment_gateway_reference = $22,
		    receipt_number = $23,
		    occupation_risk_class = $24,
		    has_existing_policies = $25,
		    claims_history_summary = $26,
		    provider_name = $27,
		    underwriting_data = $28,
		    updated_at = NOW()
		WHERE policy_id = $1 AND deleted_at IS NULL`,
		policy.PolicyId,
		policy.PolicyNumber,
		policy.ProductId,
		policy.CustomerId,
		policy.PartnerId,
		policy.AgentId,
		policy.QuoteId,
		policy.UnderwritingDecisionId,
		strings.ToUpper(policy.Status.String()),
		premiumAmount,
		sumInsured,
		policy.TenureMonths,
		startDate,
		endDate,
		issuedAt,
		premiumCurrency,
		sumInsuredCurrency,
		vatTax,
		serviceFee,
		totalPayable,
		policy.PaymentFrequency,
		policy.PaymentGatewayReference,
		policy.ReceiptNumber,
		policy.OccupationRiskClass,
		policy.HasExistingPolicies,
		policy.ClaimsHistorySummary,
		policy.ProviderName,
		underwritingData,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update policy: %w", err)
	}

	return r.GetByID(ctx, policy.PolicyId)
}

func (r *PolicyRepository) List(ctx context.Context, tenantID, customerID string, page, pageSize int) ([]*policyv1.Policy, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// Build count query
	countQuery := `SELECT COUNT(*) FROM insurance_schema.policies WHERE deleted_at IS NULL`
	var countArgs []interface{}
	argIdx := 1
	
	if tenantID != "" {
		countQuery += fmt.Sprintf(` AND tenant_id = $%d`, argIdx)
		countArgs = append(countArgs, tenantID)
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
		return nil, 0, fmt.Errorf("failed to count policies: %w", err)
	}

	// Build list query
	query := `
		SELECT policy_id, policy_number, product_id, customer_id, 
		       partner_id, agent_id, quote_id, underwriting_decision_id,
		       status, premium_amount, sum_insured, tenure_months,
		       start_date, end_date, issued_at,
		       premium_currency, sum_insured_currency,
		       vat_tax, service_fee, total_payable,
		       COALESCE(payment_frequency, '') as payment_frequency,
		       COALESCE(payment_gateway_reference, '') as payment_gateway_reference,
		       COALESCE(receipt_number, '') as receipt_number,
		       COALESCE(occupation_risk_class, '') as occupation_risk_class,
		       COALESCE(has_existing_policies, false) as has_existing_policies,
		       COALESCE(claims_history_summary, '') as claims_history_summary,
		       COALESCE(provider_name, '') as provider_name,
		       COALESCE(policy_document_url, '') as policy_document_url,
		       underwriting_data,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.policies
		WHERE deleted_at IS NULL`
	
	var queryArgs []interface{}
	argIdx = 1
	
	if tenantID != "" {
		query += fmt.Sprintf(` AND tenant_id = $%d`, argIdx)
		queryArgs = append(queryArgs, tenantID)
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
		return nil, 0, fmt.Errorf("failed to list policies: %w", err)
	}
	defer rows.Close()

	policies := make([]*policyv1.Policy, 0)
	for rows.Next() {
		var (
			p                  policyv1.Policy
			statusStr          sql.NullString
			premiumAmount      int64
			sumInsured         int64
			premiumCurrency    string
			sumInsuredCurrency string
			vatTax             sql.NullInt64
			serviceFee         sql.NullInt64
			totalPayable       sql.NullInt64
			startDate          time.Time
			endDate            time.Time
			issuedAt           sql.NullTime
			createdAt          time.Time
			updatedAt          time.Time
			deletedAt          sql.NullTime
			underwritingData   sql.NullString
			partnerId          sql.NullString
			agentId            sql.NullString
			quoteId            sql.NullString
			underwritingDecId  sql.NullString
		)

		err := rows.Scan(
			&p.PolicyId,
			&p.PolicyNumber,
			&p.ProductId,
			&p.CustomerId,
			&partnerId,
			&agentId,
			&quoteId,
			&underwritingDecId,
			&statusStr,
			&premiumAmount,
			&sumInsured,
			&p.TenureMonths,
			&startDate,
			&endDate,
			&issuedAt,
			&premiumCurrency,
			&sumInsuredCurrency,
			&vatTax,
			&serviceFee,
			&totalPayable,
			&p.PaymentFrequency,
			&p.PaymentGatewayReference,
			&p.ReceiptNumber,
			&p.OccupationRiskClass,
			&p.HasExistingPolicies,
			&p.ClaimsHistorySummary,
			&p.ProviderName,
			&p.PolicyDocumentUrl,
			&underwritingData,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan policy: %w", err)
		}

		// Set Money fields
		p.PremiumAmount = &commonv1.Money{
			Amount:   premiumAmount,
			Currency: premiumCurrency,
		}
		p.SumInsured = &commonv1.Money{
			Amount:   sumInsured,
			Currency: sumInsuredCurrency,
		}
		
		// Set currency companion fields
		p.PremiumCurrency = premiumCurrency
		p.SumInsuredCurrency = sumInsuredCurrency
		
		if vatTax.Valid {
			p.VatTax = &commonv1.Money{
				Amount:   vatTax.Int64,
				Currency: "BDT",
			}
		}
		if serviceFee.Valid {
			p.ServiceFee = &commonv1.Money{
				Amount:   serviceFee.Int64,
				Currency: "BDT",
			}
		}
		if totalPayable.Valid {
			p.TotalPayable = &commonv1.Money{
				Amount:   totalPayable.Int64,
				Currency: "BDT",
			}
		}
		
		// Set nullable foreign keys
		if partnerId.Valid {
			p.PartnerId = partnerId.String
		}
		if agentId.Valid {
			p.AgentId = agentId.String
		}
		if quoteId.Valid {
			p.QuoteId = quoteId.String
		}
		if underwritingDecId.Valid {
			p.UnderwritingDecisionId = underwritingDecId.String
		}
		
		// Set underwriting_data
		if underwritingData.Valid {
			p.UnderwritingData = underwritingData.String
		}

		// Parse status enum
		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := policyv1.PolicyStatus_value[k]; ok {
				p.Status = policyv1.PolicyStatus(v)
			}
		}

		if !startDate.IsZero() {
			p.StartDate = timestamppb.New(startDate)
		}
		if !endDate.IsZero() {
			p.EndDate = timestamppb.New(endDate)
		}
		if issuedAt.Valid {
			p.IssuedAt = timestamppb.New(issuedAt.Time)
		}
		if !createdAt.IsZero() {
			p.CreatedAt = timestamppb.New(createdAt)
		}
		if !updatedAt.IsZero() {
			p.UpdatedAt = timestamppb.New(updatedAt)
		}
		if deletedAt.Valid {
			p.DeletedAt = timestamppb.New(deletedAt.Time)
		}

		policies = append(policies, &p)
	}

	return policies, total, nil
}

func (r *PolicyRepository) Delete(ctx context.Context, policyID string) error {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.policies
		SET deleted_at = NOW()
		WHERE policy_id = $1 AND deleted_at IS NULL`,
		policyID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	return nil
}
