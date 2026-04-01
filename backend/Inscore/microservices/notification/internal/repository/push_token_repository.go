package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PushDeviceToken struct {
	TokenID     string
	UserID      string
	Provider    string
	Platform    string
	DeviceToken string
	DeviceID    string
	AppID       string
	IsActive    bool
	LastSeenAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PushTokenRepository struct {
	db *gorm.DB
}

func NewPushTokenRepository(db *gorm.DB) *PushTokenRepository {
	return &PushTokenRepository{db: db}
}

func (r *PushTokenRepository) Upsert(ctx context.Context, token *PushDeviceToken) error {
	if token == nil {
		return fmt.Errorf("push device token is required")
	}
	if strings.TrimSpace(token.TokenID) == "" {
		token.TokenID = uuid.NewString()
	}
	if strings.TrimSpace(token.UserID) == "" {
		return fmt.Errorf("user_id is required")
	}
	if strings.TrimSpace(token.Provider) == "" {
		return fmt.Errorf("provider is required")
	}
	if strings.TrimSpace(token.Platform) == "" {
		return fmt.Errorf("platform is required")
	}
	if strings.TrimSpace(token.DeviceToken) == "" {
		return fmt.Errorf("device_token is required")
	}

	return r.db.WithContext(ctx).Exec(`
		INSERT INTO notification_schema.push_device_tokens
			(token_id, user_id, provider, platform, device_token, device_id, app_id, is_active, last_seen_at)
		VALUES ($1, $2, UPPER($3), LOWER($4), $5, NULLIF($6, ''), NULLIF($7, ''), $8, COALESCE($9, NOW()))
		ON CONFLICT (provider, device_token)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			platform = EXCLUDED.platform,
			device_id = EXCLUDED.device_id,
			app_id = EXCLUDED.app_id,
			is_active = EXCLUDED.is_active,
			last_seen_at = COALESCE(EXCLUDED.last_seen_at, notification_schema.push_device_tokens.last_seen_at, NOW()),
			updated_at = NOW()`,
		token.TokenID,
		token.UserID,
		token.Provider,
		token.Platform,
		token.DeviceToken,
		token.DeviceID,
		token.AppID,
		token.IsActive,
		token.LastSeenAt,
	).Error
}

func (r *PushTokenRepository) ListActiveByUserID(ctx context.Context, userID string) ([]*PushDeviceToken, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT token_id, user_id, provider, platform, device_token, COALESCE(device_id, ''),
		       COALESCE(app_id, ''), is_active, last_seen_at, created_at, updated_at
		FROM notification_schema.push_device_tokens
		WHERE user_id = $1 AND is_active = TRUE
		ORDER BY COALESCE(last_seen_at, created_at) DESC, created_at DESC`,
		userID,
	).Rows()
	if err != nil {
		return nil, fmt.Errorf("list active push tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*PushDeviceToken
	for rows.Next() {
		token, scanErr := scanPushDeviceToken(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (r *PushTokenRepository) DeactivateByDeviceTokens(ctx context.Context, deviceTokens []string) error {
	filtered := make([]string, 0, len(deviceTokens))
	for _, token := range deviceTokens {
		if trimmed := strings.TrimSpace(token); trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}
	if len(filtered) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.push_device_tokens
		SET is_active = FALSE,
		    updated_at = NOW()
		WHERE device_token = ANY($1)`,
		pq.Array(filtered),
	).Error
}

func scanPushDeviceToken(scanner interface {
	Scan(dest ...any) error
}) (*PushDeviceToken, error) {
	var (
		token      PushDeviceToken
		lastSeenAt sql.NullTime
	)

	if err := scanner.Scan(
		&token.TokenID,
		&token.UserID,
		&token.Provider,
		&token.Platform,
		&token.DeviceToken,
		&token.DeviceID,
		&token.AppID,
		&token.IsActive,
		&lastSeenAt,
		&token.CreatedAt,
		&token.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan push device token: %w", err)
	}

	if lastSeenAt.Valid {
		t := lastSeenAt.Time
		token.LastSeenAt = &t
	}
	return &token, nil
}
