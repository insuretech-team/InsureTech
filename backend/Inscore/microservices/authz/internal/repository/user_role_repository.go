package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
)

// UserRoleRepo implements domain.UserRoleRepository using proto structs directly.
type UserRoleRepo struct{ db *gorm.DB }

func NewUserRoleRepo(db *gorm.DB) *UserRoleRepo { return &UserRoleRepo{db: db} }

func (r *UserRoleRepo) Assign(ctx context.Context, ur *entityv1.UserRole) (*entityv1.UserRole, error) {
	res := r.db.WithContext(ctx).Table("authz_schema.user_roles").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}, {Name: "domain"}},
		DoUpdates: clause.AssignmentColumns([]string{"assigned_by", "assigned_at", "expires_at"}),
	}).Create(ur)
	if res.Error != nil {
		return nil, errors.New("userRole.Assign: " + res.Error.Error())
	}
	return ur, nil
}

func (r *UserRoleRepo) Remove(ctx context.Context, userID, roleID, domain string) error {
	return r.db.WithContext(ctx).
		Table("authz_schema.user_roles").
		Where("user_id = ? AND role_id = ? AND domain = ?", userID, roleID, domain).
		Delete(&entityv1.UserRole{}).Error
}

func (r *UserRoleRepo) ListByUser(ctx context.Context, userID, domain string) ([]*entityv1.UserRole, error) {
	var urs []*entityv1.UserRole
	q := r.db.WithContext(ctx).Table("authz_schema.user_roles").Where("user_id = ?", userID)
	if domain != "" {
		q = q.Where("domain = ?", domain)
	}
	if err := q.Find(&urs).Error; err != nil {
		return nil, errors.New("userRole.ListByUser: " + err.Error())
	}
	return urs, nil
}

func (r *UserRoleRepo) ListByRole(ctx context.Context, roleID, domain string) ([]*entityv1.UserRole, error) {
	var urs []*entityv1.UserRole
	q := r.db.WithContext(ctx).Table("authz_schema.user_roles").Where("role_id = ?", roleID)
	if domain != "" {
		q = q.Where("domain = ?", domain)
	}
	if err := q.Find(&urs).Error; err != nil {
		return nil, errors.New("userRole.ListByRole: " + err.Error())
	}
	return urs, nil
}

func (r *UserRoleRepo) Revoke(ctx context.Context, userID, roleID, domain string) error {
	if err := r.db.WithContext(ctx).
		Table("authz_schema.user_roles").
		Where("user_id = ? AND role_id = ? AND domain = ?", userID, roleID, domain).
		Delete(&entityv1.UserRole{}).Error; err != nil {
		return errors.New("userRole.Revoke: " + err.Error())
	}
	return nil
}
