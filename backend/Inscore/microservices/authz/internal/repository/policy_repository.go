package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
)

// PolicyRepo implements domain.PolicyRuleRepository using proto structs directly.
type PolicyRepo struct{ db *gorm.DB }

func NewPolicyRepo(db *gorm.DB) *PolicyRepo { return &PolicyRepo{db: db} }

func (r *PolicyRepo) Create(ctx context.Context, pr *entityv1.PolicyRule) (*entityv1.PolicyRule, error) {
	if pr == nil {
		return nil, errors.New("policy.Create: nil policy")
	}
	effect := strings.ToLower(strings.TrimPrefix(pr.Effect.String(), "POLICY_EFFECT_"))
	now := time.Now()
	values := map[string]any{
		"policy_id":   pr.PolicyId,
		"subject":     pr.Subject,
		"domain":      pr.Domain,
		"object":      pr.Object,
		"action":      pr.Action,
		"effect":      effect,
		"condition":   pr.Condition,
		"description": pr.Description,
		"is_active":   pr.IsActive,
		"created_by":  nullableUUID(pr.CreatedBy),
		"created_at":  now,
		"updated_at":  now,
	}
	if pr.PolicyId == "" {
		delete(values, "policy_id")
	}
	if err := r.db.WithContext(ctx).Table("authz_schema.policy_rules").Create(values).Error; err != nil {
		return nil, errors.New("policy.Create: " + err.Error())
	}
	return pr, nil
}

func (r *PolicyRepo) GetByID(ctx context.Context, id string) (*entityv1.PolicyRule, error) {
	var pr entityv1.PolicyRule
	if err := r.db.WithContext(ctx).Table("authz_schema.policy_rules").Where("policy_id = ?", id).First(&pr).Error; err != nil {
		return nil, errors.New("policy.GetByID: " + err.Error())
	}
	return &pr, nil
}

func (r *PolicyRepo) SoftDelete(ctx context.Context, policyID string) error {
	return r.db.WithContext(ctx).Table("authz_schema.policy_rules").
		Where("policy_id = ?", policyID).
		Update("is_active", false).Error
}

func (r *PolicyRepo) List(ctx context.Context, domain string, activeOnly bool, limit, offset int) ([]*entityv1.PolicyRule, error) {
	var prs []*entityv1.PolicyRule
	q := r.db.WithContext(ctx).Table("authz_schema.policy_rules")
	if domain != "" {
		q = q.Where("domain = ?", domain)
	}
	if activeOnly {
		q = q.Where("is_active = true")
	}
	if err := q.Limit(limit).Offset(offset).Find(&prs).Error; err != nil {
		return nil, errors.New("policy.List: " + err.Error())
	}
	return prs, nil
}

func (r *PolicyRepo) Update(ctx context.Context, pr *entityv1.PolicyRule) (*entityv1.PolicyRule, error) {
	effect := strings.ToLower(strings.TrimPrefix(pr.Effect.String(), "POLICY_EFFECT_"))
	if err := r.db.WithContext(ctx).Table("authz_schema.policy_rules").
		Where("policy_id = ?", pr.PolicyId).
		Updates(map[string]any{
			"subject":     pr.Subject,
			"domain":      pr.Domain,
			"object":      pr.Object,
			"action":      pr.Action,
			"effect":      effect,
			"condition":   pr.Condition,
			"description": pr.Description,
			"is_active":   pr.IsActive,
			"created_by":  nullableUUID(pr.CreatedBy),
			"updated_at":  gorm.Expr("NOW()"),
		}).Error; err != nil {
		return nil, errors.New("policy.Update: " + err.Error())
	}
	return pr, nil
}

func (r *PolicyRepo) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Table("authz_schema.policy_rules").Where("policy_id = ?", id).Delete(&entityv1.PolicyRule{}).Error; err != nil {
		return errors.New("policy.Delete: " + err.Error())
	}
	return nil
}
