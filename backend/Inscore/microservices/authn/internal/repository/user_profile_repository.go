package repository

import (
	"context"
	"time"

	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// UserProfileRepository provides access to authn_schema.user_profiles.
// NOTE: proto UserProfile uses user_id as primary key.
type UserProfileRepository struct{ db *gorm.DB }

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (r *UserProfileRepository) Create(ctx context.Context, p *authnv1.UserProfile) error {
	// Use proto-generated model directly.
	if p.CreatedAt == nil {
		p.CreatedAt = timestamppb.Now()
	}
	p.UpdatedAt = timestamppb.Now()
	return r.db.WithContext(ctx).Table("authn_schema.user_profiles").Create(p).Error
}

func (r *UserProfileRepository) GetByUserID(ctx context.Context, userID string) (*authnv1.UserProfile, error) {
	var p authnv1.UserProfile
	err := r.db.WithContext(ctx).
		Table("authn_schema.user_profiles").
		Where("user_id = ?", userID).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *UserProfileRepository) SetKYCVerified(ctx context.Context, userID string, verified bool, at *time.Time) error {
	upd := map[string]any{"kyc_verified": verified, "updated_at": time.Now()}
	if at != nil {
		upd["kyc_verified_at"] = *at
	}
	return r.db.WithContext(ctx).
		Table("authn_schema.user_profiles").
		Where("user_id = ?", userID).
		Updates(upd).Error
}

func (r *UserProfileRepository) Update(ctx context.Context, p *authnv1.UserProfile) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.user_profiles").
		Where("user_id = ?", p.UserId).
		Updates(p).Error
}

func (r *UserProfileRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Table("authn_schema.user_profiles").Where("user_id = ?", userID).Delete(map[string]any{}).Error
}

// UserProfileRepoIface is the interface satisfied by UserProfileRepository.
type UserProfileRepoIface interface {
	Create(ctx context.Context, p *authnv1.UserProfile) error
	GetByUserID(ctx context.Context, userID string) (*authnv1.UserProfile, error)
	Update(ctx context.Context, p *authnv1.UserProfile) error
	SetKYCVerified(ctx context.Context, userID string, verified bool, at *time.Time) error
	DeleteByUserID(ctx context.Context, userID string) error
}
