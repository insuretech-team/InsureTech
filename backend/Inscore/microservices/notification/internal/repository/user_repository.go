package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/notifyprefs"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*authnentityv1.User, error) {
	row := r.db.WithContext(ctx).Raw(`
		SELECT user_id, COALESCE(mobile_number, ''), COALESCE(email, ''),
		       COALESCE(notification_preference, ''), COALESCE(preferred_language, '')
		FROM authn_schema.users
		WHERE user_id = $1 AND deleted_at IS NULL`, userID).Row()

	var user authnentityv1.User
	if err := row.Scan(
		&user.UserId,
		&user.MobileNumber,
		&user.Email,
		&user.NotificationPreference,
		&user.PreferredLanguage,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetPreferences(ctx context.Context, userID string) (*notificationv1.NotificationPreference, error) {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	rawPreference := strings.TrimSpace(user.GetNotificationPreference())
	return notifyprefs.Parse(userID, rawPreference)
}
