package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// policyRuleRow is a plain Go struct used for GORM scanning.
// We cannot use the proto-generated PolicyRule directly because its
// timestamp fields are *timestamppb.Timestamp, which GORM cannot reflect.
type policyRuleRow struct {
	PolicyID    string     `gorm:"column:policy_id"`
	Subject     string     `gorm:"column:subject"`
	Domain      string     `gorm:"column:domain"`
	Object      string     `gorm:"column:object"`
	Action      string     `gorm:"column:action"`
	Effect      string     `gorm:"column:effect"`
	Condition   string     `gorm:"column:condition"`
	Description string     `gorm:"column:description"`
	IsActive    bool       `gorm:"column:is_active"`
	CreatedBy   *string    `gorm:"column:created_by"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (policyRuleRow) TableName() string { return "authz_schema.policy_rules" }

// rowToProto converts a scanned policyRuleRow into the proto PolicyRule.
func rowToProto(r *policyRuleRow) *entityv1.PolicyRule {
	pr := &entityv1.PolicyRule{
		PolicyId:    r.PolicyID,
		Subject:     r.Subject,
		Domain:      r.Domain,
		Object:      r.Object,
		Action:      r.Action,
		Condition:   r.Condition,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
	if r.CreatedBy != nil {
		pr.CreatedBy = *r.CreatedBy
	}
	if r.CreatedAt != nil {
		pr.CreatedAt = timestamppb.New(*r.CreatedAt)
	}
	if r.UpdatedAt != nil {
		pr.UpdatedAt = timestamppb.New(*r.UpdatedAt)
	}
	if r.DeletedAt != nil {
		pr.DeletedAt = timestamppb.New(*r.DeletedAt)
	}
	// Map effect string back to enum
	switch strings.ToLower(r.Effect) {
	case "allow":
		pr.Effect = entityv1.PolicyEffect_POLICY_EFFECT_ALLOW
	case "deny":
		pr.Effect = entityv1.PolicyEffect_POLICY_EFFECT_DENY
	}
	return pr
}

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
	var row policyRuleRow
	if err := r.db.WithContext(ctx).Where("policy_id = ?", id).First(&row).Error; err != nil {
		return nil, errors.New("policy.GetByID: " + err.Error())
	}
	return rowToProto(&row), nil
}

func (r *PolicyRepo) SoftDelete(ctx context.Context, policyID string) error {
	return r.db.WithContext(ctx).Table("authz_schema.policy_rules").
		Where("policy_id = ?", policyID).
		Update("is_active", false).Error
}

func (r *PolicyRepo) List(ctx context.Context, domain string, activeOnly bool, limit, offset int) ([]*entityv1.PolicyRule, error) {
	var rows []*policyRuleRow
	q := r.db.WithContext(ctx).Model(&policyRuleRow{})
	if domain != "" {
		q = q.Where("domain = ?", domain)
	}
	if activeOnly {
		q = q.Where("is_active = true")
	}
	if err := q.Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, errors.New("policy.List: " + err.Error())
	}
	prs := make([]*entityv1.PolicyRule, len(rows))
	for i, row := range rows {
		prs[i] = rowToProto(row)
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
	if err := r.db.WithContext(ctx).Where("policy_id = ?", id).Delete(&policyRuleRow{}).Error; err != nil {
		return errors.New("policy.Delete: " + err.Error())
	}
	return nil
}
