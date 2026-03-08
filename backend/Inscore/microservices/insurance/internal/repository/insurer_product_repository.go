package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/lib/pq"

	insurerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/insurer/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type InsurerProductRepository struct {
	db *gorm.DB
}

func NewInsurerProductRepository(db *gorm.DB) *InsurerProductRepository {
	return &InsurerProductRepository{db: db}
}

func (r *InsurerProductRepository) Create(ctx context.Context, product *insurerv1.InsurerProduct) (*insurerv1.InsurerProduct, error) {
	if product.Id == "" {
		return nil, fmt.Errorf("product_id is required")
	}

	// Extract Money values
	minSumAssured := int64(0)
	minSumAssuredCurrency := "BDT"
	if product.MinSumAssured != nil {
		minSumAssured = product.MinSumAssured.Amount
		minSumAssuredCurrency = product.MinSumAssured.Currency
	}

	maxSumAssured := int64(0)
	maxSumAssuredCurrency := "BDT"
	if product.MaxSumAssured != nil {
		maxSumAssured = product.MaxSumAssured.Amount
		maxSumAssuredCurrency = product.MaxSumAssured.Currency
	}

	minPremium := int64(0)
	minPremiumCurrency := "BDT"
	if product.MinPremium != nil {
		minPremium = product.MinPremium.Amount
		minPremiumCurrency = product.MinPremium.Currency
	}

	maxPremium := int64(0)
	maxPremiumCurrency := "BDT"
	if product.MaxPremium != nil {
		maxPremium = product.MaxPremium.Amount
		maxPremiumCurrency = product.MaxPremium.Currency
	}

	medicalThreshold := int64(0)
	medicalThresholdCurrency := "BDT"
	if product.MedicalThreshold != nil {
		medicalThreshold = product.MedicalThreshold.Amount
		medicalThresholdCurrency = product.MedicalThreshold.Currency
	}

	// Handle timestamps
	var effectiveFrom time.Time
	if product.EffectiveFrom != nil {
		effectiveFrom = product.EffectiveFrom.AsTime()
	}

	var effectiveTo sql.NullTime
	if product.EffectiveTo != nil {
		effectiveTo = sql.NullTime{Time: product.EffectiveTo.AsTime(), Valid: true}
	}

	// Handle JSONB fields
	var features interface{}
	if product.Features != "" {
		features = product.Features
	}
	var exclusions interface{}
	if product.Exclusions != "" {
		exclusions = product.Exclusions
	}
	var auditInfo interface{}
	if product.AuditInfo != nil {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.insurer_products
			(product_id, insurer_id, base_product_id, code, name, status,
			 min_sum_assured, min_sum_assured_currency, max_sum_assured, max_sum_assured_currency,
			 min_premium, min_premium_currency, max_premium, max_premium_currency,
			 min_entry_age, max_entry_age, max_maturity_age, min_term_years, max_term_years,
			 premium_payment_modes, medical_required, medical_threshold, medical_threshold_currency,
			 free_look_period_days, commission_config_id, features, exclusions,
			 effective_from, effective_to, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)`,
		product.Id,
		product.InsurerId,
		product.ProductId,
		product.Code,
		product.Name,
		strings.ToUpper(product.Status.String()),
		minSumAssured,
		minSumAssuredCurrency,
		maxSumAssured,
		maxSumAssuredCurrency,
		minPremium,
		minPremiumCurrency,
		maxPremium,
		maxPremiumCurrency,
		product.MinEntryAge,
		product.MaxEntryAge,
		product.MaxMaturityAge,
		product.MinTermYears,
		product.MaxTermYears,
		pq.Array(product.PremiumPaymentModes),
		product.MedicalRequired,
		medicalThreshold,
		medicalThresholdCurrency,
		product.FreeLookPeriodDays,
		product.CommissionConfigId,
		features,
		exclusions,
		effectiveFrom,
		effectiveTo,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert insurer product: %w", err)
	}

	return r.GetByID(ctx, product.Id)
}

func (r *InsurerProductRepository) GetByID(ctx context.Context, productID string) (*insurerv1.InsurerProduct, error) {
	var (
		p                        insurerv1.InsurerProduct
		statusStr                sql.NullString
		minSumAssured            int64
		minSumAssuredCurrency    string
		maxSumAssured            int64
		maxSumAssuredCurrency    string
		minPremium               int64
		minPremiumCurrency       string
		maxPremium               int64
		maxPremiumCurrency       string
		premiumPaymentModes      pq.StringArray
		medicalThreshold         int64
		medicalThresholdCurrency string
		commissionConfigID       sql.NullString
		features                 sql.NullString
		exclusions               sql.NullString
		effectiveFrom            time.Time
		effectiveTo              sql.NullTime
		auditInfo                sql.NullString
		deletedAt                sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT product_id, insurer_id, base_product_id, code, name, status,
		       min_sum_assured, min_sum_assured_currency, max_sum_assured, max_sum_assured_currency,
		       min_premium, min_premium_currency, max_premium, max_premium_currency,
		       min_entry_age, max_entry_age, max_maturity_age, min_term_years, max_term_years,
		       premium_payment_modes, medical_required, medical_threshold, medical_threshold_currency,
		       free_look_period_days, commission_config_id, features, exclusions,
		       effective_from, effective_to, audit_info, deleted_at
		FROM insurance_schema.insurer_products
		WHERE product_id = $1 AND deleted_at IS NULL`,
		productID,
	).Row().Scan(
		&p.Id,
		&p.InsurerId,
		&p.ProductId,
		&p.Code,
		&p.Name,
		&statusStr,
		&minSumAssured,
		&minSumAssuredCurrency,
		&maxSumAssured,
		&maxSumAssuredCurrency,
		&minPremium,
		&minPremiumCurrency,
		&maxPremium,
		&maxPremiumCurrency,
		&p.MinEntryAge,
		&p.MaxEntryAge,
		&p.MaxMaturityAge,
		&p.MinTermYears,
		&p.MaxTermYears,
		&premiumPaymentModes,
		&p.MedicalRequired,
		&medicalThreshold,
		&medicalThresholdCurrency,
		&p.FreeLookPeriodDays,
		&commissionConfigID,
		&features,
		&exclusions,
		&effectiveFrom,
		&effectiveTo,
		&auditInfo,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get insurer product: %w", err)
	}

	// Set Money fields
	p.MinSumAssured = &commonv1.Money{Amount: minSumAssured, Currency: minSumAssuredCurrency}
	p.MaxSumAssured = &commonv1.Money{Amount: maxSumAssured, Currency: maxSumAssuredCurrency}
	p.MinPremium = &commonv1.Money{Amount: minPremium, Currency: minPremiumCurrency}
	p.MaxPremium = &commonv1.Money{Amount: maxPremium, Currency: maxPremiumCurrency}
	p.MedicalThreshold = &commonv1.Money{Amount: medicalThreshold, Currency: medicalThresholdCurrency}

	// Set array field
	p.PremiumPaymentModes = premiumPaymentModes

	// Set optional fields
	if commissionConfigID.Valid {
		p.CommissionConfigId = commissionConfigID.String
	}
	if features.Valid {
		p.Features = features.String
	}
	if exclusions.Valid {
		p.Exclusions = exclusions.String
	}

	// Parse enum
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := insurerv1.ProductStatus_value[k]; ok {
			p.Status = insurerv1.ProductStatus(v)
		}
	}

	// Set timestamps
	if !effectiveFrom.IsZero() {
		p.EffectiveFrom = timestamppb.New(effectiveFrom)
	}
	if effectiveTo.Valid {
		p.EffectiveTo = timestamppb.New(effectiveTo.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		p.AuditInfo = &commonv1.AuditInfo{}
	}

	return &p, nil
}

func (r *InsurerProductRepository) Update(ctx context.Context, product *insurerv1.InsurerProduct) (*insurerv1.InsurerProduct, error) {
	// Extract Money values
	minSumAssured := int64(0)
	minSumAssuredCurrency := "BDT"
	if product.MinSumAssured != nil {
		minSumAssured = product.MinSumAssured.Amount
		minSumAssuredCurrency = product.MinSumAssured.Currency
	}

	maxSumAssured := int64(0)
	maxSumAssuredCurrency := "BDT"
	if product.MaxSumAssured != nil {
		maxSumAssured = product.MaxSumAssured.Amount
		maxSumAssuredCurrency = product.MaxSumAssured.Currency
	}

	minPremium := int64(0)
	minPremiumCurrency := "BDT"
	if product.MinPremium != nil {
		minPremium = product.MinPremium.Amount
		minPremiumCurrency = product.MinPremium.Currency
	}

	maxPremium := int64(0)
	maxPremiumCurrency := "BDT"
	if product.MaxPremium != nil {
		maxPremium = product.MaxPremium.Amount
		maxPremiumCurrency = product.MaxPremium.Currency
	}

	medicalThreshold := int64(0)
	medicalThresholdCurrency := "BDT"
	if product.MedicalThreshold != nil {
		medicalThreshold = product.MedicalThreshold.Amount
		medicalThresholdCurrency = product.MedicalThreshold.Currency
	}

	// Handle timestamps
	var effectiveFrom time.Time
	if product.EffectiveFrom != nil {
		effectiveFrom = product.EffectiveFrom.AsTime()
	}

	var effectiveTo sql.NullTime
	if product.EffectiveTo != nil {
		effectiveTo = sql.NullTime{Time: product.EffectiveTo.AsTime(), Valid: true}
	}

	// Handle JSONB fields
	var features interface{}
	if product.Features != "" {
		features = product.Features
	}
	var exclusions interface{}
	if product.Exclusions != "" {
		exclusions = product.Exclusions
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.insurer_products
		SET insurer_id = $2,
		    base_product_id = $3,
		    code = $4,
		    name = $5,
		    status = $6,
		    min_sum_assured = $7,
		    min_sum_assured_currency = $8,
		    max_sum_assured = $9,
		    max_sum_assured_currency = $10,
		    min_premium = $11,
		    min_premium_currency = $12,
		    max_premium = $13,
		    max_premium_currency = $14,
		    min_entry_age = $15,
		    max_entry_age = $16,
		    max_maturity_age = $17,
		    min_term_years = $18,
		    max_term_years = $19,
		    premium_payment_modes = $20,
		    medical_required = $21,
		    medical_threshold = $22,
		    medical_threshold_currency = $23,
		    free_look_period_days = $24,
		    commission_config_id = $25,
		    features = $26,
		    exclusions = $27,
		    effective_from = $28,
		    effective_to = $29
		WHERE product_id = $1 AND deleted_at IS NULL`,
		product.Id,
		product.InsurerId,
		product.ProductId,
		product.Code,
		product.Name,
		strings.ToUpper(product.Status.String()),
		minSumAssured,
		minSumAssuredCurrency,
		maxSumAssured,
		maxSumAssuredCurrency,
		minPremium,
		minPremiumCurrency,
		maxPremium,
		maxPremiumCurrency,
		product.MinEntryAge,
		product.MaxEntryAge,
		product.MaxMaturityAge,
		product.MinTermYears,
		product.MaxTermYears,
		pq.Array(product.PremiumPaymentModes),
		product.MedicalRequired,
		medicalThreshold,
		medicalThresholdCurrency,
		product.FreeLookPeriodDays,
		product.CommissionConfigId,
		features,
		exclusions,
		effectiveFrom,
		effectiveTo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update insurer product: %w", err)
	}

	return r.GetByID(ctx, product.Id)
}

func (r *InsurerProductRepository) Delete(ctx context.Context, productID string) error {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.insurer_products
		SET deleted_at = NOW()
		WHERE product_id = $1 AND deleted_at IS NULL`,
		productID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete insurer product: %w", err)
	}

	return nil
}

func (r *InsurerProductRepository) ListByInsurerID(ctx context.Context, insurerID string) ([]*insurerv1.InsurerProduct, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT product_id, insurer_id, base_product_id, code, name, status,
		       min_sum_assured, min_sum_assured_currency, max_sum_assured, max_sum_assured_currency,
		       min_premium, min_premium_currency, max_premium, max_premium_currency,
		       min_entry_age, max_entry_age, max_maturity_age, min_term_years, max_term_years,
		       premium_payment_modes, medical_required, medical_threshold, medical_threshold_currency,
		       free_look_period_days, commission_config_id, features, exclusions,
		       effective_from, effective_to, audit_info, deleted_at
		FROM insurance_schema.insurer_products
		WHERE insurer_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC`,
		insurerID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list insurer products: %w", err)
	}
	defer rows.Close()

	products := make([]*insurerv1.InsurerProduct, 0)
	for rows.Next() {
		var (
			p                        insurerv1.InsurerProduct
			statusStr                sql.NullString
			minSumAssured            int64
			minSumAssuredCurrency    string
			maxSumAssured            int64
			maxSumAssuredCurrency    string
			minPremium               int64
			minPremiumCurrency       string
			maxPremium               int64
			maxPremiumCurrency       string
			premiumPaymentModes      pq.StringArray
			medicalThreshold         int64
			medicalThresholdCurrency string
			commissionConfigID       sql.NullString
			features                 sql.NullString
			exclusions               sql.NullString
			effectiveFrom            time.Time
			effectiveTo              sql.NullTime
			auditInfo                sql.NullString
			deletedAt                sql.NullTime
		)

		err := rows.Scan(
			&p.Id,
			&p.InsurerId,
			&p.ProductId,
			&p.Code,
			&p.Name,
			&statusStr,
			&minSumAssured,
			&minSumAssuredCurrency,
			&maxSumAssured,
			&maxSumAssuredCurrency,
			&minPremium,
			&minPremiumCurrency,
			&maxPremium,
			&maxPremiumCurrency,
			&p.MinEntryAge,
			&p.MaxEntryAge,
			&p.MaxMaturityAge,
			&p.MinTermYears,
			&p.MaxTermYears,
			&premiumPaymentModes,
			&p.MedicalRequired,
			&medicalThreshold,
			&medicalThresholdCurrency,
			&p.FreeLookPeriodDays,
			&commissionConfigID,
			&features,
			&exclusions,
			&effectiveFrom,
			&effectiveTo,
			&auditInfo,
			&deletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan insurer product: %w", err)
		}

		// Set Money fields
		p.MinSumAssured = &commonv1.Money{Amount: minSumAssured, Currency: minSumAssuredCurrency}
		p.MaxSumAssured = &commonv1.Money{Amount: maxSumAssured, Currency: maxSumAssuredCurrency}
		p.MinPremium = &commonv1.Money{Amount: minPremium, Currency: minPremiumCurrency}
		p.MaxPremium = &commonv1.Money{Amount: maxPremium, Currency: maxPremiumCurrency}
		p.MedicalThreshold = &commonv1.Money{Amount: medicalThreshold, Currency: medicalThresholdCurrency}

		// Set array field
		p.PremiumPaymentModes = premiumPaymentModes

		// Set optional fields
		if commissionConfigID.Valid {
			p.CommissionConfigId = commissionConfigID.String
		}
		if features.Valid {
			p.Features = features.String
		}
		if exclusions.Valid {
			p.Exclusions = exclusions.String
		}

		// Parse enum
		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := insurerv1.ProductStatus_value[k]; ok {
				p.Status = insurerv1.ProductStatus(v)
			}
		}

		// Set timestamps
		if !effectiveFrom.IsZero() {
			p.EffectiveFrom = timestamppb.New(effectiveFrom)
		}
		if effectiveTo.Valid {
			p.EffectiveTo = timestamppb.New(effectiveTo.Time)
		}

		// Set audit info
		if auditInfo.Valid {
			p.AuditInfo = &commonv1.AuditInfo{}
		}

		products = append(products, &p)
	}

	return products, nil
}
