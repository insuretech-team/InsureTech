package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) Create(ctx context.Context, template *notificationv1.NotificationTemplate) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO notification_schema.notification_templates
			(template_id, template_name, type, channel, subject_template, body_template, language, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, $7, NOW(), NOW(), $8)`,
		template.GetTemplateId(),
		template.GetTemplateName(),
		dbNotificationType(template.GetType()),
		dbNotificationChannel(template.GetChannel()),
		template.GetSubjectTemplate(),
		template.GetBodyTemplate(),
		template.GetLanguage(),
		template.GetIsActive(),
	).Error
}

func (r *TemplateRepository) GetByID(ctx context.Context, templateID string) (*notificationv1.NotificationTemplate, error) {
	row := r.db.WithContext(ctx).Raw(`
		SELECT template_id, template_name, type, channel, COALESCE(subject_template, ''),
		       body_template, language, created_at, updated_at, is_active
		FROM notification_schema.notification_templates
		WHERE template_id = $1`, templateID).Row()
	return scanTemplate(row)
}

func (r *TemplateRepository) Update(ctx context.Context, templateID, name, subject, body string) error {
	updates := map[string]any{
		"updated_at": time.Now(),
	}
	if strings.TrimSpace(name) != "" {
		updates["template_name"] = name
	}
	if subject != "" {
		updates["subject_template"] = subject
	}
	if body != "" {
		updates["body_template"] = body
	}
	return r.db.WithContext(ctx).Table("notification_schema.notification_templates").Where("template_id = ?", templateID).Updates(updates).Error
}

func (r *TemplateRepository) Deactivate(ctx context.Context, templateID string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.notification_templates
		SET is_active = false,
		    updated_at = NOW()
		WHERE template_id = $1`, templateID).Error
}

func scanTemplate(scanner interface {
	Scan(dest ...any) error
}) (*notificationv1.NotificationTemplate, error) {
	var (
		template   notificationv1.NotificationTemplate
		typeStr    string
		channelStr string
		createdAt  time.Time
		updatedAt  time.Time
	)

	if err := scanner.Scan(
		&template.TemplateId,
		&template.TemplateName,
		&typeStr,
		&channelStr,
		&template.SubjectTemplate,
		&template.BodyTemplate,
		&template.Language,
		&createdAt,
		&updatedAt,
		&template.IsActive,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("scan notification template: %w", err)
	}

	template.Type = parseNotificationType(typeStr)
	template.Channel = parseNotificationChannel(channelStr)
	template.CreatedAt = timestamppb.New(createdAt)
	template.UpdatedAt = timestamppb.New(updatedAt)
	return &template, nil
}
