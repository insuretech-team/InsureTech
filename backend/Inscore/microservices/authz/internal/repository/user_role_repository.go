package repository

import (
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
)

// userRoleRow is a plain Go struct for GORM operations on user_roles table.
// It avoids proto-generated field issues with GORM reflection.
type userRoleRow struct {
	UserRoleID string     `gorm:"primaryKey;column:user_role_id"`
	UserID     string     `gorm:"column:user_id"`
	RoleID     string     `gorm:"column:role_id"`
	Domain     string     `gorm:"column:domain"`
	AssignedBy *string    `gorm:"column:assigned_by"`
	AssignedAt *time.Time `gorm:"column:assigned_at"`
	ExpiresAt  *time.Time `gorm:"column:expires_at"`
}

func (userRoleRow) TableName() string {
	return "authz_schema.user_roles"
}

// userRoleRowToProto converts a userRoleRow to a proto UserRole, converting timestamps.
func userRoleRowToProto(row *userRoleRow) *entityv1.UserRole {
	ur := &entityv1.UserRole{
		UserRoleId: row.UserRoleID,
		UserId:     row.UserID,
		RoleId:     row.RoleID,
		Domain:     row.Domain,
	}
	if row.AssignedBy != nil {
		ur.AssignedBy = *row.AssignedBy
	}
	if row.AssignedAt != nil {
		ur.AssignedAt = timestamppb.New(*row.AssignedAt)
	}
	if row.ExpiresAt != nil {
		ur.ExpiresAt = timestamppb.New(*row.ExpiresAt)
	}
	return ur
}

// UserRoleRepo implements domain.UserRoleRepository using plain row structs for GORM.
type UserRoleRepo struct{ db *gorm.DB }

func NewUserRoleRepo(db *gorm.DB) *UserRoleRepo { return &UserRoleRepo{db: db} }

func (r *UserRoleRepo) Assign(ctx context.Context, ur *entityv1.UserRole) (*entityv1.UserRole, error) {
	var assignedBy *string
	var assignedAt *time.Time
	var expiresAt *time.Time
	
	if ur.AssignedBy != "" {
		assignedBy = &ur.AssignedBy
	}
	if ur.AssignedAt != nil {
		t := ur.AssignedAt.AsTime()
		assignedAt = &t
	}
	if ur.ExpiresAt != nil {
		t := ur.ExpiresAt.AsTime()
		expiresAt = &t
	}
	
	row := &userRoleRow{
		UserRoleID: ur.UserRoleId,
		UserID:     ur.UserId,
		RoleID:     ur.RoleId,
		Domain:     ur.Domain,
		AssignedBy: assignedBy,
		AssignedAt: assignedAt,
		ExpiresAt:  expiresAt,
	}
	
	res := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}, {Name: "domain"}},
		DoUpdates: clause.AssignmentColumns([]string{"assigned_by", "assigned_at", "expires_at"}),
	}).Create(row)
	if res.Error != nil {
		return nil, errors.New("userRole.Assign: " + res.Error.Error())
	}
	return ur, nil
}

func (r *UserRoleRepo) Remove(ctx context.Context, userID, roleID, domain string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ? AND domain = ?", userID, roleID, domain).
		Delete(&userRoleRow{}).Error
}

func (r *UserRoleRepo) ListByUser(ctx context.Context, userID, domain string) ([]*entityv1.UserRole, error) {
	var rows []*userRoleRow
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if domain != "" {
		q = q.Where("domain = ?", domain)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, errors.New("userRole.ListByUser: " + err.Error())
	}
	
	urs := make([]*entityv1.UserRole, len(rows))
	for i, row := range rows {
		urs[i] = userRoleRowToProto(row)
	}
	return urs, nil
}

func (r *UserRoleRepo) ListByRole(ctx context.Context, roleID, domain string) ([]*entityv1.UserRole, error) {
	var rows []*userRoleRow
	q := r.db.WithContext(ctx).Where("role_id = ?", roleID)
	if domain != "" {
		q = q.Where("domain = ?", domain)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, errors.New("userRole.ListByRole: " + err.Error())
	}
	
	urs := make([]*entityv1.UserRole, len(rows))
	for i, row := range rows {
		urs[i] = userRoleRowToProto(row)
	}
	return urs, nil
}

func (r *UserRoleRepo) Revoke(ctx context.Context, userID, roleID, domain string) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ? AND domain = ?", userID, roleID, domain).
		Delete(&userRoleRow{}).Error; err != nil {
		return errors.New("userRole.Revoke: " + err.Error())
	}
	return nil
}
