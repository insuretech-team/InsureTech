package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *notificationv1.Notification) error {
	templateData, err := marshalMap(notification.GetTemplateData())
	if err != nil {
		return fmt.Errorf("marshal template data: %w", err)
	}

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO notification_schema.notifications
			(notification_id, recipient_id, type, channel, subject, message, template_data,
			 priority, status, scheduled_at, sent_at, delivered_at, read_at, created_at,
			 retry_count, error_message)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, $7::jsonb, $8, $9, $10, $11, $12, $13, NOW(), $14, NULLIF($15, ''))`,
		notification.GetNotificationId(),
		notification.GetRecipientId(),
		dbNotificationType(notification.GetType()),
		dbNotificationChannel(notification.GetChannel()),
		notification.GetSubject(),
		notification.GetMessage(),
		string(templateData),
		dbNotificationPriority(notification.GetPriority()),
		dbNotificationStatus(notification.GetStatus()),
		tsToNullableTime(notification.GetScheduledAt()),
		tsToNullableTime(notification.GetSentAt()),
		tsToNullableTime(notification.GetDeliveredAt()),
		tsToNullableTime(notification.GetReadAt()),
		notification.GetRetryCount(),
		notification.GetErrorMessage(),
	).Error
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, notificationID string) (*notificationv1.Notification, error) {
	row := r.db.WithContext(ctx).Raw(`
		SELECT notification_id, recipient_id, type, channel, COALESCE(subject, ''), message,
		       COALESCE(template_data::text, '{}'), priority, status, scheduled_at, sent_at,
		       delivered_at, read_at, created_at, retry_count, COALESCE(error_message, '')
		FROM notification_schema.notifications
		WHERE notification_id = $1`, notificationID).Row()
	return scanNotification(row)
}

func (r *NotificationRepository) ListByRecipient(ctx context.Context, recipientID string, unreadOnly bool, limit, offset int32) ([]*notificationv1.Notification, int32, int32, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	where := "recipient_id = $1"
	args := []any{recipientID}
	if unreadOnly {
		where += " AND read_at IS NULL"
	}

	var totalCount int64
	if err := r.db.WithContext(ctx).Raw(
		`SELECT COUNT(*) FROM notification_schema.notifications WHERE `+where,
		args...,
	).Scan(&totalCount).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("count notifications: %w", err)
	}

	var unreadCount int64
	if err := r.db.WithContext(ctx).Raw(
		`SELECT COUNT(*) FROM notification_schema.notifications WHERE recipient_id = $1 AND read_at IS NULL`,
		recipientID,
	).Scan(&unreadCount).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("count unread notifications: %w", err)
	}

	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT notification_id, recipient_id, type, channel, COALESCE(subject, ''), message,
		       COALESCE(template_data::text, '{}'), priority, status, scheduled_at, sent_at,
		       delivered_at, read_at, created_at, retry_count, COALESCE(error_message, '')
		FROM notification_schema.notifications
		WHERE `+where+`
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		append(args, limit, offset)...,
	).Rows()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	notifications := make([]*notificationv1.Notification, 0)
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, 0, 0, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, int32(totalCount), int32(unreadCount), nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationIDs []string) error {
	if len(notificationIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.notifications
		SET read_at = NOW(),
		    status = CASE
		        WHEN status IN ('SENT', 'DELIVERED', 'NOTIFICATION_STATUS_SENT', 'NOTIFICATION_STATUS_DELIVERED') THEN 'READ'
		        ELSE status
		    END
		WHERE notification_id = ANY($1)`,
		pq.Array(notificationIDs),
	).Error
}

func (r *NotificationRepository) ListDue(ctx context.Context, now time.Time, limit int) ([]*notificationv1.Notification, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT notification_id, recipient_id, type, channel, COALESCE(subject, ''), message,
		       COALESCE(template_data::text, '{}'), priority, status, scheduled_at, sent_at,
		       delivered_at, read_at, created_at, retry_count, COALESCE(error_message, '')
		FROM notification_schema.notifications
		WHERE status IN ('QUEUED', 'NOTIFICATION_STATUS_QUEUED')
		  AND (scheduled_at IS NULL OR scheduled_at <= $1)
		ORDER BY priority DESC, created_at ASC
		LIMIT $2`,
		now, limit,
	).Rows()
	if err != nil {
		return nil, fmt.Errorf("list due notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*notificationv1.Notification
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (r *NotificationRepository) MarkSending(ctx context.Context, notificationID string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.notifications
		SET status = 'SENDING',
		    error_message = NULL
		WHERE notification_id = $1`, notificationID).Error
}

func (r *NotificationRepository) MarkSent(ctx context.Context, notificationID string, sentAt time.Time) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.notifications
		SET status = 'SENT',
		    sent_at = $2,
		    error_message = NULL
		WHERE notification_id = $1`,
		notificationID, sentAt,
	).Error
}

func (r *NotificationRepository) MarkDelivered(ctx context.Context, notificationID string, at time.Time) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.notifications
		SET status = 'DELIVERED',
		    sent_at = COALESCE(sent_at, $2),
		    delivered_at = $2,
		    error_message = NULL
		WHERE notification_id = $1`,
		notificationID, at,
	).Error
}

func (r *NotificationRepository) ScheduleRetry(ctx context.Context, notificationID string, retryCount int32, nextAttempt time.Time, errorMessage string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.notifications
		SET status = 'QUEUED',
		    retry_count = $2,
		    scheduled_at = $3,
		    error_message = NULLIF($4, '')
		WHERE notification_id = $1`,
		notificationID, retryCount, nextAttempt, errorMessage,
	).Error
}

func (r *NotificationRepository) MarkFailed(ctx context.Context, notificationID string, retryCount int32, errorMessage string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.notifications
		SET status = 'FAILED',
		    retry_count = $2,
		    error_message = NULLIF($3, '')
		WHERE notification_id = $1`,
		notificationID, retryCount, errorMessage,
	).Error
}

func scanNotification(scanner interface {
	Scan(dest ...any) error
}) (*notificationv1.Notification, error) {
	var (
		notification    notificationv1.Notification
		typeStr         string
		channelStr      string
		priorityStr     string
		statusStr       string
		templateDataRaw string
		scheduledAt     sql.NullTime
		sentAt          sql.NullTime
		deliveredAt     sql.NullTime
		readAt          sql.NullTime
		createdAt       time.Time
	)

	if err := scanner.Scan(
		&notification.NotificationId,
		&notification.RecipientId,
		&typeStr,
		&channelStr,
		&notification.Subject,
		&notification.Message,
		&templateDataRaw,
		&priorityStr,
		&statusStr,
		&scheduledAt,
		&sentAt,
		&deliveredAt,
		&readAt,
		&createdAt,
		&notification.RetryCount,
		&notification.ErrorMessage,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("scan notification: %w", err)
	}

	notification.Type = parseNotificationType(typeStr)
	notification.Channel = parseNotificationChannel(channelStr)
	notification.Priority = parseNotificationPriority(priorityStr)
	notification.Status = parseNotificationStatus(statusStr)
	notification.CreatedAt = timestamppb.New(createdAt)
	notification.ScheduledAt = nullTimeToTimestamp(scheduledAt)
	notification.SentAt = nullTimeToTimestamp(sentAt)
	notification.DeliveredAt = nullTimeToTimestamp(deliveredAt)
	notification.ReadAt = nullTimeToTimestamp(readAt)

	if strings.TrimSpace(templateDataRaw) == "" {
		templateDataRaw = "{}"
	}
	if err := json.Unmarshal([]byte(templateDataRaw), &notification.TemplateData); err != nil {
		return nil, fmt.Errorf("decode notification template data: %w", err)
	}

	return &notification, nil
}

func marshalMap(values map[string]string) ([]byte, error) {
	if len(values) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(values)
}

func tsToNullableTime(ts *timestamppb.Timestamp) any {
	if ts == nil {
		return nil
	}
	return ts.AsTime()
}

func nullTimeToTimestamp(value sql.NullTime) *timestamppb.Timestamp {
	if !value.Valid {
		return nil
	}
	return timestamppb.New(value.Time)
}

func parseNotificationType(value string) notificationv1.NotificationType {
	if parsed, ok := notificationv1.NotificationType_value[normalizeNotificationEnum(value, "NOTIFICATION_TYPE_")]; ok {
		return notificationv1.NotificationType(parsed)
	}
	return notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
}

func parseNotificationChannel(value string) notificationv1.NotificationChannel {
	if parsed, ok := notificationv1.NotificationChannel_value[normalizeNotificationEnum(value, "NOTIFICATION_CHANNEL_")]; ok {
		return notificationv1.NotificationChannel(parsed)
	}
	return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_UNSPECIFIED
}

func parseNotificationPriority(value string) notificationv1.NotificationPriority {
	if parsed, ok := notificationv1.NotificationPriority_value[normalizeNotificationEnum(value, "NOTIFICATION_PRIORITY_")]; ok {
		return notificationv1.NotificationPriority(parsed)
	}
	return notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_UNSPECIFIED
}

func parseNotificationStatus(value string) notificationv1.NotificationStatus {
	if parsed, ok := notificationv1.NotificationStatus_value[normalizeNotificationEnum(value, "NOTIFICATION_STATUS_")]; ok {
		return notificationv1.NotificationStatus(parsed)
	}
	return notificationv1.NotificationStatus_NOTIFICATION_STATUS_UNSPECIFIED
}

func normalizeNotificationEnum(value, prefix string) string {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if normalized == "" {
		return normalized
	}
	if strings.HasPrefix(normalized, prefix) {
		return normalized
	}
	return prefix + normalized
}

func dbNotificationType(value notificationv1.NotificationType) string {
	return strings.TrimPrefix(value.String(), "NOTIFICATION_TYPE_")
}

func dbNotificationChannel(value notificationv1.NotificationChannel) string {
	return strings.TrimPrefix(value.String(), "NOTIFICATION_CHANNEL_")
}

func dbNotificationPriority(value notificationv1.NotificationPriority) string {
	return strings.TrimPrefix(value.String(), "NOTIFICATION_PRIORITY_")
}

func dbNotificationStatus(value notificationv1.NotificationStatus) string {
	return strings.TrimPrefix(value.String(), "NOTIFICATION_STATUS_")
}
