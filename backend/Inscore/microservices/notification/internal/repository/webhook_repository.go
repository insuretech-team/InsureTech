package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type WebhookSubscription struct {
	SubscriptionID string
	SubscriberName string
	TargetURL      string
	Secret         string
	EventTypes     []string
	TopicGroups    []string
	Topics         []string
	Channels       []string
	TimeoutSeconds int
	MaxAttempts    int
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type WebhookDeliveryAttempt struct {
	AttemptID       string
	SubscriptionID  string
	NotificationID  string
	LifecycleEvent  string
	SourceTopic     string
	Payload         json.RawMessage
	Status          string
	ResponseStatus  int
	ResponseBody    string
	RetryCount      int
	ErrorMessage    string
	ScheduledAt     *time.Time
	LastAttemptedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) CreateSubscription(ctx context.Context, subscription *WebhookSubscription) error {
	if subscription == nil {
		return fmt.Errorf("webhook subscription is required")
	}
	if strings.TrimSpace(subscription.SubscriptionID) == "" {
		subscription.SubscriptionID = uuid.NewString()
	}
	if strings.TrimSpace(subscription.SubscriberName) == "" {
		return fmt.Errorf("subscriber_name is required")
	}
	if strings.TrimSpace(subscription.TargetURL) == "" {
		return fmt.Errorf("target_url is required")
	}
	if strings.TrimSpace(subscription.Secret) == "" {
		return fmt.Errorf("secret is required")
	}
	if subscription.TimeoutSeconds <= 0 {
		subscription.TimeoutSeconds = 10
	}
	if subscription.MaxAttempts <= 0 {
		subscription.MaxAttempts = 5
	}

	return r.db.WithContext(ctx).Exec(`
		INSERT INTO notification_schema.webhook_subscriptions
			(subscription_id, subscriber_name, target_url, secret, event_types, topic_groups,
			 topics, channels, timeout_seconds, max_attempts, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		subscription.SubscriptionID,
		subscription.SubscriberName,
		subscription.TargetURL,
		subscription.Secret,
		pq.Array(normalizeFilterValues(subscription.EventTypes, true)),
		pq.Array(normalizeFilterValues(subscription.TopicGroups, false)),
		pq.Array(normalizeFilterValues(subscription.Topics, false)),
		pq.Array(normalizeFilterValues(subscription.Channels, true)),
		subscription.TimeoutSeconds,
		subscription.MaxAttempts,
		subscription.IsActive,
	).Error
}

func (r *WebhookRepository) ListMatchingSubscriptions(ctx context.Context, lifecycleEvent, sourceTopic, sourceGroup, channel string) ([]*WebhookSubscription, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT subscription_id, subscriber_name, target_url, secret, event_types, topic_groups,
		       topics, channels, timeout_seconds, max_attempts, is_active, created_at, updated_at
		FROM notification_schema.webhook_subscriptions
		WHERE is_active = TRUE
		  AND (COALESCE(array_length(event_types, 1), 0) = 0 OR UPPER($1) = ANY(event_types))
		  AND (COALESCE(array_length(topics, 1), 0) = 0 OR ($2 <> '' AND $2 = ANY(topics)))
		  AND (COALESCE(array_length(topic_groups, 1), 0) = 0 OR ($3 <> '' AND $3 = ANY(topic_groups)))
		  AND (COALESCE(array_length(channels, 1), 0) = 0 OR ($4 <> '' AND UPPER($4) = ANY(channels)))
		ORDER BY created_at ASC`,
		strings.ToUpper(strings.TrimSpace(lifecycleEvent)),
		strings.TrimSpace(sourceTopic),
		strings.TrimSpace(sourceGroup),
		strings.TrimSpace(channel),
	).Rows()
	if err != nil {
		return nil, fmt.Errorf("list matching webhook subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*WebhookSubscription
	for rows.Next() {
		subscription, scanErr := scanWebhookSubscription(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

func (r *WebhookRepository) GetSubscriptionByID(ctx context.Context, subscriptionID string) (*WebhookSubscription, error) {
	row := r.db.WithContext(ctx).Raw(`
		SELECT subscription_id, subscriber_name, target_url, secret, event_types, topic_groups,
		       topics, channels, timeout_seconds, max_attempts, is_active, created_at, updated_at
		FROM notification_schema.webhook_subscriptions
		WHERE subscription_id = $1`,
		subscriptionID,
	).Row()
	subscription, err := scanWebhookSubscription(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return subscription, nil
}

func (r *WebhookRepository) EnqueueAttempt(ctx context.Context, attempt *WebhookDeliveryAttempt) error {
	if attempt == nil {
		return fmt.Errorf("webhook delivery attempt is required")
	}
	if strings.TrimSpace(attempt.AttemptID) == "" {
		attempt.AttemptID = uuid.NewString()
	}
	if strings.TrimSpace(attempt.SubscriptionID) == "" {
		return fmt.Errorf("subscription_id is required")
	}
	if strings.TrimSpace(attempt.LifecycleEvent) == "" {
		return fmt.Errorf("lifecycle_event is required")
	}
	if len(attempt.Payload) == 0 {
		attempt.Payload = json.RawMessage(`{}`)
	}
	status := strings.ToUpper(strings.TrimSpace(attempt.Status))
	if status == "" {
		status = "QUEUED"
	}

	return r.db.WithContext(ctx).Exec(`
		INSERT INTO notification_schema.webhook_delivery_attempts
			(attempt_id, subscription_id, notification_id, lifecycle_event, source_topic,
			 payload, status, retry_count, scheduled_at)
		VALUES ($1, $2, NULLIF($3, '')::uuid, $4, NULLIF($5, ''), $6::jsonb, $7, $8, $9)
		ON CONFLICT (subscription_id, notification_id, lifecycle_event) DO NOTHING`,
		attempt.AttemptID,
		attempt.SubscriptionID,
		zeroStringToNil(attempt.NotificationID),
		strings.ToUpper(strings.TrimSpace(attempt.LifecycleEvent)),
		attempt.SourceTopic,
		string(attempt.Payload),
		status,
		attempt.RetryCount,
		attempt.ScheduledAt,
	).Error
}

func (r *WebhookRepository) ListDueAttempts(ctx context.Context, now time.Time, limit int) ([]*WebhookDeliveryAttempt, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT attempt_id, subscription_id, COALESCE(notification_id::text, ''), lifecycle_event,
		       COALESCE(source_topic, ''), payload::text, status, COALESCE(response_status, 0),
		       COALESCE(response_body, ''), retry_count, COALESCE(error_message, ''), scheduled_at,
		       last_attempted_at, created_at, updated_at
		FROM notification_schema.webhook_delivery_attempts
		WHERE status = 'QUEUED'
		  AND (scheduled_at IS NULL OR scheduled_at <= $1)
		ORDER BY created_at ASC
		LIMIT $2`,
		now,
		limit,
	).Rows()
	if err != nil {
		return nil, fmt.Errorf("list due webhook attempts: %w", err)
	}
	defer rows.Close()

	var attempts []*WebhookDeliveryAttempt
	for rows.Next() {
		attempt, scanErr := scanWebhookDeliveryAttempt(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		attempts = append(attempts, attempt)
	}
	return attempts, nil
}

func (r *WebhookRepository) MarkAttemptSending(ctx context.Context, attemptID string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.webhook_delivery_attempts
		SET status = 'SENDING',
		    response_status = NULL,
		    response_body = NULL,
		    error_message = NULL,
		    last_attempted_at = NOW()
		WHERE attempt_id = $1`,
		attemptID,
	).Error
}

func (r *WebhookRepository) MarkAttemptSucceeded(ctx context.Context, attemptID string, responseStatus int, responseBody string, at time.Time) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.webhook_delivery_attempts
		SET status = 'SENT',
		    response_status = NULLIF($2, 0),
		    response_body = NULLIF($3, ''),
		    error_message = NULL,
		    last_attempted_at = $4
		WHERE attempt_id = $1`,
		attemptID,
		responseStatus,
		responseBody,
		at,
	).Error
}

func (r *WebhookRepository) ScheduleAttemptRetry(ctx context.Context, attemptID string, retryCount int, nextAttempt time.Time, responseStatus int, responseBody, errorMessage string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.webhook_delivery_attempts
		SET status = 'QUEUED',
		    retry_count = $2,
		    scheduled_at = $3,
		    response_status = NULLIF($4, 0),
		    response_body = NULLIF($5, ''),
		    error_message = NULLIF($6, ''),
		    last_attempted_at = NOW()
		WHERE attempt_id = $1`,
		attemptID,
		retryCount,
		nextAttempt,
		responseStatus,
		responseBody,
		errorMessage,
	).Error
}

func (r *WebhookRepository) MarkAttemptFailed(ctx context.Context, attemptID string, retryCount int, responseStatus int, responseBody, errorMessage string, at time.Time) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE notification_schema.webhook_delivery_attempts
		SET status = 'FAILED',
		    retry_count = $2,
		    response_status = NULLIF($3, 0),
		    response_body = NULLIF($4, ''),
		    error_message = NULLIF($5, ''),
		    last_attempted_at = $6
		WHERE attempt_id = $1`,
		attemptID,
		retryCount,
		responseStatus,
		responseBody,
		errorMessage,
		at,
	).Error
}

func (r *WebhookRepository) GetAttemptByID(ctx context.Context, attemptID string) (*WebhookDeliveryAttempt, error) {
	row := r.db.WithContext(ctx).Raw(`
		SELECT attempt_id, subscription_id, COALESCE(notification_id::text, ''), lifecycle_event,
		       COALESCE(source_topic, ''), payload::text, status, COALESCE(response_status, 0),
		       COALESCE(response_body, ''), retry_count, COALESCE(error_message, ''), scheduled_at,
		       last_attempted_at, created_at, updated_at
		FROM notification_schema.webhook_delivery_attempts
		WHERE attempt_id = $1`,
		attemptID,
	).Row()
	return scanWebhookDeliveryAttempt(row)
}

func scanWebhookSubscription(scanner interface {
	Scan(dest ...any) error
}) (*WebhookSubscription, error) {
	var (
		subscription WebhookSubscription
		eventTypes   pq.StringArray
		topicGroups  pq.StringArray
		topics       pq.StringArray
		channels     pq.StringArray
	)

	if err := scanner.Scan(
		&subscription.SubscriptionID,
		&subscription.SubscriberName,
		&subscription.TargetURL,
		&subscription.Secret,
		&eventTypes,
		&topicGroups,
		&topics,
		&channels,
		&subscription.TimeoutSeconds,
		&subscription.MaxAttempts,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan webhook subscription: %w", err)
	}

	subscription.EventTypes = []string(eventTypes)
	subscription.TopicGroups = []string(topicGroups)
	subscription.Topics = []string(topics)
	subscription.Channels = []string(channels)
	return &subscription, nil
}

func scanWebhookDeliveryAttempt(scanner interface {
	Scan(dest ...any) error
}) (*WebhookDeliveryAttempt, error) {
	var (
		attempt         WebhookDeliveryAttempt
		payloadText     string
		scheduledAt     sql.NullTime
		lastAttemptedAt sql.NullTime
	)

	if err := scanner.Scan(
		&attempt.AttemptID,
		&attempt.SubscriptionID,
		&attempt.NotificationID,
		&attempt.LifecycleEvent,
		&attempt.SourceTopic,
		&payloadText,
		&attempt.Status,
		&attempt.ResponseStatus,
		&attempt.ResponseBody,
		&attempt.RetryCount,
		&attempt.ErrorMessage,
		&scheduledAt,
		&lastAttemptedAt,
		&attempt.CreatedAt,
		&attempt.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("scan webhook delivery attempt: %w", err)
	}

	attempt.Payload = json.RawMessage(payloadText)
	if scheduledAt.Valid {
		t := scheduledAt.Time
		attempt.ScheduledAt = &t
	}
	if lastAttemptedAt.Valid {
		t := lastAttemptedAt.Time
		attempt.LastAttemptedAt = &t
	}
	return &attempt, nil
}

func normalizeFilterValues(values []string, uppercase bool) []string {
	if len(values) == 0 {
		return []string{}
	}
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if uppercase {
			trimmed = strings.ToUpper(trimmed)
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}
	return normalized
}

func zeroStringToNil(value string) any {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return value
}
