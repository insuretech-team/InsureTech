package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
)

type FraudRuleRepository struct {
	db *gorm.DB
}

func NewFraudRuleRepository(db *gorm.DB) *FraudRuleRepository {
	return &FraudRuleRepository{db: db}
}

func (r *FraudRuleRepository) Create(ctx context.Context, rule *fraudv1.FraudRule) (*fraudv1.FraudRule, error) {
	if rule.FraudRuleId == "" {
		return nil, fmt.Errorf("fraud_rule_id is required")
	}

	// Get valid user UUID for audit_info
	var createdBy string
	err := r.db.WithContext(ctx).Raw(`SELECT user_id FROM authn_schema.users LIMIT 1`).Scan(&createdBy).Error
	if err != nil || createdBy == "" {
		return nil, fmt.Errorf("failed to get valid user for created_by: %w", err)
	}

	auditInfoJSON := fmt.Sprintf(`{"created_by":"%s","created_at":"%s"}`, createdBy, time.Now().UTC().Format(time.RFC3339))

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.fraud_rules
			(fraud_rule_id, name, category, description, conditions, risk_level, 
			 score_weight, is_active, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		rule.FraudRuleId,
		rule.Name,
		strings.ToUpper(rule.Category.String()),
		rule.Description,
		rule.Conditions,
		strings.ToUpper(rule.RiskLevel.String()),
		rule.ScoreWeight,
		rule.IsActive,
		auditInfoJSON,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert fraud rule: %w", err)
	}

	return r.GetByID(ctx, rule.FraudRuleId)
}

func (r *FraudRuleRepository) GetByID(ctx context.Context, ruleID string) (*fraudv1.FraudRule, error) {
	var (
		rule         fraudv1.FraudRule
		categoryStr  sql.NullString
		riskLevelStr sql.NullString
		conditions   sql.NullString
		auditInfo    sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT fraud_rule_id, name, category, 
		       COALESCE(description, '') as description,
		       conditions, risk_level, score_weight, 
		       COALESCE(is_active, false) as is_active,
		       audit_info
		FROM insurance_schema.fraud_rules
		WHERE fraud_rule_id = $1`,
		ruleID,
	).Row().Scan(
		&rule.FraudRuleId,
		&rule.Name,
		&categoryStr,
		&rule.Description,
		&conditions,
		&riskLevelStr,
		&rule.ScoreWeight,
		&rule.IsActive,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get fraud rule: %w", err)
	}

	// Parse category enum
	if categoryStr.Valid {
		k := strings.ToUpper(categoryStr.String)
		if v, ok := fraudv1.RuleCategory_value[k]; ok {
			rule.Category = fraudv1.RuleCategory(v)
		}
	}

	// Parse risk_level enum
	if riskLevelStr.Valid {
		k := strings.ToUpper(riskLevelStr.String)
		if v, ok := fraudv1.RiskLevel_value[k]; ok {
			rule.RiskLevel = fraudv1.RiskLevel(v)
		}
	}

	// Set conditions
	if conditions.Valid {
		rule.Conditions = conditions.String
	}

	return &rule, nil
}

func (r *FraudRuleRepository) Update(ctx context.Context, rule *fraudv1.FraudRule) (*fraudv1.FraudRule, error) {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.fraud_rules
		SET name = $2,
		    category = $3,
		    description = $4,
		    conditions = $5,
		    risk_level = $6,
		    score_weight = $7,
		    is_active = $8
		WHERE fraud_rule_id = $1`,
		rule.FraudRuleId,
		rule.Name,
		strings.ToUpper(rule.Category.String()),
		rule.Description,
		rule.Conditions,
		strings.ToUpper(rule.RiskLevel.String()),
		rule.ScoreWeight,
		rule.IsActive,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update fraud rule: %w", err)
	}

	return r.GetByID(ctx, rule.FraudRuleId)
}

func (r *FraudRuleRepository) Delete(ctx context.Context, ruleID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.fraud_rules
		WHERE fraud_rule_id = $1`,
		ruleID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete fraud rule: %w", err)
	}

	return nil
}

func (r *FraudRuleRepository) List(ctx context.Context, page, pageSize int) ([]*fraudv1.FraudRule, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	err := r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM insurance_schema.fraud_rules`).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fraud rules: %w", err)
	}

	// Get rules
	query := fmt.Sprintf(`
		SELECT fraud_rule_id, name, category, 
		       COALESCE(description, '') as description,
		       conditions, risk_level, score_weight, 
		       COALESCE(is_active, false) as is_active,
		       audit_info
		FROM insurance_schema.fraud_rules
		ORDER BY fraud_rule_id DESC LIMIT %d OFFSET %d`, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list fraud rules: %w", err)
	}
	defer rows.Close()

	rules := make([]*fraudv1.FraudRule, 0)
	for rows.Next() {
		var (
			rule         fraudv1.FraudRule
			categoryStr  sql.NullString
			riskLevelStr sql.NullString
			conditions   sql.NullString
			auditInfo    sql.NullString
		)

		err := rows.Scan(
			&rule.FraudRuleId,
			&rule.Name,
			&categoryStr,
			&rule.Description,
			&conditions,
			&riskLevelStr,
			&rule.ScoreWeight,
			&rule.IsActive,
			&auditInfo,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fraud rule: %w", err)
		}

		// Parse category enum
		if categoryStr.Valid {
			k := strings.ToUpper(categoryStr.String)
			if v, ok := fraudv1.RuleCategory_value[k]; ok {
				rule.Category = fraudv1.RuleCategory(v)
			}
		}

		// Parse risk_level enum
		if riskLevelStr.Valid {
			k := strings.ToUpper(riskLevelStr.String)
			if v, ok := fraudv1.RiskLevel_value[k]; ok {
				rule.RiskLevel = fraudv1.RiskLevel(v)
			}
		}

		// Set conditions
		if conditions.Valid {
			rule.Conditions = conditions.String
		}

		rules = append(rules, &rule)
	}

	return rules, total, nil
}

func (r *FraudRuleRepository) ListActive(ctx context.Context) ([]*fraudv1.FraudRule, error) {
	query := `
		SELECT fraud_rule_id, name, category, 
		       COALESCE(description, '') as description,
		       conditions, risk_level, score_weight, 
		       COALESCE(is_active, false) as is_active,
		       audit_info
		FROM insurance_schema.fraud_rules
		WHERE is_active = true
		ORDER BY score_weight DESC`

	rows, err := r.db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list active fraud rules: %w", err)
	}
	defer rows.Close()

	rules := make([]*fraudv1.FraudRule, 0)
	for rows.Next() {
		var (
			rule         fraudv1.FraudRule
			categoryStr  sql.NullString
			riskLevelStr sql.NullString
			conditions   sql.NullString
			auditInfo    sql.NullString
		)

		err := rows.Scan(
			&rule.FraudRuleId,
			&rule.Name,
			&categoryStr,
			&rule.Description,
			&conditions,
			&riskLevelStr,
			&rule.ScoreWeight,
			&rule.IsActive,
			&auditInfo,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan fraud rule: %w", err)
		}

		// Parse category enum
		if categoryStr.Valid {
			k := strings.ToUpper(categoryStr.String)
			if v, ok := fraudv1.RuleCategory_value[k]; ok {
				rule.Category = fraudv1.RuleCategory(v)
			}
		}

		// Parse risk_level enum
		if riskLevelStr.Valid {
			k := strings.ToUpper(riskLevelStr.String)
			if v, ok := fraudv1.RiskLevel_value[k]; ok {
				rule.RiskLevel = fraudv1.RiskLevel(v)
			}
		}

		// Set conditions
		if conditions.Valid {
			rule.Conditions = conditions.String
		}

		rules = append(rules, &rule)
	}

	return rules, nil
}
