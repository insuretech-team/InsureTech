package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	underwritingv1 "github.com/newage-saint/insuretech/gen/go/insuretech/underwriting/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type QuoteRepository struct {
	db *gorm.DB
}

func NewQuoteRepository(db *gorm.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

func (r *QuoteRepository) Create(ctx context.Context, quote *underwritingv1.Quote) (*underwritingv1.Quote, error) {
	if quote.Id == "" {
		return nil, fmt.Errorf("quote_id is required")
	}

	// Extract Money values
	sumAssured := int64(0)
	sumAssuredCurrency := "BDT"
	if quote.SumAssured != nil {
		sumAssured = quote.SumAssured.Amount
		sumAssuredCurrency = quote.SumAssured.Currency
	}

	basePremium := int64(0)
	basePremiumCurrency := "BDT"
	if quote.BasePremium != nil {
		basePremium = quote.BasePremium.Amount
		basePremiumCurrency = quote.BasePremium.Currency
	}

	riderPremium := int64(0)
	riderPremiumCurrency := "BDT"
	if quote.RiderPremium != nil {
		riderPremium = quote.RiderPremium.Amount
		riderPremiumCurrency = quote.RiderPremium.Currency
	}

	taxAmount := int64(0)
	taxAmountCurrency := "BDT"
	if quote.TaxAmount != nil {
		taxAmount = quote.TaxAmount.Amount
		taxAmountCurrency = quote.TaxAmount.Currency
	}

	totalPremium := int64(0)
	totalPremiumCurrency := "BDT"
	if quote.TotalPremium != nil {
		totalPremium = quote.TotalPremium.Amount
		totalPremiumCurrency = quote.TotalPremium.Currency
	}

	// Handle JSONB fields
	var premiumCalc interface{}
	if quote.PremiumCalculation != "" {
		premiumCalc = quote.PremiumCalculation
	}

	var selectedRiders interface{}
	if quote.SelectedRiders != "" {
		selectedRiders = quote.SelectedRiders
	}

	// Handle timestamps
	var validUntil time.Time
	if quote.ValidUntil != nil {
		validUntil = quote.ValidUntil.AsTime()
	}

	var convertedAt sql.NullTime
	if quote.ConvertedAt != nil {
		convertedAt = sql.NullTime{Time: quote.ConvertedAt.AsTime(), Valid: true}
	}

	var convertedPolicyID sql.NullString
	if quote.ConvertedPolicyId != "" {
		convertedPolicyID = sql.NullString{String: quote.ConvertedPolicyId, Valid: true}
	}

	var auditInfo interface{}
	if quote.AuditInfo != nil {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.quotes
			(quote_id, quote_number, beneficiary_id, insurer_product_id, status,
			 sum_assured, sum_assured_currency, term_years, premium_payment_mode,
			 base_premium, base_premium_currency, rider_premium, rider_premium_currency,
			 tax_amount, tax_amount_currency, total_premium, total_premium_currency,
			 premium_calculation, selected_riders, applicant_age, applicant_occupation,
			 smoker, valid_until, converted_policy_id, converted_at, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)`,
		quote.Id,
		quote.QuoteNumber,
		quote.BeneficiaryId,
		quote.InsurerProductId,
		strings.ToUpper(quote.Status.String()),
		sumAssured,
		sumAssuredCurrency,
		quote.TermYears,
		quote.PremiumPaymentMode,
		basePremium,
		basePremiumCurrency,
		riderPremium,
		riderPremiumCurrency,
		taxAmount,
		taxAmountCurrency,
		totalPremium,
		totalPremiumCurrency,
		premiumCalc,
		selectedRiders,
		quote.ApplicantAge,
		quote.ApplicantOccupation,
		quote.Smoker,
		validUntil,
		convertedPolicyID,
		convertedAt,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert quote: %w", err)
	}

	return r.GetByID(ctx, quote.Id)
}

func (r *QuoteRepository) GetByID(ctx context.Context, quoteID string) (*underwritingv1.Quote, error) {
	var (
		q                      underwritingv1.Quote
		statusStr              sql.NullString
		sumAssured             int64
		sumAssuredCurrency     string
		basePremium            int64
		basePremiumCurrency    string
		riderPremium           int64
		riderPremiumCurrency   string
		taxAmount              int64
		taxAmountCurrency      string
		totalPremium           int64
		totalPremiumCurrency   string
		premiumCalc            sql.NullString
		selectedRiders         sql.NullString
		applicantOccupation    sql.NullString
		validUntil             time.Time
		convertedPolicyID      sql.NullString
		convertedAt            sql.NullTime
		auditInfo              sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT quote_id, quote_number, beneficiary_id, insurer_product_id, status,
		       sum_assured, sum_assured_currency, term_years, premium_payment_mode,
		       base_premium, base_premium_currency, rider_premium, rider_premium_currency,
		       tax_amount, tax_amount_currency, total_premium, total_premium_currency,
		       premium_calculation, selected_riders, applicant_age, 
		       applicant_occupation, smoker, valid_until, converted_policy_id, 
		       converted_at, audit_info
		FROM insurance_schema.quotes
		WHERE quote_id = $1 AND deleted_at IS NULL`,
		quoteID,
	).Row().Scan(
		&q.Id,
		&q.QuoteNumber,
		&q.BeneficiaryId,
		&q.InsurerProductId,
		&statusStr,
		&sumAssured,
		&sumAssuredCurrency,
		&q.TermYears,
		&q.PremiumPaymentMode,
		&basePremium,
		&basePremiumCurrency,
		&riderPremium,
		&riderPremiumCurrency,
		&taxAmount,
		&taxAmountCurrency,
		&totalPremium,
		&totalPremiumCurrency,
		&premiumCalc,
		&selectedRiders,
		&q.ApplicantAge,
		&applicantOccupation,
		&q.Smoker,
		&validUntil,
		&convertedPolicyID,
		&convertedAt,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	// Set Money fields
	q.SumAssured = &commonv1.Money{Amount: sumAssured, Currency: sumAssuredCurrency}
	q.BasePremium = &commonv1.Money{Amount: basePremium, Currency: basePremiumCurrency}
	q.RiderPremium = &commonv1.Money{Amount: riderPremium, Currency: riderPremiumCurrency}
	q.TaxAmount = &commonv1.Money{Amount: taxAmount, Currency: taxAmountCurrency}
	q.TotalPremium = &commonv1.Money{Amount: totalPremium, Currency: totalPremiumCurrency}

	// Set JSONB fields
	if premiumCalc.Valid {
		q.PremiumCalculation = premiumCalc.String
	}
	if selectedRiders.Valid {
		q.SelectedRiders = selectedRiders.String
	}

	// Set optional fields
	if applicantOccupation.Valid {
		q.ApplicantOccupation = applicantOccupation.String
	}
	if convertedPolicyID.Valid {
		q.ConvertedPolicyId = convertedPolicyID.String
	}

	// Parse status enum
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := underwritingv1.QuoteStatus_value[k]; ok {
			q.Status = underwritingv1.QuoteStatus(v)
		}
	}

	// Set timestamps
	if !validUntil.IsZero() {
		q.ValidUntil = timestamppb.New(validUntil)
	}
	if convertedAt.Valid {
		q.ConvertedAt = timestamppb.New(convertedAt.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		q.AuditInfo = &commonv1.AuditInfo{}
	}

	return &q, nil
}

func (r *QuoteRepository) Update(ctx context.Context, quote *underwritingv1.Quote) (*underwritingv1.Quote, error) {
	// Extract Money values
	sumAssured := int64(0)
	sumAssuredCurrency := "BDT"
	if quote.SumAssured != nil {
		sumAssured = quote.SumAssured.Amount
		sumAssuredCurrency = quote.SumAssured.Currency
	}

	basePremium := int64(0)
	basePremiumCurrency := "BDT"
	if quote.BasePremium != nil {
		basePremium = quote.BasePremium.Amount
		basePremiumCurrency = quote.BasePremium.Currency
	}

	riderPremium := int64(0)
	riderPremiumCurrency := "BDT"
	if quote.RiderPremium != nil {
		riderPremium = quote.RiderPremium.Amount
		riderPremiumCurrency = quote.RiderPremium.Currency
	}

	taxAmount := int64(0)
	taxAmountCurrency := "BDT"
	if quote.TaxAmount != nil {
		taxAmount = quote.TaxAmount.Amount
		taxAmountCurrency = quote.TaxAmount.Currency
	}

	totalPremium := int64(0)
	totalPremiumCurrency := "BDT"
	if quote.TotalPremium != nil {
		totalPremium = quote.TotalPremium.Amount
		totalPremiumCurrency = quote.TotalPremium.Currency
	}

	// Handle JSONB fields
	var premiumCalc interface{}
	if quote.PremiumCalculation != "" {
		premiumCalc = quote.PremiumCalculation
	}

	var selectedRiders interface{}
	if quote.SelectedRiders != "" {
		selectedRiders = quote.SelectedRiders
	}

	// Handle timestamps
	var validUntil time.Time
	if quote.ValidUntil != nil {
		validUntil = quote.ValidUntil.AsTime()
	}

	var convertedAt sql.NullTime
	if quote.ConvertedAt != nil {
		convertedAt = sql.NullTime{Time: quote.ConvertedAt.AsTime(), Valid: true}
	}

	var convertedPolicyID sql.NullString
	if quote.ConvertedPolicyId != "" {
		convertedPolicyID = sql.NullString{String: quote.ConvertedPolicyId, Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.quotes
		SET quote_number = $2,
		    beneficiary_id = $3,
		    insurer_product_id = $4,
		    status = $5,
		    sum_assured = $6,
		    sum_assured_currency = $7,
		    term_years = $8,
		    premium_payment_mode = $9,
		    base_premium = $10,
		    base_premium_currency = $11,
		    rider_premium = $12,
		    rider_premium_currency = $13,
		    tax_amount = $14,
		    tax_amount_currency = $15,
		    total_premium = $16,
		    total_premium_currency = $17,
		    premium_calculation = $18,
		    selected_riders = $19,
		    applicant_age = $20,
		    applicant_occupation = $21,
		    smoker = $22,
		    valid_until = $23,
		    converted_policy_id = $24,
		    converted_at = $25
		WHERE quote_id = $1 AND deleted_at IS NULL`,
		quote.Id,
		quote.QuoteNumber,
		quote.BeneficiaryId,
		quote.InsurerProductId,
		strings.ToUpper(quote.Status.String()),
		sumAssured,
		sumAssuredCurrency,
		quote.TermYears,
		quote.PremiumPaymentMode,
		basePremium,
		basePremiumCurrency,
		riderPremium,
		riderPremiumCurrency,
		taxAmount,
		taxAmountCurrency,
		totalPremium,
		totalPremiumCurrency,
		premiumCalc,
		selectedRiders,
		quote.ApplicantAge,
		quote.ApplicantOccupation,
		quote.Smoker,
		validUntil,
		convertedPolicyID,
		convertedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update quote: %w", err)
	}

	return r.GetByID(ctx, quote.Id)
}

func (r *QuoteRepository) Delete(ctx context.Context, quoteID string) error {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.quotes
		SET deleted_at = NOW()
		WHERE quote_id = $1 AND deleted_at IS NULL`,
		quoteID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete quote: %w", err)
	}

	return nil
}

func (r *QuoteRepository) List(ctx context.Context, beneficiaryID string, page, pageSize int) ([]*underwritingv1.Quote, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM insurance_schema.quotes WHERE deleted_at IS NULL`
	if beneficiaryID != "" {
		countQuery += ` AND beneficiary_id = $1`
		err := r.db.WithContext(ctx).Raw(countQuery, beneficiaryID).Scan(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count quotes: %w", err)
		}
	} else {
		err := r.db.WithContext(ctx).Raw(countQuery).Scan(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count quotes: %w", err)
		}
	}

	// Get quotes
	query := `
		SELECT quote_id, quote_number, beneficiary_id, insurer_product_id, status,
		       sum_assured, sum_assured_currency, term_years, premium_payment_mode,
		       base_premium, base_premium_currency, rider_premium, rider_premium_currency,
		       tax_amount, tax_amount_currency, total_premium, total_premium_currency,
		       premium_calculation, selected_riders, applicant_age, 
		       applicant_occupation, smoker, valid_until, converted_policy_id, 
		       converted_at, audit_info
		FROM insurance_schema.quotes
		WHERE deleted_at IS NULL`

	if beneficiaryID != "" {
		query += ` AND beneficiary_id = $1`
		query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, pageSize, offset)
	} else {
		query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, pageSize, offset)
	}

	var rows *sql.Rows
	var err error
	if beneficiaryID != "" {
		rows, err = r.db.WithContext(ctx).Raw(query, beneficiaryID).Rows()
	} else {
		rows, err = r.db.WithContext(ctx).Raw(query).Rows()
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list quotes: %w", err)
	}
	defer rows.Close()

	quotes := make([]*underwritingv1.Quote, 0)
	for rows.Next() {
		var (
			q                      underwritingv1.Quote
			statusStr              sql.NullString
			sumAssured             int64
			sumAssuredCurrency     string
			basePremium            int64
			basePremiumCurrency    string
			riderPremium           int64
			riderPremiumCurrency   string
			taxAmount              int64
			taxAmountCurrency      string
			totalPremium           int64
			totalPremiumCurrency   string
			premiumCalc            sql.NullString
			selectedRiders         sql.NullString
			applicantOccupation    sql.NullString
			validUntil             time.Time
			convertedPolicyID      sql.NullString
			convertedAt            sql.NullTime
			auditInfo              sql.NullString
		)

		err := rows.Scan(
			&q.Id,
			&q.QuoteNumber,
			&q.BeneficiaryId,
			&q.InsurerProductId,
			&statusStr,
			&sumAssured,
			&sumAssuredCurrency,
			&q.TermYears,
			&q.PremiumPaymentMode,
			&basePremium,
			&basePremiumCurrency,
			&riderPremium,
			&riderPremiumCurrency,
			&taxAmount,
			&taxAmountCurrency,
			&totalPremium,
			&totalPremiumCurrency,
			&premiumCalc,
			&selectedRiders,
			&q.ApplicantAge,
			&applicantOccupation,
			&q.Smoker,
			&validUntil,
			&convertedPolicyID,
			&convertedAt,
			&auditInfo,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan quote: %w", err)
		}

		// Set Money fields
		q.SumAssured = &commonv1.Money{Amount: sumAssured, Currency: sumAssuredCurrency}
		q.BasePremium = &commonv1.Money{Amount: basePremium, Currency: basePremiumCurrency}
		q.RiderPremium = &commonv1.Money{Amount: riderPremium, Currency: riderPremiumCurrency}
		q.TaxAmount = &commonv1.Money{Amount: taxAmount, Currency: taxAmountCurrency}
		q.TotalPremium = &commonv1.Money{Amount: totalPremium, Currency: totalPremiumCurrency}

		// Set JSONB fields
		if premiumCalc.Valid {
			q.PremiumCalculation = premiumCalc.String
		}
		if selectedRiders.Valid {
			q.SelectedRiders = selectedRiders.String
		}

		// Set optional fields
		if applicantOccupation.Valid {
			q.ApplicantOccupation = applicantOccupation.String
		}
		if convertedPolicyID.Valid {
			q.ConvertedPolicyId = convertedPolicyID.String
		}

		// Parse status enum
		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := underwritingv1.QuoteStatus_value[k]; ok {
				q.Status = underwritingv1.QuoteStatus(v)
			}
		}

		// Set timestamps
		if !validUntil.IsZero() {
			q.ValidUntil = timestamppb.New(validUntil)
		}
		if convertedAt.Valid {
			q.ConvertedAt = timestamppb.New(convertedAt.Time)
		}

		// Set audit info
		if auditInfo.Valid {
			q.AuditInfo = &commonv1.AuditInfo{}
		}

		quotes = append(quotes, &q)
	}

	return quotes, total, nil
}
