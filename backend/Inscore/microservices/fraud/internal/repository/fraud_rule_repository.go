package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	"gorm.io/gorm"
)

var ErrRuleNotFound = errors.New("fraud rule not found")

// FraudRuleRepository handles CRUD for fraud rules.
type FraudRuleRepository struct {
	db *gorm.DB
}

func NewFraudRuleRepository(db *gorm.DB) *FraudRuleRepository {
	return &FraudRuleRepository{db: db}
}

func (r *FraudRuleRepository) Create(ctx context.Context, rule *fraudv1.FraudRule) error {
	if rule.FraudRuleId == "" {
		rule.FraudRuleId = uuid.NewString()
	}
	if rule.ScoreWeight <= 0 {
		rule.ScoreWeight = 10
	}
	if !rule.IsActive {
		rule.IsActive = true
	}


	now := time.Now().UTC()
	values := map[string]any{
		"fraud_rule_id": rule.FraudRuleId,
		"name":          rule.Name,
		"category":      rule.Category.String(),
		"description":   rule.Description,
		"conditions":    rule.Conditions,
		"risk_level":    rule.RiskLevel.String(),
		"score_weight":  rule.ScoreWeight,
		"is_active":     rule.IsActive,
		"created_at":    now,
		"updated_at":    now,
	}

	return r.db.WithContext(ctx).Table("insurance_schema.fraud_rules").Create(values).Error
}

func (r *FraudRuleRepository) GetByID(ctx context.Context, ruleID string) (*fraudv1.FraudRule, error) {
	var rule fraudv1.FraudRule
	err := r.db.WithContext(ctx).Table("insurance_schema.fraud_rules").
		Where("fraud_rule_id = ?", ruleID).
		First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRuleNotFound
		}
		return nil, err
	}
	return &rule, nil
}

func (r *FraudRuleRepository) Update(ctx context.Context, ruleID string, rule *fraudv1.FraudRule) error {
	updates := map[string]any{"updated_at": time.Now().UTC()}

	if rule.Name != "" {
		updates["name"] = rule.Name
	}
	if rule.Category != fraudv1.RuleCategory_RULE_CATEGORY_UNSPECIFIED {
		updates["category"] = rule.Category.String()
	}
	if rule.Description != "" {
		updates["description"] = rule.Description
	}
	if rule.Conditions != "" {
		updates["conditions"] = rule.Conditions
	}
	if rule.RiskLevel != fraudv1.RiskLevel_RISK_LEVEL_UNSPECIFIED {
		updates["risk_level"] = rule.RiskLevel.String()
	}
	if rule.ScoreWeight > 0 {
		updates["score_weight"] = rule.ScoreWeight
	}

	res := r.db.WithContext(ctx).Table("insurance_schema.fraud_rules").
		Where("fraud_rule_id = ?", ruleID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrRuleNotFound
	}
	return nil
}

func (r *FraudRuleRepository) List(ctx context.Context, category fraudv1.RuleCategory, activeOnly bool, limit, offset int) ([]*fraudv1.FraudRule, int32, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	q := r.db.WithContext(ctx).Table("insurance_schema.fraud_rules")
	if category != fraudv1.RuleCategory_RULE_CATEGORY_UNSPECIFIED {
		q = q.Where("category = ?", category.String())
	}
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rules []*fraudv1.FraudRule
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rules).Error; err != nil {
		return nil, 0, err
	}

	return rules, int32(total), nil
}

func (r *FraudRuleRepository) SetActive(ctx context.Context, ruleID string, active bool) error {
	res := r.db.WithContext(ctx).Table("insurance_schema.fraud_rules").
		Where("fraud_rule_id = ?", ruleID).
		Updates(map[string]any{
			"is_active":  active,
			"updated_at": time.Now().UTC(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrRuleNotFound
	}
	return nil
}
