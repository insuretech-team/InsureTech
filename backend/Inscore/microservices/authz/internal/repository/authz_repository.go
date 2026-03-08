package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
)

// casbinRuleRow is a plain Go struct for GORM operations on casbin_rules table.
// It avoids proto-generated field issues with GORM reflection.
type casbinRuleRow struct {
	ID    int64  `gorm:"primaryKey;autoIncrement"`
	Ptype string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

func (casbinRuleRow) TableName() string {
	return "authz_schema.casbin_rules"
}

// CasbinRuleRepo implements domain.CasbinRuleRepository using plain row structs for GORM.
type CasbinRuleRepo struct{ db *gorm.DB }

func NewCasbinRuleRepo(db *gorm.DB) *CasbinRuleRepo { return &CasbinRuleRepo{db: db} }

func (r *CasbinRuleRepo) Upsert(ctx context.Context, rule *entityv1.CasbinRule) (*entityv1.CasbinRule, error) {
	res := r.db.WithContext(ctx).
		Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ?", rule.Ptype, rule.V0, rule.V1, rule.V2, rule.V3).
		Updates(map[string]any{"v4": rule.V4, "v5": rule.V5}).Model(&casbinRuleRow{})
	if res.Error != nil {
		return nil, errors.New("casbinRule.Upsert: " + res.Error.Error())
	}
	if res.RowsAffected > 0 {
		return rule, nil
	}

	row := &casbinRuleRow{
		Ptype: rule.Ptype,
		V0:    rule.V0,
		V1:    rule.V1,
		V2:    rule.V2,
		V3:    rule.V3,
		V4:    rule.V4,
		V5:    rule.V5,
	}
	res = r.db.WithContext(ctx).Create(row)
	if res.Error != nil {
		return nil, errors.New("casbinRule.Upsert: " + res.Error.Error())
	}
	return rule, nil
}

func (r *CasbinRuleRepo) Delete(ctx context.Context, rule *entityv1.CasbinRule) error {
	if err := r.db.WithContext(ctx).
		Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ?",
			rule.Ptype, rule.V0, rule.V1, rule.V2, rule.V3).
		Delete(&casbinRuleRow{}).Error; err != nil {
		return errors.New("casbinRule.Delete: " + err.Error())
	}
	return nil
}

func (r *CasbinRuleRepo) ListByDomain(ctx context.Context, domain string) ([]*entityv1.CasbinRule, error) {
	var rows []*casbinRuleRow
	if err := r.db.WithContext(ctx).Where("v1 = ?", domain).Find(&rows).Error; err != nil {
		return nil, errors.New("casbinRule.ListByDomain: " + err.Error())
	}
	
	rules := make([]*entityv1.CasbinRule, len(rows))
	for i, row := range rows {
		rules[i] = &entityv1.CasbinRule{
			Id:    row.ID,
			Ptype: row.Ptype,
			V0:    row.V0,
			V1:    row.V1,
			V2:    row.V2,
			V3:    row.V3,
			V4:    row.V4,
			V5:    row.V5,
		}
	}
	return rules, nil
}

func (r *CasbinRuleRepo) DeleteByDomainAndSubject(ctx context.Context, domain, subject string) error {
	if err := r.db.WithContext(ctx).
		Where("v1 = ? AND v0 = ?", domain, subject).
		Delete(&casbinRuleRow{}).Error; err != nil {
		return errors.New("casbinRule.DeleteByDomainAndSubject: " + err.Error())
	}
	return nil
}
