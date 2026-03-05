package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
)

// CasbinRuleRepo implements domain.CasbinRuleRepository using proto structs directly.
type CasbinRuleRepo struct{ db *gorm.DB }

func NewCasbinRuleRepo(db *gorm.DB) *CasbinRuleRepo { return &CasbinRuleRepo{db: db} }

func (r *CasbinRuleRepo) Upsert(ctx context.Context, rule *entityv1.CasbinRule) (*entityv1.CasbinRule, error) {
	res := r.db.WithContext(ctx).Table("authz_schema.casbin_rules").
		Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ?", rule.Ptype, rule.V0, rule.V1, rule.V2, rule.V3).
		Updates(map[string]any{"v4": rule.V4, "v5": rule.V5})
	if res.Error != nil {
		return nil, errors.New("casbinRule.Upsert: " + res.Error.Error())
	}
	if res.RowsAffected > 0 {
		return rule, nil
	}

	res = r.db.WithContext(ctx).Table("authz_schema.casbin_rules").Create(rule)
	if res.Error != nil {
		return nil, errors.New("casbinRule.Upsert: " + res.Error.Error())
	}
	return rule, nil
}

func (r *CasbinRuleRepo) Delete(ctx context.Context, rule *entityv1.CasbinRule) error {
	if err := r.db.WithContext(ctx).
		Table("authz_schema.casbin_rules").
		Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ?",
			rule.Ptype, rule.V0, rule.V1, rule.V2, rule.V3).
		Delete(&entityv1.CasbinRule{}).Error; err != nil {
		return errors.New("casbinRule.Delete: " + err.Error())
	}
	return nil
}

func (r *CasbinRuleRepo) ListByDomain(ctx context.Context, domain string) ([]*entityv1.CasbinRule, error) {
	var rules []*entityv1.CasbinRule
	if err := r.db.WithContext(ctx).Table("authz_schema.casbin_rules").Where("v1 = ?", domain).Find(&rules).Error; err != nil {
		return nil, errors.New("casbinRule.ListByDomain: " + err.Error())
	}
	return rules, nil
}

func (r *CasbinRuleRepo) DeleteByDomainAndSubject(ctx context.Context, domain, subject string) error {
	if err := r.db.WithContext(ctx).
		Table("authz_schema.casbin_rules").
		Where("v1 = ? AND v0 = ?", domain, subject).
		Delete(&entityv1.CasbinRule{}).Error; err != nil {
		return errors.New("casbinRule.DeleteByDomainAndSubject: " + err.Error())
	}
	return nil
}
