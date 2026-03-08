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

type UnderwritingDecisionRepository struct {
	db *gorm.DB
}

func NewUnderwritingDecisionRepository(db *gorm.DB) *UnderwritingDecisionRepository {
	return &UnderwritingDecisionRepository{db: db}
}

func (r *UnderwritingDecisionRepository) Create(ctx context.Context, decision *underwritingv1.UnderwritingDecision) (*underwritingv1.UnderwritingDecision, error) {
	if decision.Id == "" {
		return nil, fmt.Errorf("decision_id is required")
	}

	// Extract Money values
	adjustedPremium := int64(0)
	adjustedPremiumCurrency := "BDT"
	if decision.AdjustedPremium != nil {
		adjustedPremium = decision.AdjustedPremium.Amount
		adjustedPremiumCurrency = decision.AdjustedPremium.Currency
	}

	// Handle JSONB fields
	var riskFactors interface{}
	if decision.RiskFactors != "" {
		riskFactors = decision.RiskFactors
	}

	var conditions interface{}
	if decision.Conditions != "" {
		conditions = decision.Conditions
	}

	// Handle timestamps
	var decidedAt time.Time
	if decision.DecidedAt != nil {
		decidedAt = decision.DecidedAt.AsTime()
	} else {
		decidedAt = time.Now()
	}

	var validUntil sql.NullTime
	if decision.ValidUntil != nil {
		validUntil = sql.NullTime{Time: decision.ValidUntil.AsTime(), Valid: true}
	}

	var underwriterID sql.NullString
	if decision.UnderwriterId != "" {
		underwriterID = sql.NullString{String: decision.UnderwriterId, Valid: true}
	}

	var auditInfo interface{}
	if decision.AuditInfo != nil {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.underwriting_decisions
			(decision_id, quote_id, decision, method, risk_score, risk_level,
			 risk_factors, reason, conditions, premium_adjusted, adjusted_premium,
			 adjusted_premium_currency, adjustment_reason, underwriter_id,
			 underwriter_comments, decided_at, valid_until, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		decision.Id,
		decision.QuoteId,
		strings.ToUpper(decision.Decision.String()),
		strings.ToUpper(decision.Method.String()),
		decision.RiskScore,
		strings.ToUpper(decision.RiskLevel.String()),
		riskFactors,
		decision.Reason,
		conditions,
		decision.PremiumAdjusted,
		adjustedPremium,
		adjustedPremiumCurrency,
		decision.AdjustmentReason,
		underwriterID,
		decision.UnderwriterComments,
		decidedAt,
		validUntil,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert underwriting decision: %w", err)
	}

	return r.GetByID(ctx, decision.Id)
}

func (r *UnderwritingDecisionRepository) GetByID(ctx context.Context, decisionID string) (*underwritingv1.UnderwritingDecision, error) {
	var (
		d                       underwritingv1.UnderwritingDecision
		decisionStr             sql.NullString
		methodStr               sql.NullString
		riskScore               sql.NullString
		riskLevelStr            sql.NullString
		riskFactors             sql.NullString
		reason                  sql.NullString
		conditions              sql.NullString
		adjustedPremium         int64
		adjustedPremiumCurrency string
		adjustmentReason        sql.NullString
		underwriterID           sql.NullString
		underwriterComments     sql.NullString
		decidedAt               time.Time
		validUntil              sql.NullTime
		auditInfo               sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT decision_id, quote_id, decision, method, risk_score, risk_level,
		       risk_factors, reason, conditions, premium_adjusted, adjusted_premium,
		       adjusted_premium_currency, adjustment_reason, underwriter_id,
		       underwriter_comments, decided_at, valid_until, audit_info
		FROM insurance_schema.underwriting_decisions
		WHERE decision_id = $1`,
		decisionID,
	).Row().Scan(
		&d.Id,
		&d.QuoteId,
		&decisionStr,
		&methodStr,
		&riskScore,
		&riskLevelStr,
		&riskFactors,
		&reason,
		&conditions,
		&d.PremiumAdjusted,
		&adjustedPremium,
		&adjustedPremiumCurrency,
		&adjustmentReason,
		&underwriterID,
		&underwriterComments,
		&decidedAt,
		&validUntil,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get underwriting decision: %w", err)
	}

	// Set Money field
	d.AdjustedPremium = &commonv1.Money{Amount: adjustedPremium, Currency: adjustedPremiumCurrency}

	// Set optional fields
	if riskScore.Valid {
		d.RiskScore = riskScore.String
	}
	if riskFactors.Valid {
		d.RiskFactors = riskFactors.String
	}
	if reason.Valid {
		d.Reason = reason.String
	}
	if conditions.Valid {
		d.Conditions = conditions.String
	}
	if adjustmentReason.Valid {
		d.AdjustmentReason = adjustmentReason.String
	}
	if underwriterID.Valid {
		d.UnderwriterId = underwriterID.String
	}
	if underwriterComments.Valid {
		d.UnderwriterComments = underwriterComments.String
	}

	// Parse enums
	if decisionStr.Valid {
		k := strings.ToUpper(decisionStr.String)
		if v, ok := underwritingv1.DecisionType_value[k]; ok {
			d.Decision = underwritingv1.DecisionType(v)
		}
	}
	if methodStr.Valid {
		k := strings.ToUpper(methodStr.String)
		if v, ok := underwritingv1.DecisionMethod_value[k]; ok {
			d.Method = underwritingv1.DecisionMethod(v)
		}
	}
	if riskLevelStr.Valid {
		k := strings.ToUpper(riskLevelStr.String)
		if v, ok := underwritingv1.RiskLevel_value[k]; ok {
			d.RiskLevel = underwritingv1.RiskLevel(v)
		}
	}

	// Set timestamps
	if !decidedAt.IsZero() {
		d.DecidedAt = timestamppb.New(decidedAt)
	}
	if validUntil.Valid {
		d.ValidUntil = timestamppb.New(validUntil.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		d.AuditInfo = &commonv1.AuditInfo{}
	}

	return &d, nil
}

func (r *UnderwritingDecisionRepository) Update(ctx context.Context, decision *underwritingv1.UnderwritingDecision) (*underwritingv1.UnderwritingDecision, error) {
	// Extract Money values
	adjustedPremium := int64(0)
	adjustedPremiumCurrency := "BDT"
	if decision.AdjustedPremium != nil {
		adjustedPremium = decision.AdjustedPremium.Amount
		adjustedPremiumCurrency = decision.AdjustedPremium.Currency
	}

	// Handle JSONB fields
	var riskFactors interface{}
	if decision.RiskFactors != "" {
		riskFactors = decision.RiskFactors
	}

	var conditions interface{}
	if decision.Conditions != "" {
		conditions = decision.Conditions
	}

	// Handle timestamps
	var decidedAt time.Time
	if decision.DecidedAt != nil {
		decidedAt = decision.DecidedAt.AsTime()
	}

	var validUntil sql.NullTime
	if decision.ValidUntil != nil {
		validUntil = sql.NullTime{Time: decision.ValidUntil.AsTime(), Valid: true}
	}

	var underwriterID sql.NullString
	if decision.UnderwriterId != "" {
		underwriterID = sql.NullString{String: decision.UnderwriterId, Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.underwriting_decisions
		SET quote_id = $2,
		    decision = $3,
		    method = $4,
		    risk_score = $5,
		    risk_level = $6,
		    risk_factors = $7,
		    reason = $8,
		    conditions = $9,
		    premium_adjusted = $10,
		    adjusted_premium = $11,
		    adjusted_premium_currency = $12,
		    adjustment_reason = $13,
		    underwriter_id = $14,
		    underwriter_comments = $15,
		    decided_at = $16,
		    valid_until = $17
		WHERE decision_id = $1`,
		decision.Id,
		decision.QuoteId,
		strings.ToUpper(decision.Decision.String()),
		strings.ToUpper(decision.Method.String()),
		decision.RiskScore,
		strings.ToUpper(decision.RiskLevel.String()),
		riskFactors,
		decision.Reason,
		conditions,
		decision.PremiumAdjusted,
		adjustedPremium,
		adjustedPremiumCurrency,
		decision.AdjustmentReason,
		underwriterID,
		decision.UnderwriterComments,
		decidedAt,
		validUntil,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update underwriting decision: %w", err)
	}

	return r.GetByID(ctx, decision.Id)
}

func (r *UnderwritingDecisionRepository) Delete(ctx context.Context, decisionID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.underwriting_decisions
		WHERE decision_id = $1`,
		decisionID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete underwriting decision: %w", err)
	}

	return nil
}

func (r *UnderwritingDecisionRepository) ListByQuoteID(ctx context.Context, quoteID string) ([]*underwritingv1.UnderwritingDecision, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT decision_id, quote_id, decision, method, risk_score, risk_level,
		       risk_factors, reason, conditions, premium_adjusted, adjusted_premium,
		       adjusted_premium_currency, adjustment_reason, underwriter_id,
		       underwriter_comments, decided_at, valid_until, audit_info
		FROM insurance_schema.underwriting_decisions
		WHERE quote_id = $1
		ORDER BY decided_at DESC`,
		quoteID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list underwriting decisions: %w", err)
	}
	defer rows.Close()

	decisions := make([]*underwritingv1.UnderwritingDecision, 0)
	for rows.Next() {
		var (
			d                       underwritingv1.UnderwritingDecision
			decisionStr             sql.NullString
			methodStr               sql.NullString
			riskScore               sql.NullString
			riskLevelStr            sql.NullString
			riskFactors             sql.NullString
			reason                  sql.NullString
			conditions              sql.NullString
			adjustedPremium         int64
			adjustedPremiumCurrency string
			adjustmentReason        sql.NullString
			underwriterID           sql.NullString
			underwriterComments     sql.NullString
			decidedAt               time.Time
			validUntil              sql.NullTime
			auditInfo               sql.NullString
		)

		err := rows.Scan(
			&d.Id,
			&d.QuoteId,
			&decisionStr,
			&methodStr,
			&riskScore,
			&riskLevelStr,
			&riskFactors,
			&reason,
			&conditions,
			&d.PremiumAdjusted,
			&adjustedPremium,
			&adjustedPremiumCurrency,
			&adjustmentReason,
			&underwriterID,
			&underwriterComments,
			&decidedAt,
			&validUntil,
			&auditInfo,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan underwriting decision: %w", err)
		}

		// Set Money field
		d.AdjustedPremium = &commonv1.Money{Amount: adjustedPremium, Currency: adjustedPremiumCurrency}

		// Set optional fields
		if riskScore.Valid {
			d.RiskScore = riskScore.String
		}
		if riskFactors.Valid {
			d.RiskFactors = riskFactors.String
		}
		if reason.Valid {
			d.Reason = reason.String
		}
		if conditions.Valid {
			d.Conditions = conditions.String
		}
		if adjustmentReason.Valid {
			d.AdjustmentReason = adjustmentReason.String
		}
		if underwriterID.Valid {
			d.UnderwriterId = underwriterID.String
		}
		if underwriterComments.Valid {
			d.UnderwriterComments = underwriterComments.String
		}

		// Parse enums
		if decisionStr.Valid {
			k := strings.ToUpper(decisionStr.String)
			if v, ok := underwritingv1.DecisionType_value[k]; ok {
				d.Decision = underwritingv1.DecisionType(v)
			}
		}
		if methodStr.Valid {
			k := strings.ToUpper(methodStr.String)
			if v, ok := underwritingv1.DecisionMethod_value[k]; ok {
				d.Method = underwritingv1.DecisionMethod(v)
			}
		}
		if riskLevelStr.Valid {
			k := strings.ToUpper(riskLevelStr.String)
			if v, ok := underwritingv1.RiskLevel_value[k]; ok {
				d.RiskLevel = underwritingv1.RiskLevel(v)
			}
		}

		// Set timestamps
		if !decidedAt.IsZero() {
			d.DecidedAt = timestamppb.New(decidedAt)
		}
		if validUntil.Valid {
			d.ValidUntil = timestamppb.New(validUntil.Time)
		}

		// Set audit info
		if auditInfo.Valid {
			d.AuditInfo = &commonv1.AuditInfo{}
		}

		decisions = append(decisions, &d)
	}

	return decisions, nil
}
